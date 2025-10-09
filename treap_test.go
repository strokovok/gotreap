package util

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

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
	tr := NewAutoOrderTreap(5, 1, 3, 5)
	requireTreapValues(t, tr, 1, 3, 5, 5)

	idx := tr.InsertLeft(5)
	require.Equal(t, 2, idx)
	idx = tr.InsertRight(5)
	require.Equal(t, 5, idx)

	requireTreapValues(t, tr, 1, 3, 5, 5, 5, 5)
}

func TestNewTreapCustomLess(t *testing.T) {
	reverse := func(a, b int) bool { return a > b }
	tr := NewTreap(reverse, 1, 2, 3, 4)

	requireTreapValues(t, tr, 4, 3, 2, 1)

	require.Equal(t, 1, tr.InsertLeft(3))
	require.Equal(t, 3, tr.InsertRight(3))

	requireTreapValues(t, tr, 4, 3, 3, 3, 2, 1)
}

func TestEraseVariants(t *testing.T) {
	tr := NewAutoOrderTreap(1, 2, 2, 2, 3, 3, 4)

	require.Equal(t, 2, tr.EraseAll(3))
	requireTreapValues(t, tr, 1, 2, 2, 2, 4)

	require.Equal(t, 2, tr.EraseLeftmost(2, 2))
	require.Equal(t, 1, tr.EraseRightmost(2, 1))
	require.Equal(t, 1, tr.EraseLeftmost(1, -1))

	requireTreapValues(t, tr, 4)
}

func TestEraseRangeAndPanics(t *testing.T) {
	tr := NewAutoOrderTreap(1, 2, 3, 4, 5)

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
	tr := NewAutoOrderTreap(1, 2, 3, 4, 5, 6, 7)

	require.Equal(t, 3, tr.EraseRange(2, false, 6, false))

	requireTreapValues(t, tr, 1, 2, 6, 7)

	require.Equal(t, 0, tr.EraseRange(100, true, 200, true))
}

func TestEraseLeftmostZeroCount(t *testing.T) {
	tr := NewAutoOrderTreap(1, 1, 1, 2, 2)

	require.Equal(t, 0, tr.EraseLeftmost(1, 0))
	requireTreapValues(t, tr, 1, 1, 1, 2, 2)

	require.Equal(t, 3, tr.EraseLeftmost(1, 10))
	requireTreapValues(t, tr, 2, 2)
}

func TestEraseRightmostNegativeRemovesAll(t *testing.T) {
	tr := NewAutoOrderTreap(3, 3, 3, 4, 5)

	require.Equal(t, 3, tr.EraseRightmost(3, -1))
	requireTreapValues(t, tr, 4, 5)

	require.Equal(t, 0, tr.EraseRightmost(5, 0))
	requireTreapValues(t, tr, 4, 5)
}

func TestEraseAt(t *testing.T) {
	tr := NewAutoOrderTreap(1, 2, 3, 4, 5)
	require.Equal(t, 3, tr.EraseAt(1, 3))
	requireTreapValues(t, tr, 1, 5)
	require.Panics(t, func() { tr.EraseAt(-1, 1) })
	require.Panics(t, func() { tr.EraseAt(0, -1) })
}

func TestEraseAtZeroOrExcessCount(t *testing.T) {
	tr := NewAutoOrderTreap(1, 2, 3, 4)
	require.Equal(t, 0, tr.EraseAt(2, 0))

	require.Equal(t, 3, tr.EraseAt(1, 10))

	requireTreapValues(t, tr, 1)
}

func TestBoundsAndLookup(t *testing.T) {
	tr := NewAutoOrderTreap(1, 2, 2, 3, 4, 5)

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
	tr := NewAutoOrderTreap(1, 3, 5)

	node, idx := tr.FindLowerBound(6)
	require.Nil(t, node)
	require.Equal(t, 0, idx)

	node, idx = tr.FindUpperBound(0)
	require.Nil(t, node)
	require.Equal(t, 0, idx)

	require.Nil(t, tr.At(10))

	node = tr.At(-1)
	require.NotNil(t, node)
	require.Equal(t, 1, node.value)
}

func TestSizeEmptyClear(t *testing.T) {
	tr := NewAutoOrderTreap[int]()
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
	tr := NewAutoOrderTreap(5, 6, 7)
	tr.Clear()
	require.True(t, tr.Empty())

	require.Equal(t, 0, tr.InsertLeft(10))
	require.Equal(t, 1, tr.Size())

	requireTreapValues(t, tr, 10)
}

