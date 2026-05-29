package game

import "time"

const (
	difficultyInterval   = 15 * time.Second
	baseScrollSpeed      = 60.0  // px/s
	maxScrollSpeed       = 200.0 // px/s
	scrollSpeedIncrement = 5.0

	baseObstacleDensity      = 0.10 // 10% of tiles
	maxObstacleDensity       = 0.40
	obstacleDensityIncrement = 0.05

	baseFoxInterval      = 10 * time.Second // one fox every N seconds
	minFoxInterval       = 5 * time.Second
	foxIntervalDecrement = 2 * time.Second
)

// Difficulty holds the current game difficulty parameters.
type Difficulty struct {
	ScrollSpeed     float64
	ObstacleDensity float64
	FoxInterval     time.Duration
	level           int
}

// NewDifficulty returns difficulty at level 0.
func NewDifficulty() Difficulty {
	return Difficulty{
		ScrollSpeed:     baseScrollSpeed,
		ObstacleDensity: baseObstacleDensity,
		FoxInterval:     baseFoxInterval,
	}
}

// Level returns the current difficulty level.
func (d *Difficulty) Level() int { return d.level }

// Update advances difficulty based on elapsed time.
// It returns true if the level increased this call.
func (d *Difficulty) Update(elapsed time.Duration) bool {
	newLevel := int(elapsed / difficultyInterval)
	if newLevel <= d.level {
		return false
	}
	steps := newLevel - d.level
	d.level = newLevel
	for i := 0; i < steps; i++ {
		d.ScrollSpeed += scrollSpeedIncrement
		if d.ScrollSpeed > maxScrollSpeed {
			d.ScrollSpeed = maxScrollSpeed
		}
		d.ObstacleDensity += obstacleDensityIncrement
		if d.ObstacleDensity > maxObstacleDensity {
			d.ObstacleDensity = maxObstacleDensity
		}
		d.FoxInterval -= foxIntervalDecrement
		if d.FoxInterval < minFoxInterval {
			d.FoxInterval = minFoxInterval
		}
	}
	return true
}
