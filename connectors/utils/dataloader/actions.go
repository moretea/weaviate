package dataloader

import (
	"context"

	"github.com/go-openapi/strfmt"

	"github.com/creativesoftwarefdn/weaviate/connectors/utils"
	"github.com/creativesoftwarefdn/weaviate/models"
)

// AddAction function
func (f *DataLoader) AddAction(ctx context.Context, action *models.Action, UUID strfmt.UUID) error {
	defer f.messaging.TimeTrack(time.Now())

	return f.databaseConnector.AddAction(ctx, action, UUID)
}

// GetAction function
func (f *DataLoader) GetAction(ctx context.Context, UUID strfmt.UUID, actionResponse *models.ActionGetResponse) error {
	defer f.messaging.TimeTrack(time.Now(), fmt.Sprintf("DataLoader#GetAction: '%s'", UUID))

	// Init variables used by the dataloader
	var result interface{}
	var loader *dataloader.Loader
	var ok bool

	// Load the dataloader from the context
	if loader, ok = ctx.Value(actionsDataLoader).(*dataloader.Loader); !ok {
		return fmt.Errorf("dataloader not found in context")
	}

	// Use a thunk function to load the data based on the dataloader
	thunk := loader.Load(ctx, dataloader.StringKey(string(UUID)))
	result, err := thunk()

	// Fill the Action values retrieved from the thunk function
	if err == nil {
		actionResponse.Action = result.(*models.ActionGetResponse).Action
		actionResponse.ActionID = result.(*models.ActionGetResponse).ActionID
	}

	return err
}

// GetActions function
func (f *DataLoader) GetActions(ctx context.Context, UUIDs []strfmt.UUID, actionsResponse *models.ActionsListResponse) error {
	defer f.messaging.TimeTrack(time.Now())

	return f.databaseConnector.GetActions(ctx, UUIDs, actionsResponse)
}

// ListActions function
func (f *DataLoader) ListActions(ctx context.Context, UUID strfmt.UUID, first int, offset int, wheres []*connutils.WhereQuery, actionsResponse *models.ActionsListResponse) error {
	defer f.messaging.TimeTrack(time.Now())

	return f.databaseConnector.ListActions(ctx, UUID, first, offset, wheres, actionsResponse)
}

// UpdateAction function
func (f *DataLoader) UpdateAction(ctx context.Context, action *models.Action, UUID strfmt.UUID) error {
	defer f.messaging.TimeTrack(time.Now())

	// Init variables used by the dataloader
	var thingsLoader *dataloader.Loader
	var actionsLoader *dataloader.Loader
	var ok bool

	// Load the dataloader from the context
	if thingsLoader, ok = ctx.Value(thingsDataLoader).(*dataloader.Loader); !ok {
		return fmt.Errorf("dataloader not found in context")
	}

	// Clear the data from the Thing-dataloader cache
	thingsLoader.Clear(ctx, dataloader.StringKey(string(action.Things.Subject.NrDollarCref)))
	thingsLoader.Clear(ctx, dataloader.StringKey(string(action.Things.Object.NrDollarCref)))

	// Load the dataloader from the context
	if actionsLoader, ok = ctx.Value(actionsDataLoader).(*dataloader.Loader); !ok {
		return fmt.Errorf("dataloader not found in context")
	}

	// Clear the data from the Thing-dataloader cache
	actionsLoader.Clear(ctx, dataloader.StringKey(string(UUID)))

	// Forward request to db-connector
	return f.databaseConnector.UpdateAction(ctx, action, UUID)
}

// DeleteAction function
func (f *DataLoader) DeleteAction(ctx context.Context, action *models.Action, UUID strfmt.UUID) error {
	defer f.messaging.TimeTrack(time.Now())

	// Init variables used by the dataloader
	var thingsLoader *dataloader.Loader
	var actionsLoader *dataloader.Loader
	var ok bool

	// Load the dataloader from the context
	if thingsLoader, ok = ctx.Value(thingsDataLoader).(*dataloader.Loader); !ok {
		return fmt.Errorf("dataloader not found in context")
	}

	// Clear the data from the Thing-dataloader cache
	thingsLoader.Clear(ctx, dataloader.StringKey(string(action.Things.Subject.NrDollarCref)))
	thingsLoader.Clear(ctx, dataloader.StringKey(string(action.Things.Object.NrDollarCref)))

	// Load the dataloader from the context
	if actionsLoader, ok = ctx.Value(actionsDataLoader).(*dataloader.Loader); !ok {
		return fmt.Errorf("dataloader not found in context")
	}

	// Clear the data from the Thing-dataloader cache
	actionsLoader.Clear(ctx, dataloader.StringKey(string(UUID)))

	// Forward request to db-connector
	return f.databaseConnector.DeleteAction(ctx, action, UUID)
}

// HistoryAction fills the history of an Action based on its UUID
func (f *DataLoader) HistoryAction(ctx context.Context, UUID strfmt.UUID, history *models.ActionHistory) error {
	return f.databaseConnector.HistoryAction(ctx, UUID, history)
}

// MoveToHistoryAction moves an Action to history
func (f *DataLoader) MoveToHistoryAction(ctx context.Context, action *models.Action, UUID strfmt.UUID, deleted bool) error {
	return f.databaseConnector.MoveToHistoryAction(ctx, action, UUID, deleted)
}
