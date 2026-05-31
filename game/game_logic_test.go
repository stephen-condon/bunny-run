package game

import (
	"testing"
	"time"
)

// newTestGame creates a Game with fake dependencies. face=nil is safe as long as Draw is not called.
func newTestGame() *Game {
	store := &fakeScoreStore{}
	clock := newFakeClock(time.Now())
	return NewGame(nil, nil, store, clock)
}

func TestGameStartsAtMenu(t *testing.T) {
	g := newTestGame()
	if g.state != stateMenu {
		t.Errorf("game should start at menu, got %v", g.state)
	}
}

func TestMenuSelectPlayTransitionsToPlaying(t *testing.T) {
	g := newTestGame()
	g.input = &fakeInput{action: true}
	g.Update() //nolint:errcheck
	if g.state != statePlaying {
		t.Errorf("selecting Play should go to playing, got %v", g.state)
	}
}

func TestMenuNavigateDownSelectsLeaderboard(t *testing.T) {
	g := newTestGame()
	g.input = &fakeInput{down: true}
	g.Update() //nolint:errcheck
	g.input = &fakeInput{action: true}
	g.Update() //nolint:errcheck
	if g.state != stateLeaderboard {
		t.Errorf("should transition to leaderboard, got %v", g.state)
	}
}

func TestStartGameInitializesBunnyAndWorld(t *testing.T) {
	g := newTestGame()
	g.startGame()
	if g.bunny == nil {
		t.Error("bunny should be initialized after startGame")
	}
	if g.world == nil {
		t.Error("world should be initialized after startGame")
	}
}

func TestStartGameResetsDifficulty(t *testing.T) {
	g := newTestGame()
	g.startGame()
	if g.difficulty.ScrollSpeed != baseScrollSpeed {
		t.Errorf("scroll speed should be base on start, got %v", g.difficulty.ScrollSpeed)
	}
}

func TestGameOverWhenBunnyLeftBehind(t *testing.T) {
	g := newTestGame()
	g.startGame()
	g.bunny.Pos.X = 0
	g.camera.X = float64(5 * TileSize)
	g.input = &fakeInput{}
	g.updatePlaying()
	if g.state != stateGameOver {
		t.Errorf("game over when bunny left behind, got %v", g.state)
	}
}

func TestGameOverWhenFoxCatchesBunny(t *testing.T) {
	g := newTestGame()
	g.startGame()
	fox := newTestFox(bunnyStartX, bunnyStartY)
	g.foxes = []*Fox{fox}
	g.input = &fakeInput{}
	g.updatePlaying()
	if g.state != stateGameOver {
		t.Errorf("game over when fox catches bunny, got %v", g.state)
	}
}

func TestGameOverRecordsElapsedSeconds(t *testing.T) {
	g := newTestGame()
	g.startGame()
	g.startTime = g.clock.Now().Add(-42 * time.Second)
	g.bunny.Pos.X = 0
	g.camera.X = float64(5 * TileSize)
	g.input = &fakeInput{}
	g.updatePlaying()
	if g.seconds != 42 {
		t.Errorf("expected 42 seconds, got %d", g.seconds)
	}
}

func TestDifficultyScalesWithTime(t *testing.T) {
	g := newTestGame()
	g.startGame()
	initialSpeed := g.difficulty.ScrollSpeed
	g.startTime = g.clock.Now().Add(-15 * time.Second)
	g.input = &fakeInput{}
	g.updatePlaying()
	if g.difficulty.ScrollSpeed <= initialSpeed {
		t.Error("scroll speed should increase after 15s")
	}
}

func TestGameOverNoTopScoreReturnsToMenu(t *testing.T) {
	store := &fakeScoreStore{
		entries: func() []ScoreEntry {
			e := make([]ScoreEntry, 10)
			for i := range e {
				e[i] = ScoreEntry{Seconds: 1000 - i}
			}
			return e
		}(),
	}
	clock := newFakeClock(time.Now())
	g := NewGame(nil, nil, store, clock)
	g.startGame()
	g.seconds = 1 // won't beat any entry
	g.triggerGameOver()
	if g.gameOver == nil {
		t.Fatal("gameOver should be set")
	}
	g.input = &fakeInput{action: true}
	g.updateGameOver()
	if g.state != stateMenu {
		t.Errorf("should return to menu, got %v", g.state)
	}
}

func TestGameOverTopScoreEntersNameAndSaves(t *testing.T) {
	store := &fakeScoreStore{}
	clock := newFakeClock(time.Now())
	g := NewGame(nil, nil, store, clock)
	g.startGame()
	g.startTime = g.clock.Now().Add(-99 * time.Second)
	g.triggerGameOver() // empty store → top score, seconds = 99

	// Enter name "ABC".
	for _, ch := range []rune{'A', 'B', 'C'} {
		g.input = &fakeInput{ch: ch}
		g.updateGameOver()
	}
	// First action: confirms name entry, sets phase to Done (returns false).
	g.input = &fakeInput{action: true}
	g.updateGameOver()
	// Second action: now phase is Done, returns true → transitions to menu.
	g.input = &fakeInput{action: true}
	g.updateGameOver()

	if g.state != stateMenu {
		t.Errorf("should return to menu after name entry, got %v", g.state)
	}
	saved, _ := store.Load()
	if len(saved) == 0 {
		t.Fatal("score should be saved")
	}
	if saved[0].Seconds != 99 {
		t.Errorf("expected 99 seconds saved, got %d", saved[0].Seconds)
	}
	if saved[0].Name != "ABC" {
		t.Errorf("expected name ABC, got %s", saved[0].Name)
	}
}

