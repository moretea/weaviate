package contextionary

import (
	"testing"
//  "fmt"
)

func TestSimpleCombinedIndex(t *testing.T) {
  builder1 := InMemoryBuilder(3)
  builder2 := InMemoryBuilder(3)
  builder3 := InMemoryBuilder(3)

  builder1.AddWord("a", NewVector([]float32 { 1, 0, 0}))
  builder2.AddWord("b", NewVector([]float32 { 0, 1, 0}))
  builder3.AddWord("c", NewVector([]float32 { 0, 0, 1}))


  memory_index1 := VectorIndex(builder1.Build(3))
  memory_index2 := VectorIndex(builder2.Build(3))
  memory_index3 := VectorIndex(builder3.Build(3))

  var indices123 []VectorIndex = []VectorIndex { memory_index1, memory_index2, memory_index3 }
  var indices231 []VectorIndex = []VectorIndex { memory_index2, memory_index3, memory_index1 }
  var indices312 []VectorIndex = []VectorIndex { memory_index3, memory_index1, memory_index2 }

  t.Run("indices 123", func(t *testing.T) { test_simple_combined(t, indices123) })
  t.Run("indices 231", func(t *testing.T) { test_simple_combined(t, indices231) })
  t.Run("indices 312", func(t *testing.T) { test_simple_combined(t, indices312) })
}

func test_simple_combined(t *testing.T, indices []VectorIndex) {
  ci, err := CombineVectorIndices(indices)
  if err != nil { panic("should work") }

  a_idx := ci.WordToItemIndex("a")
  if !a_idx.IsPresent() { panic("should be present") }

  b_idx  := ci.WordToItemIndex("b")
  if !b_idx.IsPresent() { panic("should be present") }

  c_idx := ci.WordToItemIndex("c")
  if !c_idx.IsPresent() { panic("should be present") }

  items, _, err := ci.GetNnsByItem(a_idx, 3, 3)
  if err != nil { panic("should work") }

  assert_eq_idx := func(name string, a, b ItemIndex) {
    if a != b {
      t.Errorf("Expected %v to be at %v, but was at %b", name, a, b)
    }
  }

  if len(items) != 3 { t.Errorf("got length %v, expected 3", len(items)); t.FailNow() }

  assert_eq_idx("a", a_idx, items[0])
  assert_eq_idx("b", b_idx, items[1])
  assert_eq_idx("c", c_idx, items[2])
}
