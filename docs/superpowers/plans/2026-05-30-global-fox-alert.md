# Global Fox Alert Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** When any fox spots the bunny, all foxes (regardless of distance or current state) immediately converge on the bunny's position.

**Architecture:** Two surgical changes in `game/fox.go`: remove the distance cap from `alertNearby` so every fox is alerted, and simplify `Alert` to unconditionally update `lastKnown` and enter Chase — even when already chasing or wandering.

**Tech Stack:** Go, standard `testing` package

---

### Task 1: Update tests to reflect new behavior (they will fail)

Two existing tests assert the old behavior and must be updated before we touch production code.

**Files:**
- Modify: `game/fox_test.go`

- [ ] **Step 1: Replace `TestFoxAlertDoesNotReachFarFox` with the inverted assertion**

In `game/fox_test.go`, replace the entire `TestFoxAlertDoesNotReachFarFox` function (lines 114–125) with:

```go
func TestFoxAlertReachesFarFox(t *testing.T) {
	w := newFakeWorld(30, 20)
	f1 := newTestFox(5, 5)
	f2 := newTestFox(20, 5) // well beyond old alertRadius of 6

	b := NewBunny(6, 5)
	f1.Update(w, b, []*Fox{f2}, 0, foxSpeed)

	if f2.State != FoxStateChase {
		t.Errorf("f2 should be alerted even when far away, got %v", f2.State)
	}
	if f2.lastKnown != (Vec2{6, 5}) {
		t.Errorf("f2.lastKnown should be bunny pos {6,5}, got %v", f2.lastKnown)
	}
}
```

- [ ] **Step 2: Replace `TestFoxAlertNoOpWhenAlreadyChasing` with the inverted assertion**

In `game/fox_test.go`, replace the entire `TestFoxAlertNoOpWhenAlreadyChasing` function (lines 169–177) with:

```go
func TestFoxAlertUpdatesLastKnownWhenChasing(t *testing.T) {
	f := newTestFox(5, 5)
	f.State = FoxStateChase
	f.lastKnown = Vec2{3, 3}
	f.Alert(Vec2{9, 9})
	if f.lastKnown != (Vec2{9, 9}) {
		t.Errorf("alert should update last known even when already chasing, got %v", f.lastKnown)
	}
	if f.State != FoxStateChase {
		t.Errorf("fox should remain in chase state, got %v", f.State)
	}
}
```

- [ ] **Step 3: Add test for alerting a wandering fox**

Append this new test at the bottom of `game/fox_test.go`:

```go
func TestFoxAlertWakesWanderingFox(t *testing.T) {
	f := newTestFox(5, 5)
	f.State = FoxStateWander
	f.wanderTimer = 4.0
	f.Alert(Vec2{7, 7})
	if f.State != FoxStateChase {
		t.Errorf("alerted wandering fox should enter chase, got %v", f.State)
	}
	if f.lastKnown != (Vec2{7, 7}) {
		t.Errorf("lastKnown should be {7,7}, got %v", f.lastKnown)
	}
}
```

- [ ] **Step 4: Run tests to confirm failures**

```
go test ./game/... -run "TestFoxAlertReachesFarFox|TestFoxAlertUpdatesLastKnownWhenChasing|TestFoxAlertWakesWanderingFox" -v
```

Expected: all three tests FAIL (the production code hasn't changed yet).

---

### Task 2: Implement the production code changes

**Files:**
- Modify: `game/fox.go`

- [ ] **Step 1: Delete the `alertRadius` constant**

In `game/fox.go`, remove line 11:

```go
alertRadius      = 6   // tiles (Chebyshev) for alerting nearby foxes
```

The constant block (lines 8–15) should become:

```go
const (
	foxSpeed         = 3.0 // tiles per second
	visionRadius     = 6   // tiles (Chebyshev), 360°
	bushWanderMinSec = 3.0
	bushWanderMaxSec = 5.0
	lostInterestDist = 1 // tiles — fox "arrives" at last known pos within this dist
)
```

- [ ] **Step 2: Simplify `alertNearby` to alert all foxes**

Replace the `alertNearby` function body:

```go
func (f *Fox) alertNearby(others []*Fox, pos Vec2) {
	for _, o := range others {
		if o == f {
			continue
		}
		o.Alert(pos)
	}
}
```

- [ ] **Step 3: Simplify `Alert` to unconditionally update**

Replace the `Alert` function body:

```go
// Alert forces this fox into chase mode toward pos.
func (f *Fox) Alert(pos Vec2) {
	f.lastKnown = pos
	f.State = FoxStateChase
}
```

- [ ] **Step 4: Run all tests**

```
go test ./game/... -v
```

Expected: all tests PASS.

- [ ] **Step 5: Commit**

```bash
git add game/fox.go game/fox_test.go
git commit -m "feat(fox): alert all foxes globally when bunny is spotted"
```
