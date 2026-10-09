package raft

// GetCurrentTerm returns the current term safely.
func (n *Node) GetCurrentTerm() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.CurrentTerm
}
