package game

import "testing"

// moveN advances the bunny one tile per call, n times.
func moveN(b *Bunny, input *fakeInput, w *fakeWorld, n int) {
	for i := 0; i < n; i++ {
		b.Update(input, w, nil, 1000, 1.0/BunnySpeed)
	}
}

func TestBunnyMovesRight(t *testing.T) {
	w := newFakeWorld(10, 10)
	b := NewBunny(2, 5)
	moveN(b, &fakeInput{right: true}, w, 1)
	if b.Pos.X != 3 {
		t.Errorf("expected X=3, got %d", b.Pos.X)
	}
}

func TestBunnyMovesLeft(t *testing.T) {
	w := newFakeWorld(10, 10)
	b := NewBunny(5, 5)
	moveN(b, &fakeInput{left: true}, w, 1)
	if b.Pos.X != 4 {
		t.Errorf("expected X=4, got %d", b.Pos.X)
	}
}

func TestBunnyMovesUp(t *testing.T) {
	w := newFakeWorld(10, 10)
	b := NewBunny(5, 5)
	moveN(b, &fakeInput{up: true}, w, 1)
	if b.Pos.Y != 4 {
		t.Errorf("expected Y=4, got %d", b.Pos.Y)
	}
}

func TestBunnyMovesDown(t *testing.T) {
	w := newFakeWorld(10, 10)
	b := NewBunny(5, 5)
	moveN(b, &fakeInput{down: true}, w, 1)
	if b.Pos.Y != 6 {
		t.Errorf("expected Y=6, got %d", b.Pos.Y)
	}
}

func TestBunnyBlockedByTree(t *testing.T) {
	w := newFakeWorld(10, 10)
	w.set(4, 5, TileTree)
	b := NewBunny(3, 5)
	moveN(b, &fakeInput{right: true}, w, 1)
	if b.Pos.X != 3 {
		t.Errorf("bunny should be blocked by tree, got X=%d", b.Pos.X)
	}
}

func TestBunnyBlockedByBoulder(t *testing.T) {
	w := newFakeWorld(10, 10)
	w.set(4, 5, TileBoulder)
	b := NewBunny(3, 5)
	moveN(b, &fakeInput{right: true}, w, 1)
	if b.Pos.X != 3 {
		t.Errorf("bunny should be blocked by boulder, got X=%d", b.Pos.X)
	}
}

func TestBunnyPassesFallenLog(t *testing.T) {
	w := newFakeWorld(10, 10)
	w.set(4, 5, TileFallenLog)
	b := NewBunny(3, 5)
	moveN(b, &fakeInput{right: true}, w, 1)
	if b.Pos.X != 4 {
		t.Errorf("bunny should pass fallen log, got X=%d", b.Pos.X)
	}
}

func TestBunnyHiddenInBushNoFox(t *testing.T) {
	w := newFakeWorld(10, 10)
	w.set(3, 5, TileBush)
	b := NewBunny(3, 5)
	b.Update(&fakeInput{}, w, nil, 1000, 0)
	if !b.Hidden {
		t.Error("bunny should be hidden in bush with no foxes nearby")
	}
}

func TestBunnyNotHiddenOutsideBush(t *testing.T) {
	w := newFakeWorld(10, 10)
	b := NewBunny(3, 5)
	b.Update(&fakeInput{}, w, nil, 1000, 0)
	if b.Hidden {
		t.Error("bunny should not be hidden on empty tile")
	}
}

func TestBunnyConcealmentsBreaksWhenFoxClose(t *testing.T) {
	w := newFakeWorld(10, 10)
	w.set(3, 5, TileBush)
	b := NewBunny(3, 5)
	foxPositions := []Vec2{{4, 5}} // 1 tile away (within 3)
	b.Update(&fakeInput{}, w, foxPositions, 1000, 0)
	if b.Hidden {
		t.Error("bunny should not be hidden when fox is within 3 tiles")
	}
}

func TestBunnyHiddenWhenFoxFar(t *testing.T) {
	w := newFakeWorld(20, 20)
	w.set(3, 5, TileBush)
	b := NewBunny(3, 5)
	foxPositions := []Vec2{{10, 10}} // far away (Manhattan dist > 3)
	b.Update(&fakeInput{}, w, foxPositions, 1000, 0)
	if !b.Hidden {
		t.Error("bunny should be hidden when fox is far away")
	}
}

func TestBunnyHiddenLostOnBushExit(t *testing.T) {
	w := newFakeWorld(10, 10)
	w.set(3, 5, TileBush)
	b := NewBunny(3, 5)
	b.Update(&fakeInput{}, w, nil, 1000, 0)
	if !b.Hidden {
		t.Fatal("bunny should be hidden in bush")
	}
	moveN(b, &fakeInput{right: true}, w, 1)
	if b.Hidden {
		t.Error("bunny should lose hidden status after leaving bush")
	}
}

func TestBunnyCaughtByEdge(t *testing.T) {
	b := NewBunny(2, 5)
	cam := &Camera{X: float64(3 * TileSize)}
	if !b.IsCaughtByEdge(cam) {
		t.Error("bunny at tile 2 should be caught by camera at tile 3")
	}
}

func TestBunnyNotCaughtByEdge(t *testing.T) {
	b := NewBunny(5, 5)
	cam := &Camera{X: float64(2 * TileSize)}
	if b.IsCaughtByEdge(cam) {
		t.Error("bunny at tile 5 should not be caught by camera at tile 2")
	}
}

func TestBunnyBlockedByRightBound(t *testing.T) {
	w := newFakeWorld(20, 10)
	b := NewBunny(5, 5)
	// rightBound=6 means bunny cannot enter column 6 or beyond
	b.Update(&fakeInput{right: true}, w, nil, 6, 1.0/BunnySpeed)
	if b.Pos.X != 5 {
		t.Errorf("bunny should be blocked at right bound, got X=%d", b.Pos.X)
	}
}

func TestBunnyCanReachRightBoundMinusOne(t *testing.T) {
	w := newFakeWorld(20, 10)
	b := NewBunny(4, 5)
	// rightBound=6 means column 5 is the last allowed column
	b.Update(&fakeInput{right: true}, w, nil, 6, 1.0/BunnySpeed)
	if b.Pos.X != 5 {
		t.Errorf("bunny should move to column 5 (rightBound-1), got X=%d", b.Pos.X)
	}
}
