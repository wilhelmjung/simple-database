package main

import (
	"fmt"
	"log"
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
//const NumCell = 16
const NumCell = 4

// Node :
// (child, {key, val}),
type Node struct {
	Parent *Node  `json:"parent,omitempty"`
	Cells  []Cell `json:"cells,omitempty"` // Num + 1 cells, // slice
	Used   int    `json:"used,omitempty"`
}

// NewNode :
func NewNode() *Node {
	arr := make([]Cell, NumCell+1, NumCell+1) // pre-allocated?
	return &Node{Parent: nil, Cells: arr, Used: 0}
}

// check if node is full of cells;
func isFull(node *Node) bool {
	return node.Used == NumCell
}

// Cursor : search or insert position
type Cursor struct {
	Node  *Node `json:"node,omitempty"`
	Index int   `json:"index,omitempty"`
}

// Search :
func Search(key int) *Pair {
	found, cursor := searchTree(key, root)
	if found {
		return getKeyVal(cursor)
	}
	return nil
}

// InvalidCursor :
var InvalidCursor = Cursor{nil, -1}

func getKeyVal(cursor Cursor) *Pair {
	return &cursor.Node.Cells[cursor.Index].KeyVal
}

func setKeyVal(cursor Cursor, keyVal *Pair) {
	cursor.Node.Cells[cursor.Index].KeyVal = *keyVal
}

func getChildNode(cursor Cursor) *Node {
	return cursor.Node.Cells[cursor.Index].Child
}

// binary search within node;
func searchTree(key int, node *Node) (bool, Cursor) {
	if node == nil {
		panic("try to search a NIL node!")
	}
	found, cursor := binarySearchNode(node, key)
	if found { // found in current node;
		return true, cursor
	}
	child := getChildNode(cursor)
	if child != nil {
		return searchTree(key, child)
	}
	return false, cursor // return nearest insert position
}

// if found, return cell and index;
// if not found, return nil and the position for insertion.
func binarySearchNode(node *Node, key int) (bool, Cursor) {
	if node == nil {
		panic("try to search a NIL node.")
	}
	if node.Used == 0 {
		return false, Cursor{node, 0}
	}
	fmt.Printf("search %v in node: %v\n", key, *node)
	l, r := 0, node.Used-1
	if key < node.Cells[l].KeyVal.Key {
		return false, Cursor{node, 0} // not found: Cells[l] move afterward
	}
	if key > node.Cells[r].KeyVal.Key {
		return false, Cursor{node, node.Used} // not found: just append
	}
	// binary search
	for l <= r {
		m := (l + r) / 2
		mKey := node.Cells[m].KeyVal.Key
		if mKey == key {
			return true, Cursor{node, m} // found
		} else if mKey < key {
			l = m + 1
		} else if mKey > key {
			r = m - 1
		}
	}
	// not found, in between r and l, insert at l;
	return false, Cursor{node, l}
}

// Insert : a kv pair
func Insert(keyVal *Pair) (bool, error) {
	// find insert position
	found, cursor := searchTree(keyVal.Key, root)
	if found { // found dup key
		kv := getKeyVal(cursor)
		log.Printf("found dup key: %v -> %v\n", *keyVal, *kv)
		setKeyVal(cursor, keyVal) // overwrite dup key
		return true, nil
	}
	if cursor.Node == nil {
		panic("found invalid cursor!")
	}
	// insert this kv pair first to make it really full;
	ok, err := insertIntoNode(cursor, keyVal)
	if !ok {
		log.Printf("failed insertIntoNode - cursor:%v, kv:%v", cursor, keyVal)
		return false, err
	}
	return true, nil
}

// splitNode on middle kv pair;
func splitNode(node *Node) {
	if node.Used < 3 {
		log.Printf("panic: node - %v", node)
		panic("node to be split should have at least 3 kv pairs.")
	}
	// split on middle kv pair
	mid := node.Used / 2
	keyVal := node.Cells[mid].KeyVal
	//node.Cells[mid].KeyVal = Pair{0, nil} // clear kv
	lNode := node
	rNode := NewNode()
	// copy right half of node to rNode
	j := 0
	for i := mid + 1; i <= node.Used; i++ { // include last ptr;
		rNode.Cells[j] = lNode.Cells[i]
		//lNode.Cells[i] = Cell{nil, Pair{0, nil}}
		j++
	}
	// update used
	lNode.Used -= j + 1 // include mid
	rNode.Used += j
	if node.Parent == nil { // split root node
		newRoot := NewNode()
		newRoot.Used = 1
		// set children and kv
		newRoot.Cells[0].Child = lNode   // left child
		newRoot.Cells[0].KeyVal = keyVal // key val
		newRoot.Cells[1].Child = rNode   // right child
		// set parent
		lNode.Parent, rNode.Parent, root = newRoot, newRoot, newRoot
		return
	}
	// insert kv into its parent node
	pNode := node.Parent
	// to find the exact cell that points to current node
	found, cursor := binarySearchNode(pNode, keyVal.Key)
	if !found && (cursor.Index == 0 && cursor.Index == cursor.Node.Used) {
		log.Printf("panic: node: %v, key: %v", pNode, keyVal.Key)
		panic("key is not within range of node.")
	}
	ok, err := insertIntoNode(cursor, &keyVal)
	if !ok {
		log.Printf("insertIntoNode failed, err: %v", err)
		panic("insertIntoNode failed.")
	}
	return
}

// insert kv into node
func insertIntoNode(cursor Cursor, kv *Pair) (bool, error) {
	node := cursor.Node
	if node == nil {
		err := "try to insert into nil node!"
		panic(err)
	}
	if isFull(node) {
		log.Printf("try to insert into a full node: %v, kv: %v", node, kv)
		panic("insert into a full node.")
	}
	idx := cursor.Index
	// TODO	 check node.Cells[idx].Child == nil
	// move cells
	for i := node.Used + 1; i > idx; i-- {
		node.Cells[i] = node.Cells[i-1]
	}
	// set cell
	node.Cells[idx].KeyVal = *kv
	// set used
	node.Used++
	// check full
	if isFull(node) {
		splitNode(node)
	}
	return true, nil
}

// Init :
func Init() {
	if root == nil {
		root = NewNode()
	}
}

// TestDB :
func TestDB() {
	Init()
	kv1 := &Pair{10, []byte{'f', 'o', 'o'}}
	kv2 := &Pair{20, []byte{'b', 'a', 'r'}}
	kv3 := &Pair{30, []byte{'b', 'a', 'r'}}
	var ok bool
	var err error
	var p *Pair
	ok, err = Insert(kv2)
	log.Printf("ok:%v, err:%v", ok, err)
	ok, err = Insert(kv1)
	log.Printf("ok:%v, err:%v", ok, err)
	ok, err = Insert(kv3)
	log.Printf("ok:%v, err:%v", ok, err)

	p = Search(30)
	log.Printf("found p: %v, root: %v", p, root)
	p = Search(10)
	log.Printf("found p: %v, root: %v", p, root)
	p = Search(20)
	log.Printf("found p: %v, root: %v", p, root)
}
