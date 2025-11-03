package gotreap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// mustNode constructs a node with the provided value, priority and children.
// It also ensures parent pointers and sizes are kept up to date which is
// essential for verifying treap invariants in white-box tests.
func mustNode[T any](value T, priority int, left *Node[T], right *Node[T]) *Node[T] {
	n := &Node[T]{
		value:          value,
		heightPriority: priority,
		left:           left,
		right:          right,
	}
	if left != nil {
		left.parent = n
	}
	if right != nil {
		right.parent = n
	}
	n.recalcSize()
	return n
}

// mustInOrder returns the in-order traversal of a treap.
func mustInOrder[T any](root *Node[T]) []T {
	if root == nil {
		return nil
	}
	var res []T
	var visit func(*Node[T])
	visit = func(cur *Node[T]) {
		if cur == nil {
			return
		}
		visit(cur.left)
		res = append(res, cur.value)
		visit(cur.right)
	}
	visit(root)
	return res
}

func requireInOrder[T any](t *testing.T, root *Node[T], expected ...T) {
	t.Helper()
	require.Equal(t, expected, mustInOrder(root))
}

func requireNodeIntegrity[T any](t *testing.T, nodes ...*Node[T]) {
	t.Helper()
	for _, node := range nodes {
		if node == nil {
			continue
		}
		if node.left != nil {
			require.Equal(t, node, node.left.parent, "left child parent mismatch for %v", node.value)
		}
		if node.right != nil {
			require.Equal(t, node, node.right.parent, "right child parent mismatch for %v", node.value)
		}
		expectedSize := node.left.safeSize() + 1 + node.right.safeSize()
		require.Equalf(t, expectedSize, node.size, "node %v size", node.value)
	}
}

func TestNodeSafeSizeAndRecalc(t *testing.T) {
	require.Equal(t, 0, (*Node[int])(nil).safeSize())

	leaf := mustNode(5, 10, nil, nil)
	require.Equal(t, 1, leaf.safeSize())

	parent := mustNode(7, 5, leaf, nil)
	require.Equal(t, 2, parent.safeSize())

	parent.left = nil
	parent.recalcSize()
	require.Equal(t, 1, parent.safeSize())
}

func TestNodeSafeSetParent(t *testing.T) {
	child := mustNode(1, 1, nil, nil)
	parent := mustNode(2, 2, nil, nil)

	child.safeSetParent(parent)
	require.Equal(t, parent, child.parent)

	(*Node[int])(nil).safeSetParent(parent)
}

func TestNodeSafeSetParentCanClearParent(t *testing.T) {
	child := mustNode(2, 2, nil, nil)
	parent := mustNode(1, 1, child, nil)

	require.Equal(t, parent, child.parent)

	child.safeSetParent(nil)
	require.Nil(t, child.parent)
}

func TestMergeMaintainsHeapAndOrder(t *testing.T) {
	// Construct a left-heavy tree and right-heavy tree to stress parent updates.
	left := mustNode(2, 40,
		mustNode(1, 35, nil, nil),
		mustNode(3, 30, nil, nil),
	)
	right := mustNode(6, 25,
		mustNode(5, 15, nil, nil),
		mustNode(7, 20, nil, nil),
	)

	merged := merge(left, right)

	require.Nil(t, merged.parent)
	requireInOrder(t, merged, 1, 2, 3, 5, 6, 7)

	requireNodeIntegrity(t,
		merged,
		merged.left,
		merged.right,
		merged.left.left,
		merged.left.right,
		merged.right.left,
		merged.right.right,
	)
}

func TestMergeHandlesNilSubtrees(t *testing.T) {
	leaf := mustNode(1, 10, nil, nil)

	require.Same(t, leaf, merge(nil, leaf))
	require.Same(t, leaf, merge(leaf, nil))

	left := mustNode(1, 20, nil, nil)
	right := mustNode(2, 10, nil, nil)
	root := merge(left, right)

	require.Same(t, left, root)
	require.Same(t, right, root.right)
	require.Nil(t, left.parent)
	require.Equal(t, root, right.parent)
}

