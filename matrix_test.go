package golina

import (
	"math"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func GenerateRandomFloat() float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()
}

func GenerateRandomVector(size int) *Vector {
	slice := make(Vector, size, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		slice[i] = rand.Float64() - rand.Float64()
	}
	return &slice
}

func GenerateRandomSymmetric33Matrix() *Matrix {
	entries := *GenerateRandomVector(6)
	m := ZeroMatrix(3, 3)
	m.Set(0, 0, entries[0])
	m.Set(1, 1, entries[1])
	m.Set(2, 2, entries[2])
	m.Set(0, 1, entries[3])
	m.Set(1, 0, entries[3])
	m.Set(0, 2, entries[4])
	m.Set(2, 0, entries[4])
	m.Set(1, 2, entries[5])
	m.Set(2, 1, entries[5])
	return m
}

func GenerateRandomSquareMatrix(size int) *Matrix {
	return GenerateRandomMatrix(size, size)
}

func GenerateRandomMatrix(row, col int) *Matrix {
	rows := make(Data, row)
	for i := range rows {
		rows[i] = *GenerateRandomVector(col)
	}
	m := new(Matrix).Init(rows)
	return m
}

// https://blog.karenuorteva.fi/go-unit-test-setup-and-teardown-db1601a796f2#.2aherx2z5

func TestMatrix_Init(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	if matA._array == nil {
		t.Fail()
	}
}

func TestMatrix_String(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	if matA.String() != "{1.000000, 2.000000, 3.000000,\n 4.000000, 5.000000, 6.000000,\n 7.000000, 8.000000, 9.000000}\n" {
		t.Fail()
	}
}

func TestMatrix_Dims(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	row, col := matA.Dims()
	if row != 3 || col != 3 {
		t.Fail()
	}
}

func TestMatrix_At(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	if matA.At(1, 1) != 5 {
		t.Fail()
	}
}

func TestMatrix_Set(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	matA.Set(1, 1, 10)
	if matA.At(1, 1) != 10 {
		t.Fail()
	}
}

func TestMatrix_T(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	matAT := matA.T()
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if matAT.At(i, j) != matA.At(j, i) {
				t.Fail()
			}
		}
	}
}

func TestMatrix_Row(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	row := matA.Row(1)
	if !VEqual(row, &Vector{4, 5, 6}) {
		t.Fail()
	}
}

func TestMatrix_Col(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	col := matA.Col(1)
	if !VEqual(col, &Vector{2, 5, 8}) {
		t.Fail()
	}
}

func TestMatrix_Max(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	maxA := matA.Max()
	if maxA.value != 9 && maxA.row != 2 && maxA.col != 2 {
		t.Fail()
	}
}

func TestMatrix_Min(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	minA := matA.Min()
	if minA.value != 1 && minA.row != 0 && minA.col != 0 {
		t.Fail()
	}
}

func TestCopy(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	matB := Copy(matA)
	if !MEqual(matA, matB) {
		t.Fail()
	}
}

func TestEmpty(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	matB := Empty(matA)
	row1, col1 := matA.Dims()
	row2, col2 := matB.Dims()
	if [2]int{row1, col1} != [2]int{row2, col2} {
		t.Fail()
	}
}

func TestZeroMatrix(t *testing.T) {
	a := Data{{0}, {0}}
	matA := ZeroMatrix(2, 1)
	if !MEqual(matA, new(Matrix).Init(a)) {
		t.Fail()
	}
}

func TestOneMatrix(t *testing.T) {
	a := Data{{1}, {1}}
	matA := OneMatrix(2, 1)
	if !MEqual(matA, new(Matrix).Init(a)) {
		t.Fail()
	}
}

func TestIdentityMatrix(t *testing.T) {
	a := Data{{1, 0}, {0, 1}}
	matA := IdentityMatrix(2)
	if !MEqual(matA, new(Matrix).Init(a)) {
		t.Fail()
	}
}

