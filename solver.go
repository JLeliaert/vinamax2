package vinamax2

import (
	//	"fmt"
	"math"
)

func calculatedemag() {
	CalcDemag()
}

//Set the solver to use, "euler" or "heun"
func Setsolver(a string) {
	switch a {

	case "euler":
		{
			solver = "euler"
			order = 1
		}
	case "heun":
		{
			solver = "heun"
			order = 2
		}
	case "rk3":
		{
			solver = "rk3"
			order = 3
		}
	case "annelies":
		{
			solver = "annelies"
			order = 3
		}

	case "rk4":
		{
			solver = "rk4"
			order = 4
		}
	case "dopri":
		{
			solver = "dopri"
			order = 5
		}
	case "fehl56":
		{
			solver = "fehl56"
			order = 6
		}
	case "fehl67":
		{
			solver = "fehl67"
			order = 7
		}
	case "time":
		{
			solver = "time"
			order = 0
		}
	default:
		{
			Fatal(a, " is not a possible solver, \"euler\" or \"heun\" or \"rk3\"or \"rk4\"or \"dopri\"or \"fehl56\"or \"fehl67\"")
		}
	}
}

//Runs the simulation for a certain time

func Run(time float64) {
	gammaoveralpha = gamma0 / (1. + (Alpha * Alpha))
	testinput()
	syntaxrun()
	for i := range Particles {
		norm(Particles[i].M)
	}
	write(averagemoments(Particles))
	//averages is not weighted with volume, averagemoments is
	//write(averages(Particles))
	previousdemagcalc := T - demagtime
	for j := T; T < j+time; {
		if (demagevery == true) && (T-previousdemagcalc >= demagtime) {
			calculatedemag()
			previousdemagcalc = T
		}
		if Demag {
			calculatedemag()
		}
		switch solver {
		case "heun":
			{
				heunstep(Particles)
				T += Dt
			}
		case "euler":
			{
				eulerstep(Particles)
				T += Dt
			}
		case "rk3":
			{
				rk3step(Particles)
				T += Dt
			}
		case "annelies":
			{
				anneliesstep(Particles)
				T += Dt
			}
		case "rk4":
			{
				rk4step(Particles)
				T += Dt
			}
		case "dopri":
			{
				dopristep(Particles)
				T += Dt
				//fmt.Println(Dt)
				if Adaptivestep {
					if maxtauwitht > Errortolerance {
						undobadstep(Particles)
						if Dt == Mindt {
							Fatal("mindt is too small for your specified error tolerance")
						}
					}

					Dt = 0.95 * Dt * math.Pow(Errortolerance/maxtauwitht, (1./float64(order)))
					if Dt < Mindt {
						Dt = Mindt
					}
					if Dt > Maxdt {
						Dt = Maxdt
					}
					//fmt.Println("dt:   ", Dt)
					maxtauwitht = 1.e-12
				}
			}
		case "fehl56":
			{
				fehl56step(Particles)
				T += Dt
				if Adaptivestep {
					if maxtauwitht > Errortolerance {
						undobadstep(Particles)
						if Dt == Mindt {
							Fatal("mindt is too small for your specified error tolerance")
						}

					}

					Dt = 0.95 * Dt * math.Pow(Errortolerance/maxtauwitht, (1./float64(order)))
					if Dt < Mindt {
						Dt = Mindt
					}
					if Dt > Maxdt {
						Dt = Maxdt
					}
					//fmt.Println("dt:   ", Dt)
					maxtauwitht = 1.e-12
				}
			}
		case "fehl67":
			{
				fehl67step(Particles)
				T += Dt
				if Adaptivestep {
					if maxtauwitht > Errortolerance {
						undobadstep(Particles)
						if Dt == Mindt {
							Fatal("mindt is too small for your specified error tolerance")
						}

					}

					Dt = 0.95 * Dt * math.Pow(Errortolerance/maxtauwitht, (1./float64(order)))
					if Dt < Mindt {
						Dt = Mindt
					}
					if Dt > Maxdt {
						Dt = Maxdt
					}
					//	fmt.Println("dt:   ", Dt)
					maxtauwitht = 1.e-12
				}
			}

		case "time":
			{
				T += Dt
			}
		}

		//	plotswitchtime()//EXTRA
		//		if Jumpnoise {
		//			checkallswitch(Particles)
		//		}
		//fmt.Println(Dt)
		//write(averages(Particles))
		write(averagemoments(Particles))
		if (T > j+time-Dt) && (T < j+time) {
			//undobadstep(Particles)
			Dt = j + time - T + 1e-15
		}
	}

	//if suggest_timestep {
	//	printsuggestedtimestep()
	//}
}

//##################################################

//Perform a timestep using euler forward method
func eulerstep(Lijst []*Particle) {
	for _, p := range Lijst {
		temp := p.temp()

		tau := p.tau(temp)
		p.M = p.M.MAdd(Dt, tau)
		//p.M[0] += tau[0] * Dt
		//p.M[1] += tau[1] * Dt
		//p.M[2] += tau[2] * Dt
		norm(p.M)
		//if you have to save mdotH
		p.heff = p.b_eff(temp)
	}
}

//#########################################################################
//perform a timestep using heun method
//http://en.wikipedia.org/wiki/Heun_method
func heunstep(Lijst []*Particle) {
	for _, p := range Lijst {
		temp := p.temp()
		p.tempfield = temp
		tau1 := p.tau(temp)
		p.k1 = tau1

		//tau (t+1)
		p.M[0] += tau1[0] * Dt
		p.M[1] += tau1[1] * Dt
		p.M[2] += tau1[2] * Dt
		T += Dt
	}

	if Demag {
		calculatedemag()
	}

	for _, p := range Lijst {
		temp := p.tempfield
		//temp := p.temp()
		tau2 := p.tau(temp)
		tau1 := p.k1
		p.M = p.M.MAdd(-Dt*0.5, tau1).MAdd(Dt*0.5, tau2)

		//p.M[0] += ((-tau1[0] + tau2[0]) * 0.5 * Dt)
		//p.M[1] += ((-tau1[1] + tau2[1]) * 0.5 * Dt)
		//p.M[2] += ((-tau1[2] + tau2[2]) * 0.5 * Dt)

		norm(p.M)

		T -= Dt
		//when saving  mdotH
		p.heff = p.b_eff(temp)

	}
}

