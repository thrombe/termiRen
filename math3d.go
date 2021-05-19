package main

import "math"

/*returns projection matrix with 1 and 1000 as near and far*/
func projectionMat() [][]float64 {
    cot := 1/math.Tan(fov/2)
    // scrDist := float64(xlim)*cot/2
    f := 1000.0
    n := 1.0
    scale := float64(xlim)/2
    return [][]float64 {
        {cot*scale, 0, 0, 0},
        {0, cot*scale, 0, 0},
        {0, 0, (f+n)/(f-n), (2*f*n)/(f-n)},
        {0, 0, -1, 0},
    }
}

/*returns rotation matrix z angle from z axis*/
func rotMat3dz(z float64) [][]float64 {
    rotmat := [][]float64 {
        {math.Cos(z), -math.Sin(z), 0, 0},
        {math.Sin(z), math.Cos(z), 0, 0},
        {0, 0, 1, 0},
        {0, 0, 0, 1},
    }
    return rotmat
}

/*returns rotation matrix of y angle from y axis*/
func rotMat3dy(y float64) [][]float64 {
    rotmat := [][]float64 { // i cross k is -j
        {math.Cos(y), 0, math.Sin(y), 0},
        {0, 1, 0, 0},
        {-math.Sin(y), 0, math.Cos(y), 0},
        {0, 0, 0, 1},
    }
    return rotmat
}

/*returns rotation matrix of x angle from x axis*/
func rotMat3dx(x float64) [][]float64 {
    rotmat := [][]float64 {
        {1, 0, 0, 0},
        {0, math.Cos(x), -math.Sin(x), 0},
        {0, math.Sin(x), math.Cos(x), 0},
        {0, 0, 0, 1},
    }
    return rotmat
}

/*returns a scaling matrix or whatever its called
s is the scale factor and matsize is the len of square mat
multiply a vector with it(vector to the right), and everything
except the last row will be scaled*/
func scaleMat(s float64, matsize int) [][]float64 {
    result := make([][]float64, matsize)
    for i := 0; i < matsize ; i++ {
        row := make([]float64, matsize)
        row[i] = s
        result[i] = row
    }
    result[matsize-1][matsize-1] = 1
    return result
}

/*returns translation matrix (adds o to coords)*/
func transMat(o [][]float64) [][]float64 {
    if len(o[0]) != 1 {panic("transMat matrix shape error")}
    length := len(o)
    trans := make([][]float64, length)
    for i, ele := range o {
        row := []float64 {0, 0, 0, ele[0]}
        row[i] = 1
        trans[i] = row
    }
    return trans
}

/*returns translation matrix (subs o from coords*/
func transMatInv(o [][]float64) [][]float64 {
    return transMat(matMul(scaleMat(-1, len(o)), o))
}

/*adds back and forth translation to rotation matrix*/
func rotAboutPoint(rot ,o [][]float64) [][]float64 {
    return nMatMul(transMat(o), rot, transMatInv(o))
}

/*axis passes through the origin. (translate to other point to get other axes)*/
func rotAboutVec(angle float64, axis [][]float64) [][]float64 {
    // y and z are angles around y and z axes.
    // transforming the axis into x-axis
    y := math.Atan2(axis[2][0], axis[0][0]) // tan^-1(z/x)
    z := math.Atan2(axis[1][0], axis[0][0]) // tan^-1(y/x)
    rot := rotMat3dx(angle)
    ry := rotMat3dy(y)
    ryinv := rotMat3dy(-y)
    rz := rotMat3dz(-z)
    rzinv := rotMat3dz(z)
    return nMatMul(ryinv, rzinv, rot, rz, ry)
}

/*returns the cross product of two 3d vectors(4 by 1)*/
func vecCross(vec1, vec2 [][]float64) [][]float64 {
    return [][]float64 {
        {vec1[1][0]*vec2[2][0]-vec2[1][0]*vec1[2][0]},
        {-vec1[0][0]*vec2[2][0]+vec2[0][0]*vec1[2][0]},
        {vec1[0][0]*vec2[1][0]-vec2[0][0]*vec1[1][0]},
        {0},
    }
}
