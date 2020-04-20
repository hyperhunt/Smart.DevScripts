package rbang

import (
	"github.com/tidwall/geoindex/child"
)

const (
	maxEntries = 32
	minEntries = maxEntries * 40 / 100
)

type rect struct {
	min, max [2]float64
	data     interface{}
}

type node struct {
	count int
	boxes [maxEntries + 1]rect
}

// RTree ...
type RTree struct {
	height   int
	root     rect
	count    int
	reinsert []rect
}

func (r *rect) expand(b *rect) {
	if b.min[0] < r.min[0] {
		r.min[0] = b.min[0]
	}
	if b.max[0] > r.max[0] {
		r.max[0] = b.max[0]
	}
	if b.min[1] < r.min[1] {
		r.min[1] = b.min[1]
	}
	if b.max[1] > r.max[1] {
		r.max[1] = b.max[1]
	}
}

func (r *rect) area() float64 {
	return (r.max[0] - r.min[0]) * (r.max[1] - r.min[1])
}

func (r *rect) overlapArea(b *rect) float64 {
	area := 1.0
	var max, min float64
	if r.max[0] < b.max[0] {
		max = r.max[0]
	} else {
		max = b.max[0]
	}
	if r.min[0] > b.min[0] {
		min = r.min[0]
	} else {
		min = b.min[0]
	}
	if max > min {
		area *= max - min
	} else {
		return 0
	}
	if r.max[1] < b.max[1] {
		max = r.max[1]
	} else {
		max = b.max[1]
	}
	if r.min[1] > b.min[1] {
		min = r.min[1]
	} else {
		min = b.min[1]
	}
	if max > min {
		area *= max - min
	} else {
		return 0
	}
	return area
}

func (r *rect) enlargedArea(b *rect) float64 {
	area := 1.0
	if b.max[0] > r.max[0] {
		if b.min[0] < r.min[0] {
			area *= b.max[0] - b.min[0]
		} else {
			area *= b.max[0] - r.min[0]
		}
	} else {
		if b.min[0] < r.min[0] {
			area *= r.max[0] - b.min[0]
		} else {
			area *= r.max[0] - r.min[0]
		}
	}
	if b.max[1] > r.max[1] {
		if b.min[1] < r.min[1] {
			area *= b.max[1] - b.min[1]
		} else {
			area *= b.max[1] - r.min[1]
		}
	} else {
		if b.min[1] < r.min[1] {
			area *= r.max[1] - b.min[1]
		} else {
			area *= r.max[1] - r.min[1]
		}
	}
	return area
}

// Insert inserts an item into the RTree
func (tr *RTree) Insert(min, max [2]float64, value interface{}) {
	var item rect
	fit(min, max, value, &item)
	tr.insert(&item)
}

func (tr *RTree) insert(item *rect) {
	if tr.root.data == nil {
		fit(item.min, item.max, new(node), &tr.root)
	}
	grown := tr.root.insert(item, tr.height)
	if grown {
		tr.root.expand(item)
	}
	if tr.root.data.(*node).count == maxEntries+1 {
		newRoot := new(node)
		tr.root.splitLargestAxisEdgeSnap(&newRoot.boxes[1])
		newRoot.boxes[0] = tr.root
		newRoot.count = 2
		tr.root.data = newRoot
		tr.root.recalc()
		tr.height++
	}
	tr.count++
}

func (r *rect) chooseLeastEnlargement(b *rect) int {
	j, jenlargement, jarea := -1, 0.0, 0.0
	n := r.data.(*node)
	for i := 0; i < n.count; i++ {
		area := n.boxes[i].area()
		enlargement := n.boxes[i].enlargedArea(b) - area
		if j == -1 || enlargement < jenlargement {
			j, jenlargement, jarea = i, enlargement, area
		} else if enlargement == jenlargement {
			if area < jarea {
				j, jenlargement, jarea = i, enlargement, area
			}
		}
	}
	return j
}

func (r *rect) recalc() {
	n := r.data.(*node)
	r.min = n.boxes[0].min
	r.max = n.boxes[0].max
	for i := 1; i < n.count; i++ {
		r.expand(&n.boxes[i])
	}
}

// contains return struct when b is fully contained inside of n
func (r *rect) contains(b *rect) bool {
	if b.min[0] < r.min[0] || b.max[0] > r.max[0] {
		return false
	}
	if b.min[1] < r.min[1] || b.max[1] > r.max[1] {
		return false
	}
	return true
}

func (r *rect) largestAxis() (axis int, size float64) {
	if r.max[1]-r.min[1] > r.max[0]-r.min[0] {
		return 1, r.max[1] - r.min[1]
	}
	return 0, r.max[0] - r.min[0]
}

