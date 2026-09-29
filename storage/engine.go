package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"resilientkv/storage/memtable"
	"resilientkv/storage/sstable"
	"resilientkv/storage/wal"
)

// Engine combines MemTable, WAL, and SSTable storage.
type Engine struct {
	memTable *memtable.MemTable
	wal      *wal.WAL
	dataDir  string
}

// New creates a new storage engine.
func New(dataDir string) (*Engine, error) {
	err := os.MkdirAll(dataDir, 0755)
	if err != nil {
		return nil, err
	}

	walPath := filepath.Join(dataDir, "wal.log")

	w, err := wal.New(walPath)
	if err != nil {
		return nil, err
	}

	return &Engine{
		memTable: memtable.New(),
		wal:      w,
		dataDir:  dataDir,
	}, nil
}

// Put stores a key-value pair.
func (e *Engine) Put(key string, value string) error {
	operation := fmt.Sprintf("PUT %s %s", key, value)

	// Write to WAL first.
	err := e.wal.Append(operation)
	if err != nil {
		return err
	}

	// Then update MemTable.
	e.memTable.Put(key, value)

	return nil
}

// Get retrieves a value from the MemTable.
func (e *Engine) Get(key string) (string, bool) {
	return e.memTable.Get(key)
}

// Flush writes the MemTable contents to an SSTable.
func (e *Engine) Flush() error {
	entries := e.memTable.AllEntries()

	if len(entries) == 0 {
		return nil
	}

	sstablePath := filepath.Join(
		e.dataDir,
		fmt.Sprintf("sstable-%d.sst", len(entries)),
	)

	return sstable.Write(sstablePath, convertEntries(entries))
}

// convertEntries converts MemTable entries to SSTable entries.
func convertEntries(entries []memtable.Entry) []sstable.Entry {
	result := make([]sstable.Entry, 0, len(entries))

	for _, entry := range entries {
		result = append(result, sstable.Entry{
			Key:   entry.Key,
			Value: entry.Value,
		})
	}

	return result
}

// Close closes the storage engine.
func (e *Engine) Close() error {
	return e.wal.Close()
}
