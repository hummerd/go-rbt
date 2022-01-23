package rbt_test

import (
	"constraints"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"testing"

	"gotest.com/rbt"
)

func TestTreeHeight(t *testing.T) {
	n := 144

	tree := &rbt.Tree[int]{}

	for i := 0; i < n; i++ {
		v := rand.Intn(64)
		t.Log("insert", v)
		tree.Insert(v)

		err := checkTree[int](tree.Root)
		if err != nil {
			t.Fatal(err)
		}
	}

	// rbt.DrawSVGFile("test_insert.svg", tree.Root())
	h := tree.Height()

	if float64(h) > 2*math.Log2(float64(n+1)) {
		t.Fatal("height is too big: ", h)
	}
}

func TestTreeDelete(t *testing.T) {
	n := 144

	tree := &rbt.Tree[int]{}
	vs := []int{}

	for i := 0; i < n; i++ {
		v := rand.Intn(64)
		vs = append(vs, v)
		tree.Insert(v)
	}

	for i, v := range vs {
		t.Log("delete", v)
		tree.Delete(v)

		r := tree.Root
		if i == n-1 && r != nil {
			t.Fatal("non empty tree after all")
		}

		if i < n-1 && r == nil {
			t.Fatal("unexpected empty tree", i)
		}

		err := checkTree(r)
		if err != nil {
			t.Fatal(err)
		}
	}

	// rbt.DrawSVGFile("test_delete.svg", tree.Root)
}

func BenchmarkTreeInsert(b *testing.B) {
	tree := &rbt.Tree[int]{}

	vs := []int{2, 3, 6, 8, 2, 3, 7, 6, 9}

	for i := 0; i < b.N; i++ {
		tree.Insert(vs[i%len(vs)])
	}
}

// checkTree checks that all red-black propwerties are valid for tree n.
// n should be root node for tree.
// checkTree checks:
// 1. Root node is black
// 2. Checks that black height is same for left and right children
// 3. Checks that if node is red it has only black children
func checkTree[T constraints.Ordered](n *rbt.Node[T]) error {
	if n == nil {
		return nil
	}

	if n.Red {
		return errors.New("root not black")
	}

	err := checkNode(n.Left)
	if err != nil {
		return err
	}

	err = checkNode(n.Right)
	if err != nil {
		return err
	}

	return nil
}

func checkNode[T constraints.Ordered](n *rbt.Node[T]) error {
	if n == nil {
		return nil
	}

	if n.Parent != nil && n.Parent.Left != n && n.Parent.Right != n {
		return errors.New("wrong parent")
	}

	if n.Red {
		if !n.Left.Black() {
			return errors.New("left child not black")
		}

		if !n.Right.Black() {
			return errors.New("right child not black")
		}
	}

	_, err := blackHeight(n)
	if err != nil {
		return err
	}

	return nil
}

func blackHeight[T constraints.Ordered](n *rbt.Node[T]) (int, error) {
	if n == nil {
		return 0, nil
	}

	bl, err := blackHeight(n.Left)
	if err != nil {
		return 0, err
	}

	br, err := blackHeight(n.Right)
	if err != nil {
		return 0, err
	}

	if bl != br {
		return 0, fmt.Errorf("black height differs for %v; %d != %d", n.Value, bl, br)
	}

	if n.Black() {
		return bl + 1, nil
	}

	return bl, nil
}
