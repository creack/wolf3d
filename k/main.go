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
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed textures.png
var textureData []byte

//go:embed map4
var mapData []byte

func loadWorld() [][]int {
	//nolint:prealloc // False positive.
	var w [][]int
	for _, line := range strings.Split(string(mapData), "\n") {
		if line == "" || line[0] == '#' {
			continue
		}
		var worldLine []int
		for _, elem := range strings.Split(line, " ") {
			n, err := strconv.Atoi(elem)
			if err != nil {
				panic(fmt.Errorf("parse map case %q: %w", elem, err))
			}
			worldLine = append(worldLine, n)
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

func main() {
	g := &Game{
		width:  640,
		height: 480,

		world:    loadWorld(),
		textures: loadTextures(),

		pos:   Pt(22.0, 11.5),
		dir:   Pt(-1.0, 0.0),
		plane: Pt(0.0, 0.66),
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
