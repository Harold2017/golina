package golina

import (
	"math"
	"runtime"
	"sync"
)

// init function to set CPU usage
func init() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus) // Try to use all available CPUs.
}

// EPS for float number comparision
const EPS float64 = 1E-6

type Vector []float64 // 1D array
type Data []Vector    // 2D array -> backend of Matrix

// matrix entry
type Entry struct {
	value    float64
	row, col int
}

// _Matrix interface

type _Matrix interface {
	// dimensions
	Dims() (row, col int)

	// value at index(row i, col j), panic if not access
	At(i, j int) float64

	// set value at index(row i, col j), panic if not access
	Set(i, j int, value float64)

	// transpose matrix
	T() _Matrix

	// get row
	Row(i int) *Vector

	// get column
	Col(i int) *Vector

	// max entry
	Max() *Entry

	// min entry
	Min() *Entry

	// rank
	Rank() int
}

// Transpose struct implementing _Matrix interface and return transpose of input _Matrix
type Matrix struct {
	_Matrix
	_array Data // row-wise
}

func (t *Matrix) Init(array Data) *Matrix {
	return &Matrix{_array: array}
}

func (t *Matrix) Dims() (row, col int) {
	return len(t._array), len(t._array[0])
}

func (t *Matrix) At(i, j int) float64 {
	return t._array[i][j]
}

func (t *Matrix) Set(i, j int, value float64) {
	t._array[i][j] = value
}

func (t *Matrix) T() *Matrix {
	row, col := t.Dims()
	ntArray := make(Data, col)
	for i := 0; i < col; i++ {
		ntArray[i] = make([]float64, row)
		for j := 0; j < row; j++ {
			ntArray[i][j] = t._array[j][i]
		}
	}
	nt := new(Matrix).Init(ntArray)
	return nt
}

func (t *Matrix) Row(m int) *Vector {
	row, _ := t.Dims()
	if m > -1 && m < row {
		return &t._array[m]
	}
	panic("row index out of range")
}

func (t *Matrix) Col(n int) *Vector {
	_, col := t.Dims()
	if n > -1 && n < col {
		return &t.T()._array[n]
	}
	panic("column index out of range")
}

func VEqual(v1, v2 *Vector) bool {
	if len(*v1) != len(*v2) {
		return false
	}
	for i, v := range *v1 {
		if v != (*v2)[i] && v-(*v2)[i] > EPS {
			return false
		}
	}
	return true
}

// https://stackoverflow.com/questions/37884152/how-do-i-check-the-equality-of-three-values-elegantly
func Equal(mat1, mat2 *Matrix) bool {
	row1, col1 := mat1.Dims()
	row2, col2 := mat2.Dims()
	if [2]int{row1, col1} != [2]int{row2, col2} {
		return false
	}
	for i, col := range mat1._array {
		if !VEqual(&col, mat2.Row(i)) {
			return false
		}
	}
	return true
}

func Copy(t *Matrix) *Matrix {
	nt := Matrix{_array: make([]Vector, len(t._array))}
	for i := range t._array {
		nt._array[i] = make(Vector, len(t._array[i]))
		copy(nt._array[i], t._array[i])
	}
	return &nt
}

func Empty(t *Matrix) *Matrix {
	row, col := t.Dims()
	nt := Matrix{_array: make([]Vector, row)}
	for i := range t._array {
		nt._array[i] = make(Vector, col)
	}
	return &nt
}

// nil entries
func EmptyMatrix(row, col int) *Matrix {
	nt := Matrix{_array: make([]Vector, row)}
	for i := range nt._array {
		nt._array[i] = make(Vector, col)
	}
	return &nt
}

func ZeroMatrix(row, col int) *Matrix {
	nt := EmptyMatrix(row, col)
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			nt.Set(i, j, 0)
		}
	}
	return nt
}

func OneMatrix(row, col int) *Matrix {
	nt := EmptyMatrix(row, col)
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			nt.Set(i, j, 1)
		}
	}
	return nt
}

func IdentityMatrix(n int) *Matrix {
	nt := ZeroMatrix(n, n)
	for i := 0; i < n; i++ {
		nt.Set(i, i, 1)
	}
	return nt
}

// TODO: how to find all max or min?
func (t *Matrix) Max() *Entry {
	entry := Entry{}
	entry.value = math.Inf(-1)
	for r, i := range t._array {
		for c, j := range i {
			if j > entry.value {
				entry.value = j
				entry.row, entry.col = r, c
			}
		}
	}
	return &entry
}

