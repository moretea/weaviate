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
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/creativesoftwarefdn/weaviate/connectors/utils"
	"github.com/creativesoftwarefdn/weaviate/models"

	"github.com/creativesoftwarefdn/weaviate/gremlin"
)

func (f *Janusgraph) AddAction(ctx context.Context, action *models.Action, UUID strfmt.UUID) error {
	// Base settings
	q := gremlin.G.AddV(ACTION_LABEL).
		As("newAction").
		StringProperty("uuid", string(UUID)).
		StringProperty("atClass", action.AtClass).
		StringProperty("context", action.AtContext).
		Int64Property("creationTimeUnix", action.CreationTimeUnix).
		Int64Property("lastUpdateTimeUnix", action.LastUpdateTimeUnix)

	type edgeToAdd struct {
		PropertyName string
		Type         string
		Reference    string
		Location     string
	}

	var edgesToAdd []edgeToAdd

	schema, schema_ok := action.Schema.(map[string]interface{})
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
				f.messaging.ExitError(78, "The type "+reflect.TypeOf(value).String()+" is not supported for Action properties.")
			}
		}
	}

	// Add edges to all referened actions.
	for _, edge := range edgesToAdd {
		q = q.AddE("actionEdge").
			FromRef("newAction").
			ToQuery(gremlin.G.V().HasLabel(ACTION_LABEL).HasString("uuid", edge.Reference)).
			StringProperty(PROPERTY_EDGE_LABEL, edge.PropertyName).
			StringProperty("$cref", edge.Reference).
			StringProperty("type", edge.Type).
			StringProperty("locationUrl", edge.Location)
	}

	// Link to key
	q = q.AddE(KEY_LABEL).
		StringProperty("locationUrl", *action.Key.LocationURL).
		FromRef("newAction").
		ToQuery(gremlin.G.V().HasLabel(KEY_LABEL).HasString("uuid", action.Key.NrDollarCref.String()))

	// Link to Subject
	var subjectLabel string
	switch action.Things.Object.Type {
	case "Thing":
		subjectLabel = THING_LABEL
	case "Action":
		subjectLabel = ACTION_LABEL
	case "Key":
		subjectLabel = KEY_LABEL
	}
	q = q.AddE(ACTION_SUBJECT_LABEL).
		FromRef("newAction").
		ToQuery(gremlin.G.V().HasLabel(subjectLabel).HasString("uuid", string(action.Things.Subject.NrDollarCref)))

	// Link to Object
	var objectLabel string
	switch action.Things.Object.Type {
	case "Thing":
		objectLabel = THING_LABEL
	case "Action":
		objectLabel = ACTION_LABEL
	case "Key":
		objectLabel = KEY_LABEL
	}
	q = q.AddE(ACTION_OBJECT_LABEL).
		FromRef("newAction").
		ToQuery(gremlin.G.V().HasLabel(objectLabel).HasString("uuid", string(action.Things.Object.NrDollarCref)))

	_, err := f.client.Execute(q)

	return err
}

func (f *Janusgraph) GetAction(ctx context.Context, UUID strfmt.UUID, actionResponse *models.ActionGetResponse) error {
	// Fetch the action, it's key, and it's relations.
	q := gremlin.G.V().
		HasLabel(ACTION_LABEL).
		HasString("uuid", string(UUID)).
		As("action").
		OutEWithLabel(ACTION_SUBJECT_LABEL).InV().As("subject").
		V().
		HasLabel(ACTION_LABEL).
		HasString("uuid", string(UUID)).
		OutEWithLabel(ACTION_OBJECT_LABEL).InV().As("object").
		V().
		HasLabel(ACTION_LABEL).
		HasString("uuid", string(UUID)).
		OutEWithLabel(KEY_LABEL).As("keyEdge").
		InV().Path().FromRef("keyEdge").As("key"). // also get the path, so that we can learn about the location of the key.
		V().
		HasLabel(ACTION_LABEL).
		HasString("uuid", string(UUID)).
		Raw(`.optional(outE("actionEdge").as("actionEdge").as("ref")).choose(select("ref"), select("action", "object", "subject", "key", "ref"), select("action", "object", "subject", "key"))`)

	result, err := f.client.Execute(q)

	if err != nil {
		return err
	}

	if len(result.Data) == 0 {
		return errors.New(connutils.StaticActionNotFound)
	}

	// The outputs 'action' and 'key' will be repeated over all results. Just get them for one for now.
	actionVertex := result.Data[0].AssertKey("action").AssertVertex()
	objectVertex := result.Data[0].AssertKey("object").AssertVertex()
	subjectVertex := result.Data[0].AssertKey("subject").AssertVertex()
	keyPath := result.Data[0].AssertKey("key").AssertPath()

	// However, we can get multiple refs. In that case, we'll have multiple datums,
	// each with the same action & key, but a different ref.
	// Let's extract those refs.
	var refEdges []*gremlin.Edge
	for _, datum := range result.Data {
		ref, err := datum.Key("ref")
		if err == nil {
			refEdges = append(refEdges, ref.AssertEdge())
		}
	}

	actionResponse.Key = newKeySingleRefFromKeyPath(keyPath)
	return fillActionResponseFromVertexAndEdges(actionVertex, subjectVertex, objectVertex, refEdges, actionResponse)
}

