/*                          _       _
 *__      _____  __ ___   ___  __ _| |_ ___
 *\ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
 * \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
 *  \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
 *
 * Copyright © 2016 - 2018 Weaviate. All rights reserved.
 * LICENSE: https://github.com/creativesoftwarefdn/weaviate/blob/develop/LICENSE.md
 * AUTHOR: Bob van Luijt (bob@kub.design)
 * See www.creativesoftwarefdn.org for details
 * Contact: @CreativeSofwFdn / bob@kub.design
 */

/*
 * THIS IS A DEMO CONNECTOR!
 * USE IT TO LEARN HOW TO CREATE YOUR OWN CONNECTOR.
 */

/*
When starting Weaviate, functions are called in the following order;
(find the function in this document to understand what it is that they do)
 - GetName
 - SetConfig
 - SetSchema
 - SetMessaging
 - SetServerAddress
 - Connect
 - Init

All other function are called on the API request

After creating the connector, make sure to add the name of the connector to: func GetAllConnectors() in configure_weaviate.go

*/

package tinkerpop

import (
	"bytes"
	"context"
	errors_ "errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/qasaur/gremgo"

	"github.com/go-openapi/strfmt"
	"github.com/mitchellh/mapstructure"

	"github.com/creativesoftwarefdn/weaviate/config"
	"github.com/creativesoftwarefdn/weaviate/connectors/utils"
	"github.com/creativesoftwarefdn/weaviate/messages"
	"github.com/creativesoftwarefdn/weaviate/models"
	"github.com/creativesoftwarefdn/weaviate/schema"
)

// Tinkerpop has some basic variables.
// This is mandatory, only change it if you need aditional, global variables
type Tinkerpop struct {
	client *gremgo.Client
	kind   string

	config        Config
	serverAddress string
	schema        *schema.WeaviateSchema
	messaging     *messages.Messaging
}

// Config represents the config outline for Foobar. The Database config shoud be of the following form:
// "database_config" : {
//     "host": "127.0.0.1",
//     "port": 9080
// }
// Notice that the port is the GRPC-port.
type Config struct {
	Host string
	Port int
}

// thingToVertex translates a thing struct into a vertex string
func thingOrActionToVertex(thing *models.Thing, typeOfObject string) (string, error) {

	// Create message service
	if typeOfObject != "thing" && typeOfObject != "action" {
		return "", errors_.New("type of Object should be 'thing' or 'action'")
	}

	// start vertex string Buffer
	var vertex bytes.Buffer

	// add vertext class and ID
	vertex.WriteString(`g.addV("` + thing.AtClass + `").property("uuid", uuid).property("type", "` + typeOfObject + `")`)

	// set the meta values (@class will be the vector)
	vertex.WriteString(`.property("context", "` + thing.AtContext + `")`)
	vertex.WriteString(`.property("creationTimeUnix", "` + strconv.FormatInt(thing.CreationTimeUnix, 10) + `")`)
	vertex.WriteString(`.property("lastUpdateTimeUnix", "` + strconv.FormatInt(thing.LastUpdateTimeUnix, 10) + `")`)

	// reflect in Schema
	schema := reflect.ValueOf(thing.Schema)

	// fetch the schema.properties from the schema
	if schema.Kind() == reflect.Map {
		for _, e := range schema.MapKeys() {
			v := schema.MapIndex(e)
			switch t := v.Interface().(type) {
			case string:
				vertex.WriteString(`.property("schema__` + e.String() + `", "` + t + `")`)
			case int:
				vertex.WriteString(`.property("schema__` + e.String() + `", ` + strconv.Itoa(t) + `)`)
			case int8:
				vertex.WriteString(`.property("schema__` + e.String() + `", ` + strconv.Itoa(int(t)) + `)`)
			case int16:
				vertex.WriteString(`.property("schema__` + e.String() + `", ` + strconv.FormatInt(int64(t), 10) + `)`)
			case int32:
				vertex.WriteString(`.property("schema__` + e.String() + `", ` + strconv.FormatInt(int64(t), 10) + `)`)
			case int64:
				vertex.WriteString(`.property("schema__` + e.String() + `", ` + strconv.FormatInt(t, 10) + `)`)
			case bool:
				vertex.WriteString(`.property("schema__` + e.String() + `", ` + strconv.FormatBool(t) + `)`)
			case float32:
				vertex.WriteString(`.property("schema__` + e.String() + `", ` + strconv.FormatFloat(float64(t), 'g', -1, 32) + `)`)
			case float64:
				vertex.WriteString(`.property("schema__` + e.String() + `", ` + strconv.FormatFloat(t, 'g', -1, 64) + `)`)
			default:
				fmt.Println("SOMETHING WEIRD")
				//messaging.ExitError(78, "The type "+reflect.ValueOf(t).Kind().String()+" is not found.")
			}
		}
	}

	return vertex.String(), nil
}

