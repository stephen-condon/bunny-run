package game

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/stephen-condon/bunny-run/screens"
)

const (
	ScreenWidth  = 1280
	ScreenHeight = 720

	bunnyStartX = 5
	bunnyStartY = WorldHeight / 2
)

// tileColor is used as a fallback when no emoji face is available.
var tileColor = [...]color.RGBA{
	TileEmpty:     {34, 85, 34, 255},
	TileTree:      {0, 100, 0, 255},
	TileBush:      {50, 150, 50, 255},
	TileBoulder:   {120, 120, 120, 255},
	TileFallenLog: {139, 90, 43, 255},
}

type screenState int

const (
	stateMenu screenState = iota
	statePlaying
	stateGameOver
	stateLeaderboard
)

// Game implements ebiten.Game.
type Game struct {
	state   screenState
	face    text.Face // UI text (menus, HUD)
	sprites *Sprites  // tile and entity rendering
	input   InputSource
	store   ScoreStore
	clock   Clock

	menu *screens.Menu

	world      *World
	camera     Camera
	bunny      *Bunny
	foxes      []*Fox
	difficulty Difficulty
	startTime  time.Time
	foxTimer   float64

	gameOver *screens.GameOver
	seconds  int

	leaderboard *screens.Leaderboard
}

// NewGame creates a fully wired game instance.
// sprites may be nil; tile rendering falls back to colored rectangles.
func NewGame(face text.Face, sprites *Sprites, store ScoreStore, clock Clock) *Game {
	g := &Game{
		state:   stateMenu,
		face:    face,
		sprites: sprites,
		input:   KeyboardInput{},
		store:   store,
		clock:   clock,
	}
	g.menu = screens.NewMenu(face)
	return g
}

func (g *Game) Layout(_, _ int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Update() error {
	switch g.state {
	case stateMenu:
		g.updateMenu()
	case statePlaying:
		g.updatePlaying()
	case stateGameOver:
		g.updateGameOver()
	case stateLeaderboard:
		g.updateLeaderboard()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case stateMenu:
		g.menu.Draw(screen)
	case statePlaying:
		g.drawPlaying(screen)
	case stateGameOver:
		g.gameOver.Draw(screen)
	case stateLeaderboard:
		g.leaderboard.Draw(screen)
	}
}

// --- Menu ---

func (g *Game) updateMenu() {
	choice := g.menu.Update(
		g.input.IsUpPressed(),
		g.input.IsDownPressed(),
		g.input.IsActionPressed(),
	)
	switch choice {
	case screens.MenuPlay:
		g.startGame()
	case screens.MenuLeaderboard:
		g.openLeaderboard()
	}
}

// --- Game start ---

func (g *Game) startGame() {
	g.world = NewWorld(time.Now().UnixNano())
	g.world.EnsureGenerated(ScreenWidth/TileSize + genBufferCols)
	g.camera = Camera{}
	g.bunny = NewBunny(bunnyStartX, bunnyStartY)
	g.foxes = nil
	g.difficulty = NewDifficulty()
	g.startTime = g.clock.Now()
	g.foxTimer = 0
	g.state = statePlaying
}

// --- Playing ---

func (g *Game) updatePlaying() {
	const dt = 1.0 / 60.0

	elapsed := g.clock.Now().Sub(g.startTime)
	g.difficulty.Update(elapsed)
	g.world.SetDensity(g.difficulty.ObstacleDensity)

	g.camera.Update(g.difficulty.ScrollSpeed, dt)

	rightCol := g.camera.RightTile(ScreenWidth) + genBufferCols
	g.world.EnsureGenerated(rightCol)
	g.world.Evict(g.camera.LeftTile() - 5)

	foxPositions := make([]Vec2, len(g.foxes))
	for i, f := range g.foxes {
		foxPositions[i] = f.Pos
	}

	g.bunny.Update(g.input, g.world, foxPositions, dt)

	for _, f := range g.foxes {
		others := make([]*Fox, 0, len(g.foxes)-1)
		for _, o := range g.foxes {
			if o != f {
				others = append(others, o)
			}
		}
		f.Update(g.world, g.bunny, others, dt)
	}

	g.foxTimer += dt
	if g.foxTimer >= g.difficulty.FoxInterval.Seconds() {
		g.foxTimer = 0
		g.spawnFox()
	}

	if g.bunny.IsCaughtByEdge(&g.camera) {
		g.triggerGameOver()
		return
	}
	for _, f := range g.foxes {
		if f.CatchesBunny(g.bunny) {
			g.triggerGameOver()
			return
		}
	}
}

func (g *Game) spawnFox() {
	spawnCol := g.camera.RightTile(ScreenWidth) + 1
	spawnRow := rand.Intn(WorldHeight)
	pos := Vec2{spawnCol, spawnRow}
	patrolA := Vec2{spawnCol - 5, spawnRow}
	if patrolA.X < g.camera.LeftTile() {
		patrolA.X = g.camera.LeftTile()
	}
	fox := NewFox(pos, PatrolPath{A: patrolA, B: pos}, time.Now().UnixNano())
	g.foxes = append(g.foxes, fox)

	alive := g.foxes[:0]
	for _, f := range g.foxes {
		if f.Pos.X >= g.camera.LeftTile()-2 {
			alive = append(alive, f)
		}
	}
	g.foxes = alive
}

func (g *Game) triggerGameOver() {
	g.seconds = int(g.clock.Now().Sub(g.startTime).Seconds())
	entries, _ := g.store.Load()
	isTop := IsTopScore(g.seconds, entries)
	g.gameOver = screens.NewGameOver(g.seconds, isTop, g.face)
	g.state = stateGameOver
}

// --- Drawing ---

// drawSprite draws a pre-rendered sprite image at (sx, sy), optionally dimmed.
func drawSprite(screen, img *ebiten.Image, sx, sy float64, alpha float32) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(sx, sy)
	if alpha < 1.0 {
		op.ColorScale.Scale(1, 1, 1, alpha)
	}
	screen.DrawImage(img, op)
}

