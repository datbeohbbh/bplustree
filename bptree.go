package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

const t int = 400

type Node struct {
	leaf   bool
	keyNum int
	key    [2*t + 1]int
	parent *Node
	child  [2*t + 1]*Node
	value  [2*t + 1]int
	left   *Node
	right  *Node
}

type BPTree struct {
	root *Node
}

func (bptree *BPTree) emptyTree() bool {
	return (bptree.root == nil)
}

func (bptree *BPTree) findLeaf(key int) *Node {
	cur := bptree.root
	for cur != nil && cur.leaf == false {
		for i := 0; i <= cur.keyNum; i++ {
			if i == cur.keyNum || key < cur.key[i] {
				cur = cur.child[i]
				break
			}
		}
	}

	return cur
}

func leafHasKey(L *Node, key int) bool {
	if L == nil {
		return false
	}

	lo, hi := 0, L.keyNum-1
	for lo <= hi {
		mid := (lo + hi) >> 1
		if L.key[mid] < key {
			lo = mid + 1
		} else if L.key[mid] > key {
			hi = mid - 1
		} else {
			return true
		}
	}

	return false
}

func (bptree *BPTree) searchKey(key int) (*Node, bool) {
	leaf := bptree.findLeaf(key)
	return leaf, leafHasKey(leaf, key)
}

func (bptree *BPTree) insert(key, value int) bool {
	if bptree.emptyTree() {
		bptree.root = &Node{leaf: true}
	}
	leaf, has := bptree.searchKey(key)
	if has {
		return false
	}

	if leaf == nil {
		leaf = &Node{}
	}

	pos := 0
	for pos < leaf.keyNum && leaf.key[pos] < key {
		pos += 1
	}

	for i := leaf.keyNum; i >= pos+1; i-- {
		leaf.key[i] = leaf.key[i-1]
		leaf.value[i] = leaf.value[i-1]
	}
	leaf.key[pos] = key
	leaf.value[pos] = value
	leaf.keyNum += 1

	if leaf.keyNum == 2*t {
		bptree.split(leaf)
	}
	return true
}

func (bptree *BPTree) split(node *Node) {
	newNode := &Node{}

	newNode.right = node.right
	if node.right != nil {
		node.right.left = newNode
	}
	node.right = newNode
	newNode.left = node

	midKey := node.key[t]
	newNode.keyNum = t - 1
	node.keyNum = t

	for i := 0; i < newNode.keyNum; i++ {
		newNode.key[i] = node.key[i+t+1]
		newNode.value[i] = node.value[i+t+1]
		newNode.child[i] = node.child[i+t+1]
		if newNode.child[i] != nil {
			newNode.child[i].parent = newNode
		}
	}

	newNode.child[newNode.keyNum] = node.child[2*t]
	if newNode.child[newNode.keyNum] != nil {
		newNode.child[newNode.keyNum].parent = newNode
	}

	if node.leaf {
		newNode.keyNum += 1
		newNode.leaf = true

		for i := newNode.keyNum - 1; i >= 1; i-- {
			newNode.key[i] = newNode.key[i-1]
			newNode.value[i] = newNode.value[i-1]
		}

		newNode.key[0] = node.key[t]
		newNode.value[0] = node.value[t]
	}

	if node.parent == nil {
		bptree.root = &Node{}
		bptree.root.key[0] = midKey
		bptree.root.child[0] = node
		bptree.root.child[1] = newNode
		bptree.root.keyNum = 1
		node.parent = bptree.root
		newNode.parent = bptree.root
	} else {
		newNode.parent = node.parent
		parent := node.parent

		pos := 0
		for pos < parent.keyNum && parent.key[pos] < midKey {
			pos += 1
		}

		for i := parent.keyNum; i >= pos+1; i-- {
			parent.key[i] = parent.key[i-1]
		}

		for i := parent.keyNum + 1; i >= pos+2; i-- {
			parent.child[i] = parent.child[i-1]
		}

		parent.key[pos] = midKey
		parent.child[pos+1] = newNode
		parent.keyNum += 1

		if parent.keyNum == 2*t {
			bptree.split(parent)
		}
	}
}