func TestExtremaAndPops(t *testing.T) {
	tr := NewAutoOrderTreap(3, 1, 4, 1, 5)
	require.Equal(t, 1, tr.Leftmost().value)
	require.Equal(t, 5, tr.Rightmost().value)

	val, ok := tr.PopLeftmost()
	require.True(t, ok)
	require.Equal(t, 1, val)
	val, ok = tr.PopRightmost()
	require.True(t, ok)
	require.Equal(t, 5, val)

	empty := NewAutoOrderTreap[int]()
	_, ok = empty.PopLeftmost()
	require.False(t, ok)
	_, ok = empty.PopRightmost()
	require.False(t, ok)
}

func TestExtremaAndPopsOnSingleton(t *testing.T) {
	tr := NewAutoOrderTreap(42)

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
	tr := NewAutoOrderTreap(4, 1, 3, 2)

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
	tr := NewAutoOrderTreap(1, 2, 3, 4, 5)

	left, right := tr.SplitBefore(3)
	requireTreapValues(t, left, 1, 2)
	requireTreapValues(t, right, 3, 4, 5)

	tr = NewAutoOrderTreap(1, 2, 3, 4, 5)
	left, right = tr.SplitAfter(3)
	requireTreapValues(t, left, 1, 2, 3)
	requireTreapValues(t, right, 4, 5)

	tr = NewAutoOrderTreap(1, 2, 3, 4, 5)
	left, right = tr.Cut(3)
	requireTreapValues(t, left, 1, 2, 3)
	requireTreapValues(t, right, 4, 5)
}

func TestSplitBeforeAndAfterBoundaries(t *testing.T) {
	tr := NewAutoOrderTreap(2, 4, 6)
	left, right := tr.SplitBefore(1)
	require.True(t, left.Empty())
	requireTreapValues(t, right, 2, 4, 6)

	tr = NewAutoOrderTreap(2, 4, 6)
	left, right = tr.SplitBefore(10)
	requireTreapValues(t, left, 2, 4, 6)
	require.True(t, right.Empty())

	tr = NewAutoOrderTreap(2, 4, 6)
	left, right = tr.SplitAfter(1)
	require.True(t, left.Empty())
	requireTreapValues(t, right, 2, 4, 6)

	tr = NewAutoOrderTreap(2, 4, 6)
	left, right = tr.SplitAfter(10)
	requireTreapValues(t, left, 2, 4, 6)
	require.True(t, right.Empty())
}

func TestSplitClearsOriginal(t *testing.T) {
	tr := NewAutoOrderTreap(1, 2, 3, 4)
	left, right := tr.SplitBefore(3)

	require.True(t, tr.Empty())
	requireTreapValues(t, left, 1, 2)
	requireTreapValues(t, right, 3, 4)

	require.NotNil(t, left.lessFn)
	require.NotNil(t, right.lessFn)
}

func TestSplitWithDuplicateBoundaries(t *testing.T) {
	tr := NewAutoOrderTreap(1, 2, 2, 2, 3, 4)

	left, right := tr.SplitBefore(2)
	requireTreapValues(t, left, 1)
	requireTreapValues(t, right, 2, 2, 2, 3, 4)

	tr = NewAutoOrderTreap(1, 2, 2, 2, 3, 4)
	left, right = tr.SplitAfter(2)
	requireTreapValues(t, left, 1, 2, 2, 2)
	requireTreapValues(t, right, 3, 4)
}

func TestCountRangeAndCount(t *testing.T) {
	tr := NewAutoOrderTreap(1, 2, 2, 3, 4, 4, 5)

	require.Equal(t, 5, tr.CountRange(2, true, 4, true))
	require.Equal(t, 1, tr.CountRange(2, false, 4, false))
	require.Equal(t, 0, tr.CountRange(6, true, 9, true))

	require.Panics(t, func() { tr.CountRange(5, true, 4, true) })
	require.Panics(t, func() { tr.CountRange(5, true, 5, false) })
	require.Panics(t, func() { tr.CountRange(5, false, 5, true) })

	require.Equal(t, 2, tr.Count(4))
}

