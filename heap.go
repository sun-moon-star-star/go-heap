package go_heap

type Heap struct {
	data []interface{}

	less func(i, j interface{}) bool
}

func New(less func(i, j interface{}) bool) *Heap {
	return &Heap{
		less: less,
	}
}

func (heap *Heap) Less(i, j int) bool {
	return heap.less(heap.data[i], heap.data[j])
}

func (heap *Heap) Swap(i, j int) {
	(heap.data)[i], (heap.data)[j] = (heap.data)[j], (heap.data)[i]
}

func (heap *Heap) Len() int {
	return len(heap.data)
}

func (heap *Heap) Pop() (v interface{}) {
	heap.data, v = (heap.data)[:heap.Len()-1], (heap.data)[heap.Len()-1]
	return
}

func (heap *Heap) Push(v interface{}) {
	heap.data = append(heap.data, v.(int))
}
