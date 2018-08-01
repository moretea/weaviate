package gremlin

import (
  "strings"
  "fmt"
)

// Escape a string so that it can be used without risk of SQL-injection like escapes.
func escapeString(str string) string {
  return strings.Replace(str, `"`, `\"`, -1)
}

func extend_query(query *Query, format string, vals ...interface{}) *Query {
  r := Query { query: query.query + fmt.Sprintf(format, vals...) }
  return &r
}
