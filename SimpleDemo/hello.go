/*
package main

import (
	std "fmt"	//别名
	//"strconv"
)

const PI = 3.14

var name="yiworld"	//全局变量

type newType int	//一般类型声明

type gopher struct {	//结构声明

}

type golang interface {	//接口声明

}
*/
/*func main(){
	std.Println("hello world");
	var a=65;
	b:=strconv.Itoa(a)
	a,_=strconv.Atoi(b)
	std.Println(b);
	std.Println(a);
}*//*


*/
/*func test() func() {
	//x := 100
	std.Printf("测试2")

	return func() {
		std.Printf("测试1")
	}
}

func main() {
	f := test()
	f()
}*//*


func test() {
	defer recover() // 无效！
	defer std.Println(recover()) // 无效！
	defer func() {
		func() {
			println("defer inner")
			recover() // 无效！
		}()
	}()

	panic("test panic")
}

func main() {
	test()
}*/

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		defer wg.Done()

		for i := 0; i < 6; i++ {
			fmt.Println(i)

			if i == 3 {
				//fmt.Println("H3")
				runtime.Gosched()
			}
		}
	}()

	go func() {
		defer wg.Done()
		fmt.Println("Hello, World!")
	}()

	wg.Wait()
}