//#########################################################################

//perform a timestep using 3th order RK
func rk3step(Lijst []*Particle) {
	for _, p := range Lijst {
		temp := p.temp()
		p.tempfield = temp
		p.k1 = p.tau(temp)

		//k1
		p.M = p.M.MAdd(0.5*Dt, p.k1)
		//p.M[0] += tau0[0] * 1 / 2. * Dt
		//p.M[1] += tau0[1] * 1 / 2. * Dt
		//p.M[2] += tau0[2] * 1 / 2. * Dt
		T += 1 / 2. * Dt
	}

	if Demag {
		calculatedemag()
	}

	for _, p := range Lijst {
		temp := p.tempfield
		p.k2 = p.tau(temp)

		p.M = p.M.MAdd(-3./2.*Dt, p.k1).MAdd(2*Dt, p.k2)
		//p.M[0] += ((-3/2.*k1[0] + 2*k2[0]) * Dt)
		//p.M[1] += ((-3/2.*k1[1] + 2*k2[1]) * Dt)
		//p.M[2] += ((-3/2.*k1[2] + 2*k2[2]) * Dt)
		T += 1 / 2. * Dt
	}
	if Demag {
		calculatedemag()
	}
	for _, p := range Lijst {
		temp := p.tempfield
		p.k3 = p.tau(temp)
		p.M = p.M.MAdd(7./6.*Dt, p.k1).MAdd(-4./3.*Dt, p.k2).MAdd(1./6.*Dt, p.k3)
		//p.M[0] += ((7/6.*k1[0] - 4/3.*k2[0] + 1/6.*k3[0]) * Dt)
		//p.M[1] += ((7/6.*k1[1] - 4/3.*k2[1] + 1/6.*k3[1]) * Dt)
		//p.M[2] += ((7/6.*k1[2] - 4/3.*k2[2] + 1/6.*k3[2]) * Dt)

		norm(p.M)

		T -= Dt
		//when saving mdotH
		p.heff = p.b_eff(temp)

	}
}

//#########################################################################

//perform a timestep using 3th order anneliessolver
func anneliesstep(Lijst []*Particle) {
	for _, p := range Lijst {
		temp := p.temp()
		p.tempfield = temp
		p.k1 = p.tau(temp)

		p.M = p.M.MAdd(1./10.*Dt, p.k1)
		//p.M[0] += tau0[0] * 1. / 10. * Dt
		//p.M[1] += tau0[1] * 1. / 10. * Dt
		//p.M[2] += tau0[2] * 1. / 10. * Dt
		T += 1. / 10. * Dt
	}

	if Demag {
		calculatedemag()
	}

	for _, p := range Lijst {
		temp := p.tempfield
		p.k2 = p.tau(temp)

		p.M = p.M.MAdd((-1./10.-2189./5746)*Dt, p.k1).MAdd(2310./2873.*Dt, p.k2)
		//p.M[0] += (((-1./10.-2189./5746.)*k1[0] + 2310./2873.*k2[0]) * Dt)
		//p.M[1] += (((-1./10.-2189./5746.)*k1[1] + 2310./2873.*k2[1]) * Dt)
		//p.M[2] += (((-1./10.-2189./5746.)*k1[2] + 2310./2873.*k2[2]) * Dt)
		T += (-1./10. + 11./26.) * Dt
	}
	if Demag {
		calculatedemag()
	}
	for _, p := range Lijst {
		temp := p.tempfield
		p.k3 = p.tau(temp)

		p.M = p.M.MAdd((89./33.+2189./5746.)*Dt, p.k1).MAdd((-475./126-2310./2873.)*Dt, p.k2).MAdd(2873/1386.*Dt, p.k3)
		//p.M[0] += (((89./33.+2189./5746.)*k1[0] + (-475./126-2310./2873.)*k2[0] + 2873/1386.*k3[0]) * Dt)
		//p.M[1] += (((89./33.+2189./5746.)*k1[1] + (-475./126-2310./2873.)*k2[1] + 2873/1386.*k3[1]) * Dt)
		//p.M[2] += (((89./33.+2189./5746.)*k1[2] + (-475./126-2310./2873.)*k2[2] + 2873/1386.*k3[2]) * Dt)
		T += (-11/26. + 1.) * Dt

		norm(p.M)

		T -= Dt
		//when saving mdotH
		p.heff = p.b_eff(temp)

	}
}

//#########################################################################

