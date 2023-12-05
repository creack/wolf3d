// Copyright 2019 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

// import (
// 	"errors"
// 	"fmt"
// 	"image/color"
// 	_ "image/png"
// 	"log"
// 	"math"
// 	"sort"

// 	"github.com/hajimehoshi/ebiten/v2"
// 	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
// 	"github.com/hajimehoshi/ebiten/v2/inpututil"
// 	"github.com/hajimehoshi/ebiten/v2/vector"
// )

// const (
// 	screenWidth  = 240
// 	screenHeight = 240
// 	padding      = 20
// )

// // var (
// // 	// bgImage       *ebiten.Image
// // 	shadowImage   = ebiten.NewImage(screenWidth, screenHeight)
// // 	triangleImage = ebiten.NewImage(screenWidth, screenHeight)
// // )

// // func init() {
// // 	// Decode an image from the image file's byte slice.
// // 	img, _, err := image.Decode(bytes.NewReader(images.Tile_png))
// // 	if err != nil {
// // 		log.Fatal(err)
// // 	}
// // 	bgImage = ebiten.NewImageFromImage(img)
// // 	triangleImage.Fill(color.White)
// // }

// type line struct {
// 	X1, Y1, X2, Y2 float64
// }

// func (l *line) angle() float64 {
// 	return math.Atan2(l.Y2-l.Y1, l.X2-l.X1)
// }

// type object struct {
// 	walls   []line
// 	visible bool
// }

// func (o object) points() [][2]float64 {
// 	// Get one of the endpoints for all segments,
// 	// + the startpoint of the first one, for non-closed paths
// 	var points [][2]float64
// 	for _, wall := range o.walls {
// 		points = append(points, [2]float64{wall.X2, wall.Y2})
// 	}
// 	p := [2]float64{o.walls[0].X1, o.walls[0].Y1}
// 	if p[0] != points[len(points)-1][0] && p[1] != points[len(points)-1][1] {
// 		points = append(points, [2]float64{o.walls[0].X1, o.walls[0].Y1})
// 	}
// 	return points
// }

// func newRay(x, y, length, angle float64) line {
// 	return line{
// 		X1: x,
// 		Y1: y,
// 		X2: x + length*math.Cos(angle),
// 		Y2: y + length*math.Sin(angle),
// 	}
// }

// // intersection calculates the intersection of given two lines.
// func intersection(l1, l2 line) (x, y float64, intersect bool) {
// 	// https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection#Given_two_points_on_each_line
// 	denom := (l1.X1-l1.X2)*(l2.Y1-l2.Y2) - (l1.Y1-l1.Y2)*(l2.X1-l2.X2)
// 	tNum := (l1.X1-l2.X1)*(l2.Y1-l2.Y2) - (l1.Y1-l2.Y1)*(l2.X1-l2.X2)
// 	uNum := -((l1.X1-l1.X2)*(l1.Y1-l2.Y1) - (l1.Y1-l1.Y2)*(l1.X1-l2.X1))

// 	if denom == 0 {
// 		return 0, 0, false
// 	}

// 	t := tNum / denom
// 	if t > 1 || t < 0 {
// 		return 0, 0, false
// 	}

// 	u := uNum / denom
// 	if u > 1 || u < 0 {
// 		return 0, 0, false
// 	}

// 	x = l1.X1 + t*(l1.X2-l1.X1)
// 	y = l1.Y1 + t*(l1.Y2-l1.Y1)
// 	return x, y, true
// }

// type foo struct {
// 	p [2]float64
// 	o *object
// }

// // rayCasting returns a slice of line originating from point cx, cy and intersecting with objects
// func rayCasting(cx, cy float64, objects []*object) []line {
// 	const rayLength = 1000 // something large enough to reach all objects

// 	var rays []line
// 	for _, obj := range objects {
// 		for _, p := range obj.points() {
// 			l := line{cx, cy, p[0], p[1]}
// 			angle := l.angle()

// 			// if angle > 60*math.Pi/180 {
// 			// 	continue
// 			// }

