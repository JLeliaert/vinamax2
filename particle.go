package vinamax2

import (
	"fmt"
	"math"
)

type Particle struct {
	M          Vector  // magnetization
	b          Vector  // demag field
	center     Vector  // particle position
	u_anis     Vector  // Uniaxial anisotropy axis
	c1_anis    Vector  // cubic anisotropy axis
	c2_anis    Vector  // cubic anisotropy axis
	c3_anis    Vector  // cubic anisotropy axis
	r          float64 // radius
	msat       float64 // Saturation magnetisation in A/m
	flip       float64 // time of next flip event
	heff       Vector  //effective field
	tempfield  Vector  //thermal field
	tempm      Vector  //temporary magnetisation
	tempnumber float64 // prefactor for thermal field
	previousm  Vector  //previous magnetisation
	k1         Vector
	k2         Vector
	k3         Vector
	k4         Vector
	k5         Vector
	k6         Vector
	k7         Vector
	k8         Vector
	k9         Vector
	k10        Vector
	k11        Vector
	k12        Vector
	k13        Vector
}

func NewParticle(center Vector, radius float64, M Vector, Msat float64) *Particle {
	return &Particle{center: center, M: M, r: radius, msat: Msat}
}

//print position and magnetisation of a particle
func (p Particle) string() string {
	return fmt.Sprintf("Particle@(%v, %v, %v), %v %v %v", p.center[0], p.center[1], p.center[2], p.M[0], p.M[1], p.M[2])
}

//Gives all particles the same specified uniaxialanisotropy-axis
func Anisotropy_axis(x, y, z float64) {
	uaniscalled = true
	a := Vector{x, y, z}
	norm(a)
	for i := range Particles {
		Particles[i].u_anis = a
	}
}

//Gives all particles the same specified cubic1anisotropy-axis
func C1anisotropy_axis(x, y, z float64) {
	c1called = true
	a := Vector{x, y, z}
	norm(a)
	for i := range Particles {
		Particles[i].c1_anis = a
	}
}

//Gives all particles the same specified cubic2anisotropy-axis, must be orthogonal to c1
func C2anisotropy_axis(x, y, z float64) {
	c2called = true
	a := Vector{x, y, z}
	norm(a)
	for i := range Particles {
		if Particles[i].c1_anis.Dot(a) == 0 {
			Particles[i].c2_anis = a
			b := Particles[i].c1_anis.Cross(a)
			norm(b)
			Particles[i].c3_anis = b
		} else {
			Fatal("c1 and c2 should be orthogonal")
		}
	}
}

//Gives all particles a random anisotropy-axis
func Anisotropy_random() {
	uaniscalled = true
	for i := range Particles {
		phi := rng.Float64() * (2 * math.Pi)
		theta := 2 * math.Asin(math.Sqrt(rng.Float64()))
		Particles[i].u_anis = Vector{math.Sin(theta) * math.Cos(phi), math.Sin(theta) * math.Sin(phi), math.Cos(theta)}
		if math.Cos(theta) < 0. {
			Particles[i].u_anis = Particles[i].u_anis.Mul(-1.)
		}
	}
}

//Gives all particles with random magnetisation orientation
func M_random() {
	magnetisationcalled = true
	for i := range Particles {
		phi := rng.Float64() * (2 * math.Pi)
		theta := 2 * math.Asin(math.Sqrt(rng.Float64()))
		Particles[i].M = Vector{math.Sin(theta) * math.Cos(phi), math.Sin(theta) * math.Sin(phi), math.Cos(theta)}
	}
}

//Gives all particles a specified magnetisation direction
func M_uniform(x, y, z float64) {
	magnetisationcalled = true
	a := Vector{x, y, z}
	norm(a)
	for i := range Particles {
		Particles[i].M = a
	}
}

//Sets the saturation magnetisation of all Particles in A/m
func Msat(x float64) {
	msatcalled = true
	for i := range Particles {
		Particles[i].msat = x
	}
}

// TODO (j): just export b field + rename
func (p *Particle) Bdemag() Vector { return p.b }
func (p *Particle) Center() Vector { return p.center }
