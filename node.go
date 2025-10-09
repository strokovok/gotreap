package util

type leftCondition[T any] func(nodeValue T, nodeIndex int) bool

type Node[T any] struct {
	value          T
	heightPriority int
	left           *Node[T]
	right          *Node[T]
	parent         *Node[T]
	size           int
}

// newNode creates a new treap node containing value with a random heap priority.
func newNode[T any](value T, heightPriority int) *Node[T] {
	return &Node[T]{
		value:          value,
		heightPriority: heightPriority,
		left:           nil,
		right:          nil,
		parent:         nil,
		size:           1,
	}
}

// safeSize returns the subtree size stored in t, treating a nil node as zero.
func (t *Node[T]) safeSize() int {
	if t == nil {
		return 0
	}
	return t.size
}

// recalcSize recomputes t.size based on its children's sizes.
func (t *Node[T]) recalcSize() {
	t.size = t.left.safeSize() + 1 + t.right.safeSize()
}

// safeSetParent assigns parent to t when t is non-nil.
func (t *Node[T]) safeSetParent(parent *Node[T]) {
	if t == nil {
		return
	}
	t.parent = parent
}

// merge combines two priority-ordered treap subtrees preserving in-order sequence.
func merge[T any](left *Node[T], right *Node[T]) *Node[T] {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}

	if left.heightPriority >= right.heightPriority {
		left.right = merge(left.right, right)
		left.right.safeSetParent(left)
		left.recalcSize()
		return left
	}

	right.left = merge(left, right.left)
	right.left.safeSetParent(right)
	right.recalcSize()
	return right
}

// split partitions the treap into nodes satisfying leftCond (left) and the rest (right).
func (t *Node[T]) split(leftCond leftCondition[T], indexOffset int) (left, right *Node[T]) {
	if t == nil {
		return nil, nil
	}

	centralIndexOffset := indexOffset + t.left.safeSize()
	if leftCond(t.value, centralIndexOffset) {
		t.right, right = t.right.split(leftCond, centralIndexOffset+1)
		t.right.safeSetParent(t)
		right.safeSetParent(nil)
		t.recalcSize()
		return t, right
	}

	left, t.left = t.left.split(leftCond, indexOffset)
	left.safeSetParent(nil)
	t.left.safeSetParent(t)
	t.recalcSize()
	return left, t
}

// Prev returns the in-order predecessor of t within the treap.
func (t *Node[T]) Prev() *Node[T] {
	if t == nil {
		return nil
	}

	if t.left != nil {
		cur := t.left
		for cur.right != nil {
			cur = cur.right
		}
		return cur
	}

	for cur := t; cur.parent != nil; cur = cur.parent {
		if cur.parent.right == cur {
			return cur.parent
		}
	}
	return nil
}

// Next returns the in-order successor of t within the treap.
func (t *Node[T]) Next() *Node[T] {
	if t == nil {
		return nil
	}

	if t.right != nil {
		cur := t.right
		for cur.left != nil {
			cur = cur.left
		}
		return cur
	}

	for cur := t; cur.parent != nil; cur = cur.parent {
		if cur.parent.left == cur {
			return cur.parent
		}
	}
	return nil
}

// Leftmost returns the minimum node in the treap containing t.
func (t *Node[T]) Leftmost() *Node[T] {
	if t == nil {
		return nil
	}

	cur := t
	for cur.parent != nil {
		cur = cur.parent
	}
	for cur.left != nil {
		cur = cur.left
	}
	return cur
}

// Rightmost returns the maximum node in the treap containing t.
func (t *Node[T]) Rightmost() *Node[T] {
	if t == nil {
		return nil
	}

	cur := t
	for cur.parent != nil {
		cur = cur.parent
	}
	for cur.right != nil {
		cur = cur.right
	}
	return cur
}

// Index computes the zero-based position of t within an in-order traversal.
func (t *Node[T]) Index() int {
	if t == nil {
		return 0
	}

	indexOffset := t.left.safeSize()
	for cur := t; cur.parent != nil; cur = cur.parent {
		if cur.parent.right == cur {
			indexOffset += cur.parent.left.safeSize() + 1
		}
	}
	return indexOffset
}

// Valid reports whether t references an actual node.
func (t *Node[T]) Valid() bool {
	return t != nil
}

// Value returns the stored node value or the zero value if t is nil.
func (t *Node[T]) Value() (result T) {
	if t != nil {
		result = t.value
	}
	return result
}

// lookupRightmostMatch finds the rightmost node satisfying leftCond together with its index.
func (t *Node[T]) lookupRightmostMatch(leftCond leftCondition[T], indexOffset int) (node *Node[T], index int) {
	if t == nil {
		return nil, 0
	}

	centralIndexOffset := indexOffset + t.left.safeSize()
	if leftCond(t.value, centralIndexOffset) {
		res, idx := t.right.lookupRightmostMatch(leftCond, centralIndexOffset+1)
		if res != nil {
			return res, idx
		}
		return t, centralIndexOffset
	}

	return t.left.lookupRightmostMatch(leftCond, indexOffset)
}

// lookupLeftmostUnmatch finds the leftmost node that fails leftCond and returns it with its index.
func (t *Node[T]) lookupLeftmostUnmatch(leftCond leftCondition[T], indexOffset int) (node *Node[T], index int) {
	if t == nil {
		return nil, 0
	}

	centralIndexOffset := indexOffset + t.left.safeSize()
	if leftCond(t.value, centralIndexOffset) {
		return t.right.lookupLeftmostUnmatch(leftCond, centralIndexOffset+1)
	}

	res, idx := t.left.lookupLeftmostUnmatch(leftCond, indexOffset)
	if res != nil {
		return res, idx
	}
	return t, centralIndexOffset
}
