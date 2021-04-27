package main
import (
    "fmt"
    "math"
    )

func main() {
    board := genB()
    point(0, 0, board)
    x, y := rotateP(10, 10, 0, 0, 3.1 )
    point(round(x), round(y), board)
    // point(-10, -10, board )
    vector(1, 2, 30, 11, board)
    printB(board)
}

const xlim = 151
const ylim = 160

func rotateP(x, y, ox, oy, r float64) (float64, float64) {
    x, y = x-ox, y-oy
    x, y =  x*math.Cos(r)-y*math.Sin(r), x*math.Sin(r)+y*math.Cos(r)
    return x+ox, y+oy
}

func printB(board []int) {
    //xlim, ylim := 150, 160
    // nxlim, nylim := -int(xlim/2), -int(ylim/2)
    // pxlim, pylim := xlim+nxlim, ylim+nylim
    var scr string
    // for y := nylim; y <= pylim; y++ {
        // for x := nxlim; x <= pxlim; x++ {
    for y := 0; y < ylim; y++ {
        for x := 0; x < xlim; x++ {
            if board[x+y*xlim] == 1 {
                scr += "x"
            } else {scr += "."}
        }
        scr += "\n"
    }
    fmt.Printf(scr)
}

func line(x1 int, y1 int, x2 int, y2 int, board []int) {
    length := math.Sqrt(float64((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1)))
    vx, vy := float64(x2-x1), float64(y2-y1)
    ux, uy := vx/length,  vy/length
    x, y, x3 := float64(x1), float64(y1), float64(x2)
    for x <= x3 {
        //x, y = math.Round(x), math.Round(y)
        point(round(x), round(y), board)
        x, y = x+ux, y+uy
    }
}

func matMul() {
    fmt.Printf("matMul not implemented")
}

func vector(x1 int, y1 int, x2 int, y2 int, board []int) {
    line(x1, y1, x1+x2, y1+y2, board)
}

func point(x int, y int, board []int) {
    x, y = giveInd(x, y)
    board[x+y*xlim] = 1
}

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

func giveCoords(x int, y int) (int, int) {
    return x-int(xlim/2), -y+int(ylim/2)
}

func giveInd(x int, y int) (int, int) {
    return x+int(xlim/2), -y+int(ylim/2)
}

func round(i float64) int {
    return int(math.Round(i))
    //llim := int(i)
    //if i-float64(llim) >= 0.5 {return llim+1} else {return llim}
}