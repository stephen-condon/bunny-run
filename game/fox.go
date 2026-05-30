package game

import (
	"math/rand"
	"time"
)

const (
	foxSpeed          = 3.0 // tiles per second
	visionConeDepth   = 5   // tiles
	visionConeHalfAng = 1   // cone half-width at each distance step (Manhattan spread)
	peripheralRadius  = 1   // tiles (Chebyshev)
	alertRadius       = 6   // tiles (Chebyshev) for alerting nearby foxes
	bushWanderMinSec  = 3.0
	bushWanderMaxSec  = 5.0
	lostInterestDist  = 1 // tiles — fox "arrives" at last known pos within this dist
)

// FoxState represents what the fox is currently doing.
type FoxState int

const (
	FoxStatePatrol FoxState = iota
	FoxStateChase
	FoxStateWander // lost bunny in a bush
)

// Dir represents a cardinal facing direction.
type Dir int

const (
	DirRight Dir = iota
	DirDown
	DirLeft
	DirUp
)

var dirVec = [4]Vec2{
	DirRight: {1, 0},
	DirDown:  {0, 1},
	DirLeft:  {-1, 0},
	DirUp:    {0, -1},
}

// PatrolPath defines the endpoints of a fox patrol segment.
type PatrolPath struct {
	A, B Vec2
}

// Fox is an enemy that patrols and chases the bunny.
type Fox struct {
	Pos          Vec2
	Facing       Dir
	State        FoxState
	patrol       PatrolPath
	patrolTarget Vec2 // current patrol waypoint (A or B)
	lastKnown    Vec2 // last known bunny position
	wanderTimer  float64
	wanderDir    Dir
	moveAccum    float64
	rng          *rand.Rand
}

// NewFox creates a fox at pos with the given patrol path.
func NewFox(pos Vec2, patrol PatrolPath, seed int64) *Fox {
	f := &Fox{
		Pos:          pos,
		Facing:       DirRight,
		State:        FoxStatePatrol,
		patrol:       patrol,
		patrolTarget: patrol.B,
		rng:          rand.New(rand.NewSource(seed)),
	}
	return f
}

// Update runs the fox AI for one frame. Returns true if the fox spotted the bunny.
// otherFoxes is used to propagate alerts. speed is the current effective tiles-per-second.
func (f *Fox) Update(world WorldReader, bunny *Bunny, otherFoxes []*Fox, dt float64, speed float64) bool {
	spotted := false

	if !bunny.Hidden {
		if f.canSee(bunny.Pos, world) {
			f.lastKnown = bunny.Pos
			f.State = FoxStateChase
			f.alertNearby(otherFoxes, bunny.Pos)
			spotted = true
		}
	}

	f.moveAccum += speed * dt

	for f.moveAccum >= 1.0 {
		f.moveAccum -= 1.0
		switch f.State {
		case FoxStatePatrol:
			f.stepPatrol(world)
		case FoxStateChase:
			f.stepChase(bunny, world)
		case FoxStateWander:
			f.wanderTimer -= 1.0 / speed
			if f.wanderTimer <= 0 {
				f.State = FoxStatePatrol
			} else {
				f.stepWander(world)
			}
		}
	}

	return spotted
}

// Alert forces this fox into chase mode toward pos.
func (f *Fox) Alert(pos Vec2) {
	if f.State == FoxStatePatrol {
		f.lastKnown = pos
		f.State = FoxStateChase
	}
}

func (f *Fox) alertNearby(others []*Fox, pos Vec2) {
	for _, o := range others {
		if o == f {
			continue
		}
		if ChebyshevDist(f.Pos, o.Pos) <= alertRadius {
			o.Alert(pos)
		}
	}
}

