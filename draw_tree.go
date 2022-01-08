package rbt

import (
	"constraints"
	"fmt"
	"io"
	"math"
	"os"

	svg "github.com/ajstarks/svgo"
)

const (
	pad    = 4
	radius = 16
)

// DrawSVGFile generates svg for subtree n to file fileName.
func DrawSVGFile[T constraints.Ordered](fileName string, n *Node[T]) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer f.Close()

	DrawSVG(f, n)

	return nil
}

// DrawSVG generates svg for subtree n to out.
func DrawSVG[T constraints.Ordered](out io.Writer, n *Node[T]) {
	h := n.Height()

	lastRowCount := math.Pow(2, float64(h-1))

	width := (2*radius+pad)*int(lastRowCount) + pad
	height := (2*radius+pad)*h + pad

	canvas := svg.New(out)
	canvas.Start(width, height)

	drawNode(canvas, n, 0, width, 0)

	canvas.End()
}

func drawNode[T constraints.Ordered](canvas *svg.SVG, n *Node[T], left, right, height int) {
	if n == nil {
		return
	}

	h := height + pad + radius
	m := (right + left) / 2

	fill := "gray"
	if n.Red {
		fill = "red"
	}

	if n.Left != nil {
		canvas.Line(m, h, (left+m)/2, h+2*radius, "stroke-width:2;stroke:black")
	}

	if n.Right != nil {
		canvas.Line(m, h, (m+right)/2, h+2*radius, "stroke-width:2;stroke:black")
	}

	canvas.Circle(m, h, radius, "fill:"+fill)
	canvas.Text(m, h+radius/4, fmt.Sprintf("%v", n.Value), "text-anchor:middle;font-size:16px;fill:white")

	drawNode(canvas, n.Left, left, m, h+radius)
	drawNode(canvas, n.Right, m, right, h+radius)
}
