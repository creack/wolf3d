package main

import (
	"fmt"
	"image/color"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Layout implements ebiten.
func (g Game) Layout(_, _ int) (w, h int) {
	return g.width, g.height
}

// Draw implements ebiten.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	// img := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	img := ebiten.NewImageFromImage(g.frame())
	// img.Clear()

	if g.mapMod != -1 {
		scale := 0.6
		if g.mapMod == 1 {
			scale = 1.0
		}

		minimapImg := ebiten.NewImageFromImage(g.minimap(int(float64(img.Bounds().Dx())*scale), int(float64(img.Bounds().Dy())*scale)))

		opMinimap := &ebiten.DrawImageOptions{}
		// opMinimap.GeoM.Scale(scale, scale)
		opMinimap.GeoM.Translate(float64(g.width)-float64(minimapImg.Bounds().Dx()), 0)
		img.DrawImage(minimapImg, opMinimap)
	}

	ebitenutil.DebugPrint(img, fmt.Sprintf("TPS: %0.2f, FPS: %0.2f\nResolution: %dx%d", ebiten.ActualTPS(), ebiten.ActualFPS(), g.width, g.height))

	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(img, op)
}

// Update implements ebiten.
func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		if runtime.GOOS != "js" {
			return fmt.Errorf("exit")
		}
	}

	dt := time.Since(g.last).Seconds()
	g.last = time.Now()

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.showRays = !g.showRays
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		switch g.mapMod {
		case -1:
			g.mapMod = 0
		case 0:
			g.mapMod = 1
		case 1:
			g.mapMod = -1
		}
	}

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
