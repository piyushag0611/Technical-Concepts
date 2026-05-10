import random
import timeit
from heap import Heap

sizes = [100, 1000, 10000, 100000, 1000000, 10000000]
attempts = 5
for n in sizes:
    data = random.sample(range(n * 10), n)

    t1 = timeit.timeit(lambda: Heap().heapify_method1(data[:]), number=attempts) / attempts
    t2 = timeit.timeit(lambda: Heap().heapify_method2(data[:]), number=attempts) / attempts

    print(f"n={n:<8} method1={t1:.4f}s  method2={t2:.4f}s  ratio={t1/t2:.2f}x")
