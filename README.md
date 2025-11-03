# 🌳 gotreap

[![Go Reference](https://pkg.go.dev/badge/github.com/strokovok/gotreap.svg)](https://pkg.go.dev/github.com/strokovok/gotreap)
[![Go Report Card](https://goreportcard.com/badge/github.com/strokovok/gotreap)](https://goreportcard.com/report/github.com/strokovok/gotreap)
[![License: CC0-1.0](https://img.shields.io/badge/License-CC0_1.0-lightgrey.svg)](http://creativecommons.org/publicdomain/zero/1.0/)

**A complete, type-safe Treap (Tree + Heap) data structure implementation in Go with generics.**

Gotreap provides a self-balancing binary search tree that combines the properties of a binary search tree and a max-heap. Perfect for scenarios requiring fast ordered operations, range queries, and efficient splits/merges.

---

## 📋 Table of Contents

- [What is a Treap?](#-what-is-a-treap)
- [Features](#-features)
- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [Core Operations](#-core-operations)
- [Advanced Usage](#-advanced-usage)
- [API Reference](#-api-reference)
- [Performance](#-performance)
- [Examples](#-examples)
- [Contributing](#-contributing)
- [License](#-license)

---

## 🤔 What is a Treap?

A **Treap** is a randomized balanced binary search tree that maintains two properties simultaneously:

1. **Binary Search Tree (BST)** property: In-order traversal yields sorted elements
2. **Max-Heap** property: Each node has a random priority higher than its children

This dual property ensures **O(log n)** expected time for insertions, deletions, and searches without requiring complex rotation logic like AVL or Red-Black trees.

### Why Choose Treap?

- ✅ **Simple Implementation** - Easier to understand and debug than AVL/Red-Black trees
- ✅ **Self-Balancing** - Automatic balancing via random priorities
- ✅ **Fast Split/Merge** - O(log n) tree splitting and merging operations
- ✅ **Order Statistics** - Direct index-based access to the k-th element
- ✅ **Range Operations** - Efficient range queries and bulk deletions
- ✅ **Type-Safe** - Full Go generics support for any comparable type

---

## ✨ Features

- 🎯 **Type-safe generics** for any ordered or custom-comparable types
- 🔍 **Fast lookups** with O(log n) search, insert, and delete
- 📊 **Index-based access** - Get element by position like an array
- 🎲 **Randomized balancing** - No manual rebalancing needed
- ➗ **Split & Merge** - Divide and combine treaps efficiently
- 🔢 **Range operations** - Count, erase, and query ranges
- 🔄 **Iterator support** - Forward and backward iteration with Go 1.23+ iter.Seq
- 🧭 **Bidirectional navigation** - Next, Prev, Jump operations on nodes
- 📦 **Zero dependencies** - Only uses Go standard library
- 🧪 **Well-tested** - Comprehensive test suite

---

## 📦 Installation

```bash
go get github.com/strokovok/gotreap
```

**Requirements:** Go 1.23 or higher (uses generics and iter package)

---

## 🚀 Quick Start

```go
package main

import (
    "fmt"
    "github.com/strokovok/gotreap"
)

func main() {
    // Create a treap with integers (automatically sorted)
    treap := gotreap.NewAutoOrderTreap(5, 2, 8, 1, 9, 3)

    // Insert elements
    treap.InsertRight(7)
    treap.InsertLeft(4)

    // Check size
    fmt.Println("Size:", treap.Size()) // Output: 8

    // Access by index (0-based, like a sorted array)
    fmt.Println("Element at index 3:", treap.At(3).Value()) // Output: 4

    // Find minimum and maximum
    fmt.Println("Min:", treap.Leftmost().Value())  // Output: 1
    fmt.Println("Max:", treap.Rightmost().Value()) // Output: 9

    // Iterate over all values in order
    for value := range treap.Values() {
        fmt.Println(value)
    }
}
```

---

## 🔧 Core Operations

### Creating a Treap

```go
// Auto-ordered treap for built-in comparable types
treap := gotreap.NewAutoOrderTreap(3, 1, 4, 1, 5)

// Custom comparison function
reverseTreap := gotreap.NewTreap(
    func(a, b int) bool { return a > b }, // Descending order
    3, 1, 4, 1, 5,
)

// Custom type with custom comparator
type Person struct {
    Name string
    Age  int
}

byAge := gotreap.NewTreap(
    func(a, b Person) bool { return a.Age < b.Age },
    Person{"Alice", 30},
    Person{"Bob", 25},
)
```

### Insertion

```go
// Insert before duplicates
index := treap.InsertLeft(42)

// Insert after duplicates
index := treap.InsertRight(42)
```

### Deletion

```go
// Remove all occurrences of a value
count := treap.EraseAll(42)

// Remove first n occurrences
count := treap.EraseLeftmost(42, 3)

// Remove by index
count := treap.EraseAt(5, 2) // Remove 2 elements starting at index 5

// Remove by range
count := treap.EraseRange(10, true, 20, false) // Remove [10, 20)
```

### Access & Search

```go
// Access by index (supports negative indexing)
node := treap.At(0)    // First element
node := treap.At(-1)   // Last element

// Find bounds
node, idx := treap.FindLowerBound(42) // First element >= 42
node, idx := treap.FindUpperBound(42) // Last element <= 42

// Check existence
count := treap.Count(42)
exists := count > 0
```

---

## 🎓 Advanced Usage

### Range Queries

```go
// Count elements in range [10, 20]
count := treap.CountRange(10, true, 20, true)

// Count elements in range (10, 20)
count := treap.CountRange(10, false, 20, false)

// Erase range [15, 25)
erased := treap.EraseRange(15, true, 25, false)
```

### Splitting & Merging

```go
// Split before value
left, right := treap.SplitBefore(50)  // left: < 50, right: >= 50

// Split after value
left, right := treap.SplitAfter(50)   // left: <= 50, right: > 50

// Split by index
left, right := treap.Cut(5)           // First 5 elements vs rest

// Negative index cuts from end
left, right := treap.Cut(-3)          // All but last 3 vs last 3

// Merge two treaps (must have same ordering function)
merged := gotreap.Merge(left, right)
```

### Navigation & Iteration

```go
// Get extrema and pop
min := treap.Leftmost()
max := treap.Rightmost()

value, ok := treap.PopLeftmost()
value, ok := treap.PopRightmost()

// Node navigation
node := treap.At(5)
next := node.Next()
prev := node.Prev()
jumped := node.JumpRight(10)  // Jump 10 positions right
jumped := node.JumpLeft(5)    // Jump 5 positions left

// Get node index
idx := node.Index()

// Iterate forward
for node := range treap.Elements() {
    fmt.Println(node.Value())
}

// Iterate backward
for value := range treap.ValuesBackwards() {
    fmt.Println(value)
}
```

---

## 📚 API Reference

### Constructor Functions

| Function                                                   | Description                         |
| ---------------------------------------------------------- | ----------------------------------- |
| `NewAutoOrderTreap[T cmp.Ordered](values ...T)`            | Create treap with natural ordering  |
| `NewAutoOrderTreapWithRand[T cmp.Ordered](randFn, values)` | Create treap with custom RNG        |
| `NewTreap[T any](lessFn, values)`                          | Create treap with custom comparator |
| `NewTreapWithRand[T any](lessFn, randFn, values)`          | Full control over ordering and RNG  |

### Insertion Methods

| Method               | Time     | Description                  |
| -------------------- | -------- | ---------------------------- |
| `InsertLeft(value)`  | O(log n) | Insert before equal elements |
| `InsertRight(value)` | O(log n) | Insert after equal elements  |

### Deletion Methods

| Method                                       | Time     | Description                    |
| -------------------------------------------- | -------- | ------------------------------ |
| `EraseAll(value)`                            | O(log n) | Remove all occurrences         |
| `EraseLeftmost(value, n)`                    | O(log n) | Remove first n occurrences     |
| `EraseRightmost(value, n)`                   | O(log n) | Remove last n occurrences      |
| `EraseAt(index, count)`                      | O(log n) | Remove count elements at index |
| `EraseRange(start, inclStart, end, inclEnd)` | O(log n) | Remove elements in range       |
| `Clear()`                                    | O(1)     | Remove all elements            |

### Access Methods

| Method        | Time     | Description                              |
| ------------- | -------- | ---------------------------------------- |
| `At(index)`   | O(log n) | Get element at index (supports negative) |
| `Leftmost()`  | O(log n) | Get minimum element                      |
| `Rightmost()` | O(log n) | Get maximum element                      |
| `Root()`      | O(1)     | Get arbitrary node                       |

### Search Methods

| Method                                       | Time     | Description             |
| -------------------------------------------- | -------- | ----------------------- |
| `FindLowerBound(value)`                      | O(log n) | First element >= value  |
| `FindUpperBound(value)`                      | O(log n) | Last element <= value   |
| `Count(value)`                               | O(log n) | Count occurrences       |
| `CountRange(start, inclStart, end, inclEnd)` | O(log n) | Count elements in range |

### Split & Merge

| Method               | Time     | Description                       |
| -------------------- | -------- | --------------------------------- |
| `SplitBefore(value)` | O(log n) | Split at first element >= value   |
| `SplitAfter(value)`  | O(log n) | Split after last element <= value |
| `Cut(n)`             | O(log n) | Split at index n                  |
| `Merge(left, right)` | O(log n) | Combine two treaps                |

### Utility Methods

| Method           | Time     | Description               |
| ---------------- | -------- | ------------------------- |
| `Size()`         | O(1)     | Number of elements        |
| `Empty()`        | O(1)     | Check if empty            |
| `PopLeftmost()`  | O(log n) | Remove and return minimum |
| `PopRightmost()` | O(log n) | Remove and return maximum |

### Iteration

| Method                | Description                  |
| --------------------- | ---------------------------- |
| `Elements()`          | Iterate nodes left-to-right  |
| `ElementsBackwards()` | Iterate nodes right-to-left  |
| `Values()`            | Iterate values left-to-right |
| `ValuesBackwards()`   | Iterate values right-to-left |

### Node Methods

| Method         | Time     | Description                       |
| -------------- | -------- | --------------------------------- |
| `Value()`      | O(1)     | Get node's value                  |
| `Index()`      | O(log n) | Get node's position (-1 if nil)   |
| `Next()`       | O(log n) | Get next node in order            |
| `Prev()`       | O(log n) | Get previous node in order        |
| `JumpRight(n)` | O(log n) | Jump n positions right            |
| `JumpLeft(n)`  | O(log n) | Jump n positions left             |
| `Leftmost()`   | O(log n) | Get minimum from this node's tree |
| `Rightmost()`  | O(log n) | Get maximum from this node's tree |
| `Valid()`      | O(1)     | Check if node is non-nil          |

---

## ⚡ Performance

All operations have **O(log n)** expected time complexity:

| Operation       | Complexity | Notes                          |
| --------------- | ---------- | ------------------------------ |
| Insert          | O(log n)   | Amortized due to randomization |
| Delete          | O(log n)   | Including range deletions      |
| Search          | O(log n)   | Binary search on BST property  |
| Access by index | O(log n)   | Via augmented size information |
| Split           | O(log n)   | Efficient partition operation  |
| Merge           | O(log n)   | Combine two treaps             |
| Range query     | O(log n)   | Count or access ranges         |

**Space Complexity:** O(n) where n is the number of elements.

### When to Use Treap vs Other Structures

| Use Case                        | Treap              | Alternatives             |
| ------------------------------- | ------------------ | ------------------------ |
| Frequent splits/merges          | ✅ Excellent       | ❌ Slow in most BSTs     |
| Order statistics (k-th element) | ✅ O(log n)        | ⚠️ O(n) in standard BSTs |
| Range operations                | ✅ Fast            | ⚠️ Varies                |
| Simple implementation           | ✅ Yes             | ❌ AVL/RB complex        |
| Predictable worst-case          | ⚠️ Randomized      | ✅ AVL/RB guaranteed     |
| Memory efficiency               | ⚠️ Parent pointers | ✅ Some BSTs             |

---

## 💡 Examples

### Priority Queue with Order Statistics

```go
// Min-heap style priority queue that also supports "get k-th smallest"
pq := gotreap.NewAutoOrderTreap[int]()

pq.InsertRight(5)
pq.InsertRight(2)
pq.InsertRight(8)
pq.InsertRight(1)

// Get minimum (like heap.Pop)
min, _ := pq.PopLeftmost()  // 1

// Get 2nd smallest element (not possible with standard heap!)
secondSmallest := pq.At(1).Value()  // 2
```

### Leaderboard System

```go
type Player struct {
    Name  string
    Score int
}

// Descending order by score
leaderboard := gotreap.NewTreap(
    func(a, b Player) bool { return a.Score > b.Score },
)

leaderboard.InsertRight(Player{"Alice", 1000})
leaderboard.InsertRight(Player{"Bob", 850})
leaderboard.InsertRight(Player{"Charlie", 1200})

// Top 3 players
for i := 0; i < 3 && i < leaderboard.Size(); i++ {
    player := leaderboard.At(i).Value()
    fmt.Printf("%d. %s - %d points\n", i+1, player.Name, player.Score)
}

// Find player rank
bobNode, rank := leaderboard.FindLowerBound(Player{"", 850})
fmt.Printf("Bob is rank %d\n", rank+1)
```

### Time-Based Event Log

```go
type Event struct {
    Timestamp time.Time
    Message   string
}

events := gotreap.NewTreap(
    func(a, b Event) bool { return a.Timestamp.Before(b.Timestamp) },
)

// Add events
events.InsertRight(Event{time.Now(), "Server started"})
events.InsertRight(Event{time.Now().Add(5 * time.Minute), "User logged in"})

// Query events in time range
start := time.Now().Add(-1 * time.Hour)
end := time.Now()

startEvent := Event{Timestamp: start}
endEvent := Event{Timestamp: end}

count := events.CountRange(startEvent, true, endEvent, true)
fmt.Printf("Events in last hour: %d\n", count)
```

### Interval Merging

```go
type Interval struct {
    Start, End int
}

intervals := gotreap.NewTreap(
    func(a, b Interval) bool { return a.Start < b.Start },
    Interval{1, 3},
    Interval{2, 6},
    Interval{8, 10},
    Interval{15, 18},
)

// Process intervals by iterating
prev := intervals.Leftmost()
for node := prev.Next(); node != nil; node = node.Next() {
    if prev.Value().End >= node.Value().Start {
        // Overlapping intervals - merge logic here
        fmt.Printf("Merge [%d,%d] and [%d,%d]\n",
            prev.Value().Start, prev.Value().End,
            node.Value().Start, node.Value().End)
    }
    prev = node
}
```

---

## 🧪 Testing

Run the test suite:

```bash
go test -v
```

All operations are thoroughly tested with:

- ✅ Edge cases (empty treaps, single elements)
- ✅ Boundary conditions (negative indices, out-of-bounds)
- ✅ Stress tests (1000+ operations)
- ✅ Parent pointer integrity
- ✅ BST and heap invariants
- ✅ Panic conditions

---

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

1. 🐛 **Report bugs** - Open an issue with reproduction steps
2. 💡 **Suggest features** - Share your ideas for improvements
3. 📝 **Improve docs** - Fix typos or add examples
4. 🔧 **Submit PRs** - Follow the existing code style

### Development Guidelines

- Write tests for new features
- Maintain backward compatibility
- Update documentation
- Run `go fmt` and `go vet`
- Ensure all tests pass

---

## 📄 License

CC0-1.0 License - see [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

Based on the classic Treap data structure introduced by Cecilia Aragon and Raimund Seidel (1989).

---

## 📞 Support

- 📖 [Documentation](https://pkg.go.dev/github.com/strokovok/gotreap)
- 🐛 [Issue Tracker](https://github.com/strokovok/gotreap/issues)
- 💬 Discussions welcome via GitHub Issues

---

**Made with ❤️ for the Go community**

_Star ⭐ this repository if you find it useful!_