//perform a timestep using 4th order RK
func rk4step(Lijst []*Particle) {
	for _, p := range Lijst {
		temp := p.temp()
		p.tempfield = temp
		p.k1 = p.tau(temp)

		p.M = p.M.MAdd(0.5*Dt, p.k1)
		//p.M[0] += tau0[0] * 1 / 2. * Dt
		//p.M[1] += tau0[1] * 1 / 2. * Dt
		//p.M[2] += tau0[2] * 1 / 2. * Dt
		T += 1 / 2. * Dt
	}

	if Demag {
		calculatedemag()
	}

	for _, p := range Lijst {
		temp := p.tempfield
		p.k2 = p.tau(temp)

		p.M = p.M.MAdd(-0.5*Dt, p.k1).MAdd(0.5*Dt, p.k2)
		//p.M[0] += ((-1/2.*k1[0] + 1/2.*k2[0]) * Dt)
		//p.M[1] += ((-1/2.*k1[1] + 1/2.*k2[1]) * Dt)
		//p.M[2] += ((-1/2.*k1[2] + 1/2.*k2[2]) * Dt)
	}
	if Demag {
		calculatedemag()
	}
	for _, p := range Lijst {
		temp := p.tempfield
		p.k3 = p.tau(temp)

		p.M = p.M.MAdd(-0.5*Dt, p.k2).MAdd(Dt, p.k3)
		//p.M[0] += ((-1/2.*k2[0] + 1*k3[0]) * Dt)
		//p.M[1] += ((-1/2.*k2[1] + 1*k3[1]) * Dt)
		//p.M[2] += ((-1/2.*k2[2] + 1*k3[2]) * Dt)
		T += 1 / 2. * Dt
	}
	if Demag {
		calculatedemag()
	}
	for _, p := range Lijst {
		temp := p.tempfield
		p.k4 = p.tau(temp)

		p.M = p.M.MAdd(1/6.*Dt, p.k1).MAdd(1/3.*Dt, p.k2).MAdd(-2/3.*Dt, p.k3).MAdd(1/6.*Dt, p.k4)
		//p.M[0] += ((1/6.*k1[0] + 1/3.*k2[0] - 2/3.*k3[0] + 1/6.*k4[0]) * Dt)
		//p.M[1] += ((1/6.*k1[1] + 1/3.*k2[1] - 2/3.*k3[1] + 1/6.*k4[1]) * Dt)
		//p.M[2] += ((1/6.*k1[2] + 1/3.*k2[2] - 2/3.*k3[2] + 1/6.*k4[2]) * Dt)

		norm(p.M)
		T -= Dt
		//when saving mdotH
		p.heff = p.b_eff(temp)

	}
}

//#########################################################################
//perform a timestep using dormand-prince

// Gebruik maken van de FSAL (enkel bij niet-brown noise!!!)

func dopristep(Lijst []*Particle) {
	for _, p := range Lijst {
		p.tempm = p.M
		p.previousm = p.M

		temp := p.temp()
		p.tempfield = temp
		p.k1 = p.tau(temp)

		p.M = p.M.MAdd(1/5.*Dt, p.k1)
		//p.M[0] += k1[0] * 1 / 5. * Dt
		//p.M[1] += k1[1] * 1 / 5. * Dt
		//p.M[2] += k1[2] * 1 / 5. * Dt
		T += 1 / 5. * Dt

	}
	if Demag {
		calculatedemag()
	}
	for _, p := range Lijst {

		temp := p.tempfield
		p.k2 = p.tau(temp)

		p.M = p.tempm.MAdd(3/40.*Dt, p.k1).MAdd(9/40.*Dt, p.k2)
		//p.M = p.tempm
		//p.M[0] += ((3/40.*k1[0] + 9/40.*k2[0]) * Dt)
		//p.M[1] += ((3/40.*k1[1] + 9/40.*k2[1]) * Dt)
		//p.M[2] += ((3/40.*k1[2] + 9/40.*k2[2]) * Dt)
		T += 1 / 10. * Dt
	}
	if Demag {
		calculatedemag()
	}
	for _, p := range Lijst {
		temp := p.tempfield
		p.k3 = p.tau(temp)

		p.M = p.tempm.MAdd(44/45.*Dt, p.k1).MAdd(-56/15.*Dt, p.k2).MAdd(32/9.*Dt, p.k3)
		//p.M = p.tempm
		//p.M[0] += ((44/45.*k1[0] - 56/15.*k2[0] + 32/9.*k3[0]) * Dt)
		//p.M[1] += ((44/45.*k1[1] - 56/15.*k2[1] + 32/9.*k3[1]) * Dt)
		//p.M[2] += ((44/45.*k1[2] - 56/15.*k2[2] + 32/9.*k3[2]) * Dt)
		T += 1 / 2. * Dt
	}
	if Demag {
		calculatedemag()
	}
	for _, p := range Lijst {
		temp := p.tempfield
		p.k4 = p.tau(temp)

		p.M = p.tempm.MAdd(19372/6561.*Dt, p.k1).MAdd(-25360/2187.*Dt, p.k2).MAdd(64448/6561.*Dt, p.k3).MAdd(-212/729.*Dt, p.k4)
		//p.M = p.tempm
		//p.M[0] += ((19372/6561.*k1[0] - 25360/2187.*k2[0] + 64448/6561.*k3[0] - 212/729.*k4[0]) * Dt)
		//p.M[1] += ((19372/6561.*k1[1] - 25360/2187.*k2[1] + 64448/6561.*k3[1] - 212/729.*k4[1]) * Dt)
		//p.M[2] += ((19372/6561.*k1[2] - 25360/2187.*k2[2] + 64448/6561.*k3[2] - 212/729.*k4[2]) * Dt)
		T += (-4/5. + 8/9.) * Dt
	}
	if Demag {
		calculatedemag()
	}
	for _, p := range Lijst {
		temp := p.tempfield
		p.k5 = p.tau(temp)

		p.M = p.tempm.MAdd(9017/3168.*Dt, p.k1).MAdd(-355/33.*Dt, p.k2).MAdd(46732/5247.*Dt, p.k3).MAdd(49/176.*Dt, p.k4).MAdd(-5103/18656.*Dt, p.k5)
		//p.M = p.tempm
		//p.M[0] += ((9017/3168.*k1[0] - 355/33.*k2[0] + 46732/5247.*k3[0] + 49/176.*k4[0] - 5103/18656.*k5[0]) * Dt)
		//p.M[1] += ((9017/3168.*k1[1] - 355/33.*k2[1] + 46732/5247.*k3[1] + 49/176.*k4[1] - 5103/18656.*k5[1]) * Dt)
		//p.M[2] += ((9017/3168.*k1[2] - 355/33.*k2[2] + 46732/5247.*k3[2] + 49/176.*k4[2] - 5103/18656.*k5[2]) * Dt)
		T += 1 / 9. * Dt
	}
	if Demag {
		calculatedemag()
	}

	for _, p := range Lijst {
		temp := p.tempfield
		p.k6 = p.tau(temp)

		p.M = p.tempm.MAdd(35/384.*Dt, p.k1).MAdd(500/1113.*Dt, p.k3).MAdd(125/192.*Dt, p.k4).MAdd(-2187/6784.*Dt, p.k5).MAdd(11/84.*Dt, p.k6)
		//p.M = p.tempm
		//p.M[0] += ((35/384.*k1[0] + 0.*k2[0] + 500/1113.*k3[0] + 125/192.*k4[0] - 2187/6784.*k5[0] + 11/84.*k6[0]) * Dt)
		//p.M[1] += ((35/384.*k1[1] + 0.*k2[1] + 500/1113.*k3[1] + 125/192.*k4[1] - 2187/6784.*k5[1] + 11/84.*k6[1]) * Dt)
		//p.M[2] += ((35/384.*k1[2] + 0.*k2[2] + 500/1113.*k3[2] + 125/192.*k4[2] - 2187/6784.*k5[2] + 11/84.*k6[2]) * Dt)
		//and this is also the fifth order solution
	}
	if Demag {
		calculatedemag()
	}

	for _, p := range Lijst {
		temp := p.tempfield
		p.k7 = p.tau(temp)

		p.tempm = p.tempm.MAdd(5179/57600.*Dt, p.k1).MAdd(7571/16695.*Dt, p.k3).MAdd(393/640.*Dt, p.k4).MAdd(-92097/339200.*Dt, p.k5).MAdd(187/2100.*Dt, p.k6).MAdd(1/40.*Dt, p.k7)
		//p.tempm[0] += ((5179/57600.*k1[0] + 0.*k2[0] + 7571/16695.*k3[0] + 393/640.*k4[0] - 92097/339200.*k5[0] + 187/2100.*k6[0] + 1/40.*k7[0]) * Dt)
		//p.tempm[1] += ((5179/57600.*k1[1] + 0.*k2[1] + 7571/16695.*k3[1] + 393/640.*k4[1] - 92097/339200.*k5[1] + 187/2100.*k6[1] + 1/40.*k7[1]) * Dt)
		//p.tempm[2] += ((5179/57600.*k1[2] + 0.*k2[2] + 7571/16695.*k3[2] + 393/640.*k4[2] - 92097/339200.*k5[2] + 187/2100.*k6[2] + 1/40.*k7[2]) * Dt)

		//and this is also the fourth order solution
		norm(p.M)
		norm(p.tempm)

		//the error is the difference between the two solutions
		//error := math.Sqrt(sqr(p.m[0]-p.tempm[0]) + sqr(p.m[1]-p.tempm[1]) + sqr(p.m[2]-p.tempm[2]))
		error := p.M.Sub(p.tempm).Len()

		//fmt.Println("error    :", error)
		if Adaptivestep {
			if error > maxtauwitht {
				maxtauwitht = error
			}
		}
		//when saving mdotH
		p.heff = p.b_eff(temp)

		T -= Dt
	}
}

