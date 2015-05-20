package vinamax2

import "math"

// 3-component vector
type Vector [3]float64

// Returns a*v.
func (v Vector) Mul(a float64) Vector {
	return Vector{a * v[0], a * v[1], a * v[2]}
}

// point-wise multiplication of components
func (v Vector) Mul3(a Vector) Vector {
	return Vector{a[0] * v[0], a[1] * v[1], a[2] * v[2]}
}

func (v Vector) Abs() Vector {
	return Vector{math.Abs(v[0]), math.Abs(v[1]), math.Abs(v[2])}
}

// Returns (1/a)*v.
func (v Vector) Div(a float64) Vector {
	return v.Mul(1 / a)
}

// Returns a+b.
func (a Vector) Add(b Vector) Vector {
	return Vector{a[0] + b[0], a[1] + b[1], a[2] + b[2]}
}

// Returns a+s*b.
func (a Vector) MAdd(s float64, b Vector) Vector {
	return Vector{a[0] + s*b[0], a[1] + s*b[1], a[2] + s*b[2]}
}

// Returns a-b.
func (a Vector) Sub(b Vector) Vector {
	return Vector{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
}

// Returns the norm of v.
func (v Vector) Len() float64 {
	len2 := v.Dot(v)
	return math.Sqrt(len2)
}

// Returns the uniform norm of v
// (maximum of absolute values of components)
func (v Vector) MaxNorm() float64 {
	x := math.Abs(v[0])
	y := math.Abs(v[1])
	z := math.Abs(v[2])
	return math.Max(math.Max(x, y), z)
}

// Returns the dot (inner) product a.b.
func (a Vector) Dot(b Vector) float64 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

// Returns the cross (vector) product a x b
// in a right-handed coordinate system.
func (a Vector) Cross(b Vector) Vector {
	x := a[1]*b[2] - a[2]*b[1]
	y := a[2]*b[0] - a[0]*b[2]
	z := a[0]*b[1] - a[1]*b[0]
	return Vector{x, y, z}
}

const (
	X = 0
	Y = 1
	Z = 2
)
