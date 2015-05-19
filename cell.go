package main

import "fmt"

// Cell in the FMM tree
type Cell struct {
	child    [8]*Cell   // octree of child cells
	partner  []*Cell    // I receive field taylor expansion from these cells
	near     []*Cell    // I recieve brute-force field contributions form these cells
	center   Vector     // my position
	size     Vector     // my diameter (x, y, z)
	m        Vector     // sum of child+particle magnetizations
	particle []Particle // If I'm a leaf cell, particles inside me (nil otherwise)
}

// unit vectors for left-back-bottom, left-back-top, ...
var direction = [8]Vector{
	{-1, -1, -1}, {-1, -1, +1}, {-1, +1, -1}, {-1, +1, +1},
	{+1, -1, -1}, {+1, -1, +1}, {+1, +1, -1}, {+1, +1, +1}}

// Create child cells to reach nLevels of levels and add to global level array.
// nLevels == 1 stops creating children (we always already have at least 1 level),
// but still adds the cell to the global level array.
func (c *Cell) Divide(nLevels int) {

	myLevel := len(level) - nLevels
	println(nLevels, myLevel)
	level[myLevel] = append(level[myLevel], c)

	if nLevels == 1 {
		return
	}

	for i := range c.child {
		newSize := c.size.Div(2)
		newCenter := c.center.Add(direction[i].Mul3(newSize.Div(2)))
		c.child[i] = &Cell{center: newCenter, size: newSize}
	}

	// recursively divide further
	for _, c := range c.child {
		c.Divide(nLevels - 1)
	}
}

func (c *Cell) String() string {
	if c == nil {
		return "nil cell"
	} else {
		return fmt.Sprint("cell@", c.center, " size=", c.size)
	}
}
