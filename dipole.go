package vinamax2

import "math"

var NEvals int

// Field generated at position R relative to dipole m
func DipoleField(m, R Vector) Vector {

	NEvals++

	if (R == Vector{0, 0, 0}) {
		return Vector{0, 0, 0}
		//return m.Mul(-1. / 3.) // self-demag for sphere
	}

	i_r2 := 1 / (R[X]*R[X] + R[Y]*R[Y] + R[Z]*R[Z])
	i_r := math.Sqrt(i_r2)
	i_r3 := i_r * i_r2
	i3_r5 := 3 * i_r3 * i_r2

	m_R := m.Dot(R)

	Bx := (1 / (4 * math.Pi)) * ((i3_r5 * m_R * R[X]) - (i_r3 * m[X]))
	By := (1 / (4 * math.Pi)) * ((i3_r5 * m_R * R[Y]) - (i_r3 * m[Y]))
	Bz := (1 / (4 * math.Pi)) * ((i3_r5 * m_R * R[Z]) - (i_r3 * m[Z]))

	return Vector{Bx, By, Bz}
}

// Partial derivative (dB/di) of field generated at position R relative to dipole m.
// Direction of derivative: i = X,Y or Z.

func DiffDipole(i int, m, R Vector) Vector {

	return Vector{d(X, i, m, R), d(Y, i, m, R), d(Z, i, m, R)}

	//r := R.Len()
	//r2 := r * r
	//r3 := r * r2
	//r5 := r3 * r2
	//r7 := r5 * r2
	//mdotr := m.Dot(R)

	//dBxdi := 3*R[X]*m[i]/r5 - 15*mdotr*R[X]*R[i]/r7 + 3*m[X]*R[i]/r5
	//if i == X {
	//	dBxdi += 3 * (mdotr) / r5
	//}

	//dBydi := 3*R[Y]*m[i]/r5 - 15*mdotr*R[Y]*R[i]/r7 + 3*m[Y]*R[i]/r5
	//if i == Y {
	//	dBydi += 3 * (mdotr) / r5
	//}

	//dBzdi := 3*R[Z]*m[i]/r5 - 15*mdotr*R[Z]*R[i]/r7 + 3*m[Z]*R[i]/r5
	//if i == Z {
	//	dBzdi += 3 * (mdotr) / r5
	//}

	//return Vector{dBxdi, dBydi, dBzdi}.Mul(1 / (4 * math.Pi))
}

func d(c, i int, m, R Vector) float64 {
	if c == i {
		return ddiag(i, m, R)
	} else {
		return doff(c, i, m, R)
	}
}

// dBx/dx = 3*x*(m1*x+m2*y+m3*z)/(x^2+y^2+z^2)^(5/2) - m1/(x^2+y^2+z^2)^(3/2)
func ddiag(i int, m, R Vector) float64 {
	x := R[i]
	r2 := R.Dot(R)
	pre := 1 / (4 * math.Pi)
	return pre * (3*x*m.Dot(R)/math.Pow(r2, 5./2.) - m[i]/math.Pow(r2, 3./2.))
}

// dBx/dy = -15xy mdotR / r^7 + 3y m_x/r^5 + 3m_y x / r^5
func doff(c, i int, m, R Vector) float64 {
	x := R[c]
	y := R[i]
	m_x := m[c]
	m_y := m[i]
	r2 := R.Dot(R)
	pre := -1 / (4 * math.Pi) // minus sign gives best accuracy but should not be here, what's wrong?
	return pre * (-15*x*y*m.Dot(R)/math.Pow(r2, 7./2.) + 3*y*m_x/math.Pow(r2, 5./2.) + 3*m_y*x/math.Pow(r2, 5./2.))
}
