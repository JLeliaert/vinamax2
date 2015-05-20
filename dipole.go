package main

import "math"

// Field generated at position R relative to dipole m
func DipoleField(m, R Vector) Vector {

	r := R.Len()
	r2 := r * r
	r3 := r * r2
	r5 := r3 * r2

	m_R := m.Dot(R)

	Bx := (1 / (4 * math.Pi)) * ((3 * m_R * R[X] / r5) - (m[X] / r3))
	By := (1 / (4 * math.Pi)) * ((3 * m_R * R[Y] / r5) - (m[Y] / r3))
	Bz := (1 / (4 * math.Pi)) * ((3 * m_R * R[Z] / r5) - (m[Z] / r3))

	return Vector{Bx, By, Bz}
}
