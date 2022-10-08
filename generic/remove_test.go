package bplustree

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestRemoveHandTest(t *testing.T) {
	tree := &BPTree{
		rootNode: &node{
			isLeaf: true,
		},
	}

	ins := []int{7, 8, 14, 20, 21, 27, 34, 42, 43, 47, 48, 52, 64, 72, 90, 91, 93, 94, 97}
	for _, v := range ins {
		rc := &record{
			key:   []byte(fmt.Sprintf("%d", v)),
			value: v,
		}
		tree.insert(rc)
	}

	rm := []int{43, 47, 7, 8, 20, 21, 27, 97, 48, 93, 94, 64, 72, 52, 34, 42, 91, 14, 90}

	for _, v := range rm {
		rc := &record{
			key:   []byte(fmt.Sprintf("%d", v)),
			value: v,
		}
		tree.remove(rc)
	}
}

func TestRemoveStress(t *testing.T) {
	tree := &BPTree{
		rootNode: &node{
			isLeaf: true,
		},
	}
	const MAX_VALUE int32 = 2000000000

	rand.Seed(time.Now().UnixNano())
	rcList := []*record{}

	for i := 0; i <= 10000000; i++ {
		val := rand.Int31n(MAX_VALUE)
		rc := &record{
			key:   []byte(fmt.Sprintf("get value = %d at time = %v", val, time.Now())),
			value: i,
		}
		rcList = append(rcList, rc)

		tree.insert(rc)
	}

	st := time.Now()
	for _, rc := range rcList {
		if tree.findLeaf(rc).hasKey(rc) == true && tree.remove(rc) == false {
			t.Errorf("expected false on remove: %p, value = %d", rc, rc.value)
		}
	}

	if tree.Empty() == false {
		t.Errorf("tree should be empty after remove all record")
	}
	log.Printf("remove took %v", time.Since(st))
}
