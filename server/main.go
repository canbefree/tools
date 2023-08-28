package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"os"

	"github.com/SuperJourney/tools/helper"
)

// for debug 非线程安全
var clientList []net.Conn

func main() {
	var host = "localhost:9090"
	l, err := net.Listen("tcp", host)
	helper.PaincIfErr(err)

	helper.Println("服务端启动")

	go func() {
		// 从终端读取输入并发送消息
		scanner := bufio.NewScanner(os.Stdin)
		helper.Printf("准备读取输入数据")
		for scanner.Scan() {
			message := scanner.Text()
			for k, conn := range clientList {
				_, err := conn.Write([]byte(message + "\n"))
				helper.Printf("发送数据:%s", conn.RemoteAddr().String())
				if err != nil {
					helper.Printf("删除索引 %v", k)
					clientList = append(clientList[:k], clientList[k+1:]...)
					helper.Printf("当前conn #%v", conn.RemoteAddr())
				}
			}
		}
	}()

	for {
		conn, err := l.Accept()
		helper.PaincIfErr(err)
		clientList = append(clientList, conn)
		helper.Printf("新客户端连接:%s", conn.RemoteAddr())
		go Handle(conn)
	}

}

// 处理接受的消息
func Handle(conn net.Conn) {
	defer conn.Close()
	remoteAddr := conn.RemoteAddr()
	for {
		var b []byte = make([]byte, 10)
		l, err := conn.Read(b)
		if err != nil {
			if errors.Is(err, io.EOF) {
				helper.Printf("客户端断开连接 %s", conn.RemoteAddr())
				break
			}
		}
		helper.PaincIfErr(err)
		helper.Printf("接受数据 remote addr %v,len,%v,data:%s", remoteAddr, l, b)
	}
}