func (r *rect) splitLargestAxisEdgeSnap(right *rect) {
	axis, _ := r.largestAxis()
	left := r
	leftNode := left.data.(*node)
	rightNode := new(node)
	right.data = rightNode

	var equals []rect
	for i := 0; i < leftNode.count; i++ {
		minDist := leftNode.boxes[i].min[axis] - left.min[axis]
		maxDist := left.max[axis] - leftNode.boxes[i].max[axis]
		if minDist < maxDist {
			// stay left
		} else {
			if minDist > maxDist {
				// move to right
				rightNode.boxes[rightNode.count] = leftNode.boxes[i]
				rightNode.count++
			} else {
				// move to equals, at the end of the left array
				equals = append(equals, leftNode.boxes[i])
			}
			leftNode.boxes[i] = leftNode.boxes[leftNode.count-1]
			leftNode.boxes[leftNode.count-1].data = nil
			leftNode.count--
			i--
		}
	}
	for _, b := range equals {
		if leftNode.count < rightNode.count {
			leftNode.boxes[leftNode.count] = b
			leftNode.count++
		} else {
			rightNode.boxes[rightNode.count] = b
			rightNode.count++
		}
	}
	left.recalc()
	right.recalc()
}

func (r *rect) insert(item *rect, height int) (grown bool) {
	n := r.data.(*node)
	if height == 0 {
		n.boxes[n.count] = *item
		n.count++
		grown = !r.contains(item)
		return grown
	}
	// choose subtree
	index := r.chooseLeastEnlargement(item)
	child := &n.boxes[index]
	grown = child.insert(item, height-1)
	if grown {
		child.expand(item)
		grown = !r.contains(item)
	}
	if child.data.(*node).count == maxEntries+1 {
		child.splitLargestAxisEdgeSnap(&n.boxes[n.count])
		n.count++
	}
	return grown
}

// fit an external item into a box type
func fit(min, max [2]float64, value interface{}, target *rect) {
	target.min = min
	target.max = max
	target.data = value
}

type overlapsResult int

const (
	not overlapsResult = iota
	intersects
	contains
)

// overlaps detects if r insersects or contains b.
// return not, intersects, contains
func (r *rect) overlaps(b *rect) overlapsResult {
	if b.min[0] > r.max[0] || b.max[0] < r.min[0] {
		return not
	}
	if r.min[0] > b.min[0] || b.max[0] > r.max[0] {
		if b.min[1] > r.max[1] || b.max[1] < r.min[1] {
			return not
		}
		return intersects
	}
	if b.min[1] > r.max[1] || b.max[1] < r.min[1] {
		return not
	}
	if r.min[1] > b.min[1] || b.max[1] > r.max[1] {
		return intersects
	}
	return contains
}

// contains return struct when b is fully contained inside of n
func (r *rect) intersects(b *rect) bool {
	if b.min[0] > r.max[0] || b.max[0] < r.min[0] {
		return false
	}
	if b.min[1] > r.max[1] || b.max[1] < r.min[1] {
		return false
	}
	return true
}

func (r *rect) search(
	target *rect, height int,
	iter func(min, max [2]float64, value interface{}) bool,
) bool {
	n := r.data.(*node)
	if height == 0 {
		for i := 0; i < n.count; i++ {
			if target.intersects(&n.boxes[i]) {
				if !iter(n.boxes[i].min, n.boxes[i].max,
					n.boxes[i].data) {
					return false
				}
			}
		}
	} else {
		for i := 0; i < n.count; i++ {
			switch target.overlaps(&n.boxes[i]) {
			case intersects:
				if !n.boxes[i].search(target, height-1, iter) {
					return false
				}
			case contains:
				if !n.boxes[i].scan(height-1, iter) {
					return false
				}
			}
		}
	}
	return true
}

func (tr *RTree) search(
	target *rect,
	iter func(min, max [2]float64, value interface{}) bool,
) {
	if tr.root.data == nil {
		return
	}
	res := target.overlaps(&tr.root)
	if res == intersects {
		tr.root.search(target, tr.height, iter)
	} else if res == contains {
		tr.root.scan(tr.height, iter)
	}
}

// Search ...
func (tr *RTree) Search(
	min, max [2]float64,
	iter func(min, max [2]float64, value interface{}) bool,
) {
	var target rect
	fit(min, max, nil, &target)
	tr.search(&target, iter)
}

func (r *rect) scan(
	height int,
	iter func(min, max [2]float64, value interface{}) bool,
) bool {
	n := r.data.(*node)
	if height == 0 {
		for i := 0; i < n.count; i++ {
			if !iter(n.boxes[i].min, n.boxes[i].max, n.boxes[i].data) {
				return false
			}
		}
	} else {
		for i := 0; i < n.count; i++ {
			if !n.boxes[i].scan(height-1, iter) {
				return false
			}
		}
	}
	return true
}

