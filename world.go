package main

import (
	"math"
	// "fmt"
)

/*rotate x, y about ox, oy by r radians*/
func rotateP2d(x, y, ox, oy, r float64) (float64, float64) {
    x, y = x-ox, y-oy
    x, y =  x*math.Cos(r)-y*math.Sin(r), x*math.Sin(r)+y*math.Cos(r)
    return x+ox, y+oy
}

/*converts coords from 3d space to 2d so that it can be drawn on canvas*/
func project(p [][]float64) (float64, float64) {
	z := p[2][0]*distanceScaleFactor
	p = matScale(p, 1/z)
	return p[0][0], p[1][0]
}

func line3d(p, q [][]float64, board []int) {
	x1, y1 := project(p)
	x2, y2 := project(q)
	line(x1, y1, x2, y2, board)
}

/*draws a square with edges parallel to x, y axes*/
func xysquare3d(o [][]float64, side float64, board []int) {
	r2 := 1/math.Sqrt(2)
	u := [][]float64 {{r2}, {r2}, {0}}
	v := [][]float64 {{-r2}, {r2}, {0}}
	diagH := side/math.Sqrt(2)
	u = matScale(u, diagH)
	v = matScale(v, diagH)
	line3d(matAdd(o, v), matAdd(o, u), board)
	line3d(matAdd(o, v), matSub(o, u), board)
	line3d(matAdd(o, u), matSub(o, v), board)
	line3d(matSub(o, u), matSub(o, v), board)
}