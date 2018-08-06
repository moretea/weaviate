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
	"sort"
)

// CombinedIndex combines the indeces and the vectors
type CombinedIndex struct {
	indices      []combinedIndex
	totalSize    int
	vectorLength int
}

type combinedIndex struct {
	offset int
	size   int
	index  *Contextionary
}

// CombineVectorIndices combines multiple indices, present them as one.
// It assumes that each index stores unique words
func CombineVectorIndices(indices []Contextionary) (*CombinedIndex, error) {
	// We join the ItemIndex spaces the indivual indices, by
	// offsetting the 2nd ItemIndex with len(indices[0]),
	// the 3rd ItemIndex space with len(indices[0]) + len(indices[1]), etc.

	if len(indices) < 2 {
		return nil, fmt.Errorf("less than two vector indices provided")
	}

	combinedIndices := make([]combinedIndex, len(indices))

	var offset int

	vectorLength := indices[0].GetVectorLength()

	for i := 0; i < len(indices); i++ {
		size := indices[i].GetNumberOfItems()

		combinedIndices[i] = combinedIndex{
			offset: offset,
			size:   size,
			index:  &indices[i],
		}

		offset += size

		myLength := indices[i].GetVectorLength()

		if myLength != vectorLength {
			return nil, fmt.Errorf("vector length not equal")
		}
	}

	return &CombinedIndex{indices: combinedIndices, totalSize: offset, vectorLength: vectorLength}, nil
}

// VerifyDisjoint verifies that all the indices are disjoint
// Returns nil on success, an error if the words in the indices are not disjoint.
func (ci *CombinedIndex) VerifyDisjoint() error {
	for indexI, itemI := range ci.indices {
		for i := ItemIndex(0); int(i) < itemI.size; i++ {
			word, err := (*itemI.index).ItemIndexToWord(i)
			if err != nil {
				panic("should not happen; this index should always be accessible")
			}

			for indexJ, itemJ := range ci.indices {
				if indexI != indexJ {
					result := (*(itemJ.index)).WordToItemIndex(word)
					if result.IsPresent() {
						return fmt.Errorf("word %v is in more than one index", word)
					}
				}
			}
		}
	}

	return nil
}

// GetNumberOfItems counts the amount of numbers
func (ci *CombinedIndex) GetNumberOfItems() int {
	return ci.totalSize
}

// GetVectorLength returns the length of a vector
func (ci *CombinedIndex) GetVectorLength() int {
	return ci.vectorLength
}

// WordToItemIndex returns the index based in a word
func (ci *CombinedIndex) WordToItemIndex(word string) ItemIndex {
	for _, item := range ci.indices {
		itemIndex := (*item.index).WordToItemIndex(word)

		if (&itemIndex).IsPresent() {
			return itemIndex + ItemIndex(item.offset)
		}
	}

	return -1
}

func (ci *CombinedIndex) findVectorIndexForItemIndex(itemIndex ItemIndex) (ItemIndex, *Contextionary, error) {
	item := int(itemIndex)

	for _, idx := range ci.indices {
		if item >= idx.offset && item < (idx.offset+idx.size) {
			return ItemIndex(item - idx.offset), idx.index, nil
		}
	}

	return 0, nil, fmt.Errorf("out of index")
}

// ItemIndexToWord returns the item index based on a word or an error
func (ci *CombinedIndex) ItemIndexToWord(item ItemIndex) (string, error) {
	offsettedIndex, vi, err := ci.findVectorIndexForItemIndex(item)

	if err != nil {
		return "", err
	}

	word, err := (*vi).ItemIndexToWord(offsettedIndex)
	return word, err
}

// GetVectorForItemIndex returns the vector based on an index or returns an error
func (ci *CombinedIndex) GetVectorForItemIndex(item ItemIndex) (*Vector, error) {
	offsettedIndex, vi, err := ci.findVectorIndexForItemIndex(item)

	if err != nil {
		return nil, err
	}

	word, err := (*vi).GetVectorForItemIndex(offsettedIndex)
	return word, err
}

// GetDistance computes the distance between two items.
func (ci *CombinedIndex) GetDistance(a ItemIndex, b ItemIndex) (float32, error) {
	v1, err := ci.GetVectorForItemIndex(a)
	if err != nil {
		return 0.0, err
	}

	v2, err := ci.GetVectorForItemIndex(b)
	if err != nil {
		return 0.0, err
	}

	dist, err := v1.Distance(v2)
	if err != nil {
		return 0.0, err
	}

	return dist, nil
}

// GetNnsByItem gets the n nearest neighbours of item, examining k trees.
// Returns an array of indices, and of distances between item and the n-nearest neighbors.
func (ci *CombinedIndex) GetNnsByItem(item ItemIndex, n int, k int) ([]ItemIndex, []float32, error) {
	vec, err := ci.GetVectorForItemIndex(item)
	if err != nil {
		return nil, nil, err
	}

	return ci.GetNnsByVector(*vec, n, k)
}

type combinedNnSearchResult struct {
	item ItemIndex
	dist float32
}

type combinedNnSearchResults struct {
	items []combinedNnSearchResult
	ci    *CombinedIndex
}

func (a combinedNnSearchResults) Len() int      { return len(a.items) }
func (a combinedNnSearchResults) Swap(i, j int) { a.items[i], a.items[j] = a.items[j], a.items[i] }
func (a combinedNnSearchResults) Less(i, j int) bool {
	// Sort on distance first, if those are the same, sort on lexographical order of the words.
	if a.items[i].dist == a.items[j].dist {
		wi, err := a.ci.ItemIndexToWord(a.items[i].item)
		if err != nil {
			panic("should be there")
		}

		wj, err := a.ci.ItemIndexToWord(a.items[j].item)
		if err != nil {
			panic("should be there")
		}
		return wi < wj
	}
	return a.items[i].dist < a.items[j].dist
}

// Remove a certain element from the result search.
func (a *combinedNnSearchResults) Remove(i int) {
	a.items = append(a.items[:i], a.items[i+1:]...)
}

// GetNnsByVector gets the n nearest neighbours of item, examining k trees.
// Returns an array of indices, and of distances between item and the n-nearest neighbors.
func (ci *CombinedIndex) GetNnsByVector(vector Vector, n int, k int) ([]ItemIndex, []float32, error) {
	results := combinedNnSearchResults{
		items: make([]combinedNnSearchResult, 0),
		ci:    ci,
	}

	for _, item := range ci.indices {
		indices, floats, err := (*item.index).GetNnsByVector(vector, n, k)

		if err == nil {
			for i, itemIdx := range indices {
				results.items = append(results.items, combinedNnSearchResult{item: itemIdx + ItemIndex(item.offset), dist: floats[i]})
			}
		}
		return nil, nil, err
	}

	sort.Sort(results)

	// Now remove duplicates.
	for i := 1; i < len(results.items); {
		if results.items[i].item == results.items[i-1].item {
			results.Remove(i)
		} else {
			i++ // only increment if we're not removing.
		}
	}

	items := make([]ItemIndex, 0)
	floats := make([]float32, 0)

	var maxIndex int

	if n < len(results.items) {
		maxIndex = n
	} else {
		maxIndex = len(results.items)
	}

	for i := 0; i < maxIndex; i++ {
		items = append(items, results.items[i].item)
		floats = append(floats, results.items[i].dist)
	}

	return items, floats, nil
}
