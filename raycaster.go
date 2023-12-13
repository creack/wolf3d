package main

import (
	"image"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"go.creack.net/wolf3d/math2"
)

const texSize = 64

// Game holds the state.
type Game struct {
	world [][]MapPoint

	width, height int

	// NOTE: FOV is the ration of dir/plane vectors.
	dir   math2.Point // Direction vector.
	plane math2.Point // Camera plane vector.

	pos math2.Point // Current player position.

	last time.Time // Time when last frame was rendered. Used to scale movements.

	mapMod   int // -1: hidden, 0: minimap, 1: fullmap.
	showRays bool

	// Preloaded/cache data.
	textures, sideTextures           *image.RGBA
	texturesCache, sideTexturesCache [texSize][texSize * 8][3]byte
	triangleImg                      *ebiten.Image
}

// Implements the DDA algoright (Digital Differential Analysis).
func (g *Game) frame() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, g.width, g.height))
	// NOTE: Perf gain by using a buffer variable vs using img.Pix directly.
	buffer := img.Pix

	// img := image.NewRGBA(image.Rect(0, 0, g.width, g.height))
	// Go over each point along the X axis and cast a ray between the play and that point.
	for x := 0; x < g.width; x++ {
		// cameraX is the x-coordinate on the camera plane that
		// the current x-coordinate of the screen represents.
		// Done this way so that:
		//   - rightmost side gets coordinate 1
		//   - center         gets coordinate 0
		//   - leftmost  side gets coordinate -1
		cameraX := 2*float64(x)/float64(g.width) - 1 // X-coordinate in camera space.

		// The player position is a float, cast down to int to get the actual world case.
		worldPt := image.Pt(int(g.pos.X), int(g.pos.Y))

		// Run the DDA algo to cast a ray and get the distance
		// to the nearest wall as well as if we touch it from the X or Y side.
		dda := newDDA(cameraX, g.pos, g.dir, g.plane, worldPt)
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

		// x coordinate on the texture.
		texX := int(wallX * texSize)
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

			texs := &g.texturesCache
			if dda.side {
				texs = &g.sideTexturesCache
			}
			// Manually inline for perf gain (~5fps).
			off := (y*g.width + x) * 4
			buffer[off] = texs[texY][texNum*texSize+texX][0]
			buffer[off+1] = texs[texY][texNum*texSize+texX][1]
			buffer[off+2] = texs[texY][texNum*texSize+texX][2]
		}

		g.drawBackground(img, dda, x, wallX, drawEnd)
	}

	return img
}

func (g *Game) drawBackground(img *image.RGBA, dda *DDA, x int, wallX float64, drawEnd int) {
	// NOTE: 10fps gain by creating a buffer variable vs using img.Pix directly.
	buffer := img.Pix
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
	for y := drawEnd + 1; y < g.height; y++ {
		currentDist := float64(g.height) / (2.0*float64(y) - float64(g.height))

		weight := (currentDist - distPlayer) / (distWall - distPlayer)

		currentFloor := math2.Pt(
			weight*floorWall.X+(1.0-weight)*g.pos.X,
			weight*floorWall.Y+(1.0-weight)*g.pos.Y,
		)

		fx := int(currentFloor.X*float64(texSize)) % texSize
		fy := int(currentFloor.Y*float64(texSize)) % texSize
		fx2 := fx + (4 * texSize)

		// NOTE: 20fps gain by manually inlining.
		off := (y*g.width + x) * 4
		buffer[off] = g.texturesCache[fy][fx][0]
		buffer[off+1] = g.texturesCache[fy][fx][1]
		buffer[off+2] = g.texturesCache[fy][fx][2]

		off1 := ((g.height-y)*g.width + x) * 4
		buffer[off1] = g.texturesCache[fy][fx2][0]
		buffer[off1+1] = g.texturesCache[fy][fx2][1]
		buffer[off1+2] = g.texturesCache[fy][fx2][2]
	}
}

func (g *Game) getTexNum(x, y int) int {
	return g.world[y][x].wallType
}

func (g *Game) getColor(x, y int) color.Color {
	switch g.getTexNum(x, y) {
	case 1:
		return color.RGBA{A: 255, R: 255}
	case 2:
		return color.RGBA{A: 255, G: 255}
	case 3:
		return color.RGBA{A: 255, B: 255}
	case 4:
		return color.RGBA{A: 255, R: 180, G: 180, B: 180}
	case 5:
		return color.RGBA{A: 255, R: 0, G: 255, B: 255}
	case 6:
		return color.RGBA{A: 255, R: 255, G: 0, B: 255}
	case 7:
		return color.RGBA{A: 255, R: 255, G: 255, B: 0}
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
