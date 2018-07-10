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
	"bytes"
	"encoding/base64"
	"encoding/json"
	errors_ "errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/creativesoftwarefdn/weaviate/connectors/utils"
	"github.com/creativesoftwarefdn/weaviate/models"
)

// helper function to get all things
func (f *Gremlin) _getAll(what string) {

	fmt.Println("--------")
	fmt.Println(what)

	all, _ := f.client.Execute(
		`g.V().hasLabel("thing").has("uuid", "78bf463c-c992-443b-82f6-9157765c570b")`,
		map[string]string{},
		map[string]string{},
	)

	edgesBytes, _ := json.Marshal(all)

	fmt.Println(string(edgesBytes))
	fmt.Println("--------")

}

// keyToGremlin translates a thing struct into a vertex string
func (f *Gremlin) keyToGremlin(objToHandle *models.Key, token string, UUID string) (string, string) {

	// create vertex and edge buffers
	var vertex bytes.Buffer
	var edge bytes.Buffer

	// set general fields
	vertex.WriteString(`g.addV("key").property("uuid", "` + UUID + `").property("type", "key")`)

	// set boolean properties
	isRoot := objToHandle.Parent == nil

	vertex.WriteString(`.property("isRoot", ` + strconv.FormatBool(isRoot) + `)`)
	vertex.WriteString(`.property("delete", ` + strconv.FormatBool(objToHandle.Delete) + `)`)
	vertex.WriteString(`.property("execute", ` + strconv.FormatBool(objToHandle.Execute) + `)`)
	vertex.WriteString(`.property("read", ` + strconv.FormatBool(objToHandle.Read) + `)`)
	vertex.WriteString(`.property("write", ` + strconv.FormatBool(objToHandle.Write) + `)`)

	// set string properties
	vertex.WriteString(`.property("email", "` + objToHandle.Email + `")`)

	// set array properties
	vertex.WriteString(`.property("IPOrigin", "` + strings.Join(objToHandle.IPOrigin, ";") + `")`)

	// set integers
	vertex.WriteString(`.property("keyExpiresUnix", ` + strconv.FormatInt(objToHandle.KeyExpiresUnix, 10) + `)`)

	// add the secured and encrypted token
	vertex.WriteString(`.property("__token", "` + base64.StdEncoding.EncodeToString([]byte(token)) + `")`)

	// if isRoot is false, an edge needs to be added
	if isRoot == false {
		fmt.Println(objToHandle.Parent)
	}

	return vertex.String(), edge.String()
}

