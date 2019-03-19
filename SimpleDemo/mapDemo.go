package main

import "fmt"

func main() {
	var m map[int]string=make(map[int]string)
	//n:=make(map[int]string)
	m[1]="ok"
	fmt.Println(m)
	delete(m,1)
	fmt.Println(m)
	//fmt.Println(n[1])
	m[1]="ok"
	for _,v:=range m{
		fmt.Println(v)
	}
}
