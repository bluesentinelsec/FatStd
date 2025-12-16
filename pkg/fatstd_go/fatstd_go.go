package main

import "C"

//export fatstd_go_add
func fatstd_go_add(a, b C.int) C.int {
	return a + b
}

func main() {}

