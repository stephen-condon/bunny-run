package game

import "time"

// Clock abstracts time so tests can control it.
type Clock interface {
	Now() time.Time
}

// InputSource abstracts keyboard input.
type InputSource interface {
	IsUpPressed() bool
	IsDownPressed() bool
	IsLeftPressed() bool
	IsRightPressed() bool
	IsActionPressed() bool // confirm / enter
	JustPressedChar() rune // for name entry; 0 if none
	IsBackspaceJustPressed() bool
}

// ScoreStore persists and retrieves high scores.
type ScoreStore interface {
	Load() ([]ScoreEntry, error)
	Save(entries []ScoreEntry) error
}

// WorldReader gives entities read-only access to the tile grid.
type WorldReader interface {
	TileAt(x, y int) TileType
	Width() int
	Height() int
}
