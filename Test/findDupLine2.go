package main

import (
	"bufio"
	"os"
	"fmt"
)

//找出相同的行
//读取标准输入或使用 os.Open打开文件
func main() {
	counts := make(map[string]int)
	files:=os.Args[1:]
	if len(files)==0{
		conutLines(os.Stdin,counts)
	}else{
		for _,arg := range files{
			f,err:=os.Open(arg)
			if err!=nil{
				fmt.Fprintf(os.Stderr,"dup: %v\n", err)
				continue
			}
			conutLines(f,counts)
			f.Close()
		}
		for line, n := range counts {
			if n > 1 {
				fmt.Printf("%d\t%s\n", n, line)
			}
		}
	}
}

func conutLines(f *os.File,count map[string]int){
	input:=bufio.NewScanner(os.Stdin)
	for input.Scan(){
		count(input.Text())++
	}
}
