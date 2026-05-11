package heap

import (
	"errors"
)

type MinHeap struct {
	data []int
}

func NewMinHeap() *MinHeap {
	return &MinHeap{data: []int{}}
}

func (h *MinHeap) Peek() int {
	return h.data[0]
}

func (h *MinHeap) Add(val int) {
	h.data = append(h.data, val)
	index := len(h.data) - 1
	for index > 0 {
		parent := (index - 1) / 2
		if h.data[parent] <= h.data[index] {
			break
		}
		h.data[parent], h.data[index] = h.data[index], h.data[parent]
		index = parent
	}
}

func (h *MinHeap) Remove() (int, error) {
	if len(h.data) > 0 {
		elementRemoved := h.data[0]
		h.data[0] = h.data[len(h.data)-1]
		h.data = h.data[:len(h.data)-1]
		n := len(h.data)
		i := 0
		for 2*i+1 < n {
			lc := 2*i + 1
			if lc+1 < n {
				minElem, minIndex := h.data[lc], lc
				if h.data[lc+1] < h.data[lc] {
					minElem, minIndex = h.data[lc+1], lc+1
				}
				if minElem < h.data[i] {
					h.data[minIndex], h.data[i] = h.data[i], h.data[minIndex]
					i = minIndex
				} else {
					break
				}
			} else {
				if h.data[lc] < h.data[i] {
					h.data[lc], h.data[i] = h.data[i], h.data[lc]
					i = lc
				} else {
					break
				}
			}
		}
		return elementRemoved, nil
	}
	return -1, errors.New("Cannot remove from an empty heap")
}

func (h *MinHeap) HeapifyMethod1(nums []int) {
	n := len(nums)
	for i := range n {
		h.Add(nums[i])
	}
}

func (h *MinHeap) HeapifyMethod2(nums []int) {
	h.data = nums
	n := len(h.data)
	pos := (n / 2) - 1
	for pos >= 0 {
		i := pos
		for 2*i+1 < n {
			lc := 2*i + 1
			if lc+1 < n {
				minElem, minIndex := h.data[lc], lc
				if h.data[lc+1] < h.data[lc] {
					minElem, minIndex = h.data[lc+1], lc+1
				}
				if minElem < h.data[i] {
					h.data[minIndex], h.data[i] = h.data[i], h.data[minIndex]
					i = minIndex
				} else {
					break
				}
			} else {
				if h.data[lc] < h.data[i] {
					h.data[lc], h.data[i] = h.data[i], h.data[lc]
					i = lc
				} else {
					break
				}
			}
		}
		pos -= 1
	}

}
