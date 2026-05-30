package game

import "testing"

func newTestFox(x, y int) *Fox {
	pos := Vec2{x, y}
	return NewFox(pos, PatrolPath{A: Vec2{x - 3, y}, B: Vec2{x + 3, y}}, 42)
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
	f := newTestFox(5, 5)
	f.patrolTarget = Vec2{8, 5}
	b := NewBunny(0, 19) // outside vision radius

	f.Update(w, b, nil, 1.0/foxSpeed, foxSpeed)

	if f.Pos.X != 6 {
		t.Errorf("fox should have moved right to 6, got %d", f.Pos.X)
	}
}

func TestFoxPeripheralDetection(t *testing.T) {
	w := newFakeWorld(20, 20)
	f := newTestFox(5, 5)
	f.Facing = DirRight

	b := NewBunny(6, 5) // 1 tile away — within peripheral radius
	spotted := f.Update(w, b, nil, 0, foxSpeed)
	if !spotted {
		t.Error("fox should spot bunny within peripheral radius (1 tile)")
	}
}

func TestFoxDetectsNearbyBunny(t *testing.T) {
	w := newFakeWorld(20, 20)
	f := newTestFox(5, 5)
	f.Facing = DirRight

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
	f.Facing = DirRight

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
	f.Facing = DirRight

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
	f.Facing = DirRight

	b := NewBunny(6, 5) // adjacent, spotted by peripheral
	f.Update(w, b, nil, 0, foxSpeed)
	if f.State != FoxStateChase {
		t.Errorf("fox should be in chase state after spotting bunny, got %v", f.State)
	}
}

func TestFoxAlertPropagates(t *testing.T) {
	w := newFakeWorld(20, 20)
	f1 := newTestFox(5, 5)
	f2 := newTestFox(8, 5) // within alertRadius
	f1.Facing = DirRight

	b := NewBunny(6, 5)
	f1.Update(w, b, []*Fox{f2}, 0, foxSpeed)

	if f2.State != FoxStateChase {
		t.Errorf("f2 should be alerted and chasing, got %v", f2.State)
	}
}

func TestFoxAlertDoesNotReachFarFox(t *testing.T) {
	w := newFakeWorld(30, 20)
	f1 := newTestFox(5, 5)
	f2 := newTestFox(20, 5) // beyond alertRadius
	f1.Facing = DirRight

	b := NewBunny(6, 5)
	f1.Update(w, b, []*Fox{f2}, 0, foxSpeed)

	if f2.State == FoxStateChase {
		t.Error("f2 should not be alerted when too far away")
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

func TestFoxAlertNoOpWhenAlreadyChasing(t *testing.T) {
	f := newTestFox(5, 5)
	f.State = FoxStateChase
	f.lastKnown = Vec2{3, 3}
	f.Alert(Vec2{9, 9})
	if f.lastKnown != (Vec2{3, 3}) {
		t.Error("alert should not override last known when already chasing")
	}
}
