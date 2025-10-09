package util

import (
	"cmp"
	"sort"
)

type Treap[T any] struct {
	lessFn func(a T, b T) bool
	root   *Node[T]
}

func NewTreap[T any](lessFn func(a T, b T) bool, values ...T) *Treap[T] {
	t := &Treap[T]{
		lessFn: lessFn,
		root:   nil,
	}

	sort.Slice(values, func(i, j int) bool {
		return lessFn(values[i], values[j])
	})

	for _, val := range values {
		t.root = merge(t.root, newNode(val))
	}

	return t
}

func NewAutoOrderTreap[T cmp.Ordered](values ...T) *Treap[T] {
	return NewTreap(cmp.Less[T], values...)
}

func (t *Treap[T]) condLess(value T) leftCondition[T] {
	return func(nodeValue T, nodeIndex int) bool {
		return t.lessFn(nodeValue, value)
	}
}

func (t *Treap[T]) condLeq(value T) leftCondition[T] {
	return func(nodeValue T, nodeIndex int) bool {
		return !t.lessFn(value, nodeValue)
	}
}

func (t *Treap[T]) condCutN(n int) leftCondition[T] {
	return func(nodeValue T, nodeIndex int) bool {
		return nodeIndex < n
	}
}

func (t *Treap[T]) InsertLeft(value T) (index int) {
	less, greaterOrEqual := t.root.split(t.condLess(value), 0)

	index = less.safeSize()

	greaterOrEqual = merge(newNode(value), greaterOrEqual)
	t.root = merge(less, greaterOrEqual)

	return index
}

func (t *Treap[T]) InsertRight(value T) (index int) {
	lessOrEqual, greater := t.root.split(t.condLeq(value), 0)

	index = lessOrEqual.safeSize()

	lessOrEqual = merge(lessOrEqual, newNode(value))
	t.root = merge(lessOrEqual, greater)

	return index
}

func (t *Treap[T]) EraseAll(value T) (erasedCount int) {
	less, greaterOrEqual := t.root.split(t.condLess(value), 0)

	equal, greater := greaterOrEqual.split(t.condLeq(value), 0)

	t.root = merge(less, greater)

	return equal.safeSize()
}

func (t *Treap[T]) EraseLeftmost(value T, n int) (erasedCount int) {
	less, greaterOrEqual := t.root.split(t.condLess(value), 0)

	equal, greater := greaterOrEqual.split(t.condLeq(value), 0)

	if n < 0 {
		n = equal.safeSize()
	}
	equalErased, equalRemainder := equal.split(t.condCutN(n), 0)

	t.root = merge(less, merge(equalRemainder, greater))

	return equalErased.safeSize()
}

func (t *Treap[T]) EraseRightmost(value T, n int) (erasedCount int) {
	less, greaterOrEqual := t.root.split(t.condLess(value), 0)

	equal, greater := greaterOrEqual.split(t.condLeq(value), 0)

	if n < 0 {
		n = equal.safeSize()
	}
	remainderN := equal.safeSize() - n
	equalRemainder, equalErased := equal.split(t.condCutN(remainderN), 0)

	t.root = merge(less, merge(equalRemainder, greater))

	return equalErased.safeSize()
}

func (t *Treap[T]) EraseRange(startValue T, inclusiveStart bool, endValue T, inclusiveEnd bool) (erasedCount int) {
	if t.lessFn(endValue, startValue) {
		panic("provided endValue must not be lower than startValue")
	}
	if !t.lessFn(startValue, endValue) && (!inclusiveStart || !inclusiveEnd) {
		panic("when startValue == endValue, both start and end must be inclusive")
	}

	var leftRemainder, toErase, rightRemainder *Node[T]

	if inclusiveStart {
		leftRemainder, rightRemainder = t.root.split(t.condLess(startValue), 0)
	} else {
		leftRemainder, rightRemainder = t.root.split(t.condLeq(startValue), 0)
	}

	if inclusiveEnd {
		toErase, rightRemainder = rightRemainder.split(t.condLeq(endValue), 0)
	} else {
		toErase, rightRemainder = rightRemainder.split(t.condLess(endValue), 0)
	}

	t.root = merge(leftRemainder, rightRemainder)

	return toErase.safeSize()
}

