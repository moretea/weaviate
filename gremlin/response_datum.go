package gremlin

import (
	"errors"
	"fmt"
	"reflect"
)

// A single piece of data returned by a Gremlin query.
type Datum struct {
	Datum interface{}
}

// Extract the type of a datum. This method is intended to be a debugging tool,
// don't use it in a critical path.
func (d *Datum) Type() string {
	switch d.Datum.(type) {
	case int:
		return "int"
	default:
		return fmt.Sprintf("Unknown type for '%#v'", reflect.TypeOf(d.Datum).Name())
	}
}

// Attempt to extract a Vertex from this datum.
func (d *Datum) Vertex() (*Vertex, error) {
	v, ok := d.Datum.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Expected a vertex, but got something else as the result")
	}

	type_, ok := v["type"]
	if !ok || type_ != "vertex" {
		return nil, errors.New("Vertex element in result does not have type 'vertex'")
	}

	label_interface, ok := v["label"]
	if !ok {
		return nil, errors.New("Vertex element does not have a label")
	}
	label, ok := label_interface.(string)
	if !ok {
		return nil, errors.New("Vertex element does not have a string as a label")
	}

	id_interface, ok := v["id"]
	if !ok {
		return nil, errors.New("Vertex element does not have a id")
	}
	id_f, ok := id_interface.(float64)
	if !ok {
		return nil, errors.New("Vertex element does not have a string as a id")
	}
	id := int(id_f)

	properties_interface, ok := v["properties"]
	if !ok {
		return nil, errors.New("Vertex element does not have a properties key")
	}
	properties_map, ok := properties_interface.(map[string]interface{})
	if !ok {
		return nil, errors.New("Vertex element does not have an object for properties ")
	}

	properties, err := extractProperties(properties_map)
	if err != nil {
		return nil, err
	}

	vertex := Vertex{
		Id:         id,
		Label:      label,
		Properties: properties,
	}

	return &vertex, nil
}

func (d *Datum) AssertVertex() *Vertex {
	v, err := d.Vertex()

	if err != nil {
		panic(fmt.Sprintf("Expected datum to be an Vertex, but it was '%s'", d.Type()))
	}

	return v
}
