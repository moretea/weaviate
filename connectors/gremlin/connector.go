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

package gremlin

import (
	"context"
	errors_ "errors"
	"fmt"
	"runtime"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"github.com/qasaur/gremgo"

	"github.com/creativesoftwarefdn/weaviate/config"
	"github.com/creativesoftwarefdn/weaviate/messages"
	"github.com/creativesoftwarefdn/weaviate/models"
	"github.com/creativesoftwarefdn/weaviate/schema"
)

// Gremlin has some basic variables.
// This is mandatory, only change it if you need aditional, global variables
type Gremlin struct {
	client *gremgo.Client
	kind   string

	config        Config
	serverAddress string
	schema        *schema.WeaviateSchema
	messaging     *messages.Messaging
}

// Config represents the config outline for Gremlin. The Database config shoud be of the following form:
// "database_config" : {
//     "host": "127.0.0.1",
//     "port": 9080
// }
// Notice that the port is the GRPC-port.
type Config struct {
	Host string
	Port int
}

// Edge results from Gremlin
type Edges [][]struct {
	ID         int              `json:"id"`
	InV        int              `json:"inV"`
	InVLabel   string           `json:"inVLabel"`
	Label      string           `json:"label"`
	OutV       int              `json:"outV"`
	OutVLabel  string           `json:"outVLabel"`
	Properties models.SingleRef `json:"properties"`
	Type       string           `json:"type"`
}

func (f *Gremlin) trace() {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f2 := runtime.FuncForPC(pc[0])
	//file, line := f2.FileLine(pc[0])
	fmt.Printf("THIS FUNCTION RUNS: %s\n", f2.Name())
}

// GetName returns a unique connector name, this name is used to define the connector in the weaviate config
func (f *Gremlin) GetName() string {
	return "gremlin"
}

// SetConfig sets variables, which can be placed in the config file section "database_config: {}"
// can be custom for any connector, in the example below there is only host and port available.
//
// Important to bear in mind;
// 1. You need to add these to the struct Config in this document.
// 2. They will become available via f.config.[variable-name]
//
// 	"database": {
// 		"name": "gremlin",
// 		"database_config" : {
// 			"host": "127.0.0.1",
// 			"port": 9080
// 		}
// 	},
func (f *Gremlin) SetConfig(configInput *config.Environment) error {

	// Mandatory: needed to add the JSON config represented as a map in f.config
	err := mapstructure.Decode(configInput.Database.DatabaseConfig, &f.config)

	// Example to: Validate if the essential  config is available, like host and port.
	if err != nil || len(f.config.Host) == 0 || f.config.Port == 0 {
		return errors_.New("could not get Gremlin host/port from config")
	}

	// If success return nil, otherwise return the error (see above)
	return nil
}

// SetSchema takes actionSchema and thingsSchema as an input and makes them available globally at f.schema
// In case you want to modify the schema, this is the place to do so.
// Note: When this function is called, the schemas (action + things) are already validated, so you don't have to build the validation.
func (f *Gremlin) SetSchema(schemaInput *schema.WeaviateSchema) error {
	f.schema = schemaInput

	// If success return nil, otherwise return the error
	return nil
}

// SetMessaging is used to send messages to the service.
// Available message types are: f.messaging.Infomessage ...DebugMessage ...ErrorMessage ...ExitError (also exits the service) ...InfoMessage
func (f *Gremlin) SetMessaging(m *messages.Messaging) error {

	// mandatory, adds the message functions to f.messaging to make them globally accessible.
	f.messaging = m

	// If success return nil, otherwise return the error
	return nil
}

// SetServerAddress is used to fill in a global variable with the server address, but can also be used
// to do some custom actions.
// Does not return anything
func (f *Gremlin) SetServerAddress(addr string) {
	f.serverAddress = addr
}

func (f *Gremlin) Connect() error {

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
func (f *Gremlin) Init() error {

	/*
	 * 1.  If a schema is needed, you need to add the schema to the DB here.
	 * 1.1 Create the (thing or action) classes first, classes that a node (subject or object) can have (for example: Building, Person, etcetera)
	 * 2.  Create a root key.
	 */

	// Example of creating rootkey
	//
	// Add ROOT-key if not exists
	// Search for Root key

	// SEARCH FOR ROOTKEY

	//if totalResult.Root.Count == 0 {
	//	f.messaging.InfoMessage("No root-key found.")
	//
	//	// Create new object and fill it
	//	keyObject := models.Key{}
	//	token := connutils.CreateRootKeyObject(&keyObject)
	//
	//	err = f.AddKey(&keyObject, connutils.GenerateUUID(), token)
	//
	//	if err != nil {
	//		return err
	//	}
	//}
	// END KEYS

	// If success return nil, otherwise return the error
	return nil
}

// Attach can attach something to the request-context
func (f *Gremlin) Attach(ctx context.Context) (context.Context, error) {
	return ctx, nil
}
