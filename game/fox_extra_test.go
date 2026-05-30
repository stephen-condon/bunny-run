package game

import "testing"

func TestFoxStepWander(t *testing.T) {
	w := newFakeWorld(20, 20)
	f := newTestFox(5, 5)
	f.State = FoxStateWander
	f.wanderTimer = 10.0
	f.wanderDir = DirRight

	startX := f.Pos.X
	f.Update(w, NewBunny(0, 19), nil, 1.0/foxSpeed, foxSpeed) // bunny outside vision radius

	// Fox should have moved (or changed direction if blocked).
	// Since the world is clear, it should move right.
	if f.Pos.X != startX+1 {
		t.Errorf("wandering fox should move right, got x=%d (started %d)", f.Pos.X, startX)
	}
}

func TestFoxWanderChangesDirectionWhenBlocked(t *testing.T) {
	w := newFakeWorld(20, 20)
	w.set(6, 5, TileTree) // block right
	f := newTestFox(5, 5)
	f.State = FoxStateWander
	f.wanderTimer = 10.0
	f.wanderDir = DirRight
	f.rng.Seed(1) // deterministic direction change

	f.Update(w, NewBunny(0, 19), nil, 1.0/foxSpeed, foxSpeed) // bunny outside vision radius
	// Fox could not go right, must have changed direction.
	if f.Pos.X == 6 {
		t.Error("fox should not enter tree tile")
	}
}

func TestFoxPatrolAdvancesWaypoint(t *testing.T) {
	w := newFakeWorld(20, 20)
	f := newTestFox(5, 5)
	// Place fox at W1={8,5} with patrolIdx=1; on next step it should advance to idx 2.
	f.Pos = Vec2{8, 5}
	f.patrolIdx = 1

	b := NewBunny(0, 19) // outside vision radius
	f.Update(w, b, nil, 1.0/foxSpeed, foxSpeed)

	if f.patrolIdx != 2 {
		t.Errorf("patrol index should advance to 2, got %d", f.patrolIdx)
	}
}

func TestFoxChaseLosesInterestAtLastKnown(t *testing.T) {
	w := newFakeWorld(20, 20)
	f := newTestFox(5, 5)
	f.State = FoxStateChase
	f.lastKnown = Vec2{5, 5} // fox IS at last known position
	b := NewBunny(0, 19)     // bunny outside vision radius

	f.Update(w, b, nil, 1.0/foxSpeed, foxSpeed)

	if f.State != FoxStatePatrol {
		t.Errorf("fox should return to patrol at last known pos, got %v", f.State)
	}
}

func TestWanderDuration(t *testing.T) {
	import_rand := newTestFox(5, 5).rng
	d := WanderDuration(import_rand)
	if d.Seconds() < bushWanderMinSec || d.Seconds() > bushWanderMaxSec {
		t.Errorf("wander duration %v out of range [%v, %v]", d.Seconds(), bushWanderMinSec, bushWanderMaxSec)
	}
}

func TestFoxMovesProportionallyWithScrollSpeed(t *testing.T) {
	w := newFakeWorld(20, 20)
	f := newTestFox(5, 5) // patrolIdx=1, target W1={8,5}
	b := NewBunny(0, 19)  // outside vision radius
	// At 2× fox speed with same dt that normally steps 1 tile, should step 2 tiles.
	f.Update(w, b, nil, 1.0/foxSpeed, foxSpeed*2)
	if f.Pos.X != 7 {
		t.Errorf("fox at 2x speed should move 2 tiles, got x=%d", f.Pos.X)
	}
}

func TestFoxDetectsBunnyBehindAt6Tiles(t *testing.T) {
	w := newFakeWorld(20, 20)
	f := newTestFox(11, 5)
	b := NewBunny(5, 5) // 6 tiles behind fox
	spotted := f.Update(w, b, nil, 0, foxSpeed)
	if !spotted {
		t.Error("fox should spot bunny 6 tiles behind (360° radius)")
	}
}
