package vinamax2

import "math"

var NEvals int

// Field generated at position R relative to dipole m
func DipoleField(m, R Vector) Vector {

	NEvals++

	i_r2 := 1 / (R[X]*R[X] + R[Y]*R[Y] + R[Z]*R[Z])
	i_r := math.Sqrt(i_r2)
	i_r3 := i_r * i_r2
	i3_r5 := 3 * i_r3 * i_r2

	m_R := m.Dot(R)

	// TODO: use vectors
	Bx := (1 / (4 * math.Pi)) * ((i3_r5 * m_R * R[X]) - (i_r3 * m[X]))
	By := (1 / (4 * math.Pi)) * ((i3_r5 * m_R * R[Y]) - (i_r3 * m[Y]))
	Bz := (1 / (4 * math.Pi)) * ((i3_r5 * m_R * R[Z]) - (i_r3 * m[Z]))

	return Vector{Bx, By, Bz}
}

// Partial derivative (dB/di) of field generated at position R relative to dipole m.
// Direction of derivative: i = X,Y or Z.
func DiffDipole(i int, m, R Vector) Vector {
	r := R.Len()
	r2 := r * r
	r3 := r * r2
	r5 := r3 * r2
	r7 := r5 * r2
	mdotr :=m.Dot(R)

	
	dBxdi:= 3*R[X]*m[i]/r5-15*mdotr*R[X]*R[i]/r7-3*m[i]*R[i]/r5
	if(i==X){dBxdi+=3*(mdotr)/r5}

	dBydi:=3*R[Y]*m[i]/r5-15*mdotr*R[Y]*R[i]/r7-3*m[i]*R[i]/r5
	if(i==Y){dBydi+=3*(mdotr)/r5}

	dBzdi:=3*R[Z]*m[i]/r5-15*mdotr*R[Z]*R[i]/r7-3*m[i]*R[i]/r5
	if(i==Z){dBzdi+=3*(mdotr)/r5}

	return Vector{dBxdi, dBydi, dBzdi}.Mul(1 / (4 * math.Pi))
}