func (g *Game) drawPlaying(screen *ebiten.Image) {
	screen.Fill(color.RGBA{34, 85, 34, 255})

	leftCol := g.camera.LeftTile()
	rightCol := g.camera.RightTile(ScreenWidth) + 1

	for x := leftCol; x <= rightCol; x++ {
		for y := 0; y < WorldHeight; y++ {
			tile := g.world.TileAt(x, y)
			if tile == TileEmpty {
				continue
			}
			sx := g.camera.ScreenX(x)
			sy := float64(y * TileSize)
			if spr := g.sprites.TileSprite(tile); spr != nil {
				drawSprite(screen, spr, sx, sy, 1.0)
			} else {
				ebitenutil.DrawRect(screen, sx, sy, float64(TileSize-1), float64(TileSize-1), tileColor[tile])
			}
		}
	}

	// Bunny — semi-transparent when hidden in a bush.
	bsx := g.camera.ScreenX(g.bunny.Pos.X)
	bsy := float64(g.bunny.Pos.Y * TileSize)
	if g.sprites != nil {
		alpha := float32(1.0)
		if g.bunny.Hidden {
			alpha = 0.35
		}
		drawSprite(screen, g.sprites.Bunny, bsx, bsy, alpha)
	} else {
		bunnyColor := color.RGBA{255, 255, 255, 255}
		if g.bunny.Hidden {
			bunnyColor = color.RGBA{200, 200, 200, 120}
		}
		ebitenutil.DrawRect(screen, bsx+4, bsy+4, float64(TileSize-8), float64(TileSize-8), bunnyColor)
	}

	// Foxes.
	for _, f := range g.foxes {
		fsx := g.camera.ScreenX(f.Pos.X)
		fsy := float64(f.Pos.Y * TileSize)
		if g.sprites != nil {
			drawSprite(screen, g.sprites.Fox, fsx, fsy, 1.0)
		} else {
			ebitenutil.DrawRect(screen, fsx+4, fsy+4, float64(TileSize-8), float64(TileSize-8), color.RGBA{255, 140, 0, 255})
		}
	}

	// HUD.
	elapsed := int(g.clock.Now().Sub(g.startTime).Seconds())
	hud := &text.DrawOptions{}
	hud.ColorScale.ScaleWithColor(color.White)
	hud.GeoM.Translate(10, 10)
	text.Draw(screen, fmt.Sprintf("Time: %ds  Level: %d", elapsed, g.difficulty.Level()), g.face, hud)
}

// --- Game Over ---

func (g *Game) updateGameOver() {
	done := g.gameOver.Update(
		g.input.JustPressedChar(),
		g.input.IsBackspaceJustPressed(),
		g.input.IsActionPressed(),
	)
	if done {
		if g.gameOver.IsTopScore && g.gameOver.Phase == screens.GameOverDone {
			entries, _ := g.store.Load()
			entry := ScoreEntry{
				Name:    g.gameOver.GetName(),
				Seconds: g.seconds,
				Date:    g.clock.Now(),
			}
			entries = InsertScore(entry, entries)
			_ = g.store.Save(entries)
		}
		g.state = stateMenu
		g.menu = screens.NewMenu(g.face)
	}
}

// --- Leaderboard ---

func (g *Game) openLeaderboard() {
	entries, _ := g.store.Load()
	rows := make([]screens.LeaderboardRow, len(entries))
	for i, e := range entries {
		rows[i] = screens.LeaderboardRow{Name: e.Name, Seconds: e.Seconds}
	}
	g.leaderboard = screens.NewLeaderboard(rows, g.face)
	g.state = stateLeaderboard
}

func (g *Game) updateLeaderboard() {
	if g.leaderboard.Update(g.input.IsActionPressed()) {
		g.state = stateMenu
		g.menu = screens.NewMenu(g.face)
	}
}
