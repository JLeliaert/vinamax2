//This second example is a test if the demagnetising field is implemented correctly
//To check this, we let 2 particles relax in the presence of an external field
//and check the output versus mumax. We also do the same simulation without
//calculating the demagnetising field to see if this problem is suited to
//check the implementation; i.e. to see that the demagnetising field
//makes a difference.

package main

import (
	. "github.com/JLeliaert/vinamax2"
)

func main() {

	World(0, 0, 0, 2e-6)

	Particle_radius(16e-9)

	Addsingleparticle(-64.48e-9, 0, 0)
	Addsingleparticle(64.48e-9, 0, 0)

	B_ext = func(t float64) (float64, float64, float64) { return 0.001, 0., 0.0 }

	FMM = true
	Demag = true

	//set the saturation magnetisation of the particles
	Msat(860e3)

	Dt = 1e-12
	T = 0.
	Temp = 0.0
	Alpha = 0.1
	Ku1 = 0
	Anisotropy_axis(0, 0, 1)

	M_uniform(0, 1, 0)

	Output(1e-10)

	Run(100.e-9)

}
