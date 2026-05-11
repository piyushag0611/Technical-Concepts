"""
Heaps are arrays which when visualized as a binary tree has its parents larger / smaller than its children.
A min heap has its parents smaller than their children, viceversa for max heap.

Operations:

a) Push : Add a new element to the heap O(log k)
b) Pop : Removes an element from the heap O(log k)
c) Peak: Gets the root element O(1)
d) Heapify : O(n)


"""

class Heap:

    def __init__(self):

        self.elements = []

    def push_elements(self, num):

        pos = len(self.elements)
        self.elements.append(num)
        while(pos > 0 and self.elements[(pos-1)//2] > self.elements[pos]):
            self.elements[pos], self.elements[(pos-1)//2] = self.elements[(pos-1)//2], self.elements[pos]
            pos = (pos-1)//2
        
    def remove_root(self):

        n = len(self.elements)
        if (n == 0):
            raise Exception("No elements to remove")

        if (n==1):
            return self.elements.pop()

        root = self.elements[0]
        last_elem = self.elements.pop()
        n = len(self.elements)
        self.elements[0] = last_elem
        pos = 0
        while(2*pos+1 < n):
            lc = 2*pos+1 #Left Child
            if (lc+1 < n):
                min_elem, min_index = (self.elements[lc], lc) if (self.elements[lc] < self.elements[lc+1]) else (self.elements[lc+1], lc+1)
                if (min_elem < self.elements[pos]):
                    self.elements[min_index], self.elements[pos] = self.elements[pos], self.elements[min_index]
                    pos = min_index
                else:
                    break
            else:
                if (self.elements[lc] < self.elements[pos]):
                    self.elements[lc], self.elements[pos] = self.elements[pos], self.elements[lc]
                    pos = lc
                else:
                    break
        return root 


    def heapify_method1(self, array_elements):
        
        for elem in array_elements:
            self.push_elements(elem)
    
    def heapify_method2(self, array_elements):

        self.elements = array_elements
        n = len(self.elements)
        pos = n//2 - 1  # last non leaf node

        while (pos >= 0):

            current_pos = pos
            while(2*current_pos+1 < n):
                lc = 2*current_pos+1 #Left Child
                if (lc+1 < n):
                    min_elem, min_index = (self.elements[lc], lc) if (self.elements[lc] < self.elements[lc+1]) else (self.elements[lc+1], lc+1)
                    if (min_elem < self.elements[current_pos]):
                        self.elements[min_index], self.elements[current_pos] = self.elements[current_pos], self.elements[min_index]
                        current_pos = min_index
                    else:
                        break
                else:
                    if (self.elements[lc] < self.elements[current_pos]):
                        self.elements[lc], self.elements[current_pos] = self.elements[current_pos], self.elements[lc]
                        current_pos = lc
                    else:
                        break
            pos -= 1
        



