// +build ignore

package main

import (
	"flag"
	"math"

	. "."
)

var (
	flagProf    = flag.Bool("prof", false, "turn on CPU profiling")
	flagMemProf = flag.Bool("memprof", false, "turn on memory profiling")
)

func main() {
	defer Cleanup()

	flag.Parse()

	if *flagProf {
		InitCPUProfile()
	}
	if *flagMemProf {
		InitMemProfile()
	}

	worldSize := Vector{1, 1, 1}
	NLEVEL := 5

	Proximity = 1.1
	InitFMM(worldSize, NLEVEL)

	r := 1e-9
	msat := 1e6

	// place particles with m=0 , as field probes
	baseLevel := Level[NLEVEL-1]
	for _, c := range baseLevel {
		M := Vector{}
		AddParticle(NewParticle(c.Center(), r, M, msat))
	}

	// place one magneticed particle as source
	hotcell := baseLevel[0]
	M := Vector{1, 0, 0}
	AddParticle(NewParticle(hotcell.Center(), r, M, msat))

	FMMOrder = 0
	Log("Order:", FMMOrder, " Proxy:", Proximity)
	for i := 0; i < 500; i++ {
		CalcDemag()
	}
	//Log("Demag error:", DemagError())

	// output one layer
	//for _, p := range Particles {
	//	r := p.Center()
	//	//b := p.Bdemag()
	//	b := p.BruteDemag()
	//	b = b.Div(b.Len()).Mul(1. / 16.) // normalize
	//	if r[Z] == -0.46875 {
	//		fmt.Println(r[X], r[Y], r[Z], b[X], b[Y], b[Z])
	//	}
	//}

	Log("#Dipole evaluations:", NEvals)

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
