package main

import (
	"fmt"
)

func main() {
	ch:=make(chan string)
	for i:=0; i<5000 ; i++ {
		go printHelloworld(i,ch)
	}
	for  {
		msg:=<-ch
		fmt.Println(msg)
	}
}

func printHelloworld(i int,ch chan string)  {
	for  {
		ch <- fmt.Sprintf("hello world from goruntine %d!\n",i)
	}
}