// Scan iterates through all data in tree.
func (tr *RTree) Scan(iter func(min, max [2]float64, data interface{}) bool) {
	if tr.root.data == nil {
		return
	}
	tr.root.scan(tr.height, iter)
}

// Delete data from tree
func (tr *RTree) Delete(min, max [2]float64, data interface{}) {
	var item rect
	fit(min, max, data, &item)
	if tr.root.data == nil || !tr.root.contains(&item) {
		return
	}
	var removed, recalced bool
	removed, recalced, tr.reinsert =
		tr.root.delete(&item, tr.height, tr.reinsert[:0])
	if !removed {
		return
	}
	tr.count -= len(tr.reinsert) + 1
	if tr.count == 0 {
		tr.root = rect{}
		recalced = false
	} else {
		for tr.height > 0 && tr.root.data.(*node).count == 1 {
			tr.root = tr.root.data.(*node).boxes[0]
			tr.height--
			tr.root.recalc()
		}
	}
	if recalced {
		tr.root.recalc()
	}
	for i := range tr.reinsert {
		tr.insert(&tr.reinsert[i])
		tr.reinsert[i].data = nil
	}
}

func (r *rect) delete(item *rect, height int, reinsert []rect) (
	removed, recalced bool, reinsertOut []rect,
) {
	n := r.data.(*node)
	if height == 0 {
		for i := 0; i < n.count; i++ {
			if n.boxes[i].data == item.data {
				// found the target item to delete
				recalced = r.onEdge(&n.boxes[i])
				n.boxes[i] = n.boxes[n.count-1]
				n.boxes[n.count-1].data = nil
				n.count--
				if recalced {
					r.recalc()
				}
				return true, recalced, reinsert
			}
		}
	} else {
		for i := 0; i < n.count; i++ {
			if !n.boxes[i].contains(item) {
				continue
			}
			removed, recalced, reinsert =
				n.boxes[i].delete(item, height-1, reinsert)
			if !removed {
				continue
			}
			if n.boxes[i].data.(*node).count < minEntries {
				// underflow
				if !recalced {
					recalced = r.onEdge(&n.boxes[i])
				}
				reinsert = n.boxes[i].flatten(reinsert, height-1)
				n.boxes[i] = n.boxes[n.count-1]
				n.boxes[n.count-1].data = nil
				n.count--
			}
			if recalced {
				r.recalc()
			}
			return removed, recalced, reinsert
		}
	}
	return false, false, reinsert
}

// flatten flattens all leaf boxes into a single list
func (r *rect) flatten(all []rect, height int) []rect {
	n := r.data.(*node)
	if height == 0 {
		all = append(all, n.boxes[:n.count]...)
	} else {
		for i := 0; i < n.count; i++ {
			all = n.boxes[i].flatten(all, height-1)
		}
	}
	return all
}

// onedge returns true when b is on the edge of r
func (r *rect) onEdge(b *rect) bool {
	if r.min[0] == b.min[0] || r.max[0] == b.max[0] {
		return true
	}
	if r.min[1] == b.min[1] || r.max[1] == b.max[1] {
		return true
	}
	return false
}

// Len returns the number of items in tree
func (tr *RTree) Len() int {
	return tr.count
}

// Bounds returns the minimum bounding box
func (tr *RTree) Bounds() (min, max [2]float64) {
	if tr.root.data == nil {
		return
	}
	return tr.root.min, tr.root.max
}

// Children is a utility function that returns all children for parent node.
// If parent node is nil then the root nodes should be returned. The min, max,
// data, and items slices all must have the same lengths. And, each element
// from all slices must be associated. Returns true for `items` when the the
// item at the leaf level. The reuse buffers are empty length slices that can
// optionally be used to avoid extra allocations.
func (tr *RTree) Children(
	parent interface{},
	reuse []child.Child,
) []child.Child {
	children := reuse
	if parent == nil {
		if tr.Len() > 0 {
			// fill with the root
			children = append(children, child.Child{
				Min:  tr.root.min,
				Max:  tr.root.max,
				Data: tr.root.data,
				Item: false,
			})
		}
	} else {
		// fill with child items
		n := parent.(*node)
		item := true
		if n.count > 0 {
			if _, ok := n.boxes[0].data.(*node); ok {
				item = false
			}
		}
		for i := 0; i < n.count; i++ {
			children = append(children, child.Child{
				Min:  n.boxes[i].min,
				Max:  n.boxes[i].max,
				Data: n.boxes[i].data,
				Item: item,
			})
		}
	}
	return children
}
