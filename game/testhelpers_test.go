package game

import "time"

// --- FakeClock ---

type fakeClock struct{ t time.Time }

func newFakeClock(t time.Time) *fakeClock  { return &fakeClock{t} }
func (f *fakeClock) Now() time.Time        { return f.t }
func (f *fakeClock) advance(d time.Duration) { f.t = f.t.Add(d) }

// --- FakeInput ---

type fakeInput struct {
	up, down, left, right, action bool
	ch                            rune
	backspace                     bool
}

func (f *fakeInput) IsUpPressed() bool           { return f.up }
func (f *fakeInput) IsDownPressed() bool          { return f.down }
func (f *fakeInput) IsLeftPressed() bool          { return f.left }
func (f *fakeInput) IsRightPressed() bool         { return f.right }
func (f *fakeInput) IsActionPressed() bool        { return f.action }
func (f *fakeInput) JustPressedChar() rune        { return f.ch }
func (f *fakeInput) IsBackspaceJustPressed() bool { return f.backspace }

// --- FakeScoreStore ---

type fakeScoreStore struct {
	entries []ScoreEntry
	saveErr error
	loadErr error
}

func (f *fakeScoreStore) Load() ([]ScoreEntry, error) {
	if f.loadErr != nil {
		return nil, f.loadErr
	}
	out := make([]ScoreEntry, len(f.entries))
	copy(out, f.entries)
	return out, nil
}

func (f *fakeScoreStore) Save(entries []ScoreEntry) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	f.entries = make([]ScoreEntry, len(entries))
	copy(f.entries, entries)
	return nil
}

// --- FakeWorld ---

type fakeWorld struct {
	tiles [][]TileType // [y][x]
	w, h  int
}

func newFakeWorld(w, h int) *fakeWorld {
	tiles := make([][]TileType, h)
	for y := range tiles {
		tiles[y] = make([]TileType, w)
	}
	return &fakeWorld{tiles: tiles, w: w, h: h}
}

func (f *fakeWorld) TileAt(x, y int) TileType {
	if x < 0 || y < 0 || x >= f.w || y >= f.h {
		return TileTree
	}
	return f.tiles[y][x]
}

func (f *fakeWorld) Width() int  { return f.w }
func (f *fakeWorld) Height() int { return f.h }

func (f *fakeWorld) set(x, y int, t TileType) {
	f.tiles[y][x] = t
}
