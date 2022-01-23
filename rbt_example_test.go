package rbt_test

import (
	"fmt"

	"gotest.com/rbt"
)

func Example() {
	tree := &rbt.Tree[int]{}

	tree.Insert(2)
	tree.Insert(4)
	tree.Insert(7)
	tree.Insert(9)

	tree.Delete(7)

	h := tree.Height()
	fmt.Println(h)
	// Output:
	// 2
}
