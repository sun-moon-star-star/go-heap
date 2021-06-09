package go_heap

import (
	"errors"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type Task struct {
	RunUnixNano int64
	Data        interface{}
	CallBack    func(interface{})
}

type PeriodTask struct {
	Queue    *DelayQueue
	Schedule cron.Schedule

	RawData     interface{}
	RawCallBack func(interface{})
	Period      string
}

var globalParser cron.Parser

func init() {
	globalParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
}

func (t *PeriodTask) GetNextRunUnixNano() int64 {
	return t.Schedule.Next(time.Now()).UnixNano()
}

func periodTaskCallBack(data interface{}) {
	t := data.(*PeriodTask)

	t.RawCallBack(t.RawData)

	runUnixNano := t.GetNextRunUnixNano()
	if runUnixNano > 0 {
		t.Queue.Push(&Task{
			RunUnixNano: runUnixNano,
			Data: &PeriodTask{
				Queue:       t.Queue,
				Schedule:    t.Schedule,
				RawData:     t.RawData,
				RawCallBack: t.RawCallBack,
				Period:      t.Period,
			},
			CallBack: periodTaskCallBack,
		})
	}
}

func GetPeriodTask(queue *DelayQueue, period string, task *Task) (*Task, error) {
	schedule, err := globalParser.Parse(period)
	if err != nil {
		return nil, err
	}

	return &Task{
		RunUnixNano: task.RunUnixNano,
		Data: &PeriodTask{
			Queue:    queue,
			Schedule: schedule,

			RawData:     task.Data,
			RawCallBack: task.CallBack,
			Period:      period,
		},
		CallBack: periodTaskCallBack,
	}, nil
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
	isStopping    bool
	isRunning     bool

	lock sync.Mutex
	cond *sync.Cond

	heap   *Heap
	maxLen int
}

func NewDelayQueue() *DelayQueue {
	priority := func(i, j interface{}) bool {
		return i.(*Task).RunUnixNano < j.(*Task).RunUnixNano
	}

	queue := &DelayQueue{
		isAllowedPush: true,
		heap:          New(priority),
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
	queue.cond.Broadcast()

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
	queue.cond.Broadcast()

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
	queue.cond.Broadcast()
}

func (queue *DelayQueue) EndNow() {
	queue.isAllowedPush = false
	queue.isStopping = true
	queue.cond.Broadcast()
}

func (queue *DelayQueue) Run() {
	queue.isRunning = true
	var task *Task

	for !queue.isStopping {
		queue.cond.L.Lock()
		for !queue.isStopping {
			for queue.heap.Len() == 0 && queue.isAllowedPush {
				queue.cond.Wait()
			}

			if queue.isStopping || queue.heap.Len() == 0 {
				queue.cond.L.Unlock()
				queue.isRunning = false
				queue.cond.Broadcast()
				return
			}

			value, ok := queue.heap.Top()
			if ok && value.(*Task).RunUnixNano <= time.Now().UnixNano() {
				value, _ = queue.heap.Pop()
				task = value.(*Task)
				queue.cond.L.Unlock()
				break
			} else if !queue.isStopping {
				unixNano := (value.(*Task).RunUnixNano - time.Now().UnixNano() + 10)
				if unixNano > 0 {
					go func() {
						select {
						case <-time.After(time.Duration(unixNano * int64(time.Nanosecond))):
							queue.cond.Broadcast()
						}
					}()
					queue.cond.Wait()
				}
			} else {
				queue.cond.L.Unlock()
				queue.isRunning = false
				queue.cond.Broadcast()
				return
			}
		}

		if !queue.isStopping {
			task.CallBack(task.Data)
		}
	}
	queue.isRunning = false
	queue.cond.Broadcast()
}

func (queue *DelayQueue) Begin() {
	queue.isRunning = true
	go queue.Run()
}
