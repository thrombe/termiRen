package main

import (
	// "time"
	"seehuhn.de/go/ncurses"
)

/*creates a default position and direction for the camera*/
func camInit() ([][]float64, [][]float64) {
	camPos := [][]float64 {{0}, {0}, {0}, {1}}
	camDir := [][]float64 {{0}, {0}, {-tranV}, {0}}
	return camPos, camDir
}

/*detects keys and sends curresponding matrix to the channel*/
func detectKey(camPos, camDir *[][]float64, win *ncurses.Window, matchan chan [][]float64, quitchan chan bool) {
	mat := transMat(*camPos) // camPos and camDir arent really needed for this. but will have to keep track of them here to find camera's in 3d position
	up := [][]float64 {{0}, {tranV}, {0}, {0}}
	hor := [][]float64 {{tranV}, {0}, {0}, {0}}
	fwd := [][]float64 {{0}, {0}, {-tranV}, {0}}
	for {
		kee := win.GetCh()
		switch kee { // keep piling up rotations and translations on top(left matrix multiplication) of the previous ones. so i dont need to keet track of its current position and face direction for movement. though that would be helpful later
		case 'w': // move forward
			mat = matMul(transMatInv(fwd), mat)
		case 's': // move backward
			mat = matMul(transMat(fwd), mat)
		case 'a': // move left
			mat = matMul(transMat(hor), mat)
		case 'd': // move right
			mat = matMul(transMatInv(hor), mat)
		case 'W': // move up
			mat = matMul(transMatInv(up), mat)
		case 'S': // move down
			mat = matMul(transMat(up), mat)
		case 'j': // look left
			mat = matMul(rotAboutVec(-rotA, up), mat)
		case 'l': // look right
			mat = matMul(rotAboutVec(rotA, up), mat)
		case 'i': // look up
			mat = matMul(rotAboutVec(-rotA, hor), mat)
		case 'k': // look down
			mat = matMul(rotAboutVec(rotA, hor), mat)
		case 'q': // quit
			quitchan <- true
		// default: // default not needed. so any key presses other than controls are ignored
			// win.Println(up)/////////
			// win.Refresh()///////////	
		}
		matchan <- mat
	}
}
