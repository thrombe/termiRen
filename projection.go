package main
import (
    "fmt"
    "math"
    "strings"
    "time"
    )

func main() {// REMEMBER TO USE POINTERS and fix line() and printB()
    demo()
}

func demo() {
    //board := genB()
    var x1, y1, x2, y2 float64 = 5, 6, 28, 31
    var rx, ry float64 = 5, -3
    var r float64 = 0.2
    for {
        board := genB()
        line(x1, y1, x2, y2, board)
        point(rx, ry, board)
        printB(board)
        time.Sleep(time.Millisecond*150)
        _ = board
        //fmt.Println(x1, x2, y1, y2)
        x1, y1 = rotateP(x1, y1, rx, ry, r)
        x2, y2 = rotateP(x2, y2, rx, ry, r)
    }
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
    var scr strings.Builder
    for y := 0; y < ylim; y++ {
        for x := 0; x < xlim; x++ {
            if board[x+y*xlim] == 1 {
                scr.WriteString("x")
            } else {scr.WriteString(".")}
        }
        scr.WriteString("\n")
    }
    fmt.Println(scr.String())

    //fmt.Println(0)
}

/*draw line on canvas*/
func line(x1, y1, x2, y2 float64, board []int) {
    length := math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
    ux, uy := (x2-x1)/length, (y2-y1)/length
    //ux, uy := vx/length,  vy/length
    //x, y, x3 := float64(x1), float64(y1), float64(x2)
    for int(math.Sqrt((x2-x1)*(x2-x1)+(y2-y1)*(y2-y1))) != 0 {
        point(x1, y1, board)
        x1, y1 = x1+ux, y1+uy
    }
}

/*draw vector with x1, y1 as offset and x2, y2 as direction and size*/
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

