package game

import (
	"testing"
	"time"
)

func TestIsTopScoreEmptyList(t *testing.T) {
	if !IsTopScore(1, nil) {
		t.Error("any score should be top score when list is empty")
	}
}

func TestIsTopScoreUnderMax(t *testing.T) {
	entries := make([]ScoreEntry, 5)
	if !IsTopScore(0, entries) {
		t.Error("score should qualify when fewer than 10 entries")
	}
}

func TestIsTopScoreFullListBeats(t *testing.T) {
	entries := make([]ScoreEntry, 10)
	for i := range entries {
		entries[i] = ScoreEntry{Seconds: 10 - i}
	}
	// Sorted descending: 10,9,8,...,1. Last = 1. Score of 2 should qualify.
	if !IsTopScore(2, entries) {
		t.Error("score of 2 should beat last entry of 1")
	}
}

func TestIsTopScoreFullListLoses(t *testing.T) {
	entries := make([]ScoreEntry, 10)
	for i := range entries {
		entries[i] = ScoreEntry{Seconds: 100 - i}
	}
	// Last = 91. Score of 90 should not qualify.
	if IsTopScore(90, entries) {
		t.Error("score of 90 should not beat last entry of 91")
	}
}

func TestInsertScoreAddsEntry(t *testing.T) {
	entry := ScoreEntry{Name: "AAA", Seconds: 42, Date: time.Now()}
	result := InsertScore(entry, nil)
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].Name != "AAA" {
		t.Errorf("expected name AAA, got %s", result[0].Name)
	}
}

func TestInsertScoreSortDescending(t *testing.T) {
	entries := []ScoreEntry{
		{Name: "B", Seconds: 50},
		{Name: "A", Seconds: 100},
	}
	result := InsertScore(ScoreEntry{Name: "C", Seconds: 75}, entries)
	if result[0].Seconds != 100 || result[1].Seconds != 75 || result[2].Seconds != 50 {
		t.Errorf("sort order wrong: %v", result)
	}
}

func TestInsertScoreTrimToMax(t *testing.T) {
	entries := make([]ScoreEntry, 10)
	for i := range entries {
		entries[i] = ScoreEntry{Seconds: 10 - i}
	}
	// Adding a score of 100 should push out the lowest.
	result := InsertScore(ScoreEntry{Name: "TOP", Seconds: 100}, entries)
	if len(result) != 10 {
		t.Fatalf("expected 10 entries after trim, got %d", len(result))
	}
	if result[0].Name != "TOP" {
		t.Error("new top score should be first")
	}
}

func TestFakeScoreStoreRoundTrip(t *testing.T) {
	store := &fakeScoreStore{}
	entries := []ScoreEntry{
		{Name: "ACE", Seconds: 99},
	}
	if err := store.Save(entries); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(loaded) != 1 || loaded[0].Name != "ACE" {
		t.Errorf("round trip failed: %v", loaded)
	}
}

func TestFakeScoreStoreLoadEmpty(t *testing.T) {
	store := &fakeScoreStore{}
	entries, err := store.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty, got %v", entries)
	}
}
