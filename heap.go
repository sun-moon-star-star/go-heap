package go_heap

type Heap struct {
	data []interface{}

	priority func(i, j interface{}) bool
}

func New(priority func(i, j interface{}) bool) *Heap {
	return &Heap{
		priority: priority,
	}
}

func (heap *Heap) Top() (interface{}, bool) {
	if heap.Len() == 0 {
		return nil, false
	}
	return heap.data[0], true
}

func (heap *Heap) Less(i, j int) bool {
	return heap.priority(heap.data[i], heap.data[j])
}

func (heap *Heap) Swap(i, j int) {
	(heap.data)[i], (heap.data)[j] = (heap.data)[j], (heap.data)[i]
}

func (heap *Heap) Len() int {
	return len(heap.data)
}

func (heap *Heap) Pop() (v interface{}) {
	if heap.Len() == 0 {
		return nil
	}
	heap.data, v = (heap.data)[:heap.Len()-1], (heap.data)[heap.Len()-1]
	return
}

func (heap *Heap) Push(v interface{}) {
	heap.data = append(heap.data, v.(int))
}
