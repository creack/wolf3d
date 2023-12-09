package main

import (
	_ "embed"
	"errors"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"go.creack.net/wolf3d/math2"
)

const (
	blockSize = 20

	screenWidth  = 1920
	screenHeight = 1080
	padding      = 20
)

//nolint:gochecknoglobals // Expected "readonly" presets.
var (
	wallColor       = color.RGBA{A: 255, R: 152, G: 0, B: 27}
	wallBorderColor = color.RGBA{A: 90, R: 0xf0, G: 0xf0, B: 0xf0}
	skyColor        = color.RGBA{A: 255, R: 39, G: 72, B: 255}
	groundColor     = color.RGBA{A: 255, R: 117, G: 117, B: 117}
	backgroundColor = color.RGBA{A: 255, R: 30, G: 30, B: 30}
)

//go:embed maps/map1
var mapData []byte

const (
	radFactor = math.Pi / 180
)

// Game holds the state.
type Game struct {
	px, py int
	pangle int

	world [][]int
}

// Update implements the ebiten interface.
func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return errors.New("game ended by player")
	}

	curX, curY := g.px, g.py
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.py = max(0, g.py-1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.py = min(len(g.world)-1, g.py+1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.px = max(0, g.px-1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.px = min(len(g.world[0])-1, g.px+1)
	}
	if g.world[g.py][g.px] > 0 {
		g.px, g.py = curX, curY
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.pangle += 10
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.pangle -= 10
	}

	return nil
}

// GetScale returns the scale factor to fit bounds in the given screen size.
func GetScale(screenWidth, screenHeight int, bounds image.Rectangle) int {
	width := bounds.Dx()
	height := bounds.Dy()

	return min((screenWidth-screenWidth/8.0)/width, (screenHeight-screenHeight/8.0)/height)
}

// FOV angle in degrees.
const FOV = 60

func isBetween(angle, start, end math2.Angle) bool {
	// Normalize angles.
	angle = angle.Normalize()
	start = start.Normalize()
	end = end.Normalize()

	// Check if angle is between start and end.
	if start <= end {
		return angle >= start && angle <= end
	}

	// Handle the case where the range spans 0 degrees.
	return angle >= start || angle <= end
}

// Draw implements the ebiten interface.
func (g *Game) Draw(screen *ebiten.Image) {
	img := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	img.Fill(color.RGBA{A: 0xf, R: 0x10, G: 0x10, B: 0x10})
	img.Fill(backgroundColor)

	scale := GetScale(img.Bounds().Dx(), img.Bounds().Dy(), image.Rect(0, 0, len(g.world[0]), len(g.world)))

	screenOffset := math2.
		Pt(screen.Bounds().Dx(), screen.Bounds().Dy()).
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

		if isBetween(a0, a2, a1) {
			c.B = 0
		}

		vector.StrokeLine(img, float32(playerPixelCoords.X), float32(playerPixelCoords.Y), float32(cur.X), float32(cur.Y), 1, c, false)
		ebitenutil.DebugPrintAt(img, fmt.Sprintf("%.2f\n%.2f\n%.2f", a0.Degrees(), a1.Degrees(), a2.Degrees()), int(cur.X), int(cur.Y))
	}

	// Go over each point in the world.
	for y, row := range g.world {
		for x, elem := range row {
			// If the point doesn't have a wall, keep going.
			if elem == 0 {
				continue
			}

			// Cast a ray from the player to x/y.
			cast(x, y)

			// Draw a square from the scaled x/y of size scale.
			origin := math2.Scale(math2.Pt(x, y), scale).Add(screenOffset)
			// vector.DrawFilledRect(img, float32(p.X), float32(p.Y), float32(scale), float32(scale), wallColor, false)
			vector.StrokeRect(img, float32(origin.X), float32(origin.Y), float32(scale), float32(scale), 1, wallBorderColor, false)
		}
	}

	// Outer border.
	vector.StrokeRect(img, 1, 1, float32(img.Bounds().Dx())-2, float32(img.Bounds().Dy()-2), 1, color.RGBA{A: 90, R: 0xf0, G: 0xf0, B: 0xf0}, false)

	// Outer inner.
	vector.StrokeRect(img,
		float32(screenOffset.X), float32(screenOffset.Y),
		float32(len(g.world[0])*scale), float32(len(g.world)*scale),
		1, color.RGBA{A: 40, R: 0xf0, G: 0xf0, B: 0xf0}, false)

	// Draw player as a rect.
	playerRectOrigin := playerPixelCoords.Sub(math2.Pt(2, 2))
	vector.StrokeRect(img, float32(playerRectOrigin.X), float32(playerRectOrigin.Y), 5, 5, 1, color.RGBA{255, 100, 100, 255}, false)

	// Draw the FOV lines.
	fovLine1 := math2.CoordinatesFromAngleDist(playerPixelCoords, playerPixelCoords, math2.NewDegAngle(g.pangle+FOV/2), 1000)
	vector.StrokeLine(img, float32(playerPixelCoords.X), float32(playerPixelCoords.Y), float32(fovLine1.X), float32(fovLine1.Y), 1, color.RGBA{A: 0xff, R: 0xff, G: 0x0}, false)
	fovLine2 := math2.CoordinatesFromAngleDist(playerPixelCoords, playerPixelCoords, math2.NewDegAngle(g.pangle-FOV/2), 1000)
	vector.StrokeLine(img, float32(playerPixelCoords.X), float32(playerPixelCoords.Y), float32(fovLine2.X), float32(fovLine2.Y), 1, color.RGBA{A: 0xff, B: 0xff, G: 0x0}, false)

	// Draw the image to the screen.
	op := &ebiten.DrawImageOptions{
		Blend: ebiten.BlendCopy,
	}
	// op.GeoM.Translate(10, 10)
	screen.DrawImage(img, op)

	msg := fmt.Sprintf("%d/%d\n%d/%d\n%d\n", len(g.world[0])*int(scale), len(g.world)*int(scale), g.px, g.py, g.pangle)
	ebitenutil.DebugPrint(screen, msg)

	ebitenutil.DebugPrintAt(screen, "WASD: move", 160, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()), 51, 51)
}

// Layout implements the ebiten interface.
func (g *Game) Layout(_, _ int) (w, h int) { return screenWidth, screenHeight }

func main() {
	g := &Game{
		px: 1,
		py: 1,
	}
	m, err := parseMap(mapData)
	if err != nil {
		panic(err)
	}
	for _, row := range m {
		line := make([]int, 0, len(row))
		for _, elem := range row {
			if elem.isWall {
				line = append(line, 1)
			} else {
				line = append(line, 0)
			}
		}
		g.world = append(g.world, line)
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Ray casting and shadows (Ebitengine Demo)")
	if runtime.GOOS != "js" {
		ebiten.SetFullscreen(true)
	}
	println("Starting")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
