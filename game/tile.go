package game

type TileType int

const (
	TileEmpty     TileType = iota
	TileTree               // blocks bunny + fox, blocks vision
	TileBush               // hides bunny (can still move), fox loses scent if >3 tiles away
	TileBoulder            // blocks bunny only
	TileFallenLog          // blocks fox only
)

type TileProps struct {
	BlocksBunny  bool
	BlocksFox    bool
	BlocksVision bool
	Emoji        string
}

var tileProps = [...]TileProps{
	TileEmpty:     {false, false, false, "·"},
	TileTree:      {true, true, true, "🌲"},
	TileBush:      {false, false, false, "🌿"},
	TileBoulder:   {true, false, false, "🪨"},
	TileFallenLog: {false, true, false, "🪵"},
}

func (t TileType) Props() TileProps {
	if int(t) >= len(tileProps) {
		return tileProps[TileEmpty]
	}
	return tileProps[t]
}

func (t TileType) BlocksBunny() bool  { return t.Props().BlocksBunny }
func (t TileType) BlocksFox() bool    { return t.Props().BlocksFox }
func (t TileType) BlocksVision() bool { return t.Props().BlocksVision }
func (t TileType) Emoji() string      { return t.Props().Emoji }
func (t TileType) IsBush() bool       { return t == TileBush }
