package main

import (
	"flag"
	"fmt"
)

func main() {
	wordPtr := flag.String("word", "foo", "a string")
	numbPtr := flag.Int("numb", 42, "an int")
	boolPtr := flag.Bool("fork", false, "a bool")

	// It’s also possible to declare an option
	// that uses an existing var declared elsewhere in the program.
	// Note that we need to pass in a pointer to the flag declaration function.
	var svar string
	flag.StringVar(&svar, "svar", "bar", "a string var")

	flag.Parse()

	fmt.Println("word:", *wordPtr)
	fmt.Println("numb:", *numbPtr)
	fmt.Println("fork:", *boolPtr)
	fmt.Println("svar:", svar)
	fmt.Println("tail:", flag.Args())
}

/*
$ ./CommandLineFlag -word=opt -numb=7 -fork -svar=flag
word: opt
numb: 7
fork: true
svar: flag
tail: []

$ ./CommandLineFlag -word=opt
word: opt
numb: 42
fork: false
svar: bar
tail: []

$ ./CommandLineFlag -word=opt a1 a2 a3
word: opt
...
tail: [a1 a2 a3]

$ ./CommandLineFlag -word=opt a1 a2 a3 -numb=7
word: opt
numb: 42
fork: false
svar: bar
tail: [a1 a2 a3 -numb=7]

$ ./CommandLineFlag -h
Usage of ./CommandLineFlag:
-fork=false: a bool
-numb=42: an int
-svar="bar": a string var
-word="foo": a string

$ ./CommandLineFlag -wat
flag provided but not defined: -wat
Usage of ./CommandLineFlag:
...*/
