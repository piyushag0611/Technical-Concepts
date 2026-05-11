import re
import numpy as np
import matplotlib.pyplot as plt

with open("results.txt", encoding="utf-16") as f:
    content = f.read()

pattern = r"Benchmark(HeapifyMethod\d)/n=(\d+)-\d+\s+\d+\s+([\d.]+)\s+ns/op\s+([\d.]+)\s+B/op"

data = {"HeapifyMethod1": {}, "HeapifyMethod2": {}}
for method, n, time_ns, mem_b in re.findall(pattern, content):
    data[method][int(n)] = {
        "time_ms": float(time_ns) / 1_000_000,
        "mem_mb":  float(mem_b)  / 1_000_000,
    }

target = [10_000, 100_000, 1_000_000, 10_000_000]
sizes  = [n for n in target if n in data["HeapifyMethod1"] and n in data["HeapifyMethod2"]]
labels = [f"n=10$^{{{len(str(n))-1}}}$" for n in sizes]

m1_time = [data["HeapifyMethod1"][n]["time_ms"] for n in sizes]
m2_time = [data["HeapifyMethod2"][n]["time_ms"] for n in sizes]
m1_mem  = [data["HeapifyMethod1"][n]["mem_mb"]  for n in sizes]
m2_mem  = [data["HeapifyMethod2"][n]["mem_mb"]  for n in sizes]

time_ratios = [a / b for a, b in zip(m1_time, m2_time)]
mem_ratios  = [a / b for a, b in zip(m1_mem,  m2_mem)]

x = np.arange(len(sizes))
w = 0.35

def make_chart(m1_vals, m2_vals, ratios, unit, title, filename):
    fig, (ax1, ax2) = plt.subplots(
        2, 1, figsize=(10, 9), gridspec_kw={"height_ratios": [3, 1]}
    )

    b1 = ax1.bar(x - w / 2, m1_vals, w, label="Method 1 — naive O(n log n)", color="tomato",    alpha=0.9)
    b2 = ax1.bar(x + w / 2, m2_vals, w, label="Method 2 — Floyd's O(n)",     color="steelblue", alpha=0.9)
    ax1.set_xticks(x)
    ax1.set_xticklabels(labels, fontsize=12)
    ax1.set_ylabel(unit, fontsize=12)
    ax1.set_title(title, fontsize=14, fontweight="bold")
    ax1.legend(fontsize=11)
    ax1.bar_label(b1, fmt="%.2f", padding=4, fontsize=9, color="tomato",    fontweight="bold")
    ax1.bar_label(b2, fmt="%.2f", padding=4, fontsize=9, color="steelblue", fontweight="bold")

    ratio_bars = ax2.bar(x, ratios, color="gold", edgecolor="black", alpha=0.9)
    ax2.axhline(y=1, color="black", linestyle="--", linewidth=0.8)
    ax2.set_xticks(x)
    ax2.set_xticklabels(labels, fontsize=12)
    ax2.set_ylabel("Ratio (M1 / M2)", fontsize=10)
    ax2.set_ylim(0, ratios[-1]+2)
    ax2.set_title("How many times faster / leaner is Method 2?", fontsize=11)
    ax2.bar_label(ratio_bars, fmt="%.1fx", padding=4, fontsize=11, fontweight="bold")

    plt.tight_layout(pad=2.0)
    plt.savefig(filename, dpi=150)
    plt.show()
    print(f"Saved {filename}")

make_chart(m1_time, m2_time, time_ratios, "Time (ms)",    "Heapify Runtime Comparison (Go)",        "time_comparison.png")
make_chart(m1_mem,  m2_mem,  mem_ratios,  "Memory (MB)",  "Heapify Memory Comparison (Go)",         "memory_comparison.png")
