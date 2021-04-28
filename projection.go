package main
import (
    "fmt"
    "math"
    )

func main() {
    // board := genB()
    // point(0, 0, board)
    // x, y := rotateP(10, 10, 0, 0, 0.0)
    // point(x, y, board)
    // point(-10, -10, board )
    // vector(1, 2, 30, 11, board)
    // printB(board)
    mat1 := [][]float64 {
        {3, 1, 4, 6}, 
        {1, 4, 6, 3}, 
        {6, 3, 1, 9},
        {2, 8, 5, 8} }
    // fmt.Println(mat1[0][0])
    new := subMat(mat1, 2, 3)
    fmt.Println(new)
}

const xlim = 151
const ylim = 160

/*rotate x, y about ox, oy by r radians*/
func rotateP(x, y, ox, oy, r float64) (float64, float64) {
    x, y = x-ox, y-oy
    x, y =  x*math.Cos(r)-y*math.Sin(r), x*math.Sin(r)+y*math.Cos(r)
    return x+ox, y+oy
}

func printB(board []int) {
    var scr string
    for y := 0; y < ylim; y++ {
        for x := 0; x < xlim; x++ {
            if board[x+y*xlim] == 1 {
                scr += "x"
            } else {scr += "."}
        }
        scr += "\n"
    }
    fmt.Println(scr)
}

/*draw line on canvas*/
func line(x1, y1, x2, y2 float64, board []int) {
    length := math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
    vx, vy := float64(x2-x1), float64(y2-y1)
    ux, uy := vx/length,  vy/length
    x, y, x3 := float64(x1), float64(y1), float64(x2)
    for x <= x3 {
        point(x, y, board)
        x, y = x+ux, y+uy
    }
}

/*draw vector with x1, y1 as offset and x2, y2 as direction*/
func vector(x1, y1, x2, y2 float64, board []int) {
    line(x1, y1, x1+x2, y1+y2, board)
}

/*draw point on canvas*/
func point(h, k float64, board []int) {
    x, y := giveInd(h, k)
    board[x+y*xlim] = 1
}

/*make a canvas. might be useful to use float to be able to have brightness and stuff*/
func genB() []int {
    return make([]int, xlim*ylim)
    // var board [][]int
    // for y := 0; y < ylim; y++ {
    //     var temp []int
    //     for x := 0; x < xlim; x++ {
    //         temp = append(temp, 0)
    //     }
    // }
}

/*enter index of pixel in canvas and return coors in space*/
func giveCoords(x, y int) (float64, float64) {
    return float64(x-int(xlim/2)), float64(-y+int(ylim/2)) 
}

/*enter coords of points and return index in canvas*/
func giveInd(x, y float64) (int, int) {
    return round(x)+int(xlim/2), -round(y)+int(ylim/2)
}

