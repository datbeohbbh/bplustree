package bplsutree

import "bytes"

type record struct {
	key   []byte
	value int
}

func (rc *record) Less(other *record) bool {
	return bytes.Compare(rc.key, other.key) == -1
}

func (rc *record) Equal(other *record) bool {
	return bytes.Equal(rc.key, other.key)
}
