package go_heap

import (
	"errors"
	"sync"
	"time"
)

type Task struct {
	RunUnixNano int64
	Data        interface{}
	CallBack    func(interface{})
}

var globalDelayQueue *DelayQueue

func GlobalDelayQueue() *DelayQueue {
	return globalDelayQueue
}

func SetGlobalDelayQueue(queue *DelayQueue) error {
	if globalDelayQueue != nil {
		return errors.New("GlobalDelayQueue is already set")
	}
	globalDelayQueue = queue
	return nil
}

type DelayQueue struct {
	isAllowedPush bool
	isRunning     bool

	lock sync.Mutex
	cond *sync.Cond

	heap   Heap
	maxLen int
}

func NewDelayQueue(priority func(i, j interface{}) bool) *DelayQueue {
	queue := &DelayQueue{
		isAllowedPush: true,
		heap: Heap{
			priority: priority,
		},
	}
	queue.cond = sync.NewCond(&sync.Mutex{})
	return queue
}

func (queue *DelayQueue) SetMaxLen(maxLen int) {
	queue.lock.Lock()
	queue.maxLen = maxLen
	queue.lock.Unlock()
}

func (queue *DelayQueue) Push(item *Task) {
	queue.cond.L.Lock()
	for queue.maxLen != -1 && queue.heap.Len() >= queue.maxLen {
		queue.cond.Wait()
	}

	queue.heap.Push(item)
	queue.cond.Signal()

	queue.cond.L.Unlock()
}

func (queue *DelayQueue) TryPush(item *Task) bool {
	if queue.maxLen != -1 && queue.heap.Len() >= queue.maxLen {
		return false
	}

	queue.cond.L.Lock()
	defer queue.cond.L.Unlock()
	for queue.maxLen != -1 && queue.heap.Len() >= queue.maxLen {
		return false
	}

	queue.heap.Push(item)
	queue.cond.Signal()

	return true
}

func (queue *DelayQueue) Top() (*Task, bool) {
	if queue.heap.Len() == 0 {
		return nil, false
	}

	queue.cond.L.Lock()
	defer queue.cond.L.Unlock()
	task, ok := queue.heap.Top()
	return task.(*Task), ok
}

func (queue *DelayQueue) Len() int {
	var len int
	queue.lock.Lock()
	len = queue.heap.Len()
	queue.lock.Unlock()
	return len
}

func (queue *DelayQueue) Join() {
	queue.cond.L.Lock()
	for queue.isRunning {
		queue.cond.Wait()
	}
	queue.cond.L.Unlock()
}

func (queue *DelayQueue) End() {
	queue.isAllowedPush = false
	queue.cond.Signal()
}

func (queue *DelayQueue) EndNow() {
	queue.isAllowedPush = false
	queue.isRunning = false
	queue.cond.Signal()
}

func (queue *DelayQueue) Begin() {
	queue.isRunning = true
	var task *Task

	for queue.isRunning {
		queue.cond.L.Lock()
		for queue.isRunning {
			for queue.heap.Len() == 0 && queue.isAllowedPush {
				queue.cond.Wait()
			}

			if !queue.isRunning || queue.heap.Len() == 0 {
				queue.cond.L.Unlock()
				return
			}

			value, ok := queue.heap.Top()
			if ok && value.(*Task).RunUnixNano <= time.Now().UnixNano() {
				task = queue.heap.Pop().(*Task)
				queue.cond.L.Unlock()
				break
			} else if queue.isRunning {
				queue.cond.Wait()
			} else {
				queue.cond.L.Unlock()
				return
			}
		}

		if queue.isRunning {
			task.CallBack(task.Data)
		}
	}
}
