package go_heap

import "container/heap"

type Heap struct {
	base *HeapBase
}

func New(priority func(i, j interface{}) bool) *Heap {
	return &Heap{
		base: &HeapBase{
			priority: priority,
		},
	}
}

func (h *Heap) Top() (interface{}, bool) {
	return h.base.Top()
}

func (h *Heap) Len() int {
	return h.base.Len()
}

func (h *Heap) Pop() (interface{}, bool) {
	if h.base.Len() == 0 {
		return nil, false
	}
	return heap.Pop(h.base), true
}

func (h *Heap) Push(v interface{}) {
	heap.Push(h.base, v)
}
