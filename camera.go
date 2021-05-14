package main

import (
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
	camup := [][]float64 {{0}, {tranV}, {0}, {0}}
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
			camhor := matScalar(vecCross(*camDir, camup), 1/tranV) // hor == cross of fwd and up
			*camPos = matMul(transMatInv(camhor), *camPos)
		case 'd': // move right
			mat = matMul(transMatInv(hor), mat)
			camhor := matScalar(vecCross(*camDir, camup), 1/tranV)
			*camPos = matMul(transMat(camhor), *camPos)
		case 'W': // move up
			mat = matMul(transMatInv(up), mat)
			*camPos = matMul(transMat(camup), *camPos)
		case 'S': // move down
			mat = matMul(transMat(up), mat)
			*camPos = matMul(transMatInv(camup), *camPos)
		case 'j': // look left
			mat = matMul(rotAboutVec(-rotA, up), mat)
			*camDir = matMul(rotAboutPoint(rotAboutVec(rotA, camup), *camPos), *camDir)
		case 'l': // look right
			mat = matMul(rotAboutVec(rotA, up), mat)
			*camDir = matMul(rotAboutPoint(rotAboutVec(-rotA, camup), *camPos), *camDir)
		case 'i': // look up
			mat = matMul(rotAboutVec(-rotA, hor), mat)
			camhor := matScalar(vecCross(*camDir, camup), 1/tranV)
			*camDir = matMul(rotAboutPoint(rotAboutVec(rotA, camhor), *camPos), *camDir)
			camup = matMul(rotAboutPoint(rotAboutVec(rotA, camhor), *camPos), camup)
		case 'k': // look down
			mat = matMul(rotAboutVec(rotA, hor), mat)
			camhor := matScalar(vecCross(*camDir, camup), 1/tranV)
			*camDir = matMul(rotAboutPoint(rotAboutVec(-rotA, camhor), *camPos), *camDir)
			camup = matMul(rotAboutPoint(rotAboutVec(-rotA, camhor), *camPos), camup)
		case 'q': // quit
			quitchan <- true
		// default: // default not needed. so any key presses other than controls are ignored
		// win.Println(*camPos, *camDir)/////////
	}
		matchan <- matMul(proj, mat)
	}
}
