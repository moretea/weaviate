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
)

// ComputeCentroid returns a computed centroid or an error
func ComputeCentroid(vectors []Vector) (*Vector, error) {
	weights := make([]float32, len(vectors))

	for i := 0; i < len(vectors); i++ {
		weights[i] = 1.0
	}

	return ComputeWeightedCentroid(vectors, weights)
}

// ComputeWeightedCentroid returns a weighted computed centroid or an error
func ComputeWeightedCentroid(vectors []Vector, weights []float32) (*Vector, error) {

	if len(vectors) == 0 {
		return nil, fmt.Errorf("can not compute centroid of empty slice")
	} else if len(vectors) != len(weights) {
		return nil, fmt.Errorf("can not compute weighted centroid if len(vectors) != len(weights)")
	} else if len(vectors) == 1 {
		return &vectors[0], nil
	} else {
		vectorLen := vectors[0].Len()

		newVector := make([]float32, vectorLen)
		var weightSum float32

		for vectorI, v := range vectors {
			if v.Len() != vectorLen {
				return nil, fmt.Errorf("vectors have different lengths")
			}

			weightSum += weights[vectorI]

			for i := 0; i < vectorLen; i++ {
				newVector[i] += v.vector[i] * weights[vectorI]
			}
		}

		for i := 0; i < vectorLen; i++ {
			newVector[i] /= weightSum
		}

		result := NewVector(newVector)
		return &result, nil
	}
}
