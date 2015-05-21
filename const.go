//This file contains all the constants and init files

package vinamax2

import (
	"math"
)

const (
	gamma0 = 1.7595e11          // Gyromagnetic ratio of electron, in rad/Ts
	mu0    = 4 * math.Pi * 1e-7 // Permeability of vacuum in Tm/A
	muB    = 9.2740091523E-24   // Bohr magneton in J/T
	kb     = 1.3806488E-23      // Boltzmann's constant in J/K
	qe     = 1.60217646E-19     // Electron charge in C
)

var outdir string // TODO: move to io.go

// TODO: races with verbose etc.
//func init() {
//	Log(`
//vinamax: a macrospin model to simulate magnetic nanoparticles
//Copyright (C) 2013  Jonathan Leliaert
//
//This program is free software: you can redistribute it and/or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, version 3 of the License
//
//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.
//
//You should have received a copy of the GNU General Public License
//along with this program.  If not, see [http:www.gnu.org/licenses/].
//
//contact: jonathan.leliaert@gmail.com
//
//
//`)
//
//}
