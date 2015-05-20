package main

import (
	"fmt"
	"testing"
)

func TestFMM(t *testing.T) {

	NLEVEL := 5

	InitFMM(NLEVEL)

	// place particles with m=0 , as field probes
	baseLevel := Level[NLEVEL-1]
	for _, c := range baseLevel {
		c.particle = []*Particle{&Particle{m: Vector{0, 0, 0}, center: c.center}}
	}

	// place one magneticed particle as source
	hotcell := baseLevel[0]
	hotcell.particle = []*Particle{&Particle{m: Vector{1, 0, 0}, center: hotcell.center}}

	// calc B demag
	Root.UpdateM()
	Root.UpdateB(nil)

	var Btotal Vector
	for _, c := range baseLevel {
		for _, p := range c.particle {
			Btotal = Btotal.Add(p.b)
		}
	}

	fmt.Println(Btotal)

}
