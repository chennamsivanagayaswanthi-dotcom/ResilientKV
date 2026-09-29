package sstable

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Entry represents one key-value pair.
type Entry struct {
	Key   string
	Value string
}

// Write creates an SSTable from key-value entries.
func Write(path string, entries []Entry) error {
	// Sort entries by key.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, entry := range entries {
		_, err := fmt.Fprintf(writer, "%s=%s\n", entry.Key, entry.Value)
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

// ReadAll reads all entries from an SSTable.
func ReadAll(path string) ([]Entry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var entries []Entry

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.SplitN(line, "=", 2)

		if len(parts) != 2 {
			continue
		}

		entries = append(entries, Entry{
			Key:   parts[0],
			Value: parts[1],
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