func (bptree *BPTree) remove(key int) bool {
	if bptree.emptyTree() {
		return false
	}
	leaf, found := bptree.searchKey(key)
	if !found {
		return false
	}

	pos := 0
	for pos < leaf.keyNum && leaf.key[pos] != key {
		pos += 1
	}
	for ; pos+1 < leaf.keyNum; pos++ {
		leaf.key[pos] = leaf.key[pos+1]
		leaf.value[pos] = leaf.value[pos+1]
	}
	leaf.keyNum -= 1
	if leaf.parent == nil {
		return true
	}
	if leaf.keyNum < t-1 {
		bptree.rebalance(leaf)
	}
	return true
}

func (bptree *BPTree) rebalance(node *Node) {
	for node != nil && node.parent != nil && node.keyNum < t-1 {
		if bptree.redistribute(node) {
			return
		} else {
			node = bptree.merge(node)
		}
	}
	if bptree.root.keyNum == 0 {
		// fmt.Println(bptree.root.child[0].key)
		bptree.root = bptree.root.child[0]
	}
}

func (bptree *BPTree) getLeftSibling(node *Node) (*Node, int) {
	parent := node.parent

	if parent == nil {
		return nil, -1
	}

	for pos := 0; pos < parent.keyNum; pos++ {
		if parent.child[pos] == nil || parent.child[pos+1] == nil {
			return nil, -1
		}
		if parent.child[pos+1].keyNum < t-1 {
			// fmt.Println("in P", pos, pos+1, parent.child[pos].keyNum, parent.child[pos+1].keyNum)
			return parent.child[pos], pos
		}
	}

	return nil, -1
}

func (bptree *BPTree) getRightSibling(node *Node) (*Node, int) {
	parent := node.parent

	if parent == nil {
		return nil, -1
	}

	for pos := 0; pos < parent.keyNum; pos++ {
		if parent.child[pos] == nil || parent.child[pos+1] == nil {
			return nil, -1
		}
		if parent.child[pos].keyNum < t-1 {
			return parent.child[pos+1], pos
		}
	}

	return nil, -1
}

func (bptree *BPTree) redistribute(node *Node) bool {
	if leftSibling, pos := bptree.getLeftSibling(node); leftSibling != nil && leftSibling.keyNum > t-1 {
		toMove := (leftSibling.keyNum - node.keyNum) / 2
		for ; toMove > 0; toMove-- {
			for i := node.keyNum; i >= 1; i-- {
				node.key[i] = node.key[i-1]
				node.value[i] = node.value[i-1]
			}
			for i := node.keyNum + 1; i >= 1; i-- {
				node.child[i] = node.child[i-1]
			}
			node.key[0] = leftSibling.key[leftSibling.keyNum-1]
			node.value[0] = leftSibling.value[leftSibling.keyNum-1]
			node.child[0] = leftSibling.child[leftSibling.keyNum]

			node.keyNum += 1
			leftSibling.keyNum -= 1

			if !node.leaf {
				if node.child[0] != nil {
					node.child[0].parent = node
				}
			}
		}
		if node.leaf {
			node.parent.key[pos] = node.key[0]
		} else {
			c := 0
			for c < node.keyNum && node.key[c] <= node.parent.key[pos] {
				c += 1
			}
			for i := node.keyNum; i > c; i-- {
				node.key[i] = node.key[i-1]
				node.value[i] = node.value[i-1]
			}
			node.key[c] = node.parent.key[pos]
			node.keyNum += 1

			node.parent.key[pos] = node.key[0]
			for i := 0; i+1 < node.keyNum; i++ {
				node.key[i] = node.key[i+1]
				node.value[i] = node.value[i+1]
			}
			node.keyNum -= 1
		}
		return true
	}

	if rightSibling, pos := bptree.getRightSibling(node); rightSibling != nil && rightSibling.keyNum > t-1 {
		toMove := (rightSibling.keyNum - node.keyNum) / 2
		for ; toMove > 0; toMove-- {
			node.key[node.keyNum] = rightSibling.key[0]
			node.value[node.keyNum] = rightSibling.value[0]
			node.keyNum += 1
			node.child[node.keyNum] = rightSibling.child[0]
			if !node.leaf {
				if node.child[node.keyNum] != nil {
					node.child[node.keyNum].parent = node
				}
			}

			for i := 0; i+1 < rightSibling.keyNum; i++ {
				rightSibling.key[i] = rightSibling.key[i+1]
				rightSibling.value[i] = rightSibling.value[i+1]
			}
			for i := 0; i < rightSibling.keyNum; i++ {
				rightSibling.child[i] = rightSibling.child[i+1]
			}
			rightSibling.keyNum -= 1
		}

		if node.leaf {
			node.parent.key[pos] = rightSibling.key[0]
		} else {
			c := 0
			for c < node.keyNum && node.key[c] <= node.parent.key[pos] {
				c += 1
			}
			for i := node.keyNum; i > c; i-- {
				node.key[i] = node.key[i-1]
				node.value[i] = node.value[i-1]
			}
			node.key[c] = node.parent.key[pos]
			node.keyNum += 1

			node.parent.key[pos] = node.key[node.keyNum-1]
			node.keyNum -= 1
		}
		return true
	}
	return false
}

