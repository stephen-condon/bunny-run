# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

Prior to implementing anything beyond a simple bugfix or config change, always use plan mode. When planning, use Opus, when implementing, use Sonnet. When researching specific facts, use a Haiku subagent. Any time there is an opportunity to do something in parallel, spawn subagents to accomplish it.

## Commands

```bash
# Install pre-commit hooks (run once after clone)
make setup

# Run the game
go run .

# Build
go build -o bunny-run .

# Run all tests
go test ./...

# Run tests with coverage (target: >80%)
go test ./game/... -coverprofile=cover.out && go tool cover -func=cover.out

# Run a single test
go test ./game/... -run TestFoxVisionConeDetection

# Run tests verbose
go test ./game/... -v
```

## Architecture

**Package layout:**
- `main.go` — Ebiten window init, font loading (goregular TTF), wires up `game.NewGame`
- `game/` — all game logic; the only package that may import Ebiten is `game.go` and `input.go`
- `screens/` — stateless UI renderers (menu, game over, leaderboard); **must not import `game`** to avoid a cycle

**Dependency injection:** The four interfaces in `game/interfaces.go` (`Clock`, `InputSource`, `ScoreStore`, `WorldReader`) are the seams that make everything testable. Production implementations live in `game/input.go` (`KeyboardInput`, `RealClock`) and `game/score.go` (`FileScoreStore`). Test doubles are inlined in `game/testhelpers_test.go` — there is no `fakes/` sub-package, which would recreate the import cycle.

**Game loop (`game/game.go`):** `Game` implements `ebiten.Game`. It owns a `screenState` enum and delegates to `updateMenu/updatePlaying/updateGameOver/updateLeaderboard`. Tiles and entities render as Noto Emoji sprites via `drawSprite` → `screen.DrawImage`; `ebitenutil.DrawRect` colored-rect branches are the nil-sprites fallback. HUD/UI text uses `text.Draw` with the goregular face.

**Emoji sprites (`game/sprites.go` + `game/assets/emoji/`):** Six Noto Emoji PNGs (Apache-2.0) are embedded at compile time via `go:embed`. `NewSprites()` decodes and scales each to 32×32 at startup. `game.NewGame` accepts a `*Sprites`; nil falls back to colored rects.

**World generation (`game/world.go`):** Columns are generated lazily into a `map[int][]TileType` as the camera scrolls right. Every generated column guarantees at least one passable row for the bunny (corridor). Old columns are evicted via `Evict(minCol)`. Obstacle density and fox spawn rate are driven by `game/difficulty.go`, which increments every 15 seconds.

**Fox AI (`game/fox.go`):** Three states — `Patrol` (walk A↔B), `Chase` (move to last-known bunny position), `Wander` (lost bunny in bush; random walk for 3–5 s). Detection uses a 5-tile forward vision cone (blocked by trees) plus a 1-tile peripheral circle. On spot, broadcasts an alert to all foxes within 6 tiles (Chebyshev). Bush concealment: if the bunny is on a bush tile AND no fox is within 3 tiles (Manhattan), `bunny.Hidden = true` — foxes skip detection entirely for hidden bunnies.

**Tile collision rules:**
| Tile       | Blocks Bunny | Blocks Fox | Blocks Vision |
|------------|:-----------:|:----------:|:-------------:|
| Tree       | ✅ | ✅ | ✅ |
| Bush       | ❌ | ❌ | ❌ |
| Boulder    | ✅ | ❌ | ❌ |
| FallenLog  | ❌ | ✅ | ❌ |

**Scores** are persisted to `~/.bunny-run/scores.json` as a JSON array sorted descending by seconds. Top-10 only; `InsertScore` handles trim. `IsTopScore` is the gate before prompting for a name.

**Screen resolution:** 1280×720, tile size 32 px → 40×22 tile grid. `Camera.X` is the only scroll state; `IsCaughtByLeftEdge` implements the left-edge death condition.

## Release Tagging
When ready to ship a release use:
`git tag v1.0.0`
`git push origin v1.0.0`

Ensure you follow semantic versioning based on the conventional commits scheme. Ship a release when the version would increment based on semantic versioning analysis of the merged PR.