package screens

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// LeaderboardRow is a single display row on the leaderboard.
type LeaderboardRow struct {
	Name    string
	Seconds int
	Score   int
}

// Leaderboard renders the top-10 scores.
type Leaderboard struct {
	rows []LeaderboardRow
	face text.Face
}

func NewLeaderboard(rows []LeaderboardRow, face text.Face) *Leaderboard {
	return &Leaderboard{rows: rows, face: face}
}

// Update returns true when the player presses a key to return to menu.
func (l *Leaderboard) Update(actionPressed bool) bool {
	return actionPressed
}

// Draw renders the leaderboard.
func (l *Leaderboard) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 60, 20, 255})

	title := &text.DrawOptions{}
	title.ColorScale.ScaleWithColor(color.White)
	title.GeoM.Translate(440, 60)
	text.Draw(screen, "TOP SCORES", l.face, title)

	header := &text.DrawOptions{}
	header.ColorScale.ScaleWithColor(color.RGBA{200, 200, 100, 255})
	header.GeoM.Translate(380, 130)
	text.Draw(screen, fmt.Sprintf("%-4s %-6s  %8s  %s", "#", "NAME", "SCORE", "TIME"), l.face, header)

	for i, row := range l.rows {
		o := &text.DrawOptions{}
		o.ColorScale.ScaleWithColor(color.White)
		o.GeoM.Translate(380, float64(170+i*40))
		text.Draw(screen, fmt.Sprintf("%-4d %-6s  %8d  %ds", i+1, row.Name, row.Score, row.Seconds), l.face, o)
	}

	back := &text.DrawOptions{}
	back.ColorScale.ScaleWithColor(color.RGBA{200, 200, 200, 255})
	back.GeoM.Translate(440, 620)
	text.Draw(screen, "Press Enter to return", l.face, back)
}
