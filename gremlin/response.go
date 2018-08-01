package gremlin

import (
	"errors"
	"fmt"
	"log"
)

type Response struct {
	Data []interface{}
}

func (r *Response) OneInt() (int, error) {
	if len(r.Data) != 1 {
		return 0, fmt.Errorf("Query resulted in %v results, whilst we expected just 1", len(r.Data))
	}

	i, ok := r.Data[0].(float64)
	if !ok {
		return 0, fmt.Errorf("Expected to get see an int, but got %#v instead", r.Data[0])
	}

	return int(i), nil
}

func (r *Response) Vertices() ([]Vertex, error) {
	vertices := make([]Vertex, 0)

	for _, raw_v := range r.Data {
		v, ok := raw_v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Expected a vertex, but got something else as the result")
		} else {
			vertex, err := extractVertex(v)
			if err != nil {
				return nil, err
			} else {
				vertices = append(vertices, *vertex)
			}
		}
	}

	return vertices, nil
}

func extractVertex(v map[string]interface{}) (*Vertex, error) {
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

func extractProperties(props map[string]interface{}) (map[string]Property, error) {
	properties := make(map[string]Property)
	for key, prop_val := range props {
		prop_val_maps, ok := prop_val.([]interface{})

		if !ok {
			return nil, fmt.Errorf("Property is not a list %#v", prop_val)
		}

		if len(prop_val_maps) != 1 {
			log.Fatalf("should be exactly 1, but got %#v", prop_val_maps)
		}

		prop_val_map, ok := prop_val_maps[0].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Property value map is not an object %#v", prop_val)
		}

		prop_val, ok := prop_val_map["value"]
		if !ok {
			return nil, fmt.Errorf("no 'value' in property object")
		}

		prop_id_interface, ok := prop_val_map["id"]
		if !ok {
			return nil, fmt.Errorf("no 'id' in property object")
		}

		prop_id, ok := prop_id_interface.(string)
		if !ok {
			return nil, fmt.Errorf("'id' in property object is not a string")
		}

		property := Property{
			Id:   prop_id,
			Data: prop_val,
		}

		properties[key] = property
	}

	return properties, nil
}

type Property struct {
	Id   string
	Data interface{}
}

func (p *Property) String() (string, bool) {
	val, ok := p.Data.(string)
	return val, ok
}

func (p *Property) AssertString() string {
	val, ok := p.String()
	if ok {
		return val
	} else {
		panic(fmt.Sprintf("Expected a string, but got %#v", p.Data))
	}
}

func (p *Property) Float() (float64, bool) {
	val, ok := p.Data.(float64)
	return val, ok
}

func (p *Property) AssertFloat() float64 {
	val, ok := p.Float()
	if ok {
		return val
	} else {
		panic(fmt.Sprintf("Expected a float, but got %#v", p.Data))
	}
}

func (p *Property) Int() (int, bool) {
	val, ok := p.Data.(float64)
	return int(val), ok
}

func (p *Property) AssertInt() int {
	val, ok := p.Int()
	if ok {
		return val
	} else {
		panic(fmt.Sprintf("Expected a int, but got %#v", p.Data))
	}
}

func (p *Property) Int64() (int64, bool) {
	val, ok := p.Data.(float64)
	return int64(val), ok
}

func (p *Property) AssertInt64() int64 {
	val, ok := p.Int64()
	if ok {
		return val
	} else {
		panic(fmt.Sprintf("Expected a int, but got %#v", p.Data))
	}
}

func (p *Property) Bool() (bool, bool) {
	val, ok := p.Data.(bool)
	return val, ok
}

func (p *Property) AssertBool() bool {
	val, ok := p.Bool()
	if ok {
		return val
	} else {
		panic(fmt.Sprintf("Expected a bool, but got %#v", p.Data))
	}
}

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
