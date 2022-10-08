package bplsutree

const branchFactor int = 10

type BPTree struct {
	rootNode *node
}

func (bptree *BPTree) Empty() bool {
	return bptree.rootNode == nil
}

func (bptree *BPTree) findLeaf(rc *record) *node {
	curNode := bptree.rootNode
	for curNode != nil && !curNode.isLeaf {
		pos := curNode.recordList.UpperBound(rc)
		if pos == -1 {
			curNode = curNode.childrenList.Back()
		} else {
			curNode = curNode.childrenList.Get(pos)
		}

	}
	return curNode
}

func (bptree *BPTree) searchRecord(rc *record) (*node, bool) {
	leaf := bptree.findLeaf(rc)
	return leaf, leaf.hasKey(rc)
}
