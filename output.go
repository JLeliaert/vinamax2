//Contains function to control the output of the program
package vinamax2

//
import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

var (
	outputFile     *os.File
	err            error
	twrite         float64
	locations      []Vector
	filecounter    int = 0
	output_B_ext       = false
	output_Dt          = false
	output_nrmzpos     = false
	output_mdoth       = false
)

//Initialise the outputdir
func initOutput() {

	// make and clear output directory
	fname := os.Args[0]
	f2name := strings.Split(fname, "/") // TODO: use path.Split?
	outdir = fmt.Sprint(f2name[len(f2name)-1], ".out")
	os.Mkdir(outdir, 0775)
	dir, err3 := os.Open(outdir)
	files, _ := dir.Readdir(1)
	// clean output dir, copied from mumax
	if len(files) != 0 {
		filepath.Walk(outdir, func(path string, i os.FileInfo, err error) error {
			if path != outdir {
				os.RemoveAll(path)
			}
			return nil
		})
	}

	if err3 != nil {
		panic(err3)
	}

	outputFile, err = os.Create(outdir + "/table.txt")
	check(err)

}

// write to table (outputFile)
func writeTable(format string, msg ...interface{}) {
	if outputFile == nil {
		initOutput()
	}
	_, err := fmt.Fprintf(outputFile, format, msg...)
	check(err)
}

//Sets the interval at which Mul the output table has to be written
func Output(interval float64) {
	outputcalled = true

	writeheader()
	outputinterval = interval
	twrite = interval
}

//checks the error
func check(e error) {
	if e != nil {
		panic(e)
	}
}

//calculates the average magnetisation components of all Particles
func averages(lijst []*Particle) Vector {
	avgs := Vector{0, 0, 0}
	for i := range lijst {
		avgs[0] += lijst[i].M[0]
		avgs[1] += lijst[i].M[1]
		avgs[2] += lijst[i].M[2]
	}
	return avgs.Mul(1. / float64(len(lijst)))
}

//calculates the average moments of all Particles
func averagemoments(lijst []*Particle) Vector {
	avgs := Vector{0, 0, 0}
	totalvolume := 0.
	for i := range lijst {
		radius := lijst[i].r
		volume := radius * radius * radius * 4. / 3. * math.Pi
		totalvolume += volume
		avgs[0] += lijst[i].M[0] * volume
		avgs[1] += lijst[i].M[1] * volume
		avgs[2] += lijst[i].M[2] * volume
	}
	//divide by total volume
	return avgs.Mul(1. / totalvolume)
}

//calculates the dotproduct of the average moments and the effective field of all Particles
//this equals the losses
func averagemdoth(lijst []*Particle) float64 {
	avg := 0.
	for i := range lijst {
		xcomp := lijst[i].M[0] * lijst[i].heff[0]
		ycomp := lijst[i].M[1] * lijst[i].heff[1]
		zcomp := lijst[i].M[2] * lijst[i].heff[2]
		avg = (xcomp + ycomp + zcomp) / mu0
	}
	return (avg)
}

//returns the number of Particles with m_z larger than 0
func nrmzpositive(lijst []*Particle) int {
	counter := 0
	for i := range lijst {
		if lijst[i].M[2] > 0. {
			counter++
		}
	}
	return counter
}

//Writes the header in table.txt
func writeheader() {
	header := fmt.Sprintf("#t\t<mx>\t<my>\t<mz>")
	writeTable(header)
	if output_B_ext {
		writeTable("\tB_ext_x\tB_ext_y\tB_ext_z")
	}
	if output_Dt {
		writeTable("\tDt")
	}
	if output_nrmzpos {
		writeTable("\tnrmzpos")
	}
	if output_mdoth {
		writeTable("\tmdotH")
	}
	for i := range locations {
		writeTable("\t(B_x\tB_y\tB_z)@(%v,%v,%v)", locations[i][0], locations[i][1], locations[i][2])
	}
	writeTable("\n")
}

//Adds the field at a specific location to the output table
func Tableadd_b_at_location(x, y, z float64) {
	tableaddcalled = true
	if outputinterval != 0 {
		Fatal("Output() should always come AFTER Tableadd_b_at_location()")
	}
	locations = append(locations, Vector{x, y, z})

}

//Writes the time and the Vector of average magnetisation in the table
//+additional stuff if specified
func write(avg Vector) {
	if twrite >= outputinterval && outputinterval != 0 {
		writeTable("%e\t%v\t%v\t%v", T, avg[0], avg[1], avg[2])

		if output_B_ext {
			B_ext_x, B_ext_y, B_ext_z := B_ext(T)
			writeTable("\t%v\t%v\t%v", B_ext_x, B_ext_y, B_ext_z)
		}
		if output_Dt {
			writeTable("\t%v", Dt)
		}
		if output_nrmzpos {
			writeTable("\t%v", nrmzpositive(Particles))
		}
		if output_mdoth {
			writeTable("\t%v", averagemdoth(Particles))
		}
		for i := range locations {

			//TODO	string = fmt.Sprintf("\t%v\t%v\t%v", (demag(locations[i][0], locations[i][1], locations[i][2])[0]), (demag(locations[i][0], locations[i][1], locations[i][2])[1]), (demag(locations[i][0], locations[i][1], locations[i][2])[2]))
			writeTable("\t%v", i)
		}
		writeTable("\n")
		twrite = Dt
	}
	twrite += Dt
}

//Saves different quantities. At the moment only "geometry" and "m" are possible
func Save(a string) {
	//open file with unique name, used counter
	name := fmt.Sprintf("%v%06v.txt", a, filecounter)
	file, error := os.Create(outdir + "/" + name)
	check(error)
	defer file.Close()
	filecounter += 1
	switch a {

	case "geometry":
		{
			// go through entire list of Particles and print their position, radius and msat.
			header := fmt.Sprintf("#position_x\tposition_y\tposition_z\tradius\tmsat\n")
			_, err = file.WriteString(header)
			check(err)

			for i := range Particles {
				string := fmt.Sprintf("%v\t%v\t%v\t%v\t%v\n", Particles[i].center[0], Particles[i].center[1], Particles[i].center[2], Particles[i].r, Particles[i].msat)
				_, error = file.WriteString(string)
				check(error)
			}
		}
	case "m":
		{
			// loop over entire list with Particles and print location, radius, msat and mag
			header := fmt.Sprintf("#t= %v\n#position_x\tposition_y\tposition_z\tradius\tmsat\tm_x\tm_y\tm_z\n", T)
			_, err = file.WriteString(header)
			check(err)

			for i := range Particles {
				string := fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", Particles[i].center[0], Particles[i].center[1], Particles[i].center[2], Particles[i].r, Particles[i].msat, Particles[i].M[0], Particles[i].M[1], Particles[i].M[2])
				_, error = file.WriteString(string)
				check(error)
			}
		}
	default:
		{
			Fatal(a, " is not a quantitity that can be saved")
		}
	}
}

//adds a quantity to the output table
func Tableadd(a string) {
	tableaddcalled = true
	if outputinterval != 0 {
		Fatal("Output() should always come AFTER Tableadd()")
	}
	switch a {
	case "B_ext":
		{
			output_B_ext = true
		}
	case "Dt":
		{
			output_Dt = true
		}
	case "nrmzpos":
		{
			output_nrmzpos = true
		}
	case "mdoth":
		{
			output_mdoth = true
		}

	default:
		{
			Fatal(a, " is currently not addable to table")
		}
	}
}
