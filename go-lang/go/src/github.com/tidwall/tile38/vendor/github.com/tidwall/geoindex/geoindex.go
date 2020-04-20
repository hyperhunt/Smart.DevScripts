package geoindex

import (
	"fmt"

	"github.com/tidwall/geoindex/child"
)

// Interface is a tree-like structure that contains geospatial data
type Interface interface {
	// Insert an item into the structure
	Insert(min, max [2]float64, data interface{})
	// Delete an item from the structure
	Delete(min, max [2]float64, data interface{})
	// Search the structure for items that intersects the rect param
	Search(
		min, max [2]float64,
		iter func(min, max [2]float64, data interface{}) bool,
	)
	// Scan iterates through all data in tree in no specified order.
	Scan(iter func(min, max [2]float64, data interface{}) bool)
	// Len returns the number of items in tree
	Len() int
	// Bounds returns the minimum bounding box
	Bounds() (min, max [2]float64)
	// Children returns all children for parent node. If parent node is nil
	// then the root nodes should be returned.
	// The reuse buffer is an empty length slice that can optionally be used
	// to avoid extra allocations.
	Children(parent interface{}, reuse []child.Child) (children []child.Child)
}

// Index is a wrapper around Interface that provides extra features like a
// Nearby (kNN) function.
// This can be created like such:
//   var tree = rbang.New()
//   var index = index.Index{tree}
// Now you can use `index` just like tree but with the extra features.
type Index struct {
	tree Interface
}

// Wrap a tree-like geospatial interface.
func Wrap(tree Interface) *Index {
	return &Index{tree}
}

// Insert an item into the index
func (index *Index) Insert(min, max [2]float64, data interface{}) {
	index.tree.Insert(min, max, data)
}

// Search the index for items that intersects the rect param
func (index *Index) Search(
	min, max [2]float64,
	iter func(min, max [2]float64, data interface{}) bool,
) {
	index.tree.Search(min, max, iter)
}

// Delete an item from the index
func (index *Index) Delete(min, max [2]float64, data interface{}) {
	index.tree.Delete(min, max, data)
}

// Children returns all children for parent node. If parent node is nil
// then the root nodes should be returned.
// The reuse buffer is an empty length slice that can optionally be used
// to avoid extra allocations.
func (index *Index) Children(parent interface{}, reuse []child.Child) (
	children []child.Child,
) {
	return index.tree.Children(parent, reuse)
}

// Nearby performs a kNN-type operation on the index.
// It's expected that the caller provides its own the `algo` function, which
// is used to calculate a distance to data. The `add` function should be
// called by the caller to "return" the data item along with a distance.
// The `iter` function will return all items from the smallest dist to the
// largest dist.
// Take a look at the SimpleBoxAlgo function for a usage example.
func (index *Index) Nearby(
	algo func(
		min, max [2]float64, data interface{}, item bool,
		add func(min, max [2]float64, data interface{}, item bool, dist float64),
	),
	iter func(min, max [2]float64, data interface{}, dist float64) bool,
) {
	var q queue
	var parent interface{}
	var children []child.Child

	var added []qnode
	add := func(min, max [2]float64, data interface{}, item bool, dist float64) {
		added = append(added, qnode{
			dist: dist,
			child: child.Child{
				Data: data,
				Min:  min,
				Max:  max,
				Item: item,
			},
		})
	}

	for {
		// gather all children for parent
		children = index.tree.Children(parent, children[:0])
		for _, child := range children {
			added = added[:0]
			algo(child.Min, child.Max, child.Data, child.Item, add)
			for _, node := range added {
				q.push(node)
			}
		}
		for {
			node, ok := q.pop()
			if !ok {
				// nothing left in queue
				return
			}
			if node.child.Item {
				if !iter(node.child.Min, node.child.Max,
					node.child.Data, node.dist) {
					return
				}
			} else {
				// gather more children
				parent = node.child.Data
				break
			}
		}
	}
}

// Len returns the number of items in tree
func (index *Index) Len() int {
	return index.tree.Len()
}

// Bounds returns the minimum bounding box
func (index *Index) Bounds() (min, max [2]float64) {
	return index.tree.Bounds()
}

// Priority Queue ordered by dist (smallest to largest)

type qnode struct {
	dist  float64
	child child.Child
}

type queue struct {
	nodes []qnode
	len   int
	size  int
}

func (q *queue) push(node qnode) {
	if q.nodes == nil {
		q.nodes = make([]qnode, 2)
	} else {
		q.nodes = append(q.nodes, qnode{})
	}
	i := q.len + 1
	j := i / 2
	for i > 1 && q.nodes[j].dist > node.dist {
		q.nodes[i] = q.nodes[j]
		i = j
		j = j / 2
	}
	q.nodes[i] = node
	q.len++
}

