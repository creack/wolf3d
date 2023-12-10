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

	minimapState minimapState

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

	return min(screenWidth/width, screenHeight/height)
}

// FOV angle in degrees.
const FOV = 60

// Draw implements the ebiten interface.
func (g *Game) Draw(screen *ebiten.Image) {
	width, height := screen.Bounds().Dx(), screen.Bounds().Dy()
	img := ebiten.NewImage(width, height)

	// Draw the background for floor/sky.
	vector.DrawFilledRect(img, 0, 0, float32(width), float32(height)/2, skyColor, false)
	vector.DrawFilledRect(img, 0, float32(height)/2, float32(width), float32(height)/2, groundColor, false)

	// Generate the minimap.
	minimapImg := g.Minimap(int(float64(width)*0.2), int(float64(height)*0.2))
	// Draw the minimap.
	minimapOp := &ebiten.DrawImageOptions{}
	minimapOp.GeoM.Translate(float64(width-minimapImg.Bounds().Dx()), 0)
	img.DrawImage(minimapImg, minimapOp)

	// Draw the image to the screen.
	screenOp := &ebiten.DrawImageOptions{}
	screen.DrawImage(img, screenOp)

	ebitenutil.DebugPrintAt(screen, "WASD: move", 160, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()), 51, 51)
}

// Layout implements the ebiten interface.
func (g *Game) Layout(_, _ int) (w, h int) { return screenWidth, screenHeight }

func main() {
	g := &Game{
		px: 1,
		py: 1,
		minimapState: minimapState{
			showFOVLines: true,
			showRays:     true,
		},
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
