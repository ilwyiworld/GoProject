package main

import "fmt"

func main() {
	switchFunc := func(i interface{}) {
		switch i.(type){
		case bool:
			fmt.Println("bool")
		case int:
			fmt.Println("int")
		case float64:
			fmt.Println("float64")
		default:
			fmt.Println("other")
		}
	}
	switchFunc(true)
}