func (bptree *BPTree) merge(node *Node) *Node {
	// merge node into left sibling
	if leftSibling, pos := bptree.getLeftSibling(node); leftSibling != nil {
		// fmt.Println("MERGE L", leftSibling.keyNum, leftSibling.key, pos)
		/* 		fmt.Println("in L", leftSibling.keyNum, node.keyNum, leftSibling.parent.keyNum)
		   		fmt.Printf("%p %p\n", leftSibling, node) */
		for i := 0; i < node.keyNum; i++ {

			// fmt.Println(leftSibling.keyNum, node.keyNum)

			leftSibling.key[leftSibling.keyNum] = node.key[i]
			leftSibling.value[leftSibling.keyNum] = node.value[i]

			if !leftSibling.leaf {
				leftSibling.child[leftSibling.keyNum+1] = node.child[i]
				leftSibling.child[leftSibling.keyNum+1].parent = leftSibling
			}
			leftSibling.keyNum += 1
		}

		// remember to move the last child of left sibling.
		if !leftSibling.leaf {
			if node.keyNum < 0 || leftSibling.keyNum+1 < 0 {
				log.Panicf("expected non-negative. found: %d %d\n", node.keyNum, leftSibling.keyNum)
			}
			leftSibling.child[leftSibling.keyNum+1] = node.child[node.keyNum]
			leftSibling.child[leftSibling.keyNum+1].parent = leftSibling
		}

		leftSibling.right = node.right
		if node.right != nil {
			node.right.left = leftSibling
		}

		separator := leftSibling.parent.key[pos]
		for c := pos; c+1 < leftSibling.parent.keyNum; c++ {
			leftSibling.parent.key[c] = leftSibling.parent.key[c+1]
			leftSibling.parent.value[c] = leftSibling.parent.value[c+1]
		}
		for c := pos + 1; c < leftSibling.parent.keyNum; c++ {
			leftSibling.parent.child[c] = leftSibling.parent.child[c+1]
		}

		leftSibling.parent.keyNum -= 1

		if !leftSibling.leaf {
			c := 0
			for c < leftSibling.keyNum && leftSibling.key[c] <= separator {
				c += 1
			}
			// if c-1 >= 0 && leftSibling.key[c-1] != separator {
			for i := leftSibling.keyNum; i > c; i-- {
				leftSibling.key[i] = leftSibling.key[i-1]
			}
			leftSibling.key[c] = separator
			leftSibling.keyNum += 1
			// }
		}
		return leftSibling.parent
	}

	// merge right sibling into node
	if rightSibling, pos := bptree.getRightSibling(node); rightSibling != nil {
		// fmt.Println("MERGE R")
		for i := 0; i < rightSibling.keyNum; i++ {
			node.key[node.keyNum] = rightSibling.key[i]
			node.value[node.keyNum] = rightSibling.value[i]

			if !node.leaf {
				node.child[node.keyNum+1] = rightSibling.child[i]
				if node.child[node.keyNum+1] != nil {
					node.child[node.keyNum+1].parent = node
				}
			}

			node.keyNum += 1
		}

		if !node.leaf {
			node.child[node.keyNum+1] = rightSibling.child[rightSibling.keyNum]
			if node.child[node.keyNum+1] != nil {
				node.child[node.keyNum+1].parent = node
			}
		}

		if rightSibling.right != nil {
			rightSibling.right.left = node
		}
		node.right = rightSibling.right

		separator := node.parent.key[pos]
		// fmt.Println(separator)
		for c := pos; c+1 < node.parent.keyNum; c++ {
			node.parent.key[c] = node.parent.key[c+1]
			node.parent.value[c] = node.parent.value[c+1]
		}
		for c := pos + 1; c < node.parent.keyNum; c++ {
			node.parent.child[c] = node.parent.child[c+1]
		}
		node.parent.keyNum -= 1

		if !node.leaf {
			// fmt.Println(separator)
			c := 0
			for c < node.keyNum && node.key[c] <= separator {
				c += 1
			}
			// if c-1 >= 0 && node.key[c-1] != separator {
			for i := node.keyNum; i > c; i-- {
				node.key[i] = node.key[i-1]
			}
			node.key[c] = separator
			node.keyNum += 1
			// fmt.Println(node.keyNum, node.key, c)
			// }
		}

		return node.parent
	}

	return nil
}

