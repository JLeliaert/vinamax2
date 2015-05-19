package main

import "fmt"

// Cell in the FMM tree
type Cell struct {
	child    [8]*Cell   // octree of child cells
	partner  []*Cell    // I receive field taylor expansions from these cells
	near     []*Cell    // I recieve brute-force field contributions form these cells
	center   Vector     // my position
	size     Vector     // my diameter (x, y, z)
	m        Vector     // sum of child+particle magnetizations
	particle []Particle // If I'm a leaf cell, particles inside me (nil otherwise)
}

// Recursively find partner cells from a list candidates.
// To be called on the root cell with level[0] as candidates.
// Partners are selected from a cell's own level, so they have the same size
// (otherwise it's just anarchy!)
func (c *Cell) FindPartners(candidates []*Cell) {

	// select partners based on proximity,
	// the rest goes into "near":
	var near []*Cell
	for _, cand := range candidates {
		if IsFar(c, cand) {
			c.partner = append(c.partner, cand)
		} else {
			near = append(near, cand)
		}
	}

	// leaf cell uses near cells for brute-force,
	// others recursively pass near cells children as new candidates
	if c.IsLeaf() {
		c.near = near
	} else {
		// children of my near cells become parter candidates
		newCand := make([]*Cell, 0, 8*len(near))
		for _, n := range near {
			newCand = append(newCand, n.child[:]...)
		}

		// recursively find partners in new candidates
		for _, c := range c.child {
			c.FindPartners(newCand)
		}
	}
}

// Is this cell a leaf cell?
func (c *Cell) IsLeaf() bool {
	for _, c := range c.child {
		if c != nil {
			return false
		}
	}
	return true
}

// Are the cells considered far separated?
func IsFar(a, b *Cell) bool {
	// TODO: this is more or less a touch criterion: improve!
	dist := a.center.Sub(b.center).Len()
	return dist > 1.1*a.size.Len()
}

// Create child cells to reach nLevels of levels and add to global level array.
// nLevels == 1 stops creating children (we always already have at least 1 level),
// but still adds the cell to the global level array.
func (c *Cell) Divide(nLevels int) {

	// add to global level array
	myLevel := len(level) - nLevels
	level[myLevel] = append(level[myLevel], c)

	if nLevels == 1 {
		return
	}

	// create children
	for i := range c.child {
		newSize := c.size.Div(2)
		newCenter := c.center.Add(direction[i].Mul3(newSize.Div(2)))
		c.child[i] = &Cell{center: newCenter, size: newSize}
	}

	// recursively go further
	for _, c := range c.child {
		c.Divide(nLevels - 1)
	}
}

func (c *Cell) String() string {
	if c == nil {
		return "nil"
	} else {
		typ := "node"
		if c.IsLeaf() {
			typ = "leaf"
		}
		return fmt.Sprint(typ, "@", c.center, len(c.partner), "partners, ", len(c.near), "near")
	}
}

// unit vectors for left-back-bottom, left-back-top, ...
var direction = [8]Vector{
	{-1, -1, -1}, {-1, -1, +1}, {-1, +1, -1}, {-1, +1, +1},
	{+1, -1, -1}, {+1, -1, +1}, {+1, +1, -1}, {+1, +1, +1}}