///#########################################################################
//perform a timestep using fehlberg 56 method

func fehl56step(Lijst []*Particle) {
	for _, p := range Lijst {
		p.tempm = p.M
		p.previousm = p.M

		temp := p.temp()
		p.tempfield = temp
		p.k1 = p.tau(temp)

		p.M = p.M.MAdd(1/6.*Dt, p.k1)
		//p.M[0] += k1[0] * 1 / 6. * Dt
		//p.M[1] += k1[1] * 1 / 6. * Dt
		//p.M[2] += k1[2] * 1 / 6. * Dt
		T += 1 / 6. * Dt

	}
	if Demag {
		calculatedemag()
	}
	for _, p := range Lijst {

		temp := p.tempfield
		p.k2 = p.tau(temp)

		p.M = p.tempm.MAdd(4/75.*Dt, p.k1).MAdd(16/75.*Dt, p.k2)
		//p.M = p.tempm
		//p.M[0] += ((4/75.*k1[0] + 16/75.*k2[0]) * Dt)
		//p.M[1] += ((4/75.*k1[1] + 16/75.*k2[1]) * Dt)
		//p.M[2] += ((4/75.*k1[2] + 16/75.*k2[2]) * Dt)
		T += (-1/6. + 4/15.) * Dt
	}
	if Demag {
		calculatedemag()
	}
	for _, p := range Lijst {
		temp := p.tempfield
		p.k3 = p.tau(temp)

		p.M = p.tempm.MAdd(5/6.*Dt, p.k1).MAdd(-8/3.*Dt, p.k2).MAdd(5/2.*Dt, p.k3)
		//p.M = p.tempm
		//p.M[0] += ((5/6.*k1[0] - 8/3.*k2[0] + 5/2.*k3[0]) * Dt)
		//p.M[1] += ((5/6.*k1[1] - 8/3.*k2[1] + 5/2.*k3[1]) * Dt)
		//p.M[2] += ((5/6.*k1[2] - 8/3.*k2[2] + 5/2.*k3[2]) * Dt)
		T += (-4/15. + 2/3.) * Dt
	}
	if Demag {
		calculatedemag()
	}
	for _, p := range Lijst {
		temp := p.tempfield
		p.k4 = p.tau(temp)

		p.M = p.tempm.MAdd(-8/5.*Dt, p.k1).MAdd(144/25.*Dt, p.k2).MAdd(-4.*Dt, p.k3).MAdd(16/25.*Dt, p.k4)
		//p.M = p.tempm
		//p.M[0] += ((-8/5.*k1[0] + 144/25.*k2[0] - 4.*k3[0] + 16/25.*k4[0]) * Dt)
		//p.M[1] += ((-8/5.*k1[1] + 144/25.*k2[1] - 4.*k3[1] + 16/25.*k4[1]) * Dt)
		//p.M[2] += ((-8/5.*k1[2] + 144/25.*k2[2] - 4.*k3[2] + 16/25.*k4[2]) * Dt)
		T += (-2/3. + 4/5.) * Dt
	}
	if Demag {
		calculatedemag()
	}
	for _, p := range Lijst {
		temp := p.tempfield
		p.k5 = p.tau(temp)

		p.M = p.tempm.MAdd(361/320.*Dt, p.k1).MAdd(-18/5.*Dt, p.k2).MAdd(407/128.*Dt, p.k3).MAdd(-11/80.*Dt, p.k4).MAdd(55/128.*Dt, p.k5)
		//p.M = p.tempm
		//p.M[0] += ((361/320.*k1[0] - 18/5.*k2[0] + 407/128.*k3[0] - 11/80.*k4[0] + 55/128.*k5[0]) * Dt)
		//p.M[1] += ((361/320.*k1[1] - 18/5.*k2[1] + 407/128.*k3[1] - 11/80.*k4[1] + 55/128.*k5[1]) * Dt)
		//p.M[2] += ((361/320.*k1[2] - 18/5.*k2[2] + 407/128.*k3[2] - 11/80.*k4[2] + 55/128.*k5[2]) * Dt)
		T += 1 / 5. * Dt
	}
	if Demag {
		calculatedemag()
	}

	for _, p := range Lijst {
		temp := p.tempfield
		p.k6 = p.tau(temp)

		p.M = p.tempm.MAdd(-11/640.*Dt, p.k1).MAdd(11/256.*Dt, p.k3).MAdd(-11/160.*Dt, p.k4).MAdd(11/256.*Dt, p.k5)
		//p.M = p.tempm
		//p.M[0] += ((-11/640.*k1[0] + 0.*k2[0] + 11/256.*k3[0] - 11/160.*k4[0] + 11/256.*k5[0] + 0.*k6[0]) * Dt)
		//p.M[1] += ((-11/640.*k1[1] + 0.*k2[1] + 11/256.*k3[1] - 11/160.*k4[1] + 11/256.*k5[1] + 0.*k6[1]) * Dt)
		//p.M[2] += ((-11/640.*k1[2] + 0.*k2[2] + 11/256.*k3[2] - 11/160.*k4[2] + 11/256.*k5[2] + 0.*k6[2]) * Dt)
		T -= 1 / 1. * Dt
	}
	if Demag {
		calculatedemag()
	}

	for _, p := range Lijst {
		temp := p.tempfield
		p.k7 = p.tau(temp)

		p.M = p.tempm.MAdd(93/640.*Dt, p.k1).MAdd(-18/5.*Dt, p.k2).MAdd(803/256.*Dt, p.k3).MAdd(-11/160.*Dt, p.k4).MAdd(99/256.*Dt, p.k5).MAdd(Dt, p.k7)
		//p.M = p.tempm
		//p.M[0] += ((93/640.*k1[0] + -18/5.*k2[0] + 803/256.*k3[0] - 11/160.*k4[0] + 99/256.*k5[0] + 0.*k6[0] + 1/1.*k7[0]) * Dt)
		//p.M[1] += ((93/640.*k1[1] + -18/5.*k2[1] + 803/256.*k3[1] - 11/160.*k4[1] + 99/256.*k5[1] + 0.*k6[1] + 1/1.*k7[1]) * Dt)
		//p.M[2] += ((93/640.*k1[2] + -18/5.*k2[2] + 803/256.*k3[2] - 11/160.*k4[2] + 99/256.*k5[2] + 0.*k6[2] + 1/1.*k7[2]) * Dt)
		T += 1 / 1. * Dt
	}

	for _, p := range Lijst {
		temp := p.tempfield
		p.k8 = p.tau(temp)

		p.M = p.tempm.MAdd(31/384.*Dt, p.k1).MAdd(1125/2816.*Dt, p.k3).MAdd(9/32.*Dt, p.k4).MAdd(125/768.*Dt, p.k5).MAdd(5/66.*Dt, p.k6)
		//p.M = p.tempm
		//p.tempm[0] += ((31/384.*k1[0] + 0.*k2[0] + 1125/2816.*k3[0] + 9/32.*k4[0] + 125/768.*k5[0] + 5/66.*k6[0] + 0/1.*k7[0]) * Dt)
		//p.tempm[1] += ((31/384.*k1[1] + 0.*k2[1] + 1125/2816.*k3[1] + 9/32.*k4[1] + 125/768.*k5[1] + 5/66.*k6[1] + 0/1.*k7[1]) * Dt)
		//p.tempm[2] += ((31/384.*k1[2] + 0.*k2[2] + 1125/2816.*k3[2] + 9/32.*k4[2] + 125/768.*k5[2] + 5/66.*k6[2] + 0/1.*k7[2]) * Dt)
		//fifth order solution

		p.M = p.tempm.MAdd(-5/66.*Dt, p.k1).MAdd(-5/66.*Dt, p.k6).MAdd(5/66.*Dt, p.k7).MAdd(5/66.*Dt, p.k8)
		//p.M[0] = p.tempm[0] + ((-5/66.*k1[0] + -5/66.*k6[0] + 5/66.*k7[0] + 5/66.*k8[0]) * Dt)
		//p.M[1] = p.tempm[1] + ((-5/66.*k1[1] + -5/66.*k6[1] + 5/66.*k7[1] + 5/66.*k8[1]) * Dt)
		//p.M[2] = p.tempm[2] + ((-5/66.*k1[2] + -5/66.*k6[2] + 5/66.*k7[2] + 5/66.*k8[2]) * Dt)
		//sixth order solution

		norm(p.M)
		norm(p.tempm)

		//the error is the difference between the two solutions
		//error := math.Sqrt(sqr(p.M[0]-p.tempm[0]) + sqr(p.M[1]-p.tempm[1]) + sqr(p.M[2]-p.tempm[2]))
		error := p.M.Sub(p.tempm).Len()

		//fmt.Println("error    :", error)
		if Adaptivestep {
			if error > maxtauwitht {
				maxtauwitht = error
			}
		}
		//when saving mdotH
		p.heff = p.b_eff(temp)

		T -= Dt
	}
}

