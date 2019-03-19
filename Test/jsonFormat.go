package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Movie struct {
	Title string
	Year int `json:"released"`
	Color bool `json:"color,omitempty"`//omitempty表示当Go语言结构体成员为空或零值时不生成JSON对象 这里false为零值
	Actors []string
}

type Response struct {
	Page   int      `json:"page"`
	Fruits []string `json:"fruits"`
}


var movies = []Movie{
	{Title: "Casablanca", Year: 1942, Color: false,
		Actors: []string{"Humphrey Bogart", "Ingrid Bergman"}},
	{Title: "Cool Hand Luke", Year: 1967, Color: true,
		Actors: []string{"Paul Newman"}},
	{Title: "Bullitt", Year: 1968, Color: true,
		Actors: []string{"Steve McQueen", "Jacqueline Bisset"}},
	// ...
}

func main() {
	data,err:=json.Marshal(&movies)
	if err!=nil{
		log.Fatalf("JSON marshaling failed: %s", err)
	}else{
		fmt.Printf("%s\n", data)
	}
	data,err = json.MarshalIndent(movies, "", " ")	//两个额外的字符串参数用于表示每一行输出的前缀和每一个层级的缩进
	if err!=nil{
		log.Fatalf("JSON marshaling failed: %s", err)
	}else{
		fmt.Printf("%s\n", data)
	}

	bolB, _ := json.Marshal(true)
	fmt.Println(string(bolB))
	intB, _ := json.Marshal(1)
	fmt.Println(string(intB))
	fltB, _ := json.Marshal(2.34)
	fmt.Println(string(fltB))
	strB, _ := json.Marshal("gopher")
	fmt.Println(string(strB))
	slcD := []string{"apple", "peach", "pear"}
	slcB, _ := json.Marshal(slcD)
	fmt.Println(string(slcB))
	mapD := map[string]int{"apple": 5, "lettuce": 7}
	mapB, _ := json.Marshal(mapD)
	fmt.Println(string(mapB))

	byt := []byte(`{"num":6.13,"strs":["a","b"]}`)
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)		//map[num:6.13 strs:[a b]]

	//num := dat["num"].(float64)
	num := dat["num"]
	fmt.Println(num)

	strs := dat["strs"].([]interface{})
	str1 := strs[0].(string)
	fmt.Println(str1)

	str := `{"page": 1, "fruits": ["apple", "peach"]}`
	res := Response{}
	json.Unmarshal([]byte(str), &res)
	fmt.Println(res)		//{1 [apple peach]}
	fmt.Println(res.Fruits[0])	//apple


	// In the examples above we always used bytes and strings
	// as intermediates between the data and JSON representation on standard out.
	// We can also stream JSON encodings directly to os.Writers
	// like os.Stdout or even HTTP response bodies.
	enc := json.NewEncoder(os.Stdout)
	d := map[string]int{"apple": 5, "lettuce": 7}
	enc.Encode(d)			//{"apple":5,"lettuce":7}
}
