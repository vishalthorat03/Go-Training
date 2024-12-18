package main

import (
	"fmt"
)

const PI = 3.14    // Untyped Constant
const A = "Vishal" //Typed Constatnt
const (            // Block Constant with many at one place
	B = 1
	C = 3
)

func main() {

	fmt.Println(PI)
	fmt.Println(A)
	fmt.Println(B)
	fmt.Println(C)
}
