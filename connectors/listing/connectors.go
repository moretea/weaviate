package listing

import (
	dbconnector "github.com/creativesoftwarefdn/weaviate/connectors"
	"github.com/creativesoftwarefdn/weaviate/connectors/foobar"
	"github.com/creativesoftwarefdn/weaviate/connectors/janusgraph"
)

// GetAllConnectors contains all available connectors
func GetAllConnectors() []dbconnector.DatabaseConnector {
	// Set all existing connectors
	connectors := []dbconnector.DatabaseConnector{
		&foobar.Foobar{},
		&janusgraph.Janusgraph{},
	}

	return connectors
}