// thingToGremlin translates a thing struct into a vertex string
func (f *Gremlin) thingToGremlin(UUID strfmt.UUID, objToHandle *models.Thing, addOrUpdate string) error {

	typeOfObject := "thing"

	// start vertex string Buffer
	var vertex bytes.Buffer
	var update bytes.Buffer
	var edge []string

	// add vertext class and ID
	if addOrUpdate == "add" {
		vertex.WriteString(`g.addV("` + typeOfObject + `")`)
	} else if addOrUpdate == "update" {
		vertex.WriteString(`g.V().hasLabel("` + typeOfObject + `")`)
		update.WriteString(`g.V().hasLabel("` + typeOfObject + `")`)
	} else {
		return errors_.New("addOrUpdate should be 'add' or 'update'")
	}

	// define the type and the ID
	vertex.WriteString(`.property("uuid", "` + UUID.String() + `").property("type", "` + typeOfObject + `")`)
	if addOrUpdate == "update" {
		update.WriteString(`.property("uuid", "` + UUID.String() + `").property("type", "` + typeOfObject + `").outE().drop()`)
	}

	// set the meta values (@class will be the vector)
	vertex.WriteString(`.property("atClass", "` + objToHandle.AtClass + `")`)
	vertex.WriteString(`.property("context", "` + objToHandle.AtContext + `")`)
	vertex.WriteString(`.property("creationTimeUnix", "` + strconv.FormatInt(objToHandle.CreationTimeUnix, 10) + `")`)
	vertex.WriteString(`.property("lastUpdateTimeUnix", "` + strconv.FormatInt(objToHandle.LastUpdateTimeUnix, 10) + `")`)

	// reflect in Schema
	schema := reflect.ValueOf(objToHandle.Schema)

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
			case interface{}:
				if reflect.TypeOf(v.Interface()).String() == "time.Time" { // in case of time, store as date
					vertex.WriteString(`.property("schema__` + e.String() + `", "` + time.Time.String(t.(time.Time)) + `")`)
				} else if reflect.TypeOf(v.Interface()).String() == "*models.SingleRef" { // in case of SingleRef, store as relation
					// Get the singleRef values
					singleRef := v.Interface().(*models.SingleRef)
					// create the edge
					edge = append(edge, `g.addE("schema__`+e.String()+`").from(g.V().hasLabel("`+typeOfObject+`").has("uuid", uuid)).to(g.V().hasLabel("`+typeOfObject+`").has("uuid", "`+singleRef.NrDollarCref.String()+`")).property("\$cref", "`+singleRef.NrDollarCref.String()+`").property("type", "`+singleRef.Type+`").property("locationUrl", "`+*singleRef.LocationURL+`")`)
				} else {
					f.messaging.ExitError(78, "The type "+reflect.TypeOf(v.Interface()).String()+" is not found.")
				}
			default:
				f.messaging.ExitError(78, "The type "+reflect.TypeOf(v.Interface()).String()+" is not found.")
			}
		}
	}

	// execute the Vertex query with uuid as binding, result is not used because we send out "202 Accepted" and assume a succes because validation takes place before this function runs
	addResult, err := f.client.Execute(
		vertex.String(),
		map[string]string{},
		map[string]string{},
	)

	// on process error, fail
	if reflect.TypeOf(addResult.([]interface{})[0]).String() == "*errors.errorString" {
		// not returning the error because it is a go routine and the error message will arrive after the fact
		f.messaging.ErrorMessage("Gremlin [ADD]: " + "[SCRIPT EVALUATION ERROR]")
	}

	// on error, fail
	if err != nil {
		return err
	}

	// drop all edges when updating
	if addOrUpdate == "update" {

		// execute the Edge query
		updateResult, err := f.client.Execute(
			update.String(),
			map[string]string{},
			map[string]string{},
		)

		// on process error, fail
		if reflect.TypeOf(updateResult.([]interface{})[0]) != nil {
			// not returning the error because it is a go routine and the error message will arrive after the fact
			f.messaging.ErrorMessage("Gremlin [UPDATE]: " + "[SCRIPT EVALUATION ERROR]")
		}

		// on error, fail
		if err != nil {
			return err
		}
	}

	// if there are any edges set...
	if len(edge) > 0 {
		// loop over edges that need to be added
		for _, singleEdge := range edge {

			// execute the Edge query
			addEdgeResult, err := f.client.Execute(
				singleEdge,
				map[string]string{"uuid": UUID.String()},
				map[string]string{},
			)

			// on process error, fail
			if reflect.TypeOf(addEdgeResult.([]interface{})[0]).String() == "*errors.errorString" {
				// not returning the error because it is a go routine and the error message will arrive after the fact
				f.messaging.ErrorMessage("Gremlin [EDGE CREATION]: " + "[SCRIPT EVALUATION ERROR]")
			}

			// on error, fail
			if err != nil {
				return err
			}
		}
	}

	// return the vertex and the edge map
	return nil
}