func TestCountRangeOutsideBounds(t *testing.T) {
	tr := NewAutoOrderTreap(5, 10, 15)

	require.Zero(t, tr.CountRange(-10, true, 0, true))
	require.Zero(t, tr.CountRange(20, true, 30, true))
	require.Zero(t, tr.Count(7))
}

func TestCountRangeExclusiveEliminatesBoundaries(t *testing.T) {
	tr := NewAutoOrderTreap(1, 2, 3, 4)

	require.Equal(t, 0, tr.CountRange(1, false, 2, false))
	require.Equal(t, 3, tr.CountRange(1, false, 4, true))
	require.Equal(t, 3, tr.CountRange(1, true, 4, false))
}

func TestCutEdgeCounts(t *testing.T) {
	tr := NewAutoOrderTreap(1, 2, 3)
	left, right := tr.Cut(0)
	require.True(t, left.Empty())
	requireTreapValues(t, right, 1, 2, 3)

	tr = NewAutoOrderTreap(1, 2, 3)
	left, right = tr.Cut(10)
	requireTreapValues(t, left, 1, 2, 3)
	require.True(t, right.Empty())

	tr = NewAutoOrderTreap(1, 2, 3)
	left, right = tr.Cut(-5)
	require.True(t, left.Empty())
	requireTreapValues(t, right, 1, 2, 3)
}

func TestMergeTreap(t *testing.T) {
	left := NewAutoOrderTreap(1, 2, 3)
	right := NewAutoOrderTreap(4, 5, 6)
	merged := Merge(left, right)

	requireTreapValues(t, merged, 1, 2, 3, 4, 5, 6)

	require.True(t, merged.lessFn(1, 2))
	require.False(t, merged.lessFn(2, 1))
}

func TestMergeTreapSupportsFurtherInsertion(t *testing.T) {
	left := NewAutoOrderTreap(1, 3)
	right := NewAutoOrderTreap(5, 7)

	merged := Merge(left, right)
	idx := merged.InsertRight(4)
	require.Equal(t, 2, idx)
	idx = merged.InsertLeft(6)
	require.Equal(t, 4, idx)

	requireTreapValues(t, merged, 1, 3, 4, 5, 6, 7)
}

func TestMergeTreapWithEmpty(t *testing.T) {
	left := NewAutoOrderTreap(1, 2, 3)
	empty := NewAutoOrderTreap[int]()

	merged := Merge(left, empty)
	requireTreapValues(t, merged, 1, 2, 3)

	merged = Merge(empty, left)
	requireTreapValues(t, merged, 1, 2, 3)

	require.NotNil(t, merged.lessFn)
}

