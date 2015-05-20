//Contains function to control the output of the program
package vinamax2

import ()

//var f *os.File
//var err error
//var twrite float64
//var locations []vector
//var filecounter int = 0
//var output_B_ext = false
//var output_Dt = false
//var output_nrmzpos = false
//var output_mdoth = false
//
////Sets the interval at which times the output table has to be written
//func Output(interval float64) {
//	outputcalled = true
//	f, err = os.Create(outdir + "/table.txt")
//	check(err)
//	defer f.Close()
//	writeheader()
//	outputinterval = interval
//	twrite = interval
//}
//
////checks the error
//func check(e error) {
//	if e != nil {
//		panic(e)
//	}
//}
//
////calculates the average magnetisation components of all particles
//func averages(lijst []*particle) vector {
//	avgs := vector{0, 0, 0}
//	for i := range lijst {
//		avgs[0] += lijst[i].M[0]
//		avgs[1] += lijst[i].M[1]
//		avgs[2] += lijst[i].M[2]
//	}
//	return avgs.times(1. / float64(len(lijst)))
//}
//
////calculates the average moments of all particles
//func averagemoments(lijst []*particle) vector {
//	avgs := vector{0, 0, 0}
//	totalvolume := 0.
//	for i := range lijst {
//		radius := lijst[i].r
//		volume := cube(radius) * 4. / 3. * math.Pi
//		totalvolume += volume
//		avgs[0] += lijst[i].M[0] * volume
//		avgs[1] += lijst[i].M[1] * volume
//		avgs[2] += lijst[i].M[2] * volume
//	}
//	//divide by total volume
//	return avgs.times(1. / totalvolume)
//}
//
////calculates the dotproduct of the average moments and the effective field of all particles
////this equals the losses
//func averagemdoth(lijst []*particle) float64 {
//	avg := 0.
//	for i := range lijst {
//		xcomp := lijst[i].M[0] * lijst[i].heff[0]
//		ycomp := lijst[i].M[1] * lijst[i].heff[1]
//		zcomp := lijst[i].M[2] * lijst[i].heff[2]
//		avg = (xcomp + ycomp + zcomp) / mu0
//	}
//	return (avg)
//}
//
////returns the number of particles with m_z larger than 0
//func nrmzpositive(lijst []*particle) int {
//	counter := 0
//	for i := range lijst {
//		if lijst[i].M[2] > 0. {
//			counter++
//		}
//	}
//	return counter
//}
//
////Writes the header in table.txt
//func writeheader() {
//	header := fmt.Sprintf("#t\t<mx>\t<my>\t<mz>")
//	_, err = f.WriteString(header)
//	check(err)
//	if output_B_ext {
//		header := fmt.Sprintf("\tB_ext_x\tB_ext_y\tB_ext_z")
//		_, err = f.WriteString(header)
//		check(err)
//	}
//	if output_Dt {
//		header := fmt.Sprintf("\tDt")
//		_, err = f.WriteString(header)
//		check(err)
//	}
//	if output_nrmzpos {
//		header := fmt.Sprintf("\tnrmzpos")
//		_, err = f.WriteString(header)
//		check(err)
//	}
//	if output_mdoth {
//		header := fmt.Sprintf("\tmdotH")
//		_, err = f.WriteString(header)
//		check(err)
//	}
//	for i := range locations {
//
//		header = fmt.Sprintf("\t(B_x\tB_y\tB_z)@(%v,%v,%v)", locations[i][0], locations[i][1], locations[i][2])
//		_, err = f.WriteString(header)
//		check(err)
//	}
//
//	header = fmt.Sprintf("\n")
//	_, err = f.WriteString(header)
//	check(err)
//
//}
//
////Adds the field at a specific location to the output table
//func Tableadd_b_at_location(x, y, z float64) {
//	tableaddcalled = true
//	if outputinterval != 0 {
//		log.Fatal("Output() should always come AFTER Tableadd_b_at_location()")
//	}
//	locations = append(locations, vector{x, y, z})
//
//}
//
////Writes the time and the vector of average magnetisation in the table
////+additional stuff if specified
//func write(avg vector) {
//	if twrite >= outputinterval && outputinterval != 0 {
//		string := fmt.Sprintf("%e\t%v\t%v\t%v", T, avg[0], avg[1], avg[2])
//		_, err = f.WriteString(string)
//		check(err)
//
//		if output_B_ext {
//			B_ext_x, B_ext_y, B_ext_z := B_ext(T)
//			string = fmt.Sprintf("\t%v\t%v\t%v", B_ext_x, B_ext_y, B_ext_z)
//			_, err = f.WriteString(string)
//			check(err)
//		}
//		if output_Dt {
//			string = fmt.Sprintf("\t%v", Dt)
//			_, err = f.WriteString(string)
//			check(err)
//		}
//		if output_nrmzpos {
//			string = fmt.Sprintf("\t%v", nrmzpositive(universe.lijst))
//			_, err = f.WriteString(string)
//			check(err)
//		}
//		if output_mdoth {
//			string = fmt.Sprintf("\t%v", averagemdoth(universe.lijst))
//			_, err = f.WriteString(string)
//			check(err)
//		}
//		for i := range locations {
//
//			string = fmt.Sprintf("\t%v\t%v\t%v", (demag(locations[i][0], locations[i][1], locations[i][2])[0]), (demag(locations[i][0], locations[i][1], locations[i][2])[1]), (demag(locations[i][0], locations[i][1], locations[i][2])[2]))
//			_, err = f.WriteString(string)
//			check(err)
//		}
//		_, err = f.WriteString("\n")
//		check(err)
//		twrite = 0.
//	}
//	twrite += Dt
//}
//
////Saves different quantities. At the moment only "geometry" and "m" are possible
//func Save(a string) {
//	//open file with unique name, used counter
//	name := fmt.Sprintf("%v%06v.txt", a, filecounter)
//	file, error := os.Create(outdir + "/" + name)
//	check(error)
//	defer file.Close()
//	filecounter += 1
//	switch a {
//
//	case "geometry":
//		{
//			// go through entire list of particles and print their position, radius and msat.
//			header := fmt.Sprintf("#position_x\tposition_y\tposition_z\tradius\tmsat\n")
//			_, err = file.WriteString(header)
//			check(err)
//
//			for i := range universe.lijst {
//				string := fmt.Sprintf("%v\t%v\t%v\t%v\t%v\n", universe.lijst[i].x, universe.lijst[i].y, universe.lijst[i].z, universe.lijst[i].r, universe.lijst[i].msat)
//				_, error = file.WriteString(string)
//				check(error)
//			}
//		}
//	case "m":
//		{
//			// loop over entire list with particles and print location, radius, msat and mag
//			header := fmt.Sprintf("#t= %v\n#position_x\tposition_y\tposition_z\tradius\tmsat\tm_x\tm_y\tm_z\n", T)
//			_, err = file.WriteString(header)
//			check(err)
//
//			for i := range universe.lijst {
//				string := fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", universe.lijst[i].x, universe.lijst[i].y, universe.lijst[i].z, universe.lijst[i].r, universe.lijst[i].msat, universe.lijst[i].M[0], universe.lijst[i].M[1], universe.lijst[i].M[2])
//				_, error = file.WriteString(string)
//				check(error)
//			}
//		}
//	default:
//		{
//			log.Fatal(a, " is not a quantitity that can be saved")
//		}
//	}
//}
//
////adds a quantity to the output table
//func Tableadd(a string) {
//	tableaddcalled = true
//	if outputinterval != 0 {
//		log.Fatal("Output() should always come AFTER Tableadd()")
//	}
//	switch a {
//	case "B_ext":
//		{
//			output_B_ext = true
//		}
//	case "Dt":
//		{
//			output_Dt = true
//		}
//	case "nrmzpos":
//		{
//			output_nrmzpos = true
//		}
//	case "mdoth":
//		{
//			output_mdoth = true
//		}
//
//	default:
//		{
//			log.Fatal(a, " is currently not addable to table")
//		}
//	}
//}
//
//func Writeintable(a string) {
//	string := fmt.Sprintf("%v\n", a)
//	_, err = f.WriteString(string)
//	check(err)
//}
