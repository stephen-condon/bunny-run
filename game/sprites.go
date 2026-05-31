package game

import (
	"embed"
	"image"
	"image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/emoji/*.png
var emojiFS embed.FS

// Sprites holds pre-scaled 32×32 tile and entity images decoded from embedded Noto Emoji PNGs.
type Sprites struct {
	tiles  [5]*ebiten.Image // indexed by TileType; nil for TileEmpty
	Bunny  *ebiten.Image
	Fox    *ebiten.Image
	Carrot *ebiten.Image
}

// NewSprites decodes the embedded emoji PNGs and scales each to TileSize×TileSize.
// Returns nil if any asset fails to decode.
func NewSprites() *Sprites {
	load := func(name string) *ebiten.Image {
		f, err := emojiFS.Open("assets/emoji/" + name)
		if err != nil {
			return nil
		}
		defer f.Close()
		src, err := png.Decode(f)
		if err != nil {
			return nil
		}
		return scaleTo(src, TileSize)
	}

	s := &Sprites{}
	s.tiles[TileTree] = load("tree.png")
	s.tiles[TileBush] = load("bush.png")
	s.tiles[TileBoulder] = load("boulder.png")
	s.tiles[TileFallenLog] = load("log.png")
	s.Bunny = load("bunny.png")
	s.Fox = load("fox.png")
	s.Carrot = load("carrot.png")

	// Return nil if any sprite failed — caller falls back to colored rects.
	for _, img := range s.tiles {
		if img == nil && s.tiles[TileTree] == nil {
			return nil
		}
	}
	if s.tiles[TileTree] == nil || s.tiles[TileBush] == nil ||
		s.tiles[TileBoulder] == nil || s.tiles[TileFallenLog] == nil ||
		s.Bunny == nil || s.Fox == nil || s.Carrot == nil {
		return nil
	}
	return s
}

// TileSprite returns the sprite for a tile, or nil for TileEmpty / out-of-range.
func (s *Sprites) TileSprite(t TileType) *ebiten.Image {
	if s == nil || t == TileEmpty || int(t) >= len(s.tiles) {
		return nil
	}
	return s.tiles[t]
}

// scaleTo scales src to a square of side `size` pixels using Ebiten's DrawImage.
func scaleTo(src image.Image, size int) *ebiten.Image {
	srcW := src.Bounds().Dx()
	srcH := src.Bounds().Dy()
	base := ebiten.NewImageFromImage(src)
	dst := ebiten.NewImage(size, size)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(size)/float64(srcW), float64(size)/float64(srcH))
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(base, op)
	return dst
}
