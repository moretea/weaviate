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
 */package contextionary

import (
	"testing"
)

func TestComputeCentroid(t *testing.T) {

	assertCentroidEqual := func(points []Vector, expected Vector) {
		centroid, err := ComputeCentroid(points)
		if err != nil {
			t.Errorf("Could not compute centroid of %v", points)
		}
		equal, err := centroid.Equal(&expected)
		if err != nil {
			t.Errorf("Could not compare centroid with expected vector; %v", err)
		}

		if !equal {
			pointsStr := "{"
			first := true

			for _, point := range points {
				if first {
					first = false
				} else {
					pointsStr += ", "
				}

				pointsStr += point.ToString()
			}
			pointsStr += "}"

			t.Errorf("centroid of %v should be %v but was %v", pointsStr, expected.ToString(), centroid.ToString())
		}
	}

	assertWeightedCentroidEqual := func(points []Vector, weights []float32, expected Vector) {
		centroid, err := ComputeWeightedCentroid(points, weights)
		if err != nil {
			t.Errorf("Could not compute centroid of %v", points)
		}
		equal, err := centroid.EqualEpsilon(&expected, 0.01)
		if err != nil {
			t.Errorf("Could not compare centroid with expected vector; %v", err)
		}

		if !equal {
			pointsStr := "{"
			first := true

			for _, point := range points {
				if first {
					first = false
				} else {
					pointsStr += ", "
				}

				pointsStr += point.ToString()
			}
			pointsStr += "}"

			t.Errorf("centroid of %v should be %v but was %v", pointsStr, expected.ToString(), centroid.ToString())
		}
	}

	va := NewVector([]float32{1, 1, 1})
	vb := NewVector([]float32{0, 0, 0})
	vc := NewVector([]float32{-1, -1, -1})

	assertCentroidEqual([]Vector{va, vb}, NewVector([]float32{0.5, 0.5, 0.5}))
	assertCentroidEqual([]Vector{va, vb, vc}, NewVector([]float32{0.0, 0.0, 0.0}))

	assertWeightedCentroidEqual([]Vector{va, vb}, []float32{1, 0}, va)
	assertWeightedCentroidEqual([]Vector{va, vb}, []float32{0, 1}, vb)
	assertWeightedCentroidEqual([]Vector{va, vb}, []float32{0.66, 0.33}, NewVector([]float32{0.66, 0.66, 0.66}))
}
