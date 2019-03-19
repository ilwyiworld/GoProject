package main

import (
	"../pipeline"
	"fmt"
	"os"
	"bufio"
)
func main() {
	//mergeDemo()
	const filename="small.in"
	const count=64

	//生成随机数文件
	file,err:=os.Create(filename)
	if err!=nil{
		panic(err)
	}
	defer file.Close()
	p:=pipeline.RandomSource(count)
	writer:=bufio.NewWriter(file)
	pipeline.WriterSink(writer,p)
	writer.Flush()


	//写入随机数到chan int
	file,err=os.Open(filename)
	if err!=nil{
		panic(err)
	}
	defer file.Close()
	p=pipeline.ReaderSource(bufio.NewReader(file),-1)
	countNum:=0
	for value:=range p{
		fmt.Println(value)
		countNum++
		if countNum>100{	//只打印100个
			break
		}
	}
}

func mergeDemo(){
	p:=pipeline.Merge(
		pipeline.InMemSort(pipeline.ArraySource(2,58,74,21,-4,5)),
		pipeline.InMemSort(pipeline.ArraySource(42,2,72,234,-42331,523)))
	for  v:=range p{
		fmt.Println(v)
	}
}
