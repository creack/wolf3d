package main

import (
	"image"
	"image/color"
	"image/draw"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"go.creack.net/wolf3d/math2"
)

func rayVertices(x1, y1, x2, y2, x3, y3 float64) []ebiten.Vertex {
	return []ebiten.Vertex{
		{DstX: float32(x1), DstY: float32(y1), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(x2), DstY: float32(y2), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(x3), DstY: float32(y3), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
	}
}

func (g *Game) minimap(width, height int) image.Image {
	worldWidth, worldHeight := len(g.world[0]), len(g.world)
	scale := min(width/worldWidth, height/worldHeight)
	width, height = worldWidth*scale, worldHeight*scale

	var hits [100][100][2]*DDA
	shadowImage := ebiten.NewImage(width, height)
	shadowImage.Fill(color.Black)
	var img draw.Image = image.NewRGBA(image.Rect(0, 0, width, height))
	img = ebiten.NewImageFromImage(img)

	// Player position.
	spos := g.pos.Scale(float64(scale))

	// Go over each point along the X axis and cast a ray between the play and that point.
	for x := 0; x < g.width; x++ {
		cameraX := 2*float64(x)/float64(g.width) - 1 // X-coordinate in camera space.

		// Run the DDA algo to cast a ray and get the distance
		// to the nearest wall as well as if we touch it from the X or Y side.
		dda := newDDA(cameraX, g.pos, g.dir, g.plane)
		dda.run(g.world, g.pos)

		h := hits[dda.worldPt.X][dda.worldPt.Y]
		if h[0] == nil {
			h[0] = dda
			h[1] = dda
			hits[dda.worldPt.X][dda.worldPt.Y] = h
		}
		if dda.realWallDist < h[0].realWallDist {
			h[0] = dda
		}
		if dda.realWallDist > h[1].realWallDist {
			h[1] = dda
		}
	}

	opt := &ebiten.DrawTrianglesOptions{}
	opt.Address = ebiten.AddressRepeat
	opt.Blend = ebiten.BlendSourceOut

	rays := make([]*DDA, 0, len(hits)*2)
	for y := 0; y < worldHeight; y++ {
		for x := 0; x < worldWidth; x++ {
			elem := hits[x][y]
			if elem[0] != nil {
				rays = append(rays, elem[0], elem[1])
			}
		}
	}

	getAngle := func(dda *DDA) math2.Angle {
		return math2.GetAngle(spos, spos, dda.rayDir)
	}
	sort.Slice(rays, func(i int, j int) bool {
		return getAngle(rays[i]) < getAngle(rays[j])
	})

	for i, dda := range rays {
		if i+1 >= len(rays) {
			continue
		}
		next := rays[(i+1)%len(rays)]

		getLine := func(dda *DDA) math2.Point {
			a0 := math2.GetAngle(spos, spos, spos.Add(dda.rayDir))
			return math2.CoordinatesFromAngleDist(spos, spos, a0, (dda.realWallDist)*float64(scale))
		}
		line := getLine(dda)
		nextLine := getLine(next)

		if g.showRays {
			c := color.RGBA{A: 255, G: 255}
			img1, _ := img.(*ebiten.Image)
			vector.StrokeLine(img1, float32(spos.X), float32(spos.Y), float32(line.X), float32(line.Y), 1, c, true)
			img = img1
		}

		v := rayVertices(spos.X, spos.Y, nextLine.X, nextLine.Y, line.X, line.Y)
		shadowImage.DrawTriangles(v, []uint16{0, 1, 2}, g.triangleImg, opt)
	}

	{
		img1, _ := img.(*ebiten.Image)
		op := &ebiten.DrawImageOptions{}
		op.ColorScale.ScaleAlpha(0.7)
		img1.DrawImage(shadowImage, op)
		img = img1
	}

	g.drawMinimapWalls(img, scale, hits)
	img = g.drawMinimapPlayer(img, scale)

	return img
}

func (g *Game) drawMinimapPlayer(i draw.Image, scale int) draw.Image {
	img, _ := i.(*ebiten.Image)

	spos := g.pos.Scale(float64(scale))
	// Draw the player itself.
	vector.DrawFilledCircle(img, float32(spos.X), float32(spos.Y), float32(min(1, scale)), color.RGBA{A: 255, R: 255}, true)
	if g.showHighlight {
		// Highlight the current world coordinate.
		vector.StrokeRect(img, float32(int(g.pos.X)*scale), float32(int(g.pos.Y)*scale), float32(scale), float32(scale), 1, color.White, false)
	}

	return img
}

func (g *Game) drawMinimapWalls(i draw.Image, scale int, hits [100][100][2]*DDA) {
	worldWidth, worldHeight := len(g.world[0]), len(g.world)
	img, _ := i.(*ebiten.Image)

	for y := 0; y < worldHeight; y++ {
		for x := 0; x < worldWidth; x++ {
			if g.showMinimapGrid {
				vector.StrokeRect(img, float32(x*scale), float32(y*scale), float32(scale), float32(scale), 1, color.White, false)
			}
			c := g.getColor(x, y)
			if c == color.Black {
				continue
			}
			if g.hideInvisibleWalls && hits[x][y][0] == nil {
				continue
			}

			vector.DrawFilledRect(img, float32(x*scale), float32(y*scale), float32(scale), float32(scale), c, false)
		}
	}
}
