package vinamax2

import "math"

// Field generated at position R relative to dipole m
func DipoleField(m, R Vector) Vector {

	r := R.Len()
	r2 := r * r
	r3 := r * r2
	r5 := r3 * r2

	return R.Mul(3 * m.Dot(R) / r5).Sub(m.Div(r3)).Div(4 * math.Pi)

}