func TestSwapRow(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	matB := Copy(matA)
	SwapRow(matB, 1, 2)
	if !VEqual(matA.Row(1), matB.Row(2)) || !VEqual(matA.Row(2), matB.Row(1)) {
		t.Fail()
	}
}

func TestMatrix_Rank(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 10}, {30, 50, 0}}
	matA := new(Matrix).Init(a)
	if matA.Rank() != 2 {
		t.Fail()
	}

	b := Data{{0, 1, 2}, {-1, -2, 1}, {2, 7, 8}}
	matB := new(Matrix).Init(b)
	if matB.Rank() != 3 {
		t.Fail()
	}

	c := Data{{0, 1, 2}, {-1, -2, 1}, {2, 7, 8}, {3, 5, 3}}
	matC := new(Matrix).Init(c)
	if matC.Rank() != 3 {
		t.Fail()
	}
}

func TestMatrix_Det(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 10}, {30, 50, 0}}
	matA := new(Matrix).Init(a)
	if !FloatEqual(matA.Det(), 0) {
		t.Fail()
	}

	b := Data{{32, 12, 1}, {6, 3, 45}, {9, 2, 1}}
	matB := new(Matrix).Init(b)
	if matB.Det() != 1989 {
		t.Fail()
	}
}

func TestMatrix_Inverse(t *testing.T) {
	a := Data{{32, 12, 1}, {6, 3, 45}, {9, 2, 1}}
	matA := new(Matrix).Init(a)
	b := Data{{-0.04374057315233785821, -0.00502765208647561595, 0.26998491704374057313},
		{0.2006033182503770739, 0.0115635997988939167, -0.7209653092006033182},
		{-0.007541478129713423831, 0.022121669180492709902, 0.012066365007541478129}}
	if !MEqual(matA.Inverse(), new(Matrix).Init(b)) {
		t.Fail()
	}
}

func TestNaiveDet(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 10}, {30, 50, 0}}
	matA := new(Matrix).Init(a)
	if NaiveDet(matA) != 0 {
		t.Fail()
	}

	b := Data{{32, 12, 1}, {6, 3, 45}, {9, 2, 1}}
	matB := new(Matrix).Init(b)
	if NaiveDet(matB) != 1989 {
		t.Fail()
	}
}

func TestNaiveAdj(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 10}, {30, 50, 0}}
	matA := new(Matrix).Init(a)
	b := Data{{-500, 500, 500}, {300, -300, -300}, {-100, 100, 100}}
	if !MEqual(NaiveAdj(matA), new(Matrix).Init(b)) {
		t.Fail()
	}
}

func TestNaiveInverse(t *testing.T) {
	a := Data{{32, 12, 1}, {6, 3, 45}, {9, 2, 1}}
	matA := new(Matrix).Init(a)
	b := Data{{-0.04374057315233785821, -0.00502765208647561595, 0.26998491704374057313},
		{0.2006033182503770739, 0.0115635997988939167, -0.7209653092006033182},
		{-0.007541478129713423831, 0.022121669180492709902, 0.012066365007541478129}}
	if !MEqual(NaiveInverse(matA), new(Matrix).Init(b)) {
		t.Fail()
	}
}

func TestMatrix_Add(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 10}, {30, 50, 0}}
	matA := new(Matrix).Init(a)
	b := Data{{32, 12, 1}, {6, 3, 45}, {9, 2, 1}}
	matB := new(Matrix).Init(b)
	matC := matA.Add(matB)
	c := Data{{42, 32, 11}, {-14, -27, 55}, {39, 52, 1}}
	if !MEqual(matC, new(Matrix).Init(c)) {
		t.Fail()
	}
}

