package main

import (
	"fmt"
)

func main() {
	var x = 1
	fmt.Print(x)
	var a, b string = "Hello", "World"
	fmt.Println(a, b)
	var i string = "Hello"
	var j int = 15
	fmt.Printf("i has value: %v and type: %T\n", i, i)
	fmt.Printf("j has value: %v and type: %T", j, j)
}