//###########################################################################################################

//perform a timestep using fehlberg 67 method

func fehl67step(Lijst []*Particle) {
	for _, p := range Lijst {
		p.tempm = p.M
		p.previousm = p.M

		temp := p.temp()
		p.tempfield = temp
		p.k1 = p.tau(temp)

		p.M = p.M.MAdd(2/27.*Dt, p.k1)
		//p.M[0] += k1[0] * 2 / 27. * Dt
		//p.M[1] += k1[1] * 2 / 27. * Dt
		//p.M[2] += k1[2] * 2 / 27. * Dt
		T += 2 / 27. * Dt

	}
	if Demag {
		calculatedemag()
	}
	for _, p := range Lijst {

		temp := p.tempfield
		p.k2 = p.tau(temp)

		p.M = p.tempm.MAdd(1/36.*Dt, p.k1).MAdd(1/12.*Dt, p.k2)
		//p.M = p.tempm
		//p.M[0] += ((1/36.*k1[0] + 1/12.*k2[0]) * Dt)
		//p.M[1] += ((1/36.*k1[1] + 1/12.*k2[1]) * Dt)
		//p.M[2] += ((1/36.*k1[2] + 1/12.*k2[2]) * Dt)
		T += (-2/27. + 1/9.) * Dt
	}
	if Demag {
		calculatedemag()
	}
	for _, p := range Lijst {
		temp := p.tempfield
		p.k3 = p.tau(temp)

		p.M = p.tempm.MAdd(1/24.*Dt, p.k1).MAdd(1/8.*Dt, p.k3)
		//p.M = p.tempm
		//p.M[0] += ((1/24.*k1[0] + 0.*k2[0] + 1/8.*k3[0]) * Dt)
		//p.M[1] += ((1/24.*k1[1] + 0.*k2[1] + 1/8.*k3[1]) * Dt)
		//p.M[2] += ((1/24.*k1[2] + 0.*k2[2] + 1/8.*k3[2]) * Dt)
		T += (-1/9. + 1/6.) * Dt
	}
	if Demag {
		calculatedemag()
	}
	for _, p := range Lijst {
		temp := p.tempfield
		p.k4 = p.tau(temp)

		p.M = p.tempm.MAdd(5/12.*Dt, p.k1).MAdd(-25/16.*Dt, p.k3).MAdd(25/16.*Dt, p.k4)
		//p.M = p.tempm
		//p.M[0] += ((5/12.*k1[0] + 0.*k2[0] - 25/16.*k3[0] + 25/16.*k4[0]) * Dt)
		//p.M[1] += ((5/12.*k1[1] + 0.*k2[1] - 25/16.*k3[1] + 25/16.*k4[1]) * Dt)
		//p.M[2] += ((5/12.*k1[2] + 0.*k2[2] - 25/16.*k3[2] + 25/16.*k4[2]) * Dt)
		T += (-1/6. + 5/12.) * Dt
	}
	if Demag {
		calculatedemag()
	}
	for _, p := range Lijst {
		temp := p.tempfield
		p.k5 = p.tau(temp)

		p.M = p.tempm.MAdd(1/20.*Dt, p.k1).MAdd(1/4.*Dt, p.k4).MAdd(1/5.*Dt, p.k5)
		//p.M = p.tempm
		//p.M[0] += ((1/20.*k1[0] + 0.*k2[0] + 0.*k3[0] + 1/4.*k4[0] + 1/5.*k5[0]) * Dt)
		//p.M[1] += ((1/20.*k1[1] + 0.*k2[1] + 0.*k3[1] + 1/4.*k4[1] + 1/5.*k5[1]) * Dt)
		//p.M[2] += ((1/20.*k1[2] + 0.*k2[2] + 0.*k3[2] + 1/4.*k4[2] + 1/5.*k5[2]) * Dt)
		T += (-5/12. + 1/2.) * Dt
	}
	if Demag {
		calculatedemag()
	}

	for _, p := range Lijst {
		temp := p.tempfield
		p.k6 = p.tau(temp)

		p.M = p.tempm.MAdd(-25/108.*Dt, p.k1).MAdd(125/108.*Dt, p.k4).MAdd(-65/27.*Dt, p.k5).MAdd(125/54.*Dt, p.k6)
		//p.M = p.tempm
		//p.M[0] += ((-25/108.*k1[0] + 0.*k2[0] + 0.*k3[0] + 125/108.*k4[0] - 65/27.*k5[0] + 125/54.*k6[0]) * Dt)
		//p.M[1] += ((-25/108.*k1[1] + 0.*k2[1] + 0.*k3[1] + 125/108.*k4[1] - 65/27.*k5[1] + 125/54.*k6[1]) * Dt)
		//p.M[2] += ((-25/108.*k1[2] + 0.*k2[2] + 0.*k3[2] + 125/108.*k4[2] - 65/27.*k5[2] + 125/54.*k6[2]) * Dt)
		T += (-1/2. + 5/6.) * Dt
	}
	if Demag {
		calculatedemag()
	}

	for _, p := range Lijst {
		temp := p.tempfield
		p.k7 = p.tau(temp)

		p.M = p.tempm.MAdd(31/300.*Dt, p.k1).MAdd(61/225.*Dt, p.k5).MAdd(-2/9.*Dt, p.k6).MAdd(13/900.*Dt, p.k7)
		//p.M = p.tempm
		//p.M[0] += ((31/300.*k1[0] + 0.*k2[0] + 0.*k3[0] + 0.*k4[0] + 61/225.*k5[0] - 2/9.*k6[0] + +13/900.*k7[0]) * Dt)
		//p.M[1] += ((31/300.*k1[1] + 0.*k2[1] + 0.*k3[1] + 0.*k4[1] + 61/225.*k5[1] - 2/9.*k6[1] + +13/900.*k7[1]) * Dt)
		//p.M[2] += ((31/300.*k1[2] + 0.*k2[2] + 0.*k3[2] + 0.*k4[2] + 61/225.*k5[2] - 2/9.*k6[2] + +13/900.*k7[2]) * Dt)
		T += (-5/6. + 1/6.) * Dt
	}
	if Demag {
		calculatedemag()
	}
	for _, p := range Lijst {
		temp := p.tempfield
		p.k8 = p.tau(temp)

		p.M = p.tempm.MAdd(2.*Dt, p.k1).MAdd(-53/6.*Dt, p.k4).MAdd(704/45.*Dt, p.k5).MAdd(-107/9.*Dt, p.k6).MAdd(67/90.*Dt, p.k7).MAdd(3.*Dt, p.k8)
		//p.M = p.tempm
		//p.M[0] += ((2.*k1[0] + 0.*k2[0] + 0.*k3[0] - 53/6.*k4[0] + 704/45.*k5[0] - 107/9.*k6[0] + 67/90.*k7[0] + 3.*k8[0]) * Dt)
		//p.M[1] += ((2.*k1[1] + 0.*k2[1] + 0.*k3[1] - 53/6.*k4[1] + 704/45.*k5[1] - 107/9.*k6[1] + 67/90.*k7[1] + 3.*k8[1]) * Dt)
		//p.M[2] += ((2.*k1[2] + 0.*k2[2] + 0.*k3[2] - 53/6.*k4[2] + 704/45.*k5[2] - 107/9.*k6[2] + 67/90.*k7[2] + 3.*k8[2]) * Dt)
		T += (-1/6. + 2/3.) * Dt
	}
	if Demag {
		calculatedemag()
	}
	for _, p := range Lijst {
		temp := p.tempfield
		p.k9 = p.tau(temp)

		p.M = p.tempm.MAdd(-91/108.*Dt, p.k1).MAdd(23/108.*Dt, p.k4).MAdd(-976/135.*Dt, p.k5).MAdd(311/54.*Dt, p.k6).MAdd(-19/60.*Dt, p.k7).MAdd(17/6.*Dt, p.k8).MAdd(-1/12.*Dt, p.k9)
		//p.M = p.tempm
		//p.M[0] += ((-91/108.*k1[0] + 0.*k2[0] + 0.*k3[0] + 23/108.*k4[0] - 976/135.*k5[0] + 311/54.*k6[0] - 19/60.*k7[0] + 17/6.*k8[0] - 1/12.*k9[0]) * Dt)
		//p.M[1] += ((-91/108.*k1[1] + 0.*k2[1] + 0.*k3[1] + 23/108.*k4[1] - 976/135.*k5[1] + 311/54.*k6[1] - 19/60.*k7[1] + 17/6.*k8[1] - 1/12.*k9[1]) * Dt)
		//p.M[2] += ((-91/108.*k1[2] + 0.*k2[2] + 0.*k3[2] + 23/108.*k4[2] - 976/135.*k5[2] + 311/54.*k6[2] - 19/60.*k7[2] + 17/6.*k8[2] - 1/12.*k9[2]) * Dt)
		T += (-1 / 3.) * Dt
	}
	if Demag {
		calculatedemag()
	}
	for _, p := range Lijst {
		temp := p.tempfield
		p.k10 = p.tau(temp)

		p.M = p.tempm.MAdd(2383/4100.*Dt, p.k1).MAdd(-341/164.*Dt, p.k4).MAdd(4496/1025.*Dt, p.k5).MAdd(-301/82.*Dt, p.k6).MAdd(2133/4100.*Dt, p.k7).MAdd(45/82.*Dt, p.k8).MAdd(45/164.*Dt, p.k9).MAdd(18/41.*Dt, p.k10)
		//p.M = p.tempm
		//p.M[0] += ((2383/4100.*k1[0] + 0.*k2[0] + 0.*k3[0] - 341/164.*k4[0] + 4496/1025.*k5[0] - 301/82.*k6[0] + 2133/4100.*k7[0] + 45/82.*k8[0] + 45/164.*k9[0] + 18/41.*k10[0]) * Dt)
		//p.M[1] += ((2383/4100.*k1[1] + 0.*k2[1] + 0.*k3[1] - 341/164.*k4[1] + 4496/1025.*k5[1] - 301/82.*k6[1] + 2133/4100.*k7[1] + 45/82.*k8[1] + 45/164.*k9[1] + 18/41.*k10[1]) * Dt)
		//p.M[2] += ((2383/4100.*k1[2] + 0.*k2[2] + 0.*k3[2] - 341/164.*k4[2] + 4496/1025.*k5[2] - 301/82.*k6[2] + 2133/4100.*k7[2] + 45/82.*k8[2] + 45/164.*k9[2] + 18/41.*k10[2]) * Dt)
		T += (2 / 3.) * Dt
	}
	if Demag {
		calculatedemag()
	}
	for _, p := range Lijst {
		temp := p.tempfield
		p.k11 = p.tau(temp)

		p.M = p.tempm.MAdd(3/205.*Dt, p.k1).MAdd(-6/41.*Dt, p.k6).MAdd(-3/205.*Dt, p.k7).MAdd(-3/41.*Dt, p.k8).MAdd(3/41.*Dt, p.k9).MAdd(6/41.*Dt, p.k10)
		//p.M = p.tempm
		//p.M[0] += ((3/205.*k1[0] + 0.*k2[0] + 0.*k3[0] + 0.*k4[0] + 0.*k5[0] - 6/41.*k6[0] - 3/205.*k7[0] - 3/41.*k8[0] + 3/41.*k9[0] + 6/41.*k10[0]) * Dt)
		//p.M[1] += ((3/205.*k1[1] + 0.*k2[1] + 0.*k3[1] + 0.*k4[1] + 0.*k5[1] - 6/41.*k6[1] - 3/205.*k7[1] - 3/41.*k8[1] + 3/41.*k9[1] + 6/41.*k10[1]) * Dt)
		//p.M[2] += ((3/205.*k1[2] + 0.*k2[2] + 0.*k3[2] + 0.*k4[2] + 0.*k5[2] - 6/41.*k6[2] - 3/205.*k7[2] - 3/41.*k8[2] + 3/41.*k9[2] + 6/41.*k10[2]) * Dt)
		T += (-1.) * Dt
	}
	if Demag {
		calculatedemag()
	}
	for _, p := range Lijst {
		temp := p.tempfield
		p.k12 = p.tau(temp)

		p.M = p.tempm.MAdd(-1777/4100.*Dt, p.k1).MAdd(-341/164.*Dt, p.k4).MAdd(4496/1025.*Dt, p.k5).MAdd(-289/82.*Dt, p.k6).MAdd(2193/4100.*Dt, p.k7).MAdd(51/82.*Dt, p.k8).MAdd(33/164.*Dt, p.k9).MAdd(12/41.*Dt, p.k10).MAdd(Dt, p.k12)
		//p.M = p.tempm
		//p.M[0] += ((-1777/4100.*k1[0] + 0.*k2[0] + 0.*k3[0] - 341/164.*k4[0] + 4496/1025.*k5[0] - 289/82.*k6[0] + 2193/4100.*k7[0] + 51/82.*k8[0] + 33/164.*k9[0] + 12/41.*k10[0] + 1.*k12[0]) * Dt)
		//p.M[1] += ((-1777/4100.*k1[1] + 0.*k2[1] + 0.*k3[1] - 341/164.*k4[1] + 4496/1025.*k5[1] - 289/82.*k6[1] + 2193/4100.*k7[1] + 51/82.*k8[1] + 33/164.*k9[1] + 12/41.*k10[1] + 1.*k12[1]) * Dt)
		//p.M[2] += ((-1777/4100.*k1[2] + 0.*k2[2] + 0.*k3[2] - 341/164.*k4[2] + 4496/1025.*k5[2] - 289/82.*k6[2] + 2193/4100.*k7[2] + 51/82.*k8[2] + 33/164.*k9[2] + 12/41.*k10[2] + 1.*k12[2]) * Dt)
		T += (1.) * Dt
	}
	for _, p := range Lijst {
		temp := p.tempfield
		p.k13 = p.tau(temp)

		p.M = p.tempm.MAdd(41/840.*Dt, p.k1).MAdd(34/105.*Dt, p.k6).MAdd(9/35.*Dt, p.k7).MAdd(9/35.*Dt, p.k8).MAdd(9/280.*Dt, p.k9).MAdd(9/280.*Dt, p.k10).MAdd(41/840.*Dt, p.k11)
		//p.M = p.tempm
		//p.tempm[0] += ((41/840.*k1[0] + 34/105.*k6[0] + 9/35.*k7[0] + 9/35.*k8[0] + 9/280.*k9[0] + 9/280.*k10[0] + 41/840.*k11[0]) * Dt)
		//p.tempm[1] += ((41/840.*k1[1] + 34/105.*k6[1] + 9/35.*k7[1] + 9/35.*k8[1] + 9/280.*k9[1] + 9/280.*k10[1] + 41/840.*k11[1]) * Dt)
		//p.tempm[2] += ((41/840.*k1[2] + 34/105.*k6[2] + 9/35.*k7[2] + 9/35.*k8[2] + 9/280.*k9[2] + 9/280.*k10[2] + 41/840.*k11[2]) * Dt)
		//sixth order solution

		p.M = p.tempm.MAdd(-41/840.*Dt, p.k1).MAdd(-41/840.*Dt, p.k11).MAdd(41/840.*Dt, p.k12).MAdd(41/840.*Dt, p.k13)
		//p.M[0] = p.tempm[0] + ((-41/840.*k1[0] - 41/840.*k11[0] + 41/840.*k12[0] + 41/840.*k13[0]) * Dt)
		//p.M[1] = p.tempm[1] + ((-41/840.*k1[1] - 41/840.*k11[1] + 41/840.*k12[1] + 41/840.*k13[1]) * Dt)
		//p.M[2] = p.tempm[2] + ((-41/840.*k1[2] - 41/840.*k11[2] + 41/840.*k12[2] + 41/840.*k13[2]) * Dt)
		//seventh order solution

		norm(p.M)
		norm(p.tempm)

		//the error is the difference between the two solutions
		//error := math.Sqrt(sqr(p.M[0]-p.tempm[0]) + sqr(p.M[1]-p.tempm[1]) + sqr(p.M[2]-p.tempm[2]))
		error := p.M.Sub(p.tempm).Len()

		//fmt.Println("error    :", error)
		if Adaptivestep {
			if error > maxtauwitht {
				maxtauwitht = error
			}
		}
		//when saving mdotH
		p.heff = p.b_eff(temp)

		T -= Dt
	}
}

//###########################################################################################################

func undobadstep(Lijst []*Particle) {
	for _, p := range Lijst {
		p.M = p.previousm
	}
	T -= Dt
}
