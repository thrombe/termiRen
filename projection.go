package main

import (
    "fmt"
    // "math"
    "seehuhn.de/go/ncurses"
    )

//returns a generator style func that prints the board and also clears zbuf and the board
func perint(cam *camera) (*ncurses.Window, func()) {
    xlim, ylim, board, rawboard, zbuf := cam.xlim, cam.ylim, cam.board, cam.rawboard, cam.zbuf
    if cam.ncursed {
        win := ncurses.Init()
        // xlim, ylim = win.GetMaxYX() // donno if its x, y or y, x
        ncurses.CursSet(0)
        return win, func() {
            win.Erase()
            win.AddStr(string(rawboard))
            win.Refresh()
            for y := 0; y < ylim; y++ { // clear board and zbuf
                for x := 0; x < xlim; x++ {
                    board[y][x] = ' '
                    zbuf[y][x] = -2 // -1 should be the clip limit (1000 i think, as defined in projection matrix)
                }
            }
        }
    } else {
        return nil, func() {
            fmt.Println(string(rawboard))
            for y := 0; y < ylim; y++ {
                for x := 0; x < xlim; x++ {
                    board[y][x] = ' '
                    zbuf[y][x] = -2
                }
            }
        }
    }
}

//make a canvas
func genB(cam *camera){
    xlim, ylim := cam.xlim, cam.ylim
    cam.rawboard = make([]byte, (xlim+1)*ylim)
    cam.board = make([][]byte, ylim)
    zbufrray := make([]float64, xlim*ylim)
    cam.zbuf = make([][]float64, ylim)
    var xini, xfin int = 0, xlim-1
    for y := 0; y < ylim; y++ {
        cam.board[y] = cam.rawboard[xini : xfin+1]
        cam.rawboard[xfin+1] = '\n'
        xini += xlim-1+2
        xfin = xini+xlim-1
        cam.zbuf[y] = zbufrray[y*xlim : (y+1)*xlim]
        for x := 0; x < xlim; x++ {
            cam.zbuf[y][x] = -2
        }
    }
}

// //USELESS /*enter index of pixel in canvas and return coors in space*/
// func giveCoords(x, y int) (float64, float64) {
//     return float64(x-int(xlim/2)), float64(-y+int(ylim/2)) 
// }

//enter coords of points and return index in canvas
func giveInd(x, y float64, xlim, ylim int) (int, int) {
    return round(x)+xlim/2, -round(y)+ylim/2
}

//draw point on canvas
func point(h, k float64, cam *camera, texture byte) {
    x, y := giveInd(h, k*cam.charRatio, cam.xlim, cam.ylim)
    if 0 <= x && x < cam.xlim && 0 <= y && y < cam.ylim {
        cam.board[y][x] = texture
    }
}

//draws a 3d point on canvas using zbuf
func point3d(p [][]float64, cam *camera, texture byte) {
    x, y := giveInd(p[0][0], p[1][0]*cam.charRatio, cam.xlim, cam.ylim)
    if !(0 <= x && x < cam.xlim && 0 <= y && y < cam.ylim) {return}
    if cam.zbuf[y][x] < p[2][0] {
        cam.zbuf[y][x] = p[2][0]
        cam.board[y][x] = texture
    }
}

//draw line on canvas
func line(v1, v2 [][]float64, cam *camera, texture byte) {
    x, y, x2, y2 := round(v1[0][0]), round(v1[1][0]), round(v2[0][0]), round(v2[1][0])
    var steepx, steepy, ydir int
    if absVal(float64(x2-x)) > absVal(float64(y2-y)) {steepx, steepy = 1, 0} else {steepx, steepy, x, x2, y, y2 = 0, 1, y, y2, x, x2} // if slope > 1
    if x > x2 {x, x2, y, y2 = x2, x, y2, y} // x2 to the right of x
    dx, dydy := x2-x, (y2-y)*2 // we only need 2*dy
    if dydy > 0 {ydir = 1} else {ydir, dydy = -1, -dydy} // we only need magnitude of dy. ydir is slope > 0 or slope < 0
    eror := 0
    for ; x <= x2; x++ {
        point(float64(x*steepx+y*steepy), float64(y*steepx+x*steepy), cam, texture)
        eror += dydy // similar to err += dy/dx and then checking if err > .5, if yes err -= 1 (hint- multiplying by 2*dx eliminates need of floats)
        if eror > dx {
            y += ydir
            eror -= 2*dx
        }
    }
}

// draw vector with v1 as offset and v2 as direction and size
func drawVec(v1, v2 [][]float64, cam *camera, texture  byte) {
    line(v1, matAdd(v1, v2), cam, texture)
}

// /*rotate x, y about ox, oy by r radians*/
// func rotateP2d(x, y, ox, oy, r float64) (float64, float64) {
//     x, y = x-ox, y-oy
//     x, y =  x*math.Cos(r)-y*math.Sin(r), x*math.Sin(r)+y*math.Cos(r)
//     return x+ox, y+oy
// }
