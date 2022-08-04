package minheap

func New[T comparable]() Heap[T] {
	return Heap[T]{
		Map: make(map[T]int),
	}
}

type Node[T comparable] struct {
	Score float64
	Value T
}

type Heap[T comparable] struct {
	Tree []Node[T]
	Map  map[T]int
}

func (h Heap[T]) Len() int {
	return len(h.Tree)
}

func (h *Heap[T]) Extract() (T, bool) {
	if len(h.Tree) == 0 {
		var t T
		return t, false
	}

	out := h.Tree[0].Value

	h.Tree[0] = h.Tree[len(h.Tree)-1] // set last value to first
	h.Tree = h.Tree[:len(h.Tree)-1]   // remove last value

	downHeap(h, 0)

	if len(h.Tree) > 0 {
		h.Map[h.Tree[0].Value] = 0
	}

	delete(h.Map, out)

	return out, true
}

func downHeap[T comparable](h *Heap[T], i int) {
	left := i*2 + 1
	right := i*2 + 2
	smallest := i

	if len(h.Tree) > left && h.Tree[left].Score < h.Tree[smallest].Score {
		smallest = left
	}

	if len(h.Tree) > right && h.Tree[right].Score < h.Tree[smallest].Score {
		smallest = right
	}

	if smallest != i {
		h.Tree[i], h.Tree[smallest] = h.Tree[smallest], h.Tree[i]
		h.Map[h.Tree[i].Value] = i
		h.Map[h.Tree[smallest].Value] = smallest
		downHeap(h, smallest)
	}
}

func (h *Heap[T]) Insert(score float64, value T) {
	if _, ok := h.Map[value]; ok {
		return
	}

	h.Tree = append(h.Tree, Node[T]{score, value})
	h.Map[value] = len(h.Tree) - 1

	upHeap(h, len(h.Tree)-1)
}

func upHeap[T comparable](h *Heap[T], i int) {
	if i == 0 {
		return
	}

	parent := (i - 1) / 2

	if h.Tree[parent].Score > h.Tree[i].Score {
		h.Tree[parent], h.Tree[i] = h.Tree[i], h.Tree[parent]
		h.Map[h.Tree[i].Value] = i
		h.Map[h.Tree[parent].Value] = parent
		upHeap(h, parent)
	}
}
