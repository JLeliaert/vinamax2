//This second example is a test if the demagnetising field is implemented correctly
//To check this, we let 2 particles relax in the presence of an external field
//and check the output versus mumax. We also do the same simulation without
//calculating the demagnetising field to see if this problem is suited to
//check the implementation; i.e. to see that the demagnetising field
//makes a difference.

package main

import (
	. ".."
)

func main() {

	worldsize := 2e-6
	nLevels := 5
	InitFMM(Vector{worldsize, worldsize, worldsize}, nLevels)

	r := 16e-9
	M := Vector{0, 1, 0}
	msat := 860e3

	p1 := NewParticle(Vector{-64.48e-9, 0, 0}, r, M, msat)
	p2 := NewParticle(Vector{64.48e-9, 0, 0}, r, M, msat)

	AddParticle(p1)
	AddParticle(p2)

	B_ext = func(t float64) (float64, float64, float64) { return 0.001, 0., 0.0 }

	//	FMM = true
	//	Demag = true
	//
	//	//set the saturation magnetisation of the particles
	//	Msat()

	Dt = 1e-12
	T = 0.
	Temp = 0.0
	Alpha = 0.1
	Ku1 = 0
	Anisotropy_axis(0, 0, 1)

	Output(1e-10)

	Run(100.e-9)

}
