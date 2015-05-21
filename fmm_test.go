package vinamax2

import (
	"testing"
)

func init() {
	verbose = false
}

// Test average field of one spin
// against 0-order solution from 2015-05-20 (Arne)
func TestFMM(t *testing.T) {

	worldSize := Vector{1, 1, 1}
	NLEVEL := 5

	InitFMM(worldSize, NLEVEL)

	// place particles with m=0 , as field probes
	baseLevel := Level[NLEVEL-1]
	for _, c := range baseLevel {
		AddParticle(&Particle{M: Vector{0, 0, 0}, center: c.center})
	}

	// place one magneticed particle as source
	hotcell := baseLevel[0]
	AddParticle(&Particle{M: Vector{1, 2, 3}, center: hotcell.center})

	CalcDemag()

	var Btotal Vector
	for _, p := range Particles {
		Btotal = Btotal.Add(p.b)
	}

	solution := Vector{5850.136490409946, 4680.109192327974, 3510.08189424605}

	tol := 1e-6
	if Btotal.Sub(solution).Len() > tol {
		t.Error("got:", Btotal, "expected:", solution)
	}

}

func BenchmarkFMM5Levels(b *testing.B) {

	b.StopTimer()

	worldSize := Vector{1, 1, 1}
	NLEVEL := 5

	InitFMM(worldSize, NLEVEL)

	// place particles with m=0 , as field probes
	baseLevel := Level[NLEVEL-1]
	for _, c := range baseLevel {
		AddParticle(&Particle{M: Vector{0, 0, 0}, center: c.center})
	}

	// place one magneticed particle as source
	hotcell := baseLevel[0]
	AddParticle(&Particle{M: Vector{1, 2, 3}, center: hotcell.center})

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		CalcDemag()
	}
}

func BenchmarkAddPartnerFields(b *testing.B) {

	b.StopTimer()

	worldSize := Vector{1, 1, 1}
	NLEVEL := 5

	InitFMM(worldSize, NLEVEL)

	// place particles with m=0 , as field probes
	baseLevel := Level[NLEVEL-1]
	for _, c := range baseLevel {
		AddParticle(&Particle{M: Vector{0, 0, 0}, center: c.center})
	}

	// place one magneticed particle as source
	hotcell := baseLevel[30]
	AddParticle(&Particle{M: Vector{1, 2, 3}, center: hotcell.center})

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		hotcell.addPartnerFields()
	}
}
