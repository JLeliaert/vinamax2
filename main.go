package main

import "fmt"

var (
	root  Cell      // roots the entire FMM tree
	level [][]*Cell // for each level of the FMM tree: all cells on that level. Root = level 0
)

func main() {
	NLEVEL := 3
	level = make([][]*Cell, NLEVEL)

	root = Cell{size: Vector{1, 1, 1}}
	root.Divide(NLEVEL)

	l := level[2]
	for _, c := range l {
		fmt.Println("level", l, ":", c)
	}
}
