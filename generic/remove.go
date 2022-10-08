package bplustree

func (bptree *BPTree) remove(rc *record) bool {
	if bptree.Empty() {
		return false
	}

	leaf, hasRecord := bptree.searchRecord(rc)
	if !hasRecord {
		return false
	}

	pos := leaf.recordList.BinarySearch(rc)
	leaf.recordList.Remove(pos)

	if leaf.parentNode == nil {
		return true
	}

	if leaf.recordList.Size() < branchFactor-1 {
		bptree.rebalance(leaf)
	}

	return true
}

func (bptree *BPTree) rebalance(curNode *node) {
	for curNode != nil && curNode.parentNode != nil && curNode.recordList.Size() < branchFactor-1 {
		if bptree.redistribute(curNode) {
			return
		} else {
			curNode = bptree.merge(curNode)
		}
	}

	if bptree.rootNode.recordList.Size() == 0 {
		if bptree.rootNode.childrenList.Size() == 0 {
			bptree.rootNode = nil
		} else {
			bptree.rootNode = bptree.rootNode.childrenList.Back()
		}
	}
}

func (bptree *BPTree) getLeafSibling(curNode *node) (*node, int) {
	parent := curNode.parentNode
	if parent == nil {
		return nil, -1
	}

	for pos := 0; pos+1 < parent.childrenList.Size(); pos++ {
		if parent.childrenList.Get(pos) == nil || parent.childrenList.Get(pos+1) == nil {
			return nil, -1
		}

		if parent.childrenList.Get(pos+1).recordList.Size() < branchFactor-1 {
			return parent.childrenList.Get(pos), pos
		}
	}

	return nil, -1
}

func (bptree *BPTree) getRightSibling(curNode *node) (*node, int) {
	parent := curNode.parentNode
	if parent == nil {
		return nil, -1
	}

	for pos := 0; pos+1 < parent.childrenList.Size(); pos++ {
		if parent.childrenList.Get(pos) == nil || parent.childrenList.Get(pos+1) == nil {
			return nil, -1
		}

		if parent.childrenList.Get(pos).recordList.Size() < branchFactor-1 {
			return parent.childrenList.Get(pos + 1), pos
		}
	}

	return nil, -1
}

func (bptree *BPTree) redistribute(curNode *node) bool {
	if leftSibling, pos := bptree.getLeafSibling(curNode); leftSibling != nil && leftSibling.recordList.Size() > branchFactor-1 {
		toMove := (leftSibling.recordList.Size() - curNode.recordList.Size()) / 2
		for ; toMove > 0; toMove-- {
			curNode.recordList.Insert(0, leftSibling.recordList.Back())
			leftSibling.recordList.PopBack()

			if !curNode.isLeaf {
				curNode.childrenList.Insert(0, leftSibling.childrenList.Back())
				curNode.childrenList.Get(0).setParentNode(curNode)
				leftSibling.childrenList.PopBack()
			}
		}

		if curNode.isLeaf {
			curNode.parentNode.recordList.Set(pos, curNode.recordList.Get(0))
		} else {
			separator := curNode.parentNode.recordList.Get(pos)

			c := curNode.recordList.UpperBound(separator)
			if c == -1 {
				c = curNode.recordList.Size()
			}

			curNode.recordList.Insert(c, separator)

			curNode.parentNode.recordList.Set(pos, curNode.recordList.Get(0))
			curNode.recordList.Remove(0)
		}
		return true
	}

	if rightSibling, pos := bptree.getRightSibling(curNode); rightSibling != nil && rightSibling.recordList.Size() > branchFactor-1 {
		toMove := (rightSibling.recordList.Size() - curNode.recordList.Size()) / 2
		for ; toMove > 0; toMove-- {
			curNode.recordList.PushBack(rightSibling.recordList.Get(0))
			rightSibling.recordList.Remove(0)

			if !curNode.isLeaf {
				curNode.childrenList.PushBack(rightSibling.childrenList.Get(0))
				curNode.childrenList.Back().setParentNode(curNode)
				rightSibling.childrenList.Remove(0)
			}
		}

		if curNode.isLeaf {
			curNode.parentNode.recordList.Set(pos, rightSibling.recordList.Get(0))
		} else {
			separator := curNode.parentNode.recordList.Get(pos)

			c := curNode.recordList.UpperBound(separator)
			if c == -1 {
				c = curNode.recordList.Size()
			}

			curNode.recordList.Insert(c, separator)
			curNode.parentNode.recordList.Set(pos, curNode.recordList.Back())
			curNode.recordList.PopBack()
		}
		return true
	}

	return false
}

func (bptree *BPTree) merge(curNode *node) *node {
	if leftSibling, pos := bptree.getLeafSibling(curNode); leftSibling != nil {
		for i := 0; i < curNode.recordList.Size(); i++ {
			leftSibling.recordList.PushBack(curNode.recordList.Get(i))
			if !leftSibling.isLeaf {
				leftSibling.childrenList.PushBack(curNode.childrenList.Get(i))
				leftSibling.childrenList.Back().setParentNode(leftSibling)
			}
		}

		if !leftSibling.isLeaf {
			leftSibling.childrenList.PushBack(curNode.childrenList.Back())
			leftSibling.childrenList.Back().setParentNode(leftSibling)
		}

		leftSibling.setNextNode(curNode.nextNode)
		curNode.nextNode.setPreviousNode(leftSibling)

		separator := leftSibling.parentNode.recordList.Get(pos)
		leftSibling.parentNode.recordList.Remove(pos)
		leftSibling.parentNode.childrenList.Remove(pos + 1)

		if !leftSibling.isLeaf {
			c := leftSibling.recordList.UpperBound(separator)
			if c == -1 {
				c = leftSibling.recordList.Size()
			}
			leftSibling.recordList.Insert(c, separator)
		}

		return leftSibling.parentNode
	}

	if rightSibling, pos := bptree.getRightSibling(curNode); rightSibling != nil {
		for i := 0; i < rightSibling.recordList.Size(); i++ {
			curNode.recordList.PushBack(rightSibling.recordList.Get(i))

			if !curNode.isLeaf {
				curNode.childrenList.PushBack(rightSibling.childrenList.Get(i))
				curNode.childrenList.Back().setParentNode(curNode)
			}
		}

		if !curNode.isLeaf {
			curNode.childrenList.PushBack(rightSibling.childrenList.Back())
			curNode.childrenList.Back().setParentNode(curNode)
		}

		rightSibling.nextNode.setPreviousNode(curNode)
		curNode.setNextNode(rightSibling.nextNode)

		separator := curNode.parentNode.recordList.Get(pos)
		curNode.parentNode.recordList.Remove(pos)
		curNode.parentNode.childrenList.Remove(pos + 1)

		if !curNode.isLeaf {
			c := curNode.recordList.UpperBound(separator)
			curNode.recordList.Insert(c, separator)
		}

		return curNode.parentNode
	}

	return nil
}
