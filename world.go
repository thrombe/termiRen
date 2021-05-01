package main

import (
	"math"
	// "fmt"
)

/*converts coords from 3d space to 2d so that it can be drawn on canvas*/
func projectP(p [][]float64) (float64, float64) { // 3 by 1 vectors
	z := p[2][0]
	scrDist := xlim/(math.Tan(fov/2)*2) // this is essentially how far is the screen from eye
	p = matScale(p, scrDist/(z*math.Tan(fov/2)))
	return p[0][0], p[1][0]
}

/*draws a projected line from 3d to 2d*/
func line3d(p, q [][]float64, board []int) { // 3 by 1 vectors
	x1, y1 := projectP(p)
	x2, y2 := projectP(q)
	line(x1, y1, x2, y2, board)
}

/*draws a square with edges parallel to x, y axes*/
func xysquare3d(o [][]float64, side float64, board []int) { // change this to draw any kind of rectangle
    o = o[ : 3]
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

/*returns rotation matrix z angle from z axis*/
func rotateP3dz(z float64) [][]float64 {
    rotMat := [][]float64 {
        {math.Cos(z), -math.Sin(z), 0, 0},
        {math.Sin(z), math.Cos(z), 0, 0},
        {0, 0, 1, 0},
        {0, 0, 0, 1},
    }
    return rotMat
}

/*returns rotation matrix of y angle from y axis*/
func rotateP3dy(y float64) [][]float64 {
    rotMat := [][]float64 { // i cross k is -j
        {math.Cos(y), 0, math.Sin(y), 0},
        {0, 1, 0, 0},
        {-math.Sin(y), 0, math.Cos(y), 0},
        {0, 0, 0, 1},
    }
    return rotMat
}

/*returns rotation matrix of x angle from x axis*/
func rotateP3dx(x float64) [][]float64 {
    rotMat := [][]float64 {
        {1, 0, 0, 0},
        {0, math.Cos(x), -math.Sin(x), 0},
        {0, math.Sin(x), math.Cos(x), 0},
        {0, 0, 0, 1},
    }
    return rotMat
}

/*returns translation matrix (subtracts o from coords)*/
func transnMat(o [][]float64) [][]float64 {
    return transpMat(matScale(o, -1)) // transpMat ignores the last row of o ans always replaces it with 1
}

/*returns translation matrix (adds o to coords)*/
func transpMat(o [][]float64) [][]float64 {
    length := len(o)
    trans := make([][]float64, length)
    for i, ele := range o {
        row := []float64 {0, 0, 0, ele[0]}
        row[i] = 1
        trans[i] = row
    }
    row := []float64 {0, 0, 0, 1}
    trans[length-1] = row
    return trans
}

/*adds back and forth translation to rotation matrix*/
func rotAboutPoint(rot ,o [][]float64) [][]float64 {
    return matMul(transpMat(o), matMul(rot, transnMat(o)))
}

type cuboid struct{
    coords [][]float64
}

/*use this to represent multiple coords in 1 matrix*/
func matAppend(mat [][]float64, mats... [][]float64) {
    for c := 0; c < len(mat); c++ {
        for _, mat1 := range mats {
            mat[c] = append(mat[c], mat1[c]...)
        }
    }
}

/*returns a 3d vector from the given column no.*/
func getCoord3d(mat [][]float64, n int) [][]float64 {
    return [][]float64 { // convert this into a n dimentional instead of jusst 3
        {mat[0][n]},
        {mat[1][n]},
        {mat[2][n]},
    }
}

/*creates vertices of cube from given centre and half diagonal vectors.
sores them in cube.coords in a single 4 by 8 matrix*/
func (cu *cuboid) create(o, u [][]float64) {
    // creating cube parallel to axes by default
    u = [][]float64 { // another method for this is rotating the original u by some angle (pi/2 in case of cubes) in different planes
        {u[0][0], u[0][0], -u[0][0], -u[0][0], u[0][0], u[0][0], -u[0][0], -u[0][0]},
        {u[1][0], u[1][0], u[1][0], u[1][0], -u[1][0], -u[1][0], -u[1][0], -u[1][0]},
        {u[2][0], -u[2][0], -u[2][0], u[2][0], u[2][0], -u[2][0], -u[2][0], u[2][0]},
        {0, 0, 0, 0, 0, 0, 0, 0},
    }
    cu.coords = make([][]float64, 4)
    matAppend(cu.coords, o, o, o, o, o, o, o, o)
    cu.coords = matAdd(cu.coords, u)
}

/*draws the cuboid on canvas*/
func (cu *cuboid) draw(board []int) {
    vertices := make([][][]float64, 8)
    for i := 0; i < 8; i++ {
        vertices[i] = getCoord3d(cu.coords, i)
    }
    for i := 0; i < 4; i++ { // connecting vertices by lines
        line3d(vertices[i], vertices[(i+1)%4], board)
        line3d(vertices[i+4], vertices[4+(i+1)%4], board)
        line3d(vertices[i], vertices[i+4], board)
    }
}