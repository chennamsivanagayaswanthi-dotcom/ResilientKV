
import csv
import matplotlib.pyplot as plt

times = []
cpu = []
memory = []

with open("monitoring_results.csv", "r") as file:
    reader = csv.DictReader(file)
    for i, row in enumerate(reader):
        times.append(i + 1)
        cpu.append(float(row["CPU_Percent"]))
        memory.append(float(row["Memory_Percent"]))

plt.plot(times, cpu, marker="o", label="CPU Usage (%)")
plt.plot(times, memory, marker="s", label="Memory Usage (%)")

plt.xlabel("Monitoring Sample")
plt.ylabel("Usage (%)")
plt.title("ResilientKV Resource Monitoring")
plt.legend()
plt.grid(True)
plt.tight_layout()
plt.savefig("monitoring_graph.png", dpi=150)
print("Graph saved as monitoring_graph.png")
plt.show()
