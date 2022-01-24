package rbt

import (
	"fmt"
)

// TreeCmp represents red-black tree with more flexible approach using Cmp function.
type TreeCmp[T any] struct {
	Root *NodeCmp[T]
	Cmp func(a, b T) int
}

func (t *TreeCmp[T]) Insert(v T) {
	if t.Root == nil {
		t.Root = &NodeCmp[T]{
			Value: v,
		}
		return
	}

	top := t.Root.insert(v, t.Cmp)

	// insert can replace root - so check it
	if top.Parent == nil {
		t.Root = top
	} else if top.Parent.Parent == nil {
		t.Root = top.Parent
	}
}

func (t *TreeCmp[T]) Delete(v T) bool {
	n := t.Root.Find(v, t.Cmp)
	if n == nil {
		return false
	}

	c := n.delete()

	// delete can replace root - so check it
	if c == nil || c.Parent == nil {
		t.Root = c
	} else if c.Parent.Parent == nil {
		t.Root = c.Parent
	} else if c.Parent.Parent.Parent == nil {
		t.Root = c.Parent.Parent
	}

	return true
}

func (t *TreeCmp[T]) Height() int {
	if t.Root == nil {
		return 0
	}

	return t.Root.Height()
}

type NodeCmp[T any] struct {
	Left   *NodeCmp[T]
	Right  *NodeCmp[T]
	Parent *NodeCmp[T]
	Red    bool
	Value  T
}

// Black returns true if node is black. Nil node is considered black.
func (n *NodeCmp[T]) Black() bool {
	if n == nil {
		return true
	}

	return !n.Red
}

// Find finds node with value v in subtree n.
func (n *NodeCmp[T]) Find(v T, cmp func(a, b T) int) *NodeCmp[T] {
	for n != nil {
		if cmp(n.Value, v) == 0 {
			return n
		} else if cmp(v, n.Value) > 0 {
			n = n.Right
		} else {
			n = n.Left
		}
	}

	return n
}

// Finds node Successor or nil if there is no successor.
func (n *NodeCmp[T]) Successor() *NodeCmp[T] {
	if n == nil {
		return nil
	}

	if n.Right != nil {
		return n.Right.Min()
	}

	p := n.Parent
	for p != nil && n == p.Right {
		n = p
		p = p.Parent
	}

	return p
}

// Finds min (most left) value in tree. Or nil if tree is empty.
func (n *NodeCmp[T]) Min() *NodeCmp[T] {
	if n == nil {
		return nil
	}

	for n.Left != nil {
		n = n.Left
	}

	return n
}

// Finds max (most right) value in tree. Or nil if tree is empty.
func (n *NodeCmp[T]) Max() *NodeCmp[T] {
	if n == nil {
		return nil
	}

	for n.Right != nil {
		n = n.Right
	}

	return n
}

// delete deletes node n from subtree n and then resore broken red-black properties.
func (n *NodeCmp[T]) delete() *NodeCmp[T] {
	if n == nil {
		panic("can not delete nil node")
	}

	var d *NodeCmp[T] // node that will be physically deleted
	if n.Left == nil || n.Right == nil {
		d = n
	} else {
		d = n.Successor()
	}

	var c *NodeCmp[T] // child node that will replace deleted
	if d.Left != nil {
		c = d.Left
	} else {
		c = d.Right
	}

	cfake := c == nil
	if !cfake {
		c.Parent = d.Parent
	} else {
		c = &NodeCmp[T]{
			Red:    false,
			Parent: d.Parent,
		}
	}

	if d.Parent != nil {
		if d.Parent.Left == d {
			d.Parent.Left = c
		} else {
			d.Parent.Right = c
		}
	}

	if d != n {
		n.Value = d.Value
	}

	pp := c
	if !d.Red {
		pp = c.deleteFixup()
	}

	if cfake {
		if c.Parent != nil {
			if c.Parent.Left == c {
				c.Parent.Left = nil
			} else {
				c.Parent.Right = nil
			}
		} else {
			return nil
		}
	}

	return pp
}

// deleteFixup
func (n *NodeCmp[T]) deleteFixup() *NodeCmp[T] {
	for n.Parent != nil && n.Black() {
		if n == n.Parent.Left {
			// case 1 - transform it to case 2, 3 or 4
			r := n.Parent.Right
			if r.Red {
				r.Red = false
				r.Parent.Red = true
				n.Parent.RotateLeft()
				r = n.Parent.Right
			}

			if r.Right.Black() && r.Left.Black() {
				// case 2: turn r to red and repeat fixup for n parent
				r.Red = true
				n = n.Parent
			} else {
				if r.Right.Black() {
					// case 3: r.Right is black
					// transform it to case 4
					r.Left.Red = false
					r.Red = true
					r.RotateRight()
					r = n.Parent.Right
				}

				// case 4: final case
				// copy color from n's parent to r
				// color n's parent and and r's right child to black
				// make left rotation against n's parent
				// case 4 is final step in fixing after that all properties
				// are restored
				r.Red = n.Parent.Red
				n.Parent.Red = false
				r.Right.Red = false
				n.Parent.RotateLeft()
				break
			}
		} else {
			l := n.Parent.Left
			if l.Red {
				l.Red = false
				l.Parent.Red = true
				n.Parent.RotateRight()
				l = n.Parent.Left
			}

			if l.Left.Black() && l.Right.Black() {
				l.Red = true
				n = n.Parent
			} else {
				if l.Left.Black() {
					l.Right.Red = false
					l.Red = true
					l.RotateLeft()
					l = n.Parent.Left
				}

				l.Red = n.Parent.Red
				n.Parent.Red = false
				l.Left.Red = false
				n.Parent.RotateRight()
				break
			}
		}
	}

	n.Red = false
	return n
}

