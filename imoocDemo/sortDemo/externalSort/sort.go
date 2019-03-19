package main

import (
	"os"
	"../pipeline"
	"bufio"
	"fmt"
	"strconv"
)

func main() {
	p:=createNetworkPipeline("small.in", 512, 4)
	//p := createPipeline("small.in", 512, 4)
	writeToFile(p, "small.out")
	printFile("small.out")
}

func printFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	p := pipeline.ReaderSource(file, -1)        //所有数据都读
	for v := range p {
		fmt.Println(v)
	}

}

func writeToFile(p<-chan int, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	pipeline.WriterSink(writer, p)
}

func createPipeline(filename string,
		fileSize, chunkCount int) <-chan int {
	//fileSize文件总大小 chunkCount分块的数量
	chunkSize := fileSize / chunkCount        //每块文件的大小
	pipeline.Init()

	sortResults := [] <-chan int{}
	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}

		file.Seek(int64(i * chunkSize), 0)        //从第几块开始读
		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)
		sortResults = append(sortResults, pipeline.InMemSort(source))
	}
	return pipeline.MergeN(sortResults...)
}

func createNetworkPipeline(filename string,
		fileSize, chunkCount int) <-chan int {
	//fileSize文件总大小 chunkCount分块的数量
	chunkSize := fileSize / chunkCount        //每块文件的大小
	pipeline.Init()

	sortAddr := [] string{}
	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}

		file.Seek(int64(i * chunkSize), 0)        //从第几块开始读
		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)
		address := "localhost:" + strconv.Itoa(7000+i)
		pipeline.NetworkSink(address, pipeline.InMemSort(source))
		sortAddr = append(sortAddr, address)
	}

	sortResults := [] <- chan int{}
	for _, addr := range sortAddr {
		sortResults = append(sortResults, pipeline.NetworkSource(addr))
	}
	return pipeline.MergeN(sortResults...)
}


