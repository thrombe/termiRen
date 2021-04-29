package main
import "testing"

func Test_matDet(t *testing.T) {
    mat := [][]float64 {{0, 1, 0}, {4, 66, 3}, {2, 22, 3}}
    out := matDet(mat)
    if out != -6 {
        t.Errorf("input: %v expected: -6 but receieved %v", mat, out)
    }
}

func Test_vecDot(t *testing.T) {
    vec1 := [][]float64 {{10}, {2}}
    vec2 := [][]float64 {{0}, {6}}
    out := vecDot(vec1, vec2)
    if out != 12 {
        t.Errorf("(10, 2) dot (0, 2) should be 12 and not %v", out)
    }
}

func TestTable_vecSize(t *testing.T) {
    var tests = []struct {
        input [][]float64
        output float64
    } {
        { [][]float64 {{3}, {4}}, 5 },
        { [][]float64 {{12}, {5}}, 13 },
        { [][]float64 {{3}, {4}}, 5 },
        { [][]float64 {{3}, {4}}, 5 },
    }
    for _, test := range tests {
        if out := vecSize(test.input); out != test.output {
            t.Errorf("for input: %v output should be: %v but was:%v", test.input, test.output, out)
        }
    }
}