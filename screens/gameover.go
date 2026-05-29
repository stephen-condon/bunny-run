package screens

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// GameOverPhase tracks what the game over screen is doing.
type GameOverPhase int

const (
	GameOverEnteringName GameOverPhase = iota
	GameOverDone
)

// GameOver handles the game-over screen including optional name entry.
type GameOver struct {
	Seconds    int
	IsTopScore bool
	Name       [3]rune
	nameLen    int
	Phase      GameOverPhase
	face       text.Face
}

func NewGameOver(seconds int, isTop bool, face text.Face) *GameOver {
	g := &GameOver{Seconds: seconds, IsTopScore: isTop, face: face}
	if !isTop {
		g.Phase = GameOverDone
	}
	return g
}

// GetName returns the entered name as a 3-char string.
func (g *GameOver) GetName() string {
	name := make([]rune, 3)
	for i := 0; i < 3; i++ {
		if i < g.nameLen && g.Name[i] != 0 {
			name[i] = g.Name[i]
		} else {
			name[i] = '_'
		}
	}
	return string(name)
}

// Update processes name-entry input. Returns true when ready to proceed to menu.
func (g *GameOver) Update(ch rune, backspace, action bool) bool {
	if g.Phase == GameOverDone {
		return action
	}
	if backspace && g.nameLen > 0 {
		g.nameLen--
		g.Name[g.nameLen] = 0
	}
	if ch >= 'A' && ch <= 'Z' || ch >= 'a' && ch <= 'z' || ch >= '0' && ch <= '9' {
		if g.nameLen < 3 {
			if ch >= 'a' && ch <= 'z' {
				ch -= 32
			}
			g.Name[g.nameLen] = ch
			g.nameLen++
		}
	}
	if action && g.nameLen == 3 {
		g.Phase = GameOverDone
	}
	return false
}

// Draw renders the game over screen.
func (g *GameOver) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{80, 20, 20, 255})

	title := &text.DrawOptions{}
	title.ColorScale.ScaleWithColor(color.RGBA{255, 80, 80, 255})
	title.GeoM.Translate(490, 180)
	text.Draw(screen, "GAME OVER", g.face, title)

	score := &text.DrawOptions{}
	score.ColorScale.ScaleWithColor(color.White)
	score.GeoM.Translate(440, 260)
	text.Draw(screen, fmt.Sprintf("Survived: %d seconds", g.Seconds), g.face, score)

	if g.IsTopScore {
		top := &text.DrawOptions{}
		top.ColorScale.ScaleWithColor(color.RGBA{255, 220, 50, 255})
		top.GeoM.Translate(420, 320)
		text.Draw(screen, "NEW HIGH SCORE!", g.face, top)

		if g.Phase == GameOverEnteringName {
			prompt := &text.DrawOptions{}
			prompt.ColorScale.ScaleWithColor(color.White)
			prompt.GeoM.Translate(380, 390)
			text.Draw(screen, fmt.Sprintf("Enter your name: %s", g.GetName()), g.face, prompt)

			hint := &text.DrawOptions{}
			hint.ColorScale.ScaleWithColor(color.RGBA{200, 200, 200, 255})
			hint.GeoM.Translate(380, 440)
			text.Draw(screen, "Type 3 letters then Enter", g.face, hint)
			return
		}
	}

	cont := &text.DrawOptions{}
	cont.ColorScale.ScaleWithColor(color.RGBA{200, 200, 200, 255})
	cont.GeoM.Translate(420, 430)
	text.Draw(screen, "Press Enter to continue", g.face, cont)
}
