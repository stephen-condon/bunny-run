# Global Fox Alert Design

**Date:** 2026-05-30  
**Status:** Approved

## Problem

When one fox spots the bunny, only foxes within 6 tiles (Chebyshev) are alerted. Foxes farther away continue patrolling, and foxes in Wander state are never interrupted. The desired behavior is that any fox spotting the bunny triggers all foxes to converge.

## Solution

Two targeted changes in `game/fox.go`.

### 1. Remove the radius cap from `alertNearby`

`alertNearby` currently skips foxes beyond `alertRadius` (6 tiles). Remove that distance check so every other fox receives the alert. The `alertRadius` constant becomes unused and is deleted.

### 2. Simplify `Alert` to always update

`Alert` currently only acts when the receiver is in Patrol state. Replace the guard with an unconditional update:

```go
func (f *Fox) Alert(pos Vec2) {
    f.lastKnown = pos
    f.State = FoxStateChase
}
```

The `pos` passed to `Alert` is always the bunny's live position (set by `canSee` in the same frame), so it is never stale. All states — Patrol, Wander, and Chase — get the freshest known position and enter Chase.

## Unchanged Behavior

- A fox that personally sees the bunny still sets its own `lastKnown` and transitions to Chase before alerting others.
- The Wander → Patrol fallback (timer expiry) is still intact; it just won't fire if the fox is re-alerted first.
- Bush concealment (`bunny.Hidden`) still prevents detection entirely when active.

## Tests to Add

1. Two foxes separated by more than 6 tiles; one spots the bunny; verify the far fox enters `FoxStateChase` with the correct `lastKnown`.
2. Fox in `FoxStateWander` receives `Alert`; verify it enters `FoxStateChase` with updated `lastKnown`.
