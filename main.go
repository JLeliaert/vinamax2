package main

import (
	"fmt"
	"log"
	"math"
	"time"
)

var (
	root  Cell      // roots the entire FMM tree
	level [][]*Cell // for each level of the FMM tree: all cells on that level. Root = level 0
)

func main() {

	for NLEVEL := 1; NLEVEL <= 256; NLEVEL++ {
		fmt.Print(math.Pow(8, float64(NLEVEL-1)), " ")
		flops = 0
		level = make([][]*Cell, NLEVEL)

		root = Cell{size: Vector{1, 1, 1}}
		log.Println("dividing")
		root.Divide(NLEVEL)
		log.Println("finding partners")
		root.FindPartners(level[0])
		log.Println("start")
		start := time.Now()
		root.UpdateM()
		root.UpdateB()

		fmt.Println(flops, " ", float64(time.Since(start).Nanoseconds())/1e9)

		//for l := range level {
		//	fmt.Println("level", l)
		//	for _, c := range level[l] {
		//		fmt.Println(c)
		//	}
		//	fmt.Println()
		//}

	}
}
