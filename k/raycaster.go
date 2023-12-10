package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"runtime"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed textures.png
var textureData []byte

//go:embed map4
var mapData []byte

type Point struct{ X, Y float64 }

func Pt(x, y float64) Point { return Point{x, y} }

const texSize = 64

func loadWorld() [][]int {
	//nolint:prealloc // False positive.
	var w [][]int
	for _, line := range strings.Split(string(mapData), "\n") {
		if line == "" {
			continue
		}
		var worldLine []int
		for _, elem := range strings.ReplaceAll(line, " ", "") {
			worldLine = append(worldLine, int(elem-'0'))
		}
		w = append(w, worldLine)
	}
	w2 := make([][]int, len(w[0]))
	for i := range w2 {
		w2[i] = make([]int, len(w))
	}
	for y, line := range w {
		for x, elem := range line {
			w2[x][y] = elem
		}
	}
	return w2
}

func loadTextures() *image.RGBA {
	p, err := png.Decode(bytes.NewReader(textureData))
	if err != nil {
		panic(err)
	}

	m := image.NewRGBA(p.Bounds())

	draw.Draw(m, m.Bounds(), p, image.Point{}, draw.Src)

	return m
}

func (g *Game) getTexNum(x, y int) int {
	return g.world[x][y]
}

func (g *Game) frame() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, g.width, g.height))

	for x := 0; x < g.width; x++ {
		var (
			step         image.Point
			sideDist     Point
			perpWallDist float64
			hit, side    bool

			rayPos, worldX, worldY = g.pos, int(g.pos.X), int(g.pos.Y)

			cameraX = 2*float64(x)/float64(g.width) - 1

			rayDir = Pt(
				g.dir.X+g.plane.X*cameraX,
				g.dir.Y+g.plane.Y*cameraX,
			)

			deltaDist = Pt(
				math.Sqrt(1.0+(rayDir.Y*rayDir.Y)/(rayDir.X*rayDir.X)),
				math.Sqrt(1.0+(rayDir.X*rayDir.X)/(rayDir.Y*rayDir.Y)),
			)
		)

		if rayDir.X < 0 {
			step.X = -1
			sideDist.X = (rayPos.X - float64(worldX)) * deltaDist.X
		} else {
			step.X = 1
			sideDist.X = (float64(worldX) + 1.0 - rayPos.X) * deltaDist.X
		}

		if rayDir.Y < 0 {
			step.Y = -1
			sideDist.Y = (rayPos.Y - float64(worldY)) * deltaDist.Y
		} else {
			step.Y = 1
			sideDist.Y = (float64(worldY) + 1.0 - rayPos.Y) * deltaDist.Y
		}

		for !hit {
			if sideDist.X < sideDist.Y {
				sideDist.X += deltaDist.X
				worldX += step.X
				side = false
			} else {
				sideDist.Y += deltaDist.Y
				worldY += step.Y
				side = true
			}

			if g.world[worldX][worldY] > 0 {
				hit = true
			}
		}

		var wallX float64

		if side {
			perpWallDist = (float64(worldY) - rayPos.Y + (1-float64(step.Y))/2) / rayDir.Y
			wallX = rayPos.X + perpWallDist*rayDir.X
		} else {
			perpWallDist = (float64(worldX) - rayPos.X + (1-float64(step.X))/2) / rayDir.X
			wallX = rayPos.Y + perpWallDist*rayDir.Y
		}

		if x == g.width/2 {
			g.wallDistance = perpWallDist
		}

		wallX -= math.Floor(wallX)

		// texX := int(wallX * float64(texSize))

		lineHeight := int(float64(g.height) / perpWallDist)

		if lineHeight < 1 {
			lineHeight = 1
		}

		drawStart := -lineHeight/2 + g.height/2
		if drawStart < 0 {
			drawStart = 0
		}

		drawEnd := lineHeight/2 + g.height/2
		if drawEnd >= g.height {
			drawEnd = g.height - 1
		}

		// if !side && rayDir.X > 0 {
		// 	texX = texSize - texX - 1
		// }

		// if side && rayDir.Y < 0 {
		// 	texX = texSize - texX - 1
		// }

		texNum := g.getTexNum(worldX, worldY)

		for y := drawStart; y < drawEnd+1; y++ {
			// d := y*256 - g.height*128 + lineHeight*128
			// texY := ((d * texSize) / lineHeight) / 256

			// c := g.textures.RGBAAt(
			// 	texX+texSize*(texNum),
			// 	texY%texSize,
			// )

			c := color.RGBA{A: 255, R: 255}
			if texNum == 0 {
				c = color.RGBA{A: 255, B: 255}
			}
			// if side {
			// 	c.R /= 2
			// 	c.G /= 2
			// 	c.B /= 2
			// }

			img.Set(x, y, c)
		}

		continue
		var floorWall Point

		if !side && rayDir.X > 0 {
			floorWall.X = float64(worldX)
			floorWall.Y = float64(worldY) + wallX
		} else if !side && rayDir.X < 0 {
			floorWall.X = float64(worldX) + 1.0
			floorWall.Y = float64(worldY) + wallX
		} else if side && rayDir.Y > 0 {
			floorWall.X = float64(worldX) + wallX
			floorWall.Y = float64(worldY)
		} else {
			floorWall.X = float64(worldX) + wallX
			floorWall.Y = float64(worldY) + 1.0
		}

		continue
		distWall, distPlayer := perpWallDist, 0.0

		for y := drawEnd + 1; y < g.height; y++ {
			currentDist := float64(g.height) / (2.0*float64(y) - float64(g.height))

			weight := (currentDist - distPlayer) / (distWall - distPlayer)

			currentFloor := Pt(
				weight*floorWall.X+(1.0-weight)*g.pos.X,
				weight*floorWall.Y+(1.0-weight)*g.pos.Y,
			)

			fx := int(currentFloor.X*float64(texSize)) % texSize
			fy := int(currentFloor.Y*float64(texSize)) % texSize

			img.Set(x, y, g.textures.At(fx, fy))

			img.Set(x, g.height-y-1, g.textures.At(fx+(4*texSize), fy))
			img.Set(x, g.height-y, g.textures.At(fx+(4*texSize), fy))
		}
	}

	return img
}

