package vinamax2

import "fmt"

// Cell in the FMM tree
type Cell struct {
	parent           *Cell       //parent cell, or nil
	child            [8]*Cell    // octree of child cells
	partner          []*Cell     // I receive field taylor expansions from these cells
	near             []*Cell     // I recieve brute-force field contributions form these cells
	center           Vector      // my position
	size             Vector      // my diameter (x, y, z)
	m                Vector      // sum of child+particle magnetizations
	b0               Vector      // Field in cell center
	dbdx, dbdy, dbdz Vector      // Derivatives for Taylor expansion of field around center
	particles        []*Particle // If I'm a leaf cell, particles inside me (nil otherwise)
}

// find the leaf cell enclosing Particle p and add
// p to that cell. Called by AddParticle.
func (c *Cell) addParticle(p *Particle) {
	if c.IsLeaf() {
		if !c.contains(p.center) {
			panic("add particle: is outside cell")
		}
		c.particles = append(c.particles, p)
	} else {

		for _, c := range c.child {
			if c.contains(p.center) {
				c.addParticle(p)
				return
			}
		}
		panic("add particle: is outside cell")
	}
}

// is position x inside cell?
// used by AddParticle
func (c *Cell) contains(x Vector) bool {
	size := c.size.Div(2)
	d := x.Sub(c.center).Abs()
	return d[X] <= size[X] && d[Y] <= size[Y] && d[Z] <= size[Z]
}

// Like updateBdemag1, but 0th-order.
func (c *Cell) updateBdemag0(parent *Cell) {
	if c == nil {
		return
	}

	// use parent field in this cell
	c.b0 = parent.b0
	c.addPartnerFields0()

	if !c.IsLeaf() {
		for _, ch := range c.child {
			ch.updateBdemag0(c)
		}
	} else {
		c.addNearFields0()
	}
}

func CalcDemagIter() {
	Root.updateM()
	Root.b0 = Vector{0, 0, 0}
	for i := 1; i < len(Level)-1; i++ {
		for _, c := range Level[i] {
			c.b0 = c.parent.b0
			c.addPartnerFields0()
		}
	}
	for _, c := range Level[len(Level)-1] {
		c.b0 = c.parent.b0
		c.addPartnerFields0()
		c.addNearFields0()
	}
}

// like addPartnerFields1, but 0th order.
func (c *Cell) addPartnerFields0() {
	for _, p := range c.partner {
		r := Vector{c.center[X] - p.center[X], c.center[Y] - p.center[Y], c.center[Z] - p.center[Z]}
		b := DipoleField(p.m, r)
		c.b0[X] += b[X]
		c.b0[Y] += b[Y]
		c.b0[Z] += b[Z]
	}
}

// Like addNearFields1, but 0th-order.
func (c *Cell) addNearFields0() {
	for _, dst := range c.particles {
		dst.b = c.b0
		c.addNearFields(dst)
	}
}

// add, in a brute-force way, the near particle's fields, to p
func (c *Cell) addNearFields(dst *Particle) {
	for _, n := range c.near {
		for _, src := range n.particles {
			r := dst.center.Sub(src.center)
			dst.b = dst.b.Add(DipoleField(src.M, r))
		}
	}
}

// recursively update this cell's m as the sum
// of its children's m.
func (c *Cell) updateM() {
	c.m = Vector{0, 0, 0}

	// leaf node: sum particle m's.
	if c.particles != nil {
		for _, p := range c.particles {
			c.m = c.m.Add(p.M)
		}
		return
	}

	// non-leaf: update children, then add to me.
	for _, ch := range c.child {
		if ch != nil {
			ch.updateM()
			c.m = c.m.Add(ch.m)
		}
	}
}

// Recursively find partner cells from a list candidates.
// To be called on the root cell with level[0] as candidates.
// Partners are selected from a cell's own level, so they have the same size
// (otherwise it's just anarchy!)
// TODO: when we have unit tests, it can be optimized not to do so many allocations
func (c *Cell) FindPartners(candidates []*Cell) {

	// select partners based on proximity,
	// the rest goes into "near":
	var near []*Cell
	for _, cand := range candidates {
		if IsFar(c, cand) {
			c.partner = append(c.partner, cand)
			totalPartners++
		} else {
			near = append(near, cand)
		}
	}

	// leaf cell uses near cells for brute-force,
	// others recursively pass near cells children as new candidates
	if c.IsLeaf() {
		c.near = near
		totalNear += len(near)
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
	// TODO:  improve!
	dist := a.center.Sub(b.center).Len()
	return dist > Proximity*a.size.Len()
}

// Create child cells to reach nLevels of levels and add to global level array.
// nLevels == 1 stops creating children (we always already have at least 1 level),
// but still adds the cell to the global level array.
func (c *Cell) Divide(nLevels int) {

	// add to global level array
	myLevel := len(Level) - nLevels
	Level[myLevel] = append(Level[myLevel], c)
	totalCells++

	if nLevels == 1 {
		return
	}

	// create children
	for i := range c.child {
		newSize := c.size.Div(2)
		newCenter := c.center.Add(direction[i].Mul3(newSize.Div(2)))
		c.child[i] = &Cell{parent: c, center: newCenter, size: newSize}
	}

	// recursively go further
	for _, c := range c.child {
		c.Divide(nLevels - 1)
	}
}

func (c *Cell) Center() Vector { return c.center }

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
