package rbt_test

import (
	"testing"

	"gotest.com/rbt"
)

func FuzzMutateTree(f *testing.F) {
	tree := &rbt.Tree[int]{}

	f.Fuzz(func(t *testing.T, add, del int) {
		t.Log("fuzz add/del", add, del)

		tree.Insert(add)

		err := checkTree[int](tree.Root)
		if err != nil {
			t.Fatal(err)
		}

		tree.Delete(del)
		err = checkTree[int](tree.Root)
		if err != nil {
			t.Fatal(err)
		}
	})
}
