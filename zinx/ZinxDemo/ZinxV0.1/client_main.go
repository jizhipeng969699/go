package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	Conn, err := net.Dial("tcp", ":8999")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}
	defer Conn.Close()

	for {
		_, err := Conn.Write([]byte("hello zinx ........."))
		if err != nil {
			fmt.Println("client write err:", err)
			return
		}

		buf := make([]byte, 4096)
		n, err := Conn.Read(buf)
		if err != nil {
			fmt.Println("client read err:", err)
			return
		}
		fmt.Println(string(buf[:n]))

		time.Sleep(time.Second)
	}
}
