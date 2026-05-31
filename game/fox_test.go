package game

import "testing"

func newTestFox(x, y int) *Fox {
	pos := Vec2{x, y}
	// Rectangle: W0={x,y}, W1={x+3,y}, W2={x+3,y+2}, W3={x,y+2}
	// patrolIdx starts at 1 so the fox heads toward W1={x+3,y} on spawn.
	patrol := PatrolPath{Waypoints: [4]Vec2{
		{x, y}, {x + 3, y}, {x + 3, y + 2}, {x, y + 2},
	}}
	return NewFox(pos, patrol, 42)
}

func TestFoxCatchesBunny(t *testing.T) {
	f := newTestFox(5, 5)
	b := NewBunny(5, 5)
	if !f.CatchesBunny(b) {
		t.Error("fox should catch bunny on same tile")
	}
}

func TestFoxDoesNotCatchBunnyOnDifferentTile(t *testing.T) {
	f := newTestFox(5, 5)
	b := NewBunny(6, 5)
	if f.CatchesBunny(b) {
		t.Error("fox should not catch bunny on different tile")
	}
}

func TestFoxPatrolMovesAlongPath(t *testing.T) {
	w := newFakeWorld(20, 20)
	f := newTestFox(5, 5) // patrolIdx=1, target W1={8,5}
	b := NewBunny(0, 19)  // outside vision radius

	f.Update(w, b, nil, 1.0/foxSpeed, foxSpeed)

	if f.Pos.X != 6 {
		t.Errorf("fox should have moved right to 6, got %d", f.Pos.X)
	}
}

func TestFoxDetectsNearbyBunny(t *testing.T) {
	w := newFakeWorld(20, 20)
	f := newTestFox(5, 5)

	b := NewBunny(9, 5) // 4 tiles away
	spotted := f.Update(w, b, nil, 0, foxSpeed)
	if !spotted {
		t.Error("fox should spot bunny within vision radius")
	}
}

func TestFoxVisionBlockedByTree(t *testing.T) {
	w := newFakeWorld(20, 20)
	w.set(7, 5, TileTree)
	f := newTestFox(5, 5)

	b := NewBunny(9, 5)
	spotted := f.Update(w, b, nil, 0, foxSpeed)
	if spotted {
		t.Error("fox vision should be blocked by tree")
	}
}

func TestFoxDoesNotSeeHiddenBunny(t *testing.T) {
	w := newFakeWorld(20, 20)
	w.set(7, 5, TileBush)
	f := newTestFox(5, 5)

	b := NewBunny(7, 5)
	b.Hidden = true
	spotted := f.Update(w, b, nil, 0, foxSpeed)
	if spotted {
		t.Error("fox should not spot hidden bunny in bush")
	}
}

func TestFoxChasesAfterSpotting(t *testing.T) {
	w := newFakeWorld(20, 20)
	f := newTestFox(5, 5)

	b := NewBunny(6, 5) // adjacent, spotted by peripheral
	f.Update(w, b, nil, 0, foxSpeed)
	if f.State != FoxStateChase {
		t.Errorf("fox should be in chase state after spotting bunny, got %v", f.State)
	}
}

func TestFoxAlertPropagates(t *testing.T) {
	w := newFakeWorld(20, 20)
	f1 := newTestFox(5, 5)
	f2 := newTestFox(8, 5) // alert propagates regardless of distance

	b := NewBunny(6, 5)
	f1.Update(w, b, []*Fox{f2}, 0, foxSpeed)

	if f2.State != FoxStateChase {
		t.Errorf("f2 should be alerted and chasing, got %v", f2.State)
	}
}

func TestFoxAlertReachesFarFox(t *testing.T) {
	w := newFakeWorld(30, 20)
	f1 := newTestFox(5, 5)
	f2 := newTestFox(20, 5) // well beyond old alertRadius of 6

	b := NewBunny(6, 5)
	f1.Update(w, b, []*Fox{f2}, 0, foxSpeed)

	if f2.State != FoxStateChase {
		t.Errorf("f2 should be alerted even when far away, got %v", f2.State)
	}
	if f2.lastKnown != (Vec2{6, 5}) {
		t.Errorf("f2.lastKnown should be bunny pos {6,5}, got %v", f2.lastKnown)
	}
}

func TestFoxWandersWhenBunnyHidesInBush(t *testing.T) {
	w := newFakeWorld(20, 20)
	w.set(9, 5, TileBush)
	f := newTestFox(5, 5)
	f.State = FoxStateChase
	f.lastKnown = Vec2{9, 5}

	b := NewBunny(9, 5)
	b.Hidden = true // fox is >3 tiles away (Manhattan dist = 4)
	// Pass enough dt to tick the movement accumulator once.
	f.Update(w, b, nil, 1.0/foxSpeed, foxSpeed)

	if f.State != FoxStateWander {
		t.Errorf("fox should enter wander state when bunny hides in bush, got %v", f.State)
	}
}

func TestFoxReturnsToPatrolAfterWander(t *testing.T) {
	w := newFakeWorld(20, 20)
	f := newTestFox(5, 5)
	f.State = FoxStateWander
	f.wanderTimer = 0.1 // just about to expire

	b := NewBunny(0, 19) // outside vision radius
	f.Update(w, b, nil, 1.0/foxSpeed, foxSpeed)

	if f.State != FoxStatePatrol {
		t.Errorf("fox should return to patrol after wander, got %v", f.State)
	}
}

func TestFoxAlert(t *testing.T) {
	f := newTestFox(5, 5)
	f.Alert(Vec2{8, 8})
	if f.State != FoxStateChase {
		t.Errorf("alerted fox should be chasing, got %v", f.State)
	}
	if f.lastKnown != (Vec2{8, 8}) {
		t.Errorf("last known should be {8,8}, got %v", f.lastKnown)
	}
}

func TestFoxAlertUpdatesLastKnownWhenChasing(t *testing.T) {
	f := newTestFox(5, 5)
	f.State = FoxStateChase
	f.lastKnown = Vec2{3, 3}
	f.Alert(Vec2{9, 9})
	if f.lastKnown != (Vec2{9, 9}) {
		t.Errorf("alert should update last known even when already chasing, got %v", f.lastKnown)
	}
	if f.State != FoxStateChase {
		t.Errorf("fox should remain in chase state, got %v", f.State)
	}
}

func TestFoxAlertWakesWanderingFox(t *testing.T) {
	f := newTestFox(5, 5)
	f.State = FoxStateWander
	f.wanderTimer = 4.0
	f.Alert(Vec2{7, 7})
	if f.State != FoxStateChase {
		t.Errorf("alerted wandering fox should enter chase, got %v", f.State)
	}
	if f.lastKnown != (Vec2{7, 7}) {
		t.Errorf("lastKnown should be {7,7}, got %v", f.lastKnown)
	}
}

func TestFoxEntersBushToCatchBunny(t *testing.T) {
	w := newFakeWorld(20, 20)
	w.set(6, 5, TileBush)
	f := newTestFox(5, 5)

	b := NewBunny(6, 5)
	b.Hidden = false // fox is within 3 tiles, so bunny is not hidden

	// One move step: fox detects bunny (Chebyshev dist 1 ≤ 6, not hidden),
	// enters Chase with lastKnown={6,5}, then must step onto the bush tile.
	f.Update(w, b, nil, 1.0/foxSpeed, foxSpeed)

	if !f.CatchesBunny(b) {
		t.Error("fox should have entered bush and caught bunny")
	}
}
