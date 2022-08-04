package minheap_test

import (
	"fmt"
	"testing"

	"git.ultraware.nl/martin/temp/minheap"
)

func TestMinHeap(t *testing.T) {
	heap := minheap.New[int]()

	fmt.Println(heap)
	heap.Insert(1, 1)
	fmt.Println(heap)
	heap.Insert(-1, 2)
	fmt.Println(heap)
	heap.Insert(-3, 3)
	fmt.Println(heap)
	heap.Insert(-4, 4)
	fmt.Println(heap)
	heap.Insert(2, 5)
	fmt.Println(heap)
	heap.Insert(3, 6)
	fmt.Println(heap)

	for _, expected := range [...]int{4, 3, 2, 1, 5, 6} {
		v, ok := heap.Extract()
		fmt.Println(heap)

		if !ok {
			t.Error(ok)
		}
		fmt.Println(v)
		if v != expected {
			t.Error(v, expected)
		}
	}

	if _, ok := heap.Extract(); ok {
		t.Error(ok)
	}
}
