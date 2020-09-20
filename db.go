package db

import (
	"log"
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
//const NumCell = 16
const NumCell = 4

// Node :
// (child, {key, val}),
type Node struct {
	Parent *Node
	Cells  []Cell // Num + 1 cells, // slice
	Used   int
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
	Node  *Node
	Index int
}

// Search :
func Search(key int) *Pair {
	kv, _ := binSearch(key, root)
	return kv
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
	// empty node
	if end == -1 {
		return nil, Cursor{node, 0}
	}
	for {
		if key < node.Cells[beg].KeyVal.Key { // search left
			return binSearch(key, node.Cells[beg].Child)
		}
		if key > node.Cells[end].KeyVal.Key { // search right
			return binSearch(key, node.Cells[end].Child)
		}
		mid = (beg + end) / 2
		// search middle, not found in this node;
		if mid == beg {
			child := node.Cells[mid].Child
			if child == nil {
				// not found, but return cursor
				return nil, Cursor{node, beg}
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

// if found, return cell and index;
// if not found, return nil and the position for insertion.
func searchInNode(node *Node, key int) (*Cell, int) {
	l, r := 0, node.Used-1
	if key < node.Cells[l].KeyVal.Key {
		return nil, 0 // Cells[l] move afterward
	}
	if key > node.Cells[r].KeyVal.Key {
		return nil, node.Used // just append
	}
	// binary search
	for l < r {
		m := (l + r) / 2
		mKey := node.Cells[m].KeyVal.Key
		if mKey == key {
			return &node.Cells[m], m // found
		} else if mKey < key {
			l = m
		} else if mKey > key {
			r = m
		}
		if l-r == -1 { // in between l and r, so no need to search
			return nil, r
		}
	}
	// used == 1, l == r == 0, and key == cells[0].keyval.key
	return &node.Cells[l], 0
}

// Insert : a kv pair
func Insert(keyVal *Pair) (bool, error) {
	// find insert position
	kv, cursor := binSearch(keyVal.Key, root)
	if kv != nil { // found dup key
		log.Printf("found dup key: @%v with kv %v\n", cursor, kv)
		kv.Val = keyVal.Val // overwrite dup key
		return true, nil
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
	cell, idx := searchInNode(pNode, keyVal.Key)
	if cell == nil {
		log.Printf("panic: node: %v, key: %v", pNode, keyVal.Key)
		panic("key is not within range of node.")
	}
	cur := Cursor{pNode, idx}
	ok, err := insertIntoNode(cur, &keyVal)
	if !ok {
		log.Printf("insertIntoNode failed, err: %v", err)
		panic("insertIntoNode failed.")
	}
	return
}

// insert kv into node
func insertIntoNode(cursor Cursor, kv *Pair) (bool, error) {
	node := cursor.Node
	if isFull(cursor.Node) {
		log.Printf("try to insert into a full node: %v, kv: %v", cursor.Node, kv)
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
	// check full
	if isFull(cursor.Node) {
		splitNode(cursor.Node)
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
	log.Printf("found p: %v", p)
	p = Search(10)
	log.Printf("found p: %v", p)
	p = Search(20)
	log.Printf("found p: %v", p)
}

func hello() {
	it := NewNode()
	log.Printf("hello,world. %v, \n sz: %v\n", it, unsafe.Sizeof(it))
}
