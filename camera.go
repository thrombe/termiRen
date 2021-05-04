package main

import (
	// "time"
	"seehuhn.de/go/ncurses"
)

var tranV = float64(4)
var rotA = 0.1

/*creates a default position and direction for the camera*/
func camInit() ([][]float64, [][]float64) {
	camPos := [][]float64 {{0}, {0}, {0}, {1}}
	camDir := [][]float64 {{0}, {0}, {-tranV}, {0}}
	return camPos, camDir
}

/*detects keys and sends curresponding matrix to the channel*/
func detectKey(camPos, camDir *[][]float64, win *ncurses.Window, matchan chan [][]float64) {
	trick := [][]float64 {{0}, {0}, {0}, {1}}
	mat := transMat(trick)
	
	for {
		kee := win.GetCh()
		switch kee {
		case 'w':
			//*camPos = matSub(*camPos, *camDir)
			// win.Println(*camPos)/////////
			// win.Refresh()///////////	
			//mat = matMul(transMat(*camPos), mat)
			mat = matMul(transMatInv(*camDir), mat)
			*camPos = matMul(mat, *camPos)
			matchan <- mat
		case 's':
			//*camPos = matAdd(*camPos, *camDir)
			//mat = matMul(transMat(*camPos), mat)
			mat = matMul(transMat(*camDir), mat)
			*camPos = matMul(mat, *camPos)
			matchan <- mat
		case 'a':
			up := [][]float64 {{0}, {1}, {0}, {0}}
			mat = matMul(rotAboutVec(-rotA, up), mat)
			// angle += rotA
			// *camDir = matMul(mat, *camDir)
			//mat = rotAboutPoint(mat, *camPos) // i dont understand why this isnt needed
			matchan <- mat
		case 'd':
			up := [][]float64 {{0}, {1}, {0}, {0}}
			mat = matMul(rotAboutVec(rotA, up), mat)
			// angle += -rotA
			// *camDir = matMul(mat, *camDir)
			//mat = rotAboutPoint(mat, *camPos)
			matchan <- mat
		default:
		}
	}
}
