package vinamax2

import (
	"fmt"
	"math"
	"runtime"
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

var (
	ch   chan *Cell
	done chan struct{}
)

// EXPERIMENTAL
func CalcDemagParallel() {

	NCPU := 3

	if ch == nil {
		ch = make(chan *Cell, 8*8)
		done = make(chan struct{}, 8*8)
		runtime.GOMAXPROCS(NCPU)
		for i := 0; i < NCPU; i++ {
			go func() {
				Log("spinning up worker")
				for c := range ch {
					c.updateBdemag0(&Root)
					done <- struct{}{}
				}
			}()
		}
	}

	Root.updateM()

	Root.b0 = Vector{0, 0, 0}

	for _, c := range Root.child {
		for _, c := range c.child {
			ch <- c
		}
	}
	for _, c := range Root.child {
		for _, _ = range c.child {
			<-done
		}
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
	//TODO these have to be called when  they work
	//PruneTree()
	//CalculateCenterOfMags()

	printFMMStats()
}

func printFMMStats() {
	nLeaf := int(math.Pow(8, float64(len(Level)-1)) + 0.5)
	Log(totalCells, "cells, avg", totalPartners/totalCells, "partners/cell, avg", totalNear/nLeaf, "near/leaf")
}

//Prunes all the empty cells from the fmm tree
func PruneTree() {
	prune(&Root)
}

//ALL THESE ARE TODO /////////////////////////////////////////////////////////////////////////////////////
//recursively checks if a child cell contains particles and if not prunes them from the FMMtree
//TODO is not really recursive yet!!
func prune(c *Cell) {
	for _, c := range c.child {
		if c.IsLeaf() == false {
			if len(c.particles) != 0 {
				prune(c)
			} else {
				c = nil
			}
		}
	}
}

//TODO make recursive
func CalculateCenterOfMags() {
	updatecom(&Root)
}

//calculates com of a cell and than calls its child cells to do the same
func updatecom(c *Cell) {
	c.centerofmag = Vector{0, 0, 0}
	totalmoment := 0.
	for _, p := range c.particles {
		totalmoment += p.volume() * p.msat
		c.centerofmag.MAdd(p.volume()*p.msat, p.center)
		c.centerofmag.Div(totalmoment)
	}
	if c.IsLeaf() == false {
		for _, d := range c.child {
			if c != nil {
				updatecom(d)
			}
		}
	}
}
