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
	"strconv"

	"github.com/go-openapi/strfmt"

	"github.com/creativesoftwarefdn/weaviate/connectors/utils"
	"github.com/creativesoftwarefdn/weaviate/models"
)

// AddAction adds an action to the Foobar database with the given UUID.
// Takes the action and a UUID as input.
// Action is already validated against the ontology
func (f *Gremlin) AddAction(ctx context.Context, action *models.Action, UUID strfmt.UUID) error {

	// convert the thing to a Vertex and and Edge.
	err := f.actionToGremlin(UUID, action, "add")

	// on error fail
	if err != nil {
		return err
	}

	err = f.connectToKey(UUID, action.Key.NrDollarCref, "action")

	// on error fail
	if err != nil {
		return err
	}

	// If success return nil, otherwise return the error
	return nil
}

// GetAction fills the given ActionGetResponse with the values from the database, based on the given UUID.
func (f *Gremlin) GetAction(ctx context.Context, UUID strfmt.UUID, actionResponse *models.ActionGetResponse) error {
	// define the ID vertex and the UUID to fetch
	result, err := f.client.Execute(
		`g.V().has("uuid", uuid).has("type", objectType)`,
		map[string]string{"uuid": UUID.String(), "objectType": "action"},
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

	// process the thing repsonse for a single thing. thingNo = 0, because we want the first one
	err = f.processSingleAction(result, 0, actionResponse)

	// in case of error, return the error
	if err != nil {
		return err
	}

	// If success return nil, otherwise return the error
	return nil
}

// GetActions fills the given ActionsListResponse with the values from the database, based on the given UUIDs.
func (f *Gremlin) GetActions(ctx context.Context, UUIDs []strfmt.UUID, actionsResponse *models.ActionsListResponse) error {
	// If success return nil, otherwise return the error
	return nil
}

// ListActions fills the given ActionListResponse with the values from the database, based on the given parameters.
func (f *Gremlin) ListActions(ctx context.Context, UUID strfmt.UUID, first int, offset int, wheres []*connutils.WhereQuery, actionsResponse *models.ActionsListResponse) error {
	// find the edges (if any)
	result, err := f.client.Execute(
		`g.V().has("type", objectType).range(`+strconv.Itoa((first*offset))+`, `+strconv.Itoa(first)+`)`,
		map[string]string{"objectType": "action"},
		map[string]string{},
	)

	// nothing is found
	if result.([]interface{})[0] == nil {
		return nil
	}

	// Loop over the Gremlin results
	for actionKey := range result.([]interface{})[0].([]interface{}) {

		// define singleThing
		var singleAction models.ActionGetResponse

		// get the individual thing
		err := f.processSingleAction(result, actionKey, &singleAction)

		// in case of error, return the error
		if err != nil {
			return err
		}

		// add the thing to the response
		actionsResponse.Actions = append(actionsResponse.Actions, &singleAction)
	}

	// add the final results
	actionsResponse.TotalResults = int64(len(result.([]interface{})[0].([]interface{})))

	// in case of error, return the error
	if err != nil {
		return err
	}

	// If success return nil, otherwise return the error
	return nil
}

// UpdateAction updates the Thing in the DB at the given UUID.
func (f *Gremlin) UpdateAction(ctx context.Context, action *models.Action, UUID strfmt.UUID) error {

	// get the vertexes
	err := f.actionToGremlin(UUID, action, "update")

	if err != nil {
		return err
	}

	// If success return nil, otherwise return the error
	return nil

}

// DeleteAction deletes the Action in the DB at the given UUID.
func (f *Gremlin) DeleteAction(ctx context.Context, action *models.Action, UUID strfmt.UUID) error {

	// Remove based on type and uuid
	_, err := f.client.Execute(
		`g.V().has("uuid", uuid).has("type", objectType).drop()`,
		map[string]string{"uuid": UUID.String(), "objectType": "action"},
		map[string]string{},
	)

	// return error
	if err != nil {
		return err
	}

	// If success return nil, otherwise return the error
	return nil
}

// HistoryAction fills the history of a Action based on its UUID
func (f *Gremlin) HistoryAction(ctx context.Context, UUID strfmt.UUID, history *models.ActionHistory) error {
	return nil
}

// MoveToHistoryAction moves an action to history
func (f *Gremlin) MoveToHistoryAction(ctx context.Context, action *models.Action, UUID strfmt.UUID, deleted bool) error {
	return nil
}