// 			// Cast two rays per point so we have a triagle.
// 			for _, offset := range []float64{-0.005, 0.005} {
// 				var points []foo
// 				ray := newRay(cx, cy, rayLength, angle+offset)

// 				// Unpack all objects.
// 				for _, o := range objects {
// 					for _, wall := range o.walls {
// 						if px, py, ok := intersection(ray, wall); ok {
// 							points = append(points, foo{[2]float64{px, py}, o})
// 						}
// 					}
// 				}
// 				// Find the point closest to start of ray.
// 				min := math.Inf(1)
// 				minI := -1
// 				for i, p := range points {
// 					d2 := (cx-p.p[0])*(cx-p.p[0]) + (cy-p.p[1])*(cy-p.p[1])
// 					if d2 < min {
// 						min = d2
// 						minI = i
// 					}
// 				}
// 				minP := points[minI]
// 				minP.o.visible = true

// 				rays = append(rays, line{cx, cy, minP.p[0], minP.p[1]})
// 			}
// 		}
// 	}

// 	// Sort rays based on angle, otherwise light triangles will not come out right
// 	sort.Slice(rays, func(i int, j int) bool {
// 		return rays[i].angle() < rays[j].angle()
// 	})
// 	return rays
// }

// func (g *Game) handleMovement() {
// 	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
// 		g.px += 4
// 	}

// 	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
// 		g.py += 4
// 	}

// 	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
// 		g.px -= 4
// 	}

// 	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
// 		g.py -= 4
// 	}

// 	// +1/-1 is to stop player before it reaches the border
// 	if g.px >= screenWidth-padding {
// 		g.px = screenWidth - padding - 1
// 	}

// 	if g.px <= padding {
// 		g.px = padding + 1
// 	}

// 	if g.py >= screenHeight-padding {
// 		g.py = screenHeight - padding - 1
// 	}

// 	if g.py <= padding {
// 		g.py = padding + 1
// 	}
// }

// // func rayVertices(x1, y1, x2, y2, x3, y3 float64) []ebiten.Vertex {
// // 	return []ebiten.Vertex{
// // 		{DstX: float32(x1), DstY: float32(y1), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
// // 		{DstX: float32(x2), DstY: float32(y2), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
// // 		{DstX: float32(x3), DstY: float32(y3), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
// // 	}
// // }

// type Game struct {
// 	showRays bool
// 	px, py   int
// 	objects  []*object
// }

// func (g *Game) Update() error {
// 	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
// 		return errors.New("game ended by player")
// 	}

// 	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
// 		g.showRays = !g.showRays
// 	}

// 	g.handleMovement()

// 	return nil
// }

// func (g *Game) Draw(screen *ebiten.Image) {
// 	for _, elem := range g.objects {
// 		elem.visible = false
// 	}
// 	// Reset the shadowImage
// 	// shadowImage.Fill(color.Black)
// 	rays := rayCasting(float64(g.px), float64(g.py), g.objects)

// 	// Subtract ray triangles from shadow
// 	// opt := &ebiten.DrawTrianglesOptions{}
// 	// opt.Address = ebiten.AddressRepeat
// 	// opt.Blend = ebiten.BlendSourceOut
// 	// for i, line := range rays {
// 	// 	nextLine := rays[(i+1)%len(rays)]

// 	// 	// Draw triangle of area between rays
// 	// 	v := rayVertices(float64(g.px), float64(g.py), nextLine.X2, nextLine.Y2, line.X2, line.Y2)
// 	// 	shadowImage.DrawTriangles(v, []uint16{0, 1, 2}, triangleImage, opt)
// 	// }

// 	// Draw background
// 	// screen.DrawImage(bgImage, nil)

// 	// // Draw rays
// 	// for _, r := range rays {
// 	// 	vector.StrokeLine(screen, float32(r.X1), float32(r.Y1), float32(r.X2), float32(r.Y2), 1, color.RGBA{255, 255, 0, 150}, true)
// 	// }