func (t *Matrix) Min() *Entry {
	entry := Entry{}
	entry.value = math.Inf(1)
	for r, i := range t._array {
		for c, j := range i {
			if j < entry.value {
				entry.value = j
				entry.row, entry.col = r, c
			}
		}
	}
	return &entry
}

// TODO: need optimize
// Gaussian elimination (row echelon form)
func (t *Matrix) Rank() (rank int) {
	mat := Copy(t)
	rowN, colN := mat.Dims()
	rank = colN
	for row := 0; row < rank; row++ {
		// diagonal entry is not zero
		if mat.At(row, row) != 0 {
			for col := 0; col < rowN; col++ {
				if col != row {
					// makes all entries of current column as 0 except entry `mat[row][row]`
					multipler := mat.At(col, row) / mat.At(row, row)
					for i := 0; i < rank; i++ {
						mat.Set(col, i, mat.At(col, i)-multipler*mat.At(row, i))
					}
				}
			}
		} else {
			// diagonal entry is already zero, now two cases
			// 1) if there is a row below it with non-zero entry, then swap this row with that row and process that row
			// 2) if all elements in current column below mat[row][row] are 0,
			// 	  then remove this column by swapping it with last column and reducing rank by 1
			reduce := true

			for i := row + 1; i < rowN; i++ {
				// swap the row with non-zero entry with this row
				if mat.At(i, row) > EPS {
					SwapRow(mat, row, i)
					reduce = false
					break
				}
			}

			// if no row with non-zero entry in current column, then all values in this column are 0
			if reduce {
				// reduce rank
				rank--
				// copy the last column here
				for i := 0; i < rowN; i++ {
					mat.Set(i, row, mat.At(i, rank))
				}
			}

			// process this row again
			row--
		}
	}
	return rank
}

func SwapRow(t *Matrix, row1, row2 int) {
	t._array[row1], t._array[row2] = *t.Row(row2), *t.Row(row1)
}

// TODO: need optimize
// Determinant of N x N matrix recursively
func (t *Matrix) Det() float64 {
	row, col := t.Dims()
	if row != col {
		panic("need N x N matrix for determinant calculation")
	}
	return _det(t, row)
}
func _det(t *Matrix, n int) float64 {
	det := 0.

	// base case: if matrix only contains one entry
	if n == 1 {
		return t.At(0, 0)
	}

	// template matrix to store coefficients
	matTmp := Empty(t)
	// sign of multiplier
	sign := 1.
	// iterate for each entry of first row
	for f := 0; f < n; f++ {
		// get coefficient of mat[0][f]
		getCoeff(t, matTmp, 0, f, n)
		det += sign * t.At(0, f) * _det(matTmp, n-1)
		sign = -sign
	}
	return det
}

// func to get coefficients of mat[p][q] in matTmp, n is dimension of current matrix (to avoid re-calculation)
func getCoeff(t, matTmp *Matrix, p, q, n int) {
	i, j := 0, 0

	// looping for each entries of the matrix
	for row := 0; row < n; row++ {
		for col := 0; col < n; col++ {
			// fill template matrix
			if row != p && col != q {
				matTmp.Set(i, j, t.At(row, col))
				j++
				// row is filled, so increase row index and rest col index
				if j == n-1 {
					j = 0
					i++
				}
			}
		}
	}
}

// Adjugate Matrix
// https://en.wikipedia.org/wiki/Adjugate_matrix
func (t *Matrix) Adj() (adj *Matrix) {
	row, col := t.Dims()
	if row != col {
		panic("need N x N matrix for adjugate calculation")
	}

	adj = Empty(t)

	if row == 1 {
		adj.Set(0, 0, 1)
		return
	}

	// temp to store coefficients
	matTmp := Empty(t)
	sign := 1.

	for i := 0; i < row; i++ {
		for j := 0; j < row; j++ {
			// get coefficient of t[i][j]
			getCoeff(t, matTmp, i, j, row)
			// sign of adj[j][i] is positive if sum of row and column indexes is even
			sign = Ternary((i+j)%2 == 0, 1., -1.).(float64)
			// interchanging rows and columns to get transpose
			adj.Set(j, i, sign*_det(matTmp, row-1))
		}
	}
	return
}

// Inverse Matrix
// inverse(t) = adj(t) / det(t)
func (t *Matrix) Inverse() *Matrix {
	det := t.Det()
	if det == 0 {
		panic("this matrix is not invertible")
	}
	adj := t.Adj()
	inverse := Empty(t)
	n, _ := t.Dims()
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			inverse.Set(i, j, adj.At(i, j)/det)
		}
	}
	return inverse
}

