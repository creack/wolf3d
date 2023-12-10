package main

import "math"

func (g *Game) moveForward(s float64) {
	if g.world[int(g.pos.X+g.dir.X*s)][int(g.pos.Y)] == 0 {
		g.pos.X += g.dir.X * s
	}

	if g.world[int(g.pos.X)][int(g.pos.Y+g.dir.Y*s)] == 0 {
		g.pos.Y += g.dir.Y * s
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
