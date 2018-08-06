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

	annoy "github.com/creativesoftwarefdn/weaviate/contextionary/annoyindex"
)

type mmappedIndex struct {
	wordIndex *Wordlist
	knn       annoy.AnnoyIndex
}

func (m *mmappedIndex) GetNumberOfItems() int {
	return int(m.wordIndex.numberOfWords)
}

// Returns the length of the used vectors.
func (m *mmappedIndex) GetVectorLength() int {
	return int(m.wordIndex.vectorWidth)
}

func (m *mmappedIndex) WordToItemIndex(word string) ItemIndex {
	return m.wordIndex.FindIndexByWord(word)
}

func (m *mmappedIndex) ItemIndexToWord(item ItemIndex) (string, error) {
	if item >= 0 && item <= m.wordIndex.GetNumberOfWords() {
		return m.wordIndex.getWord(item), nil
	}
	return "", fmt.Errorf("index out of bounds")
}

func (m *mmappedIndex) GetVectorForItemIndex(item ItemIndex) (*Vector, error) {
	if item >= 0 && item <= m.wordIndex.GetNumberOfWords() {
		var floats []float32
		m.knn.GetItem(int(item), &floats)

		return &Vector{floats}, nil
	}
	return nil, fmt.Errorf("index out of bounds")
}

// Compute the distance between two items.
func (m *mmappedIndex) GetDistance(a ItemIndex, b ItemIndex) (float32, error) {
	if a >= 0 && b >= 0 && a <= m.wordIndex.GetNumberOfWords() && b <= m.wordIndex.GetNumberOfWords() {
		return m.knn.GetDistance(int(a), int(b)), nil
	}
	return 0, fmt.Errorf("index out of bounds")
}

func (m *mmappedIndex) GetNnsByItem(item ItemIndex, n int, k int) ([]ItemIndex, []float32, error) {
	if item >= 0 && item <= m.wordIndex.GetNumberOfWords() {
		var items []int
		var distances []float32

		m.knn.GetNnsByItem(int(item), n, k, &items, &distances)

		indices := make([]ItemIndex, len(items))
		for i, x := range items {
			indices[i] = ItemIndex(x)
		}

		return indices, distances, nil
	}
	return nil, nil, fmt.Errorf("index out of bounds")
}

func (m *mmappedIndex) GetNnsByVector(vector Vector, n int, k int) ([]ItemIndex, []float32, error) {
	if len(vector.vector) == m.GetVectorLength() {
		var items []int
		var distances []float32

		m.knn.GetNnsByVector(vector.vector, n, k, &items, &distances)

		indices := make([]ItemIndex, len(items))
		for i, x := range items {
			indices[i] = ItemIndex(x)
		}

		return indices, distances, nil
	}
	return nil, nil, fmt.Errorf("wrong vector length provided")
}

// LoadVectorFromDisk loads the vector file from disk
func LoadVectorFromDisk(annoyIndex string, wordIndexFileName string) (*Contextionary, error) {
	wordIndex, err := LoadWordlist(wordIndexFileName)

	if err != nil {
		return nil, fmt.Errorf("could not load vector: %+v", err)
	}

	knn := annoy.NewAnnoyIndexEuclidean(int(wordIndex.vectorWidth))
	knn.Load(annoyIndex)

	idx := new(mmappedIndex)
	idx.wordIndex = wordIndex
	idx.knn = knn

	contextionary := Contextionary(idx)
	return &contextionary, nil
}
