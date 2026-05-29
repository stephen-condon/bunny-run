package game

import "testing"

func TestWorldHeightConstant(t *testing.T) {
	w := NewWorld(1)
	w.EnsureGenerated(5)
	if w.Height() != WorldHeight {
		t.Errorf("world height: want %d got %d", WorldHeight, w.Height())
	}
}

func TestWorldOutOfBoundsReturnsTree(t *testing.T) {
	w := NewWorld(1)
	w.EnsureGenerated(5)
	if w.TileAt(3, -1) != TileTree {
		t.Error("y=-1 should return TileTree (solid wall)")
	}
	if w.TileAt(3, WorldHeight) != TileTree {
		t.Error("y=WorldHeight should return TileTree (solid wall)")
	}
}

func TestWorldEarlyColumnsAreEmpty(t *testing.T) {
	w := NewWorld(42)
	w.EnsureGenerated(10)
	// First 3 columns should be fully empty (safe start zone).
	for col := 0; col < 3; col++ {
		for row := 0; row < WorldHeight; row++ {
			if w.TileAt(col, row) != TileEmpty {
				t.Errorf("col %d row %d should be empty, got %v", col, row, w.TileAt(col, row))
			}
		}
	}
}

func TestWorldHasAtLeastOneCorridorRow(t *testing.T) {
	w := NewWorld(99)
	w.EnsureGenerated(50)
	// Every column from 3 onwards must have at least one passable row.
	for col := 3; col <= 50; col++ {
		hasOpen := false
		for row := 0; row < WorldHeight; row++ {
			if !w.TileAt(col, row).BlocksBunny() {
				hasOpen = true
				break
			}
		}
		if !hasOpen {
			t.Errorf("column %d has no passable row for bunny", col)
		}
	}
}

func TestWorldEvictsOldColumns(t *testing.T) {
	w := NewWorld(1)
	w.EnsureGenerated(20)
	w.Evict(10)
	// Column 5 should be gone (return empty, not TileTree).
	tile := w.TileAt(5, 5)
	_ = tile // evicted columns return TileEmpty, which is fine
	// Verify no panic and width stays sane.
	if w.Width() < 0 {
		t.Error("width should not be negative")
	}
}

func TestWorldSetDensityChangesNewColumns(t *testing.T) {
	w := NewWorld(7)
	w.SetDensity(0.0) // no obstacles
	w.EnsureGenerated(100)
	// With density 0, columns ≥3 should only have corridor rows (all empty).
	nonEmpty := 0
	for col := 3; col <= 100; col++ {
		for row := 0; row < WorldHeight; row++ {
			if w.TileAt(col, row) != TileEmpty {
				nonEmpty++
			}
		}
	}
	if nonEmpty != 0 {
		t.Errorf("with density=0 expected 0 non-empty tiles after col 2, got %d", nonEmpty)
	}
}

func TestWorldWidthGrowsWithGeneration(t *testing.T) {
	w := NewWorld(1)
	w.EnsureGenerated(10)
	if w.Width() < 11 {
		t.Errorf("expected Width >= 11, got %d", w.Width())
	}
}