// createCrefObject is a helper function to create a cref-object. This function is used for Cassandra only.
func (f *Tinkerpop) createCrefObject(UUID strfmt.UUID, location string, refType connutils.RefType) *models.SingleRef {
	// Create the 'cref'-node for the response.
	crefObj := models.SingleRef{}

	// Get the given node properties to generate response object
	crefObj.NrDollarCref = UUID
	crefObj.Type = string(refType)
	url := location
	crefObj.LocationURL = &url

	return &crefObj
}

// GetName returns a unique connector name, this name is used to define the connector in the weaviate config
func (f *Tinkerpop) GetName() string {
	return "tinkerpop"
}

// SetConfig sets variables, which can be placed in the config file section "database_config: {}"
// can be custom for any connector, in the example below there is only host and port available.
//
// Important to bear in mind;
// 1. You need to add these to the struct Config in this document.
// 2. They will become available via f.config.[variable-name]
//
// 	"database": {
// 		"name": "foobar",
// 		"database_config" : {
// 			"host": "127.0.0.1",
// 			"port": 9080
// 		}
// 	},
func (f *Tinkerpop) SetConfig(configInput *config.Environment) error {

	// Mandatory: needed to add the JSON config represented as a map in f.config
	err := mapstructure.Decode(configInput.Database.DatabaseConfig, &f.config)

	// Example to: Validate if the essential  config is available, like host and port.
	if err != nil || len(f.config.Host) == 0 || f.config.Port == 0 {
		return errors_.New("could not get tinkerpop host/port from config")
	}

	// If success return nil, otherwise return the error (see above)
	return nil
}

// SetSchema takes actionSchema and thingsSchema as an input and makes them available globally at f.schema
// In case you want to modify the schema, this is the place to do so.
// Note: When this function is called, the schemas (action + things) are already validated, so you don't have to build the validation.
func (f *Tinkerpop) SetSchema(schemaInput *schema.WeaviateSchema) error {
	f.schema = schemaInput

	// If success return nil, otherwise return the error
	return nil
}

// SetMessaging is used to send messages to the service.
// Available message types are: f.messaging.Infomessage ...DebugMessage ...ErrorMessage ...ExitError (also exits the service) ...InfoMessage
func (f *Tinkerpop) SetMessaging(m *messages.Messaging) error {

	// mandatory, adds the message functions to f.messaging to make them globally accessible.
	f.messaging = m

	// If success return nil, otherwise return the error
	return nil
}

// SetServerAddress is used to fill in a global variable with the server address, but can also be used
// to do some custom actions.
// Does not return anything
func (f *Tinkerpop) SetServerAddress(addr string) {
	f.serverAddress = addr
}

// Connect creates a connection to the database and tables if not already available.
// The connections could not be closed because it is used more often.
func (f *Tinkerpop) Connect() error {

	messaging := &messages.Messaging{}

	// listen for errors
	errs := make(chan error)
	go func(chan error) {
		err := <-errs
		messaging.ExitError(78, err)
	}(errs) // Example of connection error handling logic

	// dial the websocket
	dialer := gremgo.NewDialer("ws://" + f.config.Host + ":" + strconv.Itoa(f.config.Port)) // Returns a WebSocket dialer to connect to Gremlin Server
	client, err := gremgo.Dial(dialer, errs)                                                // Returns a gremgo client to interact with
	if err != nil {
		return err
	}

	// return the client
	f.client = &client

	// If success return nil, otherwise return the error (also see above)
	return nil
}

