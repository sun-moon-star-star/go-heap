package go_heap

import (
	"testing"
	"time"
)

func TestDelayQueue(t *testing.T) {
	queue := NewDelayQueue()
	queue.SetMaxLen(-1)
	var arr []int

	queue.Push(&Task{
		RunUnixNano: time.Now().UnixNano() + 0,
		CallBack: func(interface{}) {
			arr = append(arr, 1)
		},
	})

	queue.Begin()
	queue.Push(&Task{
		RunUnixNano: time.Now().UnixNano() + 100000,
		CallBack: func(interface{}) {
			arr = append(arr, 2)
			queue.Push(&Task{
				RunUnixNano: time.Now().UnixNano() + 100000,
				CallBack: func(interface{}) {
					arr = append(arr, 4)
				},
			})
		},
	})
	queue.Push(&Task{
		RunUnixNano: time.Now().UnixNano() + 5000,
		CallBack: func(interface{}) {
			arr = append(arr, 3)
		},
	})

	periodTask, err := GetPeriodTask(queue, "*/3 * * * * *", &Task{
		RunUnixNano: time.Now().UnixNano(),
		CallBack: func(interface{}) {
			t.Log("hello", time.Now())
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	queue.Push(periodTask)

	queue.End()
	queue.Join()
	if len(arr) != 4 || arr[0] != 1 || arr[1] != 3 || arr[2] != 2 || arr[3] != 4 {
		t.Fatal()
	}
}
