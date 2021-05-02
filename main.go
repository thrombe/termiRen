package main

import (
	"time"
	//"fmt"
	"math"
)

// REMEMBER TO USE POINTERS and fix printB()
// KEEP Z COORD -ve (camera is facing -z)
const xlim = 151
const ylim = 163
var fov = math.Pi/2.5 //*horizontal // keep this between 0 and pi
var charRatio = 1.4/2.55 // used only in point() // width/height of a character

func main() {
    demo3()
}

func demo4() { // morphing cube 3d
	o := [][]float64 {{0}, {0}, {30}, {1}} // 1 for 4 by 1 matrix
	u := [][]float64 {{5}, {5}, {5}, {0}} // 0 dosent matter
	//centre := [][]float64 {{4}, {15}, {30}}
	rot := matMul(rotMat3dx(0.1), rotMat3dy(0.3))
	rot = matMul(rot, rotMat3dz(0.2))
	for {
		u = matMul(rot, u)
		b := cuboid{}
    	b.create(o, u)
		board := genB()
		b.draw(board)
		printB(board)
	}
}

func demo3() { // rotating cube 3d
	o := [][]float64 {{0}, {0}, {-30}, {1}} // 1 for 4 by 1 matrix
	u := [][]float64 {{5}, {5}, {5}, {0}} // 0 dosent matter here
	axis := [][]float64 {{9}, {9}, {-9}, {0}}
	rot := rotAboutVec(0.2, axis)
	rot = rotAboutPoint(rot, o)
	b := cuboid{}
	b.create(o, u)
	for {
		b.coords = matMul(rot, b.coords)
		board := genB()
		b.draw(board)
		printB(board)
	}
}

func demo2() { // 2 moving squares 3d
	fov = math.Pi/2.7
	var x float64 = -20
	var y float64 = 20
	var z float64 = 70
	for i := -15; i <= 15; i++ {
		x = float64(i)
		p := [][]float64 {
			{0+x},
			{0+y},
			{z},
			{1}, // 1 for 4 by 1 matrix
		}
		board := genB()
		point(0, 0, board)
		// p[2][0] = z
		xysquare3d(p, 15, board)
		p[2][0] = z+15
		xysquare3d(p, 15, board)
		printB(board)
	}
}

func demo() { // 2d rotating line
    fov = math.Pi/1.15
    var x1, y1, x2, y2 float64 = 5, 6, 28, 31
    var rx, ry float64 = 5, -3
    var r float64 = 0.2
    for i:= 0; i < 250; i++ {
        board := genB()
        line(x1, y1, x2, y2, board)
        point(rx, ry, board)
        printB(board)
        time.Sleep(time.Millisecond*100)
        x1, y1 = rotateP2d(x1, y1, rx, ry, r)
        x2, y2 = rotateP2d(x2, y2, rx, ry, r)
    }
}
