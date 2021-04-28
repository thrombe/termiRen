package main
import (
    // "fmt"
    "math"
)

/*multiply matrices uf arbitary sizes(legal only ofc)*/
func matMul(mat1, mat2 [][]float64) [][]float64 {
    m1rows, m1cols, m2rows, m2cols := len(mat1), len(mat1[0]), len(mat2), len(mat2[0])
    if m1rows != m2cols {panic("matMul shape error")}
    result := make([][]float64, m1rows)
    for r := 0; r < m1rows; r++ {
        row := make([]float64, m2rows)
        for c := 0; c < m2cols; c++ {
            for item := 0; item < m1cols; item++ {
                row[item] += mat1[r][item]*mat2[item][c]
            }
        }
        result[r] = row
    }
    return result
}

func matAdd(mat1, mat2 [][]float64) [][]float64 {
    m1rows, m1cols, m2rows, m2cols := len(mat1), len(mat1[0]), len(mat2), len(mat2[0])
    if !(m1rows == m2rows && m1cols == m2cols) {panic("matAdd shape error")}
    result := make([][]float64, m1rows)
    for r := 0; r < m1rows; r++ {
        row := make([]float64, m1cols)
        for c := 0; c < m1cols; c++ {
            row[c] = mat1[r][c] + mat2[r][c]
        }
        result[r] = row
    }
    return result
}

func matSub(mat1, mat2 [][]float64) [][]float64 {
    m1rows, m1cols, m2rows, m2cols := len(mat1), len(mat1[0]), len(mat2), len(mat2[0])
    if !(m1rows == m2rows && m1cols == m2cols) {panic("matSub shape error")}
    result := make([][]float64, m1rows)
    for r := 0; r < m1rows; r++ {
        row := make([]float64, m1cols)
        for c := 0; c < m1cols; c++ {
            row[c] = mat1[r][c] - mat2[r][c]
        }
        result[r] = row
    }
    return result
}

func matScale(mat [][]float64, scale float64) [][]float64 {
    mrows, mcols := len(mat), len(mat[0])
    for r := 0; r < mrows; r++ {
        for c := 0; c < mcols; c++ {
            mat[r][c] = mat[r][c]*scale
        }
    }
    return mat
}

/*remember to input y and x index resp. *///returns -1^(x+y)*submat */
func subMat(mat [][]float64, y, x int) [][]float64 {
    
    mrows := len(mat)
    var submat [][]float64
    for r := 0; r < mrows; r++ {
        if r == y {continue}
        var row []float64
        for c := 0; c < mrows; c++ {

            if c == x {continue}
            row = append(row, mat[r][c])
        }
        
        submat = append(submat, row)
    }
    return submat
    //return matScale(submat, math.Pow(-1, float64(x+y)))
}

/*returns determinant of a square matrix*/
func matDet(mat [][]float64) float64 {
    mrows, mcols := len(mat), len(mat[0])
    if mrows != mcols {panic("matDet non square matrix")}
    if mrows == 2 {return mat[0][0]*mat[1][1] - mat[0][1]*mat[1][0]}
    submat := make([][]float64, mrows-1)
    var result float64
    for r := 0; r < mrows; r++ {
        for c := 0; c < mrows; c++ {
            result += mat[r][c] * math.Pow(-1, float64(c)) * matDet(submat)
        }
    }

    panic("matDet not imllemented")
    return 0
}

/*returns the size of a n-dimentional vector*/
func vecSize(vec [][]float64) float64 {
    if len(vec[0]) > 1 {panic("vecSize not a vector")}
    vecDimensions := len(vec)
    var result float64
    for r := 0; r < vecDimensions; r++ {
        result += vec[r][0] * vec[r][0]
    }
    return math.Sqrt(result)
}

/*returns dot of 2 vectors*/
func vecDot(vec1, vec2 [][]float64) float64 {
    if len(vec1[0]) > 1 || len(vec2) > 1 {panic("vecDot not a vector")}
    if len(vec1) != len(vec2) {panic("vecDot vectors of different dimentions")}
    vecDimensions := len(vec1)
    var result float64
    for r := 0; r < vecDimensions; r++ {
        result += vec1[r][0] * vec2[r][0]
    }
    return result
}

func absVal(n float64) float64 {
    if n >= 0 {return n} else {return -n}
}

/*rounds to nearest int and retrun int*/
func round(i float64) int {
    return int(math.Round(i))
    //llim := int(i)
    //if i-float64(llim) >= 0.5 {return llim+1} else {return llim}
}