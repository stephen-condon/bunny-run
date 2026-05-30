package game

// Vec2 is a tile-grid coordinate.
type Vec2 struct{ X, Y int }

// InBounds returns true if v is within the world dimensions.
func InBounds(v Vec2, w WorldReader) bool {
	return v.X >= 0 && v.Y >= 0 && v.X < w.Width() && v.Y < w.Height()
}

// CanBunnyEnter returns true if the bunny may move onto tile (x,y).
func CanBunnyEnter(x, y int, w WorldReader) bool {
	if x < 0 || y < 0 || x >= w.Width() || y >= w.Height() {
		return false
	}
	return !w.TileAt(x, y).BlocksBunny()
}

// CanFoxEnter returns true if a fox may move onto tile (x,y).
func CanFoxEnter(x, y int, w WorldReader) bool {
	if x < 0 || y < 0 || x >= w.Width() || y >= w.Height() {
		return false
	}
	return !w.TileAt(x, y).BlocksFox()
}

// ChebyshevDist returns the Chebyshev (chessboard) distance between two tiles.
func ChebyshevDist(a, b Vec2) int {
	dx := a.X - b.X
	if dx < 0 {
		dx = -dx
	}
	dy := a.Y - b.Y
	if dy < 0 {
		dy = -dy
	}
	if dx > dy {
		return dx
	}
	return dy
}

// ManhattanDist returns the Manhattan distance between two tiles.
func ManhattanDist(a, b Vec2) int {
	dx := a.X - b.X
	if dx < 0 {
		dx = -dx
	}
	dy := a.Y - b.Y
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

// foxSegmentClear returns true when every tile on the axis-aligned segment from→to is fox-passable.
func foxSegmentClear(from, to Vec2, w WorldReader) bool {
	if from.X == to.X {
		minY, maxY := from.Y, to.Y
		if minY > maxY {
			minY, maxY = maxY, minY
		}
		for y := minY; y <= maxY; y++ {
			if !CanFoxEnter(from.X, y, w) {
				return false
			}
		}
	} else {
		minX, maxX := from.X, to.X
		if minX > maxX {
			minX, maxX = maxX, minX
		}
		for x := minX; x <= maxX; x++ {
			if !CanFoxEnter(x, from.Y, w) {
				return false
			}
		}
	}
	return true
}

// StepToward returns a unit step from src toward dst (one tile at a time).
// Prefers horizontal movement when dx == dy.
func StepToward(src, dst Vec2) Vec2 {
	dx := dst.X - src.X
	dy := dst.Y - src.Y
	step := src
	if dx != 0 {
		if dx > 0 {
			step.X++
		} else {
			step.X--
		}
	} else if dy != 0 {
		if dy > 0 {
			step.Y++
		} else {
			step.Y--
		}
	}
	return step
}
