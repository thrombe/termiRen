package main
import (
    "fmt"
    "math"
    "strings"
    "seehuhn.de/go/ncurses"
    "time"
    )

func perint(board []int) (*ncurses.Window, func()) {
    if ncursed == 1 {
        win := ncurses.Init()
        // xlim, ylim = win.GetMaxYX() // donno if its x, y or y, x
        ncurses.CursSet(0)
        return win, func() {
            scr := printB(board)
            time.Sleep(time.Millisecond*50)
            win.Erase()
            win.AddStr(scr)
            win.Refresh()
        }
    } else {
        return nil, func() {
            fmt.Println(printB(board))
        }
    }
}

func printB(board []int) string {
    var scr strings.Builder
    for y := 0; y < ylim; y++ {
        for x := 0; x < xlim; x++ {
            if board[x+y*xlim] == 1 {
                scr.WriteString(".")
                board[x+y*xlim] = 0
            } else {scr.WriteString(" ")}
        }
        scr.WriteString("\n")
    }
    // try returning the rune slice instead. this method copies the string while returning
    return scr.String()
    //fmt.Println(scr.String())
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

//USELESS /*enter index of pixel in canvas and return coors in space*/
func giveCoords(x, y int) (float64, float64) {
    return float64(x-int(xlim/2)), float64(-y+int(ylim/2)) 
}

/*enter coords of points and return index in canvas*/
func giveInd(x, y float64) (int, int) {
    return round(x)+int(xlim/2), -round(y)+int(ylim/2)
}

/*draw point on canvas*/
func point(h, k float64, board []int) {
    x, y := giveInd(h, k*charRatio)
    if 0 <= x && x < xlim && 0 <= y && y < ylim {
        board[x+y*xlim] = 1
    }
}

/*draw line on canvas*/
func line(x1, y1, x2, y2 float64, board []int) {
    length := math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
    ux, uy := (x2-x1)/length, (y2-y1)/length
    for math.Round(2*((x2-x1)*(x2-x1)+(y2-y1)*(y2-y1))) != 0 {
        point(x1, y1, board)
        x1, y1 = x1+ux, y1+uy
    }
}

/*draw vector with x1, y1 as offset and x2, y2 as direction and size*/
func vector(x1, y1, x2, y2 float64, board []int) {
    line(x1, y1, x1+x2, y1+y2, board)
}

/*rotate x, y about ox, oy by r radians*/
func rotateP2d(x, y, ox, oy, r float64) (float64, float64) {
    x, y = x-ox, y-oy
    x, y =  x*math.Cos(r)-y*math.Sin(r), x*math.Sin(r)+y*math.Cos(r)
    return x+ox, y+oy
}
