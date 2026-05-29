package game

const TileSize = 32 // pixels per tile

// Camera tracks the horizontal scroll position of the viewport.
type Camera struct {
	// X is the left edge of the viewport in pixel coordinates.
	X float64
}

// LeftTile returns the leftmost visible tile column (inclusive).
func (c *Camera) LeftTile() int {
	return int(c.X) / TileSize
}

// RightTile returns the rightmost visible tile column (inclusive) given screen width.
func (c *Camera) RightTile(screenW int) int {
	return (int(c.X) + screenW) / TileSize
}

// Update advances the camera by scrollSpeed pixels per second.
func (c *Camera) Update(scrollSpeed, dt float64) {
	c.X += scrollSpeed * dt
}

// IsCaughtByLeftEdge returns true if the bunny tile column is left of the camera.
func (c *Camera) IsCaughtByLeftEdge(bunnyTileX int) bool {
	return float64(bunnyTileX*TileSize+TileSize) <= c.X
}

// WorldToScreen converts a world pixel X to a screen pixel X.
func (c *Camera) WorldToScreen(worldX float64) float64 {
	return worldX - c.X
}

// ScreenX returns the screen-space X for a tile column.
func (c *Camera) ScreenX(tileX int) float64 {
	return float64(tileX*TileSize) - c.X
}
