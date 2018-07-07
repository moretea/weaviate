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

package tinkerpop

import (
	"context"

	"github.com/go-openapi/strfmt"

	"github.com/creativesoftwarefdn/weaviate/connectors/utils"
	"github.com/creativesoftwarefdn/weaviate/models"
)

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
