package janusgraph

import (
	"strings"

	"github.com/creativesoftwarefdn/weaviate/gremlin"
	"github.com/creativesoftwarefdn/weaviate/models"

	"github.com/go-openapi/strfmt"
)

func fillKeyResponseFromVertex(vertex *gremlin.Vertex, keyResponse *models.KeyGetResponse) {
	keyResponse.KeyID = strfmt.UUID(vertex.AssertProperty("uuid").AssertString())
	keyResponse.KeyExpiresUnix = vertex.AssertProperty("keyExpiresUnix").AssertInt64()
	keyResponse.Write = vertex.AssertProperty("write").AssertBool()
	keyResponse.Email = vertex.AssertProperty("email").AssertString()
	keyResponse.Read = vertex.AssertProperty("read").AssertBool()
	keyResponse.Delete = vertex.AssertProperty("delete").AssertBool()
	keyResponse.Execute = vertex.AssertProperty("execute").AssertBool()
	keyResponse.IPOrigin = strings.Split(vertex.AssertProperty("IPOrigin").AssertString(), ";")

	isRoot := vertex.AssertProperty("isRoot").AssertBool()
	keyResponse.IsRoot = &isRoot
}

func fillKeySingleRefFromVertex(vertex *gremlin.Vertex, keyRef *models.SingleRef) {
	// TODO
}