// insert inserts v to search tree and restore broken red-black properties.
// insert returns node that can be new root, or it's parent can be new root.
func (n *NodeCmp[T]) insert(v T, cmp func(a, b T) int) *NodeCmp[T] {
	if n == nil {
		panic("can not insert into nil node")
	}

	var p *NodeCmp[T]

	for n != nil {
		p = n

		if cmp(v, p.Value) > 0 {
			n = n.Right
		} else {
			n = n.Left
		}
	}

	nn := &NodeCmp[T]{
		Value:  v,
		Red:    true,
		Parent: p,
	}

	if cmp(v, p.Value) > 0 {
		p.Right = nn
	} else {
		p.Left = nn
	}

	return nn.insertFixup()
}

// insertFixup restores red-black properties that could be broken after inserting red node.
func (n *NodeCmp[T]) insertFixup() *NodeCmp[T] {
	for n.Parent != nil && n.Parent.Red {
		parentLeft := n.Parent.Parent.Left == n.Parent

		var uncle *NodeCmp[T]
		if parentLeft {
			uncle = n.Parent.Parent.Right
		} else {
			uncle = n.Parent.Parent.Left
		}

		if uncle != nil && uncle.Red {
			// case 1: we got red uncle
			// makes uncle and parent black
			// and repaet fixup for grand parent
			uncle.Red = false
			n.Parent.Red = false
			n.Parent.Parent.Red = true
			n = n.Parent.Parent

			if n.Parent == nil {
				n.Red = false
			}

			continue
		}

		if parentLeft {
			if n.Parent.Right == n {
				// case 2: n is right child
				// make right roatation and go to case 3
				n = n.Parent
				n.RotateLeft()
			}

			// case 3: rotate to right
			// then parent black and sibling red
			n.Parent.Parent.RotateRight()
			n.Parent.Red = false
			n.Parent.Right.Red = true
		} else {
			if n.Parent.Left == n {
				n = n.Parent
				n.RotateRight()
			}

			n.Parent.Parent.RotateLeft()
			n.Parent.Red = false
			n.Parent.Left.Red = true
		}

		if n.Parent.Parent == nil {
			n.Parent.Red = false
		}
	}

	return n
}

// Height returns max height for subtree n.
func (n *NodeCmp[T]) Height() int {
	if n == nil {
		return 0
	}

	lh := n.Left.Height()
	rh := n.Right.Height()

	if lh > rh {
		return lh + 1
	}

	return rh + 1
}

// String returns string representation for node.
func (n *NodeCmp[T]) String() string {
	if n == nil {
		return "<nil>"
	}

	p := ""
	if n.Parent != nil {
		p = "; Parent " + fmt.Sprint(n.Parent.Value)
	}

	c := "b"
	if n.Red {
		c = "r"
	}

	return "Node " + fmt.Sprint(n.Value) + c + p
}

// RotateLeft makes left rotation for node n.
// Left rotation:
//     N    <-
//   B   C
//      D E
// ------------
//     C
//   N   E
//  B D
func (n *NodeCmp[T]) RotateLeft() {
	if n == nil {
		return
	}

	c := n.Right
	if c == nil {
		return
	}

	p := n.Parent
	if p != nil {
		p.ReplaceChild(n, c)
	} else {
		c.Parent = nil
	}

	d := c.Left

	c.SetLeft(n)
	n.SetRight(d)
}

// RotateRight makes right rotation for node n.
// Right rotation:
//     N    ->
//   B   C
//  D E
// ------------
//      B
//   D     N
//        E C
func (n *NodeCmp[T]) RotateRight() {
	if n == nil {
		return
	}

	b := n.Left
	if b == nil {
		return
	}

	p := n.Parent
	if p != nil {
		p.ReplaceChild(n, b)
	} else {
		b.Parent = nil
	}

	e := b.Right

	b.SetRight(n)
	n.SetLeft(e)
}

// ReplaceChild replaces left or right child old with new.
// Old must be left or right child.
func (n *NodeCmp[T]) ReplaceChild(old, new *NodeCmp[T]) {
	if n == nil {
		return
	}

	if n.Left == old {
		n.Left = new
	} else {
		n.Right = new
	}

	if new != nil {
		new.Parent = n
	}
}

// SetLeft sets  l as left child for n.
func (n *NodeCmp[T]) SetLeft(l *NodeCmp[T]) {
	if n == nil {
		return
	}

	n.Left = l
	if l != nil {
		l.Parent = n
	}
}

// SetRight sets r as right child for n.
func (n *NodeCmp[T]) SetRight(r *NodeCmp[T]) {
	if n == nil {
		return
	}

	n.Right = r
	if r != nil {
		r.Parent = n
	}
}
