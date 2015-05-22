package vinamax2

import (
	"math"
)

// RMS difference between CalcDemag and brute-force method
func DemagError() float64 {
	CalcDemag()

	error := 0.0
	for _, p := range Particles {
		error += p.BruteDemag().Sub(p.Bdemag()).Len2()
	}

	error /= float64(len(Particles))
	return math.Sqrt(error)
}
