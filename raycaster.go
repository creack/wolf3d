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

// getDeltaDist returns the distance the ray has to travel to go
// from one x-side or ine y-side to the next.
//
// For delaDist.X:
//   - the x side is `1` because we go from one case to the next.
//   - the y side is `rayDir.Y / rayDir.X` because it is exactly
//     the amount of units the ray goes in the Y-direction when
//     taking 1 step in the X-direction, i.e. it is the slope of the ray.
//     Slope formula: (y2-y1)/(x2-x1).
//     We set the first point to be the origin (0,0), so we have:
//     slope = (rayDir.Y-0)/(rayDir.X-0) = rayDir.Y/rayDir.X
//
// For deltaDist.Y, same idea but the slope is X/Y.
//
// Formula:
//
//	xSide1, ySide1 = 1, slope
//	xSide2, ySide2 = slope, 1
//	deltaDist.X = math.Hypot(xSide1, ySide1) = math.Hypot(1, slopeX) = sqrt(1*1 + ((rayDir.Y * rayDir.Y) / (rayDirX * rayDir.X)))
//	deltaDist.Y = math.Hypot(xSide2, ySide2) = math.Hypot(slopeY, 1) = sqrt(1*1 + ((rayDir.X * rayDir.X) / (rayDirY * rayDir.Y)))
//
// Simplified formula:
//
//	|rayDir| = math.Hypot(rayDir.X, rayDir.Y) = sqrt(rayDir.X * rayDir.X + rayDir.Y * rayDir.Y)
//	deltaDist.X = abs(|rayDir| / rayDir.X)
//	deltaDist.Y = abs(|rayDir| / rayDir.Y)
//
// And for our purpose, we only consider the ratio, not the actual size.
//
//	deltaDist.X = abs(1 / rayDir.X)
//	deltaDist.Y = abs(1 / rayDir.Y)
//
// NOTE: If rayDir.X or rayDir.Y are 0, the division through zero is avoided
// by setting it to a very high value math.Inf.
//
// For reference, the actual formula code:
//
//	func getDeltaDist(rayDir math2.Point) math2.Point {
//		slopeX, slopeY := rayDir.X/rayDir.Y, rayDir.Y/rayDir.X
//		return math2.Point{math.Hypot(1, slopeY), math.Hypot(slopeX, 1)}
//	}
//
// and the general simplified version:
//
//	func getDeltaDist(rayDir math2.Point) math2.Point {
//		rayDirLen := math.Hypot(rayDir.X, rayDir.Y)
//		out := math2.Point{math.Inf(1), math.Inf(1)} // Default values in case the denominator is 0.
//		if rayDir.X != 0 {
//			out.X = math.Abs(rayDirLen / rayDir.X)
//		}
//		if rayDir.Y != 0 {
//			out.Y = math.Abs(rayDirLen / rayDir.Y)
//		}
//		return out
//	}
func getDeltaDist(rayDir math2.Point) math2.Point {
	// For our purpose we don't need the length, just the ratio.
	// Use the simplified version with a length of 1.
	const rayDirLen = 1.
	out := math2.Pt(math.Inf(1), math.Inf(1)) // Default values in case the denominator is 0.
	if rayDir.X != 0 {
		out.X = math.Abs(rayDirLen / rayDir.X)
	}
	if rayDir.Y != 0 {
		out.Y = math.Abs(rayDirLen / rayDir.Y)
	}
	return out
}

// getInitialSideDist returns the length of ray from current
// position to next x or y-side.
//
// sideDist.X and sideDist.Y are initially the distance the ray has to travel
// from its start position to the first x-side and the first y-side.
//
// Later in the code they will be incremented while steps are taken.
//
// Note that we always have:
//   - worldPt.X <= g.pos.X <= worldPt.X+1
//   - worldPt.Y <= g.pos.Y <= worldPt.Y+1
func getInitialSideDist(rayDir, deltaDist, pos math2.Point, worldPt image.Point) math2.Point {
	var sideDist math2.Point

	if rayDir.X < 0 {
		// Relative X position within the world case from the left.
		sideDist.X = pos.X - float64(worldPt.X)
	} else {
		// Relative X position within the world case from the right.
		sideDist.X = float64(worldPt.X+1) - pos.X
	}
	if rayDir.Y < 0 {
		// Relative Y position within the world case from the top.
		sideDist.Y = pos.Y - float64(worldPt.Y)
	} else {
		// Relative Y position within the world case from the bottom.
		sideDist.Y = float64(worldPt.Y+1) - pos.Y
	}

	// Multiply the relative position by deltaDist to get the side dist.
	return sideDist.Mul(deltaDist)
}

// getStep returns a vector with either +1 or -1
// to describe the direction each step should take.
func getStep(rayDir math2.Point) math2.Point {
	// What direction to step in x or y-direction (either +1 or -1).
	step := math2.Pt(1, 1)
	if rayDir.X < 0 {
		step.X = -1
	}
	if rayDir.Y < 0 {
		step.Y = -1
	}
	return step
}

