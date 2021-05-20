package main
import (
    "fmt"
    "math"
    "seehuhn.de/go/ncurses"
    )

/*returns a generator style func that prints the board and also clears zbuf and the board*/
func perint(rawboard []rune, board [][]rune, zbuf []float64) (*ncurses.Window, func()) {
    if ncursed == 1 {
        win := ncurses.Init()
        // xlim, ylim = win.GetMaxYX() // donno if its x, y or y, x
        ncurses.CursSet(0)
        return win, func() {
            scr := string(rawboard)
            win.Erase()
            win.AddStr(scr)
            win.Refresh()
            for y := 0; y < ylim; y++ { // clear board and zbuf
                for x := 0; x < xlim; x++ {
                    board[y][x] = blank
                    zbuf[y*xlim + x] = -2 // -1 should be the clip limit (1000 i think, as defined in projection matrix)
                }
            }
        }
    } else {
        return nil, func() {
            fmt.Println(string(rawboard))
            for y := 0; y < ylim; y++ {
                for x := 0; x < xlim; x++ {
                    board[y][x] = blank
                    zbuf[y*xlim + x] = -2
                }
            }
        }
    }
}

/*make a canvas*/
func genB() ([]rune, [][]rune, []float64) {
    rawboard := make([]rune, (xlim+1)*ylim)
    board := make([][]rune, ylim)
    zbuf := make([]float64, xlim*ylim)
    var xini, xfin int = 0, xlim-1
    for y := 0; y < ylim; y++ {
        board[y] = rawboard[xini : xfin+1]
        rawboard[xfin+1] = '\n'
        xini += xlim-1+2
        xfin = xini+xlim-1
        for x := 0; x < xlim; x++ {
            zbuf[y*xlim + x] = -2
        }
    }
    return rawboard, board, zbuf
}

// //USELESS /*enter index of pixel in canvas and return coors in space*/
// func giveCoords(x, y int) (float64, float64) {
//     return float64(x-int(xlim/2)), float64(-y+int(ylim/2)) 
// }

/*enter coords of points and return index in canvas*/
func giveInd(x, y float64) (int, int) {
    return round(x)+int(xlim/2), -round(y)+int(ylim/2)
}

/*draw point on canvas*/
func point(h, k float64, board [][]rune, texture rune) {
    x, y := giveInd(h, k*charRatio)
    if 0 <= x && x < xlim && 0 <= y && y < ylim {
        board[y][x] = texture
    }
}

/*draws a 3d point on canvas using zbuf*/
func point3d(p [][]float64, board [][]rune, zbuf []float64, texture rune) {
    x, y := giveInd(p[0][0], p[1][0]*charRatio)
    if !(0 <= x && x < xlim && 0 <= y && y < ylim) {return}
    if zbuf[y*xlim + x] < p[2][0] {
        zbuf[y*xlim + x] = p[2][0]
        board[y][x] = texture
    }
}

/*draw line on canvas*/
func line(v1, v2 [][]float64, board [][]rune, texture rune) {
    x1, y1, x2, y2 := v1[0][0], v1[1][0], v2[0][0], v2[1][0]
    length := math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
    if round(length) == 0 { // if line is too small, just draw a point
        point(x1, y1, board, texture)
        return
    }
    ux, uy := (x2-x1)/length, (y2-y1)/length
    for math.Round(2*((x2-x1)*(x2-x1)+(y2-y1)*(y2-y1))) != 0 { // that 2 there is for resolution of line
        point(x1, y1, board, texture)
        x1, y1 = x1+ux, y1+uy
    }
}

/*draw vector with v1 as offset and v2 as direction and size*/
func vector(v1, v2 [][]float64, board [][]rune, texture rune) {
    line(v1, matAdd(v1, v2), board, texture)
}

// /*rotate x, y about ox, oy by r radians*/
// func rotateP2d(x, y, ox, oy, r float64) (float64, float64) {
//     x, y = x-ox, y-oy
//     x, y =  x*math.Cos(r)-y*math.Sin(r), x*math.Sin(r)+y*math.Cos(r)
//     return x+ox, y+oy
// }
