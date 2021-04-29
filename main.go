package main
import (
	"time"
	//"fmt"
	"math"
)
// REMEMBER TO USE POINTERS and fix printB()
// KEEP Z COORD +VE
const xlim = 151
const ylim = 163
var fov = math.Pi/2.5 //*horizontal // keep this between 0 and pi
var charRatio = 1.4/2.55 // used only in point() // width/height of a character

func main() {
	demo2()
	//fmt.Println(rotateP3d(1, 0, 0))
}

func demo2() {
	fov = math.Pi/2.7
	var x float64 = -20
	var y float64 = 20
	var z float64 = 70
	for i := -15; i <= 15; i++ {
		x = float64(i)
		p := [][]float64 {
			{0+x},
			{0+y},
			{z},
		}
		board := genB()
		point(0, 0, board)
		// p[2][0] = z
		xysquare3d(p, 15, board)
		p[2][0] = z+15
		xysquare3d(p, 15, board)
		printB(board)
	}
}

func demo() {
    fov = math.Pi/1.15
    var x1, y1, x2, y2 float64 = 5, 6, 28, 31
    var rx, ry float64 = 5, -3
    var r float64 = 0.2
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
