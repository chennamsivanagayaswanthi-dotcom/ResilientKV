package sstable

import (
	"os"
	"testing"
)

func TestSSTableWriteAndRead(t *testing.T) {
	file := "test.sst"

	defer os.Remove(file)

	entries := []Entry{
		{Key: "name", Value: "Yaswanthi"},
		{Key: "city", Value: "Guntur"},
		{Key: "course", Value: "MTech"},
	}

	err := Write(file, entries)

	if err != nil {
		t.Fatalf("failed to write SSTable: %v", err)
	}

	result, err := ReadAll(file)

	if err != nil {
		t.Fatalf("failed to read SSTable: %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}

	if result[0].Key != "city" {
		t.Fatalf("expected first key to be city, got %s", result[0].Key)
	}

	if result[1].Key != "course" {
		t.Fatalf("expected second key to be course, got %s", result[1].Key)
	}

	if result[2].Key != "name" {
		t.Fatalf("expected third key to be name, got %s", result[2].Key)
	}
}
