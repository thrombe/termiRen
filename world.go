package main

import (
	"bufio"
	// "math"
	"os"
	"strconv"
	"strings"
	// "fmt"
)

// helper functions

//multiplies matrix with each coord.
// enter 1 or 2 vector lists, multiplication of mat and 1st is saved in second (or 1st if second not entered)
func transform(mat [][]float64, verlists ...[][][]float64) {
    var vertices, savein [][][]float64
    if len(verlists) == 2 {
        vertices = verlists[0]
        savein = verlists[1]
    } else {
        vertices = verlists[0]
        savein = verlists[0]
    }
    length := len(vertices)
    for i := 0; i < length; i++ {
        // savein[i] = matMul(mat, vertices[i])
        vecMul(mat, vertices[i], savein[i])
        if savein[i][3][0] != 1 { // if the forth column if vector isnt 1, then scale the vector
            if round(savein[i][3][0]*1000) == 0 {panic("divide by 0 in transform()")}
            savein[i] = matScalar(savein[i], 1/savein[i][3][0])//*math.Tan(fov/2))) // is the cot extra??????????????
        }
    }
}

// objects

type triangle struct{
    vertices []*[][]float64
    camtices []*[][]float64 // the world matrix is multiplied with vertices and is stored here. for methods on triangke
}

func (tri *triangle) create(vers ...*[][]float64) {
    if len(vers) == 3 {
        tri.vertices = []*[][]float64 {vers[0], vers[1], vers[2]}
        tri.camtices = []*[][]float64 {vers[0], vers[1], vers[2]}
    } else if len(vers) == 6 {
        tri.vertices = []*[][]float64 {vers[0], vers[1], vers[2]}
        tri.camtices = []*[][]float64 {vers[3], vers[4], vers[5]}
    } else {
        panic("tri.create strange no. of vertices given")
    }
}

// returns the normal of the triangle in 3d space (tri.vertices).
// normal in the direction of the visible face (anticlocck)
func (tri *triangle) normal() [][]float64 {
    return vecCross(matSub(*tri.vertices[1], *tri.vertices[0]), matSub(*tri.vertices[2], *tri.vertices[0]))
}

//draws triangle using camtices
func (tri *triangle) draw(cam *camera, board [][] byte, texture  byte) {
    if vecDot(matSub(*tri.vertices[0], cam.camPos), tri.normal()) >= 0 {return} // if the front(anticlockwise) face of triangle faces away from/perpendicular to cam, dont draw 
    // >= cuz both vectors have different origin

    charRatio, xlim, ylim := cam.charRatio, cam.xlim, cam.ylim
    for i := 0; i < 3; i++ {
        line(*tri.camtices[i], *tri.camtices[(i+1)%3], charRatio, xlim, ylim, board, texture)
    }
}

//fills up triangle using camtices
func (tri *triangle) fill(cam *camera, board [][] byte, zbuf [][]float64) {
    if vecDot(matSub(*tri.vertices[0], cam.camPos), tri.normal()) >= 0 {return} // if the front(anticlockwise) face of triangle faces away from/perpendicular to cam, dont draw 
    // >= cuz both vectors have different origin
    
    i := 0
    for _, vert := range tri.camtices { // if traingle is outside the screen, dont draw
        if absVal((*vert)[0][0]) > float64(cam.xlim)/2 || absVal((*vert)[1][0])*cam.charRatio > float64(cam.ylim)/2 {i++}
    }
    if i == 3 {return}

    xlim, ylim, charRatio := cam.xlim, cam.ylim, cam.charRatio
    
    lightDir := matSub(*tri.vertices[0], cam.camPos)
    // lightDir := [][]float64 {{-1}, {-1}, {-2}, {0}} // from +z to -z
    tex := vecDot(vecUnit(lightDir), vecUnit(tri.normal())) // 0 to 1
    // textures := ".`^,:;Il!i~+_-?][}{!)(|/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
    textures := "':;!vx1WnOM@B"
    tex = -tex*float64(len(textures)-1)
    if tex < 0 {tex = 0}
    texture :=  byte(textures[round(tex)])
    
    v1 := *tri.camtices[0]
    v21 := matSub(*tri.camtices[1], v1)
    v31 := matSub(*tri.camtices[2], v1)
    dw2, dw3 := 1.0/vecSize(v21), 1.0/vecSize(v31)
    for w2 := 1.0; w2 >= 0; w2 -= dw2 {
        for w3 := 0.0; w3 <= 1-w2; w3 += dw3 {
            poin := nMatAdd(v1, matScalar(v21, w2), matScalar(v31, w3))
            point3d(poin, charRatio, xlim, ylim, board, zbuf, texture)
        }
        poin := nMatAdd(v1, matScalar(v21, w2), matScalar(v31, 1-w2)) // just to make sure there are no holes (works pretty well)
        point3d(poin, charRatio, xlim, ylim, board, zbuf, texture)
    }
}

