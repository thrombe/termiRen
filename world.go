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

type cuboid struct{
    coords [][]float64
}

func matAppend(mat [][]float64, mats... [][]float64) {
    // use this to represent multiple coords in 1 matrix
    for c := 0; c < len(mat); c++ {
        for _, mat1 := range mats {
            mat[c] = append(mat[c], mat1[c]...)
        }
    }
}

func getCoord(mat [][]float64, n int) [][]float64 {
    return [][]float64 {
        {mat[0][n]},
        {mat[1][n]},
        {mat[2][n]},
    }
}

func (cu *cuboid) create(o, u [][]float64) {
    // use displace from centre method, so to rotate the cube,
    // only that displacement vector needs to be rotated
    // creating cube parallel to axes by default
    //size := vecSize(u)
    //size = size/2
    u = [][]float64 {
        {u[0][0], u[0][0], -u[0][0], -u[0][0], u[0][0], u[0][0], -u[0][0], -u[0][0]},
        {u[1][0], u[1][0], u[1][0], u[1][0], -u[1][0], -u[1][0], -u[1][0], -u[1][0]},
        {u[2][0], -u[2][0], -u[2][0], u[2][0], u[2][0], -u[2][0], -u[2][0], u[2][0]},
    }
    cu.coords = make([][]float64, 3)
    matAppend(cu.coords, o, o, o, o, o, o, o, o)
    cu.coords = matAdd(cu.coords, u)
}

func (cu *cuboid) draw(board []int) {
    vertices := make([][][]float64, 8)
    for i := 0; i < 8; i++ {
        vertices[i] = getCoord(cu.coords, i)
    }
    for i := 0; i < 4; i++ {
        line3d(vertices[i], vertices[(i+1)%4], board)
        line3d(vertices[i+4], vertices[4+(i+1)%4], board)
        line3d(vertices[i], vertices[i+4], board)
    }
}