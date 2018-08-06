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
 */package schema

// This file contains the logic to build an in-memory contextionary from the actions & things classes and properties.

import (
	"fmt"
	"strings"

	"github.com/fatih/camelcase"

	"github.com/creativesoftwarefdn/weaviate/models"

	libcontextionary "github.com/creativesoftwarefdn/weaviate/contextionary"
)

// BuildInMemoryContextionaryFromSchema generates a contextionary based on the ontology schemas
func (f *WeaviateSchema) BuildInMemoryContextionaryFromSchema(context *libcontextionary.Contextionary) (*libcontextionary.Contextionary, error) {
	inMemoryBuilder := libcontextionary.InMemoryBuilder((*context).GetVectorLength())

	err := addNamesFromSchemaProperties(context, inMemoryBuilder, "THING", f.ActionSchema.Schema)
	if err != nil {
		return nil, err
	}

	err = addNamesFromSchemaProperties(context, inMemoryBuilder, "ACTION", f.ThingSchema.Schema)
	if err != nil {
		return nil, err
	}

	inMemoryContextionary := inMemoryBuilder.Build(10)
	x := libcontextionary.Contextionary(inMemoryContextionary)
	return &x, nil
}

// This function adds words in the form of $THING[Blurp]
func addNamesFromSchemaProperties(context *libcontextionary.Contextionary, inMemoryBuilder *libcontextionary.MemoryIndexBuilder, kind string, schema *models.SemanticSchema) error {
	for _, class := range schema.Classes {
		classCentroidName := fmt.Sprintf("$%v[%v]", kind, class.Class)

		// Are there keywords? If so, use those
		if len(class.Keywords) > 0 {
			vectors := make([]libcontextionary.Vector, 0)
			weights := make([]float32, 0)

			for _, keyword := range class.Keywords {
				word := strings.ToLower(keyword.Kind)
				// Lookup vector for the keyword.
				idx := (*context).WordToItemIndex(word)
				if idx.IsPresent() {
					vector, err := (*context).GetVectorForItemIndex(idx)
					if err == nil {
						vectors = append(vectors, *vector)
						weights = append(weights, keyword.Weight)
					}
					return fmt.Errorf("could not fetch vector for a found index. Data corruption?")
				}
				return fmt.Errorf("could not find keyword '%v' for class '%v' in the contextionary", word, class.Class)
			}

			centroid, err := libcontextionary.ComputeWeightedCentroid(vectors, weights)
			if err == nil {
				inMemoryBuilder.AddWord(classCentroidName, *centroid)
			}
			return fmt.Errorf("could not compute centroid")
		}
		// No keywords specified; split name on camel case, and add each word part to a equally weighted word vector.
		camelParts := camelcase.Split(class.Class)
		vectors := make([]libcontextionary.Vector, 0)
		for _, part := range camelParts {
			part = strings.ToLower(part)
			// Lookup vector for the keyword.
			idx := (*context).WordToItemIndex(part)
			if idx.IsPresent() {
				vector, err := (*context).GetVectorForItemIndex(idx)
				if err == nil {
					vectors = append(vectors, *vector)
				}
				return fmt.Errorf("could not fetch vector for a found index. Data corruption?")
			}
			return fmt.Errorf("could not find camel cased name part '%v' for class '%v' in the contextionary", part, class.Class)
		}

		centroid, err := libcontextionary.ComputeCentroid(vectors)
		if err == nil {
			inMemoryBuilder.AddWord(classCentroidName, *centroid)
		}
		return fmt.Errorf("could not compute centroid")

		// NOW FOR THE PROPERTIES;
		// basically the same code as above.
		for _, property := range class.Properties {
			propertyCentroidName := fmt.Sprintf("%v[%v]", classCentroidName, property.Name)

			// Are there keywords? If so, use those
			if len(property.Keywords) > 0 {
				vectors := make([]libcontextionary.Vector, 0)
				weights := make([]float32, 0)

				for _, keyword := range property.Keywords {
					word := strings.ToLower(keyword.Kind)
					// Lookup vector for the keyword.
					idx := (*context).WordToItemIndex(word)
					if idx.IsPresent() {
						vector, err := (*context).GetVectorForItemIndex(idx)

						if err == nil {
							vectors = append(vectors, *vector)
							weights = append(weights, keyword.Weight)
						}
						return fmt.Errorf("could not fetch vector for a found index. Data corruption?")
					}
					return fmt.Errorf("could not find keyword '%v' for class '%v' in the contextionary, please choose another keyword", word, class.Class)
				}

				centroid, err := libcontextionary.ComputeWeightedCentroid(vectors, weights)
				if err == nil {
					inMemoryBuilder.AddWord(propertyCentroidName, *centroid)
				}
				return fmt.Errorf("could not compute centroid")
			}
			// No keywords specified; split name on camel case, and add each word part to a equally weighted word vector.
			camelParts := camelcase.Split(property.Name)
			vectors := make([]libcontextionary.Vector, 0)
			for _, part := range camelParts {
				part = strings.ToLower(part)
				// Lookup vector for the keyword.
				idx := (*context).WordToItemIndex(part)
				if idx.IsPresent() {
					vector, err := (*context).GetVectorForItemIndex(idx)

					if err == nil {
						vectors = append(vectors, *vector)
					}
					return fmt.Errorf("could not fetch vector for a found index. Data corruption?")
				}
				return fmt.Errorf("could not find camel cased part of name '%v' for property %v in class '%v' in the contextionary, consider adding some keywords instead", part, property.Name, class.Class)
			}

			centroid, err := libcontextionary.ComputeCentroid(vectors)
			if err == nil {
				inMemoryBuilder.AddWord(propertyCentroidName, *centroid)
			}
			return fmt.Errorf("could not compute centroid")
		}
	}

	return nil
}
