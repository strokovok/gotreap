package gotreap

import (
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func staticRand() func() int {
	return rand.New(rand.NewPCG(528, 491)).Int
}

// mustValues returns a sorted slice of all values stored in the treap.
func mustValues[T any](t *Treap[T]) []T {
	var res []T
	var walk func(*Node[T])
	walk = func(cur *Node[T]) {
		if cur == nil {
			return
		}
		walk(cur.left)
		res = append(res, cur.value)
		walk(cur.right)
	}
	walk(t.root)
	if res == nil {
		return []T{}
	}
	return res
}

func requireTreapValues[T any](t *testing.T, tr *Treap[T], expected ...T) {
	t.Helper()
	require.Equal(t, expected, mustValues(tr))
}

func TestNewTreapAndInsertions(t *testing.T) {
	tr := NewAutoOrderTreapWithRand(staticRand(), 5, 1, 3, 5)
	requireTreapValues(t, tr, 1, 3, 5, 5)

	idx := tr.InsertLeft(5)
	require.Equal(t, 2, idx)
	idx = tr.InsertRight(5)
	require.Equal(t, 5, idx)

	requireTreapValues(t, tr, 1, 3, 5, 5, 5, 5)
}

func TestNewTreapCustomLess(t *testing.T) {
	reverse := func(a, b int) bool { return a > b }
	tr := NewTreapWithRand(reverse, staticRand(), 1, 2, 3, 4)

	requireTreapValues(t, tr, 4, 3, 2, 1)

	require.Equal(t, 1, tr.InsertLeft(3))
	require.Equal(t, 3, tr.InsertRight(3))

	requireTreapValues(t, tr, 4, 3, 3, 3, 2, 1)
}

func TestEraseVariants(t *testing.T) {
	tr := NewAutoOrderTreapWithRand(staticRand(), 1, 2, 2, 2, 3, 3, 4)

	require.Equal(t, 2, tr.EraseAll(3))
	requireTreapValues(t, tr, 1, 2, 2, 2, 4)

	require.Equal(t, 2, tr.EraseLeftmost(2, 2))
	require.Equal(t, 1, tr.EraseRightmost(2, 1))
	require.Equal(t, 1, tr.EraseLeftmost(1, -1))

	requireTreapValues(t, tr, 4)
}

func TestEraseRangeAndPanics(t *testing.T) {
	tr := NewAutoOrderTreapWithRand(staticRand(), 1, 2, 3, 4, 5)

	require.Equal(t, 2, tr.EraseRange(2, true, 4, false))
	requireTreapValues(t, tr, 1, 4, 5)

	tr.InsertRight(3)
	tr.InsertRight(3)
	require.Equal(t, 2, tr.EraseRange(3, true, 3, true))

	require.Panics(t, func() { tr.EraseRange(5, true, 4, true) })
	require.Panics(t, func() { tr.EraseRange(5, true, 5, false) })
	require.Panics(t, func() { tr.EraseRange(5, false, 5, true) })
}

func TestEraseRangeExclusiveAndOutOfBounds(t *testing.T) {
	tr := NewAutoOrderTreapWithRand(staticRand(), 1, 2, 3, 4, 5, 6, 7)

	require.Equal(t, 3, tr.EraseRange(2, false, 6, false))

	requireTreapValues(t, tr, 1, 2, 6, 7)

	require.Equal(t, 0, tr.EraseRange(100, true, 200, true))
}

func TestEraseLeftmostZeroCount(t *testing.T) {
	tr := NewAutoOrderTreapWithRand(staticRand(), 1, 1, 1, 2, 2)

	require.Equal(t, 0, tr.EraseLeftmost(1, 0))
	requireTreapValues(t, tr, 1, 1, 1, 2, 2)

	require.Equal(t, 3, tr.EraseLeftmost(1, 10))
	requireTreapValues(t, tr, 2, 2)
}

func TestEraseRightmostNegativeRemovesAll(t *testing.T) {
	tr := NewAutoOrderTreapWithRand(staticRand(), 3, 3, 3, 4, 5)

	require.Equal(t, 3, tr.EraseRightmost(3, -1))
	requireTreapValues(t, tr, 4, 5)

	require.Equal(t, 0, tr.EraseRightmost(5, 0))
	requireTreapValues(t, tr, 4, 5)
}

func TestEraseAt(t *testing.T) {
	tr := NewAutoOrderTreapWithRand(staticRand(), 1, 2, 3, 4, 5)
	require.Equal(t, 3, tr.EraseAt(1, 3))
	requireTreapValues(t, tr, 1, 5)
	require.Panics(t, func() { tr.EraseAt(-1, 1) })
	require.Panics(t, func() { tr.EraseAt(0, -1) })
}

func TestEraseAtZeroOrExcessCount(t *testing.T) {
	tr := NewAutoOrderTreapWithRand(staticRand(), 1, 2, 3, 4)
	require.Equal(t, 0, tr.EraseAt(2, 0))

	require.Equal(t, 3, tr.EraseAt(1, 10))

	requireTreapValues(t, tr, 1)
}

func TestBoundsAndLookup(t *testing.T) {
	tr := NewAutoOrderTreapWithRand(staticRand(), 1, 2, 2, 3, 4, 5)

	node, idx := tr.FindLowerBound(2)
	require.NotNil(t, node)
	require.Equal(t, 2, node.value)
	require.Equal(t, 1, idx)
	node, idx = tr.FindUpperBound(2)
	require.NotNil(t, node)
	require.Equal(t, 2, node.value)
	require.Equal(t, 2, idx)
	node = tr.At(4)
	require.NotNil(t, node)
	require.Equal(t, 4, node.value)
	require.Nil(t, tr.At(100))
}

func TestBoundsAndLookupEdgeCases(t *testing.T) {
	tr := NewAutoOrderTreapWithRand(staticRand(), 1, 3, 5)

	node, idx := tr.FindLowerBound(6)
	require.Nil(t, node)
	require.Equal(t, 0, idx)

	node, idx = tr.FindUpperBound(0)
	require.Nil(t, node)
	require.Equal(t, 0, idx)

	require.Nil(t, tr.At(10))
	require.Nil(t, tr.At(-10))
	require.Nil(t, tr.At(3))
	require.Nil(t, tr.At(-4))

	node = tr.At(-3)
	require.NotNil(t, node)
	require.Equal(t, 1, node.value)

	node = tr.At(-1)
	require.NotNil(t, node)
	require.Equal(t, 5, node.value)

	node = tr.At(0)
	require.NotNil(t, node)
	require.Equal(t, 1, node.value)

	node = tr.At(1)
	require.NotNil(t, node)
	require.Equal(t, 3, node.value)

	node = tr.At(2)
	require.NotNil(t, node)
	require.Equal(t, 5, node.value)
}

func TestSizeEmptyClear(t *testing.T) {
	tr := NewAutoOrderTreapWithRand[int](staticRand())
	require.True(t, tr.Empty())
	require.Zero(t, tr.Size())
	tr.InsertRight(1)
	tr.InsertRight(2)
	require.Equal(t, 2, tr.Size())
	require.False(t, tr.Empty())
	tr.Clear()
	require.True(t, tr.Empty())
}

func TestClearAllowsReuse(t *testing.T) {
	tr := NewAutoOrderTreapWithRand(staticRand(), 5, 6, 7)
	tr.Clear()
	require.True(t, tr.Empty())

	require.Equal(t, 0, tr.InsertLeft(10))
	require.Equal(t, 1, tr.Size())

	requireTreapValues(t, tr, 10)
}

func TestExtremaAndPops(t *testing.T) {
	tr := NewAutoOrderTreapWithRand(staticRand(), 3, 1, 4, 1, 5)
	require.Equal(t, 1, tr.Leftmost().value)
	require.Equal(t, 5, tr.Rightmost().value)

	val, ok := tr.PopLeftmost()
	require.True(t, ok)
	require.Equal(t, 1, val)
	val, ok = tr.PopRightmost()
	require.True(t, ok)
	require.Equal(t, 5, val)

	empty := NewAutoOrderTreapWithRand[int](staticRand())
	_, ok = empty.PopLeftmost()
	require.False(t, ok)
	_, ok = empty.PopRightmost()
	require.False(t, ok)
}

func TestExtremaAndPopsOnSingleton(t *testing.T) {
	tr := NewAutoOrderTreapWithRand(staticRand(), 42)

	node := tr.Leftmost()
	require.NotNil(t, node)
	require.Equal(t, 42, node.value)
	node = tr.Rightmost()
	require.NotNil(t, node)
	require.Equal(t, 42, node.value)

	val, ok := tr.PopLeftmost()
	require.True(t, ok)
	require.Equal(t, 42, val)
	require.True(t, tr.Empty())
	_, ok = tr.PopRightmost()
	require.False(t, ok)
}

func TestPopUntilEmptyMaintainsOrder(t *testing.T) {
	tr := NewAutoOrderTreapWithRand(staticRand(), 4, 1, 3, 2)

	var popped []int
	for !tr.Empty() {
		val, ok := tr.PopLeftmost()
		require.True(t, ok)
		popped = append(popped, val)
	}

	require.Equal(t, []int{1, 2, 3, 4}, popped)
	require.True(t, tr.Empty())

	_, ok := tr.PopRightmost()
	require.False(t, ok)
}

func TestSplitVariants(t *testing.T) {
	tr := NewAutoOrderTreapWithRand(staticRand(), 1, 2, 3, 4, 5)

	left, right := tr.SplitBefore(3)
	requireTreapValues(t, left, 1, 2)
	requireTreapValues(t, right, 3, 4, 5)

	tr = NewAutoOrderTreapWithRand(staticRand(), 1, 2, 3, 4, 5)
	left, right = tr.SplitAfter(3)
	requireTreapValues(t, left, 1, 2, 3)
	requireTreapValues(t, right, 4, 5)

	tr = NewAutoOrderTreapWithRand(staticRand(), 1, 2, 3, 4, 5)
	left, right = tr.Cut(3)
	requireTreapValues(t, left, 1, 2, 3)
	requireTreapValues(t, right, 4, 5)
}

func TestSplitBeforeAndAfterBoundaries(t *testing.T) {
	tr := NewAutoOrderTreapWithRand(staticRand(), 2, 4, 6)
	left, right := tr.SplitBefore(1)
	require.True(t, left.Empty())
	requireTreapValues(t, right, 2, 4, 6)

	tr = NewAutoOrderTreapWithRand(staticRand(), 2, 4, 6)
	left, right = tr.SplitBefore(10)
	requireTreapValues(t, left, 2, 4, 6)
	require.True(t, right.Empty())

	tr = NewAutoOrderTreapWithRand(staticRand(), 2, 4, 6)
	left, right = tr.SplitAfter(1)
	require.True(t, left.Empty())
	requireTreapValues(t, right, 2, 4, 6)

	tr = NewAutoOrderTreapWithRand(staticRand(), 2, 4, 6)
	left, right = tr.SplitAfter(10)
	requireTreapValues(t, left, 2, 4, 6)
	require.True(t, right.Empty())
}

func TestSplitClearsOriginal(t *testing.T) {
	tr := NewAutoOrderTreapWithRand(staticRand(), 1, 2, 3, 4)
	left, right := tr.SplitBefore(3)

	require.True(t, tr.Empty())
	requireTreapValues(t, left, 1, 2)
	requireTreapValues(t, right, 3, 4)

	require.NotNil(t, left.lessFn)
	require.NotNil(t, right.lessFn)
}

func TestSplitWithDuplicateBoundaries(t *testing.T) {
	tr := NewAutoOrderTreapWithRand(staticRand(), 1, 2, 2, 2, 3, 4)

	left, right := tr.SplitBefore(2)
	requireTreapValues(t, left, 1)
	requireTreapValues(t, right, 2, 2, 2, 3, 4)

	tr = NewAutoOrderTreapWithRand(staticRand(), 1, 2, 2, 2, 3, 4)
	left, right = tr.SplitAfter(2)
	requireTreapValues(t, left, 1, 2, 2, 2)
	requireTreapValues(t, right, 3, 4)
}

func TestCountRangeAndCount(t *testing.T) {
	tr := NewAutoOrderTreapWithRand(staticRand(), 1, 2, 2, 3, 4, 4, 5)

	require.Equal(t, 5, tr.CountRange(2, true, 4, true))
	require.Equal(t, 1, tr.CountRange(2, false, 4, false))
	require.Equal(t, 0, tr.CountRange(6, true, 9, true))

	require.Panics(t, func() { tr.CountRange(5, true, 4, true) })
	require.Panics(t, func() { tr.CountRange(5, true, 5, false) })
	require.Panics(t, func() { tr.CountRange(5, false, 5, true) })

	require.Equal(t, 2, tr.Count(4))
}

func TestCountRangeOutsideBounds(t *testing.T) {
	tr := NewAutoOrderTreapWithRand(staticRand(), 5, 10, 15)

	require.Zero(t, tr.CountRange(-10, true, 0, true))
	require.Zero(t, tr.CountRange(20, true, 30, true))
	require.Zero(t, tr.Count(7))
}

func TestCountRangeExclusiveEliminatesBoundaries(t *testing.T) {
	tr := NewAutoOrderTreapWithRand(staticRand(), 1, 2, 3, 4)

	require.Equal(t, 0, tr.CountRange(1, false, 2, false))
	require.Equal(t, 3, tr.CountRange(1, false, 4, true))
	require.Equal(t, 3, tr.CountRange(1, true, 4, false))
}

func TestCutEdgeCounts(t *testing.T) {
	tr := NewAutoOrderTreapWithRand(staticRand(), 1, 2, 3)
	left, right := tr.Cut(0)
	require.True(t, left.Empty())
	requireTreapValues(t, right, 1, 2, 3)

	tr = NewAutoOrderTreapWithRand(staticRand(), 1, 2, 3)
	left, right = tr.Cut(10)
	requireTreapValues(t, left, 1, 2, 3)
	require.True(t, right.Empty())

	tr = NewAutoOrderTreapWithRand(staticRand(), 1, 2, 3)
	left, right = tr.Cut(-5)
	require.True(t, left.Empty())
	requireTreapValues(t, right, 1, 2, 3)
}

func TestMergeTreap(t *testing.T) {
	left := NewAutoOrderTreapWithRand(staticRand(), 1, 2, 3)
	right := NewAutoOrderTreapWithRand(staticRand(), 4, 5, 6)
	merged := Merge(left, right)

	requireTreapValues(t, merged, 1, 2, 3, 4, 5, 6)

	require.True(t, merged.lessFn(1, 2))
	require.False(t, merged.lessFn(2, 1))
}

func TestMergeTreapSupportsFurtherInsertion(t *testing.T) {
	left := NewAutoOrderTreapWithRand(staticRand(), 1, 3)
	right := NewAutoOrderTreapWithRand(staticRand(), 5, 7)

	merged := Merge(left, right)
	idx := merged.InsertRight(4)
	require.Equal(t, 2, idx)
	idx = merged.InsertLeft(6)
	require.Equal(t, 4, idx)

	requireTreapValues(t, merged, 1, 3, 4, 5, 6, 7)
}

func TestMergeTreapWithEmpty(t *testing.T) {
	left := NewAutoOrderTreapWithRand(staticRand(), 1, 2, 3)
	empty := NewAutoOrderTreapWithRand[int](staticRand())

	merged := Merge(left, empty)
	requireTreapValues(t, merged, 1, 2, 3)

	merged = Merge(empty, left)
	requireTreapValues(t, merged, 1, 2, 3)

	require.NotNil(t, merged.lessFn)
}

func TestJumpRight(t *testing.T) {
	arr := []int{3255, 0, 12}
	tr := NewAutoOrderTreapWithRand(staticRand(), arr...)

	rnd := rand.New(rand.NewPCG(8800, 5553535))
	for range 1000 {
		newValue := rnd.Int()
		tr.InsertRight(newValue)
		arr = append(arr, newValue)
		slices.Sort(arr)

		assert.Equal(t, len(arr), tr.Size())

		pos1 := rnd.IntN(len(arr))
		node1 := tr.At(pos1)
		if assert.NotNil(t, node1) {
			assert.Equal(t, arr[pos1], node1.Value())
		}

		assert.Same(t, node1, node1.JumpRight(0))
		assert.Same(t, node1.Next(), node1.JumpRight(1))
		assert.Same(t, node1.Prev(), node1.JumpRight(-1))

		pos2 := rnd.IntN(len(arr))
		node2 := node1.JumpRight(pos2 - pos1)
		assert.Same(t, tr.At(pos2), node2)
		if assert.NotNil(t, node2) {
			assert.Equal(t, arr[pos2], node2.value)
		}

		for _, pos3 := range []int{len(arr) * -10, len(arr) * -1, -10, -1, len(arr), len(arr) + 10, len(arr) * 2, len(arr) * 10} {
			node3 := node1.JumpRight(pos3 - pos1)
			assert.Nil(t, node3)
		}
	}
}

// TODO: fuzzing
