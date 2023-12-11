package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Layout implements ebiten.
func (g Game) Layout(_, _ int) (w, h int) { return g.width, g.height }

// Draw implements ebiten.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	img := ebiten.NewImageFromImage(g.frame())

	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(img, op)
}

// Update implements ebiten.
func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return fmt.Errorf("exit")
	}

	dt := time.Since(g.last).Seconds()
	g.last = time.Now()

	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.moveForward(3.5 * dt)
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.moveLeft(3.5 * dt)
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.moveBackwards(3.5 * dt)
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.moveRight(3.5 * dt)
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.turnRight(1.2 * dt)
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.turnLeft(1.2 * dt)
	}

	return nil
}
