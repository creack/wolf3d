// Package main is the entrypoint.
//
// Ref: https://lodev.org/cgtutor/raycasting.html
package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"

	"go.creack.net/wolf3d/math2"
)

//go:embed textures.png
var textureData []byte

//go:embed maps/map4
var mapData []byte

func loadTextures(textureData []byte) (*image.RGBA, error) {
	p, err := png.Decode(bytes.NewReader(textureData))
	if err != nil {
		return nil, fmt.Errorf("png.Decode: %w", err)
	}
	m := image.NewRGBA(p.Bounds())
	draw.Draw(m, m.Bounds(), p, image.Point{}, draw.Src)
	return m, nil
}

func main() {
	world, err := parseMap(mapData)
	if err != nil {
		log.Fatal(err)
	}
	textures, err := loadTextures(textureData)
	if err != nil {
		log.Fatal(err)
	}

	g := &Game{
		width:  1280,
		height: 1024,

		world:    world,
		textures: textures,

		pos:   math2.Pt(12.0, 11.5),
		dir:   math2.Pt(-1.0, 0.0),
		plane: math2.Pt(0.0, 0.66),
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
