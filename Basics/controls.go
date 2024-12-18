package main

import (
	"fmt"
)

func cond(a int) {
	if a%2 == 0 {
		fmt.Println("Even")
	} else {
		fmt.Println("Odd")
	}
}

func forl() {
	for i := 0; i <= 20; i += 10 {
		fmt.Println(i)
	}
}

func switchs(day int) {
	switch day {
	case 1:
		fmt.Println("Monday")
	case 2:
		fmt.Println("Tuesday")
	case 3:
		fmt.Println("Wednesday")
	case 4:
		fmt.Println("Thursday")
	case 5:
		fmt.Println("Friday")
	case 6:
		fmt.Println("Saturday")
	case 7:
		fmt.Println("Sunday")
	}
}

func main() {
	cond(4)
	forl()
	switchs(4)
}
