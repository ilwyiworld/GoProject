package main

import (
	"net/http"
	"log"
	"fmt"
)

func main() {
	mux:=http.NewServeMux()
	mux.Handle("/",&myHandler{})
	mux.HandleFunc("/bye",sayBye)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello,this is verison 1"))
		r.ParseForm()  //解析参数，默认是不会解析的
		fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
		fmt.Println("path", r.URL.Path)
		fmt.Println("scheme", r.URL.Scheme)
		fmt.Println(r.Form["url_long"])
	})

	http.HandleFunc("/bye", sayBye)

	log.Println("Starting server version 1")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

type myHandler struct{

}

func (_ *myHandler) ServeHTTP( w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Hello,this is verison 2"))
}


func sayBye(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Bye bye,this is verison 1"))
}

