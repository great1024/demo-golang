package main

type BT struct {
	parent *BT
	keyNum int
	key []int64
	child []*BT
}


