package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Id int
	Name string
	Age int
}

type Manager struct {
	User
	title string
}

func (u User) Hello(){
	fmt.Println("Hello world.")
}

func main() {
	u:=User{1,"ok",12}
	info(u)

	m:=Manager{User:User{1,"bad",20},title:"This is title"}
	t:=reflect.TypeOf(m)
	fmt.Printf("%#v\n",t.FieldByIndex([]int{0,0}))

	x:=123
	fmt.Println(x)
	v:=reflect.ValueOf(&x)
	v.Elem().SetInt(1234)
	fmt.Println(x)
}

func info(o interface{}){
	t:=reflect.TypeOf(o)
	fmt.Println("Type:"+t.Name())

	if k:=t.Kind() ;k!=reflect.Struct{
		return
	}

	v:=reflect.ValueOf(o)
	fmt.Println("Fields:")

	for i := 0; i < t.NumField(); i++ {
		f:=t.Field(i)
		val:=v.Field(i).Interface()
		fmt.Printf("%6s: %v = %v\n",f.Name,f.Type,val)
	}

	for i := 0; i < t.NumMethod(); i++ {
		m:=t.Method(i)
		fmt.Printf("%6s: %v\n",m.Name,m.Type)
	}
}
