package vinamax2

// Recursively calculate magnetostatic field via the FMM method, 1st-order.
// Precondition: Root.updateM() has been called.
func (c *Cell) updateBdemag1(parent *Cell) {
	if c == nil {
		return
	}

	// propagate parent field expansion to this cell,
	// (applies shift to Taylor expansion)
	sh := parent.center.Sub(c.center)
	c.b0 = parent.b0.MAdd(sh[X], parent.dbdx).MAdd(sh[Y], parent.dbdy).MAdd(sh[Z], parent.dbdz)
	c.dbdx = parent.dbdx
	c.dbdy = parent.dbdy
	c.dbdz = parent.dbdz

	c.addPartnerFields1()

	if !c.IsLeaf() {
		for _, ch := range c.child {
			ch.updateBdemag1(c)
		}
	} else {
		c.addNearFields1()
	}
}

// Add demag of nearby particles by brute force,
// start with 1st order evaluation of field in cell.
func (c *Cell) addNearFields1() {
	for _, dst := range c.particles {
		sh := dst.center.Sub(c.center)
		dst.b = c.b0.MAdd(sh[X], c.dbdx).MAdd(sh[Y], c.dbdy).MAdd(sh[Z], c.dbdz)
		c.addNearFields(dst)
	}
}

// add 1st-order expansions of fields of partner sources
func (c *Cell) addPartnerFields1() {
	for _, p := range c.partner {
		r := c.center.Sub(p.center)
		c.b0 = c.b0.Add(DipoleField(p.m, r))
		c.dbdx = c.dbdx.Add(DiffDipole(X, p.m, r))
		c.dbdy = c.dbdy.Add(DiffDipole(Y, p.m, r))
		c.dbdz = c.dbdz.Add(DiffDipole(Z, p.m, r))
	}
}