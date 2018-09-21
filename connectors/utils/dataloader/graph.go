package dataloader

import (
	"context"
)

// This is where the GraphQL endpoint resolvers point to
func (f *Janusgraph) GetGraph(ctx context.Context, request interface{}) (string, error) {
	return "{}", nil
}