// thingToVertex translates a thing struct into a vertex string
func (f *Gremlin) actionToGremlin(UUID strfmt.UUID, objToHandle *models.Action, addOrUpdate string) error {

	typeOfObject := "action"

	// start vertex string Buffer
	var vertex bytes.Buffer
	var update bytes.Buffer
	var edge []string

	// add vertext class and ID
	if addOrUpdate == "add" {
		vertex.WriteString(`g.addV("` + typeOfObject + `")`)
	} else if addOrUpdate == "update" {
		vertex.WriteString(`g.V().hasLabel("` + typeOfObject + `")`)
		update.WriteString(`g.V().hasLabel("` + typeOfObject + `")`)
	} else {
		return errors_.New("addOrUpdate should be 'add' or 'update'")
	}

	// define the type and the ID
	vertex.WriteString(`.property("uuid", uuid).property("type", "` + typeOfObject + `")`)
	if addOrUpdate == "update" {
		update.WriteString(`.property("uuid", uuid).property("type", "` + typeOfObject + `").outE().drop()`)
	}

	// set the meta values (@class will be the vector)
	vertex.WriteString(`.property("context", "` + objToHandle.AtContext + `")`)
	vertex.WriteString(`.property("creationTimeUnix", "` + strconv.FormatInt(objToHandle.CreationTimeUnix, 10) + `")`)
	vertex.WriteString(`.property("lastUpdateTimeUnix", "` + strconv.FormatInt(objToHandle.LastUpdateTimeUnix, 10) + `")`)

	// reflect in Schema
	schema := reflect.ValueOf(objToHandle.Schema)

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
			case interface{}:
				if reflect.TypeOf(v.Interface()).String() == "time.Time" { // in case of time, store as date
					vertex.WriteString(`.property("schema__` + e.String() + `", "` + time.Time.String(t.(time.Time)) + `")`)
				} else if reflect.TypeOf(v.Interface()).String() == "*models.SingleRef" { // in case of SingleRef, store as relation
					// Get the singleRef values
					singleRef := v.Interface().(*models.SingleRef)
					// create the edge
					edge = append(edge, `g.addE("schema__`+e.String()+`").from(g.V().hasLabel("`+typeOfObject+`").has("uuid", uuid)).to(g.V().hasLabel("`+typeOfObject+`").has("uuid", "`+singleRef.NrDollarCref.String()+`")).property("\$cref", "`+singleRef.NrDollarCref.String()+`").property("type", "`+singleRef.Type+`").property("locationUrl", "`+*singleRef.LocationURL+`")`)
				} else {
					f.messaging.ExitError(78, "The type "+reflect.TypeOf(v.Interface()).String()+" is not found.")
				}
			default:
				f.messaging.ExitError(78, "The type "+reflect.TypeOf(v.Interface()).String()+" is not found.")
			}
		}
	}

	// execute the Vertex query with uuid as binding, result is not used because we send out "202 Accepted" and assume a succes because validation takes place before this function runs
	addResult, err := f.client.Execute(
		vertex.String(),
		map[string]string{"uuid": UUID.String()},
		map[string]string{},
	)

	// on process error, fail
	if reflect.TypeOf(addResult.([]interface{})[0]).String() == "*errors.errorString" {
		// not returning the error because it is a go routine and the error message will arrive after the fact
		f.messaging.ErrorMessage("Gremlin [ADD]: " + "[SCRIPT EVALUATION ERROR]")
	}

	// on error, fail
	if err != nil {
		return err
	}

	// drop all edges when updating
	if addOrUpdate == "update" {

		// execute the Edge query
		updateResult, err := f.client.Execute(
			update.String(),
			map[string]string{"uuid": UUID.String()},
			map[string]string{},
		)

		// on process error, fail
		if reflect.TypeOf(updateResult.([]interface{})[0]) != nil {
			// not returning the error because it is a go routine and the error message will arrive after the fact
			f.messaging.ErrorMessage("Gremlin [UPDATE]: " + "[SCRIPT EVALUATION ERROR]")
		}

		// on error, fail
		if err != nil {
			return err
		}
	}

	// if there are any edges set...
	if len(edge) > 0 {
		// loop over edges that need to be added
		for _, singleEdge := range edge {

			// execute the Edge query
			addEdgeResult, err := f.client.Execute(
				singleEdge,
				map[string]string{"uuid": UUID.String()},
				map[string]string{},
			)

			// on process error, fail
			if reflect.TypeOf(addEdgeResult.([]interface{})[0]).String() == "*errors.errorString" {
				// not returning the error because it is a go routine and the error message will arrive after the fact
				f.messaging.ErrorMessage("Gremlin [EDGE CREATION]: " + "[SCRIPT EVALUATION ERROR]")
			}

			// on error, fail
			if err != nil {
				return err
			}
		}
	}

	// return the vertex and the edge map
	return nil
}

// createCrefObject is a helper function to create a cref-object. This function is used for Cassandra only.
func (f *Gremlin) createCrefObject(UUID strfmt.UUID, location string, refType connutils.RefType) *models.SingleRef {
	// Create the 'cref'-node for the response.
	crefObj := models.SingleRef{}

	// Get the given node properties to generate response object
	crefObj.NrDollarCref = UUID
	crefObj.Type = string(refType)
	url := location
	crefObj.LocationURL = &url

	return &crefObj
}

