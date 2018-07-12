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
	"encoding/json"
	"strconv"

	"github.com/go-openapi/strfmt"

	"github.com/creativesoftwarefdn/weaviate/connectors/utils"
	"github.com/creativesoftwarefdn/weaviate/models"
)

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
		// find properties if all properties are selected (i.e., no `.properties(...)` in the Gremlin query)
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

// connects a vertex to a key
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

// get a Key UUID from an Edge
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