type object struct {
    vertices [][][]float64
    camtices [][][]float64
    triangles []triangle
    center [][]float64
}

func (ob *object) create(path string, o [][]float64) {
    file, err := os.Open(path)
	defer file.Close()
	if err != nil {panic(err)}

	sc := bufio.NewScanner(file) // parsing file
	var faces [][]float64
    var bboxmin, bboxmax [][]float64
	for sc.Scan() {
		text := sc.Text()
        if len(text) == 0 {continue}
		switch text[ : 2] {
		case "v ":
			texx := strings.Split(text[2:], " ")
            vertex := matrix(4, 1)
			for i := 0; i < 3; i++ {
				num, err := strconv.ParseFloat(texx[i], 64)
				if err != nil {panic(err)}
				vertex[i][0] = num
			}
			vertex[3][0] = 1
            vertex = matAdd(vertex, o)
			vertex[3][0] = 1 // if 0 has a 1 in 4th col, it could cause probs
			ob.vertices = append(ob.vertices, vertex)
            ob.camtices = append(ob.camtices, vector(vertex[0][0], vertex[1][0], vertex[2][0], vertex[3][0]))
            if bboxmin == nil {bboxmin, bboxmax = vector(vertex[0][0], vertex[1][0], vertex[2][0], vertex[3][0]), vector(vertex[0][0], vertex[1][0], vertex[2][0], vertex[3][0])}
            if bboxmin[0][0] > vertex[0][0] {bboxmin[0][0] = vertex[0][0]}
            if bboxmin[1][0] > vertex[1][0] {bboxmin[1][0] = vertex[1][0]}
            if bboxmin[2][0] > vertex[2][0] {bboxmin[2][0] = vertex[2][0]}
            if bboxmax[0][0] < vertex[0][0] {bboxmax[0][0] = vertex[0][0]}
            if bboxmax[1][0] < vertex[1][0] {bboxmax[1][0] = vertex[1][0]}
            if bboxmax[2][0] < vertex[2][0] {bboxmax[2][0] = vertex[2][0]}
		case "f ":
			texx := strings.Split(text[2:], " ")
			face := make([]float64, 3)
			for i := 0; i < 3; i++ {
				face[i], err = strconv.ParseFloat(strings.Split(texx[i], "/")[0], 64)
				if err != nil {panic(err)}
				face[i]--
			}
			// face[0], face[2] = face[2], face[0] // clockwise to anticlockwise
			faces = append(faces, face)
		}
 	}

    ob.center = matScalar(matAdd(bboxmin, bboxmax), 0.5)

 	for _, face := range faces { // creating triangles
        a := &ob.vertices[int(face[0])]
        b := &ob.vertices[int(face[1])]
        c := &ob.vertices[int(face[2])]
        d := &ob.camtices[int(face[0])]
        e := &ob.camtices[int(face[1])]
        f := &ob.camtices[int(face[2])]
       tri := triangle{}
 	    tri.create(a, b, c, d, e, f)
 	    ob.triangles = append(ob.triangles, tri)
 	}
}

func (ob *object) draw(cam *camera, cammat [][]float64, board [][] byte, texture  byte) {
    transform(cammat, ob.vertices, ob.camtices)
    for _, tri := range ob.triangles {
        tri.draw(cam, board, texture)
    }
}

func (ob *object) fill(cam *camera, cammat [][]float64, board [][] byte, zbuf [][]float64) {
    transform(cammat, ob.vertices, ob.camtices)
    for _, tri := range ob.triangles {
        tri.fill(cam, board, zbuf)
    }
}

//multiplies the mat with coords of each triangle
func (ob *object) transform(mat [][]float64) {
    transform(mat, ob.vertices)
}