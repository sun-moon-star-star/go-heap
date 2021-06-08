package go_heap

import (
	"sync"
)

type SyncHeap struct {
	lock sync.Mutex
	cond *sync.Cond

	heap   Heap
	maxLen int

	unfinishedTaskCnt int
}

func NewSyncHeap(priority func(i, j interface{}) bool) *SyncHeap {
	heap := &SyncHeap{
		heap: Heap{
			priority: priority,
		},
	}
	heap.cond = sync.NewCond(&sync.Mutex{})
	return heap
}

func (heap *SyncHeap) SetMaxLen(maxLen int) {
	heap.lock.Lock()
	heap.maxLen = maxLen
	heap.lock.Unlock()
}

func (heap *SyncHeap) Push(item interface{}) {
	heap.cond.L.Lock()
	for heap.maxLen != -1 && heap.heap.Len() >= heap.maxLen {
		heap.cond.Wait()
	}

	heap.heap.Push(item)
	heap.unfinishedTaskCnt++
	heap.cond.Signal()

	heap.cond.L.Unlock()
}

func (heap *SyncHeap) TryPush(item interface{}) bool {
	if heap.maxLen != -1 && heap.heap.Len() >= heap.maxLen {
		return false
	}

	heap.cond.L.Lock()
	defer heap.cond.L.Unlock()
	for heap.maxLen != -1 && heap.heap.Len() >= heap.maxLen {
		return false
	}

	heap.heap.Push(item)
	heap.unfinishedTaskCnt++
	heap.cond.Signal()

	return true
}

func (heap *SyncHeap) Top() (interface{}, bool) {
	if heap.heap.Len() == 0 {
		return nil, false
	}

	heap.cond.L.Lock()
	defer heap.cond.L.Unlock()
	return heap.heap.Top()
}

func (heap *SyncHeap) Remove() (interface{}, bool) {
	heap.cond.L.Lock()
	for heap.heap.Len() == 0 && heap.unfinishedTaskCnt > 0 {
		heap.cond.Wait()
	}

	if heap.unfinishedTaskCnt <= 0 {
		heap.cond.L.Unlock()
		return nil, false
	}

	item := heap.heap.Pop()
	heap.cond.L.Unlock()
	return item, true
}

func (heap *SyncHeap) TryRemove() (interface{}, bool) {
	heap.cond.L.Lock()
	defer heap.cond.L.Unlock()
	if heap.heap.Len() == 0 {
		return nil, false
	}

	item := heap.heap.Pop()
	return item, true
}

func (heap *SyncHeap) Len() int {
	var len int
	heap.lock.Lock()
	len = heap.heap.Len()
	heap.lock.Unlock()
	return len
}

func (heap *SyncHeap) TaskDone() {
	heap.cond.L.Lock()
	newCnt := heap.unfinishedTaskCnt - 1
	heap.unfinishedTaskCnt = newCnt
	heap.cond.L.Unlock()

	if newCnt <= 0 {
		heap.cond.Broadcast()
	}
}

func (heap *SyncHeap) Join() {
	heap.cond.L.Lock()
	for heap.unfinishedTaskCnt > 0 {
		heap.cond.Wait()
	}
	heap.cond.L.Unlock()
}
