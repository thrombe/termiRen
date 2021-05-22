package main

import (
	// "fmt"
	"seehuhn.de/go/ncurses"
)

/*creates a default position and direction for the camera*/
func camInit() ([][]float64, [][]float64) {
	camPos := vector(0, 0, 0, 1)
	camDir := vector(0, 0, -tranV, 0)
	return camPos, camDir
}

/*detects keys and sends curresponding matrix to the channel*/
func detectKey(camPos, camDir *[][]float64, win *ncurses.Window, matchan chan [][]float64, quitchan chan bool) {
	mat := transMat(*camPos)
	up := vector(0, tranV, 0, 0)
	hor := vector(tranV, 0, 0, 0)
	fwd := vector(0, 0, -tranV, 0)
	camup := vector(0, tranV, 0, 0)
	camhor := vector(tranV, 0, 0, 0)
	proj := projectionMat()
	for {
		kee := win.GetCh()
		switch kee { // keep piling up rotations and translations on top(left matrix multiplication) of the previous ones. so i dont need to keet track of its current position and face direction for movement. though that would be helpful later
		case 'w': // move forward
			mat = matMul(transMatInv(fwd), mat)
			*camPos = matMul(transMat(*camDir), *camPos)
		case 's': // move backward
			mat = matMul(transMat(fwd), mat)
			*camPos = matMul(transMatInv(*camDir), *camPos)
		case 'a': // move left
			mat = matMul(transMat(hor), mat)
			*camPos = matMul(transMatInv(camhor), *camPos)
		case 'd': // move right
			mat = matMul(transMatInv(hor), mat)
			*camPos = matMul(transMat(camhor), *camPos)
		case 'W': // move up
			mat = matMul(transMatInv(up), mat)
			*camPos = matMul(transMat(camup), *camPos)
		case 'S': // move down
			mat = matMul(transMat(up), mat)
			*camPos = matMul(transMatInv(camup), *camPos)
		case 'j': // look left
			mat = matMul(rotAboutVec(-rotA, up), mat)
			*camDir = matMul(rotAboutVec(rotA, camup), *camDir)
			camhor = matMul(rotAboutVec(rotA, camup), camhor)
		case 'l': // look right
			mat = matMul(rotAboutVec(rotA, up), mat)
			*camDir = matMul(rotAboutVec(-rotA, camup), *camDir)
			camhor = matMul(rotAboutVec(-rotA, camup), camhor)
		case 'i': // look up
			mat = matMul(rotAboutVec(-rotA, hor), mat)
			*camDir = matMul(rotAboutVec(rotA, camhor), *camDir)
			camup = matMul(rotAboutVec(rotA, camhor), camup)
		case 'k': // look down
			mat = matMul(rotAboutVec(rotA, hor), mat)
			*camDir = matMul(rotAboutVec(-rotA, camhor), *camDir)
			camup = matMul(rotAboutVec(-rotA, camhor), camup)
		case 'q': // quit
			quitchan <- true
			// default: // default not needed. so any key presses other than controls are ignored
			// win.Println(*camPos, *camDir)/////////
		}
		// fmt.Println(camup, camhor)
		matchan <- matMul(proj, mat)
	}
}