func (q *queue) pop() (qnode, bool) {
	if q.len == 0 {
		return qnode{}, false
	}
	n := q.nodes[1]
	q.nodes[1] = q.nodes[q.len]
	q.len--
	var j, k int
	i := 1
	for i != q.len+1 {
		k = q.len + 1
		j = 2 * i
		if j <= q.len && q.nodes[j].dist < q.nodes[k].dist {
			k = j
		}
		if j+1 <= q.len && q.nodes[j+1].dist < q.nodes[k].dist {
			k = j + 1
		}
		q.nodes[i] = q.nodes[k]
		i = k
	}
	return n, true
}

// Scan iterates through all data in tree in no specified order.
func (index *Index) Scan(
	iter func(min, max [2]float64, data interface{}) bool,
) {
	index.tree.Scan(iter)
}

// SimpleBoxAlgo performs box-distance algorithm on rectangles.
func SimpleBoxAlgo(targetMin, targetMax [2]float64) (
	algo func(
		min, max [2]float64, data interface{}, item bool,
		add func(min, max [2]float64, data interface{}, item bool, dist float64),
	),
) {
	return func(
		min, max [2]float64, data interface{}, item bool,
		add func(min, max [2]float64, data interface{}, item bool, dist float64),
	) {
		add(min, max, data, item, boxDist(targetMin, targetMax, min, max))
	}
}

func boxDist(amin, amax, bmin, bmax [2]float64) float64 {
	var dist float64
	var min, max float64
	if amin[0] > bmin[0] {
		min = amin[0]
	} else {
		min = bmin[0]
	}
	if amax[0] < bmax[0] {
		max = amax[0]
	} else {
		max = bmax[0]
	}
	squared := min - max
	if squared > 0 {
		dist += squared * squared
	}
	if amin[1] > bmin[1] {
		min = amin[1]
	} else {
		min = bmin[1]
	}
	if amax[1] < bmax[1] {
		max = amax[1]
	} else {
		max = bmax[1]
	}
	squared = min - max
	if squared > 0 {
		dist += squared * squared
	}
	return dist
}

func (index *Index) svg(child child.Child, height int) []byte {
	var out []byte
	point := true
	for i := 0; i < 2; i++ {
		if child.Min[i] != child.Max[i] {
			point = false
			break
		}
	}
	if point { // is point
		out = append(out, fmt.Sprintf(
			"<rect x=\"%.0f\" y=\"%0.f\" width=\"%0.f\" height=\"%0.f\" "+
				"stroke=\"%s\" fill=\"purple\" "+
				"fill-opacity=\"0\" stroke-opacity=\"1\" "+
				"rx=\"15\" ry=\"15\"/>\n",
			(child.Min[0])*svgScale,
			(child.Min[1])*svgScale,
			(child.Max[0]-child.Min[0]+1/svgScale)*svgScale,
			(child.Max[1]-child.Min[1]+1/svgScale)*svgScale,
			strokes[height%len(strokes)])...)
	} else { // is rect
		out = append(out, fmt.Sprintf(
			"<rect x=\"%.0f\" y=\"%0.f\" width=\"%0.f\" height=\"%0.f\" "+
				"stroke=\"%s\" fill=\"purple\" "+
				"fill-opacity=\"0\" stroke-opacity=\"1\"/>\n",
			(child.Min[0])*svgScale,
			(child.Min[1])*svgScale,
			(child.Max[0]-child.Min[0]+1/svgScale)*svgScale,
			(child.Max[1]-child.Min[1]+1/svgScale)*svgScale,
			strokes[height%len(strokes)])...)
	}
	if !child.Item {
		children := index.tree.Children(child.Data, nil)
		for _, child := range children {
			out = append(out, index.svg(child, height+1)...)
		}
	}
	return out
}

const (
	// Continue to first child rectangle and/or next sibling.
	Continue = iota
	// Ignore child rectangles but continue to next sibling.
	Ignore
	// Stop iterating
	Stop
)

const svgScale = 4.0

var strokes = [...]string{"black", "#cccc00", "green", "red", "purple"}

// SVG prints 2D rtree in wgs84 coordinate space
func (index *Index) SVG() string {
	var out string
	out += fmt.Sprintf("<svg viewBox=\"%.0f %.0f %.0f %.0f\" "+
		"xmlns =\"http://www.w3.org/2000/svg\">\n",
		-190.0*svgScale, -100.0*svgScale,
		380.0*svgScale, 190.0*svgScale)

	out += fmt.Sprintf("<g transform=\"scale(1,-1)\">\n")

	var outb []byte
	for _, child := range index.Children(nil, nil) {
		outb = append(outb, index.svg(child, 1)...)
	}

	out += string(outb)
	out += fmt.Sprintf("</g>\n")
	out += fmt.Sprintf("</svg>\n")
	return out
}
