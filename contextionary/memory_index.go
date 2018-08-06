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
	"fmt"
	"sort"

	annoy "github.com/creativesoftwarefdn/weaviate/contextionary/annoyindex"
)

// MemoryIndex contains the dimensions, words and knn's
type MemoryIndex struct {
	dimensions int
	words      []string
	knn        annoy.AnnoyIndex
}

// GetNumberOfItems returns the number of items that is stored in the index.
func (mi *MemoryIndex) GetNumberOfItems() int {
	return len(mi.words)
}

// GetVectorLength returns the length of the used vectors.
func (mi *MemoryIndex) GetVectorLength() int {
	return mi.dimensions
}

// WordToItemIndex looks up a word, return an index.
// Perform binary search.
func (mi *MemoryIndex) WordToItemIndex(word string) ItemIndex {
	for idx, w := range mi.words {
		if word == w {
			return ItemIndex(idx)
		}
	}

	return -1
}

// ItemIndexToWord returns, based on an index, the assosiated word.
func (mi *MemoryIndex) ItemIndexToWord(item ItemIndex) (string, error) {
	if item >= 0 && int(item) <= len(mi.words) {
		return mi.words[item], nil
	}
	return "", fmt.Errorf("index out of bounds")
}

// GetVectorForItemIndex gets the vector of an item index.
func (mi *MemoryIndex) GetVectorForItemIndex(item ItemIndex) (*Vector, error) {
	if item >= 0 && int(item) <= len(mi.words) {
		var floats []float32
		mi.knn.GetItem(int(item), &floats)

		return &Vector{floats}, nil
	}
	return nil, fmt.Errorf("index out of bounds")
}

// GetDistance computes the distance between two items.
func (mi MemoryIndex) GetDistance(a ItemIndex, b ItemIndex) (float32, error) {
	if a >= 0 && b >= 0 && int(a) <= len(mi.words) && int(b) <= len(mi.words) {
		return mi.knn.GetDistance(int(a), int(b)), nil
	}
	return 0, fmt.Errorf("index out of bounds")
}

// GetNnsByItem gets the n nearest neighbours of item, examining k trees.
// Returns an array of indices, and of distances between item and the n-nearest neighbors.
func (mi *MemoryIndex) GetNnsByItem(item ItemIndex, n int, k int) ([]ItemIndex, []float32, error) {
	if item >= 0 && int(item) <= len(mi.words) {
		var items []int
		var distances []float32

		mi.knn.GetNnsByItem(int(item), n, k, &items, &distances)

		indices := make([]ItemIndex, len(items))
		for i, x := range items {
			indices[i] = ItemIndex(x)
		}

		return indices, distances, nil
	}
	return nil, nil, fmt.Errorf("index out of bounds")
}

// GetNnsByVector gets the n nearest neighbours of item, examining k trees.
// Returns an array of indices, and of distances between item and the n-nearest neighbors.
func (mi *MemoryIndex) GetNnsByVector(vector Vector, n int, k int) ([]ItemIndex, []float32, error) {
	if len(vector.vector) == mi.dimensions {
		var items []int
		var distances []float32

		mi.knn.GetNnsByVector(vector.vector, n, k, &items, &distances)

		indices := make([]ItemIndex, len(items))
		for i, x := range items {
			indices[i] = ItemIndex(x)
		}

		return indices, distances, nil
	}
	return nil, nil, fmt.Errorf("wrong vector length provided")
}

// The rest of this file concerns itself with building the Memory Index.
// This is done from the MemoryIndexBuilder struct.

// MemoryIndexBuilder contains dimensions and word vectors
type MemoryIndexBuilder struct {
	dimensions  int
	wordVectors mibPairs
}

type mibPair struct {
	word   string
	vector Vector
}

// Define custom type, and implement functions required for sort.Sort.
type mibPairs []mibPair

func (a mibPairs) Len() int           { return len(a) }
func (a mibPairs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a mibPairs) Less(i, j int) bool { return a[i].word < a[j].word }

// InMemoryBuilder constructs a new builder.
func InMemoryBuilder(dimensions int) *MemoryIndexBuilder {
	mib := MemoryIndexBuilder{
		dimensions:  dimensions,
		wordVectors: make([]mibPair, 0),
	}

	return &mib
}

// AddWord adds a word and it's vector to the builder.
func (mib *MemoryIndexBuilder) AddWord(word string, vector Vector) {
	wv := mibPair{word: word, vector: vector}
	mib.wordVectors = append(mib.wordVectors, wv)
}

// Build an efficient lookup iddex from the builder.
func (mib *MemoryIndexBuilder) Build(trees int) *MemoryIndex {
	mi := MemoryIndex{
		dimensions: mib.dimensions,
		words:      make([]string, 0),
		knn:        annoy.NewAnnoyIndexEuclidean(mib.dimensions),
	}

	// First sort the words; this way we can do binary search on the words.
	sort.Sort(mib.wordVectors)

	// Then fill up the data in the MemoryIndex
	for i, pair := range mib.wordVectors {
		mi.words = append(mi.words, pair.word)
		mi.knn.AddItem(i, pair.vector.vector)
	}

	// And instruct Annoy to build it's index
	mi.knn.Build(trees)

	return &mi
}
