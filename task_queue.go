package go_heap

import "container/heap"

type TaskQueue struct {
	heap *Heap
}

func NewTaskQueue(less func(i, j interface{}) bool) *TaskQueue {
	return &TaskQueue{heap: New(less)}
}

func (t *TaskQueue) Push(task interface{}) {
	heap.Push(t.heap, task)
}

func (t *TaskQueue) Pop() (interface{}, bool) {
	if t.Len() == 0 {
		return nil, false
	}
	return heap.Pop(t.heap), true
}

func (t *TaskQueue) Len() int {
	return t.heap.Len()
}

func (t *TaskQueue) Top() (interface{}, bool) {
	return t.heap.Top()
}
