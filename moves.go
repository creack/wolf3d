package main

import (
	"go.creack.net/wolf3d/math2"
)

func (g *Game) moveForward(s float64) {
	if g.world[int(g.pos.X+g.dir.X*s)][int(g.pos.Y)].wallType == 0 {
		g.pos.X += g.dir.X * s
	}

	if g.world[int(g.pos.X)][int(g.pos.Y+g.dir.Y*s)].wallType == 0 {
		g.pos.Y += g.dir.Y * s
	}
}

func (g *Game) moveLeft(s float64) {
	if g.world[int(g.pos.X-g.plane.X*s)][int(g.pos.Y)].wallType == 0 {
		g.pos.X -= g.plane.X * s
	}

	if g.world[int(g.pos.X)][int(g.pos.Y-g.plane.Y*s)].wallType == 0 {
		g.pos.Y -= g.plane.Y * s
	}
}

func (g *Game) moveBackwards(s float64) {
	if g.world[int(g.pos.X-g.dir.X*s)][int(g.pos.Y)].wallType == 0 {
		g.pos.X -= g.dir.X * s
	}

	if g.world[int(g.pos.X)][int(g.pos.Y-g.dir.Y*s)].wallType == 0 {
		g.pos.Y -= g.dir.Y * s
	}
}

func (g *Game) moveRight(s float64) {
	if g.world[int(g.pos.X+g.plane.X*s)][int(g.pos.Y)].wallType == 0 {
		g.pos.X += g.plane.X * s
	}

	if g.world[int(g.pos.X)][int(g.pos.Y+g.plane.Y*s)].wallType == 0 {
		g.pos.Y += g.plane.Y * s
	}
}

func (g *Game) turnRight(s float64) {
	g.dir = g.dir.Rotate(math2.Angle(-s))
	g.plane = g.plane.Rotate(math2.Angle(-s))
}

func (g *Game) turnLeft(s float64) {
	g.dir = g.dir.Rotate(math2.Angle(s))
	g.plane = g.plane.Rotate(math2.Angle(s))
}
