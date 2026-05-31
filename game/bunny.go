package game

// BunnySpeed is tiles per second for held-key movement.
const BunnySpeed = 10.0

// Bunny is the player character.
type Bunny struct {
	Pos    Vec2
	Hidden bool // true when overlapping a bush and no fox within 3 tiles

	moveAccum float64
	lastDir   Vec2 // previous frame's direction; direction change triggers instant move
}

func NewBunny(x, y int) *Bunny {
	return &Bunny{Pos: Vec2{x, y}}
}

// Update processes input and moves the bunny. dt is seconds since last frame.
// rightBound is the exclusive maximum tile column the bunny may occupy.
// foxPositions is used to determine whether bush concealment holds.
// Returns true if the bunny moved this frame.
func (b *Bunny) Update(input InputSource, world WorldReader, foxPositions []Vec2, rightBound int, dt float64) bool {
	moved := false

	dx, dy := 0, 0
	if input.IsUpPressed() {
		dy = -1
	} else if input.IsDownPressed() {
		dy = 1
	} else if input.IsLeftPressed() {
		dx = -1
	} else if input.IsRightPressed() {
		dx = 1
	}

	newDir := Vec2{dx, dy}

	if dx == 0 && dy == 0 {
		b.moveAccum = 0
		b.lastDir = Vec2{}
	} else if newDir != b.lastDir {
		// Direction changed (or key first pressed): move immediately this frame.
		b.lastDir = newDir
		b.moveAccum = 0
		nx, ny := b.Pos.X+dx, b.Pos.Y+dy
		if CanBunnyEnter(nx, ny, world) && nx < rightBound {
			b.Pos.X = nx
			b.Pos.Y = ny
			moved = true
		}
	} else {
		// Same direction held: accumulate at BunnySpeed tiles/sec.
		b.moveAccum += BunnySpeed * dt
		for b.moveAccum >= 1.0 {
			b.moveAccum -= 1.0
			nx, ny := b.Pos.X+dx, b.Pos.Y+dy
			if CanBunnyEnter(nx, ny, world) && nx < rightBound {
				b.Pos.X = nx
				b.Pos.Y = ny
				moved = true
			}
		}
	}

	b.updateHiddenState(world, foxPositions)
	return moved
}

func (b *Bunny) updateHiddenState(world WorldReader, foxPositions []Vec2) {
	if !world.TileAt(b.Pos.X, b.Pos.Y).IsBush() {
		b.Hidden = false
		return
	}
	for _, fp := range foxPositions {
		if ManhattanDist(b.Pos, fp) <= 3 {
			b.Hidden = false
			return
		}
	}
	b.Hidden = true
}

// IsCaughtByEdge returns true if the camera's left edge has overtaken the bunny.
func (b *Bunny) IsCaughtByEdge(cam *Camera) bool {
	return cam.IsCaughtByLeftEdge(b.Pos.X)
}