// gets a single value from and id value pair
func (f *Gremlin) getSinglePropertyValue(haystack interface{}, needle string, no int) interface{} {
	// Loop over the Gremlin results
	for k, v := range haystack.([]interface{})[0].([]interface{})[no].(map[string]interface{}) {
		// add schema properties
		if k == "properties" {
			// find the property in the schema
			for propK, propV := range v.(map[string]interface{}) {
				if propK == needle {
					// find the value
					for propKSingle, propVSingle := range propV.([]interface{})[0].(map[string]interface{}) {
						if propKSingle == "value" {
							return propVSingle
						}
					}
				}
			}
		}
	}
	return nil
}

func (f *Gremlin) connectToKey(UUID strfmt.UUID, KeyUUID strfmt.UUID, typeToConnect string) error {

	// execute the Edge query
	_, err := f.client.Execute(
		`g.addE("key").from(g.V().hasLabel("`+typeToConnect+`").has("uuid", "`+UUID.String()+`")).to(g.V().hasLabel("key").has("uuid", "`+KeyUUID.String()+`")).property("fromThing", "`+UUID.String()+`").property("toKeyUUID", "`+KeyUUID.String()+`")`,
		map[string]string{},
		map[string]string{},
	)

	// on error, fail
	if err != nil {
		return err
	}

	return nil
}

func (f *Gremlin) keyUUIDFromEdge(edgeValue interface{}) (strfmt.UUID, error) {

	keyEdges := KeyEdge{}

	keyEdgesBytes, err := json.Marshal(edgeValue)

	if err != nil {
		return strfmt.UUID(0), err
	}

	err = json.Unmarshal(keyEdgesBytes, &keyEdges)

	if err != nil {
		return strfmt.UUID(0), err
	}

	// find the edges (if any)
	keyEdgesResult, err := f.client.Execute(
		`g.V().hasId("`+strconv.Itoa(keyEdges.InV)+`")`,
		map[string]string{},
		map[string]string{},
	)

	// in case of error, return the error
	if err != nil {
		return strfmt.UUID(0), err
	}

	return strfmt.UUID(f.getSinglePropertyValue(keyEdgesResult, "uuid", 0).(string)), nil
}

