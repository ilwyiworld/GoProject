package pipeline

import (
	"sort"
	"io"
	"encoding/binary"
	"math/rand"
	"time"
	"fmt"
)
var startTime time.Time

func Init(){
	startTime=time.Now()
}

func ArraySource(a ...int) <-chan int{
	out:=make(chan int)
	go func(){
		for _,v:=range a  {
			out <- v
		}
		close(out)
	}()
	return out
}

func InMemSort(in <-chan int) <-chan int{
	out:=make(chan int,1024)
	go func(){
		//read into memory
		a:=[]int{}
		for v:=range in{
			a=append(a,v)
		}
		fmt.Println("Read done",time.Now().Sub(startTime))

		//Sort
		sort.Ints(a)
		fmt.Println("InMemSort done",time.Now().Sub(startTime))

		//Output
		for _,v:=range a{
			out<-v
		}
		close(out)
	}()
	return out
}

//两两归并
func Merge(in1,in2 <-chan int) <-chan int {
	out:=make(chan int,1024)
	go func() {
		v1,ok1 :=<-in1
		v2,ok2 :=<-in2
		for ok1 || ok2{
			if !ok2 || (ok1 && v1<=v2){
				// 1.in2无值，in1有值 将v1传入out
				// 2.in2有值，in1也有值 并且v1<=v2 将v1传入out
				out<-v1
				v1,ok1=<-in1
			}else{
				// 1.in2有值，in1无值，将v2传入out
				// 2.in2有值，in1有值，且v2<v1 将v2传入out
				out<-v2
				v2,ok2=<-in2
			}
		}
		close(out)
		fmt.Println("Merge done",time.Now().Sub(startTime))
	}()
	return out
}

//N个归并
func MergeN(inputs ... <-chan int) <-chan int{
	if len(inputs)==1{
		return inputs[0]
	}
	m:=len(inputs)/2
	//merge inputs[0,m) and inputs [m ...end)
	return Merge(MergeN(inputs[:m]...),MergeN(inputs[m:]...))
}

//从文件中最多读取chunkSize个数据
func ReaderSource(reader io.Reader,chunkSize int) <-chan int{
	out:=make(chan int,1024)
	go func() {
		buffer:=make([]byte,8)
		bytesRead:=0
		for  {
			n,err:=reader.Read(buffer)
			bytesRead+=n
			if n>0{
				v:=int(binary.BigEndian.Uint32(buffer))
				out<-v
			}
			if err!=nil ||
				(chunkSize!=-1 && bytesRead>=chunkSize){
				break
			}
		}
		close(out)
	}()
	return out
}

//从文件（网络）中写入数据
func WriterSink(writer io.Writer,in <-chan int){
	for v:=range in{
		buffer:=make([]byte,8)
		binary.BigEndian.PutUint32(buffer,uint32(v))
		writer.Write(buffer)
	}
}

//随机数据源
func RandomSource(count int ) <-chan int{
	out:=make(chan int)
	go func() {
		for  i:=0;i<count;i++{
			out<- rand.Int()
		}
		close(out)
	}()
	return out
}