
import csv
import socket
import time
from datetime import datetime
from pathlib import Path


def server_status(host, port):
    try:
        with socket.create_connection((host, port), timeout=1):
            return "UP"
    except OSError:
        return "DOWN"


def read_cpu_memory():
    try:
        with open("/proc/stat") as file:
            cpu = list(map(int, file.readline().split()[1:]))
        total = sum(cpu)
        idle = cpu[3] + (cpu[4] if len(cpu) > 4 else 0)

        with open("/proc/meminfo") as file:
            info = {
                line.split(":")[0]: int(line.split()[1])
                for line in file
            }

        total_memory = info["MemTotal"]
        available = info["MemAvailable"]
        memory_percent = (total_memory - available) * 100 / total_memory

        return total, idle, memory_percent
    except (OSError, KeyError, ValueError):
        return 0, 0, 0.0


servers = [
    ("node1", "127.0.0.1", 50051),
    ("node2", "127.0.0.1", 50052),
    ("node3", "127.0.0.1", 50053),
]

print("Monitoring ResilientKV for 30 seconds...")
Path("monitoring_results.csv").parent.mkdir(exist_ok=True)

with open("monitoring_results.csv", "w", newline="") as file:
    writer = csv.writer(file)
    writer.writerow([
        "Timestamp", "Server", "Port", "Status",
        "CPU_Percent", "Memory_Percent"
    ])

    previous_total, previous_idle, _ = read_cpu_memory()

    for _ in range(6):
        time.sleep(5)

        total, idle, memory = read_cpu_memory()
        difference = total - previous_total
        idle_difference = idle - previous_idle
        cpu = (
            100 * (difference - idle_difference) / difference
            if difference > 0 else 0
        )
        previous_total, previous_idle = total, idle

        for name, host, port in servers:
            status = server_status(host, port)
            writer.writerow([
                datetime.now().isoformat(timespec="seconds"),
                name, port, status, round(cpu, 2), round(memory, 2)
            ])
            print(f"{name}:{port} = {status}, CPU={cpu:.1f}%, Memory={memory:.1f}%")

        file.flush()

print("Results saved to monitoring_results.csv")
