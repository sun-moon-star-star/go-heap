package go_heap

import "testing"

func TestTaskQueue(t *testing.T) {
	priority := func(i, j interface{}) bool {
		return i.(int) < j.(int)
	}
	queue := NewTaskQueue(priority)
	queue.Push(2)
	queue.Push(5)
	queue.Push(23)
	queue.Push(-13)
	if queue.Top() != -13 {
		t.Fatal()
	}
}
