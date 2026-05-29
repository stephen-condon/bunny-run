package game

import "math/rand"

const (
	WorldHeight   = 22 // tiles tall (matches 720px / 32px)
	genBufferCols = 10 // columns to generate ahead of right viewport edge
)

// World holds the tile grid and handles procedural column generation.
type World struct {
	tiles       map[int][]TileType // column index → column of tiles (length WorldHeight)
	generated   int                // rightmost column index generated so far
	density     float64
	rng         *rand.Rand
}

// NewWorld creates an empty world seeded with the given rng source.
func NewWorld(seed int64) *World {
	return &World{
		tiles:     make(map[int][]TileType),
		generated: -1,
		density:   0.10,
		rng:       rand.New(rand.NewSource(seed)),
	}
}

func (w *World) SetDensity(d float64) { w.density = d }

// TileAt implements WorldReader.
func (w *World) TileAt(x, y int) TileType {
	if y < 0 || y >= WorldHeight {
		return TileTree
	}
	col, ok := w.tiles[x]
	if !ok {
		return TileEmpty
	}
	return col[y]
}

func (w *World) Width() int  { return w.generated + 1 }
func (w *World) Height() int { return WorldHeight }

// EnsureGenerated generates columns up to at least targetCol.
func (w *World) EnsureGenerated(targetCol int) {
	for w.generated < targetCol {
		w.generated++
		w.tiles[w.generated] = w.generateColumn(w.generated)
	}
}

// Evict removes columns older than minCol to free memory.
func (w *World) Evict(minCol int) {
	for col := range w.tiles {
		if col < minCol {
			delete(w.tiles, col)
		}
	}
}

// generateColumn creates a random column, guaranteeing at least one clear corridor row.
func (w *World) generateColumn(col int) []TileType {
	// Leave the first two columns clear so the bunny starts in open space.
	if col < 3 {
		column := make([]TileType, WorldHeight)
		return column
	}

	column := make([]TileType, WorldHeight)

	// Pick the guaranteed clear row for this column (corridor meanders slightly).
	corridorRow := w.rng.Intn(WorldHeight)

	obstacleTiles := []TileType{TileTree, TileBoulder, TileFallenLog, TileBush}

	for y := 0; y < WorldHeight; y++ {
		if y == corridorRow {
			column[y] = TileEmpty
			continue
		}
		if w.rng.Float64() < w.density {
			column[y] = obstacleTiles[w.rng.Intn(len(obstacleTiles))]
		}
	}
	return column
}
