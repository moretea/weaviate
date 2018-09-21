package dataloader

import (
	"context"

	"github.com/go-openapi/strfmt"

	"github.com/creativesoftwarefdn/weaviate/gremlin"
	"github.com/creativesoftwarefdn/weaviate/models"

	"encoding/base64"
	"fmt"
	"strings"
)

// AddKey function
func (f *DataLoader) AddKey(ctx context.Context, key *models.Key, UUID strfmt.UUID, token string) error {
	defer f.messaging.TimeTrack(time.Now())

	return f.databaseConnector.AddKey(ctx, key, UUID, token)
}

// GetKey function
func (f *DataLoader) GetKey(ctx context.Context, UUID strfmt.UUID, keyResponse *models.KeyGetResponse) error {
	defer f.messaging.TimeTrack(time.Now(), fmt.Sprintf("DataLoader#GetKey: '%s'", UUID))

	// Init variables used by the dataloader
	var result interface{}
	var loader *dataloader.Loader
	var ok bool

	// Load the dataloader from the context
	if loader, ok = ctx.Value(keysDataLoader).(*dataloader.Loader); !ok {
		return fmt.Errorf("dataloader not found in context")
	}

	// Use a thunk function to load the data based on the dataloader
	thunk := loader.Load(ctx, dataloader.StringKey(string(UUID)))
	result, err := thunk()

	fmt.Println(result)
	fmt.Println(keyResponse)

	// Fill the Key values retrieved from the thunk function.
	if err == nil {
		keyResponse.Key = result.(*models.KeyGetResponse).Key
		keyResponse.KeyID = result.(*models.KeyGetResponse).KeyID
	}

	return err
}

// GetKeys function
func (f *DataLoader) GetKeys(ctx context.Context, UUIDs []strfmt.UUID, keysResponse *[]*models.KeyGetResponse) error {
	defer f.messaging.TimeTrack(time.Now())

	return f.databaseConnector.GetKeys(ctx, UUIDs, keysResponse)
}

// DeleteKey function
func (f *DataLoader) DeleteKey(ctx context.Context, key *models.Key, UUID strfmt.UUID) error {
	defer f.messaging.TimeTrack(time.Now())

	// Init variables used by the dataloader
	var loader *dataloader.Loader
	var ok bool

	// Load the dataloader from the context
	if loader, ok = ctx.Value(keysDataLoader).(*dataloader.Loader); !ok {
		return fmt.Errorf("dataloader not found in context")
	}

	// Clear the data from the Thing-dataloader cache
	loader.Clear(ctx, dataloader.StringKey(string(UUID)))

	return f.databaseConnector.DeleteKey(ctx, key, UUID)
}

// GetKeyChildren function
func (f *DataLoader) GetKeyChildren(ctx context.Context, UUID strfmt.UUID, children *[]*models.KeyGetResponse) error {
	defer f.messaging.TimeTrack(time.Now())

	return f.databaseConnector.GetKeyChildren(ctx, UUID, children)
}

// UpdateKey updates the Key in the DB at the given UUID.
func (f *DataLoader) UpdateKey(ctx context.Context, key *models.Key, UUID strfmt.UUID, token string) error {
	defer f.messaging.TimeTrack(time.Now())

	// Init variables used by the dataloader
	var loader *dataloader.Loader
	var ok bool

	// Load the dataloader from the context
	if loader, ok = ctx.Value(keysDataLoader).(*dataloader.Loader); !ok {
		return fmt.Errorf("dataloader not found in context")
	}

	// Clear the data from the Thing-dataloader cache
	loader.Clear(ctx, dataloader.StringKey(string(UUID)))

	return f.databaseConnector.UpdateKey(ctx, key, UUID, token)
}
