package game

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// KeyboardInput is the real keyboard-backed InputSource.
type KeyboardInput struct{}

func (k KeyboardInput) IsUpPressed() bool {
	return ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW)
}
func (k KeyboardInput) IsDownPressed() bool {
	return ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS)
}
func (k KeyboardInput) IsLeftPressed() bool {
	return ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA)
}
func (k KeyboardInput) IsRightPressed() bool {
	return ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD)
}
func (k KeyboardInput) IsActionPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace)
}

func (k KeyboardInput) JustPressedChar() rune {
	keys := inpututil.AppendJustPressedKeys(nil)
	for _, key := range keys {
		if key >= ebiten.KeyA && key <= ebiten.KeyZ {
			return rune('A' + int(key-ebiten.KeyA))
		}
		if key >= ebiten.Key0 && key <= ebiten.Key9 {
			return rune('0' + int(key-ebiten.Key0))
		}
	}
	return 0
}

func (k KeyboardInput) IsBackspaceJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyBackspace)
}

// RealClock implements Clock using the real system time.
type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now() }
