package vinamax2

import (
	"math"
	"math/rand"
)

var rng = rand.New(rand.NewSource(0))

//Sums the individual fields to the effective field working on the particle
func (p *Particle) b_eff(temp Vector) Vector {
	return p.anis().Add(p.zeeman().Add(temp))
}

//Set the randomseed for the temperature
func Setrandomseed(a int64) {
	randomseedcalled = true
	rng = rand.New(rand.NewSource(a))
}

// factor 4/3pi in "number" because they are spherical
func (p *Particle) calculatetempnumber() {
	p.tempnumber = math.Sqrt((2. * kb * Alpha * Temp) / (gamma0 * p.msat * 4. / 3. * math.Pi * p.r * p.r * p.r * Dt))
}

func calculatetempnumbers(lijst []*Particle) {
	for i := range lijst {
		lijst[i].calculatetempnumber()
	}
}

//Calculates the the thermal field B_therm working on a Particle
func (p *Particle) temp() Vector {
	B_therm := Vector{0., 0., 0.}
	if Brown {
		if Temp != 0. {
			etax := rng.NormFloat64()
			etay := rng.NormFloat64()
			etaz := rng.NormFloat64()

			B_therm = Vector{etax, etay, etaz}
			B_therm = B_therm.Mul(p.tempnumber)
		}
	}
	return B_therm
}

//Calculates the Zeeman field working on a Particle
func (p *Particle) zeeman() Vector {
	x, y, z := B_ext(T)
	x2, y2, z2 := B_ext_space(T, p.center[0], p.center[1], p.center[2])
	return Vector{x + x2, y + y2, z + z2}
}

//Calculates the anisotropy field working on a Particle
func (p *Particle) anis() Vector {
	//2*Ku1*(m.u)*u/p.msat

	mdotu := p.M.Dot(p.u_anis)
	uniax := p.u_anis.Mul(2. * Ku1 * mdotu / p.msat)

	cubic := Vector{0., 0., 0.}
	if Kc1 != 0 {
		c1m := p.M.Dot(p.c1_anis)
		c2m := p.M.Dot(p.c2_anis)
		c3m := p.M.Dot(p.c3_anis)
		firstterm := p.c1_anis.Mul(c1m * (c3m*c3m + c2m*c2m))
		secondterm := p.c2_anis.Mul(c2m * (c3m*c3m + c1m*c1m))
		thirdterm := p.c3_anis.Mul(c3m * (c2m*c2m + c1m*c1m))

		cubic = firstterm.Add(secondterm.Add(thirdterm)).Mul(-2. * Kc1 / p.msat)
	}
	return uniax.Add(cubic)
}
