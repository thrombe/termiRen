package main

import (
	"math"
)

/*draws a projected line from 3d to 2d*/
func line3d(p, q [][]float64, board []int) { // 3 by 1 vectors
	x1, y1 := projectP(p)
	x2, y2 := projectP(q)
	line(x1, y1, x2, y2, board)
}

type cuboid struct{
    coords [][]float64
}

/*creates vertices of cube from given centre and half diagonal vectors.
stores them in cube.coords in a single 4 by 8 matrix
o is centre and u is half diagonal vector*/
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
func (cu *cuboid) draw(board []int, coords [][]float64) {
    vertices := make([][][]float64, 8)
    for i := 0; i < 8; i++ {
        vertices[i] = getCoord3d(coords, i)
    }
    for i := 0; i < 4; i++ { // connecting vertices by lines
        line3d(vertices[i], vertices[(i+1)%4], board)
        line3d(vertices[i+4], vertices[4+(i+1)%4], board)
        line3d(vertices[i], vertices[i+4], board)
    }
}

type triangle struct{
    vertices [][][]float64
    camtices [][][]float64 // the world matrix is multiplied with vertices and is storesd here. for methods on triangke
}

/*multiplies matrix with each coord*/
func (tri *triangle) transform(mat [][]float64, vertices [][][]float64) {
    length := len(vertices)
    for i := 0; i < length; i++ {
        vertices[i] = matMul(mat, vertices[i])
    }
}

func (tri *triangle) create(a, b, c [][]float64) {
    tri.vertices = [][][]float64 {a, b, c}
}

func (tri *triangle) normal(vertices [][][]float64) [][]float64 {
    return vecCross(matSub(vertices[2], vertices[0]), matSub(vertices[1], vertices[0]))
}

func (tri *triangle) draw(board []int, vertices [][][]float64) {
    if vecDot(vertices[0], tri.normal(vertices)) <= 0 {return} // if the front(clockwise) face of triangle faces away from/perpendicular to cam, dont draw
    for i := 0; i < 3; i++ {
        line3d(vertices[i], vertices[(i+1)%3], board)
    }
}

/*fills up triangle*/
func (tri *triangle) fill(board []int, vertices [][][]float64) {
    if vecDot(vertices[0], tri.normal(vertices)) <= 0 {return} // if the front(clockwise) face of triangle faces away from/perpendicular to cam, dont draw
    for i := 0; i < 3; i++ { // convert world coords into screen coords
        x, y := projectP(vertices[i])
        vertices[i] = [][]float64 {{x}, {y}, {0}}
    }
    minx, miny, maxx, maxy := vertices[0][0][0], vertices[0][1][0], vertices[0][0][0], vertices[0][1][0]
    // add condition for if any coord goes outside screen, then chop
    for _, vertex := range vertices {
        if vertex[0][0] > maxx {maxx = vertex[0][0]}
        if vertex[0][0] < minx {minx = vertex[0][0]}
        if vertex[1][0] > maxy {maxy = vertex[1][0]}
        if vertex[1][0] < miny {miny = vertex[1][0]}
    }
    triangle := inTriangle(vertices)
    for y := miny; y <= maxy; y++ {
        for x := minx; x <= maxx; x++ {
            pp := [][]float64 {{x}, {y}, {0}}
            if triangle(pp) {point(x, y, board)}
        }
    }
}

/*returns a func that returns if a point lies in a triangle*/
func inTriangle(vertices [][][]float64) func([][]float64) bool {
    v1v2 := matSub(vertices[1], vertices[0])
    v1v3 := matSub(vertices[2], vertices[0])
    v2v3 := matSub(vertices[2], vertices[1])
    return func(point [][]float64) bool {
        v1p := matSub(point, vertices[0])
        v2p := matSub(point, vertices[1])
        if vecCross(v1v2, v1p)[2][0] < 0 {return false}
        if vecCross(v1v3, v1p)[2][0] > 0 {return false}
        if vecCross(v2v3, v2p)[2][0] < 0 {return false}
        return true
    }
}

func (tri *triangle) fill2(board []int, vertices [][][]float64) {
    // if vecDot(vertices[0], tri.normal(vertices)) <= 0 {return} // if the front(clockwise) face of triangle faces away from/perpendicular to cam, dont draw
    for i := 0; i < 3; i++ { // convert world coords into screen coords
        x, y := projectP(vertices[i])
        vertices[i] = [][]float64 {{x}, {y}}
    }
    if vertices[0][1][0] < vertices[1][1][0] {vertices[0], vertices[1] = vertices[1], vertices[0]} // arrange coords in descending order (y)
    if vertices[0][1][0] < vertices[2][1][0] {vertices[0], vertices[2] = vertices[2], vertices[0]}
    if vertices[1][1][0] < vertices[2][1][0] {vertices[1], vertices[2] = vertices[2], vertices[1]}
    longline := matSub(vertices[2], vertices[0])
    shortline1 := matSub(vertices[1], vertices[0])
    shortline2 := matSub(vertices[2], vertices[1])
    // what if? - 1) all lines have 0 dy, 2) either of the shortline has 0 dy
    // if round(longline[1][0]) == 0 { // for case 1
        // if vertices[0][0][0] < vertices[1][0][0] {vertices[0], vertices[1] = vertices[1], vertices[0]} // arrange coords in descending order (x)
        // if vertices[0][0][0] < vertices[2][0][0] {vertices[0], vertices[2] = vertices[2], vertices[0]}
        // if vertices[1][0][0] < vertices[2][0][0] {vertices[1], vertices[2] = vertices[2], vertices[1]}
        // line(vertices[0][0][0], vertices[0][1][0], vertices[2][0][0], vertices[2][1][0], board)
    // } else if round(shortline1[1][0]) == 0 { // for case 2.1
        // line(vertices[0][0][0], vertices[0][1][0], vertices[1][0][0], vertices[1][1][0], board)        
    // } else if round(shortline2[1][0]) == 0 { // for case 2.2
        // line(vertices[1][0][0], vertices[1][1][0], vertices[2][0][0], vertices[2][1][0], board)        
    // }
    longdxdy := longline[0][0]/longline[1][0]
    short1dxdy := shortline1[0][0]/shortline1[1][0]
    short2dxdy := shortline2[0][0]/shortline2[1][0]
    longx, shortx := vertices[0][0][0], vertices[0][0][0]
    var lolim, hilim float64
    for y := math.Round(vertices[0][1][0]); y >= vertices[1][1][0]; y-- {
        // if longx > shortx {hilim, lolim = longx, shortx} else {lolim, hilim = longx, shortx}    
        hilim, lolim = longx, shortx
        for x := lolim; x < hilim; x++ {
            point(x, y, board)
        }
        longx -= longdxdy
        shortx -= short1dxdy
    }
    for y := math.Round(vertices[1][1][0]); y >= vertices[2][1][0]; y-- {
        // if longx > shortx {hilim, lolim = longx, shortx} else {lolim, hilim = longx, shortx}
        hilim, lolim = longx, shortx
        for x := lolim; x < hilim; x++ {
            point(x, y, board)
        }
        longx -= longdxdy
        shortx -= short2dxdy
    }
}