var a []int

func printTree(t *Node, dep int) {
	if t == nil {
		return
	}
	if t.leaf == true || t.leaf == false {
		fmt.Printf("level: %d\n", dep)
		for i := 0; i < t.keyNum; i++ {
			fmt.Printf("[key = %d, value = %d]\n", t.key[i], t.value[i])
			a[t.key[i]] = 1
		}
	}

	for i := 0; i <= t.keyNum; i++ {
		printTree(t.child[i], dep+1)
	}
}

func main() {
	tree := &BPTree{
		root: &Node{
			leaf: true,
		},
	}

	MAX_N := 5000000

	var MAX_VALUE int32 = 2000000000

	ma := make(map[int]bool)

	rand.Seed(time.Now().UnixNano())
	st := time.Now()
	for i := 1; i <= MAX_N; i++ {
		val := int(rand.Int31n(1 + MAX_VALUE))

		for ma[val] == true {
			val = int(rand.Int31n(1 + MAX_VALUE))
		}

		tree.insert(val, val)
		ma[val] = true
	}
	en := time.Now().Sub(st)

	fmt.Println("insert took:", en)

	for k, _ := range ma {
		if tree.insert(k, k) == true {
			log.Panicf("expected: false, found: true. %d already be inserted\n", k)
		}
	}

	st = time.Now()
	for k, _ := range ma {
		if leafHasKey(tree.findLeaf(k), k) == false {
			log.Panicf("failed on check: %d\n", k)
		}
	}
	en = time.Now().Sub(st)
	fmt.Println("search took:", en)

	st = time.Now()
	testIdx := 0
	for k, _ := range ma {
		if leafHasKey(tree.findLeaf(k), k) == true && tree.remove(k) == false {
			log.Panicf("failed on remove test #%d: %d\n", testIdx+1, k)
		} else {
			// log.Printf("ok on remove test #%d: %d\n", testIdx+1, k)
		}
		testIdx += 1
	}
	en = time.Now().Sub(st)
	fmt.Println("remove took:", en)
	/*
		for i := 1; i <= 10; i++ {
			fmt.Println(i, tree.insert(i, i))
		}

		printTree(tree.root, 0) */

	/*
		 	ins := []int{7, 8, 14, 20, 21, 27, 34, 42, 43, 47, 48, 52, 64, 72, 90, 91, 93, 94, 97}
			for _, v := range ins {
				tree.insert(v, v)
			}

			// printTree(tree.root, 0)

			tree.remove(43)
			tree.remove(47)
			tree.remove(7)
			tree.remove(8)
			tree.remove(20)
			tree.remove(21)
			tree.remove(27)
			tree.remove(97)
			tree.remove(48)
			tree.remove(93)
			tree.remove(94)
			tree.remove(64)
			tree.remove(72)

			tree.remove(52)
			tree.remove(34)

			// fmt.Println("DEBUG")
			tree.remove(42)
			tree.remove(91)
			tree.remove(14)
			tree.remove(90)

			printTree(tree.root, 0)
	*/
}
