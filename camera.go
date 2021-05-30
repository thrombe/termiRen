package main

import (
	// "fmt"
	"seehuhn.de/go/ncurses"
)

type camera struct {
	camPos, camDir [][]float64
	xlim, ylim int
	tranV, rotA, fov, charRatio float64
	
	ncursed bool
}

//detects keys and sends curresponding matrix to the channel
func detectKey(cam *camera, win *ncurses.Window, matchan chan [][]float64, quitchan chan bool) {
	mat := transMat(cam.camPos)
	up := vector(0, cam.tranV, 0, 0)
	hor := vector(cam.tranV, 0, 0, 0)
	fwd := vector(0, 0, -cam.tranV, 0)
	camup := vector(0, cam.tranV, 0, 0)
	camhor := vector(cam.tranV, 0, 0, 0)
	proj := projectionMat(cam)
	for {
		kee := win.GetCh()
		switch kee { // keep piling up rotations and translations on top(left matrix multiplication) of the previous ones. so i dont need to keet track of its current position and face direction for movement. though that would be helpful later
		case 'w': // move forward
			mat = matMul(transMatInv(fwd), mat)
			cam.camPos = matMul(transMat(cam.camDir), cam.camPos)
		case 's': // move backward
			mat = matMul(transMat(fwd), mat)
			cam.camPos = matMul(transMatInv(cam.camDir), cam.camPos)
		case 'a': // move left
			mat = matMul(transMat(hor), mat)
			cam.camPos = matMul(transMatInv(camhor), cam.camPos)
		case 'd': // move right
			mat = matMul(transMatInv(hor), mat)
			cam.camPos = matMul(transMat(camhor), cam.camPos)
		case 'W': // move up
			mat = matMul(transMatInv(up), mat)
			cam.camPos = matMul(transMat(camup), cam.camPos)
		case 'S': // move down
			mat = matMul(transMat(up), mat)
			cam.camPos = matMul(transMatInv(camup), cam.camPos)
		case 'j': // look left
			mat = matMul(rotAboutVec(-cam.rotA, up), mat)
			cam.camDir = matMul(rotAboutVec(cam.rotA, camup), cam.camDir)
			camhor = matMul(rotAboutVec(cam.rotA, camup), camhor)
		case 'l': // look right
			mat = matMul(rotAboutVec(cam.rotA, up), mat)
			cam.camDir = matMul(rotAboutVec(-cam.rotA, camup), cam.camDir)
			camhor = matMul(rotAboutVec(-cam.rotA, camup), camhor)
		case 'i': // look up
			mat = matMul(rotAboutVec(-cam.rotA, hor), mat)
			cam.camDir = matMul(rotAboutVec(cam.rotA, camhor), cam.camDir)
			camup = matMul(rotAboutVec(cam.rotA, camhor), camup)
		case 'k': // look down
			mat = matMul(rotAboutVec(cam.rotA, hor), mat)
			cam.camDir = matMul(rotAboutVec(-cam.rotA, camhor), cam.camDir)
			camup = matMul(rotAboutVec(-cam.rotA, camhor), camup)
		case 'q': // quit
			quitchan <- true
			// default: // default not needed. so any key presses other than controls are ignored
			// win.Println(cam.camPos, cam.camDir)/////////
		}
		// fmt.Println(camup, camhor)
		matchan <- matMul(proj, mat)
	}
}
