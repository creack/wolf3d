package main

import (
	"image"
	"image/color"
	"math"
	"time"

	"go.creack.net/wolf3d/math2"
)

// Game holds the state.
type Game struct {
	world    [][]MapPoint
	textures *image.RGBA

	width, height int

	// NOTE: FOV is the ration of dir/plane vectors.
	dir   math2.Point // Direction vector.
	plane math2.Point // Camera plane vector.

	pos math2.Point // Current player position.

	last time.Time // Time when last frame was rendered. Used to scale movements.
}

// Implements the DDA algoright (Digital Differential Analysis).
func (g *Game) frame() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, g.width, g.height))

	// Go over each point along the X axis and cast a ray between the play and that point.
	for x := 0; x < g.width; x++ {
		// Run the DDA algo to cast a ray and get the distance
		// to the nearest wall as well as if we touch it from the X or Y side.
		dda := newDDA(x, g.width, g.pos, g.dir, g.plane)
		dda.run(g.world, g.pos)

		// Calculate height of line to draw on screen.
		lineHeight := max(1, int(float64(g.height)/dda.perpWallDist))

		// Calculate lowest and highest pixel to fill in current stripe.
		//
		// The center of the wall should be at the center of the screen,
		// and if these points lie outside the screen, they're capped to 0 or g.height-1.
		//
		// The y center of the screen is g.height/2. Start from there -1/2 length to there +1/2 length.
		drawStart, drawEnd := max(0, g.height/2-lineHeight/2), min(g.height-1, g.height/2+lineHeight/2)

		// The value wallX represents the exact value where the
		// wall was hit, not just the integer coordinates of the wall.
		// This is required to know which x-coordinate of the texture
		// we have to use.
		//
		// This is calculated by first calculating the exact
		// x or y coordinate in the world, and then subtracting
		// the integer value of the wall off it.
		//
		// Note that even if it's called wallX, it's actually an
		// y-coordinate of the wall if side==1, but it's always
		// the x-coordinate of the texture.
		var wallX float64 // Where exactly the wall was hit.
		if !dda.side {
			wallX = g.pos.Y + dda.perpWallDist*dda.rayDir.Y
		} else {
			wallX = g.pos.X + dda.perpWallDist*dda.rayDir.X
		}
		wallX -= math.Floor(wallX)

		const texSize = 64.
		const texWidth, texHeight = texSize, texSize

		// x coordinate on the texture.
		texX := int(wallX * texWidth)
		if !dda.side && dda.rayDir.X > 0 {
			texX = texSize - texX - 1
		}
		if dda.side && dda.rayDir.Y < 0 {
			texX = texSize - texX - 1
		}

		texNum := g.getTexNum(dda.worldPt.X, dda.worldPt.Y)
		for y := drawStart; y < drawEnd; y++ {
			d := y - (g.height/2 - lineHeight/2)
			texY := (d * texSize) / lineHeight

			c := g.textures.RGBAAt(
				texNum*texWidth+texX,
				texY,
			)

			if dda.side {
				c.R /= 2
				c.G /= 2
				c.B /= 2
			}

			img.Set(x, y, c)
		}

		g.drawBackground(img, dda, wallX, drawEnd)
	}

	return img
}

func (g *Game) drawBackground(img *image.RGBA, dda *DDA, wallX float64, drawEnd int) {
	var floorWall math2.Point

	switch {
	case !dda.side && dda.rayDir.X > 0:
		floorWall.X = float64(dda.worldPt.X)
		floorWall.Y = float64(dda.worldPt.Y) + wallX
	case !dda.side && dda.rayDir.X < 0:
		floorWall.X = float64(dda.worldPt.X) + 1.0
		floorWall.Y = float64(dda.worldPt.Y) + wallX
	case dda.side && dda.rayDir.Y > 0:
		floorWall.X = float64(dda.worldPt.X) + wallX
		floorWall.Y = float64(dda.worldPt.Y)
	case dda.side && dda.rayDir.Y < 0:
		floorWall.X = float64(dda.worldPt.X) + wallX
		floorWall.Y = float64(dda.worldPt.Y) + 1.0
	}

	distWall, distPlayer := dda.perpWallDist, 0.0
	const texSize = 64

	for y := drawEnd + 1; y < g.height; y++ {
		currentDist := float64(g.height) / (2.0*float64(y) - float64(g.height))

		weight := (currentDist - distPlayer) / (distWall - distPlayer)

		currentFloor := math2.Pt(
			weight*floorWall.X+(1.0-weight)*g.pos.X,
			weight*floorWall.Y+(1.0-weight)*g.pos.Y,
		)

		fx := int(currentFloor.X*float64(texSize)) % texSize
		fy := int(currentFloor.Y*float64(texSize)) % texSize

		img.Set(dda.x, y, g.textures.At(fx, fy))

		img.Set(dda.x, g.height-y-1, g.textures.At(fx+(4*texSize), fy))
		img.Set(dda.x, g.height-y-0, g.textures.At(fx+(4*texSize), fy))
	}
}

func (g *Game) getTexNum(x, y int) int {
	return g.world[x][y].wallType
}

func (g *Game) getColor(x, y int) color.Color {
	switch g.getTexNum(x, y) {
	case 1:
		return color.RGBA{R: 255}
	case 2:
		return color.RGBA{G: 255}
	case 3:
		return color.RGBA{B: 255}
	case 4:
		return color.White
	case 5:
		return color.RGBA{R: 0, G: 255, B: 255}
	case 6:
		return color.RGBA{R: 255, G: 0, B: 255}
	case 7:
		return color.RGBA{R: 255, G: 255, B: 0}
	default:
		return color.Black
	}
}

func dimColor(in color.Color) color.Color {
	r, g, b, a := in.RGBA()
	return color.RGBA64{
		A: uint16(a),
		R: uint16(r / 2),
		G: uint16(g / 2),
		B: uint16(b / 2),
	}
}
