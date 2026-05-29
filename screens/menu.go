package screens

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// MenuChoice represents what the player selected.
type MenuChoice int

const (
	MenuNone        MenuChoice = iota
	MenuPlay
	MenuLeaderboard
)

// Menu is the main menu screen.
type Menu struct {
	selected int // 0 = Play, 1 = Top Scores
	face     text.Face
}

func NewMenu(face text.Face) *Menu {
	return &Menu{face: face}
}

// Update processes input and returns the player's choice (or MenuNone).
func (m *Menu) Update(upPressed, downPressed, actionPressed bool) MenuChoice {
	if upPressed && m.selected > 0 {
		m.selected--
	}
	if downPressed && m.selected < 1 {
		m.selected++
	}
	if actionPressed {
		if m.selected == 0 {
			return MenuPlay
		}
		return MenuLeaderboard
	}
	return MenuNone
}

// Draw renders the main menu.
func (m *Menu) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{34, 85, 34, 255})

	title := &text.DrawOptions{}
	title.ColorScale.ScaleWithColor(color.White)
	title.GeoM.Translate(440, 280)
	text.Draw(screen, "BUNNY RUN", m.face, title)

	items := []string{"> Play", "  Top Scores"}
	if m.selected == 1 {
		items = []string{"  Play", "> Top Scores"}
	}

	for i, item := range items {
		o := &text.DrawOptions{}
		o.ColorScale.ScaleWithColor(color.White)
		o.GeoM.Translate(490, float64(380+i*50))
		text.Draw(screen, item, m.face, o)
	}

	hint := &text.DrawOptions{}
	hint.ColorScale.ScaleWithColor(color.RGBA{180, 180, 180, 255})
	hint.GeoM.Translate(420, 530)
	text.Draw(screen, "Up/Down to move  Enter to select", m.face, hint)
}
