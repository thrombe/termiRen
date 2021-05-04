package main

import (
	// "time"
	// "fmt"
	"math"
	"seehuhn.de/go/ncurses"
)

// REMEMBER TO USE POINTERS and fix printB()
// KEEP Z COORD -ve (camera is facing -z)
const xlim = 151
const ylim = 163
var fov = math.Pi/2.5 //*horizontal // keep this between 0 and pi
var charRatio = 1.4/2.55 // used only in point() // width/height of a character
const ncursed = 1 // print using ncurses if 1 else just fmt.print

func main() {
    if ncursed == 1 {defer ncurses.EndWin()}
    demo5()
    // return
    // win := ncurses.Init()
    // a := string(win.GetCh())
    // win.AddStr(a+ a+ a+ a+ a+ a)
    // win.Println(a, a, a, a)
    // win.GetCh()
    // win.Refresh()
}

func demo5() { // rotating cube 3d with a cam
	o := [][]float64 {{5}, {5}, {-30}, {1}} // 1 for 4 by 1 matrix
	u := [][]float64 {{5}, {5}, {5}, {0}} // 0 dosent matter here
	axis := [][]float64 {{1}, {1}, {-1}, {0}}
	rot := rotAboutVec(0.2, axis)
	//rot := rotMat3dy(0.2)
	rot = rotAboutPoint(rot, o)
	b := cuboid{}
	b.create(o, u)
	board := genB()
	win, printy := perint(board)

	matchan := make(chan [][]float64)
	camPos, camDir := camInit()
	go detectKey(&camPos, &camDir, win, matchan)
	rotD := transMat(camPos)


	for i := 0; i < 500; i++ {
		// time.Sleep(time.Millisecond*50)/////////
		select {
		case mat := <- matchan:
			rotD = mat//matMul(mat, rotD)
			// win.Println(0, 1, camPos)/////////
			// win.Refresh()///////////	
		default:
		}
		b.coords = matMul(rot, b.coords)
		coords := matMul(rotD, b.coords)
		b.draw(board, coords)
		printy()
		// _ = printy
	}
}

/*
func demo4() { // morphing cube 3d
	o := [][]float64 {{0}, {0}, {30}, {1}} // 1 for 4 by 1 matrix
	u := [][]float64 {{5}, {5}, {5}, {0}} // 0 dosent matter
	//centre := [][]float64 {{4}, {15}, {30}}
	rot := matMul(rotMat3dx(0.1), rotMat3dy(0.3))
	rot = matMul(rot, rotMat3dz(0.2))
	board := genB()
	_, printy := perint(board)
	for i := 0; i < 200; i++ {
		u = matMul(rot, u)
		b := cuboid{}
    	b.create(o, u)
		b.draw(board)
		printy()
	}
}

func demo3() { // rotating cube 3d
	o := [][]float64 {{5}, {5}, {-30}, {1}} // 1 for 4 by 1 matrix
	u := [][]float64 {{5}, {5}, {5}, {0}} // 0 dosent matter here
	axis := [][]float64 {{1}, {1}, {-1}, {0}}
	rot := rotAboutVec(0.2, axis)
	//rot := rotMat3dy(0.2)
	rot = rotAboutPoint(rot, o)
	b := cuboid{}
	b.create(o, u)
	board := genB()
	_, printy := perint(board)
	for i := 0; i < 200; i++ {
		b.coords = matMul(rot, b.coords)
		b.draw(board)
		printy()
	}
}

func demo2() { // 2 moving squares 3d
	fov = math.Pi/2.7
	var x float64 = -20
	var y float64 = 20
	var z float64 = 70
	board := genB()
	_, printy := perint(board)
	for i := -15; i <= 15; i++ {
		x = float64(i)
		p := [][]float64 {
			{0+x},
			{0+y},
			{z},
			{1}, // 1 for 4 by 1 matrix
		}
		point(0, 0, board)
		// p[2][0] = z
		xysquare3d(p, 15, board)
		p[2][0] = z+15
		xysquare3d(p, 15, board)
		printy()
	}
}

func demo() { // 2d rotating line
    fov = math.Pi/1.15
    var x1, y1, x2, y2 float64 = 5, 6, 28, 31
    var rx, ry float64 = 5, -3
    var r float64 = 0.2
	board := genB()
	_, printy := perint(board)
    for i:= 0; i < 250; i++ {
        line(x1, y1, x2, y2, board)
        point(rx, ry, board)
		printy()
        x1, y1 = rotateP2d(x1, y1, rx, ry, r)
        x2, y2 = rotateP2d(x2, y2, rx, ry, r)
    }
}
*/