func (f *Janusgraph) GetActions(ctx context.Context, UUIDs []strfmt.UUID, response *models.ActionsListResponse) error {
	// TODO: Optimize query to perform just _one_ JanusGraph lookup.

	response.TotalResults = 0
	response.Actions = make([]*models.ActionGetResponse, 0)

	for _, uuid := range UUIDs {
		var actionResponse models.ActionGetResponse
		err := f.GetAction(ctx, uuid, &actionResponse)

		if err == nil {
			response.TotalResults += 1
			response.Actions = append(response.Actions, &actionResponse)
		} else {
			return fmt.Errorf("%s: action with UUID '%v' not found", connutils.StaticActionNotFound, uuid)
		}
	}

	return nil
}

func (f *Janusgraph) ListActions(ctx context.Context, UUID strfmt.UUID, first int, offset int, wheres []*connutils.WhereQuery, response *models.ActionsListResponse) error {
	//TODO rewrite to one query.
	if len(wheres) > 0 {
		return errors.New("Wheres are not supported in ListActions")
	}

	q := gremlin.G.V().
		HasLabel(ACTION_LABEL).
		Range(offset, first).
		Values([]string{"uuid"})

	result, err := f.client.Execute(q)

	if err != nil {
		return err
	}

	response.TotalResults = 0
	response.Actions = make([]*models.ActionGetResponse, 0)

	// Get the UUIDs from the first query.
	UUIDs := result.AssertStringSlice()

	for _, uuid := range UUIDs {
		var actionResponse models.ActionGetResponse
		err := f.GetAction(ctx, strfmt.UUID(uuid), &actionResponse)

		if err == nil {
			response.TotalResults += 1
			response.Actions = append(response.Actions, &actionResponse)
		} else {
			// skip silently; it's probably deleted.
		}
	}

	return nil
}

func (f *Janusgraph) UpdateAction(ctx context.Context, action *models.Action, UUID strfmt.UUID) error {
	// Base settings
	q := gremlin.G.V().HasLabel(ACTION_LABEL).
		HasString("uuid", string(UUID)).
		As("action").
		StringProperty("atClass", action.AtClass).
		StringProperty("context", action.AtContext).
		Int64Property("creationTimeUnix", action.CreationTimeUnix).
		Int64Property("lastUpdateTimeUnix", action.LastUpdateTimeUnix)

	type expectedEdge struct {
		PropertyName string
		Type         string
		Reference    string
		Location     string
	}

	var expectedEdges []expectedEdge

	schema, schema_ok := action.Schema.(map[string]interface{})
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
				expectedEdges = append(expectedEdges, expectedEdge{
					PropertyName: janusgraphPropertyName,
					Reference:    t.NrDollarCref.String(),
					Type:         t.Type,
					Location:     *t.LocationURL,
				})
			default:
				f.messaging.ExitError(78, "The type "+reflect.TypeOf(value).String()+" is not supported for Action properties.")
			}
		}
	}

	// Update all edges to all referened actions.
	// TODO: verify what to if we're not mentioning some reference? how should we remove such a reference?
	for _, edge := range expectedEdges {
		// First drop the edge
		q = q.Optional(gremlin.Current().OutEWithLabel("actionEdge").HasString(PROPERTY_EDGE_LABEL, edge.PropertyName).Drop()).
			AddE("actionEdge").
			FromRef("action").
			ToQuery(gremlin.G.V().HasLabel(ACTION_LABEL).HasString("uuid", edge.Reference)).
			StringProperty(PROPERTY_EDGE_LABEL, edge.PropertyName).
			StringProperty("$cref", edge.Reference).
			StringProperty("type", edge.Type).
			StringProperty("locationUrl", edge.Location)
	}

	// Don't update the key.
	// TODO verify that indeed this is the desired behaviour.

	// Don't update the subject & object
	// TODO verify that indeed this is the desired behaviour.

	_, err := f.client.Execute(q)

	return err
}