func (t *Matrix) Add(mat2 *Matrix) *Matrix {
	row1, col1 := t.Dims()
	row2, col2 := mat2.Dims()
	if [2]int{row1, col1} != [2]int{row2, col2} {
		panic("both matrices should have the same dimension")
	}
	nt := Empty(t)
	for r, i := range t._array {
		for c, j := range i {
			nt.Set(r, c, j+mat2.At(r, c))
		}
	}
	return nt
}

func (t *Matrix) Sub(mat2 *Matrix) *Matrix {
	row1, col1 := t.Dims()
	row2, col2 := mat2.Dims()
	if [2]int{row1, col1} != [2]int{row2, col2} {
		panic("both matrices should have the same dimension")
	}
	nt := Empty(t)
	for r, i := range t._array {
		for c, j := range i {
			nt.Set(r, c, j-mat2.At(r, c))
		}
	}
	return nt
}

// TODO: need optimize
// Matrix multiplication (dot | inner)
// https://en.wikipedia.org/wiki/Matrix_multiplication
func (t *Matrix) Mul(mat2 *Matrix) *Matrix {
	row1, col1 := t.Dims()
	row2, col2 := mat2.Dims()
	if col1 != row2 {
		panic("matrix multiplication need M x N and N x L matrices to get M x L matrix")
	}
	out := ZeroMatrix(row1, col2)
	for i := 0; i < row1; i++ {
		for j := 0; j < col2; j++ {
			for k := 0; k < row2; k++ {
				out.Set(i, j, out.At(i, j)+t.At(i, k)*mat2.At(k, j))
			}
		}
	}
	return out
}

func (t *Matrix) MulNum(n interface{}) *Matrix {
	multiplier := getFloat64(n)
	row, col := t.Dims()
	out := ZeroMatrix(row, col)
	for i := range t._array {
		for j, v := range t._array[i] {
			out.Set(i, j, v*multiplier)
		}
	}
	return out
}

func getFloat64(x interface{}) float64 {
	switch x := x.(type) {
	case uint8:
		return float64(x)
	case int8:
		return float64(x)
	case uint16:
		return float64(x)
	case int16:
		return float64(x)
	case uint32:
		return float64(x)
	case int32:
		return float64(x)
	case uint64:
		return float64(x)
	case int64:
		return float64(x)
	case int:
		return float64(x)
	case float32:
		return float64(x)
	case float64:
		return x
	}
	panic("invalid numeric type of input")
}

// TODO: need optimize and deal with negative condition (invertible)
// Matrix Power of square matrix
// Precondition: n >= 0
func (t *Matrix) Pow(n int) *Matrix {
	row, col := t.Dims()
	if row != col {
		panic("only square matrix has power")
	}
	if n == 0 {
		return IdentityMatrix(row)
	} else if n == 1 {
		return t
	} else {
		nt := Copy(t)
		for n > 1 {
			nt = nt.Mul(Copy(t))
			n--
		}
		return nt
	}
}

// trace: sum of all diagonal values
// https://en.wikipedia.org/wiki/Trace_(linear_algebra)
func (t *Matrix) Trace() float64 {
	row, col := t.Dims()
	if row != col {
		panic("square matrix only")
	}
	res := 0.
	for i := range t._array {
		res += t._array[i][i]
	}
	return res
}

// Eigenvalues and Eigenvectors
// https://en.wikipedia.org/wiki/Eigenvalue_algorithm
// this is for 3 x 3 symmetric matrix only
// TODO: general case
func Eigen(t *Matrix) (eig_val *Vector, eig_vec *Matrix) {
	// Eigenvalues
	eig_val = EigenValues(t)

	// Eigenvectors
	// eig_vec.Sub(IdentityMatrix(3).MulNum((*eig_val)[1]))
	if math.Pow(t._array[0][1], 2.)+math.Pow(t._array[0][2], 2.)+math.Pow(t._array[1][2], 2.) == 0 {
		eig_vec = IdentityMatrix(3)
		return
	}
	eig_vec = EigenVector(t, eig_val)
	return
}