func TestSplitProducesExpectedPartitions(t *testing.T) {
	root := mustNode(4, 150,
		mustNode(2, 120,
			mustNode(1, 100, nil, nil),
			mustNode(3, 110, nil, nil),
		),
		mustNode(6, 90,
			mustNode(5, 80, nil, nil),
			mustNode(7, 70, nil, nil),
		),
	)

	left, right := root.split(func(value int, index int) bool {
		return value < 4
	}, 0)

	requireInOrder(t, left, 1, 2, 3)
	requireInOrder(t, right, 4, 5, 6, 7)

	require.Nil(t, right.parent)
}

func TestSplitEntireTreeGoesOneSide(t *testing.T) {
	root := mustNode(2, 60,
		mustNode(1, 50, nil, nil),
		mustNode(3, 40, nil, nil),
	)

	alwaysTrue := func(value int, index int) bool { return true }
	left, right := root.split(alwaysTrue, 0)
	require.Nil(t, right)
	require.Nil(t, left.parent)

	alwaysFalse := func(value int, index int) bool { return false }
	left, right = root.split(alwaysFalse, 0)
	require.Nil(t, left)
	require.Nil(t, right.parent)
}

func TestSplitHonorsIndexOffset(t *testing.T) {
	root := mustNode(3, 110,
		mustNode(2, 100,
			mustNode(1, 90, nil, nil),
			nil,
		),
		mustNode(5, 85,
			mustNode(4, 80, nil, nil),
			mustNode(6, 70, nil, nil),
		),
	)

	// Provide a non-zero offset to emulate the tree living in a larger structure.
	left, right := root.split(func(_ int, idx int) bool { return idx < 4 }, 2)

	requireInOrder(t, left, 1, 2)
	requireInOrder(t, right, 3, 4, 5, 6)
	require.Nil(t, left.parent)
	require.Nil(t, right.parent)
}

func TestPrevAndNext(t *testing.T) {
	root := mustNode(2, 200,
		mustNode(1, 100, nil, nil),
		mustNode(3, 150, nil, nil),
	)

	require.Equal(t, root, root.left.Next())
	require.Equal(t, root, root.right.Prev())
	require.Nil(t, root.left.Prev())
	require.Nil(t, root.right.Next())
}

func TestPrevAndNextAcrossAncestors(t *testing.T) {
	root := mustNode(10, 150,
		mustNode(5, 130,
			mustNode(2, 100, nil, nil),
			mustNode(7, 120,
				mustNode(6, 110, nil, nil),
				mustNode(8, 90, nil, nil),
			),
		),
		mustNode(15, 95,
			mustNode(12, 80, nil, nil),
			mustNode(20, 70, nil, nil),
		),
	)

	got := root.left.right.left.Next()
	require.NotNil(t, got)
	require.Equal(t, 7, got.value)
	got = root.left.right.right.Prev()
	require.NotNil(t, got)
	require.Equal(t, 7, got.value)
	got = root.left.right.right.Next()
	require.NotNil(t, got)
	require.Equal(t, 10, got.value)
	require.Nil(t, root.left.left.Prev())
	require.Nil(t, root.right.right.Next())
}

func TestLeftmostRightmost(t *testing.T) {
	root := mustNode(3, 70,
		mustNode(2, 60,
			mustNode(1, 50, nil, nil),
			nil,
		),
		mustNode(5, 45,
			mustNode(4, 40, nil, nil),
			mustNode(6, 35, nil, nil),
		),
	)
	require.Equal(t, 1, root.Leftmost().value)
	require.Equal(t, 6, root.Rightmost().value)
}

func TestIndexAndValue(t *testing.T) {
	root := mustNode(3, 110,
		mustNode(2, 100,
			mustNode(1, 90, nil, nil),
			nil,
		),
		mustNode(5, 85,
			mustNode(4, 80, nil, nil),
			mustNode(6, 70, nil, nil),
		),
	)
	inorder := mustInOrder(root)
	for idx, val := range inorder {
		node, _ := root.lookupLeftmostUnmatch(func(nodeValue int, nodeIndex int) bool {
			return nodeIndex < idx
		}, 0)
		require.NotNilf(t, node, "expected node at index %d", idx)
		require.Equal(t, idx, node.Index())
		require.Equal(t, val, node.Value())
	}

	require.Zero(t, (*Node[int])(nil).Value())
}

