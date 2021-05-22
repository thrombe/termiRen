package main

import (
	"math"
    "strings"
    "strconv"
    "bufio"
    "os"
    // "fmt"
)

// helper functions

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
    return vector( // convert this into a n dimentional instead of jusst 3
        mat[0][n],
        mat[1][n],
        mat[2][n],
        mat[3][n],
    )
}

/*multiplies matrix with each coord (edits original vertices)*/
func transform(mat [][]float64, vertices [][][]float64) {
    length := len(vertices)
    for i := 0; i < length; i++ {
        vertices[i] = matMul(mat, vertices[i])
        if vertices[i][3][0] != 1 { // if the forth column if vector isnt 1, then scale the vector
            if round(vertices[i][3][0]*1000) == 0 {panic("divide by 0 in transform()")}
            vertices[i] = matScalar(vertices[i], 1/vertices[i][3][0])//*math.Tan(fov/2))) // is the cot extra??????????????
        }
    }
}

// objects

type cuboid struct{
    coords [][][]float64
    camoords [][][]float64
}

/*creates vertices of cube from given centre and half diagonal vectors.
stores them in cube.coords in a single 4 by 8 matrix
o is centre and u is half diagonal vector*/
func (cu *cuboid) create(o, u [][]float64) {
    cu.camoords = make([][][]float64, 8) // initiallise camoords
    // creating cube parallel to axes by default
    u = [][]float64 { // another method for this is rotating the original u by some angle (pi/2 in case of cubes) in different planes
        {u[0][0], u[0][0], -u[0][0], -u[0][0], u[0][0], u[0][0], -u[0][0], -u[0][0]},
        {u[1][0], u[1][0], u[1][0], u[1][0], -u[1][0], -u[1][0], -u[1][0], -u[1][0]},
        {u[2][0], -u[2][0], -u[2][0], u[2][0], u[2][0], -u[2][0], -u[2][0], u[2][0]},
        {0, 0, 0, 0, 0, 0, 0, 0},
    }
    coords := make([][]float64, 4)
    matAppend(coords, o, o, o, o, o, o, o, o)
    coords = matAdd(coords, u)
    cu.coords = make([][][]float64, 8)
    for i := 0; i < 8; i++ {
        cu.coords[i] = getCoord3d(coords, i)
    }
}

/*draws the cuboid on canvas using camoords*/
func (cu *cuboid) draw(board [][]rune, texture rune) {
    for i := 0; i < 4; i++ { // connecting vertices by lines
        line(cu.camoords[i], cu.camoords[(i+1)%4], board, texture)
        line(cu.camoords[i+4], cu.camoords[4+(i+1)%4], board, texture)
        line(cu.camoords[i], cu.camoords[i+4], board, texture)
    }
}

type triangle struct{
    vertices [][][]float64
    camtices [][][]float64 // the world matrix is multiplied with vertices and is stored here. for methods on triangke
}

func (tri *triangle) create(a, b, c [][]float64) {
    tri.vertices = [][][]float64 {a, b, c}
    tri.camtices = [][][]float64 {a, b, c}
}

/*returns the normal of the triangle in 3d space (tri.vertices).
normal in the direction of the visible face (anticlocck)*/
func (tri *triangle) normal() [][]float64 {
    return vecCross(matSub(tri.vertices[1], tri.vertices[0]), matSub(tri.vertices[2], tri.vertices[0]))
}

/*draws triangle using camtices*/
func (tri *triangle) draw(camPos *[][]float64, board [][]rune, texture rune) {
    if vecDot(matSub(tri.vertices[0], *camPos), tri.normal()) >= 0 {return} // if the front(anticlockwise) face of triangle faces away from/perpendicular to cam, dont draw 
    // >= cuz both vectors have different origin

    for i := 0; i < 3; i++ {
        line(tri.camtices[i], tri.camtices[(i+1)%3], board, texture)
    }
}

/*fills up triangle using camtices*/
func (tri *triangle) fill(camPos *[][]float64, board [][]rune, zbuf [][]float64, texture rune) {
    if vecDot(matSub(tri.vertices[0], *camPos), tri.normal()) >= 0 {return} // if the front(anticlockwise) face of triangle faces away from/perpendicular to cam, dont draw 
    // >= cuz both vectors have different origin
    
    i := 0
    for _, vert := range tri.camtices { // if traingle is outside the screen, dont draw
        if absVal(vert[0][0]) > float64(xlim)/2 || absVal(vert[1][0])*charRatio > float64(ylim)/2 {i++}
    }
    if i == 3 {return}
    
    lightDir := matSub(tri.vertices[0], *camPos)
    // lightDir := [][]float64 {{-1}, {-1}, {-2}, {0}} // from +z to -z
    tex := vecDot(vecUnit(lightDir), vecUnit(tri.normal())) // 0 to 1
    // textures := ".`^,:;Il!i~+_-?][}{!)(|/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
    textures := "':;!vx1WnOM@B"
    tex = -tex*float64(len(textures)-1)
    if tex < 0 {tex = 0}
    texture = rune(textures[round(tex)])
    
    v1 := tri.camtices[0]
    v21 := matSub(tri.camtices[1], v1)
    v31 := matSub(tri.camtices[2], v1)
    dw2, dw3 := 1.0/vecSize(v21), 1.0/vecSize(v31)
    for w2 := 1.0; w2 >= 0; w2 -= dw2 {
        for w3 := 0.0; w3 <= 1-w2; w3 += dw3 {
            poin := nMatAdd(v1, matScalar(v21, w2), matScalar(v31, w3))
            point3d(poin, board, zbuf, texture)
        }
        poin := nMatAdd(v1, matScalar(v21, w2), matScalar(v31, 1-w2)) // just to make sure there are no holes (works pretty well)
        point3d(poin, board, zbuf, texture)
    }
}