func (f *Janusgraph) DeleteAction(ctx context.Context, action *models.Action, UUID strfmt.UUID) error {
	q := gremlin.G.V().HasLabel(ACTION_LABEL).
		HasString("uuid", string(UUID)).
		Drop()

	_, err := f.client.Execute(q)

	return err
}

func (f *Janusgraph) HistoryAction(ctx context.Context, UUID strfmt.UUID, history *models.ActionHistory) error {
	return nil
}

func (f *Janusgraph) MoveToHistoryAction(ctx context.Context, action *models.Action, UUID strfmt.UUID, deleted bool) error {
	return nil
}

func fillActionResponseFromVertexAndEdges(vertex *gremlin.Vertex, subject *gremlin.Vertex, object *gremlin.Vertex, refEdges []*gremlin.Edge, actionResponse *models.ActionGetResponse) error {
	// TODO: We should actually read stuff from the database schema, then get only that stuff from JanusGraph.
	// At this moment, we're just parsing whetever there is in JanusGraph, which might not agree with the database schema
	// that is defined in Weaviate.

	actionResponse.ActionID = strfmt.UUID(vertex.AssertPropertyValue("uuid").AssertString())
	actionResponse.AtClass = vertex.AssertPropertyValue("atClass").AssertString()
	actionResponse.AtContext = vertex.AssertPropertyValue("context").AssertString()

	actionResponse.CreationTimeUnix = vertex.AssertPropertyValue("creationTimeUnix").AssertInt64()
	actionResponse.LastUpdateTimeUnix = vertex.AssertPropertyValue("lastUpdateTimeUnix").AssertInt64()

	// TODO: convert this into helper function.
	var objectType string
	var subjectType string

	switch object.Label {
	case KEY_LABEL:
		objectType = "Key"
	case THING_LABEL:
		objectType = "Thing"
	case ACTION_LABEL:
		objectType = "Action"
	}

	switch subject.Label {
	case KEY_LABEL:
		subjectType = "Key"
	case THING_LABEL:
		subjectType = "Thing"
	case ACTION_LABEL:
		subjectType = "Action"
	}

	actionResponse.Things = &models.ObjectSubject{
		Object: &models.SingleRef{
			NrDollarCref: strfmt.UUID(object.AssertPropertyValue("uuid").AssertString()),
			Type:         objectType,
		},
		Subject: &models.SingleRef{
			NrDollarCref: strfmt.UUID(subject.AssertPropertyValue("uuid").AssertString()),
			Type:         subjectType,
		},
	}

	schema := make(map[string]interface{})

	// Walk through all properties, check if they start with 'schema__', and then consider them to be 'schema' properties.
	// Just copy in the value directly. We're not doing any sanity check/casting to proper types for now.
	for key, val := range vertex.Properties {
		if strings.HasPrefix(key, "schema__") {
			key = key[8:len(key)]
			schema[key] = val.Value.Value
		}
	}

	// For each of the connected edges, get the property values,
	// and store the reference.
	for _, edge := range refEdges {
		locationUrl := edge.AssertPropertyValue("locationUrl").AssertString()
		type_ := edge.AssertPropertyValue("type").AssertString()
		edgeName := edge.AssertPropertyValue("propertyEdge").AssertString()
		uuid := edge.AssertPropertyValue("$cref").AssertString()

		key := edgeName[8:len(edgeName)]
		ref := make(map[string]interface{})
		ref["$cref"] = uuid
		ref["locationUrl"] = locationUrl
		ref["type"] = type_
		schema[key] = ref
	}

	actionResponse.Schema = schema

	return nil
}
