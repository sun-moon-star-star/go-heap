package go_heap

import (
	"testing"
	"time"
)

func TestDelayQueue(t *testing.T) {
	priority := func(i, j interface{}) bool {
		return i.(*Task).RunUnixNano < j.(*Task).RunUnixNano
	}
	queue := NewDelayQueue(priority)
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

	queue.End()
	queue.Join()
	if len(arr) != 4 || arr[0] != 1 || arr[1] != 3 || arr[2] != 2 || arr[3] != 4 {
		t.Fatal()
	}
}
