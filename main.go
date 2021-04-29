package main
import (
	"time"
	// "fmt"
)
// REMEMBER TO USE POINTERS and fix line() and printB()
const xlim = 151
const ylim = 160
var distanceScaleFactor = 1.0 // scales according to this in z axis

func main() {
	demo2()
}

func demo2() {
	distanceScaleFactor = 0.3
	var x float64 = -20
	var y float64 = 0
	var z float64 = 2
	for i := -20; i <= 20; i++ {
		x = float64(i)
		p := [][]float64 {
			{0+x},
			{20+y},
			{z},
		}
		board := genB()
		point(0, 0, board)
		// p[2][0] = z
		xysquare3d(p, 20, board)
		p[2][0] = 4
		xysquare3d(p, 20, board)
		printB(board)
	}
}

func demo() {
    var x1, y1, x2, y2 float64 = 5, 6, 28, 31
    var rx, ry float64 = 5, -3
    var r float64 = 0.03
    for i:= 0; i < 250; i++ {
        board := genB()
        line(x1, y1, x2, y2, board)
        point(rx, ry, board)
        printB(board)
        time.Sleep(time.Millisecond*100)
        x1, y1 = rotateP2d(x1, y1, rx, ry, r)
        x2, y2 = rotateP2d(x2, y2, rx, ry, r)
    }
}
