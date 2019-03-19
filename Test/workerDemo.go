package main

import (
	"fmt"
	"time"
)

func main() {
	jobs:=make(chan int,100)
	results:=make(chan int,100)

	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	for a := 1; a <= 5; a++ {
		fmt.Println(<-results)
	}

}

func worker(id int,jobs <-chan int,result chan<- int){
	for i:=range jobs{
		fmt.Println("Work",id,"started")
		time.Sleep(time.Second)
		fmt.Println("Work",id,"finished")
		result<-i*2
	}
}
