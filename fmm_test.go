package vinamax2

import (
	"testing"
)

// Test average field of one spin
// against 0-order solution from 2015-05-20 (Arne)
func TestFMM(t *testing.T) {

	NLEVEL := 5

	InitFMM(NLEVEL)

	// place particles with m=0 , as field probes
	baseLevel := Level[NLEVEL-1]
	for _, c := range baseLevel {
		c.particle = []*Particle{&Particle{M: Vector{0, 0, 0}, center: c.center}}
	}

	// place one magneticed particle as source
	hotcell := baseLevel[0]
	hotcell.particle = []*Particle{&Particle{M: Vector{1, 2, 3}, center: hotcell.center}}

	// calc B demag
	Root.UpdateM()
	Root.UpdateB(nil)

	var Btotal Vector
	for _, c := range baseLevel {
		for _, p := range c.particle {
			Btotal = Btotal.Add(p.b)
		}
	}

	solution := Vector{5850.136490409946, 4680.109192327974, 3510.08189424605}

	tol := 1e-6
	if Btotal.Sub(solution).Len() > tol {
		t.Error("got:", Btotal, "expected:", solution)
	}

}
