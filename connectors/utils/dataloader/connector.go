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

package dataloader

import (
	"context"
	"fmt"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/graph-gophers/dataloader"

	"github.com/creativesoftwarefdn/weaviate/config"
	"github.com/creativesoftwarefdn/weaviate/connectors"
	"github.com/creativesoftwarefdn/weaviate/connectors/utils"
	"github.com/creativesoftwarefdn/weaviate/messages"
	"github.com/creativesoftwarefdn/weaviate/models"
	"github.com/creativesoftwarefdn/weaviate/schema"
)

// DataLoader has some basic variables.
type DataLoader struct {
	databaseConnector dbconnector.DatabaseConnector
	messaging         *messages.Messaging
}

const thingsDataLoader string = "thingsDataLoader"
const keysDataLoader string = "keysDataLoader"
const actionsDataLoader string = "actionsDataLoader"

// SetDatabaseConnector sets the used DB-connector
func (f *DataLoader) SetDatabaseConnector(dbConnector dbconnector.DatabaseConnector) {
	f.databaseConnector = dbConnector
}

// GetName returns a unique connector name
func (f *DataLoader) GetName() string {
	return "dataloader"
}

// Connect function
func (f *DataLoader) Connect() error {
	return f.databaseConnector.Connect()
}

// Init function
func (f *DataLoader) Init() error {
	return f.databaseConnector.Init()
}

// Attach function
func (f *DataLoader) Attach(ctx context.Context) (context.Context, error) {
	// Setup batch function
	thingBatchFunc := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		// Do some async work to get data for specified keys
		results := []*dataloader.Result{}

		// Init variables
		things := &models.ThingsListResponse{}
		UUIDs := []strfmt.UUID{}

		// Append resolved values to this list
		for _, v := range keys {
			UUIDs = append(UUIDs, strfmt.UUID(v.String()))
		}

		// Get the batch of Things from the database
		err := f.databaseConnector.GetThings(ctx, UUIDs, things)

		// If the Get failed, append an error for each key to the result
		if err != nil {
			for _ = range keys {
				results = append(results, &dataloader.Result{Error: err})
			}
			return results
		}

		// Generate a response of the same length and type as the batch request
		for _, k := range keys {
			found := false
			for _, v := range things.Things {
				// If the key is same as the ID of the object, append this object
				if k.String() == string(v.ThingID) {
					results = append(results, &dataloader.Result{Data: v})
					found = true
					break
				}
			}
			// If the object ID was not found add an empty object containing an error in its place
			if !found {
				results = append(results, &dataloader.Result{Data: &models.ThingGetResponse{}, Error: fmt.Errorf(connutils.StaticThingNotFound)})
			}
		}

		return results
	}

	// Setup batch function
	actionBatchFunc := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		// Do some async work to get data for specified keys
		results := []*dataloader.Result{}

		// Init variables
		actions := &models.ActionsListResponse{}
		UUIDs := []strfmt.UUID{}

		// Append resolved values to this list
		for _, v := range keys {
			UUIDs = append(UUIDs, strfmt.UUID(v.String()))
		}

		// Get the batch of Actions from the database
		err := f.databaseConnector.GetActions(ctx, UUIDs, actions)

		// If the Get failed, append an error for each key to the result
		if err != nil {
			for _ = range keys {
				results = append(results, &dataloader.Result{Error: err})
			}
			return results
		}

		// Generate a response of the same length and type as the batch request
		for _, k := range keys {
			found := false
			for _, v := range actions.Actions {
				// If the key is same as the ID of the object, append this object
				if k.String() == string(v.ActionID) {
					results = append(results, &dataloader.Result{Data: v})
					found = true
					break
				}
			}
			// If the object ID was not found add an empty object containing an error in its place
			if !found {
				results = append(results, &dataloader.Result{Data: &models.ActionGetResponse{}, Error: fmt.Errorf(connutils.StaticActionNotFound)})
			}
		}

		return results
	}

	// Setup batch function
	keyBatchFunc := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		// Do some async work to get data for specified keys
		results := []*dataloader.Result{}

		// Init variables
		keysList := []*models.KeyGetResponse{}
		UUIDs := []strfmt.UUID{}

		// Append resolved values to this list
		for _, v := range keys {
			UUIDs = append(UUIDs, strfmt.UUID(v.String()))
		}

		// Get the batch of keys from the database
		err := f.databaseConnector.GetKeys(ctx, UUIDs, &keysList)

		// If the Get failed, append an error for each key to the result
		if err != nil {
			for _ = range keys {
				results = append(results, &dataloader.Result{Error: err})
			}
			return results
		}

		// Generate a response of the same length and type as the batch request
		for _, k := range keys {
			found := false
			for _, v := range keysList {
				// If the key is same as the ID of the object, append this object
				if k.String() == string(v.KeyID) {
					results = append(results, &dataloader.Result{Data: v})
					found = true
					break
				}
			}
			// If the object ID was not found add an empty object containing an error in its place
			if !found {
				results = append(results, &dataloader.Result{Data: &models.KeyGetResponse{}, Error: fmt.Errorf(connutils.StaticKeyNotFound)})
			}
		}

		return results
	}

	// create Loader with an in-memory cache
	thingsLoader := dataloader.NewBatchedLoader(
		thingBatchFunc,
		dataloader.WithWait(50*time.Millisecond),
		dataloader.WithBatchCapacity(100),
	)
	ctx = context.WithValue(ctx, thingsDataLoader, thingsLoader)

	// create Loader with an in-memory cache
	actionsLoader := dataloader.NewBatchedLoader(
		actionBatchFunc,
		dataloader.WithWait(50*time.Millisecond),
		dataloader.WithBatchCapacity(100),
	)
	ctx = context.WithValue(ctx, actionsDataLoader, actionsLoader)

	// create Loader with an in-memory cache
	keysLoader := dataloader.NewBatchedLoader(
		keyBatchFunc,
		dataloader.WithWait(50*time.Millisecond),
		dataloader.WithBatchCapacity(100),
	)
	ctx = context.WithValue(ctx, keysDataLoader, keysLoader)

	return f.databaseConnector.Attach(ctx)
}

// SetServerAddress function
func (f *DataLoader) SetServerAddress(serverAddress string) {
	f.databasConneector.SetServerAddress(serverAddress)
}

// SetConfig function
func (f *DataLoader) SetConfig(configInput *config.Environment) error {
	return f.databaseConnector.SetConfig(configInput)
}

// SetMessaging is used to fill the messaging object
func (f *DataLoader) SetMessaging(m *messages.Messaging) error {
	f.messaging = m
	f.databaseConnector.SetMessaging(m)

	return nil
}

// SetSchema function
func (f *DataLoader) SetSchema(schemaInput *schema.WeaviateSchema) error {
	return f.databaseConnector.SetSchema(schemaInput)
}
