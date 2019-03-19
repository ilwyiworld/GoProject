package main

import "fmt"

type person struct{
	Name string
	Age int
}
func main() {
	a:=&person{"yiworld",12}
	b:=person{
		Name:"yiiii",
		Age:15,
	}
	fmt.Println(a)
	fmt.Println(b)
	a.Name="yi"
	a.Age=14
	fmt.Println(a)
	changeAge(a)
	fmt.Println(a)

	c:= &struct {
		Name string
		Age int
	}{
		Name:"cYiworld",
		Age:22,
	}
	fmt.Println(c)
}

func changeAge(per *person){
	per.Age=20
	fmt.Println("changeAge",per)
}