func (t *Treap[T]) EraseAt(index int, count int) (erasedCount int) {
	if index < 0 || count < 0 {
		panic("index and count must not be negative")
	}

	leftRemainder, rightRemainder := t.root.split(t.condCutN(index), 0)

	toErase, rightRemainder := rightRemainder.split(t.condCutN(count), 0)

	t.root = merge(leftRemainder, rightRemainder)

	return toErase.safeSize()
}

// Find leftmost element that is equal to or greater than provided
func (t *Treap[T]) FindLowerBound(value T) (node *Node[T], index int) {
	return t.root.lookupLeftmostUnmatch(t.condLess(value), 0)
}

// Find rightmost element that is equal to or lower than provided
func (t *Treap[T]) FindUpperBound(value T) (node *Node[T], index int) {
	return t.root.lookupRightmostMatch(t.condLeq(value), 0)
}

func (t *Treap[T]) At(index int) *Node[T] {
	node, _ := t.root.lookupLeftmostUnmatch(t.condCutN(index), 0)
	return node
}

func (t *Treap[T]) Size() int {
	return t.root.safeSize()
}

func (t *Treap[T]) Empty() bool {
	return t.root.safeSize() == 0
}

func (t *Treap[T]) Clear() {
	t.root = nil
}

func (t *Treap[T]) Leftmost() *Node[T] {
	return t.root.Leftmost()
}

func (t *Treap[T]) Rightmost() *Node[T] {
	return t.root.Rightmost()
}

func (t *Treap[T]) PopLeftmost() (value T, ok bool) {
	if t.root == nil {
		return value, false
	}

	var leftmost *Node[T]
	leftmost, t.root = t.root.split(t.condCutN(1), 0)

	return leftmost.value, true
}

func (t *Treap[T]) PopRightmost() (value T, ok bool) {
	if t.root == nil {
		return value, false
	}

	var rightmost *Node[T]
	cutN := t.root.safeSize() - 1
	t.root, rightmost = t.root.split(t.condCutN(cutN), 0)

	return rightmost.value, true
}

func (t *Treap[T]) split(leftCond leftCondition[T]) (left *Treap[T], right *Treap[T]) {
	less, greaterOrEqual := t.root.split(leftCond, 0)

	left = &Treap[T]{
		lessFn: t.lessFn,
		root:   less,
	}

	right = &Treap[T]{
		lessFn: t.lessFn,
		root:   greaterOrEqual,
	}

	t.root = nil

	return left, right
}

func (t *Treap[T]) SplitBefore(value T) (left *Treap[T], right *Treap[T]) {
	return t.split(t.condLess(value))
}

func (t *Treap[T]) SplitAfter(value T) (left *Treap[T], right *Treap[T]) {
	return t.split(t.condLeq(value))
}

func (t *Treap[T]) Cut(n int) (left *Treap[T], right *Treap[T]) {
	return t.split(t.condCutN(n))
}

func (t *Treap[T]) CountRange(startValue T, inclusiveStart bool, endValue T, inclusiveEnd bool) int {
	if t.lessFn(endValue, startValue) {
		panic("provided endValue must not be lower than startValue")
	}
	if !t.lessFn(startValue, endValue) && (!inclusiveStart || !inclusiveEnd) {
		panic("when startValue == endValue, both start and end must be inclusive")
	}

	var startNode, endNode *Node[T]
	var startIdx, endIdx int

	if inclusiveStart {
		startNode, startIdx = t.root.lookupLeftmostUnmatch(t.condLess(startValue), 0)
	} else {
		startNode, startIdx = t.root.lookupLeftmostUnmatch(t.condLeq(startValue), 0)
	}
	if startNode == nil {
		return 0
	}

	if inclusiveEnd {
		endNode, endIdx = t.root.lookupRightmostMatch(t.condLeq(endValue), 0)
	} else {
		endNode, endIdx = t.root.lookupRightmostMatch(t.condLess(endValue), 0)
	}
	if endNode == nil {
		return 0
	}

	if endIdx < startIdx {
		return 0
	}
	return endIdx - startIdx + 1
}

func (t *Treap[T]) Count(val T) int {
	return t.CountRange(val, true, val, true)
}

func Merge[T any](left *Treap[T], right *Treap[T]) *Treap[T] {
	return &Treap[T]{
		lessFn: left.lessFn,
		root:   merge(left.root, right.root),
	}
}