// processes a single thing
func (f *Gremlin) processSingleThing(result interface{}, thingNo int, thingResponse *models.ThingGetResponse) error {

	// Create the schema Map, this map will contain all the results
	responseSchema := make(map[string]interface{})

	// set meta values (String)
	thingResponse.ThingID = strfmt.UUID(f.getSinglePropertyValue(result, "uuid", thingNo).(string))
	thingResponse.AtClass = f.getSinglePropertyValue(result, "atClass", thingNo).(string)
	thingResponse.AtContext = f.getSinglePropertyValue(result, "context", thingNo).(string)

	// set meta values (int64)
	thingResponse.CreationTimeUnix, _ = strconv.ParseInt(f.getSinglePropertyValue(result, "creationTimeUnix", thingNo).(string), 10, 64)
	thingResponse.LastUpdateTimeUnix, _ = strconv.ParseInt(f.getSinglePropertyValue(result, "lastUpdateTimeUnix", thingNo).(string), 10, 64)

	// Loop over the Gremlin schema, results
	for key, value := range result.([]interface{})[0].([]interface{})[thingNo].(map[string]interface{}) {
		// add schema properties
		if key == "properties" {

			for propKey, propValue := range value.(map[string]interface{}) {

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

		// find the edges (if any)
		result, err := f.client.Execute(
			`g.V().hasLabel("thing").has("uuid", "`+string(thingResponse.ThingID)+`").outE()`,
			map[string]string{},
			map[string]string{},
		)

		// in case of error, return the error
		if err != nil {
			return err
		}

		// define Edges
		var edges Edges

		// edges to JSON
		edgesBytes, err := json.Marshal(result)

		// in case of error, return the error
		if err != nil {
			return err
		}

		// merge edges
		err = json.Unmarshal(edgesBytes, &edges)

		// in case of error, return the error
		if err != nil {
			return err
		}

		// add the properties to the edges. Note that the NrDollarCref is _not_ the Gremlin/Gremlin ID but the Weaviate UUID
		for _, edgeValue := range edges[0] {
			if len(edgeValue.Label) >= 8 { // should be larger than 8
				if edgeValue.Label[0:8] == "schema__" { // only handle schema edges
					responseSchema[edgeValue.Label[8:]] = models.SingleRef{
						NrDollarCref: edgeValue.Properties.NrDollarCref,
						Type:         edgeValue.Properties.Type,
						LocationURL:  edgeValue.Properties.LocationURL,
					}
				}
			} else if edgeValue.Label == "key" { // smaller then 8 and = "key"
				// get the related key and return
				keyUUID, err := f.keyUUIDFromEdge(edgeValue)

				// in case of error, return the error
				if err != nil {
					return err
				}

				thingResponse.Key = f.createCrefObject(keyUUID, f.serverAddress, connutils.RefTypeKey)

			}
		}

		// in case of error, return the error
		if err != nil {
			return err
		}

		// schema = responseSchema
		thingResponse.Schema = responseSchema

	}

	// success, return nil
	return nil

}

// processes a single action
func (f *Gremlin) processSingleAction(result interface{}, actionNo int, actionResponse *models.ActionGetResponse) error {

	// This is a temporary key [FIX THIS]
	//actionResponse.Key = f.createCrefObject("29ece1d8-c433-4757-b258-0b278478e17a", f.serverAddress, connutils.RefTypeKey)

	// Create the schema Map, this map will contain all the results
	responseSchema := make(map[string]interface{})

	// set meta values (String)
	actionResponse.ActionID = strfmt.UUID(f.getSinglePropertyValue(result, "uuid", 0).(string))
	actionResponse.AtClass = f.getSinglePropertyValue(result, "atClass", 0).(string)
	actionResponse.AtContext = f.getSinglePropertyValue(result, "context", 0).(string)

	// set meta values (int64)
	actionResponse.CreationTimeUnix, _ = strconv.ParseInt(f.getSinglePropertyValue(result, "creationTimeUnix", 0).(string), 10, 64)
	actionResponse.LastUpdateTimeUnix, _ = strconv.ParseInt(f.getSinglePropertyValue(result, "lastUpdateTimeUnix", 0).(string), 10, 64)

	// Loop over the Gremlin results
	for key, value := range result.([]interface{})[0].([]interface{})[actionNo].(map[string]interface{}) {

		// add schema properties
		if key == "properties" {
			for propKey, propValue := range value.(map[string]interface{}) {

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

		// find the edges (if any)
		result, err := f.client.Execute(
			`g.V().hasLabel("key").has("uuid", `+string(actionResponse.ActionID)+`).outE()`,
			map[string]string{},
			map[string]string{},
		)

		// in case of error, return the error
		if err != nil {
			return err
		}

		// define Edges
		var edges Edges

		// edges to JSON
		edgesBytes, err := json.Marshal(result)

		// in case of error, return the error
		if err != nil {
			return err
		}

		// merge edges
		err = json.Unmarshal(edgesBytes, &edges)

		// in case of error, return the error
		if err != nil {
			return err
		}

		// add the properties to the edges. Note that the NrDollarCref is _not_ the Gremlin/Gremlin ID but the Weaviate UUID
		for _, edgeValue := range edges[0] {
			if len(edgeValue.Label) >= 8 { // should be larger than 8
				if edgeValue.Label[0:8] == "schema__" { // only handle schema edges
					responseSchema[edgeValue.Label[8:]] = models.SingleRef{
						NrDollarCref: edgeValue.Properties.NrDollarCref,
						Type:         edgeValue.Properties.Type,
						LocationURL:  edgeValue.Properties.LocationURL,
					}
				}
			} else if edgeValue.Label == "key" { // smaller then 8 and = "key"
				// get the related key and return
				keyUUID, err := f.keyUUIDFromEdge(edgeValue)

				// in case of error, return the error
				if err != nil {
					return err
				}

				// return key
				actionResponse.Key = f.createCrefObject(keyUUID, f.serverAddress, connutils.RefTypeKey)

			}
		}

		// in case of error, return the error
		if err != nil {
			return err
		}

		// schema = responseSchema
		actionResponse.Schema = responseSchema

	}

	// success, return nil
	return nil

}