func TestMatrix_Sub(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 10}, {30, 50, 0}}
	matA := new(Matrix).Init(a)
	b := Data{{32, 12, 1}, {6, 3, 45}, {9, 2, 1}}
	matB := new(Matrix).Init(b)
	matC := matA.Sub(matB)
	c := Data{{-22, 8, 9}, {-26, -33, -35}, {21, 48, -1}}
	if !MEqual(matC, new(Matrix).Init(c)) {
		t.Fail()
	}
}

func TestMatrix_Mul(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 10}, {30, 50, 0}}
	matA := new(Matrix).Init(a)
	b := Data{{32, 12}, {6, 3}, {9, 2}}
	matB := new(Matrix).Init(b)
	matC := matA.Mul(matB)
	c := Data{{530, 200}, {-730, -310}, {1260, 510}}
	if !MEqual(matC, new(Matrix).Init(c)) {
		t.Fail()
	}
}

func TestMatrix_MulVec(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 10}, {30, 50, 0}}
	matA := new(Matrix).Init(a)
	b := &Vector{1, 2, 3}
	v := matA.MulVec(b)
	if !VEqual(v, &Vector{80, -50, 130}) {
		t.Fail()
	}
}

func TestMatrix_MulNum(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 10}, {30, 50, 0}}
	matA := new(Matrix).Init(a)
	b := Data{{30, 60, 30}, {-60, -90, 30}, {90, 150, 0}}
	if !MEqual(matA.MulNum(3), new(Matrix).Init(b)) {
		t.Fail()
	}
}

func TestMatrix_GetDiagonalElements(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 10}, {30, 50, 0}}
	matA := new(Matrix).Init(a)
	v := &Vector{10, -30, 0}
	if !VEqual(matA.GetDiagonalElements(), v) {
		t.Fail()
	}
}

func TestMatrix_Pow(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 10}, {30, 50, 0}}
	matA := new(Matrix).Init(a)
	if !MEqual(matA.Pow(0), IdentityMatrix(3)) {
		t.Fail()
	}
	if !MEqual(matA.Pow(1), matA) {
		t.Fail()
	}
	b := Data{{0, 100, 300}, {700, 1000, -500}, {-700, -900, 800}}
	if !MEqual(matA.Pow(2), new(Matrix).Init(b)) {
		t.Fail()
	}
}

func TestNaivePow(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 10}, {30, 50, 0}}
	matA := new(Matrix).Init(a)
	if !MEqual(NaivePow(matA, 0), IdentityMatrix(3)) {
		t.Fail()
	}
	if !MEqual(NaivePow(matA, 1), matA) {
		t.Fail()
	}
	b := Data{{0, 100, 300}, {700, 1000, -500}, {-700, -900, 800}}
	if !MEqual(NaivePow(matA, 2), new(Matrix).Init(b)) {
		t.Fail()
	}
}

func TestMatrix_Trace(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 10}, {30, 50, 0}}
	matA := new(Matrix).Init(a)
	if matA.Trace() != -20 {
		t.Fail()
	}
}

func TestMatrix_Norm(t *testing.T) {
	a := Data{{2, -2, 1}, {-1, 3, -1}, {2, -4, 1}}
	matA := new(Matrix).Init(a)
	if !FloatEqual(matA.Norm(), 6.403124237) {
		t.Fail()
	}
}

func TestMatrix_Flat(t *testing.T) {
	a := Data{{2, -2, 1}, {-1, 3, -1}, {2, -4, 1}}
	matA := new(Matrix).Init(a)
	if !VEqual(matA.Flat(), &Vector{2, -2, 1, -1, 3, -1, 2, -4, 1}) {
		t.Fail()
	}
}

func TestMatrix_GetSubMatrix(t *testing.T) {
	a := Data{{2, -2, 1}, {-1, 3, -1}, {2, -4, 1}}
	matA := new(Matrix).Init(a)
	b := Data{{3, -1}, {-4, 1}}
	if !MEqual(matA.GetSubMatrix(1, 1, 2, 2), new(Matrix).Init(b)) {
		t.Fail()
	}
}