type Game struct {
	world    [][]int
	textures *image.RGBA

	width, height int

	dir          Point
	plane        Point
	pos          Point
	wallDistance float64
}

func (g Game) Layout(_, _ int) (w, h int) { return g.width, g.height }

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	img := ebiten.NewImageFromImage(g.frame())

	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(img, op)
}

var last = time.Now()

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return fmt.Errorf("exit")
	}

	dt := time.Since(last).Seconds()
	last = time.Now()

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

func (g *Game) moveForward(s float64) {
	if g.wallDistance > 0.3 {
		if g.world[int(g.pos.X+g.dir.X*s)][int(g.pos.Y)] == 0 {
			g.pos.X += g.dir.X * s
		}

		if g.world[int(g.pos.X)][int(g.pos.Y+g.dir.Y*s)] == 0 {
			g.pos.Y += g.dir.Y * s
		}
	}
}

func (g *Game) moveLeft(s float64) {
	if g.world[int(g.pos.X-g.plane.X*s)][int(g.pos.Y)] == 0 {
		g.pos.X -= g.plane.X * s
	}

	if g.world[int(g.pos.X)][int(g.pos.Y-g.plane.Y*s)] == 0 {
		g.pos.Y -= g.plane.Y * s
	}
}

func (g *Game) moveBackwards(s float64) {
	if g.world[int(g.pos.X-g.dir.X*s)][int(g.pos.Y)] == 0 {
		g.pos.X -= g.dir.X * s
	}

	if g.world[int(g.pos.X)][int(g.pos.Y-g.dir.Y*s)] == 0 {
		g.pos.Y -= g.dir.Y * s
	}
}

func (g *Game) moveRight(s float64) {
	if g.world[int(g.pos.X+g.plane.X*s)][int(g.pos.Y)] == 0 {
		g.pos.X += g.plane.X * s
	}

	if g.world[int(g.pos.X)][int(g.pos.Y+g.plane.Y*s)] == 0 {
		g.pos.Y += g.plane.Y * s
	}
}

func (g *Game) turnRight(s float64) {
	oldDirX := g.dir.X

	g.dir.X = g.dir.X*math.Cos(-s) - g.dir.Y*math.Sin(-s)
	g.dir.Y = oldDirX*math.Sin(-s) + g.dir.Y*math.Cos(-s)

	oldPlaneX := g.plane.X

	g.plane.X = g.plane.X*math.Cos(-s) - g.plane.Y*math.Sin(-s)
	g.plane.Y = oldPlaneX*math.Sin(-s) + g.plane.Y*math.Cos(-s)
}

func (g *Game) turnLeft(s float64) {
	oldDirX := g.dir.X

	g.dir.X = g.dir.X*math.Cos(s) - g.dir.Y*math.Sin(s)
	g.dir.Y = oldDirX*math.Sin(s) + g.dir.Y*math.Cos(s)

	oldPlaneX := g.plane.X

	g.plane.X = g.plane.X*math.Cos(s) - g.plane.Y*math.Sin(s)
	g.plane.Y = oldPlaneX*math.Sin(s) + g.plane.Y*math.Cos(s)
}

func main() {
	g := &Game{
		width:  320,
		height: 200,

		world:    loadWorld(),
		textures: loadTextures(),

		pos:          Pt(12.0, 14.5),
		dir:          Pt(-1.0, 0.0),
		plane:        Pt(0.0, 0.66),
		wallDistance: 8.0,
	}

	ebiten.SetWindowSize(g.width*2, g.height*2)
	ebiten.SetWindowTitle("Ray casting and shadows (Ebitengine Demo)")
	if runtime.GOOS != "js" {
		ebiten.SetFullscreen(true)
	}
	println("Starting")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
