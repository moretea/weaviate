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
	"encoding/base64"
	"reflect"
	"strings"

	"github.com/go-openapi/strfmt"

	"github.com/creativesoftwarefdn/weaviate/models"
)

// ValidateToken validates/gets a key to the Grelmin database with the given token (=UUID)
func (f *Janusgraph) ValidateToken(ctx context.Context, UUID strfmt.UUID, keyResponse *models.KeyGetResponse) (token string, err error) {

	// find the key
	getKey, err := f.client.Execute(
		`g.V().hasLabel("key").has("type", "key").has("uuid", "`+UUID.String()+`")`,
		map[string]string{},
		map[string]string{},
	)

	// nothing is found
	if getKey.([]interface{})[0] == nil {
		return "", nil
	}

	// on process error, fail
	if reflect.TypeOf(getKey.([]interface{})[0]).String() == "*errors.errorString" {
		// not returning the error because it is a go routine and the error message will arrive after the fact
		f.messaging.ErrorMessage("Janusgraph [ADD]: " + "[SCRIPT EVALUATION ERROR]")
	}

	// on error, fail
	if err != nil {
		return "", err
	}

	// set key responses
	keyResponse.KeyID = strfmt.UUID(f.getSinglePropertyValue(getKey, "uuid", 0).(string))
	keyResponse.KeyExpiresUnix = int64(f.getSinglePropertyValue(getKey, "keyExpiresUnix", 0).(float64))
	keyResponse.Write = f.getSinglePropertyValue(getKey, "write", 0).(bool)
	keyResponse.Email = f.getSinglePropertyValue(getKey, "email", 0).(string)
	keyResponse.Read = f.getSinglePropertyValue(getKey, "read", 0).(bool)
	keyResponse.Delete = f.getSinglePropertyValue(getKey, "delete", 0).(bool)
	keyResponse.Execute = f.getSinglePropertyValue(getKey, "execute", 0).(bool)
	keyResponse.IPOrigin = strings.Split(f.getSinglePropertyValue(getKey, "IPOrigin", 0).(string), ";")

	// check if root
	isRoot := f.getSinglePropertyValue(getKey, "isRoot", 0).(bool)
	keyResponse.IsRoot = &isRoot

	// get the token
	tokenToReturn, err := base64.StdEncoding.DecodeString(f.getSinglePropertyValue(getKey, "__token", 0).(string))
	if err != nil {
		return "", err
	}

	// If success return nil, otherwise return the error
	return string(tokenToReturn), nil

}