// 	// Draw shadow.
// 	// op := &ebiten.DrawImageOptions{}
// 	// op.ColorScale.ScaleAlpha(0.7)
// 	// screen.DrawImage(shadowImage, op)

// 	dist := func(x, y int, wallLine line) float64 {
// 		dx := math.Abs(float64(x) - wallLine.X2)
// 		dy := math.Abs(float64(y) - wallLine.Y2)

// 		return math.Sqrt(dx * dx * dy * dy)
// 	}

// 	// Draw walls
// 	for _, obj := range g.objects {
// 		c := color.RGBA{255, 0, 0, 255}
// 		if obj.visible {
// 			c = color.RGBA{0, 255, 0, 255}
// 		}
// 		for _, w := range obj.walls {
// 			maxH := 64.
// 			_ = maxH
// 			projectionPlanDepth := 277.
// 			_ = projectionPlanDepth
// 			h := (maxH / dist(g.px, g.py, w)) * projectionPlanDepth
// 			_ = h

// 			// x1/y1---------------------x2/y2
// 			// |
// 			// |
// 			// |
// 			// x1/y1+h------------------x2/y2+h

// 			// Bottom.
// 			if !obj.visible {
// 				continue
// 			}
// 			// vector.DrawFilledRect(screen, float32(w.X1), float32(w.Y1), float32(math.Abs(w.X1-w.X2)), float32(math.Abs(w.Y1-w.Y2)), c, true)
// 			vector.StrokeLine(screen, float32(w.X1), float32(w.Y1), float32(w.X2), float32(w.Y2), 1, c, true)

// 			// // Top.
// 			// c.G, c.R = 0, 100
// 			vector.StrokeLine(screen, float32(w.X1), float32(w.Y1+h), float32(w.X2), float32(w.Y2+h), 1, c, true)

// 			// // Left.
// 			// c.G, c.R = 200, 50
// 			vector.StrokeLine(screen, float32(w.X1), float32(w.Y1), float32(w.X1), float32(w.Y1+h), 1, c, true)

// 			// // Right.
// 			// c.G, c.R = 255, 255
// 			vector.StrokeLine(screen, float32(w.X2), float32(w.Y2), float32(w.X2), float32(w.Y2+h), 1, c, true)
// 		}
// 	}

// 	// Draw player as a rect
// 	vector.DrawFilledRect(screen, float32(g.px)-2, float32(g.py)-2, 4, 4, color.Black, true)
// 	vector.DrawFilledRect(screen, float32(g.px)-1, float32(g.py)-1, 2, 2, color.RGBA{255, 100, 100, 255}, true)

// 	ebitenutil.DebugPrintAt(screen, "WASD: move", 160, 0)
// 	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()), 51, 51)
// 	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Rays: 2*%d", len(rays)/2), padding, 222)
// 	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d/%d", g.px, g.py), 180, 222)
// }

// func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
// 	return screenWidth, screenHeight
// }

// func rect(x, y, w, h float64) []line {
// 	return []line{
// 		{x, y, x, y + h},
// 		{x, y + h, x + w, y + h},
// 		{x + w, y + h, x + w, y},
// 		{x + w, y, x, y},
// 	}
// }

// func main() {
// 	g := &Game{
// 		px: 40,
// 		py: 160,
// 	}

// 	// Add outer walls
// 	g.objects = append(g.objects, &object{walls: rect(padding, padding, screenWidth-2*padding, screenHeight-2*padding)})

// 	// Angled wall
// 	g.objects = append(g.objects, &object{walls: []line{{50, 110, 100, 150}}})

// 	// Rectangles
// 	for _, l := range rect(45, 50, 70, 20) {
// 		g.objects = append(g.objects, &object{walls: []line{l}})
// 	}
// 	for _, l := range rect(150, 50, 30, 60) {
// 		g.objects = append(g.objects, &object{walls: []line{l}})
// 	}

// 	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
// 	ebiten.SetWindowTitle("Ray casting and shadows (Ebitengine Demo)")
// 	if err := ebiten.RunGame(g); err != nil {
// 		log.Fatal(err)
// 	}
// }
