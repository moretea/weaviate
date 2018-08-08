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
	"errors"
	"reflect"
	"time"

	"github.com/go-openapi/strfmt"

	connutils "github.com/creativesoftwarefdn/weaviate/connectors/utils"
	"github.com/creativesoftwarefdn/weaviate/models"

	"github.com/creativesoftwarefdn/weaviate/gremlin"
)

func (f *Janusgraph) AddThing(ctx context.Context, thing *models.Thing, UUID strfmt.UUID) error {
	// Base settings
	q := gremlin.G.AddV(THING_LABEL).
		As("newThing").
		StringProperty("uuid", string(UUID)).
		StringProperty("atClass", thing.AtClass).
		StringProperty("context", thing.AtContext).
		Int64Property("creationTimeUnix", thing.CreationTimeUnix).
		Int64Property("lastUpdateTimeUnix", thing.LastUpdateTimeUnix)

	type edgeToAdd struct {
		PropertyName string
		Type         string
		Reference    string
		Location     string
	}

	var edgesToAdd []edgeToAdd

	schema, schema_ok := thing.Schema.(map[string]interface{})
	if schema_ok {
		for key, value := range schema {
			janusgraphPropertyName := "schema__" + key
			switch t := value.(type) {
			case string:
				q = q.StringProperty(janusgraphPropertyName, t)
			case int:
				q = q.Int64Property(janusgraphPropertyName, int64(t))
			case int8:
				q = q.Int64Property(janusgraphPropertyName, int64(t))
			case int16:
				q = q.Int64Property(janusgraphPropertyName, int64(t))
			case int32:
				q = q.Int64Property(janusgraphPropertyName, int64(t))
			case int64:
				q = q.Int64Property(janusgraphPropertyName, t)
			case bool:
				q = q.BoolProperty(janusgraphPropertyName, t)
			case float32:
				q = q.Float64Property(janusgraphPropertyName, float64(t))
			case float64:
				q = q.Float64Property(janusgraphPropertyName, t)
			case time.Time:
				q = q.StringProperty(janusgraphPropertyName, time.Time.String(t))
			case *models.SingleRef:
				// Postpone creation of edges
				edgesToAdd = append(edgesToAdd, edgeToAdd{
					PropertyName: janusgraphPropertyName,
					Reference:    t.NrDollarCref.String(),
					Type:         t.Type,
					Location:     *t.LocationURL,
				})
			default:
				f.messaging.ExitError(78, "The type "+reflect.TypeOf(value).String()+" is not supported for Thing properties.")
			}
		}
	}

	// Add edges to all referened things.
	for _, edge := range edgesToAdd {
		q = q.AddE("thingEdge").
			FromRef("newThing").
			ToQuery(gremlin.G.V().HasLabel(THING_LABEL).HasString("uuid", edge.Reference)).
			StringProperty(PROPERTY_EDGE_LABEL, edge.PropertyName).
			StringProperty("$cref", edge.Reference).
			StringProperty("type", edge.Type).
			StringProperty("locationUrl", edge.Location)
	}

	// Link to key
	q = q.AddE(KEY_LABEL).
		FromRef("newThing").
		ToQuery(gremlin.G.V().HasLabel(KEY_LABEL).HasString("uuid", thing.Key.NrDollarCref.String()))

	_, err := f.client.Execute(q)

	return err
}

func (f *Janusgraph) GetThing(ctx context.Context, UUID strfmt.UUID, thingResponse *models.ThingGetResponse) error {
	// Fetch the thing, it's key, and it's relations.
	q := gremlin.G.V().
		HasLabel(THING_LABEL).
		HasString("uuid", string(UUID)).
		As("thing").
		OutWithLabel(KEY_LABEL).As("key").
		V().
		HasLabel(THING_LABEL).
		HasString("uuid", string(UUID)).
		OutEWithLabel("thingEdge").As("ref").
		Select([]string{"thing", "key", "ref"})

	result, err := f.client.Execute(q)

	if err != nil {
		return err
	}

	if len(result.Data) == 0 {
		return errors.New("No thing found")
	}

	// We should have at least two pieces of data;
	// 1)   The Thing
	// 2)   The Key
	// 3-N) Outgoing edges
	if len(result.Data) < 2 {
		return errors.New("No thing (or it's linked key) found!")
	}

	thing_vertex := result.Data[0].AssertVertex()
	key_vertex := result.Data[1].AssertVertex()

	var edges []*gremlin.Edge

	if len(result.Data) > 2 {
		edges = result.AssertEdgeSlice(2, len(result.Data))
	}

	fillKeySingleRefFromVertex(key_vertex, thingResponse.Key)
	return f.fillThingResponseFromVertexAndEdges(thing_vertex, edges, thingResponse)
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
