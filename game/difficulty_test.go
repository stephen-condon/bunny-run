package game

import (
	"testing"
	"time"
)

func TestDifficultyInitial(t *testing.T) {
	d := NewDifficulty()
	if d.ScrollSpeed != baseScrollSpeed {
		t.Errorf("initial scroll speed: want %v got %v", baseScrollSpeed, d.ScrollSpeed)
	}
	if d.ObstacleDensity != baseObstacleDensity {
		t.Errorf("initial density: want %v got %v", baseObstacleDensity, d.ObstacleDensity)
	}
	if d.FoxInterval != baseFoxInterval {
		t.Errorf("initial fox interval: want %v got %v", baseFoxInterval, d.FoxInterval)
	}
	if d.Level() != 0 {
		t.Errorf("initial level: want 0 got %d", d.Level())
	}
}

func TestDifficultyNoChangeBeforeInterval(t *testing.T) {
	d := NewDifficulty()
	changed := d.Update(14 * time.Second)
	if changed {
		t.Error("difficulty should not change before 15s")
	}
	if d.Level() != 0 {
		t.Errorf("level should still be 0, got %d", d.Level())
	}
}

func TestDifficultyLevel1At15s(t *testing.T) {
	d := NewDifficulty()
	changed := d.Update(15 * time.Second)
	if !changed {
		t.Error("difficulty should change at 15s")
	}
	if d.Level() != 1 {
		t.Errorf("level should be 1, got %d", d.Level())
	}
	if d.ScrollSpeed != baseScrollSpeed+scrollSpeedIncrement {
		t.Errorf("scroll speed at level 1: want %v got %v", baseScrollSpeed+scrollSpeedIncrement, d.ScrollSpeed)
	}
}

func TestDifficultyLevel2At30s(t *testing.T) {
	d := NewDifficulty()
	d.Update(30 * time.Second)
	if d.Level() != 2 {
		t.Errorf("level should be 2 at 30s, got %d", d.Level())
	}
}

func TestDifficultyCapsScrollSpeed(t *testing.T) {
	d := NewDifficulty()
	d.Update(10000 * time.Second)
	if d.ScrollSpeed > maxScrollSpeed {
		t.Errorf("scroll speed exceeded cap: %v > %v", d.ScrollSpeed, maxScrollSpeed)
	}
}

func TestDifficultyCapsObstacleDensity(t *testing.T) {
	d := NewDifficulty()
	d.Update(10000 * time.Second)
	if d.ObstacleDensity > maxObstacleDensity {
		t.Errorf("density exceeded cap: %v > %v", d.ObstacleDensity, maxObstacleDensity)
	}
}

func TestDifficultyCapsfoxInterval(t *testing.T) {
	d := NewDifficulty()
	d.Update(10000 * time.Second)
	if d.FoxInterval < minFoxInterval {
		t.Errorf("fox interval below min: %v < %v", d.FoxInterval, minFoxInterval)
	}
}

func TestDifficultyIdempotentUpdate(t *testing.T) {
	d := NewDifficulty()
	d.Update(15 * time.Second)
	speed1 := d.ScrollSpeed
	changed := d.Update(15 * time.Second)
	if changed {
		t.Error("second update at same elapsed time should not change level")
	}
	if d.ScrollSpeed != speed1 {
		t.Error("scroll speed should not change on idempotent update")
	}
}
