package dataloader

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/go-openapi/strfmt"

	connutils "github.com/creativesoftwarefdn/weaviate/connectors/utils"
	"github.com/creativesoftwarefdn/weaviate/models"

	"encoding/json"
)

// AddThing function
func (f *DataLoader) AddThing(ctx context.Context, thing *models.Thing, UUID strfmt.UUID) error {
	defer f.messaging.TimeTrack(time.Now())

	return f.databaseConnector.AddThing(ctx, thing, UUID)
}

// GetThing function
func (f *DataLoader) GetThing(ctx context.Context, UUID strfmt.UUID, thingResponse *models.ThingGetResponse) error {
	defer f.messaging.TimeTrack(time.Now(), fmt.Sprintf("DataLoader#GetThing: '%s'", UUID))

	// Init variables used by the dataloader
	var result interface{}
	var loader *dataloader.Loader
	var ok bool

	// Load the dataloader from the context
	if loader, ok = ctx.Value(thingsDataLoader).(*dataloader.Loader); !ok {
		return fmt.Errorf("dataloader not found in context")
	}

	// Use a thunk function to load the data based on the dataloader
	thunk := loader.Load(ctx, dataloader.StringKey(string(UUID)))
	result, err := thunk()

	// Fill the Thing values retrieved from the thunk function.
	if err == nil {
		thingResponse.Thing = result.(*models.ThingGetResponse).Thing
		thingResponse.ThingID = result.(*models.ThingGetResponse).ThingID
	}

	return err
}

// GetThings function
func (f *DataLoader) GetThings(ctx context.Context, UUIDs []strfmt.UUID, thingsResponse *models.ThingsListResponse) error {
	defer f.messaging.TimeTrack(time.Now(), fmt.Sprintf("DataLoader#GetThings: '%s'", UUIDs))
	return f.databaseConnector.GetThings(ctx, UUIDs, thingsResponse)
}

// ListThings function
func (f *DataLoader) ListThings(ctx context.Context, first int, offset int, keyID strfmt.UUID, wheres []*connutils.WhereQuery, thingsResponse *models.ThingsListResponse) error {
	defer f.messaging.TimeTrack(time.Now())

	return f.databaseConnector.ListThings(ctx, first, offset, keyID, wheres, thingsResponse)
}

// UpdateThing function
func (f *DataLoader) UpdateThing(ctx context.Context, thing *models.Thing, UUID strfmt.UUID) error {
	defer f.messaging.TimeTrack(time.Now())

	// Init variables used by the dataloader
	var loader *dataloader.Loader
	var ok bool

	// Load the dataloader from the context
	if loader, ok = ctx.Value(thingsDataLoader).(*dataloader.Loader); !ok {
		return fmt.Errorf("dataloader not found in context")
	}

	// Clear the data from the Thing-dataloader cache
	loader.Clear(ctx, dataloader.StringKey(string(UUID)))

	// Forward request to db-connector
	return f.databaseConnector.UpdateThing(ctx, thing, UUID)
}

// DeleteThing function
func (f *DataLoader) DeleteThing(ctx context.Context, thing *models.Thing, UUID strfmt.UUID) error {
	defer f.messaging.TimeTrack(time.Now())

	// Init variables used by the dataloader
	var loader *dataloader.Loader
	var ok bool

	// Load the dataloader from the context
	if loader, ok = ctx.Value(thingsDataLoader).(*dataloader.Loader); !ok {
		return fmt.Errorf("dataloader not found in context")
	}

	// Clear the data from the Thing-dataloader cache
	loader.Clear(ctx, dataloader.StringKey(string(UUID)))

	// Forward request to db-connector
	return f.databaseConnector.DeleteThing(ctx, thing, UUID)
}

// HistoryThing fills the history of a Thing based on its UUID
func (f *DataLoader) HistoryThing(ctx context.Context, UUID strfmt.UUID, history *models.ThingHistory) error {
	return f.databaseConnector.HistoryThing(ctx, UUID, history)
}

// MoveToHistoryThing moves a Thing to history
func (f *DataLoader) MoveToHistoryThing(ctx context.Context, thing *models.Thing, UUID strfmt.UUID, deleted bool) error {
	return f.databaseConnector.MoveToHistoryThing(ctx, thing, UUID, deleted)
}
