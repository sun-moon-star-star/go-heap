package go_heap

import "testing"

func TestSyncHeap(t *testing.T) {
	priority := func(i, j interface{}) bool {
		return i.(int) < j.(int)
	}
	queue := NewSyncHeap(priority)
	queue.SetMaxLen(-1)
	queue.Push(2)
	queue.Push(5)
	queue.Push(23)
	queue.Push(-13)
	value, exist := queue.Top()
	if !exist || value != -13 {
		t.Fatal()
	}
}
