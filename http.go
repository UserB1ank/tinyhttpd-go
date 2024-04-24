package main

import (
	"bufio"
	"fmt"
	"net"
)

// TODO 套接字、请求处理、执行cgi、多线程、进程通信
func acceptRequest(conn net.Conn) {
	reader := bufio.NewReader(conn)
	var buf [128]byte
	n, err := reader.Read(buf[:])
	if err != nil {
		fmt.Println("Read from client failed, error:", err)
		return
	}
	recvStr := string(buf[:n])
	fmt.Println("收到client发来的信息:", recvStr)
	_, err = conn.Write([]byte(recvStr))
	if err != nil {
		fmt.Println("Send data to client failed, error:", err)
		return
	}
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
