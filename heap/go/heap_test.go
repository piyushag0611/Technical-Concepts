package heap

import (
	"fmt"
	"math/rand"
	"testing"
)

func setupHeap() *MinHeap {
	return NewMinHeap()
}

func isHeap(h *MinHeap) bool {
	n := len(h.data)
	if n == 0 || n == 1 {
		return true
	}

	i := 0
	for i < n {
		if 2*i+1 < n {
			if h.data[2*i+1] < h.data[i] {
				return false
			}
		}
		if 2*i+2 < n {
			if h.data[2*i+2] < h.data[i] {
				return false
			}
		}
		i += 1
	}
	return true
}

func generateList(seed int64) []int {
	r := rand.New(rand.NewSource(seed))
	nums := make([]int, 100)
	for i := range nums {
		nums[i] = r.Intn(1000)
	}
	return nums
}

func TestPushingNewElement(t *testing.T) {

	h := setupHeap()

	h.Add(10)
	h.Add(1)
	h.Add(6)
	h.Add(3)
	h.Add(7)

	if h.Peek() != 1 {
		t.Errorf("expected root to be 1, got %d", h.Peek())
	}

	if !isHeap(h) {
		t.Errorf("Array does not correspond to a heap")
	}
}

func TestRemovingElement(t *testing.T) {

	h := setupHeap()

	h.Add(10)
	h.Add(1)
	h.Add(6)
	h.Add(3)
	h.Add(7)
	elem, err := h.Remove()
	if err != nil {
		t.Errorf("Element not removed successfully, error thrown")
	}
	if elem != 1 {
		t.Errorf("expected element to be 1, got %d", elem)
	}
	if !isHeap(h) {
		t.Errorf("Array does not correspond to a heap")
	}
}

func TestRemovingFromEmptyArrayRaisesError(t *testing.T) {
	h := setupHeap()
	_, err := h.Remove()
	if err == nil {
		t.Errorf("expected an error, got nil")
	}
}

func TestHeapifyMethod1(t *testing.T) {
	seeds := []int64{42, 7, 99, 31, 10000}

	for _, seed := range seeds {
		t.Run(fmt.Sprintf("seed=%d", seed), func(t *testing.T) {
			h := NewMinHeap()
			nums := generateList(seed)
			h.HeapifyMethod1(nums)

			if len(h.data) != 100 {
				t.Errorf("expected length 100, got %d", len(h.data))
			}
			if !isHeap(h) {
				t.Error("heap property violated")
			}
		})
	}
}

func TestHeapifyMethod2(t *testing.T) {
	seeds := []int64{42, 7, 99, 31, 10000}

	for _, seed := range seeds {
		t.Run(fmt.Sprintf("seed=%d", seed), func(t *testing.T) {
			h := NewMinHeap()
			nums := generateList(seed)
			h.HeapifyMethod2(nums)

			if len(h.data) != 100 {
				t.Errorf("expected length 100, got %d", len(h.data))
			}
			if !isHeap(h) {
				t.Error("heap property violated")
			}
		})
	}
}

func BenchmarkHeapifyMethod1(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000, 1000000, 10000000}
	for _, n := range sizes {
		data := make([]int, n)
		for i := range data {
			data[i] = rand.Intn(n * 10)
		}
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cp := make([]int, n)
				copy(cp, data)
				NewMinHeap().HeapifyMethod1(cp)
			}
		})
	}
}

func BenchmarkHeapifyMethod2(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000, 1000000, 10000000}
	for _, n := range sizes {
		data := make([]int, n)
		for i := range data {
			data[i] = rand.Intn(n * 10)
		}
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cp := make([]int, n)
				copy(cp, data)
				NewMinHeap().HeapifyMethod2(cp)
			}
		})
	}
}
