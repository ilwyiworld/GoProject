package main

//TCP客户端
import (
	"io"
	"log"
	"net"
	"os"
)
func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go mustCopy(os.Stdout, conn)	//读取并打印服务端的响应
	mustCopy(conn, os.Stdin)	//从标准输入流中读取内容并将其发送给服务器
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}