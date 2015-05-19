package main

type Particle struct {
	m      Vector  // magnetization
	b      Vector  // demag field
	u      Vector  // anisotropy direction (normalized)
	size   float64 // radius
	center Vector  // particle position
}
