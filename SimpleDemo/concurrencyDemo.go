package main

import (
	"fmt"
)

func main() {
	//go Go()
	//time.Sleep(2*time.Second)

	//c:=make(chan bool)  //双向通道

	c:=make(chan bool)
	go func(){
		fmt.Println("go go go!!!")
		c <-true
		close(c)
	}()
	for v:=range c{
		fmt.Println(v)
	}
}

func Go(){
	fmt.Println("Go Go Go")
}
