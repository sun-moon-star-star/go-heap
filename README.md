# go_heap

```golang
package main

import (
	"container/heap"
	"fmt"

	"github.com/sun-moon-star-star/go_heap"
)

func main() {
	h := go_heap.New(func(i, j interface{}) bool {
		return i.(int) > j.(int)
	})

	heap.Init(h)

	heap.Push(h, 3)
	heap.Push(h, 1)
	heap.Push(h, 2)
	heap.Push(h, 9)

	len := h.Len()

	element := heap.Pop(h)

	elementInt := element.(int)

	fmt.Println(len, elementInt) // 4, 9
}
```