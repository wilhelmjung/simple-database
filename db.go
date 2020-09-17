package main

import (
	"fmt"
	"unsafe"
)

// btree root node
var root *Node

// Pair : kv
type Pair struct {
	Key int    `json:"key,omitempty"`
	Val []byte `json:"val,omitempty"`
}

// Cell :
type Cell struct {
	Child  *Node `json:"child,omitempty"`
	KeyVal Pair  `json:"item,omitempty"`
}

// NumCell : num of cells in one page
// 341 * 3 * 8 = 8184, ~ 8196
//const NumCell = 340
const NumCell = 3

// Node :
// (child, {key, val}),
type Node struct {
	Cells []Cell // Num + 1 cells, // slice
	Used  int
}

// NewNode :
func NewNode() Node {
	arr := make([]Cell, NumCell, NumCell)
	return Node{Cells: arr, Used: 0}
}

// Cursor : search or insert position
type Cursor struct {
	Node  *Node
	Index int
}

// Search :
func Search(key int) (*Pair, Cursor) {
	return binSearch(key, root)
}

// InvalidCursor :
var InvalidCursor = Cursor{nil, -1}

// binary search within node;
func binSearch(key int, node *Node) (*Pair, Cursor) {
	if node == nil {
		return nil, InvalidCursor
	}
	var beg = 0
	var end = node.Used - 1
	var mid = -1
	for {
		if key < node.Cells[beg].KeyVal.Key { // search left
			return binSearch(key, node.Cells[beg].Child)
		}
		if key > node.Cells[end].KeyVal.Key { // search right
			return binSearch(key, node.Cells[end].Child)
		}
		mid = (beg + end) / 2
		if mid == beg { // search middle, not found in this node;
			child := node.Cells[mid].Child
			if child == nil {
				return nil, Cursor{node, beg} // not found, but return cursor
			}
			return binSearch(key, child)
		}
		midKey := node.Cells[mid].KeyVal.Key
		if key == midKey {
			return &node.Cells[mid].KeyVal, Cursor{node, mid}
		} else if key > midKey {
			beg = mid
		} else if key < midKey {
			end = mid
		}
	}
	// not possible!
	//return nil
}

func hello() {
	it := NewNode()
	fmt.Printf("hello,world. %v, \n sz: %v\n", it, unsafe.Sizeof(it))
}
