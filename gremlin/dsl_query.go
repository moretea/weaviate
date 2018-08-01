package gremlin

// A query represents the (partial) query build with the DSL
type Query struct {
	query string
}

// Return the string representation of this Query.
func (q *Query) Query() string {
	return q.query
}

func RawQuery(query string) *Query {
	return &Query{query: query}
}

// Count how many vertices or edges are selected by the previous query.
func (q *Query) Count() *Query {
	return extend_query(q, ".count()")
}

func (q *Query) AddV(label string) *Query {
	return extend_query(q, `.addV("%s")`, escapeString(label))
}

func (q *Query) AddE(label string) *Query {
	return extend_query(q, `.addE("%s")`, escapeString(label))
}

// Set the expected label of the vertex/edge.
func (q *Query) HasLabel(label string) *Query {
	return extend_query(q, `.hasLabel("%s")`, escapeString(label))
}

func (q *Query) HasString(key string, value string) *Query {
	return extend_query(q, `.has("%s", "%s")`, escapeString(key), escapeString(value))
}

func (q *Query) HasBool(key string, value bool) *Query {
	return extend_query(q, `.has("%s", %v)`, escapeString(key), value)
}

func (q *Query) StringProperty(key string, value string) *Query {
	return extend_query(q, `.property("%s", "%s")`, escapeString(key), escapeString(value))
}

func (q *Query) BoolProperty(key string, value bool) *Query {
	return extend_query(q, `.property("%s", %v)`, escapeString(key), value)
}

func (q *Query) Int64Property(key string, value int64) *Query {
	return extend_query(q, `.property("%s", (long) %v)`, escapeString(key), value)
}

func (q *Query) InE() *Query {
	return extend_query(q, ".inE()")
}

func (q *Query) InEWithLabel(label string) *Query {
	return extend_query(q, `.inE("%s")`, escapeString(label))
}

func (q *Query) OutV() *Query {
	return extend_query(q, ".outV()")
}

// Create a reference
func (q *Query) As(name string) *Query {
	return extend_query(q, `.as("%s")`, escapeString(name))
}

// Point to a reference
func (q *Query) FromRef(reference string) *Query {
	return extend_query(q, `.from("%s")`, escapeString(reference))
}

func (q *Query) ToQuery(query *Query) *Query {
	return extend_query(q, `.to(%s)`, query.Query())
}

func (q *Query) Drop() *Query {
	return extend_query(q, ".drop()")
}
