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

	"github.com/go-openapi/strfmt"

	"github.com/creativesoftwarefdn/weaviate/models"
)

// AddKey adds a key to the Foobar database with the given UUID and token.
// UUID  = reference to the key
// token = is the actual access token used in the API's header
func (f *Gremlin) AddKey(ctx context.Context, key *models.Key, UUID strfmt.UUID, token string) error {

	// Key struct should be stored

	// If success return nil, otherwise return the error
	return nil
}

// GetKey fills the given KeyGetResponse with the values from the database, based on the given UUID.
func (f *Gremlin) GetKey(ctx context.Context, UUID strfmt.UUID, keyResponse *models.KeyGetResponse) error {

	return nil
}

// GetKeys fills the given []KeyGetResponse with the values from the database, based on the given UUIDs.
func (f *Gremlin) GetKeys(ctx context.Context, UUIDs []strfmt.UUID, keysResponse *[]*models.KeyGetResponse) error {
	return nil
}

// DeleteKey deletes the Key in the DB at the given UUID.
func (f *Gremlin) DeleteKey(ctx context.Context, key *models.Key, UUID strfmt.UUID) error {
	return nil
}

// GetKeyChildren fills the given KeyGetResponse array with the values from the database, based on the given UUID.
func (f *Gremlin) GetKeyChildren(ctx context.Context, UUID strfmt.UUID, children *[]*models.KeyGetResponse) error {

	// for examle: `children = [OBJECT-A, OBJECT-B, OBJECT-C]`
	// Where an OBJECT = models.KeyGetResponse

	return nil
}

// UpdateKey updates the Key in the DB at the given UUID.
func (f *Gremlin) UpdateKey(ctx context.Context, key *models.Key, UUID strfmt.UUID, token string) error {
	return nil
}
