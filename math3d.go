package main

import "math"

//returns projection matrix with 1 and 1000 as near and far
func projectionMat() [][]float64 {
    cot := 1/math.Tan(fov/2)
    // scrDist := float64(xlim)*cot/2
    f := 1000.0
    n := 1.0
    scale := float64(xlim)/2
    return matrix(4, 4, 
        cot*scale, 0, 0, 0,
        0, cot*scale, 0, 0,
        0, 0, (f+n)/(f-n), (2*f*n)/(f-n),
        0, 0, -1, 0,
    )
}

//returns rotation matrix z angle from z axis
func rotMat3dz(z float64) [][]float64 {
    return matrix(4, 4,
        math.Cos(z), -math.Sin(z), 0, 0,
        math.Sin(z), math.Cos(z), 0, 0,
        0, 0, 1, 0,
        0, 0, 0, 1,
    )
}

//returns rotation matrix of y angle from y axis
func rotMat3dy(y float64) [][]float64 {
    return matrix(4, 4, // i cross k is -j
        math.Cos(y), 0, math.Sin(y), 0,
        0, 1, 0, 0,
        -math.Sin(y), 0, math.Cos(y), 0,
        0, 0, 0, 1,
    )
}

//returns rotation matrix of x angle from x axis
func rotMat3dx(x float64) [][]float64 {
    return matrix(4, 4,
        1, 0, 0, 0,
        0, math.Cos(x), -math.Sin(x), 0,
        0, math.Sin(x), math.Cos(x), 0,
        0, 0, 0, 1,
    )
}

//returns a scaling matrix or whatever its called
// s is the scale factor and matsize is the len of square mat
// multiply a vector with it(vector to the right), and everything
// except the last row will be scaled
func scaleMat(s float64, matsize int) [][]float64 {
    result := matrix(matsize, matsize)
    for i := 0; i < matsize-1 ; i++ {
        result[i][i] = s
    }
    result[matsize-1][matsize-1] = 1
    return result
}

//returns 3d translation matrix (adds o to coords)
func transMat(o [][]float64) [][]float64 {
    if len(o[0]) != 1 {panic("transMat matrix shape error")}
    length := len(o)
    return matrix(length, 4,
        1, 0, 0, o[0][0],
        0, 1, 0, o[1][0],
        0, 0, 1, o[2][0],
        0, 0, 0, 1,
    )
    // trans := make([][]float64, length)
    // for i, ele := range o {
    //     row := []float64 {0, 0, 0, ele[0]}
    //     row[i] = 1
    //     trans[i] = row
    // }
    // return trans
}

//returns translation matrix (subs o from coords
func transMatInv(o [][]float64) [][]float64 {
    return transMat(matMul(scaleMat(-1, len(o)), o))
}

//adds back and forth translation to rotation matrix
func rotAboutPoint(rot ,o [][]float64) [][]float64 {
    return nMatMul(transMat(o), rot, transMatInv(o))
}

//axis passes through the origin. (translate to other point to get other axes)
func rotAboutVec(angle float64, axis [][]float64) [][]float64 {
    // y and z are angles around y and z axes.
    // transforming the axis into x-axis
    y := math.Atan2(axis[2][0], axis[0][0]) // tan^-1(z/x) // turn around y
    ry := rotMat3dy(y)
    axis = matMul(ry, axis) // turn the axis around y
    z := math.Atan2(axis[1][0], axis[0][0]) // tan^-1(y/x) // turn around z but according to the new axis
    rot := rotMat3dx(angle)
    ryinv := rotMat3dy(-y)
    rz := rotMat3dz(-z)
    rzinv := rotMat3dz(z)
    return nMatMul(ryinv, rzinv, rot, rz, ry)
}

//returns the cross product of two 3d vectors(4 by 1)
func vecCross(vec1, vec2 [][]float64) [][]float64 {
    return matrix(4, 1,
        vec1[1][0]*vec2[2][0]-vec2[1][0]*vec1[2][0],
        -vec1[0][0]*vec2[2][0]+vec2[0][0]*vec1[2][0],
        vec1[0][0]*vec2[1][0]-vec2[0][0]*vec1[1][0],
        0,
    )
}

// multiplies matrix to 3D vector, but dosent make a copy. 
// enter 1 or 2 vectors, multiplication of mat and 1st is saved in second (or 1st if second not entered)
func vecMul(mat [][]float64, veclist ...[][]float64) {
    var vec, savein [][]float64
    switch len(veclist) {
    case 1:
        vec, savein = veclist[0], veclist[0]
    case 2:
        vec, savein = veclist[0], veclist[1]
    }

    // this had to be done in 1 line
    savein[0][0], savein[1][0], savein[2][0], savein[3][0] = (mat[0][0]*vec[0][0] + mat[0][1]*vec[1][0] + mat[0][2]*vec[2][0] + mat[0][3]*vec[3][0]), (mat[1][0]*vec[0][0] + mat[1][1]*vec[1][0] + mat[1][2]*vec[2][0] + mat[1][3]*vec[3][0]), (mat[2][0]*vec[0][0] + mat[2][1]*vec[1][0] + mat[2][2]*vec[2][0] + mat[2][3]*vec[3][0]), (mat[3][0]*vec[0][0] + mat[3][1]*vec[1][0] + mat[3][2]*vec[2][0] + mat[3][3]*vec[3][0])
}