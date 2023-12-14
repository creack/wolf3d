// Package main is the entrypoint.
//
// Ref: https://lodev.org/cgtutor/raycasting.html
package main

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"go.creack.net/wolf3d/math2"
)

//go:embed textures.png
var textureData []byte

//go:embed maps/*
var mapData embed.FS

func loadTextures(textureData []byte) (front, side *image.RGBA, err error) {
	p, err := png.Decode(bytes.NewReader(textureData))
	if err != nil {
		return nil, nil, fmt.Errorf("png.Decode: %w", err)
	}
	front = image.NewRGBA(p.Bounds())
	draw.Draw(front, front.Bounds(), p, image.Point{}, draw.Src)

	side = image.NewRGBA(p.Bounds())
	draw.Draw(side, side.Bounds(), p, image.Point{}, draw.Src)
	for y := 0; y < side.Rect.Dy(); y++ {
		for x := 0; x < side.Rect.Dx(); x++ {
			side.Set(x, y, dimColor(side.At(x, y)))
		}
	}

	return front, side, nil
}

func main() {
	textures, sideTextures, err := loadTextures(textureData)
	if err != nil {
		log.Fatal(err)
	}
	g := &Game{
		width:  1280,
		height: 720,
		// width:  3840,
		// height: 2160,

		textures:     textures,
		sideTextures: sideTextures,

		dir:   math2.Pt(1, 0),
		plane: math2.Pt(0, 0.66),

		last: time.Now(),

		showRays: false,
		mapMod:   0,
	}
	if err := g.loadMap("maps/map4"); err != nil {
		log.Fatal(err)
	}
	// g.texturesCache = make([][][3]byte, g.height)
	for y := range g.texturesCache {
		// g.texturesCache[y] = make([][3]byte, g.width)
		for x := range g.texturesCache[y] {
			// g.texturesCache[x][y] = make([]byte, 3)
			r1, g1, b1, _ := g.textures.At(x, y).RGBA()
			g.texturesCache[y][x][0] = byte(r1)
			g.texturesCache[y][x][1] = byte(g1)
			g.texturesCache[y][x][2] = byte(b1)
			// g.texturesCache[y][x] = g.textures.At(x, y)
		}
	}
	// g.sideTexturesCache = make([][][3]byte, g.height)
	for y := range g.sideTexturesCache {
		// g.sideTexturesCache[y] = make([][3]byte, g.width)
		for x := range g.sideTexturesCache[y] {
			// g.sideTexturesCache[x][y] = make([]byte, 3)
			r1, g1, b1, _ := g.sideTextures.At(x, y).RGBA()
			g.sideTexturesCache[y][x][0] = byte(r1)
			g.sideTexturesCache[y][x][1] = byte(g1)
			g.sideTexturesCache[y][x][2] = byte(b1)
			// g.sideTexturesCache[y][x] = g.sideTextures.At(x, y)
		}
	}
	g.triangleImg = ebiten.NewImage(g.width, g.height)
	g.triangleImg.Fill(color.White)

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
