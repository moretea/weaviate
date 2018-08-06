/*                          _       _
 *__      _____  __ ___   ___  __ _| |_ ___
 *\ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
 * \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
 *  \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
 *
 * Copyright © 2016 - 2018 Weaviate. All rights reserved.
 * LICENSE: https://github.com/creativesoftwarefdn/weaviate/blob/develop/LICENSE.md
 * AUTHOR: Bob van Luijt (bob@kub.design)
 * See www.creativesoftwarefdn.org for details
 * Contact: @CreativeSofwFdn / bob@kub.design
 */

package contextionary

import (
	"fmt"
	"math"
)

// Vector is an opque type that models a fixed-length vector.
type Vector struct {
	vector []float32
}

// NewVector returns a new Vector struct
func NewVector(vector []float32) Vector {
	return Vector{vector}
}

// Equal returns if vectors are equal in the form of a boolean value or an error
func (v *Vector) Equal(other *Vector) (bool, error) {
	if len(v.vector) != len(other.vector) {
		return false, fmt.Errorf("vectors have different dimensions; %v vs %v", len(v.vector), len(other.vector))
	}

	for i, v := range v.vector {
		if other.vector[i] != v {
			return false, nil
		}
	}

	return true, nil
}

// EqualEpsilon returns if vectors have an equal epsilon in the form of a boolean value or an error
func (v *Vector) EqualEpsilon(other *Vector, epsilon float32) (bool, error) {
	if len(v.vector) != len(other.vector) {
		return false, fmt.Errorf("vectors have different dimensions; %v vs %v", len(v.vector), len(other.vector))
	}

	for i, v := range v.vector {
		vMin := v - epsilon
		vMax := v + epsilon
		if other.vector[i] < vMin && other.vector[i] > vMax {
			return false, nil
		}
	}

	return true, nil
}

// Len returns the lenth of a vector
func (v *Vector) Len() int {
	return len(v.vector)
}

// ToString returns the vector as a string
func (v *Vector) ToString() string {
	str := "["
	first := true
	for _, i := range v.vector {
		if first {
			first = false
		} else {
			str += ", "
		}

		str += fmt.Sprintf("%.3f", i)
	}

	str += "]"

	return str
}

// Distance returns a float32 of the vector distance or an error
func (v *Vector) Distance(other *Vector) (float32, error) {
	var sum float32

	if len(v.vector) != len(other.vector) {
		return 0.0, fmt.Errorf("vectors have different dimensions")
	}

	for i := 0; i < len(v.vector); i++ {
		x := v.vector[i] - other.vector[i]
		sum += x * x
	}

	return float32(math.Sqrt(float64(sum))), nil
}
