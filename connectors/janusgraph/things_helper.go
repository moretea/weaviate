package janusgraph

import (
	//	"strings"

	"github.com/creativesoftwarefdn/weaviate/gremlin"
	"github.com/creativesoftwarefdn/weaviate/models"

	"github.com/go-openapi/strfmt"
)

func (j *Janusgraph) fillThingResponseFromVertex(vertex *gremlin.Vertex, thingResponse *models.ThingGetResponse) error {
	thingResponse.ThingID = strfmt.UUID(vertex.AssertProperty("uuid").AssertString())
	thingResponse.AtClass = vertex.AssertProperty("atClass").AssertString()
	thingResponse.AtContext = vertex.AssertProperty("context").AssertString()

	thingResponse.CreationTimeUnix = vertex.AssertProperty("creationTimeUnix").AssertInt64()
	thingResponse.LastUpdateTimeUnix = vertex.AssertProperty("lastUpdateTimeUnix").AssertInt64()

	thingResponse.Schema = make(map[string]interface{})
	//thingResponse.Key = models.SingleRef {
	//  NrDollarCref: key
	//}
	//  j.createCrefObject(keyUUID, j.serverAddress, connutils.RefTypeKey)

	//keyResponse.KeyExpiresUnix = vertex.AssertProperty("keyExpiresUnix").AssertInt64()
	//keyResponse.Write = vertex.AssertProperty("write").AssertBool()
	//keyResponse.Email = vertex.AssertProperty("email").AssertString()
	//keyResponse.Read = vertex.AssertProperty("read").AssertBool()
	//keyResponse.Delete = vertex.AssertProperty("delete").AssertBool()
	//keyResponse.Execute = vertex.AssertProperty("execute").AssertBool()
	//keyResponse.IPOrigin = strings.Split(vertex.AssertProperty("IPOrigin").AssertString(), ";")

	//isRoot := vertex.AssertProperty("isRoot").AssertBool()
	//keyResponse.IsRoot = &isRoot

	return nil
}
