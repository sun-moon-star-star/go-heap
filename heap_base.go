package go_heap

type HeapBase struct {
	data []interface{}

	priority func(i, j interface{}) bool
}

func NewHeapBase(priority func(i, j interface{}) bool) *HeapBase {
	return &HeapBase{
		priority: priority,
	}
}

func (heap *HeapBase) Top() (interface{}, bool) {
	if heap.Len() == 0 {
		return nil, false
	}
	return heap.data[0], true
}

func (heap *HeapBase) Less(i, j int) bool {
	return heap.priority(heap.data[i], heap.data[j])
}

func (heap *HeapBase) Swap(i, j int) {
	(heap.data)[i], (heap.data)[j] = (heap.data)[j], (heap.data)[i]
}

func (heap *HeapBase) Len() int {
	return len(heap.data)
}

func (heap *HeapBase) Pop() (v interface{}) {
	heap.data, v = (heap.data)[:heap.Len()-1], (heap.data)[heap.Len()-1]
	return
}

func (heap *HeapBase) Push(v interface{}) {
	heap.data = append(heap.data, v)
}
