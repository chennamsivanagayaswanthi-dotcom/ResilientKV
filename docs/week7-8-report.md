# Week 7–8 — Distributed Primitives

## 1. Objective

The objective of Weeks 7–8 was to implement distributed communication
primitives for ResilientKV using gRPC.

## 2. Protocol Buffers

A Protocol Buffers API was created for the key-value service.

Operations:

- PUT
- GET
- DELETE

Request IDs were added to support duplicate-request detection.

## 3. gRPC Server

A gRPC server was implemented on port 50051.

The server connects incoming requests to the ResilientKV storage engine.

Architecture:

Client
  |
  v
gRPC Server
  |
  v
Storage Engine
  |
  +-- MemTable
  +-- WAL
  +-- SSTable

## 4. gRPC Client

A Go client was implemented to communicate with the server.

The client successfully performed:

- PUT
- GET
- DELETE

## 5. At-Least-Once Semantics

Request IDs were introduced to identify duplicate requests.

When the same request ID is received again, the server recognizes
the request as already processed.

## 6. Client Routing

A basic routing component was implemented.

The router selects a server node based on the key.

Example:

key -> hash -> node selection

## 7. Testing

The following tests were performed:

go test ./...

go test -race ./...

The project passed the existing unit tests and race detector tests.

## 8. Result

Weeks 7–8 established the distributed communication layer required
for the next milestone.

The system can now communicate through gRPC and access the storage
engine through the service interface.

## 9. Next Milestone

Weeks 9–10 will implement Raft consensus:

- Leader election
- Terms
- Voting
- Log replication
- Commit index
- Safety tests
