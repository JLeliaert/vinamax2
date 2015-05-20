package vinamax2

type Particle struct {
	M       Vector  // magnetization
	b       Vector  // demag field
	center  Vector  // particle position
	u_anis  Vector  // Uniaxial anisotropy axis
	c1_anis Vector  // cubic anisotropy axis
	c2_anis Vector  // cubic anisotropy axis
	c3_anis Vector  // cubic anisotropy axis
	r       float64 // radius
	msat    float64 // Saturation magnetisation in A/m
	flip    float64 // time of next flip event

	heff      Vector //effective field
	tempfield Vector //thermal field
	tempm     Vector //temporary magnetisation
	tempnumber float64 // prefactor for thermal field
	previousm Vector //previous magnetisation
	k1        Vector
	k2        Vector
	k3        Vector
	k4        Vector
	k5        Vector
	k6        Vector
	k7        Vector
	k8        Vector
	k9        Vector
	k10       Vector
	k11       Vector
	k12       Vector
	k13       Vector
}
