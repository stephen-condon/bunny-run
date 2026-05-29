package game

import "testing"

func TestTileProps(t *testing.T) {
	tests := []struct {
		tile         TileType
		blocksBunny  bool
		blocksFox    bool
		blocksVision bool
	}{
		{TileEmpty, false, false, false},
		{TileTree, true, true, true},
		{TileBush, false, false, false},
		{TileBoulder, true, false, false},
		{TileFallenLog, false, true, false},
	}
	for _, tc := range tests {
		if got := tc.tile.BlocksBunny(); got != tc.blocksBunny {
			t.Errorf("%v BlocksBunny: want %v got %v", tc.tile, tc.blocksBunny, got)
		}
		if got := tc.tile.BlocksFox(); got != tc.blocksFox {
			t.Errorf("%v BlocksFox: want %v got %v", tc.tile, tc.blocksFox, got)
		}
		if got := tc.tile.BlocksVision(); got != tc.blocksVision {
			t.Errorf("%v BlocksVision: want %v got %v", tc.tile, tc.blocksVision, got)
		}
	}
}

func TestTileIsBush(t *testing.T) {
	if !TileBush.IsBush() {
		t.Error("TileBush.IsBush() should be true")
	}
	if TileTree.IsBush() {
		t.Error("TileTree.IsBush() should be false")
	}
}

func TestTileEmoji(t *testing.T) {
	if TileTree.Emoji() == "" {
		t.Error("TileTree.Emoji() should not be empty")
	}
}

func TestOutOfBoundsTileType(t *testing.T) {
	var invalid TileType = 99
	props := invalid.Props()
	if props.BlocksBunny || props.BlocksFox {
		t.Error("out-of-bounds tile should default to empty props")
	}
}