// Init 1st initializes the schema in the database and 2nd creates a root key.
func (f *Tinkerpop) Init() error {

	//
	// tinkerpop does not need an initiation process
	//

	return nil
}

// Attach can attach something to the request-context
func (f *Tinkerpop) Attach(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

// AddThing adds a thing to the Foobar database with the given UUID.
// Takes the thing and a UUID as input.
// Thing is already validated against the ontology
func (f *Tinkerpop) AddThing(ctx context.Context, thing *models.Thing, UUID strfmt.UUID) error {

	// convert the thing to a Vertex
	vertex, err := thingOrActionToVertex(thing, "thing")

	// on error fail
	if err != nil {
		return err
	}

	// execute the query with uuid as binding, result is not used because we send out "202 Accepted" and assume a succes because validation takes place before this function runs
	_, err = f.client.Execute(
		vertex,
		map[string]string{"uuid": UUID.String()},
		map[string]string{},
	)

	// on error, fail
	if err != nil {
		return err
	}

	// If success return nil, otherwise return the error
	return nil
}

// GetThing fills the given ThingGetResponse with the values from the database, based on the given UUID.
func (f *Tinkerpop) GetThing(ctx context.Context, UUID strfmt.UUID, thingResponse *models.ThingGetResponse) error {

	// define the ID vertex and the UUID to fetch
	result, err := f.client.Execute(
		`g.V().has("uuid", uuid).has("type", objectType)`,
		map[string]string{"uuid": UUID.String(), "objectType": "thing"},
		map[string]string{},
	)

	// in case of error, return the error
	if err != nil {
		return err
	}

	// if there is no result, send not found by sending a nil
	if result.([]interface{})[0] == nil {
		return nil
	}

	// set UUID
	thingResponse.ThingID = UUID

	// This is a temporary key [FIX THIS]
	thingResponse.Key = f.createCrefObject("29ece1d8-c433-4757-b258-0b278478e17a", f.serverAddress, connutils.RefTypeKey)

	// Create the schema Map, this map will contain all the results
	responseSchema := make(map[string]interface{})

	// Loop over the Tinkerpop results
	for key, value := range result.([]interface{})[0].([]interface{})[0].(map[string]interface{}) {

		// set class name (called Vertex Label)
		if key == "label" {
			thingResponse.AtClass = value.(string)
		}

		// add schema properties
		if key == "properties" {
			for propKey, propValue := range value.(map[string]interface{}) {

				// set class name (called Vertex Label)
				if propKey == "context" {
					for _, propValueValue := range propValue.([]interface{}) {
						for propValueKeySingle, propValueValueSingle := range propValueValue.(map[string]interface{}) {
							if propValueKeySingle == "value" {
								thingResponse.AtContext = propValueValueSingle.(string)
							}
						}
					}
				}

				// set class name (called Vertex Label)
				if propKey == "creationTimeUnix" {
					for _, propValueValue := range propValue.([]interface{}) {
						for propValueKeySingle, propValueValueSingle := range propValueValue.(map[string]interface{}) {
							if propValueKeySingle == "value" {
								convertToInt64, err := strconv.ParseInt(propValueValueSingle.(string), 10, 64)
								if err != nil {
									return err
								}
								thingResponse.CreationTimeUnix = convertToInt64
							}
						}
					}
				}

				// set class name (called Vertex Label)
				if propKey == "lastUpdateTimeUnix" {
					for _, propValueValue := range propValue.([]interface{}) {
						for propValueKeySingle, propValueValueSingle := range propValueValue.(map[string]interface{}) {
							if propValueKeySingle == "value" {
								convertToInt64, err := strconv.ParseInt(propValueValueSingle.(string), 10, 64)
								if err != nil {
									return err
								}
								thingResponse.LastUpdateTimeUnix = convertToInt64
							}
						}
					}
				}

				// check if the key starts with schema__ prefix
				if strings.HasPrefix(propKey, "schema__") {
					// Grab the value and valueType of the properties
					for _, propValueSingle := range propValue.([]interface{}) {
						// loop over the id's and the values, add the interface{} type to the response schema
						for propValueSingleKey, propValueSingleValue := range propValueSingle.(map[string]interface{}) {
							if propValueSingleKey == "value" {
								responseSchema[propKey[8:]] = propValueSingleValue
							}
						}
					}
				}
			}
		}

		// schema = responseSchema
		thingResponse.Schema = responseSchema

	}

	// If success return nil, otherwise return the error
	return nil
}

// GetThings fills the given ThingsListResponse with the values from the database, based on the given UUIDs.
func (f *Tinkerpop) GetThings(ctx context.Context, UUIDs []strfmt.UUID, thingResponse *models.ThingsListResponse) error {
	f.messaging.DebugMessage(fmt.Sprintf("GetThings: %s", UUIDs))

	// If success return nil, otherwise return the error
	return nil
}

// ListThings fills the given ThingsListResponse with the values from the database, based on the given parameters.
func (f *Tinkerpop) ListThings(ctx context.Context, first int, offset int, keyID strfmt.UUID, wheres []*connutils.WhereQuery, thingsResponse *models.ThingsListResponse) error {

	// thingsResponse should be populated with the response that comes from the DB.
	// thingsResponse = based on the ontology

	// If success return nil, otherwise return the error
	return nil
}

// UpdateThing updates the Thing in the DB at the given UUID.
func (f *Tinkerpop) UpdateThing(ctx context.Context, thing *models.Thing, UUID strfmt.UUID) error {

	// Run the query to update the thing based on its UUID.

	// If success return nil, otherwise return the error
	return nil
}

// DeleteThing deletes the Thing in the DB at the given UUID.
func (f *Tinkerpop) DeleteThing(ctx context.Context, thing *models.Thing, UUID strfmt.UUID) error {

	// Remove based on type and uuid
	_, err := f.client.Execute(
		`g.V().has("uuid", uuid).has("type", objectType).drop()`,
		map[string]string{"uuid": UUID.String(), "objectType": "thing"},
		map[string]string{},
	)

	// return error
	if err != nil {
		return err
	}

	// If success return nil, otherwise return the error
	return nil
}

// HistoryThing fills the history of a thing based on its UUID
func (f *Tinkerpop) HistoryThing(ctx context.Context, UUID strfmt.UUID, history *models.ThingHistory) error {
	return nil
}

// MoveToHistoryThing moves a thing to history
func (f *Tinkerpop) MoveToHistoryThing(ctx context.Context, thing *models.Thing, UUID strfmt.UUID, deleted bool) error {
	return nil
}

// AddAction adds an action to the Foobar database with the given UUID.
// Takes the action and a UUID as input.
// Action is already validated against the ontology
func (f *Tinkerpop) AddAction(ctx context.Context, action *models.Action, UUID strfmt.UUID) error {

	// If success return nil, otherwise return the error
	return nil
}

// GetAction fills the given ActionGetResponse with the values from the database, based on the given UUID.
func (f *Tinkerpop) GetAction(ctx context.Context, UUID strfmt.UUID, actionResponse *models.ActionGetResponse) error {
	// actionResponse should be populated with the response that comes from the DB.
	// actionResponse = based on the ontology

	// If success return nil, otherwise return the error
	return nil
}

// GetActions fills the given ActionsListResponse with the values from the database, based on the given UUIDs.
func (f *Tinkerpop) GetActions(ctx context.Context, UUIDs []strfmt.UUID, actionsResponse *models.ActionsListResponse) error {
	// If success return nil, otherwise return the error
	return nil
}

// ListActions fills the given ActionListResponse with the values from the database, based on the given parameters.
func (f *Tinkerpop) ListActions(ctx context.Context, UUID strfmt.UUID, first int, offset int, wheres []*connutils.WhereQuery, actionsResponse *models.ActionsListResponse) error {
	// actionsResponse should be populated with the response that comes from the DB.
	// actionsResponse = based on the ontology

	// If success return nil, otherwise return the error
	return nil
}

// UpdateAction updates the Thing in the DB at the given UUID.
func (f *Tinkerpop) UpdateAction(ctx context.Context, action *models.Action, UUID strfmt.UUID) error {

	// If success return nil, otherwise return the error
	return nil
}

// DeleteAction deletes the Action in the DB at the given UUID.
func (f *Tinkerpop) DeleteAction(ctx context.Context, action *models.Action, UUID strfmt.UUID) error {

	// Run the query to delete the action based on its UUID.

	// If success return nil, otherwise return the error
	return nil
}

// HistoryAction fills the history of a Action based on its UUID
func (f *Tinkerpop) HistoryAction(ctx context.Context, UUID strfmt.UUID, history *models.ActionHistory) error {
	return nil
}

// MoveToHistoryAction moves an action to history
func (f *Tinkerpop) MoveToHistoryAction(ctx context.Context, action *models.Action, UUID strfmt.UUID, deleted bool) error {
	return nil
}

// AddKey adds a key to the Foobar database with the given UUID and token.
// UUID  = reference to the key
// token = is the actual access token used in the API's header
func (f *Tinkerpop) AddKey(ctx context.Context, key *models.Key, UUID strfmt.UUID, token string) error {

	// Key struct should be stored

	// If success return nil, otherwise return the error
	return nil
}

// ValidateToken validates/gets a key to the Foobar database with the given token (=UUID)
func (f *Tinkerpop) ValidateToken(ctx context.Context, UUID strfmt.UUID, keyResponse *models.KeyGetResponse) (token string, err error) {

	// key (= models.KeyGetResponse) should be populated with the response that comes from the DB.

	// in case the key is not found, return an error like:
	// return errors_.New("Key not found in database.")

	isRoot := true
	keyResponse.IsRoot = &isRoot
	keyResponse.KeyID = strfmt.UUID(UUID.String())
	keyResponse.Delete = true
	keyResponse.Email = "hello@weaviate.com"
	keyResponse.Execute = true
	//keyResponse.IPOrigin = "127.0.0.1".([]string)
	keyResponse.KeyExpiresUnix = -1
	keyResponse.Read = true
	keyResponse.Write = true

	// If success return nil, otherwise return the error
	return connutils.TokenHasher(UUID), nil
}

// GetKey fills the given KeyGetResponse with the values from the database, based on the given UUID.
func (f *Tinkerpop) GetKey(ctx context.Context, UUID strfmt.UUID, keyResponse *models.KeyGetResponse) error {

	return nil
}

// GetKeys fills the given []KeyGetResponse with the values from the database, based on the given UUIDs.
func (f *Tinkerpop) GetKeys(ctx context.Context, UUIDs []strfmt.UUID, keysResponse *[]*models.KeyGetResponse) error {
	return nil
}

// DeleteKey deletes the Key in the DB at the given UUID.
func (f *Tinkerpop) DeleteKey(ctx context.Context, key *models.Key, UUID strfmt.UUID) error {
	return nil
}

// GetKeyChildren fills the given KeyGetResponse array with the values from the database, based on the given UUID.
func (f *Tinkerpop) GetKeyChildren(ctx context.Context, UUID strfmt.UUID, children *[]*models.KeyGetResponse) error {

	// for examle: `children = [OBJECT-A, OBJECT-B, OBJECT-C]`
	// Where an OBJECT = models.KeyGetResponse

	return nil
}

// UpdateKey updates the Key in the DB at the given UUID.
func (f *Tinkerpop) UpdateKey(ctx context.Context, key *models.Key, UUID strfmt.UUID, token string) error {
	return nil
}