func TestMatrix_SetSubMatrix(t *testing.T) {
	a := Data{{2, -2, 1}, {-1, 3, -1}, {2, -4, 1}}
	matA := new(Matrix).Init(a)
	b := Data{{0, 5}, {-5, 8}}
	c := Data{{2, -2, 1}, {-1, 0, 5}, {2, -5, 8}}
	matA.SetSubMatrix(1, 1, new(Matrix).Init(b))
	if !MEqual(matA, new(Matrix).Init(c)) {
		t.Fail()
	}
}

func TestMatrix_SumCol(t *testing.T) {
	a := Data{{2, -2, 1}, {-1, 3, -1}, {2, -4, 1}}
	matA := new(Matrix).Init(a)
	if !FloatEqual(matA.SumCol(0), 3) {
		t.Fail()
	}
}

func TestMatrix_SumRow(t *testing.T) {
	a := Data{{2, -2, 1}, {-1, 3, -1}, {2, -4, 1}}
	matA := new(Matrix).Init(a)
	if !FloatEqual(matA.SumRow(0), 1) {
		t.Fail()
	}
}

func TestMatrix_Sum(t *testing.T) {
	a := Data{{2, -2, 1}, {-1, 3, -1}, {2, -4, 1}}
	matA := new(Matrix).Init(a)
	b := &Vector{3, -3, 1}
	if !VEqual(matA.Sum(0), b) {
		t.Fail()
	}
	b = &Vector{1, 1, -1}
	if !VEqual(matA.Sum(1), b) {
		t.Fail()
	}
	b = &Vector{1}
	if !VEqual(matA.Sum(-1), b) {
		t.Fail()
	}
}

func TestMatrix_Mean(t *testing.T) {
	a := Data{{2, -2, 1}, {-1, 3, -1}, {2, -4, 1}}
	matA := new(Matrix).Init(a)
	b := &Vector{1, -1, 1. / 3}
	if !VEqual(matA.Mean(0), b) {
		t.Fail()
	}
	b = &Vector{1. / 3, 1. / 3, -1. / 3}
	if !VEqual(matA.Mean(1), b) {
		t.Fail()
	}
	b = &Vector{1. / 9}
	if !VEqual(matA.Mean(-1), b) {
		t.Fail()
	}
}

func TestMatrix_CovMatrix(t *testing.T) {
	a := Data{{2, -2, 1}, {-1, 3, -1}, {2, -4, 1}}
	matA := new(Matrix).Init(a)
	b := Data{{3, -6, 2}, {-6, 13, -4}, {2, -4, 1.33333333333333}}
	if !MEqual(matA.CovMatrix(), new(Matrix).Init(b)) {
		t.Fail()
	}
}

func TestMatrix_IsSymmetric(t *testing.T) {
	if new(Matrix).Init(Data{{2, -2, 1}, {-1, 3, -1}, {2, -4, 1}}).IsSymmetric() {
		t.Fail()
	}
	if !GenerateRandomSymmetric33Matrix().IsSymmetric() {
		t.Fail()
	}
}

func TestCrossCovMatrix(t *testing.T) {
	a := Data{{2, -2, 1}, {-1, 3, -1}, {2, -4, 1}}
	matA := new(Matrix).Init(a)
	if !MEqual(CrossCovMatrix(matA, matA), matA.CovMatrix()) {
		t.Fail()
	}
}

// Vector
func TestVector_At(t *testing.T) {
	v := &Vector{1, 2, 3}
	if v.At(0) != 1 || v.At(-1) != 3 {
		t.Fail()
	}
}

func TestVector_Add(t *testing.T) {
	v1 := &Vector{1, 2, 3}
	v2 := &Vector{1, 2, 3}
	if !VEqual(v1.Add(v2), &Vector{2, 4, 6}) {
		t.Fail()
	}
}

