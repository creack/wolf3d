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

	"go.creack.net/fdf/math3"
	"go.creack.net/fdf/projection"
	"go.creack.net/wolf3d/math2"
)

const (
	blockSize = 20

	screenWidth  = 1024
	screenHeight = 768
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

func (g *Game) Main(width, height int) *ebiten.Image {
	worldSize := image.Rect(0, 0, len(g.world[0]), len(g.world))
	// Scale the world so it fits in the given bounds.
	scale := GetScale(width, height, worldSize)
	// // Clamp the smaller dimension.
	// width = worldSize.Dx() * scale
	// height = worldSize.Dy() * scale
	img := ebiten.NewImage(width, height)
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

	distanceFromProjectionPlan := (float64(width) / 2) / math.Tan(math2.NewDegAngle(FOV/2).Radians())
	_ = distanceFromProjectionPlan
	blockSize := float64(width) * 0.2
	_ = blockSize

	proj := projection.NewIsomorphic(1)
	_ = proj // vector.DrawFilledRect(img, float32(origin.X), float32(origin.Y), float32(scale), float32(scale), wallColor, false)
	// proj.SetOffset(math3.Vec{X: screenOffset.X, Y: screenOffset.Y})
	// proj.SetAngle(math3.Vec{})

	// Go over each point in the world.
	for y, row := range g.world {
		for x, elem := range row {
			origin := math2.Scale(math2.Pt(x, y), scale).Add(screenOffset)
			_ = origin
			cur := centerPixelCoords(x, y)

			// proj.SetAngle(math3.Vec{})
			// proj.SetOffset(math3.Vec{})
			// proj.SetScale(1)
			// If the point doesn't have a wall, keep going.
			if elem == 0 {
				a0 := math2.GetAngle(
					playerPixelCoords,
					playerPixelCoords,
					cur,
				)

				a1 := math2.NewDegAngle(float64(g.pangle) + FOV/2.)
				a2 := math2.NewDegAngle(float64(g.pangle) - FOV/2.)

				c := color.RGBA{A: 0xff, B: 0xff, G: 0xff}
				_ = c
				if !a0.Between(a2, a1) {
					continue
				}

				// vector.DrawFilledRect(img, float32(origin.X), float32(origin.Y), float32(scale), float32(scale), wallBorderColor, true)
				for i := 0; i < scale; i++ {
					p := proj.Project(math3.Vec{X: float64(i + int(origin.X)), Y: origin.Y, Z: float64(1)})

					// p.X, p.Y = origin.X, origin.Y
					vector.StrokeLine(img,
						float32(p.X), float32(p.Y),
						float32(p.X), float32(p.Y)+float32(p.Z),
						1, color.White, false)
				}

				continue
			}

			a0 := math2.GetAngle(
				playerPixelCoords,
				playerPixelCoords,
				cur,
			)

			a1 := math2.NewDegAngle(float64(g.pangle) + FOV/2.)
			a2 := math2.NewDegAngle(float64(g.pangle) - FOV/2.)

			c := color.RGBA{A: 0xff, B: 0xff, G: 0xff}
			_ = c
			if !a0.Between(a2, a1) {
				continue
			}

			A := math2.Point{}
			A.Y = (playerPixelCoords.Y/blockSize)*blockSize + blockSize
			A.X = playerPixelCoords.X + (playerPixelCoords.Y-A.Y)/math.Tan(a0.Radians())

			ya := scale
			// TODO: If cast is down, ya = -scale.
			xa := 64. / math.Tan(a0.Radians())
			_, _ = ya, xa

			// vector.StrokeRect(img, float32(origin.X), float32(origin.Y), float32(scale), float32(scale), 1, wallBorderColor, false)

			pA := math.Abs(playerPixelCoords.X-origin.X) / math.Cos(float64(math2.NewDegAngle(FOV)))
			wallHeight := (blockSize / pA) * distanceFromProjectionPlan

			for i := 0; i < scale; i++ {
				p := proj.Project(math3.Vec{X: float64(i + x*scale), Y: float64(y * scale), Z: wallHeight})
				vector.StrokeLine(img,
					float32(p.X), float32(p.Y),
					float32(p.X), float32(p.Y)+float32(wallHeight),
					1, wallColor, false)
			}

			// vector.StrokeLine(img, float32(origin.X), float32(origin.Y), float32(origin.X), float32(origin.Y)+float32(wallHeight), 1, wallBorderColor, false)
			// for i := 0; i < int(wallHeight); i++ {
			// 	for j := 0; j < scale; j++ {
			// 		x, y := origin.X+float64(j), origin.Y+float64(i)
			// 		p := proj.Project(math3.Vec{X: x, Y: y, Z: wallHeight})
			// 		// fmt.Println(i, j, x, y, p)
			// 		img.Set(int(p.X), int(p.Y), wallColor)
			// 		// img.Set(int(x), int(y), wallColor)
			// 	}
			// }
			// vector.DrawFilledRect(img, float32(origin.X), float32(origin.Y), float32(scale), float32(wallHeight), wallColor, false)

			ebitenutil.DebugPrintAt(img, fmt.Sprintf("A: %v, H: %v", A, wallHeight), int(origin.X), int(origin.Y))
			// alpha = atan(blocksize/xa)

			//-mur trouver de façon horizontale : A(x,y)
			//-mur trouver de façon verticale : B(x1,By1)
			//-PA = abs (PlayerPtX-x) / cos(AngleDeVision)
			//-PB = abs (PlayerPtX-x1) / cos(AngleDeVision)
			// distanceIncorrecte=min(PA,PB)
			// correcteDistance = distanceIncorrecte * cos(Beta)
			//   Beta = +30 pour les rayons à gauche du rayon du milieu
			//   Beta = -30 pour les rayons à droite du rayon du milieu

			// PA := math.Abs(playerPixelCoords.X - x/cos(FOV))
			// PB := math.Abs(playerPixelCoords.X - x1/cos(FOV))
			// wallHeight := (blockSize / min(PA, PB)) * distanceFromProjectionPlan
			// _ = wallHeight

		}
	}

	// Bx = int (playerPixelCoords.X/64) * (64) + 64.
	// // Si les rayons ont été confrontés à gauche
	// Bx = int (PlayerPtX/64) * (64) - 1.
	// Et on a :
	// By = playplayerPixelCoords.Y + (plaplayerPixelCoords.X-x) * math.Tan (FOV);

	return img
}

// Draw implements the ebiten interface.
func (g *Game) Draw(screen *ebiten.Image) {
	width, height := screen.Bounds().Dx(), screen.Bounds().Dy()
	img := ebiten.NewImage(width, height)

	// Draw the background for floor/sky.
	// vector.DrawFilledRect(img, 0, 0, float32(width), float32(height)/2, skyColor, false)
	// vector.DrawFilledRect(img, 0, float32(height)/2, float32(width), float32(height)/2, groundColor, false)

	// Main.
	mainImg := g.Main(width, height)
	mainOp := &ebiten.DrawImageOptions{}
	img.DrawImage(mainImg, mainOp)

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
