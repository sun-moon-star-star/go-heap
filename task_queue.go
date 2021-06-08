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

func (t *TaskQueue) Pop() interface{} {
	return heap.Pop(t.heap)
}

func (t *TaskQueue) Len() int {
	return t.heap.Len()
}

func (t *TaskQueue) Top() interface{} {
	return t.heap.Top()
}
