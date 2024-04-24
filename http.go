package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

// TODO 套接字、请求处理、执行cgi、多线程、进程通信
func acceptRequest(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("socket write failed,error:", err)
		}
		return
	}(conn)
	reader := bufio.NewReader(conn)
	//var data map[string]string
	//data = make(map[string]string)
	//var buf [1024]byte
	//var method [512]byte
	//length, err := reader.Read(buf[:])
	//if err != nil {
	//	fmt.Println("Read from client failed, error:", err)
	//	return
	//}
	_, err := resolveHttp(reader, conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(data)

}

func executeCGI() {
	//TODO 执行CGI
}

func renderFile() {
	//TODO 渲染网页文件
}

/*处理报文，返回一个map，存储了http报文的method,URL,http版本(暂未获取),headers,body
 */
func resolveHttp(reader *bufio.Reader, conn net.Conn) (map[string]string, any) {
	var data map[string]string
	data = make(map[string]string)

	//获取method,URL
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	parts := strings.Fields(string(line))
	if len(parts) < 2 {
		//TODO 抛出 500 error
		return nil, "server error 500"
	}
	data["method"] = parts[0]
	data["url"] = parts[1]
	data["headers"] = ""
	data["body"] = ""
	//获取headers
	line, err = reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	for {
		data["headers"] += string(line)
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		} else if string(line) == "\r\n" {
			break
		}
	}

	fmt.Printf("%q\n", data)
	return data, nil
}

//func readLine(reader *bufio.Reader) {
//	//TODO 读取一行http报文数据
//	//bytes, err := reader.ReadBytes('\n')
//	//line, _, err := reader.ReadLine()
//	line, _, err := reader.ReadLine()
//	if err != nil {
//		return
//	}
//	if err != nil {
//		return
//	}
//	if err != nil {
//		return
//	}
//	fmt.Printf("%q\n", string(line))
//}

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
	fmt.Println("start http: 8080")
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Accept failed, error:", err)
			continue
		}
		go acceptRequest(conn)
	}
}
