package game

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const maxScores = 10

// ScoreEntry holds a single leaderboard entry.
type ScoreEntry struct {
	Name    string    `json:"name"`
	Seconds int       `json:"seconds"`
	Date    time.Time `json:"date"`
}

// IsTopScore returns true if s would make the top-10 list given existing entries.
func IsTopScore(s int, entries []ScoreEntry) bool {
	if len(entries) < maxScores {
		return true
	}
	return s > entries[len(entries)-1].Seconds
}

// InsertScore inserts a new entry, sorts descending, and trims to top 10.
func InsertScore(entry ScoreEntry, entries []ScoreEntry) []ScoreEntry {
	entries = append(entries, entry)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Seconds > entries[j].Seconds
	})
	if len(entries) > maxScores {
		entries = entries[:maxScores]
	}
	return entries
}

// FileScoreStore persists scores to ~/.bunny-run/scores.json.
type FileScoreStore struct {
	path string
}

func NewFileScoreStore() *FileScoreStore {
	home, _ := os.UserHomeDir()
	return &FileScoreStore{path: filepath.Join(home, ".bunny-run", "scores.json")}
}

func (s *FileScoreStore) Load() ([]ScoreEntry, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []ScoreEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func (s *FileScoreStore) Save(entries []ScoreEntry) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}
