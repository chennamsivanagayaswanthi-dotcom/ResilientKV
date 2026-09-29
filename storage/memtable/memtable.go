package memtable

import "sync"

// MemTable stores key-value pairs in memory.
type MemTable struct {
	mu   sync.RWMutex
	data map[string]string
}

// New creates a new MemTable.
func New() *MemTable {
	return &MemTable{
		data: make(map[string]string),
	}
}

// Put inserts or updates a key-value pair.
func (m *MemTable) Put(key string, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = value
}

// Get retrieves a value using a key.
func (m *MemTable) Get(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, exists := m.data[key]
	return value, exists
}

// Delete removes a key from the MemTable.
func (m *MemTable) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, key)
}

// Size returns the number of entries.
func (m *MemTable) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.data)
}
