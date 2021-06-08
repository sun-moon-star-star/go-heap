package go_heap

import (
	"sync"
)

type SyncHeap struct {
	lock sync.Mutex
	cond *sync.Cond

	heap   *Heap
	maxLen int

	unfinishedTaskCnt int
}

func NewSyncHeap(priority func(i, j interface{}) bool) *SyncHeap {
	heap := &SyncHeap{
		heap: New(priority),
	}
	heap.cond = sync.NewCond(&sync.Mutex{})
	return heap
}

func (syncHeap *SyncHeap) SetMaxLen(maxLen int) {
	syncHeap.lock.Lock()
	syncHeap.maxLen = maxLen
	syncHeap.lock.Unlock()
}

func (syncHeap *SyncHeap) Push(item interface{}) {
	syncHeap.cond.L.Lock()
	for syncHeap.maxLen != -1 && syncHeap.heap.Len() >= syncHeap.maxLen {
		syncHeap.cond.Wait()
	}

	syncHeap.heap.Push(item)
	syncHeap.unfinishedTaskCnt++
	syncHeap.cond.Signal()

	syncHeap.cond.L.Unlock()
}

func (syncHeap *SyncHeap) TryPush(item interface{}) bool {
	if syncHeap.maxLen != -1 && syncHeap.heap.Len() >= syncHeap.maxLen {
		return false
	}

	syncHeap.cond.L.Lock()
	defer syncHeap.cond.L.Unlock()
	for syncHeap.maxLen != -1 && syncHeap.heap.Len() >= syncHeap.maxLen {
		return false
	}

	syncHeap.heap.Push(item)
	syncHeap.unfinishedTaskCnt++
	syncHeap.cond.Signal()

	return true
}

func (syncHeap *SyncHeap) Top() (interface{}, bool) {
	if syncHeap.heap.Len() == 0 {
		return nil, false
	}

	syncHeap.cond.L.Lock()
	defer syncHeap.cond.L.Unlock()
	return syncHeap.heap.Top()
}

func (syncHeap *SyncHeap) Remove() (interface{}, bool) {
	syncHeap.cond.L.Lock()
	for syncHeap.heap.Len() == 0 && syncHeap.unfinishedTaskCnt > 0 {
		syncHeap.cond.Wait()
	}

	if syncHeap.unfinishedTaskCnt <= 0 {
		syncHeap.cond.L.Unlock()
		return nil, false
	}

	item, exist := syncHeap.heap.Pop()
	syncHeap.cond.L.Unlock()
	return item, exist
}

func (syncHeap *SyncHeap) TryRemove() (interface{}, bool) {
	syncHeap.cond.L.Lock()
	defer syncHeap.cond.L.Unlock()
	for syncHeap.heap.Len() == 0 {
		return nil, false
	}

	return syncHeap.heap.Pop()
}

func (syncHeap *SyncHeap) Len() int {
	var len int
	syncHeap.lock.Lock()
	len = syncHeap.heap.Len()
	syncHeap.lock.Unlock()
	return len
}

func (syncHeap *SyncHeap) TaskDone() {
	syncHeap.cond.L.Lock()
	newCnt := syncHeap.unfinishedTaskCnt - 1
	syncHeap.unfinishedTaskCnt = newCnt
	syncHeap.cond.L.Unlock()

	if newCnt <= 0 {
		syncHeap.cond.Broadcast()
	}
}

func (syncHeap *SyncHeap) Join() {
	syncHeap.cond.L.Lock()
	for syncHeap.unfinishedTaskCnt > 0 {
		syncHeap.cond.Wait()
	}
	syncHeap.cond.L.Unlock()
}