func FuzzTreapInsertEraseConsistency(f *testing.F) {
	seeds := [][]byte{
		{0, 10, 1, 12, 2, 0, 15, 1, 3},
		{1, 20, 0, 3, 2, 4, 8, 1, 1, 3, 0, 7, 1, 4},
		{2, 6, 2, 15, 1, 1, 0, 9, 3, 0, 2, 5, 1, 1, 3},
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		tr := NewAutoOrderTreap[int]()
		reference := make([]int, 0)

		readByte := func(idx *int) (byte, bool) {
			if *idx >= len(data) {
				return 0, false
			}
			b := data[*idx]
			*idx = *idx + 1
			return b, true
		}

		steps := 0
		for i := 0; i < len(data) && steps < 512; steps++ {
			b, ok := readByte(&i)
			if !ok {
				break
			}
			op := b % 4
			switch op {
			case 0:
				b, ok := readByte(&i)
				if !ok {
					return
				}
				val := int(int8(b))
				if b&1 == 0 {
					tr.InsertLeft(val)
				} else {
					tr.InsertRight(val)
				}
				reference = append(reference, val)
			case 1:
				if len(reference) == 0 {
					continue
				}
				b, ok := readByte(&i)
				if !ok {
					return
				}
				var (
					removedTreap int
					val          int
				)
				if len(reference) > 0 {
					val = reference[int(b)%len(reference)]
				}
				removedTreap = tr.EraseAll(val)
				kept := make([]int, 0, len(reference))
				for _, v := range reference {
					if v != val {
						kept = append(kept, v)
					}
				}
				require.Equal(t, len(reference)-len(kept), removedTreap)
				reference = kept
			case 2:
				if len(reference) < 2 {
					continue
				}
				lowRaw, ok := readByte(&i)
				if !ok {
					return
				}
				highRaw, ok := readByte(&i)
				if !ok {
					return
				}
				inclRaw, ok := readByte(&i)
				if !ok {
					return
				}
				low := int(int8(lowRaw))
				high := int(int8(highRaw))
				if low > high {
					low, high = high, low
				}
				inclusiveLow := inclRaw&1 == 0
				inclusiveHigh := inclRaw&2 == 0
				if low == high && !(inclusiveLow && inclusiveHigh) {
					inclusiveLow, inclusiveHigh = true, true
				}
				removed := tr.EraseRange(low, inclusiveLow, high, inclusiveHigh)
				filtered := make([]int, 0, len(reference))
				for _, v := range reference {
					if v < low || (!inclusiveLow && v == low) {
						filtered = append(filtered, v)
						continue
					}
					if v > high || (!inclusiveHigh && v == high) {
						filtered = append(filtered, v)
						continue
					}
				}
				require.Equal(t, len(reference)-len(filtered), removed)
				reference = filtered
			case 3:
				if len(reference) == 0 {
					continue
				}
				dirRaw, ok := readByte(&i)
				if !ok {
					return
				}
				if dirRaw&1 == 0 {
					val, ok := tr.PopLeftmost()
					require.True(t, ok)
					sort.Ints(reference)
					require.Equal(t, reference[0], val)
					reference = reference[1:]
				} else {
					val, ok := tr.PopRightmost()
					require.True(t, ok)
					sort.Ints(reference)
					require.Equal(t, reference[len(reference)-1], val)
					reference = reference[:len(reference)-1]
				}
			}

			sort.Ints(reference)
			if reference == nil {
				reference = []int{}
			}
			require.Equal(t, reference, mustValues(tr))
			require.Equal(t, len(reference), tr.Size())
		}
	})
}

func FuzzTreapIndexAndBounds(f *testing.F) {
	seeds := [][]byte{
		{0, 5, 1, 0, 2, 3, 3, 1, 4, 2},
		{0, 10, 0, 12, 0, 14, 1, 2, 2, 1, 3, 0, 4, 0},
		{0, 2, 0, 4, 0, 6, 1, 3, 2, 0, 3, 0, 4, 1},
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		tr := NewAutoOrderTreap[int]()
		reference := make([]int, 0)

		readByte := func(idx *int) (byte, bool) {
			if *idx >= len(data) {
				return 0, false
			}
			b := data[*idx]
			*idx = *idx + 1
			return b, true
		}

		steps := 0
		for i := 0; i < len(data) && steps < 512; steps++ {
			b, ok := readByte(&i)
			if !ok {
				break
			}
			op := b % 5
			switch op {
			case 0:
				b, ok := readByte(&i)
				if !ok {
					return
				}
				val := int(int8(b))
				if b&1 == 0 {
					tr.InsertLeft(val)
				} else {
					tr.InsertRight(val)
				}
				reference = append(reference, val)
			case 1:
				if len(reference) == 0 {
					continue
				}
				b, ok := readByte(&i)
				if !ok {
					return
				}
				idx := int(b) % len(reference)
				removed := tr.EraseAt(idx, 1)
				if removed > 0 {
					reference = append(reference[:idx], reference[idx+1:]...)
				}
			case 2:
				if len(reference) == 0 {
					continue
				}
				b, ok := readByte(&i)
				if !ok {
					return
				}
				val := reference[int(b)%len(reference)]
				tr.EraseLeftmost(val, 1)
				for idx, v := range reference {
					if v == val {
						reference = append(reference[:idx], reference[idx+1:]...)
						break
					}
				}
			case 3:
				if len(reference) == 0 {
					continue
				}
				sort.Ints(reference)
				b, ok := readByte(&i)
				if !ok {
					return
				}
				idx := int(b) % len(reference)
				node := tr.At(idx)
				require.NotNil(t, node)
				require.Equal(t, reference[idx], node.Value())
			case 4:
				if len(reference) == 0 {
					continue
				}
				sort.Ints(reference)
				low := reference[0]
				high := reference[len(reference)-1]
				count := tr.CountRange(low, true, high, true)
				require.Equal(t, len(reference), count)
			}

			sort.Ints(reference)
			if reference == nil {
				reference = []int{}
			}
			require.Equal(t, len(reference), tr.Size())
		}
	})
}
