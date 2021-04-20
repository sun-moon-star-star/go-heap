package go_heap

import (
	"errors"
	"sync"
)

type SyncHeap struct {
	lock sync.Mutex
	cond *sync.Cond

	heap   Heap
	maxLen int

	unfinishedTaskCnt int
}

func NewSyncHeap(less func(i, j interface{}) bool) *SyncHeap {
	return &SyncHeap{
		heap: Heap{
			less: less,
		},
	}
}

func (heap *SyncHeap) SetMaxLen(maxLen int) {
	heap.lock.Lock()
	heap.maxLen = maxLen
	heap.lock.Unlock()
}

func (heap *SyncHeap) Append(item interface{}) error {
	var err error
	heap.cond.L.Lock()
	if heap.maxLen != -1 && heap.heap.Len() >= heap.maxLen {
		err = errors.New("Heap is full")
	} else {
		heap.heap.Push(item)
		heap.unfinishedTaskCnt++
		heap.cond.Signal()
		err = nil
	}
	heap.cond.L.Unlock()
	return err
}

func (heap *SyncHeap) Remove() interface{} {
	heap.cond.L.Lock()
	for heap.heap.Len() == 0 && heap.unfinishedTaskCnt > 0 {
		heap.cond.Wait()
	}

	if heap.unfinishedTaskCnt <= 0 {
		heap.cond.L.Unlock()
		return nil
	}

	item := heap.heap.Pop()
	heap.cond.L.Unlock()
	return item
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
