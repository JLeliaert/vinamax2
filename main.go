package main

import (
	"fmt"
	"log"
	"math"
)

func main() {

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

	// output one layer
	for _, c := range baseLevel {
		for _, p := range c.particle {
			r := p.center
			b := p.b.Div(p.b.Len()).Mul(c.size[X]) // normalize
			if r[Z] == -0.46875 {
				fmt.Println(r[X], r[Y], r[Z], b[X], b[Y], b[Z])
			}
		}
	}

	//for l := range level {
	//	fmt.Println("level", l)
	//	for _, c := range level[l] {
	//		fmt.Println(c)
	//	}
	//	fmt.Println()
	//}

}

func printStats() {
	nLeaf := int(math.Pow(8, float64(len(Level)-1)) + 0.5)
	log.Println(totalCells, "cells, avg", totalPartners/totalCells, "partners/cell, avg", totalNear/nLeaf, "near/leaf")
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
