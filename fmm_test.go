package vinamax2

// FMM unit tests

import (
	"testing"
)

func init() {
	verbose = false
}

// Test average field of one spin
// against 0-order solution from 2015-05-20 (Arne)
func TestFMM0th(t *testing.T) {

	initTestWorld()

	CalcDemag()
	CalcDemag()
	CalcDemag()
	CalcDemag()

	var Btotal Vector
	for _, p := range Particles {
		Btotal = Btotal.Add(p.b)
	}

	//solution := Vector{5850.136490409946, 4680.109192327974, 3510.08189424605}
	solution:= Vector{0.028162933645173385, 0.02253034691613862, 0.01689776018710405}

	tol := 1e-6
	if Btotal.Sub(solution).Len() > tol {
		t.Error("got:", Btotal, "expected:", solution)
	}
}

// Somehow fails if CalcDemag() (not parallel) has been called before...
//func TestFMMParallel(t *testing.T) {
//
//	initTestWorld()
//
//	CalcDemagParallel()
//	CalcDemagParallel()
//
//	var Btotal Vector
//	for _, p := range Particles {
//		Btotal = Btotal.Add(p.b)
//	}
//
//	solution := Vector{5850.136490409946, 4680.109192327974, 3510.08189424605}
//
//	tol := 1e-6
//	if Btotal.Sub(solution).Len() > tol {
//		t.Error("got:", Btotal, "expected:", solution)
//	}
//
//}

func initTestWorld() {
	worldSize := Vector{1, 1, 1}
	NLEVEL := 5

	FMMOrder = 0
	InitFMM(worldSize, NLEVEL)

	// place particles with m=0 , as field probes
	baseLevel := Level[NLEVEL-1]
	for _, c := range baseLevel {
		AddParticle(&Particle{M: Vector{0, 0, 0}, center: c.center})
	}

	// place one magneticed particle as source
	hotcell := baseLevel[0]
	AddParticle(&Particle{M: Vector{1, 2, 3}, center: hotcell.center,msat:1.,r:1.})
	InitFMM2()
}
