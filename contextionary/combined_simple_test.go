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

func TestSimpleCombinedIndex(t *testing.T) {
	builder1 := InMemoryBuilder(3)
	builder2 := InMemoryBuilder(3)
	builder3 := InMemoryBuilder(3)

	builder1.AddWord("a", NewVector([]float32{1, 0, 0}))
	builder2.AddWord("b", NewVector([]float32{0, 1, 0}))
	builder3.AddWord("c", NewVector([]float32{0, 0, 1}))

	memoryIndex1 := Contextionary(builder1.Build(3))
	memoryIndex2 := Contextionary(builder2.Build(3))
	memoryIndex3 := Contextionary(builder3.Build(3))

	indices123 := []Contextionary{memoryIndex1, memoryIndex2, memoryIndex3}
	indices231 := []Contextionary{memoryIndex2, memoryIndex3, memoryIndex1}
	indices312 := []Contextionary{memoryIndex3, memoryIndex1, memoryIndex2}

	t.Run("indices 123", func(t *testing.T) { testSimpleCombined(t, indices123) })
	t.Run("indices 231", func(t *testing.T) { testSimpleCombined(t, indices231) })
	t.Run("indices 312", func(t *testing.T) { testSimpleCombined(t, indices312) })
}

func testSimpleCombined(t *testing.T, indices []Contextionary) {
	ci, err := CombineVectorIndices(indices)
	if err != nil {
		panic("should work")
	}

	aIdx := ci.WordToItemIndex("a")
	if !aIdx.IsPresent() {
		panic("should be present")
	}

	bIdx := ci.WordToItemIndex("b")
	if !bIdx.IsPresent() {
		panic("should be present")
	}

	cIdx := ci.WordToItemIndex("c")
	if !cIdx.IsPresent() {
		panic("should be present")
	}

	items, _, err := ci.GetNnsByItem(aIdx, 3, 3)
	if err != nil {
		panic("should work")
	}

	assertEqIdx := func(name string, a, b ItemIndex) {
		if a != b {
			t.Errorf("Expected %v to be at %v, but was at %b", name, a, b)
		}
	}

	if len(items) != 3 {
		t.Errorf("got length %v, expected 3", len(items))
		t.FailNow()
	}

	// assert lexicographical order, if distances are equal

	assertEqIdx("a", aIdx, items[0])
	assertEqIdx("b", bIdx, items[1])
	assertEqIdx("c", cIdx, items[2])
}
