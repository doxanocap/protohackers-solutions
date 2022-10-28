package main

type Point struct {
	val int
	idx []int
}

type MaxHeap struct {
	points []Point
}

func (h *MaxHeap) Init(arr [][]int) {
	for _, v := range arr {
		h.Push(v)
	}
}

func (h *MaxHeap) Push(v []int) {
	h.points = append(h.points, Point{v[0]*v[0] + v[1]*v[1], v})
	h.minHeapUps(len(h.points) - 1)
}

func (h *MaxHeap) Pop() []int {
	ln := len(h.points)
	minimum := h.points[0]
	h.points[0] = h.points[ln-1]
	h.points = h.points[:ln-1]
	h.minHeapDown(0)
	return minimum.idx
}

func (h *MaxHeap) minHeapDown(index int) {
	ln := len(h.points) - 1
	l, r := left(index), right(index)
	child := 0

	for l <= ln {
		if l == ln {
			child = l
		} else if h.points[l].val > h.points[r].val {
			child = l
		} else {
			child = r
		}

		if h.points[child].val > h.points[index].val {
			h.swap(child, index)
			index = child
			l, r = left(index), right(index)
		} else {
			return
		}
	}

}

func (h *MaxHeap) minHeapUp(index int) {
	for h.points[parent(index)].val < h.points[index].val {
		h.swap(parent(index), index)
		index = parent(index)
	}
}

func (h *MaxHeap) swap(index1, index2 int) {
	h.points[index1], h.points[index2] = h.points[index2], h.points[index1]
}

func parent(index int) int {
	return (index - 1) / 2
}

func left(index int) int {
	return index*2 + 1
}

func right(index int) int {
	return index*2 + 2
}

func kClosest(points [][]int, k int) [][]int {
	h := &MaxHeap{}
	h.Init(points)
	ans := [][]int{}
	for ; k > 0; k-- {
		ans = append(ans, h.Pop())
	}
	return ans
}
