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

	"github.com/go-openapi/strfmt"

	"github.com/creativesoftwarefdn/weaviate/models"
	connutils "github.com/creativesoftwarefdn/weaviate/connectors/utils"
)

func (f *Janusgraph) AddThing(ctx context.Context, thing *models.Thing, UUID strfmt.UUID) error {
	return nil
}

func (f *Janusgraph) GetThing(ctx context.Context, UUID strfmt.UUID, thingResponse *models.ThingGetResponse) error {
	return nil
}

func (f *Janusgraph) GetThings(ctx context.Context, UUIDs []strfmt.UUID, thingResponse *models.ThingsListResponse) error {
	return nil
}

func (f *Janusgraph) ListThings(ctx context.Context, first int, offset int, keyID strfmt.UUID, wheres []*connutils.WhereQuery, thingsResponse *models.ThingsListResponse) error {
	return nil
}

func (f *Janusgraph) UpdateThing(ctx context.Context, thing *models.Thing, UUID strfmt.UUID) error {
	return nil
}

func (f *Janusgraph) DeleteThing(ctx context.Context, thing *models.Thing, UUID strfmt.UUID) error {
	return nil
}

func (f *Janusgraph) HistoryThing(ctx context.Context, UUID strfmt.UUID, history *models.ThingHistory) error {
	return nil
}

func (f *Janusgraph) MoveToHistoryThing(ctx context.Context, thing *models.Thing, UUID strfmt.UUID, deleted bool) error {
	return nil
}