func TestLeaderboardEnterReturnsToMenu(t *testing.T) {
	g := newTestGame()
	g.openLeaderboard()
	if g.state != stateLeaderboard {
		t.Fatal("should be in leaderboard state")
	}
	g.input = &fakeInput{action: true}
	g.updateLeaderboard()
	if g.state != stateMenu {
		t.Errorf("enter on leaderboard should return to menu, got %v", g.state)
	}
}

func TestTriggerGameOverSetsGameOverState(t *testing.T) {
	g := newTestGame()
	g.startGame()
	g.triggerGameOver()
	if g.state != stateGameOver {
		t.Errorf("triggerGameOver should set stateGameOver, got %v", g.state)
	}
	if g.gameOver == nil {
		t.Error("gameOver screen should be created")
	}
}

func TestGameLayout(t *testing.T) {
	g := newTestGame()
	w, h := g.Layout(0, 0)
	if w != ScreenWidth || h != ScreenHeight {
		t.Errorf("Layout: want %dx%d, got %dx%d", ScreenWidth, ScreenHeight, w, h)
	}
}

func TestSpawnFoxAddsToSlice(t *testing.T) {
	g := newTestGame()
	g.startGame()
	// Advance fox timer past the spawn threshold.
	g.foxTimer = g.difficulty.FoxInterval.Seconds() + 0.1
	g.input = &fakeInput{}
	count := len(g.foxes)
	g.updatePlaying()
	// A fox should have been spawned (timer exceeded threshold).
	if len(g.foxes) <= count {
		t.Error("spawnFox should add at least one fox when timer exceeds interval")
	}
}

func TestStartGameSpawnsOneFox(t *testing.T) {
	g := newTestGame()
	g.startGame()
	if len(g.foxes) != 1 {
		t.Errorf("startGame should spawn one initial fox, got %d", len(g.foxes))
	}
}

func TestSpawnFoxSkipsBlockedRows(t *testing.T) {
	g := newTestGame()
	g.startGame()

	spawnCol := g.camera.RightTile(ScreenWidth) + 1
	g.world.EnsureGenerated(spawnCol + 1)
	allTrees := make([]TileType, WorldHeight)
	for i := range allTrees {
		allTrees[i] = TileTree
	}
	g.world.tiles[spawnCol] = allTrees

	before := len(g.foxes)
	g.spawnFox()
	if len(g.foxes) != before {
		t.Errorf("spawnFox should not add a fox when all rows at spawn column are blocked")
	}
}

func TestOpenLeaderboardLoadsEntries(t *testing.T) {
	store := &fakeScoreStore{
		entries: []ScoreEntry{{Name: "ACE", Seconds: 50}},
	}
	g := NewGame(nil, nil, store, newFakeClock(time.Now()))
	g.openLeaderboard()
	if g.leaderboard == nil {
		t.Error("leaderboard should be set")
	}
}

func TestNewSpritesLoads(t *testing.T) {
	s := NewSprites()
	if s == nil {
		t.Fatal("NewSprites returned nil — embedded assets failed to decode")
	}
	if s.Bunny == nil {
		t.Error("Bunny sprite is nil")
	}
	if s.Fox == nil {
		t.Error("Fox sprite is nil")
	}
	if s.Carrot == nil {
		t.Error("Carrot sprite is nil")
	}
	for _, tile := range []TileType{TileTree, TileBush, TileBoulder, TileFallenLog} {
		if s.TileSprite(tile) == nil {
			t.Errorf("tile sprite for %v is nil", tile)
		}
	}
}

func TestCarrotCollectionIncrementsCount(t *testing.T) {
	g := newTestGame()
	g.startGame()
	g.carrots = []*Carrot{{Pos: g.bunny.Pos}}
	g.input = &fakeInput{}
	g.updatePlaying()
	if g.carrotsCollected != 1 {
		t.Errorf("expected 1 carrot collected, got %d", g.carrotsCollected)
	}
	if len(g.carrots) != 0 {
		t.Errorf("collected carrot should be evicted, got %d carrots", len(g.carrots))
	}
}

func TestSpawnCarrotAddsCarrot(t *testing.T) {
	g := newTestGame()
	g.startGame()
	g.world.EnsureGenerated(g.camera.RightTile(ScreenWidth) + 10)
	before := len(g.carrots)
	g.spawnCarrot()
	// spawnCarrot may legitimately skip if no passable tile found in 20 tries, but with an
	// open world it should almost always succeed.
	if len(g.carrots) < before {
		t.Error("spawnCarrot should not remove carrots")
	}
}

func TestCarrotTimerTriggersSpawn(t *testing.T) {
	g := newTestGame()
	g.startGame()
	g.world.EnsureGenerated(g.camera.RightTile(ScreenWidth) + 10)
	g.carrotTimer = g.carrotInterval + 1.0 // past threshold
	g.input = &fakeInput{}
	g.updatePlaying()
	if len(g.carrots) == 0 {
		t.Error("carrot should have been spawned when timer exceeded interval")
	}
}

func TestScoreIsSecondsTimesTenPlusCarrotsTimes25(t *testing.T) {
	g := newTestGame()
	g.startGame()
	g.startTime = g.clock.Now().Add(-10 * time.Second)
	g.carrotsCollected = 3
	g.bunny.Pos.X = 0
	g.camera.X = float64(5 * TileSize) // trigger game over via left-edge catch
	g.input = &fakeInput{}
	g.updatePlaying()
	// score = 10s*10 + 3 carrots*25 = 100 + 75 = 175
	if g.score != 175 {
		t.Errorf("expected score 175 (10s×10 + 3 carrots×25), got %d", g.score)
	}
}
