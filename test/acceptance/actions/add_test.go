package test

// Acceptance tests for actions.
import (
	//"fmt"
	"testing"

	//	"sort"
	"time"

	//	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/stretchr/testify/assert"

	"github.com/creativesoftwarefdn/weaviate/client/actions"
	"github.com/creativesoftwarefdn/weaviate/client/things"
	"github.com/creativesoftwarefdn/weaviate/models"
	"github.com/creativesoftwarefdn/weaviate/test/acceptance/helper"
	//"github.com/creativesoftwarefdn/weaviate/validation"
	//connutils "github.com/creativesoftwarefdn/weaviate/connectors/utils"
)

const fakeActionId strfmt.UUID = "11111111-1111-1111-1111-111111111111"

var localLocation string = "http://localhost:8080"

// Create a thing
func createThing(t *testing.T) strfmt.UUID {
	params := things.NewWeaviateThingsCreateParams().WithBody(&models.ThingCreate{
		AtContext: "http://example.org",
		AtClass:   "TestThing",
		Schema:    map[string]interface{}{},
	})

	resp, err := helper.Client(t).Things.WeaviateThingsCreate(params, helper.RootAuth)

	var uuid strfmt.UUID

	// Ensure that the response is OK
	helper.AssertRequestOk(t, resp, err, func() {
		thing := resp.Payload
		uuid = thing.ThingID
	})

	// Now we poll until success.
	for {
		params := things.NewWeaviateThingsGetParams().WithThingID(uuid)
		_, err := helper.Client(t).Things.WeaviateThingsGet(params, helper.RootAuth)
		if err == nil {
			break
		} else {
			t.Log("Getting the thing failed; sleeping for a bit")
			time.Sleep(10 * time.Millisecond)
		}
	}

	return uuid
}

// Check if we can create a action, and that it's properties are stored correctly.
// Then also check that we can get this action with an external HTTP request.
func TestCreateAndGetAction(t *testing.T) {
	t.Parallel()

	subjectId := createThing(t)
	objectId := createThing(t)

	actionTestString := "Test string"
	actionTestInt := 1
	actionTestBoolean := true
	actionTestNumber := 1.337
	actionTestDate := "2017-10-06T08:15:30+01:00"

	params := actions.NewWeaviateActionsCreateParams().WithBody(&models.ActionCreate{
		AtContext: "http://example.org",
		AtClass:   "TestAction",
		Schema: map[string]interface{}{
			"testString":   actionTestString,
			"testInt":      actionTestInt,
			"testBoolean":  actionTestBoolean,
			"testNumber":   actionTestNumber,
			"testDateTime": actionTestDate,
		},
		Things: &models.ObjectSubject{
			Object: &models.SingleRef{
				NrDollarCref: subjectId,
				Type:         "Thing",
			},
			Subject: &models.SingleRef{
				NrDollarCref: objectId,
				Type:         "Thing",
			},
		},
	})

	resp, err := helper.Client(t).Actions.WeaviateActionsCreate(params, helper.RootAuth)

	var action *models.ActionGetResponse

	// Ensure that the response is OK
	helper.AssertRequestOk(t, resp, err, func() {
		action = resp.Payload
		assert.Regexp(t, strfmt.UUIDPattern, action.ActionID)

		schema, ok := action.Schema.(map[string]interface{})
		if !ok {
			t.Fatal("The returned schema is not an JSON object")
		}

		// Check whether the returned information is the same as the data added
		assert.Equal(t, actionTestString, schema["testString"])
		assert.Equal(t, actionTestInt, int(schema["testInt"].(float64)))
		assert.Equal(t, actionTestBoolean, schema["testBoolean"])
		assert.Equal(t, actionTestNumber, schema["testNumber"])
		assert.Equal(t, actionTestDate, schema["testDateTime"])
	})

	// Now perform an HTTP request, and that the correct values have been stored.
	params2 := actions.NewWeaviateActionsGetParams().WithActionID(action.ActionID)
	resp2, err := helper.Client(t).Actions.WeaviateActionsGet(params2, helper.RootAuth)

	helper.AssertRequestOk(t, resp2, err, func() {
		action2 := resp2.Payload
		assert.Regexp(t, strfmt.UUIDPattern, action2.ActionID)

		schema, ok := action2.Schema.(map[string]interface{})
		if !ok {
			t.Fatal("The returned schema is not an JSON object")
		}

		// Check whether the returned information is the same as the data added
		assert.Equal(t, actionTestString, schema["testString"])
		assert.Equal(t, actionTestInt, int(schema["testInt"].(float64)))
		assert.Equal(t, actionTestBoolean, schema["testBoolean"])
		assert.Equal(t, actionTestNumber, schema["testNumber"])

		// TODO: Janus recognizes dates and modifies them.
		//assert.Equal(t, actionTestDate, schema["testDateTime"])
	})
}
