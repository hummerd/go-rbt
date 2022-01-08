# go-rbt
Red-black tree written in Golang with (go 1.18) generics.

It is just practicing in both - algorithms and go's generics.

Manual for installing go 1.18 [here](https://go.dev/dl/#go1.18beta1)

Example:
``` go
	tree := &rbt.Tree[int]{}

	tree.Insert(2)
	tree.Insert(4)
	tree.Insert(7)
	tree.Insert(9)

	tree.Delete(7)

	h := tree.Height()
	fmt.Println(h)
```
