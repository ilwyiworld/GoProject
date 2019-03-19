package main

import "fmt"

func main() {
	x,y:=1,2
	b:=[...]*int{&x,&y}
	a:=[20] int{19:1}
	var p *[20] int =&a
	fmt.Println(p)
	fmt.Println(b)

	c:=[2]int {1,2}
	d:=[2]int {1,3}
	fmt.Println(c==d)
}
