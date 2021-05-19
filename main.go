package main

import (
	"time"
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

var tranV = 4.0 // camera movement velocity (metres?)
var rotA = 0.1 // camera rotation angular vemocity (radians)

// CONTROLS - wasd to move, ijkl to turn camera, WS to move up and down, q to QUIT

var blank = ' '
// var unblank = '.'

func main() {
    if ncursed == 1 {defer ncurses.EndWin()}
    demo7()
}

func demo7() { // sphere with n*n triangles
	o := [][]float64 {{0}, {0}, {-30}, {1}}
	// axis := [][]float64 {{1}, {1}, {-1}, {0}}
	// rot := rotAboutVec(0.2, axis)
	rot := rotMat3dy(0.0)
	rot = rotAboutPoint(rot, o)
	sp := sphere{}
	sp.create(o, 5, 20)
	rawboard, board, zbuf := genB()
	win, printy:= perint(rawboard, board, zbuf)

	matchan := make(chan [][]float64)
	quitchan := make(chan bool)
	camPos, camDir := camInit()
	go detectKey(&camPos, &camDir, win, matchan, quitchan)
	camMat := matMul(projectionMat(), transMat(camPos)) // default value

	loop: for { // press q to quit
		select {
		case camMat = <- matchan:
		case <- quitchan: break loop
		default:
		}
		sp.transform(rot)

		// sp.draw(&camPos, camMat, board, '.')
		// time.Sleep(time.Millisecond*50)
		sp.fill(&camPos, camMat, board, zbuf, '#')
		printy()
	}
}

func demo6() { // testing the tringle. if it is drawn when it should not be. (it faces away from camera)
	o := [][]float64 {{0}, {0}, {-30}, {1}}
	a := [][]float64 {{5}, {5}, {-30}, {1}}
	b := [][]float64 {{-5}, {5}, {-30}, {1}}
	c := [][]float64 {{0}, {-5}, {-30}, {1}}
	// axis := [][]float64 {{1}, {1}, {-1}, {0}}
	// rot := rotAboutVec(0.2, axis)
	rot := rotMat3dy(0.0)
	rot = rotAboutPoint(rot, o)
	t := triangle{}
	t.create(a, b, c)
	rawboard, board, zbuf := genB()
	win, printy := perint(rawboard, board, zbuf)

	matchan := make(chan [][]float64)
	quitchan := make(chan bool)
	camPos, camDir := camInit()
	go detectKey(&camPos, &camDir, win, matchan, quitchan)
	camMat := matMul(projectionMat(), transMat(camPos)) // default value

	loop: for { // press q to quit
		select {
		case camMat = <- matchan:
		case <- quitchan: break loop
		default:
		}
		transform(rot, t.vertices)
		copy(t.camtices, t.vertices)

		transform(camMat, t.camtices)
		// t.draw(board)
		t.fill(&camPos, board, zbuf, '#')
		// poin := matMul(camMat, camPos)
		// vec := matMul(camMat, camDir)
		// vector(poin, vec, board)
		printy()
	}
}

func demo5() { // rotating cube 3d with a cam
	o := [][]float64 {{5}, {5}, {-30}, {1}} // 1 for 4 by 1 matrix
	u := [][]float64 {{5}, {5}, {5}, {0}} // 0 dosent matter here
	// axis := [][]float64 {{1}, {1}, {-1}, {0}}
	// rot := rotAboutVec(0.2, axis)
	rot := rotMat3dy(0.0)
	rot = rotAboutPoint(rot, o)
	b := cuboid{}
	b.create(o, u)
	rawboard, board, zbuf := genB()
	win, printy := perint(rawboard, board, zbuf)

	matchan := make(chan [][]float64)
	quitchan := make(chan bool)
	camPos, camDir := camInit()
	go detectKey(&camPos, &camDir, win, matchan, quitchan)
	camMat := matMul(projectionMat(), transMat(camPos)) // default value

	loop: for { // press q to quit
		select {
		case camMat = <- matchan:
		case <- quitchan: break loop
		default:
		}
		transform(rot, b.coords)
		copy(b.camoords, b.coords)
		transform(camMat, b.camoords)
		b.draw(board, '.')
		printy()
		time.Sleep(time.Millisecond*50)
	}
}