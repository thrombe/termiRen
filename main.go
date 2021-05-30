package main

import (
	"time"
	"fmt"
	"math"
	"seehuhn.de/go/ncurses"
	// "github.com/pkg/profile"
)

//creates a default position and direction for the camera
func camInit() *camera {
	cam := camera{}

	// can play with these
	cam.xlim = 151 // screen width in characters
	cam.ylim = 163 // screen height in characters
	cam.fov = math.Pi/2.5 //*horizontal // keep this between 0 and pi
	cam.charRatio = 1.4/2.55 // width/height of a character
	cam.tranV = 4.0 // camera movement velocity (metre/s?)
	cam.rotA = 0.1 // camera rotation angular vemocity (radians)

	// default values
	cam.camPos = vector(0, 0, 0, 1)
	cam.camDir = vector(0, 0, -cam.tranV, 0)
	cam.camMat = matMul(projectionMat(&cam), transMat(cam.camPos))
	
	cam.ncursed = true // print using ncurses if true else just fmt.print
	
	return &cam
}
// CONTROLS - wasd to move, ijkl to turn camera, WS to move up and down, q to QUIT

func main() {
	// defer profile.Start(profile.MemProfile).Stop()
	// defer profile.Start().Stop()
    a := demo8(true)
	if (camInit()).ncursed {ncurses.EndWin()}
	fmt.Println(a)
}

func demo8(taime bool) time.Duration { // load obj files
	o := vector(0, 0, -10, 1) // offset the object by that vector

	// create/import objects
	obj := importObj("./objects/teapot.obj", o) // import .obj files
	// obj := sphere(o, 5, 50) // sphere with r as radius and n*n triangles
	// obj := cuboid(o, vector(2, 2, 2, 0))
	
	// rotate objects per frame
	// rot := rotAboutVec(0.2, vector(1, 1, -1, 0))
	rot := rotMat3dy(0.2)
	rot = rotAboutPoint(rot, obj.center)

	matchan := make(chan [][]float64)
	quitchan := make(chan bool)
	cam := camInit()
	genB(cam)
	win, printy:= perint(cam)
	go detectKey(cam, win, matchan, quitchan)

	start := time.Now()
	loop: for i := 0; (i < 200 || !taime); i++ { // press q to quit
		select {
		case cam.camMat = <- matchan:
		case <- quitchan: break loop
		default:
		}
		obj.transform(rot)

		// obj.draw(cam, '.')
		// time.Sleep(time.Millisecond*50)
		obj.fill(cam)
		printy()
	}
	return time.Now().Sub(start)
}