type sphere struct {
    triangles []triangle
}

/*creates and joins vertices of a sphere from triangles*/
func (sp *sphere) create(o [][]float64, r float64, n int) {
    vertices := make([][][][]float64, n+1)
    dtheta := math.Pi/float64(n)
    dphi := dtheta*2
    var theta, phi float64
    for j := 0; j < n+1; j++ {
        vertices[j] = make([][][]float64, n)
        for i := 0; i < n; i++ {
            vertices[j][i] = vector(r*math.Sin(theta)*math.Cos(phi), r*math.Cos(theta), r*math.Sin(theta)*math.Sin(phi), 0) // 4th col is 0 cuz o has 1 there
            vertices[j][i] = matAdd(vertices[j][i], o)
            phi += dphi
        }
        theta += dtheta
    }
    sp.triangles = make([]triangle, n*(2*n))
    for j := 0; j < n; j++ {
        for i := 0; i < n; i++ {
            sp.triangles[j*n*2+i*2] = triangle{}
            sp.triangles[j*n*2+i*2].create(vertices[j][i], vertices[j][(i+1)%n], vertices[j+1][i])
            sp.triangles[j*n*2+i*2+1] = triangle{}
            sp.triangles[j*n*2+i*2+1].create(vertices[j+1][i], vertices[j][(i+1)%n], vertices[j+1][(i+1)%n])
        }
    }
}

func (sp *sphere) draw(camPos *[][]float64, cammat [][]float64, board [][]rune, texture rune) {
    for _, tri := range sp.triangles {
        copy(tri.camtices, tri.vertices)
        transform(cammat, tri.camtices)
        tri.draw(camPos, board, texture)
    }
}

func (sp *sphere) fill(camPos *[][]float64, cammat [][]float64, board [][]rune, zbuf [][]float64, texture rune) {
    for _, tri := range sp.triangles {
        copy(tri.camtices, tri.vertices)
        transform(cammat, tri.camtices)
        tri.fill(camPos, board, zbuf, texture)
    }
}

/*multiplies the mat with coords of each triangle*/
func (sp *sphere) transform(mat [][]float64) {
    for _, tri := range sp.triangles {
        transform(mat, tri.vertices)
    }
}

type object struct {
    triangles []triangle
}

func (ob *object) create(path string, o [][]float64) {
    file, err := os.Open(path)
	defer file.Close()
	if err != nil {panic(err)}

	sc := bufio.NewScanner(file) // parsing file
	var vertices [][][]float64
	var faces [][]float64
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
			vertices = append(vertices, vertex)
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
    
    // // finding the centre of a object
    // i := 0
    // cen := vector(0, 0, 0, 0)
    // for _, vertex := range vertices {
    //     i++
    //     cen = matAdd(cen, vertex)
    // }
    // cen = matScalar(cen, 1/float64(i))
    // fmt.Println(cen)
    // //
 	
 	for _, face := range faces { // creating triangles
 	    a := vertices[int(face[0])]
 	    b := vertices[int(face[1])]
 	    c := vertices[int(face[2])]
 	    tri := triangle{}
 	    tri.create(a, b, c)
 	    ob.triangles = append(ob.triangles, tri)
 	}
}

func (ob *object) draw(camPos *[][]float64, cammat [][]float64, board [][]rune, texture rune) {
    for _, tri := range ob.triangles {
        copy(tri.camtices, tri.vertices)
        transform(cammat, tri.camtices)
        tri.draw(camPos, board, texture)
    }
}

func (ob *object) fill(camPos *[][]float64, cammat [][]float64, board [][]rune, zbuf [][]float64, texture rune) {
    for _, tri := range ob.triangles {
        copy(tri.camtices, tri.vertices)
        transform(cammat, tri.camtices)
        tri.fill(camPos, board, zbuf, texture)
    }
}

/*multiplies the mat with coords of each triangle*/
func (ob *object) transform(mat [][]float64) {
    for _, tri := range ob.triangles {
        transform(mat, tri.vertices)
    }
}