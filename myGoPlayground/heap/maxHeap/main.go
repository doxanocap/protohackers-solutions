type MaxHeap struct {
	array []int
}

func (h *MaxHeap) Init(arr []int) {
	for _,v := range arr {
		h.Insert(v)
	}
}

func (h *MaxHeap) Insert(key int) {
	h.array = append(h.array, key)
	h.maxHeapifyUp(len(h.array)-1)
}

func (h *MaxHeap) Push(key int) {
	h.array = append(h.array, key)
	h.maxHeapifyUp(len(h.array)-1)
}

func (h *MaxHeap) Pop() int {
	largest := h.array[0]
	ln := len(h.array)
	h.array[0] = h.array[ln-1]
	h.array = h.array[:ln-1]

	h.maxHeapifyDown(0)
	return largest
}

func (h *MaxHeap) maxHeapifyDown(index int) {
	ln := len(h.array)-1
	l,r := left(index), right(index)
	child := 0

	for l <= ln {
		if l == ln {
			child = l
		} else if h.array[l] > h.array[r] {
			child = l
		} else {
			child = r
		}

		if h.array[child] > h.array[index] {
			h.swap(index,child)
			index = child
			l, r = left(index), right(index)
		} else {
			return
		}
	}
}

func (h *MaxHeap) maxHeapifyUp(index int) {
	for h.array[parent(index)] < h.array[index] {
		h.swap(parent(index),index)
		index = parent(index)
	}
}

func (h *MaxHeap) swap(i1, i2 int) {
	h.array[i1], h.array[i2] = h.array[i2], h.array[i1]
}

func parent(i int) int {
	return (i - 1) / 2
}

func left(i int) int {
	return 2 * i + 1
}

func right(i int) int {
	return 2 * i + 2
}

func lastStoneWeight(stones []int) int {
	h := &MaxHeap{}
	h.Init(stones)

	for len(h.array) > 1 {
		r,l := h.Pop(), h.Pop()
		if r != l {
			h.Push(r-l)
		}
	}

	if len(h.array) == 0 {
		return 0
	}
	return h.Pop()
}