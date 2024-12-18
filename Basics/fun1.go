package main

import (
	"fmt"
)

func myfirst() { //simple function
	fmt.Println("This is my first go function")
}

func NameFirst(fname string) { // single parameter
	fmt.Println("Hello ", fname, "\n Welcome to BigFix")
}

func NameAge(fname string, age int) { // single parameter
	fmt.Println("Hello ", fname, "and age is ", age)
}

func additions(num1 int, num2 int) int { // with return type only
	return num1 + num2
}

func myFunction(x int, y int) (result int) { // return with name value
	result = x + y
	return
}

func myFunctionS(x int, y int) (result int) { // return with name value
	result = x + y
	return
}

// Multiple return values in single function
func multireturn(x int, y string) (result int, txt1 string) {
	result = x + x
	txt1 = y + " World!"
	return
}

func main() {
	myfirst() // simple call
	NameFirst("Vishal")
	NameFirst("Mohini")
	NameFirst("Paras")
	NameFirst("Avadhi")
	NameAge("Vishal", 25)
	NameAge("BigFix", 26)
	fmt.Println(additions(1, 6))
	fmt.Println(myFunction(4, 5))
	sum := myFunctionS(12, 14) // storing return value of function in variable
	fmt.Println(sum)
	a, b := multireturn(1, "Hello")
	fmt.Println(a, b)

}
