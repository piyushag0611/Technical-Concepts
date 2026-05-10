import pytest
import random
from heap import Heap

def isHeap(elements) -> bool:

    n = len(elements)
    if (n==0):
        return True

    else:

        for i in range(n):
            
            if (2*i+1 < n):
                if (elements[2*i+1] < elements[i]):
                    return False
            
            if (2*i+2 < n):
                if (elements[2*i+2] < elements[i]):
                    return False
        
        return True
            
def isAscending(elements) -> bool:

    for i in range(len(elements)-1):
        if (elements[i+1] < elements[i]):
            return False
    return True



class TestHeap:

    def generate_list(self, seed):

        random.seed(seed)
        return [random.randint(1, 1000) for _ in range(100)]

    def setup_method(self):
        
        self.heap = Heap()
        self.heap.push_elements(10)
        self.heap.push_elements(2)
        self.heap.push_elements(15)
        self.heap.push_elements(18)
        self.heap.push_elements(9)
    
    # Adding an element ensures heap structure is retained
    def test_push_new_elements(self):

        self.heap.push_elements(3)
        assert isHeap(self.heap.elements)
    
    # Adding an element ensures heap structure is retained
    def test_push_existing_elements(self):
        self.heap.push_elements(10)
        assert isHeap(self.heap.elements)

    # Removing en element ensures heap structure is retained
    def test_remove_root(self):
        elem = self.heap.remove_root()
        assert elem==2
        assert isHeap(self.heap.elements)

    # Removing an element from an empty array throws
    def test_remove_all(self):
        n = len(self.heap.elements)
        removed_elems = []
        for i in range(n):

            elem = self.heap.remove_root()
            removed_elems.append(elem)
            assert isHeap(self.heap.elements)
        
        assert isAscending(removed_elems)

        with pytest.raises(Exception):
            self.heap.remove_root()


    # Heapifying a random array gives a heap
    @pytest.mark.parametrize("seed", [42, 7, 99, 31, 10000])
    def test_heapify_method1(self, seed):
        
        self.heap = Heap()
        nums = self.generate_list(seed)
        self.heap.heapify_method1(nums)
        assert len(self.heap.elements) == 100
        assert isHeap(self.heap.elements)


    # Heapifying a random array gives a heap
    @pytest.mark.parametrize("seed", [42, 7, 99, 31, 10000])
    def test_heapify_method2(self, seed):
        
        self.heap = Heap()
        nums = self.generate_list(seed)
        self.heap.heapify_method2(nums)
        assert len(self.heap.elements) == 100
        assert isHeap(self.heap.elements)






