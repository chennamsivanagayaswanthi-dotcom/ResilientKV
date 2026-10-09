package raft

import (
	"errors"
)

// AppendEntriesRequest contains the information sent by a Raft leader.
type AppendEntriesRequest struct {
	Term         int
	LeaderID     string
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []LogEntry
	LeaderCommit int
}

// AppendEntriesResponse reports whether the follower accepted the request.
type AppendEntriesResponse struct {
	Term    int
	Success bool
}

// HandleAppendEntries processes a leader heartbeat or log replication request.
func (n *Node) HandleAppendEntries(req AppendEntriesRequest) (AppendEntriesResponse, error) {
	if n == nil {
		return AppendEntriesResponse{}, errors.New("Raft node is nil")
	}
	if req.LeaderID == "" {
		return AppendEntriesResponse{}, errors.New("leader ID cannot be empty")
	}
	if req.Term < 0 || req.PrevLogIndex < 0 ||
		req.PrevLogTerm < 0 || req.LeaderCommit < 0 {
		return AppendEntriesResponse{}, errors.New("Raft request contains a negative value")
	}

	for _, entry := range req.Entries {
		if entry.Term < 0 || entry.Term > req.Term {
			return AppendEntriesResponse{}, errors.New("log entry term is invalid")
		}
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	// Reject requests from an outdated leader.
	if req.Term < n.CurrentTerm {
		return AppendEntriesResponse{
			Term: n.CurrentTerm, Success: false,
		}, nil
	}

	// A newer term, or a valid leader in the current term, makes this
	// node a follower.
	if req.Term > n.CurrentTerm {
		n.CurrentTerm = req.Term
		n.VotedFor = ""
	}
	n.State = Follower

	// Raft log indexes in this request are zero-based positions.
	if req.PrevLogIndex > len(n.Log) {
		return AppendEntriesResponse{
			Term: n.CurrentTerm, Success: false,
		}, nil
	}

	// Check that the preceding log entry matches the leader's log.
	if req.PrevLogIndex > 0 {
		localPrev := n.Log[req.PrevLogIndex-1]
		if localPrev.Term != req.PrevLogTerm {
			return AppendEntriesResponse{
				Term: n.CurrentTerm, Success: false,
			}, nil
		}
	}

	// Keep matching entries. Replace conflicting entries and everything
	// after them with the leader's entries.
	for i, incoming := range req.Entries {
		index := req.PrevLogIndex + i

		if index < len(n.Log) {
			if n.Log[index].Term != incoming.Term {
				if index < n.CommitIndex {
					return AppendEntriesResponse{
						Term: n.CurrentTerm, Success: false,
					}, nil
				}
				n.Log = n.Log[:index]
				n.Log = append(n.Log, req.Entries[i:]...)
				break
			}
			continue
		}

		n.Log = append(n.Log, req.Entries[i:]...)
		break
	}

	// Never move the commit index backwards.
	if req.LeaderCommit > n.CommitIndex {
		if req.LeaderCommit < len(n.Log) {
			n.CommitIndex = req.LeaderCommit
		} else {
			n.CommitIndex = len(n.Log)
		}
	}

	return AppendEntriesResponse{
		Term: n.CurrentTerm, Success: true,
	}, nil
}
