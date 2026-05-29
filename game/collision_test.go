package game

import "testing"

func newTestWorld() *fakeWorld { return newFakeWorld(10, 10) }

func TestCanBunnyEnter(t *testing.T) {
	w := newTestWorld()
	w.set(3, 3, TileTree)
	w.set(4, 4, TileBoulder)
	w.set(5, 5, TileBush)
	w.set(6, 6, TileFallenLog)

	if CanBunnyEnter(3, 3, w) {
		t.Error("bunny should not enter tree")
	}
	if CanBunnyEnter(4, 4, w) {
		t.Error("bunny should not enter boulder")
	}
	if !CanBunnyEnter(5, 5, w) {
		t.Error("bunny should enter bush")
	}
	if !CanBunnyEnter(6, 6, w) {
		t.Error("bunny should enter fallen log")
	}
	if CanBunnyEnter(-1, 0, w) {
		t.Error("bunny should not enter out of bounds")
	}
}

func TestCanFoxEnter(t *testing.T) {
	w := newTestWorld()
	w.set(3, 3, TileTree)
	w.set(4, 4, TileBoulder)
	w.set(5, 5, TileFallenLog)

	if CanFoxEnter(3, 3, w) {
		t.Error("fox should not enter tree")
	}
	if !CanFoxEnter(4, 4, w) {
		t.Error("fox should enter boulder")
	}
	if CanFoxEnter(5, 5, w) {
		t.Error("fox should not enter fallen log")
	}
	if CanFoxEnter(-1, 0, w) {
		t.Error("fox should not enter out of bounds")
	}
}

func TestChebyshevDist(t *testing.T) {
	tests := []struct {
		a, b Vec2
		want int
	}{
		{Vec2{0, 0}, Vec2{3, 1}, 3},
		{Vec2{0, 0}, Vec2{1, 3}, 3},
		{Vec2{5, 5}, Vec2{5, 5}, 0},
		{Vec2{0, 0}, Vec2{-2, 2}, 2},
	}
	for _, tc := range tests {
		if got := ChebyshevDist(tc.a, tc.b); got != tc.want {
			t.Errorf("ChebyshevDist(%v,%v) = %d, want %d", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestManhattanDist(t *testing.T) {
	tests := []struct {
		a, b Vec2
		want int
	}{
		{Vec2{0, 0}, Vec2{3, 4}, 7},
		{Vec2{5, 5}, Vec2{5, 5}, 0},
		{Vec2{0, 0}, Vec2{-1, -1}, 2},
	}
	for _, tc := range tests {
		if got := ManhattanDist(tc.a, tc.b); got != tc.want {
			t.Errorf("ManhattanDist(%v,%v) = %d, want %d", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestStepToward(t *testing.T) {
	tests := []struct {
		src, dst, want Vec2
	}{
		{Vec2{0, 0}, Vec2{5, 0}, Vec2{1, 0}},
		{Vec2{0, 0}, Vec2{0, 5}, Vec2{0, 1}},
		{Vec2{5, 0}, Vec2{0, 0}, Vec2{4, 0}},
		{Vec2{0, 5}, Vec2{0, 0}, Vec2{0, 4}},
		{Vec2{3, 3}, Vec2{3, 3}, Vec2{3, 3}}, // same tile, no movement
	}
	for _, tc := range tests {
		got := StepToward(tc.src, tc.dst)
		if got != tc.want {
			t.Errorf("StepToward(%v,%v) = %v, want %v", tc.src, tc.dst, got, tc.want)
		}
	}
}

func TestInBounds(t *testing.T) {
	w := newTestWorld()
	if !InBounds(Vec2{0, 0}, w) {
		t.Error("(0,0) should be in bounds")
	}
	if InBounds(Vec2{-1, 0}, w) {
		t.Error("(-1,0) should be out of bounds")
	}
	if InBounds(Vec2{10, 0}, w) {
		t.Error("(10,0) should be out of bounds")
	}
}
