package game

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileScoreStoreRoundTrip(t *testing.T) {
	dir := t.TempDir()
	store := &FileScoreStore{path: filepath.Join(dir, "scores.json")}

	entries := []ScoreEntry{
		{Name: "ACE", Seconds: 100, Date: time.Now()},
		{Name: "BOB", Seconds: 50, Date: time.Now()},
	}
	if err := store.Save(entries); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(loaded))
	}
	if loaded[0].Name != "ACE" || loaded[0].Seconds != 100 {
		t.Errorf("first entry wrong: %+v", loaded[0])
	}
}

func TestFileScoreStoreLoadMissing(t *testing.T) {
	dir := t.TempDir()
	store := &FileScoreStore{path: filepath.Join(dir, "missing.json")}
	entries, err := store.Load()
	if err != nil {
		t.Fatalf("Load of missing file should not error, got: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestFileScoreStoreCreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	store := &FileScoreStore{path: filepath.Join(dir, "subdir", "scores.json")}
	if err := store.Save(nil); err != nil {
		t.Fatalf("Save should create missing directories: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "subdir")); err != nil {
		t.Errorf("subdir should have been created: %v", err)
	}
}

func TestFileScoreStoreLoadInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "scores.json")
	if err := os.WriteFile(path, []byte("not json"), 0600); err != nil {
		t.Fatal(err)
	}
	store := &FileScoreStore{path: path}
	_, err := store.Load()
	if err == nil {
		t.Error("Load should error on invalid JSON")
	}
}
