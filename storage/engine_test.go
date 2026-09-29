package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnginePutGet(t *testing.T) {
	dataDir := t.TempDir()

	engine, err := New(dataDir)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	defer engine.Close()

	err = engine.Put("name", "Yaswanthi")
	if err != nil {
		t.Fatalf("failed to put name: %v", err)
	}

	err = engine.Put("city", "Guntur")
	if err != nil {
		t.Fatalf("failed to put city: %v", err)
	}

	value, exists := engine.Get("name")

	if !exists {
		t.Fatal("expected name to exist")
	}

	if value != "Yaswanthi" {
		t.Fatalf("expected Yaswanthi, got %s", value)
	}

	value, exists = engine.Get("city")

	if !exists {
		t.Fatal("expected city to exist")
	}

	if value != "Guntur" {
		t.Fatalf("expected Guntur, got %s", value)
	}
}

func TestEngineFlush(t *testing.T) {
	dataDir := t.TempDir()

	engine, err := New(dataDir)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	defer engine.Close()

	engine.Put("name", "Yaswanthi")
	engine.Put("city", "Guntur")
	engine.Put("course", "MTech")

	err = engine.Flush()
	if err != nil {
		t.Fatalf("failed to flush: %v", err)
	}

	files, err := os.ReadDir(dataDir)
	if err != nil {
		t.Fatalf("failed to read data directory: %v", err)
	}

	foundSSTable := false

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sst" {
			foundSSTable = true
			break
		}
	}

	if !foundSSTable {
		t.Fatal("expected SSTable file to be created")
	}
}
