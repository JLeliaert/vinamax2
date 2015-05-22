package vinamax2

// Benchmarks

import (
	"testing"
)

// demag, 5 levels, 4096 particles, 0th order
func BenchmarkFMM5Levels0th(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()

	initBenchWorld(5)
	FMMOrder = 0

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		CalcDemag()
	}
}

// demag, 5 levels, 4096 particles, 0th order, iterative
func BenchmarkFMM5LevelsIter(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()

	initBenchWorld(5)
	FMMOrder = 0

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		CalcDemagIter()
	}
}

// demag, 5 levels, 4096 particles, 0th order, iterative
func BenchmarkFMM5Iter0th(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()

	initBenchWorld(5)
	FMMOrder = 0

	b.StartTimer()
	for i := 0; i < b.N; i++ {

	}
}

// demag, 5 levels, 4096 particles, 1st order
func BenchmarkFMM5Levels1st(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()

	initBenchWorld(5)
	FMMOrder = 1

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		CalcDemag()
	}
}

func BenchmarkDipoleField(b *testing.B) {
	m := Vector{1, 2, 3}
	r := Vector{3, 4, 5}
	for i := 0; i < b.N; i++ {
		DipoleField(m, r)
	}
}

// benchmark AddPartnerFields, the most time-consuming FMM stage, 0th-order
func BenchmarkAddPartnerFields0th(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()

	hotcell := initBenchWorld(5)
	FMMOrder = 0

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		hotcell.addPartnerFields0()
	}
}

// benchmark AddPartnerFields, the most time-consuming FMM stage, 1st order
func BenchmarkAddPartnerFields1st(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()

	hotcell := initBenchWorld(5)
	FMMOrder = 1

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		hotcell.addPartnerFields1()
	}
}

func BenchmarkInitFMM(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		initBenchWorld(5)
	}
}

// init standard world for benchmarking: with NLEVEL levels,
// one particle at each base cell center.
func initBenchWorld(NLEVEL int) *Cell {
	worldSize := Vector{1, 1, 1}

	InitFMM(worldSize, NLEVEL)

	// place particles with m=0 , as field probes
	baseLevel := Level[NLEVEL-1]
	for _, c := range baseLevel {
		AddParticle(&Particle{M: Vector{0, 0, 0}, center: c.center})
	}

	// place one magneticed particle as source
	hotcell := baseLevel[30]
	AddParticle(&Particle{M: Vector{1, 2, 3}, center: hotcell.center})
	return hotcell
}
