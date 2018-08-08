package gremlin

import (
	"fmt"
)

type Edge struct {
	Id         int
	Label      string
	Properties map[string]Property
}

func (e *Edge) AssertProperty(name string) *Property {
	prop := e.Property(name)

	if prop == nil {
		panic(fmt.Sprintf("Expected to find a property '%v' on edge '%v', but no such property exists!", name, e.Id))
	}

	return prop
}

func (e *Edge) Property(name string) *Property {
	val, ok := e.Properties[name]
	if !ok {
		return nil
	} else {
		return &val
	}
}
