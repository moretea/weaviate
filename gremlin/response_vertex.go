package gremlin

import (
	"fmt"
)

type Vertex struct {
	Id         int
	Label      string
	Properties map[string]Property
}

func (v *Vertex) AssertProperty(name string) *Property {
	prop := v.Property(name)

	if prop == nil {
		panic(fmt.Sprintf("Expected to find a property '%v' on vertex '%v', but no such property exists!", name, v.Id))
	}

	return prop
}

func (v *Vertex) Property(name string) *Property {
	val, ok := v.Properties[name]
	if !ok {
		return nil
	} else {
		return &val
	}
}
