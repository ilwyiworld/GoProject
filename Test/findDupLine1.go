package main

import (
	"bufio"
	"os"
	"fmt"
)

//找出相同的行
func main() {
	counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)	//读入下一行，并移除行末的换行符
	for input.Scan() {
		counts[input.Text()]++
	}
	// NOTE: ignoring potential errors from input.Err()
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}
