package main

import "fmt"

func main() {
	root := Cell{size: Vector{1, 1, 1}}
	root.divide(2)

	fmt.Println(root)
	for i, c := range root.child[0].child {
		fmt.Println("child", i, ":", c)
	}
}