func (f *Fox) canSee(target Vec2, world WorldReader) bool {
	// Peripheral circle: always detect within 1 tile.
	if ChebyshevDist(f.Pos, target) <= peripheralRadius {
		return true
	}
	// Vision cone: cast tiles in facing direction.
	dv := dirVec[f.Facing]
	for depth := 1; depth <= visionConeDepth; depth++ {
		cx := f.Pos.X + dv.X*depth
		cy := f.Pos.Y + dv.Y*depth

		// Check a spread of tiles perpendicular to the cone at this depth.
		for spread := -depth / 2; spread <= depth/2; spread++ {
			var tx, ty int
			if dv.X != 0 {
				tx = cx
				ty = cy + spread
			} else {
				tx = cx + spread
				ty = cy
			}
			if tx == target.X && ty == target.Y {
				if !f.visionBlocked(target, world) {
					return true
				}
			}
		}

		// If the center tile blocks vision, stop the cone.
		if world.TileAt(cx, cy).BlocksVision() {
			break
		}
	}
	return false
}

// visionBlocked returns true if any tile along the straight line from fox to target blocks vision.
func (f *Fox) visionBlocked(target Vec2, world WorldReader) bool {
	x, y := f.Pos.X, f.Pos.Y
	dx := target.X - x
	dy := target.Y - y
	steps := dx
	if dx < 0 {
		steps = -dx
	}
	if abs(dy) > steps {
		steps = abs(dy)
	}
	for i := 1; i < steps; i++ {
		ix := x + dx*i/steps
		iy := y + dy*i/steps
		if world.TileAt(ix, iy).BlocksVision() {
			return true
		}
	}
	return false
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (f *Fox) stepPatrol(world WorldReader) {
	if f.Pos == f.patrolTarget {
		// Flip waypoint.
		if f.patrolTarget == f.patrol.B {
			f.patrolTarget = f.patrol.A
		} else {
			f.patrolTarget = f.patrol.B
		}
	}
	next := StepToward(f.Pos, f.patrolTarget)
	f.updateFacing(next)
	if CanFoxEnter(next.X, next.Y, world) {
		f.Pos = next
	}
}

func (f *Fox) stepChase(bunny *Bunny, world WorldReader) {
	// If bunny just entered a bush and we're not within 3 tiles, switch to wander.
	if bunny.Hidden && ManhattanDist(f.Pos, bunny.Pos) > 3 {
		dur := bushWanderMinSec + f.rng.Float64()*(bushWanderMaxSec-bushWanderMinSec)
		f.wanderTimer = dur
		f.State = FoxStateWander
		return
	}
	// Move toward last known position.
	if ManhattanDist(f.Pos, f.lastKnown) <= lostInterestDist {
		f.State = FoxStatePatrol
		return
	}
	next := StepToward(f.Pos, f.lastKnown)
	f.updateFacing(next)
	if CanFoxEnter(next.X, next.Y, world) {
		f.Pos = next
	}
}

func (f *Fox) stepWander(world WorldReader) {
	// Try the current wander direction; if blocked, pick a new random one.
	dv := dirVec[f.wanderDir]
	next := Vec2{f.Pos.X + dv.X, f.Pos.Y + dv.Y}
	if CanFoxEnter(next.X, next.Y, world) {
		f.Pos = next
		f.updateFacing(next)
	} else {
		f.wanderDir = Dir(f.rng.Intn(4))
	}
}

func (f *Fox) updateFacing(next Vec2) {
	dx := next.X - f.Pos.X
	dy := next.Y - f.Pos.Y
	switch {
	case dx > 0:
		f.Facing = DirRight
	case dx < 0:
		f.Facing = DirLeft
	case dy > 0:
		f.Facing = DirDown
	case dy < 0:
		f.Facing = DirUp
	}
}

// CatchesBunny returns true if the fox is on the same tile as the bunny.
func (f *Fox) CatchesBunny(bunny *Bunny) bool {
	return f.Pos == bunny.Pos
}

// WanderDuration returns a random wander duration for testing purposes.
func WanderDuration(rng *rand.Rand) time.Duration {
	secs := bushWanderMinSec + rng.Float64()*(bushWanderMaxSec-bushWanderMinSec)
	return time.Duration(secs * float64(time.Second))
}
