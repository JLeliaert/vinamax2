package vinamax2

import (
	"fmt"
	"math"
	"time"
)

// FMM globals
var (
	Root      Cell        // roots the entire FMM tree
	Level     [][]*Cell   // for each level of the FMM tree: all cells on that level. Root = level 0
	Particles []*Particle // all particles, to be manipulated via Root.AddParticle

	FMMOrder  = 0
	Proximity = 1.1

	// statistics:
	totalPartners int
	totalNear     int
	totalCells    int
)

// Add a particle to the global FMM tree and Particles list
func AddParticle(p *Particle) {
	Root.addParticle(p)
	Particles = append(Particles, p)
}

// Calculates the magnetostatic field of all Particles.
func CalcDemag() {
	Root.updateM()

	Root.b0 = Vector{0, 0, 0}
	Root.dbdx = Vector{0, 0, 0}
	Root.dbdy = Vector{0, 0, 0}
	Root.dbdz = Vector{0, 0, 0}

	switch FMMOrder {
	default:
		panic(fmt.Sprint("invalid FMMOrder:", FMMOrder))
	case 0:
		Root.updateBdemag0(&Root) // we abuse root as parent, it only propagetes zero fields
	case 1:
		Root.updateBdemag1(&Root)
	case -1:
		CalcDemagBrute()
	}
}

func CalcDemagBrute() {
	for _, p := range Particles {
		p.b = p.BruteDemag()
	}
}

// Initializes the global FMM variables Root, Level
// with an FMM octree, nLevels deep (8^(nLevels-1)) base cells.
func InitFMM(worldSize Vector, nLevels int) {
	Level = make([][]*Cell, nLevels)
	Root = Cell{size: worldSize}

	start := time.Now()
	Log("dividing", nLevels, "levels", math.Pow(8, float64(nLevels-1)), "base cells...")
	Root.Divide(nLevels)
	Log(time.Since(start))

	start = time.Now()
	Log("finding partners...")
	Root.FindPartners(Level[0])
	Log(time.Since(start))

	printFMMStats()
}

func printFMMStats() {
	nLeaf := int(math.Pow(8, float64(len(Level)-1)) + 0.5)
	Log(totalCells, "cells, avg", totalPartners/totalCells, "partners/cell, avg", totalNear/nLeaf, "near/leaf")
}
