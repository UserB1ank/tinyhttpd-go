package main

import (
	"bufio"
	"fmt"
	"net"
)

// TODO 套接字、请求处理、执行cgi、多线程、进程通信
func acceptRequest(conn net.Conn) {
	reader := bufio.NewReader(conn)
	var buf [1024]byte
	var method [512]byte
	len, err := reader.Read(buf[:])
	if err != nil {
		fmt.Println("Read from client failed, error:", err)
		return
	}

}

func executeCGI() {
	//TODO 执行CGI
}

func renderFile() {
	//TODO 渲染网页文件
}

func readHeaders(conn net.Conn, len int) {
	//TODO 读取http头
}

func readBody() {
	//TODO 读取报文体
}

func readLine(conn net.Conn, len int) {
	//TODO 读取一行http报文数据
}

func notFound() {
	//TODO 404error
}

func serverInternalError() {
	//TODO 500
}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Listen failed, error:", err)
		return
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Accept failed, error:", err)
			continue
		}
		go acceptRequest(conn)
	}
}