func EigenValues(t *Matrix) (eig_val *Vector) {
	// Eigenvalues
	eig0, eig1, eig2 := 0., 0., 0.
	// upper triangle
	p1 := math.Pow(t._array[0][1], 2.) + math.Pow(t._array[0][2], 2.) + math.Pow(t._array[1][2], 2.)
	if p1 == 0 {
		// t is diagonal
		eig0 = t._array[0][0]
		eig1 = t._array[1][1]
		eig2 = t._array[2][2]
	} else {
		q := t.Trace() / 3
		p2 := math.Pow(t._array[0][0]-q, 2.) + math.Pow(t._array[1][1]-q, 2.) + math.Pow(t._array[2][2]-q, 2.) + 2*p1
		p := math.Sqrt(p2 / 6)
		B := t.Sub(IdentityMatrix(3).MulNum(q)).MulNum(1 / p)
		r := B.Det() / 2
		// in exact arithmetic for a symmetric matrix: -1 <= r <= 1
		// but computation error can leave it slightly outside this range
		phi := 0.
		if r <= -1 {
			phi = math.Pi / 3
		} else if r >= 1 {
			phi = 0.
		} else {
			phi = math.Acos(r) / 3
		}
		// eigenvalues satisfy eig2 <= eig1 <= eig0
		eig0 = q + 2*p*math.Cos(phi)
		eig2 = q + 2*p*math.Cos(phi+(2*math.Pi/3))
		eig1 = 3*q - eig0 - eig2 // since t.Trace() = eig0 + eig1 + eig2
	}
	eig_val = &Vector{eig0, eig1, eig2}
	return
}

func EigenVector(t *Matrix, eig_val *Vector) (eig_vec *Matrix) {
	eig_vec = ZeroMatrix(3, 3)
	// algebraic multiplicity 1
	eig_vec._array[0] = *computeEigenVector0(Copy(t), (*eig_val)[0])
	// algebraic multiplicity 2
	eig_vec._array[1] = *computeEigenVector1(Copy(t), &eig_vec._array[0], (*eig_val)[1])
	// TODO: here multiple - 1 to keep align with numpy result, but i think the sign does not matter, an i right?
	eig_vec._array[2] = *(eig_vec.Row(0).Cross(eig_vec.Row(1)).MulNum(-1))
	return
}

// A Robust Eigensolver for 3 3 Symmetric Matrices
func computeEigenVector0(t *Matrix, val0 float64) (vec0 *Vector) {
	// Move RHS to LHS
	t.Set(0, 0, t.At(0, 0)-val0)
	t.Set(1, 1, t.At(1, 1)-val0)
	t.Set(2, 2, t.At(2, 2)-val0)

	r0r1 := t.Row(0).Cross(t.Col(1))
	r0r2 := t.Row(0).Cross(t.Col(2))
	r1r2 := t.Col(1).Cross(t.Col(2))

	d0 := r0r1.Dot(r0r1)
	d1 := r0r2.Dot(r0r2)
	d2 := r1r2.Dot(r1r2)

	dmax := d0
	imax := 0

	if d1 > dmax {
		dmax = d1
		imax = 1
	}
	if d2 > dmax {
		imax = 2
	}
	if imax == 0 {
		vec0 = r0r1.MulNum(1 / math.Sqrt(d0))
	} else if imax == 1 {
		vec0 = r0r2.MulNum(1 / math.Sqrt(d1))
	} else {
		vec0 = r1r2.MulNum(1 / math.Sqrt(d2))
	}
	return
}

func ComputeOrthogonalComplement(W *Vector) (U, V *Vector) {
	invLength := 0.
	if math.Abs((*W)[0]) > math.Abs((*W)[1]) {
		invLength = 1 / math.Sqrt(math.Pow((*W)[0], 2.)+math.Pow((*W)[2], 2.))
		U = &Vector{-(*W)[2] * invLength, 0, (*W)[0] * invLength}
	} else {
		invLength = 1 / math.Sqrt(math.Pow((*W)[1], 2.)+math.Pow((*W)[2], 2.))
		U = &Vector{0, (*W)[2] * invLength, -(*W)[1] * invLength}
	}
	V = W.Cross(U)
	return
}

func computeEigenVector1(t *Matrix, vec0 *Vector, val1 float64) (vec1 *Vector) {
	// compute a right-handed orthonormal set {U, V, vec0}
	U, V := ComputeOrthogonalComplement(vec0)
	AU := &Vector{
		t.Row(0).Dot(U),
		t.Col(1).Dot(U),
		t.Col(2).Dot(U),
	}
	AV := &Vector{
		t.Row(0).Dot(V),
		t.Col(1).Dot(V),
		t.Col(2).Dot(V),
	}

	m00 := U.Dot(AU) - val1
	m01 := U.Dot(AV)
	m11 := V.Dot(AV) - val1

	absM00 := math.Abs(m00)
	absM01 := math.Abs(m01)
	absM11 := math.Abs(m11)
	maxAbsComp := 0.

	if absM00 >= absM11 {
		maxAbsComp = math.Max(absM00, absM01)
		if maxAbsComp > 0 {
			if absM00 >= absM01 {
				m01 /= m00
				m00 = 1 / math.Sqrt(1+math.Pow(m01, 2.))
				m01 *= m00
			} else {
				m00 /= m01
				m01 = 1 / math.Sqrt(1+math.Pow(m00, 2.))
				m00 *= m01
			}
			vec1 = U.MulNum(m01).Sub(V.MulNum(m00))
		} else {
			vec1 = U
		}
	} else {
		maxAbsComp = math.Max(absM11, absM01)
		if maxAbsComp > 0 {
			if absM11 >= absM01 {
				m01 /= m11
				m11 = 1 / math.Sqrt(1+math.Pow(m01, 2.))
				m01 *= m11
			} else {
				m11 /= m01
				m01 = 1 / math.Sqrt(1+math.Pow(m11, 2.))
				m11 *= m01
			}
			vec1 = U.MulNum(m11).Sub(V.MulNum(m01))
		} else {
			vec1 = U
		}
	}
	return
}

