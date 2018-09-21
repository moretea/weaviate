package dataloader

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/go-openapi/strfmt"

	"github.com/creativesoftwarefdn/weaviate/models"

	connutils "github.com/creativesoftwarefdn/weaviate/connectors/utils"
)

// ValidateToken function
func (f *DataLoader) ValidateToken(ctx context.Context, UUID strfmt.UUID, keyResponse *models.KeyGetResponse) (token string, err error) {
	defer f.messaging.TimeTrack(time.Now())

	token, err = f.databaseConnector.ValidateToken(ctx, UUID, keyResponse)

	return token, err
}
