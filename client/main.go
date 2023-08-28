package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"os"

	"github.com/SuperJourney/tools/helper"
)

func main() {
	address := "localhost:9090"
	conn, err := net.Dial("tcp", address)
	helper.PaincIfErr(err)
	defer conn.Close()
	go handleConn(conn)

	// 从终端读取输入并发送消息
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		_, err := conn.Write([]byte(message + "\n"))
		helper.Println("发送数据")
		if err != nil {
			log.Println("连接已关闭")
			break
		}
	}
}

func handleConn(conn net.Conn) {
	for {
		var b = make([]byte, 20)
		_, err := conn.Read(b)
		if errors.Is(err, io.EOF) {
			break
		}
		helper.Printf("recv from server data:%s", string(b))
	}
}
