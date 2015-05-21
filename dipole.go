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

	return Vector{dBxdi, dBydi, dBzdi}
}