func TestVector_AddNum(t *testing.T) {
	v1 := &Vector{1, 2, 3}
	n := 1
	if !VEqual(v1.AddNum(n), &Vector{2, 3, 4}) {
		t.Fail()
	}
}

func TestVector_Sub(t *testing.T) {
	v1 := &Vector{1, 2, 3}
	v2 := &Vector{1, 2, 3}
	if !VEqual(v1.Sub(v2), &Vector{0, 0, 0}) {
		t.Fail()
	}
}

func TestVector_SubNum(t *testing.T) {
	v1 := &Vector{1, 2, 3}
	n := 1
	if !VEqual(v1.SubNum(n), &Vector{0, 1, 2}) {
		t.Fail()
	}
}

func TestVector_MulNum(t *testing.T) {
	v1 := &Vector{1, 2, 3}
	n := 2
	if !VEqual(v1.MulNum(n), &Vector{2, 4, 6}) {
		t.Fail()
	}
}

func TestVector_Dot(t *testing.T) {
	v1 := &Vector{1, 2, 3, 4, 5, 6}
	v2 := &Vector{6, 5, 4, 3, 2, 1}
	if v1.Dot(v2) != 56 || v1.Dot(v2) != v2.Dot(v1) {
		t.Fail()
	}
}

func TestVector_OuterProduct(t *testing.T) {
	v1 := &Vector{1, 2, 3}
	v2 := &Vector{4, 5}
	a := Data{{1, 2}, {2, 4}, {3, 6}}
	if !MEqual(v1.OuterProduct(v2), new(Matrix).Init(a)) {
		t.Fail()
	}
}

func TestVector_Cross(t *testing.T) {
	v1 := &Vector{1, 2, 3}
	v2 := &Vector{5, 8, 6}
	if !VEqual(v1.Cross(v2), &Vector{-12, 9, -2}) {
		t.Fail()
	}
}

func TestVector_SquareSum(t *testing.T) {
	v1 := &Vector{1, 2, 3}
	if v1.SquareSum() != 14 {
		t.Fail()
	}
}

func TestVector_Norm(t *testing.T) {
	v1 := &Vector{1, 2, 3}
	if !FloatEqual(v1.Norm(), 3.741657387) {
		t.Fail()
	}
}

func TestVector_Normalize(t *testing.T) {
	v1 := &Vector{1, 2, 3}
	if !VEqual(v1.Normalize(), &Vector{0.2672612419124244, 0.5345224838248488, 0.8017837257372732}) {
		t.Fail()
	}
}

func TestVector_ToMatrix(t *testing.T) {
	v := &Vector{1, 2, 3, 4, 5, 6}
	m := Data{{1, 2, 3}, {4, 5, 6}}
	if !MEqual(v.ToMatrix(2, 3), new(Matrix).Init(m)) {
		t.Fail()
	}
}

func TestVector_Sum(t *testing.T) {
	v := &Vector{1, 2, 3, 4, 5, 6}
	if !FloatEqual(v.Sum(), 21) {
		t.Fail()
	}
}

func TestVector_AbsSum(t *testing.T) {
	v := &Vector{1, -2, 3, -4, 5, -6}
	if !FloatEqual(v.AbsSum(), 21) {
		t.Fail()
	}
}

func TestVector_Mean(t *testing.T) {
	v := &Vector{1, 2, 3, 4, 5, 6}
	if !FloatEqual(v.Mean(), 3.5) {
		t.Fail()
	}
}

func TestVector_Tile(t *testing.T) {
	v := &Vector{1, 2, 3}
	m := Data{{1, 2, 3}, {1, 2, 3}}
	n := Data{{1, 1}, {2, 2}, {3, 3}}
	if !MEqual(v.Tile(0, 2), new(Matrix).Init(m)) {
		t.Fail()
	}
	if !MEqual(v.Tile(1, 2), new(Matrix).Init(n)) {
		t.Fail()
	}
}

