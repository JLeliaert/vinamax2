// +build ignore

package main

import (
	"flag"
	"fmt"
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
	msat := 1.

	//place particles with m=0 as field probes
		for i :=-0.5;i<0.5;i+=0.05{
		for j :=-0.5;j<0.5;j+=0.05{
		M := Vector{0, 0, 0}
			AddParticle(NewParticle(Vector{i,j,-0.03125}, r, M, msat))
}}
	r = 1e-9
	msat = 1e6

	// place one magneticed particle as source
	M := Vector{1, 0, 0}
	AddParticle(NewParticle(Vector{-0.03125, -0.03125, -0.03125}, r, M, msat))

	InitFMM2()
	//for _, c := range Level[1] {
	//	fmt.Println(c.Center(), c.Moment(), c.CenterOfMag())
	//}

	FMMOrder = 1
	Log("Order:", FMMOrder, " Proxy:", Proximity)
	for i := 0; i < 1; i++ {
		println(i)
		CalcDemagParallel()
	}
	Log("Demag error:", DemagError())

	// output one layer
	for _, p := range Particles {
		r := p.Center()
		b := p.Bdemag()
		b = b.Div(b.Len()).Mul(1. / 16.) // normalize
		if r[Z] == -0.03125 {
			fmt.Println(r[X], r[Y], r[Z], b[X], b[Y], b[Z])
		}
	}
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
