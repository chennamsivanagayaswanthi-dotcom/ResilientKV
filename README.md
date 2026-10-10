# ResilientKV
### A Distributed Key–Value Store with Raft Consensus

ResilientKV is an educational distributed key–value store project built to explore how storage engines, service communication, consensus, testing, and operations work together in a distributed system.

The project uses **Go**, an **LSM-tree-style storage engine**, **gRPC**, and **Raft-related replication components**. The goal is to build a system that stores key–value data, communicates through a client/server API, and explores how data can be coordinated across multiple nodes.

> **Project status:** Core storage, gRPC, and Raft-related functionality has been implemented and tested in the development workflow. Kubernetes deployment, complete observability, extensive fault-injection testing, and comparative performance reports should be marked complete only after their results have been produced and documented.

---

## 1. Project at a Glance

| Item | Details |
|---|---|
| Project | ResilientKV |
| Category | Distributed Systems / Operating Systems |
| Main language | Go |
| API communication | gRPC and Protocol Buffers |
| Storage design | Write-Ahead Log (WAL), Memtable, SSTable, compaction |
| Distributed coordination | Raft-related election, replication, and commit handling |
| Testing | Go unit/integration tests; race detector where available |
| Repository | [ResilientKV on GitHub](https://github.com/chennamsivanagayaswanthi-dotcom/ResilientKV) |

## 2. Problem Statement

Applications need to store data and retrieve it reliably. In a distributed system, multiple nodes may need to coordinate updates, and failures can interrupt communication or access to a node.

ResilientKV is designed as a learning project to investigate these challenges:

- How can key–value data be stored efficiently?
- How can recent writes be recovered after a process failure?
- How can clients communicate with a storage service?
- How can nodes coordinate replicated log entries?
- How can correctness and performance be tested?

## 3. Project Objectives

1. Implement key–value operations such as `PUT` and `GET`.
2. Build storage components using a WAL, memtable, SSTables, and compaction.
3. Provide client/server communication through gRPC.
4. Implement and test Raft-related leader election, log replication, and commit handling.
5. Test correctness and concurrency, and document reproducible performance measurements.
6. Prepare deployment, monitoring, and fault-injection artifacts as part of the full project roadmap.

## 4. Architecture

```text
                    +----------------------+
                    |      Go Client       |
                    |      PUT / GET       |
                    +----------+-----------+
                               |
                               | gRPC
                               v
                    +----------------------+
                    |   Service / Routing  |
                    +----------+-----------+
                               |
                               v
                    +----------------------+
                    | Raft-related Layer   |
                    | Election / Replication|
                    | Commit Handling      |
                    +----------+-----------+
                               |
                               v
                    +----------------------+
                    |    Storage Engine    |
                    | WAL -> Memtable      |
                    |      -> SSTables     |
                    |      -> Compaction   |
                    +----------------------+
```

**Important concepts**

- **Client:** Sends requests to store or retrieve a key/value pair.
- **gRPC:** Carries typed requests and responses between client and service.
- **Raft:** Provides rules for coordinating replicated log entries among nodes.
- **WAL:** Records changes so recovery can replay them.
- **Memtable:** Holds recent updates in memory.
- **SSTable:** Stores sorted, immutable data on disk.
- **Compaction:** Merges storage files and handles obsolete entries.

The diagram is a conceptual overview. The exact request path and guarantees depend on the implemented configuration and code.

## 5. Technology Stack

- **Go:** Main implementation language.
- **gRPC:** Client/server remote procedure calls.
- **Protocol Buffers:** Service and message definitions.
- **WAL:** Write-ahead logging for storage recovery.
- **Memtable and SSTable:** In-memory and on-disk storage components.
- **Go testing tools:** Unit and integration tests.
- **Linux performance tools:** `perf` and `bpftrace` for profiling, where supported.
- **Planned operations tools:** Kubernetes, Helm, Prometheus, Grafana, and fault-injection tooling.

## 6. Development Roadmap — Weeks 1 to 12

| Weeks | Focus | Main activities / evidence |
|---|---|---|
| 1–2 | Foundation | Lock topic, create repository and Go scaffold, establish Linux profiling baseline |
| 3–4 | Storage engine | Implement WAL, memtable, SSTables, recovery and compaction; plan a fair RocksDB comparison |
| 5–6 | Concurrency | Build concurrent structures, run race detection, benchmark several worker counts |
| 7–8 | Service layer | Build gRPC service/client, routing, request IDs and retry behavior |
| 9–10 | Consensus | Implement Raft-related election, replication and commit handling; add safety tests |
| 11–12 | Operations | Prepare Kubernetes/Helm deployment, Prometheus/Grafana dashboards and Jepsen-style fault testing |

This table describes the full project roadmap. It is not a claim that every milestone or deliverable is complete.

## 7. Storage Engine Overview

### Write-Ahead Log (WAL)
A write is recorded in the log so that it can be replayed during recovery. WAL recovery behavior should be validated with restart and failure tests.

### Memtable
Recent key/value updates are held in memory for fast access before being flushed to disk.

### SSTable
The memtable can be flushed into sorted, immutable files called SSTables. Indexing helps locate entries.

### Compaction
Compaction merges storage files and can remove obsolete versions and tombstones when it is safe to do so.

### Simplified write flow

```text
PUT request
    |
    v
Record write in WAL
    |
    v
Update Memtable
    |
    v
Flush sorted data to SSTable
    |
    v
Compact files when needed
```

The precise ordering and durability guarantees should be confirmed against the current implementation and tests.

## 8. gRPC and Client Example

A simple key/value example:

```text
Key   = product
Value = laptop
```

A client can send a `PUT` request to store the pair and a `GET` request to retrieve it. Request IDs can help identify retries, but retry safety depends on the server's deduplication behavior and the operation semantics.

Example demonstration sequence:

1. Start the server or the required set of nodes.
2. Run the client.
3. Send a `PUT` request for `product = laptop`.
4. Retry the request with the same request ID, if supported by the demo.
5. Send a `GET` request and inspect the returned value.

Use the commands and configuration from the current source tree when running the demo; exact options can change as the project evolves.

## 9. Raft Concepts

Raft is a consensus algorithm used to coordinate a replicated log.

- **Leader election:** Nodes vote for a candidate; a majority is needed to elect a leader for a term.
- **Log replication:** The leader sends log entries to followers.
- **Commit handling:** Entries are committed according to Raft rules and then applied to the state machine.
- **Safety testing:** Tests should check voting rules, terms, replication, commit decisions, and behavior when nodes or network connections fail.

A successful leader-election test alone does not establish that all fault-tolerance or production requirements are met.

## 10. Build and Test

Run these commands from the repository root, where `go.mod` is located.

### Check Go installation

```bash
go version
```

### Build all packages

```bash
go build ./...
```

### Run all tests

```bash
go test -count=1 ./...
```

### Run the race detector

```bash
go test -race ./...
```

The race detector may increase runtime and memory use. If a command fails, retain the complete error output and the Go version to help reproduce the issue.

## 11. Verified Test Evidence

During the development workflow, the following results were observed:

- `go test -count=1 ./...` passed on a teammate's cloned repository for the listed gRPC, Raft, and storage packages.
- A client demonstration reported a successful `PUT`, a successful retry, and a `GET` returning `product = laptop`.
- A separate Raft election client reported `elected=true` and `state=Leader`.

These are specific development results, not a claim of production readiness. Add dates, environment details, terminal screenshots, and test output to the final report when available.

## 12. Performance Evaluation Plan

For a fair performance report, record:

- CPU, RAM, operating system/kernel, Go version, storage device, node count, and configuration.
- Workload type: read-heavy, write-heavy, or mixed.
- Number of operations, key/value sizes, concurrency level, warm-up, and measurement duration.
- **Throughput:** operations per second.
- **Latency:** p50, p95, and p99 where measured.
- CPU/memory use, errors, and recovery time where relevant.

Repeat runs and report the method and variation. Compare against RocksDB, etcd, or TiKV only when the workloads, environment, and semantics are reasonably comparable. Do not insert benchmark numbers until they have actually been measured.

## 13. Deployment and Observability Roadmap

The full project brief calls for the following deliverables:

- **Kubernetes:** Run the service in a multi-node or multi-replica deployment.
- **Helm:** Package deployment resources and configuration.
- **Prometheus:** Collect service and system metrics.
- **Grafana:** Provide importable dashboards for health, request rate, latency, errors, and resource usage.
- **Jepsen-style harness:** Inject node or network failures, record operation histories, and check consistency properties.
- **Live demo:** Demonstrate the deployment and explain the tested fault scenarios.

Treat these as pending until the relevant artifacts and results are present in the repository.

## 14. Known Limitations and Next Steps

This is an educational implementation and should not be described as production-ready without further engineering and validation. Areas requiring careful review include:

- Durable Raft term/vote metadata and complete crash recovery.
- Automatic election timeouts and re-election behavior.
- Follower catch-up after disconnection and restart.
- Broader safety testing under delays, partitions, and node failures.
- Reproducible comparisons with established storage/distributed systems.
- Completed Kubernetes/Helm deployment, dashboards, and fault-injection evidence.

Recommended next steps are to close these gaps, document the behavior that is actually guaranteed, and retain reproducible test and benchmark evidence.

## 15. Final Deliverables Checklist

- [x] Source code hosted in GitHub.
- [x] Storage components including WAL, memtable, SSTable and compaction code.
- [x] gRPC service/client components.
- [x] Raft-related election, replication and commit-handling code.
- [x] Go test suite passes in the recorded development workflow.
- [ ] Separate LSM-tree library packaging, if required by the course.
- [ ] Full Raft safety and failure-recovery test report.
- [ ] RocksDB / etcd / TiKV comparative benchmark report.
- [ ] Kubernetes deployment and Helm chart.
- [ ] Prometheus metrics and importable Grafana dashboard JSON.
- [ ] Jepsen-style fault-injection harness and documented results.
- [ ] Architecture document (10–15 pages, ADR-style).
- [ ] Live multi-node Kubernetes demo.

> Update the checklist as each item is completed and verified. The checked items reflect the implementation and test evidence currently recorded in this README.

## 16. Conclusion

ResilientKV is a learning project that brings together storage-engine concepts, service communication, distributed coordination, testing, and operations. It provides a foundation for understanding how distributed key–value systems are designed and evaluated.

The project should be presented with reproducible evidence: show the source code, run the tests, demonstrate the client, explain the design, and clearly distinguish verified functionality from remaining work.

---

**Repository:** https://github.com/chennamsivanagayaswanthi-dotcom/ResilientKV