// Vector
func (v *Vector) Add(v1 *Vector) *Vector {
	if len(*v) != len(*v1) {
		panic("dot product requires equal-length vectors")
	}
	res := make(Vector, len(*v))
	for i := range *v {
		res[i] = (*v)[i] + (*v1)[i]
	}
	return &res
}

func (v *Vector) AddNum(n interface{}) *Vector {
	res := make(Vector, len(*v))
	for i := range *v {
		res[i] = (*v)[i] + getFloat64(n)
	}
	return &res
}

func (v *Vector) Sub(v1 *Vector) *Vector {
	if len(*v) != len(*v1) {
		panic("dot product requires equal-length vectors")
	}
	res := make(Vector, len(*v))
	for i := range *v {
		res[i] = (*v)[i] - (*v1)[i]
	}
	return &res
}

func (v *Vector) SubNum(n interface{}) *Vector {
	res := make(Vector, len(*v))
	for i := range *v {
		res[i] = (*v)[i] - getFloat64(n)
	}
	return &res
}

func (v *Vector) MulNum(n interface{}) *Vector {
	res := make(Vector, len(*v))
	for i := range *v {
		res[i] = (*v)[i] * getFloat64(n)
	}
	return &res
}

func (v *Vector) Dot(v1 *Vector) float64 {
	if len(*v) != len(*v1) {
		panic("dot product requires equal-length vectors")
	}
	res := 0.
	for i := range *v {
		res += (*v)[i] * (*v1)[i]
	}
	return res
}

// only for 3d
func (v *Vector) Cross(v1 *Vector) *Vector {
	if len(*v) != len(*v1) || len(*v) != 3 {
		panic("cross product requires 2d or 3d vectors")
	}
	return &Vector{(*v)[1]*(*v1)[2] - (*v)[2]*(*v1)[1], (*v)[2]*(*v1)[0] - (*v)[0]*(*v1)[2], (*v)[0]*(*v1)[1] - (*v)[1]*(*v1)[0]}
}

func (v *Vector) SquareSum() float64 {
	res := 0.
	for i := range *v {
		res += math.Pow((*v)[i], 2.) // TODO: v[i].Dot(v[i]) ? which one is faster?
	}
	return res
}

func (v *Vector) Norm() *Vector {
	ss := v.SquareSum()
	if ss == 0 {
		panic("invalid input vector with square sum equal to 0")
	}
	res := make(Vector, len(*v))
	for i := range *v {
		res[i] = (*v)[i] / math.Sqrt(ss)
	}
	return &res
}

// Vector convolve

func Ternary(statement bool, a, b interface{}) interface{} {
	if statement {
		return a
	}
	return b
}

func min(x, y int) int {
	return Ternary(x > y, y, x).(int)
}

func max(x, y int) int {
	return Ternary(x > y, x, y).(int)
}

func mul(u, v *Vector, k int) (res float64) {
	n := min(k+1, len(*u))
	j := min(k, len(*v)-1)

	for i := k - j; i < n; i, j = i+1, j-1 {
		res += (*u)[i] * (*v)[j]
	}
	return res
}

// Convolve computes w = u * v, where w[k] = Σ u[i]*v[j], i + j = k.
// Precondition: len(u) > 0, len(v) > 0.
func Convolve(u, v *Vector) *Vector {
	n := len(*u) + len(*v) - 1
	w := make(Vector, n)

	// Divide w into work units that take ~100μs-1ms to compute.
	size := max(1, 100000/n)

	var wg sync.WaitGroup
	for i, j := 0, size; i < n; i, j = j, j+size {
		if j > n {
			j = n
		}

		// The goroutines share memory, but only for reading.
		wg.Add(1)

		go func(i, j int) {
			for k := i; k < j; k++ {
				w[k] = mul(u, v, k)
			}
			wg.Done()
		}(i, j)
	}

	wg.Wait()

	return &w
}
