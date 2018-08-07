package janusgraph

import (
	"context"
	"fmt"

	"github.com/graphql-go/graphql"
)

// GetLocalGraph handles the Local GraphQL response for this connector
func (f *Janusgraph) GetLocalGraph(ctx context.Context, request graphql.ResolveParams) (interface{}, error) {
	return nil, fmt.Errorf("Not supported")
}
