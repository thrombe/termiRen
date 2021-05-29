package main

import (
	"time"
	// "fmt"
	"math"
	"seehuhn.de/go/ncurses"
	// "github.com/pkg/profile"
)

const xlim = 151
const ylim = 163
const fov = math.Pi/2.5 //*horizontal // keep this between 0 and pi
const charRatio = 1.4/2.55 // used only in point() // width/height of a character
const ncursed = 1 // print using ncurses if 1 else just fmt.print

const tranV = 4.0 // camera movement velocity (metres?)
const rotA = 0.1 // camera rotation angular vemocity (radians)

// CONTROLS - wasd to move, ijkl to turn camera, WS to move up and down, q to QUIT

func main() {
	// defer profile.Start(profile.MemProfile).Stop()
	// defer profile.Start().Stop()
    if ncursed == 1 {defer ncurses.EndWin()}
    demo5()
}

func demo8() { // load obj files
	o := vector(0, 0, -10, 1) // offset the object by that vector
	obj := object{}
	obj.create("./objects/teapot.obj", o) // 47.2k triangles 23.36k vertices in big chungus
	rawboard, board, zbuf := genB()
	win, printy:= perint(rawboard, board, zbuf)
	
	// axis := [][]float64 {{1}, {1}, {-1}, {0}} // rotations
	// rot := rotAboutVec(0.2, axis)
	rot := rotMat3dy(0.2)
	rot = rotAboutPoint(rot, obj.center)

	// big := scaleMat(5, 4) // scale object size
	// big = rotAboutPoint(big, obj.center)
	// obj.transform(big)

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
	sp := sphere(o, 5, 200)
	// sp.create // performance goal(for now)(without multicore) - should be smooth for 200 at just the draw
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
	// rot := rotMat3dy(0.2)
	rot = rotAboutPoint(rot, o)
	b := cuboid(o, u)
	// b.create(o, u)
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
		// b.draw(&camPos, camMat, board, '.')
		b.fill(&camPos, camMat, board, zbuf, '#')
		printy()
		time.Sleep(time.Millisecond*50)
	}
}
