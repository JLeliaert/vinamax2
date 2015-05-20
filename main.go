// +build ignore

package main

import (
	"flag"
	"fmt"
	"math"
)

func main() {

	flag.Parse()

	defer Cleanup()

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

	// output one layer
	for _, p := range Particles {
		r := p.center
		b := p.b.Div(p.b.Len()).Mul(1. / 16.) // normalize
		if r[Z] == -0.46875 {
			fmt.Println(r[X], r[Y], r[Z], b[X], b[Y], b[Z])
		}
	}

}

func checkNaNs(x Vector) {
	checkNaN(x[X])
	checkNaN(x[Y])
	checkNaN(x[Z])
}

func checkNaN(x float64) {
	if math.IsNaN(x) {
		panic("NaN")
	}
}
