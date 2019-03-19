package main

import "os"
import "fmt"

func main() {
	argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]
	arg := os.Args[3]
	fmt.Println(argsWithProg)
	fmt.Println(argsWithoutProg)
	fmt.Println(arg)
}

/*
$ ./CommandLineArguments a b c d
[./CommandLineArguments a b c d]
[a b c d]
c*/
