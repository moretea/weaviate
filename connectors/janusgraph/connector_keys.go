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

package janusgraph

import (
	"context"
	errors_ "errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-openapi/strfmt"

	"github.com/creativesoftwarefdn/weaviate/models"
)

// AddKey adds a key to the Foobar database with the given UUID and token.
// UUID  = reference to the key
// token = is the actual access token used in the API's header
func (f *Janusgraph) AddKey(ctx context.Context, key *models.Key, UUID strfmt.UUID, token string) error {

	vertex, edge := f.keyToJanusgraph(key, token, UUID.String(), "add")

	addVertexResult, err := f.client.Execute(
		vertex,
		map[string]string{},
		map[string]string{},
	)

	fmt.Println("ADDRESULT KEYS")
	fmt.Println(vertex)
	fmt.Println(addVertexResult)

	// on process error, fail
	if reflect.TypeOf(addVertexResult.([]interface{})[0]).String() == "*errors.errorString" {
		// not returning the error because it is a go routine and the error message will arrive after the fact
		f.messaging.ErrorMessage("Janusgraph [ADD]: " + "[SCRIPT EVALUATION ERROR]")
	}

	// on error, fail
	if err != nil {
		return err
	}

	// check if there is a parent
	if edge != "" {
		addEdgeResult, err := f.client.Execute(
			edge,
			map[string]string{},
			map[string]string{},
		)

		// on process error, fail
		if reflect.TypeOf(addEdgeResult.([]interface{})[0]).String() == "*errors.errorString" {
			// not returning the error because it is a go routine and the error message will arrive after the fact
			f.messaging.ErrorMessage("Janusgraph [ADD]: " + "[SCRIPT EVALUATION ERROR]")
		}

		// on error, fail
		if err != nil {
			return err
		}
	}

	// If success return nil, otherwise return the error
	return nil
}

// GetKey fills the given KeyGetResponse with the values from the database, based on the given UUID.
func (f *Janusgraph) GetKey(ctx context.Context, UUID strfmt.UUID, keyResponse *models.KeyGetResponse) error {

	// find the edges (if any)
	getKey, err := f.client.Execute(
		`g.V().hasLabel("key").has("type", "key").has("uuid", "`+UUID.String()+`")`,
		map[string]string{},
		map[string]string{},
	)

	// return error on error
	if err != nil {
		return err
	}

	// nothing is found
	if getKey.([]interface{})[0] == nil {
		return errors_.New("No key found")
	}

	// set keyResonse values
	keyResponse.KeyID = strfmt.UUID(f.getSinglePropertyValue(getKey, "uuid", 0).(string))
	keyResponse.KeyExpiresUnix = int64(f.getSinglePropertyValue(getKey, "keyExpiresUnix", 0).(float64))
	keyResponse.Write = f.getSinglePropertyValue(getKey, "write", 0).(bool)
	keyResponse.Email = f.getSinglePropertyValue(getKey, "email", 0).(string)
	keyResponse.Read = f.getSinglePropertyValue(getKey, "read", 0).(bool)
	keyResponse.Delete = f.getSinglePropertyValue(getKey, "delete", 0).(bool)
	keyResponse.Execute = f.getSinglePropertyValue(getKey, "execute", 0).(bool)
	keyResponse.IPOrigin = strings.Split(f.getSinglePropertyValue(getKey, "IPOrigin", 0).(string), ";")

	return nil

}

// GetKeys fills the given []KeyGetResponse with the values from the database, based on the given UUIDs.
func (f *Janusgraph) GetKeys(ctx context.Context, UUIDs []strfmt.UUID, keysResponse *[]*models.KeyGetResponse) error {

	fmt.Println("GETKEYS")

	return errors_.New("No key found")
}

// DeleteKey deletes the Key in the DB at the given UUID.
func (f *Janusgraph) DeleteKey(ctx context.Context, key *models.Key, UUID strfmt.UUID) error {
	// Remove based on type and uuid
	_, err := f.client.Execute(
		`g.V().hasLabel("key").has("uuid", `+UUID.String()+`).drop()`,
		map[string]string{},
		map[string]string{},
	)

	// return error
	if err != nil {
		return err
	}

	// If success return nil, otherwise return the error
	return nil
}

// GetKeyChildren fills the given KeyGetResponse array with the values from the database, based on the given UUID.
func (f *Janusgraph) GetKeyChildren(ctx context.Context, UUID strfmt.UUID, children *[]*models.KeyGetResponse) error {

	// find the edges (if any)
	getKeys, err := f.client.Execute(
		`g.V().hasLabel("key").has("type", "key").has("uuid", "`+UUID.String()+`").inE().outV().properties("uuid")`,
		map[string]string{},
		map[string]string{},
	)

	// return error on error
	if err != nil {
		return err
	}

	// validate if there are any results
	if reflect.TypeOf(getKeys.([]interface{})[0]) != nil {
		// loop over the results and add to the child
		for _, value := range getKeys.([]interface{})[0].([]interface{}) {
			for propK, propV := range value.(map[string]interface{}) {
				if propK == "value" {
					// create empty child struct
					child := models.KeyGetResponse{}
					// get the key values
					f.GetKey(ctx, strfmt.UUID(propV.(string)), &child)
					// append
					*children = append(*children, &child)
				}
			}
		}
	}

	return nil
}

// UpdateKey updates the Key in the DB at the given UUID.
func (f *Janusgraph) UpdateKey(ctx context.Context, key *models.Key, UUID strfmt.UUID, token string) error {

	// getting the vertex and the edge
	vertex, edge := f.keyToJanusgraph(key, token, UUID.String(), "update")

	// update the key
	_, err := f.client.Execute(
		vertex,
		map[string]string{},
		map[string]string{},
	)

	// return error on error
	if err != nil {
		return err
	}

	// update the edge of the key
	_, err = f.client.Execute(
		edge,
		map[string]string{},
		map[string]string{},
	)

	// return error on error
	if err != nil {
		return err
	}

	return nil
}
