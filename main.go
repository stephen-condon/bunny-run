package main

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
	"github.com/stephen-condon/bunny-run/game"
)

func main() {
	uiSrc, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
	}
	uiFace := &text.GoTextFace{Source: uiSrc, Size: 20}

	sprites := game.NewSprites()
	if sprites == nil {
		log.Println("emoji sprites failed to load — falling back to colored rectangles")
	}

	store := game.NewFileScoreStore()
	clock := game.RealClock{}
	g := game.NewGame(uiFace, sprites, store, clock)

	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("Bunny Run")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
