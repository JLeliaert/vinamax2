package vinamax2

import (
	"log"
	"math"
	"time"
)

// FMM globals
var (
	Root      Cell        // roots the entire FMM tree
	Level     [][]*Cell   // for each level of the FMM tree: all cells on that level. Root = level 0
	Particles []*Particle // all particles, to be manipulated via Root.AddParticle

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
	Root.updateBdemag(nil)
}

// Initializes the global FMM variables Root, Level
// with an FMM octree, nLevels deep (8^(nLevels-1)) base cells.
func InitFMM(worldSize Vector, nLevels int) {
	Level = make([][]*Cell, nLevels)
	Root = Cell{size: worldSize}

	start := time.Now()
	log.Println("dividing", nLevels, "levels", math.Pow(8, float64(nLevels-1)), "base cells...")
	Root.Divide(nLevels)
	log.Println(time.Since(start))

	start = time.Now()
	log.Println("finding partners...")
	Root.FindPartners(Level[0])
	log.Println(time.Since(start))

	printFMMStats()
}

func printFMMStats() {
	nLeaf := int(math.Pow(8, float64(len(Level)-1)) + 0.5)
	log.Println(totalCells, "cells, avg", totalPartners/totalCells, "partners/cell, avg", totalNear/nLeaf, "near/leaf")
}
