package main

import "fmt"

func main() {
	a:=A{}
	a.print()

	b:=B{}
	b.print()
}

func(a A) print(){
	fmt.Println("A")
}

func(b B) print(){
	fmt.Println("B")
}

type A struct {
	Name string
}

type B struct {
	Name string
}