// getWallDist returns the distance between the wall and the camera plane.
//
// We don't use the Euclidean distance to the point
// representing player, but instead the distance to
// the camera plane (or, the distance of the point projected
// on the camera direction to the player),
// to avoid the fisheye effect.
//
// The fisheye effect is an effect you see if you use the real distance,
// where all the walls become rounded, and can make you sick if you rotate.
func getWallDist(rayDir, step, pos math2.Point, worldPt image.Point, side bool) float64 {
	if side {
		return (float64(worldPt.Y) - pos.Y + (1-step.Y)/2) / rayDir.Y
	}
	return (float64(worldPt.X) - pos.X + (1-step.X)/2) / rayDir.X
}

// Implements the DDA algoright (Digital Differential Analysis).
func (g *Game) frame() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, g.width, g.height))

	// Go over each point along the X axis and cast a ray between the play and that point.
	for x := 0; x < g.width; x++ {
		// cameraX is the x-coordinate on the camera plane that
		// the current x-coordinate of the screen represents.
		// Done this way so that:
		//   - rightmost side gets coordinate 1
		//   - center         gets coordinate 0
		//   - leftmost  side gets coordinate -1
		cameraX := 2*float64(x)/float64(g.width) - 1 // X-coordinate in camera space.

		// The direction of the ray is the sum of
		//   - the direction vector of the camera
		//   - a part of the plane vector of the camera (g.plane scaled to cameraX).
		rayDir := g.dir.Add(g.plane.Scale(cameraX))

		// The player position is a float, cast down to int to get the actual world case.
		worldPt := image.Pt(int(g.pos.X), int(g.pos.Y))

		// Compute deltaDist, the distance to travel to go from
		// one case to the next on each axis.
		deltaDist := getDeltaDist(rayDir)

		// This represents the distance between the player and the edges
		// of the current world case.
		sideDist := getInitialSideDist(rayDir, deltaDist, g.pos, worldPt)

		// This represents which direction along the ray we travel to find a wall.
		step := getStep(rayDir)

		// Now the actual DDA starts.
		//
		// It's a loop that increments the ray with 1 square every time,
		// until a wall is hit.
		//
		// Each time, either it jumps a square in the x-direction (with step.X)
		// or a square in the y-direction (with stepY),
		// it always jumps 1 square at once.
		//
		// If the ray's direction would be the x-direction,
		// the loop will only have to jump a square in the x-direction everytime,
		// because the ray will never change its y-direction.
		// If the ray is a bit sloped to the y-direction,
		// then every so many jumps in the x-direction,
		// the ray will have to jump one square in the y-direction.
		// If the ray is exactly the y-direction,
		// it never has to jump in the x-direction, etc.
		//
		// sideDistX and sideDistY get incremented with deltaDistX with
		// every jump in their direction,
		// and mapX and mapY get incremented with stepX and stepY respectively.
		//
		// When the ray has hit a wall, the loop ends,
		// and then we'll know whether an x-side or y-side of
		// a wall was hit in the variable "side",
		// and what wall was hit with mapX and mapY.
		//
		// We won't know exactly where the wall was hit however,
		// but that's not needed in this case because we won't use textured walls for now.

		var side bool // Was a North-South or a East-West wall hit?

		for worldPt.X < len(g.world) && worldPt.Y < len(g.world[worldPt.X]) { // Sanity checks.
			if g.world[worldPt.X][worldPt.Y].wallType != 0 {
				break
			}
			if sideDist.X < sideDist.Y {
				sideDist.X += deltaDist.X
				worldPt.X += int(step.X)
				side = false
			} else {
				sideDist.Y += deltaDist.Y
				worldPt.Y += int(step.Y)
				side = true
			}
		}

		// After the DDA is done, we have to calculate the distance
		// of the ray to the wall, so that we can calculate how high
		// the wall has to be drawn after this.
		perpWallDist := getWallDist(rayDir, step, g.pos, worldPt, side)

		// Calculate height of line to draw on screen.
		lineHeight := max(1, int(float64(g.height)/perpWallDist))

		// Calculate lowest and highest pixel to fill in current stripe.
		//
		// The center of the wall should be at the center of the screen,
		// and if these points lie outside the screen, they're capped to 0 or g.height-1.
		//
		// The y center of the screen is g.height/2. Start from there -1/2 length to there +1/2 length.
		drawStart, drawEnd := max(0, g.height/2-lineHeight/2), min(g.height-1, g.height/2+lineHeight/2)

		// Draw the line.
		c := g.getColor(worldPt.X, worldPt.Y)
		for y := drawStart; y < drawEnd+1; y++ {
			if side {
				img.Set(x, y, dimColor(c))
			} else {
				img.Set(x, y, c)
			}
		}
	}

	return img
}
