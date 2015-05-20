package vinamax2

//Calculates the torque working on the magnetisation of a particle
//using the Landau Lifshitz equation
func (p *Particle) tau(temp Vector) Vector {
	pm := &p.M
	mxB := pm.Cross(p.b_eff(temp))
	amxmxB := pm.Cross(mxB).Mul(Alpha)
	mxB = mxB.Add(amxmxB)
	return mxB.Mul(-gammaoveralpha)
}
