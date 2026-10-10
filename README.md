# ResilientKV: Design and Implementation of a Fault-Tolerant Distributed Key-Value Store Using Go, gRPC, and Raft Consensus

**Project Type:** M.Tech Computer Science and Engineering Project  
**Domain:** Operating Systems and Distributed Systems  
**Implementation Language:** Go

---

## 1. Problem Statement

Distributed applications require storage systems that support persistent data storage, efficient access, concurrent operations, and consistency across multiple nodes. Failures in individual components can affect data availability and reliability.

ResilientKV aims to address these challenges by developing a distributed key-value store using Go, persistent storage structures, gRPC communication, and the Raft consensus algorithm.

## 2. Objectives

- Design and implement a key-value storage engine using Go.
- Provide persistent storage using a Write-Ahead Log (WAL), Memtable, and SSTables.
- Support GET and PUT operations for storing and retrieving data.
- Enable client-server communication through gRPC and Protocol Buffers.
- Implement Raft-based leader election and log replication.
- Validate storage correctness, persistence, concurrency, and replication through automated tests.
- Identify limitations and explore improvements in fault tolerance and performance.

## 3. Technologies Used

| Technology | Purpose |
|---|---|
| Go | Core system implementation |
| gRPC | Communication between clients and servers |
| Protocol Buffers | Service and message definitions |
| Write-Ahead Log (WAL) | Recording storage operations for recovery |
| Memtable | Managing recent updates in memory |
| SSTables | Storing sorted data persistently |
| Raft Consensus | Leader election and log replication |
| Go Testing | Functional and integration testing |
| Go Race Detector | Detecting data races |

## 4. System Architecture

ResilientKV follows a modular architecture that separates client requests, communication, consensus, and storage.

### Components

1. **Client Layer:** Sends GET and PUT requests to the server.
2. **gRPC Communication Layer:** Transfers requests and responses between clients and servers.
3. **Raft Consensus Layer:** Coordinates leader election, voting, and log replication.
4. **Storage Engine:** Manages key-value data using the WAL, Memtable, SSTables, and compaction.
5. **Testing Layer:** Verifies storage operations, recovery behavior, replication, and concurrency.

### Request Flow

```text
        Client
          |
          v
     gRPC Server
          |
          v
     Raft Processing
          |
          v
     Storage Engine
          |
          v
  WAL + Memtable + SSTables
          |
          v
   Response to Client
```

*The diagram illustrates the logical request flow. The exact execution path depends on the implemented leader-handling and replication logic.*

## 5. Development Roadmap: Week 1–12

### Week 1: Requirement Analysis
Defined the project problem, objectives, scope, and expected outcomes. Identified the need for persistent storage, efficient data access, and reliable communication between distributed nodes.

### Week 2: System Architecture and Setup
Designed the modular architecture and organized the Go project into storage, communication, consensus, and testing components. Prepared the development environment and project structure.

### Week 3: Write-Ahead Log (WAL)
Implemented the WAL to record storage operations before applying them. This provides a foundation for recovering storage state after an unexpected shutdown.

### Week 4: Memtable Implementation
Developed the Memtable to maintain recent key-value updates in memory. Implemented basic data insertion and retrieval operations for efficient access.

### Week 5: SSTable Implementation
Developed SSTable components to store data in sorted, persistent files. Added supporting storage functionality and tests.

### Week 6: Storage Integration and Recovery
Integrated the storage components and worked on compaction and recovery behavior. Tested persistent data handling, including tombstone recovery, to improve storage correctness.

### Week 7: gRPC Communication
Defined service interfaces using Protocol Buffers and implemented gRPC-based communication. Established the foundation for clients and servers to exchange requests and responses.

### Week 8: Client-Server Operations
Integrated client requests with the storage service and tested GET and PUT operations. Verified request handling, retries, and retrieval of stored key-value pairs.

### Week 9: Raft Leader Election
Implemented core Raft election functionality, including voting, terms, and leader selection. Tested the election process to verify successful leader election.

### Week 10: Raft Log Replication
Worked on replicating log entries between nodes and integrating Raft with the gRPC service. Tested follower replication to validate distributed request handling.

### Week 11: Testing and Validation
Executed functional, integration, storage, and concurrency tests. Used Go's testing framework and race detector to identify correctness issues and check concurrent data access.

### Week 12: Documentation and Final Evaluation
Organized project documentation, recorded test evidence, and identified remaining work. Planned performance benchmarking, fault-injection testing, deployment, and monitoring as further enhancements.

> **Progress note:** The roadmap summarizes the intended development sequence. Activities should be marked completed only when supported by implementation, test results, or other project evidence.

## 6. Key Features

- Persistent key-value storage using WAL, Memtable, and SSTables.
- Modular storage and communication components.
- GET and PUT operations.
- gRPC-based client-server communication.
- Raft leader election and log replication functionality.
- Storage recovery and tombstone handling tests.
- Automated functional and integration testing.
- Concurrency checking using Go's race detector.

## 7. Implementation and Execution

### Run the Test Suite

```bash
go test -count=1 ./...
```

### Run Race Detection

```bash
go test -race ./...
```

### Start the Server

From the project root, run:

```bash
go run ./cmd/server
```

The server uses the project's configured defaults and environment variables for node identity, port, data directory, and Raft peers.

## 8. Testing and Verification

The following results have been reported from project testing:

- The complete Go test suite passed in the tested environment.
- Earlier race-detector testing completed successfully.
- Raft election testing verified successful leader election.
- A gRPC integration test verified PUT replication to a follower.
- Client testing verified PUT, retry, and GET operations.
- Storage tests covered persistence and tombstone recovery.

These results provide evidence of core functionality. They do not, by themselves, establish production-level fault tolerance or performance.

## 9. Current Limitations

- Durable persistence of Raft metadata requires further work.
- Automatic election timeout and re-election behavior need additional validation.
- Follower catch-up after reconnection requires improvement.
- Large-scale throughput and latency benchmarking remain to be completed.
- Fault-injection and extended distributed consistency testing remain future work.
- Kubernetes deployment, Helm configuration, and Prometheus/Grafana monitoring require completion and verification if not already implemented.

## 10. Future Enhancements

- Strengthen Raft recovery and re-election mechanisms.
- Improve replication and follower catch-up after failures.
- Benchmark throughput, latency, and resource utilization under different workloads.
- Integrate Kubernetes-based deployment and monitoring.
- Conduct fault-injection and Jepsen-style consistency testing.
- Improve operational observability and distributed-system reliability.

## 11. Conclusion

ResilientKV demonstrates the design and implementation of a distributed key-value storage system using Go, persistent storage structures, gRPC, and Raft consensus. The project brings together storage-engine design, inter-node communication, distributed coordination, and automated testing.

The implementation provides a foundation for further research and development in fault-tolerant storage systems, with additional work needed to strengthen recovery, large-scale performance, and production readiness.

---