func TestVector_Length(t *testing.T) {
	v := &Vector{1, 2, 3}
	if v.Length() != 3 {
		t.Fail()
	}
}

func TestVector_Max(t *testing.T) {
	v := &Vector{1, 2, 3}
	idx, value := v.Max()
	if idx != 2 || value != v.At(2) {
		t.Fail()
	}
}

func TestVector_Min(t *testing.T) {
	v := &Vector{1, 2, 3}
	idx, value := v.Min()
	if idx != 0 || value != v.At(0) {
		t.Fail()
	}
}

func TestVector_Sorted(t *testing.T) {
	v := &Vector{1, 2, 3}
	sorted := v.Sorted()
	for i := 1; i < v.Length(); i++ {
		if sorted[i].value < sorted[i-1].value {
			t.Fail()
		}
	}
}

// Vector convolve
func TestConvolve(t *testing.T) {
	size := 10000
	u := GenerateRandomVector(size)
	v := GenerateRandomVector(size)

	res := Convolve(u, v)
	if len(*res) != size+size-1 {
		t.Fail()
	}
}

func TestVector_String(t *testing.T) {
	v := &Vector{1, 2, 3}
	if v.String() != "{1.000000, 2.000000, 3.000000}\n" {
		t.Fail()
	}
}

/*
BenchmarkMatrix_Mul/size-10-8     500000              3418 ns/op
BenchmarkMatrix_Mul/size-100-8       500           2845713 ns/op
BenchmarkMatrix_MulParaFor/size-10-8              200000             11248 ns/op
BenchmarkMatrix_MulParaFor/size-100-8               1000           2273291 ns/op
BenchmarkMatrix_MulParaFor/size-300-8                100          75323710 ns/op
*/
func BenchmarkMatrix_Mul(b *testing.B) {
	for k := 1; k <= 5; k++ {
		n := 0
		if k < 3 {
			n = int(math.Pow(10, float64(k)))
		} else {
			n = k * 100
		}
		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			m := GenerateRandomSquareMatrix(n)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				m.Mul(m)
			}
		})
	}
}

/*
BenchmarkNaivePow/size-10-8                30000             40736 ns/op
BenchmarkNaivePow/size-100-8                 100         212962377 ns/op
BenchmarkMatrix_Pow/size-10-8              30000             40647 ns/op
BenchmarkMatrix_Pow/size-100-8               100          22311996 ns/op
*/
func BenchmarkMatrix_Pow(b *testing.B) {
	for k := 1.0; k <= 2; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			m := GenerateRandomSquareMatrix(n)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				m.Pow(n)
			}
		})
	}
}

func BenchmarkMatrix_MulNum(b *testing.B) {
	for k := 1; k <= 10; k++ {
		n := 0
		if k < 3 {
			n = int(math.Pow(10, float64(k)))
		} else {
			n = k * 100
		}
		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			m := GenerateRandomSquareMatrix(n)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				m.MulNum(n)
			}
		})
	}
}

func BenchmarkVector_SquareSum(b *testing.B) {
	for k := 1.0; k <= 5; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			v := GenerateRandomVector(n)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				v.SquareSum()
			}
		})
	}
}

func BenchmarkMatrix_Rank(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			m := GenerateRandomSquareMatrix(n)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				m.Rank()
			}
		})
	}
}

/*
// BenchmarkMatrix_Det/size-10-8                      100        2121988644 ns/op
// too slow
func BenchmarkMatrix_Det(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			for i := 1; i < b.N; i++ {
				m := GenerateRandomSquareMatrix(n)
				m.Det()
			}
		})
	}
}
*/

func BenchmarkOneMatrix(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			for i := 1; i < b.N; i++ {
				OneMatrix(n, n)
			}
		})
	}
}

func BenchmarkConvolve(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			u := GenerateRandomVector(n)
			v := GenerateRandomVector(n)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				Convolve(u, v)
			}
		})
	}
}
