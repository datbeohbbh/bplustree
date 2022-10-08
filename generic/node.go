package bplustree

import vector "github.com/datbeohbbh/go-utils/vector"

type records = vector.Vector[*record]
type children = vector.Vector[*node]

type node struct {
	isLeaf       bool
	recordList   records
	childrenList children
	parentNode   *node
	previousNode *node
	nextNode     *node
}

func (cur *node) hasKey(rc *record) bool {
	return cur.recordList.BinarySearch(rc) != -1
}

func (cur *node) setNextNode(nextNode *node) {
	if cur != nil {
		cur.nextNode = nextNode
	}
}

func (cur *node) setPreviousNode(prevNode *node) {
	if cur != nil {
		cur.previousNode = prevNode
	}
}

func (cur *node) setParentNode(parent *node) {
	if cur != nil {
		cur.parentNode = parent
	}
}

// dummy function
func (*node) Less(other *node) bool  { return true }
func (*node) Equal(other *node) bool { return true }
