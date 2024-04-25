package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

var SERVER_STRING string = "Server: tinyhttpd-go/1.1\r\n"

// TODO 套接字、请求处理、执行cgi、多线程、进程通信
func acceptRequest(conn net.Conn) {
	defer conn.Close()
	//处理http请求数据
	reader := bufio.NewReader(conn)
	var data map[string]string
	data = make(map[string]string)
	data, err := resolveHttp(reader, conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	//判断url是否正确
	filePath := "." + data["url"]
	fileInfo, err1 := os.Stat(filePath)
	if os.IsNotExist(err1) {
		notFound(conn)
		return
	}
	if fileInfo.IsDir() {
		filePath = filepath.Join(filePath, "index.html")
	}
	//依据数据做出处理
	if strings.ToLower(data["method"]) == "get" {
		_, err = renderFile(conn, filePath)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {

	}

	//var test []byte
	//test = append(test, []byte("haha")...)
	//conn.Write(test)
}

func headers(conn net.Conn) (bool, any) {
	var buf []byte
	s := "HTTP/2 200 OK\r\n" +
		"Content-Type: text/html\r\n" +
		SERVER_STRING +
		"\r\n"
	buf = append(buf, []byte(s)...)
	_, err := conn.Write(buf[:])
	if err != nil {
		return false, nil
	}
	return true, nil
}

func executeCGI() {
	//TODO 执行CGI
}

func renderFile(conn net.Conn, path string) (bool, error) {

	file, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	headers(conn)
	_, err = conn.Write(file)
	if err != nil {
		return false, err
	}
	return true, nil
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
	if err != nil {
		return nil, err
	}
	for { //TODO 如果报文是错误的，比如只有一半报文被传输，如何添加\r\n使得报文整体不出错
		line, err = reader.ReadBytes('\n')
		//fmt.Printf("current line-> %q\n", line)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if string(line) == "\r\n" || string(line) == "" {
			break
		}
		//fmt.Printf("header-> %q\n", line)
		data["headers"] = data["headers"] + string(line)
	}
	if strings.ToLower(data["method"]) != "post" && strings.ToLower(data["method"]) != "get" {
		//TODO 抛出405 not allowed method
		err := "405 Method not allowed"
		return nil, err
	}
	//var buf [1024]byte
	//length, _ := reader.Read(buf[:])
	//fmt.Printf("剩余内容%q\n", string(buf[:length]))

	if strings.ToLower(data["method"]) == "post" {
		var buf []byte
		n, err := reader.Read(buf)
		if err != nil {
			return nil, err
		}
		data["body"] = string(buf[:n])
	}
	//fmt.Printf("HTTP 报文：\n%q\n", data)
	return data, nil
}

func notFound(conn net.Conn) (bool, error) {
	//TODO 404error
	s := "HTTP/1.0 404 NOT FOUND\r\n" +
		"Content-Type: text/html\r\n" +
		SERVER_STRING +
		"\r\n" +
		"<HTML><TITLE>Not Found</TITLE>\r\n" +
		"<BODY><P>The server could not fulfill\r\n" +
		"your request because the resource specified\r\n" +
		"is unavailable or nonexistent.\r\n" +
		"</BODY></HTML>\r\n"
	var buf []byte
	buf = append(buf, []byte(s)...)
	//fmt.Println(string(buf))
	_, err := conn.Write(buf)
	if err != nil {
		return false, err
	}
	return true, nil
}

func serverInternalError() {
	//TODO 500
}

func main() {
	listen, err := net.Listen("tcp", "0.0.0.0:8080")
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
		acceptRequest(conn)
	}
}
