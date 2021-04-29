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
func projectP(p [][]float64) (float64, float64) {
	z := p[2][0]
	scrDist := xlim/(math.Tan(fov/2)*2) // this is essentially how far is the screen from eye
	p = matScale(p, scrDist/(z*math.Tan(fov/2)))
	return p[0][0], p[1][0]
}

/*draws a projected line from 3d to 2d*/
func line3d(p, q [][]float64, board []int) {
	x1, y1 := projectP(p)
	x2, y2 := projectP(q)
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

/*returns rotation matrix of x, y, z angle from respective axes*/
func rotateP3d(x, y, z float64) [][]float64 {
    rotMat := [][]float64 {
        {math.Cos(z), -math.Sin(z), 0},
        {math.Sin(z), math.Cos(z), 0},
        {0, 0, 1},
    }
    otherMat := [][]float64 { // i cross k is -j
        {math.Cos(y), 0, math.Sin(y)},
        {0, 1, 0},
        {-math.Sin(y), 0, math.Cos(y)},
    }
    rotMat = matMul(rotMat, otherMat)
    otherMat = [][]float64 {
        {1, 0, 0},
        {0, math.Cos(x), -math.Sin(x)},
        {0, math.Sin(x), math.Cos(x)},
    }
    return matMul(rotMat, otherMat)
}

