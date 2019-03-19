package main

import "fmt"

func main() {
	a:=A
	a()

	b:=func(){	//匿名函数
		fmt.Println("Func B")
	}
	b()

	c:=closure(1)
	fmt.Println(c(2))
	fmt.Println(c(3))

	for i:=0;i<3; i++ {
		defer func(){
			fmt.Println(i)
		}()
	}

	A()
	B()
	C()
}

func closure(x int) func(int) int{
	fmt.Println("%p",&x)
	return func(y int) int{
		fmt.Println("%p",&x)
		return x+y
	}
}

func A(){
	fmt.Println("Func A")
}

func B(){
	defer func(){
		if err:=recover();err!=nil {
			fmt.Println("Recover in B")
		}
	}()
	panic("Panic in B")
}

func C(){
	fmt.Println("Panic in C")
}

