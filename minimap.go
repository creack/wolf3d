package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"go.creack.net/wolf3d/math2"
)

type minimapState struct {
	showFOVLines bool
	showRays     bool
	showGrid     bool
	showBorders  bool

	hidePlayer bool
}

// Minimap draws the world in 2D.
func (g *Game) Minimap(width, height int) *ebiten.Image {
	worldSize := image.Rect(0, 0, len(g.world[0]), len(g.world))
	// Scale the world so it fits in the given bounds.
	scale := GetScale(width, height, worldSize)
	// Clamp the smaller dimension.
	width = worldSize.Dx() * scale
	height = worldSize.Dy() * scale

	img := ebiten.NewImage(width, height)
	img.Fill(color.RGBA{A: 0xf, R: 0x10, G: 0x10, B: 0x10})
	img.Fill(backgroundColor)

	screenOffset := math2.
		Pt(img.Bounds().Dx(), img.Bounds().Dy()).
		Sub(math2.Scale(math2.Pt(len(g.world[0]), len(g.world)), scale)).
		Scale(0.5)

	centerPixelCoords := func(x, y int) math2.Point {
		return math2.
			Scale(math2.Pt(x, y), scale).   // Apply the scale.
			Add(screenOffset).              // Apply the screen offset.
			Add(math2.Pt(scale/2, scale/2)) // Translate to be in the middle of the case.
	}

	playerPixelCoords := centerPixelCoords(g.px, g.py)
	cast := func(x, y int) {
		cur := centerPixelCoords(x, y)

		a0 := math2.GetAngle(
			playerPixelCoords,
			playerPixelCoords,
			cur,
		)

		a1 := math2.NewDegAngle(float64(g.pangle) + FOV/2.)
		a2 := math2.NewDegAngle(float64(g.pangle) - FOV/2.)

		c := color.RGBA{A: 0xff, B: 0xff, G: 0xff}

		if a0.Between(a2, a1) {
			c.B = 0
		}

		vector.StrokeLine(img, float32(playerPixelCoords.X), float32(playerPixelCoords.Y), float32(cur.X), float32(cur.Y), 1, c, false)
		// ebitenutil.DebugPrintAt(img, fmt.Sprintf("%.2f\n%.2f\n%.2f", a0.Degrees(), a1.Degrees(), a2.Degrees()), int(cur.X), int(cur.Y))
	}

	// Go over each point in the world.
	for y, row := range g.world {
		for x, elem := range row {
			origin := math2.Scale(math2.Pt(x, y), scale).Add(screenOffset)

			if g.minimapState.showGrid {
				vector.StrokeRect(img, float32(origin.X), float32(origin.Y), float32(scale), float32(scale), 1, wallBorderColor, false)
			}

			// If the point doesn't have a wall, keep going.
			if elem == 0 {
				continue
			}

			// Cast a ray from the player to x/y.
			if g.minimapState.showRays {
				cast(x, y)
			}

			// Draw a square from the scaled x/y of size scale.
			vector.DrawFilledRect(img, float32(origin.X), float32(origin.Y), float32(scale), float32(scale), wallColor, false)
		}
	}

	if g.minimapState.showBorders {
		// Outer border.
		vector.StrokeRect(img, 1, 1, float32(img.Bounds().Dx())-2, float32(img.Bounds().Dy()-2), 1, color.RGBA{A: 90, R: 0xf0, G: 0xf0, B: 0xf0}, false)

		// Inner border.
		vector.StrokeRect(img,
			float32(screenOffset.X), float32(screenOffset.Y),
			float32(len(g.world[0])*scale), float32(len(g.world)*scale),
			1, color.RGBA{A: 40, R: 0xf0, G: 0xf0, B: 0xf0}, false)
	}

	// Draw player as a rect.
	if !g.minimapState.hidePlayer {
		playerRectOrigin := playerPixelCoords.Sub(math2.Pt(2, 2))
		vector.StrokeRect(img, float32(playerRectOrigin.X), float32(playerRectOrigin.Y), 5, 5, 1, color.RGBA{255, 100, 100, 255}, false)
	}

	// Draw the FOV lines.
	if g.minimapState.showFOVLines {
		fovLine1 := math2.CoordinatesFromAngleDist(playerPixelCoords, playerPixelCoords, math2.NewDegAngle(g.pangle+FOV/2), 1000)
		vector.StrokeLine(img,
			float32(playerPixelCoords.X), float32(playerPixelCoords.Y),
			float32(fovLine1.X), float32(fovLine1.Y),
			1, color.RGBA{A: 0xff, R: 0xff, G: 0x0}, false)
		fovLine2 := math2.CoordinatesFromAngleDist(playerPixelCoords, playerPixelCoords, math2.NewDegAngle(g.pangle-FOV/2), 1000)
		vector.StrokeLine(img,
			float32(playerPixelCoords.X), float32(playerPixelCoords.Y),
			float32(fovLine2.X), float32(fovLine2.Y),
			1, color.RGBA{A: 0xff, B: 0xff, G: 0x0}, false)
	}

	return img
}
