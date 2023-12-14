package main

import "go.creack.net/wolf3d/math2"

func (g *Game) moveForward(s float64) {
	if newX := g.pos.X + g.dir.X*s; g.getTexNum(int(newX), int(g.pos.Y)) == 0 {
		g.pos.X = newX
	}
	if newY := g.pos.Y + g.dir.Y*s; g.getTexNum(int(g.pos.X), int(newY)) == 0 {
		g.pos.Y = newY
	}
}

func (g *Game) moveLeft(s float64) {
	if newX := g.pos.X - g.plane.X*s; g.getTexNum(int(newX), int(g.pos.Y)) == 0 {
		g.pos.X = newX
	}
	if newY := g.pos.Y - g.plane.Y*s; g.getTexNum(int(g.pos.X), int(newY)) == 0 {
		g.pos.Y = newY
	}
}

func (g *Game) moveBackwards(s float64) {
	if newX := g.pos.X - g.dir.X*s; g.getTexNum(int(newX), int(g.pos.Y)) == 0 {
		g.pos.X = newX
	}
	if newY := g.pos.Y - g.dir.Y*s; g.getTexNum(int(g.pos.X), int(newY)) == 0 {
		g.pos.Y = newY
	}
}

func (g *Game) moveRight(s float64) {
	if newX := g.pos.X + g.plane.X*s; g.getTexNum(int(newX), int(g.pos.Y)) == 0 {
		g.pos.X = newX
	}
	if newY := g.pos.Y + g.plane.Y*s; g.getTexNum(int(g.pos.X), int(newY)) == 0 {
		g.pos.Y = newY
	}
}

func (g *Game) turnRight(s float64) {
	g.dir = g.dir.Rotate(math2.Angle(s))
	g.plane = g.plane.Rotate(math2.Angle(s))
}

func (g *Game) turnLeft(s float64) {
	g.dir = g.dir.Rotate(math2.Angle(-s))
	g.plane = g.plane.Rotate(math2.Angle(-s))
}
