package bplsutree

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

func printTree(t *node, dep int) {
	if t == nil {
		return
	}
	if t.isLeaf == true || t.isLeaf == false {
		fmt.Printf("level: %d\n", dep)
		for i := 0; i < t.recordList.Size(); i++ {
			fmt.Printf("[key = %s, value = %d]\n", string(t.recordList.Get(i).key), t.recordList.Get(i).value)
		}
	}

	for i := 0; i < t.childrenList.Size(); i++ {
		printTree(t.childrenList.Get(i), dep+1)
	}
}

func TestInsert(t *testing.T) {
	tree := &BPTree{
		rootNode: &node{
			isLeaf: true,
		},
	}

	const MAX_VALUE int32 = 2000000000

	rand.Seed(time.Now().UnixNano())
	st := time.Now()
	rcList := []*record{}
	for i := 0; i <= 1000000; i++ {
		val := rand.Int31n(MAX_VALUE)
		rc := &record{
			key:   []byte(fmt.Sprintf("get value = %d at time = %v", val, time.Now())),
			value: i,
		}
		rcList = append(rcList, rc)

		tree.insert(rc)
	}
	log.Printf("insert took %v", time.Since(st))

	for _, rc := range rcList {
		if tree.findLeaf(rc).hasKey(rc) == false {
			t.Errorf("expected true on finding: %p, value = %d", rc, rc.value)
		}
	}

	// printTree(tree.rootNode, 0)
}