func TestLookupHelpers(t *testing.T) {
	root := mustNode(4, 150,
		mustNode(2, 140,
			mustNode(1, 120, nil, nil),
			mustNode(3, 130, nil, nil),
		),
		mustNode(6, 115,
			mustNode(5, 110, nil, nil),
			mustNode(7, 105, nil, nil),
		),
	)

	node, idx := root.lookupRightmostMatch(func(value int, index int) bool {
		return value <= 4
	}, 0)
	require.NotNil(t, node)
	require.Equal(t, 4, node.value)
	require.Equal(t, 3, idx)

	node, idx = root.lookupLeftmostUnmatch(func(value int, index int) bool {
		return value < 6
	}, 0)
	require.NotNil(t, node)
	require.Equal(t, 6, node.value)
	require.Equal(t, 5, idx)
}

func TestLookupHelpersWithIndexOffset(t *testing.T) {
	root := mustNode(3, 120,
		mustNode(2, 100,
			mustNode(1, 90, nil, nil),
			nil,
		),
		mustNode(5, 95,
			mustNode(4, 80, nil, nil),
			mustNode(6, 70, nil, nil),
		),
	)

	node, idx := root.lookupRightmostMatch(func(_ int, index int) bool { return index < 6 }, 2)
	require.NotNil(t, node)
	require.Equal(t, 4, node.value)
	require.Equal(t, 5, idx)

	node, idx = root.lookupLeftmostUnmatch(func(_ int, index int) bool { return index < 5 }, 2)
	require.NotNil(t, node)
	require.Equal(t, 4, node.value)
	require.Equal(t, 5, idx)
}

func TestLookupOnEmptyTree(t *testing.T) {
	node, idx := (*Node[int])(nil).lookupRightmostMatch(func(int, int) bool { return true }, 7)
	require.Nil(t, node)
	require.Equal(t, 0, idx)
	node, idx = (*Node[int])(nil).lookupLeftmostUnmatch(func(int, int) bool { return true }, 3)
	require.Nil(t, node)
	require.Equal(t, 0, idx)
}

func TestValidReportsPresence(t *testing.T) {
	require.False(t, (*Node[int])(nil).Valid())
	require.True(t, mustNode(1, 1, nil, nil).Valid())
}

func TestNilNodeHelpersReturnIdentityValues(t *testing.T) {
	var node *Node[int]

	require.Nil(t, node.Prev())
	require.Nil(t, node.Next())
	require.Nil(t, node.Leftmost())
	require.Nil(t, node.Rightmost())
	require.Equal(t, -1, node.Index())
}

func TestSplitDetachesParentsAcrossPartitions(t *testing.T) {
	root := mustNode(10, 90,
		mustNode(5, 70,
			mustNode(2, 50, nil, nil),
			mustNode(8, 60, nil, nil),
		),
		mustNode(15, 55,
			mustNode(12, 40, nil, nil),
			mustNode(18, 35, nil, nil),
		),
	)

	left, right := root.split(func(value int, index int) bool { return value < 10 }, 0)

	requireInOrder(t, left, 2, 5, 8)
	requireInOrder(t, right, 10, 12, 15, 18)

	require.Nil(t, left.parent)
	require.Nil(t, right.parent)

	var verifyParents func(cur *Node[int])
	verifyParents = func(cur *Node[int]) {
		if cur == nil {
			return
		}
		if cur.left != nil {
			require.Equal(t, cur, cur.left.parent)
		}
		if cur.right != nil {
			require.Equal(t, cur, cur.right.parent)
		}
		verifyParents(cur.left)
		verifyParents(cur.right)
	}

	verifyParents(left)
	verifyParents(right)
}
