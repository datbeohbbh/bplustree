package bplsutree

func (bptree *BPTree) insert(rc *record) bool {
	if bptree.Empty() {
		bptree.rootNode = &node{isLeaf: true}
	}

	leaf, hasRecord := bptree.searchRecord(rc)
	if hasRecord {
		return false
	}

	if leaf == nil {
		leaf = &node{}
	}

	pos := leaf.recordList.LowerBound(rc)
	if pos == -1 {
		pos = leaf.recordList.Size()
	}
	leaf.recordList.Insert(pos, rc)

	if leaf.recordList.Size() == 2*branchFactor {
		bptree.split(leaf)
	}

	return true
}

func (bptree *BPTree) split(curNode *node) {
	newNode := &node{}

	newNode.setNextNode(curNode.nextNode)
	curNode.nextNode.setPreviousNode(newNode)
	curNode.setNextNode(newNode)
	newNode.setPreviousNode(curNode)

	midRecord := curNode.recordList.Get(branchFactor)

	for i := 0; i < branchFactor-1; i++ {
		newNode.recordList.PushBack(curNode.recordList.Get(i + branchFactor + 1))
		if !curNode.isLeaf {
			newNode.childrenList.PushBack(curNode.childrenList.Get(i + branchFactor + 1))
			newNode.childrenList.Back().setParentNode(newNode)
		}
	}

	if !curNode.isLeaf {
		newNode.childrenList.PushBack(curNode.childrenList.Get(2 * branchFactor))
		newNode.childrenList.Back().setParentNode(newNode)
	}

	for i := 0; i < branchFactor-1; i++ {
		curNode.recordList.PopBack()
		if !curNode.isLeaf {
			curNode.childrenList.PopBack()
		}
	}
	curNode.recordList.PopBack()

	if !curNode.isLeaf {
		curNode.childrenList.PopBack()
	}

	if curNode.isLeaf {
		newNode.isLeaf = true
		newNode.recordList.Insert(0, midRecord)
	}

	if curNode == bptree.rootNode { // is root node
		bptree.rootNode = &node{}
		bptree.rootNode.recordList.PushBack(midRecord)

		bptree.rootNode.childrenList.PushBack(curNode)
		bptree.rootNode.childrenList.PushBack(newNode)

		curNode.setParentNode(bptree.rootNode)
		newNode.setParentNode(bptree.rootNode)
	} else {
		newNode.setParentNode(curNode.parentNode)
		parent := curNode.parentNode

		pos := parent.recordList.LowerBound(midRecord)
		if pos == -1 {
			pos = parent.recordList.Size()
		}

		parent.recordList.Insert(pos, midRecord)
		parent.childrenList.Insert(pos+1, newNode)

		if parent.recordList.Size() == 2*branchFactor {
			bptree.split(parent)
		}
	}
}
