package main

import "fmt"

func main() {
	a:=A{B:B{Name:"B"},Name:"A"}
	fmt.Println(a)
	fmt.Println(a.Name,a.B.Name)
}

type A struct {
	B
	Name string
}

type B struct {
	Name string
}
