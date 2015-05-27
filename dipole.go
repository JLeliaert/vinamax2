package vinamax2

import (
	"math"
	"fmt"
)

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

func DiffBxdx(m, R Vector) float64{
	x := R[X]
	m_x := m[X]
	r2 := R.Dot(R)
	pre := 1 / (4 * math.Pi) 
	
	return pre * (-15*x*x*m.Dot(R)/math.Pow(r2, 7./2.) + 3*x*m_x/math.Pow(r2, 5./2.) + 3*m_x*x/math.Pow(r2, 5./2.))+pre*3*m.Dot(R)/math.Pow(r2, 5./2.)
}

func DiffBxdy(m, R Vector) float64{
	x := R[X]
	y := R[Y]
	m_x := m[X]
	m_y := m[Y]
	r2 := R.Dot(R)
	pre := 1 / (4 * math.Pi) 
	
	return pre * (-15*x*y*m.Dot(R)/math.Pow(r2, 7./2.) + 3*y*m_x/math.Pow(r2, 5./2.) + 3*m_y*x/math.Pow(r2, 5./2.))
}

func DiffBxdz(m, R Vector) float64{
	x := R[X]
	z := R[Z]
	m_x := m[X]
	m_z := m[Z]
	r2 := R.Dot(R)
	pre := 1 / (4 * math.Pi) 
	
	return pre * (-15*x*z*m.Dot(R)/math.Pow(r2, 7./2.) + 3*z*m_x/math.Pow(r2, 5./2.) + 3*m_z*x/math.Pow(r2, 5./2.))
}



func DiffBydx(m, R Vector) float64{
	x := R[X]
	y := R[Y]
	m_x := m[X]
	m_y := m[Y]
	r2 := R.Dot(R)
	pre := 1 / (4 * math.Pi) 
	
	return pre * (-15*x*y*m.Dot(R)/math.Pow(r2, 7./2.) + 3*x*m_y/math.Pow(r2, 5./2.) + 3*m_x*y/math.Pow(r2, 5./2.))
}

func DiffBydy(m, R Vector) float64{
	y := R[Y]
	m_y := m[Y]
	r2 := R.Dot(R)
	pre := 1 / (4 * math.Pi) 
	
	return pre * (-15*y*y*m.Dot(R)/math.Pow(r2, 7./2.) + 3*y*m_y/math.Pow(r2, 5./2.) + 3*m_y*y/math.Pow(r2, 5./2.))+pre*3*m.Dot(R)/math.Pow(r2, 5./2.)
}

func DiffBydz(m, R Vector) float64{
	y := R[Y]
	z := R[Z]
	m_y := m[Y]
	m_z := m[Z]
	r2 := R.Dot(R)
	pre := 1 / (4 * math.Pi) 
	
	return pre * (-15*y*z*m.Dot(R)/math.Pow(r2, 7./2.) + 3*z*m_y/math.Pow(r2, 5./2.) + 3*m_z*y/math.Pow(r2, 5./2.))
}


func DiffBzdx(m, R Vector) float64{
	x := R[X]
	z := R[Z]
	m_x := m[X]
	m_z := m[Z]
	r2 := R.Dot(R)
	pre := 1 / (4 * math.Pi) 
	
	return pre * (-15*x*z*m.Dot(R)/math.Pow(r2, 7./2.) + 3*x*m_z/math.Pow(r2, 5./2.) + 3*m_x*z/math.Pow(r2, 5./2.))
}

func DiffBzdz(m, R Vector) float64{
	z := R[Z]
	m_z := m[Z]
	r2 := R.Dot(R)
	pre := 1 / (4 * math.Pi) 
	
	return pre * (-15*z*z*m.Dot(R)/math.Pow(r2, 7./2.) + 3*z*m_z/math.Pow(r2, 5./2.) + 3*m_z*z/math.Pow(r2, 5./2.))+pre*3*m.Dot(R)/math.Pow(r2, 5./2.)
}

func DiffBzdy(m, R Vector) float64{
	y := R[Y]
	z := R[Z]
	m_y := m[Y]
	m_z := m[Z]
	r2 := R.Dot(R)
	pre := 1 / (4 * math.Pi) 
	
	return pre * (-15*y*z*m.Dot(R)/math.Pow(r2, 7./2.) + 3*z*m_y/math.Pow(r2, 5./2.) + 3*m_z*y/math.Pow(r2, 5./2.))
}



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
	r2 := R.Dot(R)
	pre := 1 / (4 * math.Pi)
	//	return pre * (3*x*m.Dot(R)/math.Pow(r2, 5./2.) - m[i]/math.Pow(r2, 3./2.))
	fmt.Println("diag",pre*3*m.Dot(R)/math.Pow(r2, 5./2.))
	fmt.Println("offdiag",doff(i, i, m, R))
	return pre*3*m.Dot(R)/math.Pow(r2, 5./2.) + doff(i, i, m, R)
}

// dBx/dy = -15xy mdotR / r^7 + 3y m_x/r^5 + 3m_y x / r^5
func doff(c, i int, m, R Vector) float64 {
	x := R[c]
	y := R[i]
	m_x := m[c]
	m_y := m[i]
	r2 := R.Dot(R)
	pre := 1 / (4 * math.Pi) // minus sign gives best accuracy but should not be here, what's wrong?
	return pre * (-15*x*y*m.Dot(R)/math.Pow(r2, 7./2.) + 3*y*m_x/math.Pow(r2, 5./2.) + 3*m_y*x/math.Pow(r2, 5./2.))
}
