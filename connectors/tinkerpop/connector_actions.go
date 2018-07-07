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
