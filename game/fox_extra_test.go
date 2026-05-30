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

func TestFoxPatrolFlipsWaypoint(t *testing.T) {
	w := newFakeWorld(20, 20)
	f := newTestFox(5, 5)
	// Place fox exactly at patrol.B so it should flip to A.
	f.Pos = Vec2{8, 5}
	f.patrolTarget = Vec2{8, 5} // already at B

	b := NewBunny(0, 19) // outside vision radius
	f.Update(w, b, nil, 1.0/foxSpeed, foxSpeed)

	// After flipping, patrolTarget should be A.
	if f.patrolTarget != f.patrol.A {
		t.Errorf("patrol target should flip to A, got %v", f.patrolTarget)
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
	f := newTestFox(5, 5)
	f.patrolTarget = Vec2{8, 5}
	b := NewBunny(0, 19) // outside vision radius
	// At 2× fox speed with same dt that normally steps 1 tile, should step 2 tiles.
	f.Update(w, b, nil, 1.0/foxSpeed, foxSpeed*2)
	if f.Pos.X != 7 {
		t.Errorf("fox at 2x speed should move 2 tiles, got x=%d", f.Pos.X)
	}
}

func TestFoxDetectsBunnyBehindAt6Tiles(t *testing.T) {
	w := newFakeWorld(20, 20)
	f := newTestFox(11, 5)
	f.Facing = DirRight // facing away from bunny
	b := NewBunny(5, 5) // 6 tiles behind fox
	spotted := f.Update(w, b, nil, 0, foxSpeed)
	if !spotted {
		t.Error("fox should spot bunny 6 tiles behind (360° radius)")
	}
}

func TestFoxUpdateFacingAllDirections(t *testing.T) {
	w := newFakeWorld(20, 20)
	// Test each facing direction by placing bunny in a clear path.
	tests := []struct {
		foxPos    Vec2
		lastKnown Vec2
		wantFacing Dir
	}{
		{Vec2{5, 5}, Vec2{8, 5}, DirRight}, // lastKnown must be >lostInterestDist away
		{Vec2{5, 5}, Vec2{2, 5}, DirLeft},
		{Vec2{5, 5}, Vec2{5, 8}, DirDown},
		{Vec2{5, 5}, Vec2{5, 2}, DirUp},
	}
	for _, tc := range tests {
		f := newTestFox(tc.foxPos.X, tc.foxPos.Y)
		f.State = FoxStateChase
		f.lastKnown = tc.lastKnown
		b := NewBunny(0, 19) // out of range
		f.Update(w, b, nil, 1.0/foxSpeed, foxSpeed)
		if f.Facing != tc.wantFacing {
			t.Errorf("pos %v → %v: want facing %v, got %v", tc.foxPos, tc.lastKnown, tc.wantFacing, f.Facing)
		}
	}
}
