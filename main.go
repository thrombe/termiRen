package main

import (
	"time"
	// "fmt"
	"math"
	"seehuhn.de/go/ncurses"
	// "github.com/pkg/profile"
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

var blank byte = ' '
// var unblank = '.'

func main() {
	// defer profile.Start(profile.MemProfile).Stop()
	// defer profile.Start().Stop()
    if ncursed == 1 {defer ncurses.EndWin()}
    demo8()
}

func demo8() { // load obj files
	o := vector(0, 0, -10, 1)
	// axis := [][]float64 {{1}, {1}, {-1}, {0}}
	// rot := rotAboutVec(0.2, axis)
	rot := rotMat3dy(0.0)
	rot = rotAboutPoint(rot, o)
	obj := object{}
	obj.create("./objects/dragon2.obj", o) // 47.2k triangles 23.36k vertices in big chungus
	rawboard, board, zbuf := genB()
	win, printy:= perint(rawboard, board, zbuf)

	big := scaleMat(5, 4) // scale object size
	// big = rotAboutPoint(big, o)
	obj.transform(big)

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
		obj.transform(rot)

		// obj.draw(&camPos, camMat, board, '.')
		// time.Sleep(time.Millisecond*50)
		obj.fill(&camPos, camMat, board, zbuf, '#')
		printy()
	}
}


func demo7() { // sphere with n*n triangles
	o := [][]float64 {{0}, {0}, {-30}, {1}}
	// axis := [][]float64 {{1}, {1}, {-1}, {0}}
	// rot := rotAboutVec(0.2, axis)
	rot := rotMat3dy(0.0)
	rot = rotAboutPoint(rot, o)
	sp := sphere{}
	sp.create(o, 5, 200) // performance goal(for now)(without multicore) - should be smooth for 200 at just the draw
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

		sp.draw(&camPos, camMat, board, '.')
		// time.Sleep(time.Millisecond*50)
		// sp.fill(&camPos, camMat, board, zbuf, '#')
		printy()
	}
}

func demo5() { // rotating cube 3d with a cam
	o := [][]float64 {{0}, {5}, {-30}, {1}} // 1 for 4 by 1 matrix
	u := [][]float64 {{5}, {5}, {5}, {0}} // 0 dosent matter here
	axis := [][]float64 {{1}, {1}, {-1}, {0}}
	rot := rotAboutVec(0.2, axis)
	// rot := rotMat3dy(0.0)
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
		b.transform(rot)
		b.draw(camMat, board, '.')
		printy()
		time.Sleep(time.Millisecond*50)
	}
}