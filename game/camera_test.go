package game

import "testing"

func TestCameraUpdate(t *testing.T) {
	c := Camera{}
	c.Update(60, 1.0)
	if c.X != 60 {
		t.Errorf("expected X=60, got %v", c.X)
	}
	c.Update(60, 0.5)
	if c.X != 90 {
		t.Errorf("expected X=90, got %v", c.X)
	}
}

func TestCameraLeftTile(t *testing.T) {
	c := Camera{X: 64}
	if c.LeftTile() != 2 {
		t.Errorf("expected LeftTile=2, got %d", c.LeftTile())
	}
}

func TestCameraRightTile(t *testing.T) {
	c := Camera{X: 0}
	// 1280 / 32 = 40
	if c.RightTile(1280) != 40 {
		t.Errorf("expected RightTile=40, got %d", c.RightTile(1280))
	}
}

func TestCameraIsCaughtByLeftEdge(t *testing.T) {
	c := Camera{X: 100}
	// Tile at column 3 = pixels 96..127; left edge at X=100 → tile 3 right edge is 128 > 100, so not caught
	if c.IsCaughtByLeftEdge(3) {
		t.Error("tile 3 should not be caught when camera is at X=100")
	}
	// Tile at column 2 = pixels 64..95; right edge 96 ≤ 100 → caught
	if !c.IsCaughtByLeftEdge(2) {
		t.Error("tile 2 should be caught when camera is at X=100")
	}
}

func TestCameraWorldToScreen(t *testing.T) {
	c := Camera{X: 100}
	if c.WorldToScreen(150) != 50 {
		t.Errorf("expected 50, got %v", c.WorldToScreen(150))
	}
}

func TestCameraScreenX(t *testing.T) {
	c := Camera{X: 64}
	// tile col 3 → world pixel 96 → screen pixel 32
	if c.ScreenX(3) != 32 {
		t.Errorf("expected 32, got %v", c.ScreenX(3))
	}
}
