package main

import (
	"zinx/net"
)

func main() {
	z := net.NewServer("zinx")
	z.Server()
}
