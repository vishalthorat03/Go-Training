package main

import (
	"fmt"
)

var a bool = true
var b int = 1
var c float32 = 3.14
var d string = "Vishal"
var e bool // will return false because the default value is false
var x int = 500
var y int = -4500
var w uint = 500
var s uint = 4500
var l float32 = 123.78
var m float32 = 3.4e+38
var t float64 = 1.7e+308
var txt1 string = "Hello!"
var txt2 string

func main() {
	txt3 := "Worl d 1"
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(d)
	fmt.Println(e)
	fmt.Printf("\n Type: %T, value: %v", x, x)
	fmt.Printf("\n Type: %T, value: %v", y, y)
	fmt.Printf("\n Type: %T, value: %v", w, w)
	fmt.Printf("\n Type: %T, value: %v", s, s)
	fmt.Printf("\n Type: %T, value: %v", l, l)
	fmt.Printf("\n Type: %T, value: %v", m, m)
	fmt.Printf("\n Type: %T, value: %v", t, t)
	fmt.Printf("Type: %T, value: %v\n", txt1, txt1)
	fmt.Printf("Type: %T, value: %v\n", txt2, txt2)
	fmt.Printf("Type: %T, value: %v\n", txt3, txt3)

}
