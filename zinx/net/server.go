package net

import (
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

//server 模块的实现层

type Server struct {
	//服务器ip
	IPVersion string
	IP        string
	//服务器port
	Port int
	//服务器名称
	Name string
}

//初始化的New方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}

	return s
}

//启动服务器
func (s *Server) Start() {
	//创建套接子   得到一个tcp的addr
	TCPAddr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("resolve tcp addr error:", err)
		return
	}
	//监听服务器地址
	TCPListener, err := net.ListenTCP(s.IPVersion, TCPAddr)
	if err != nil {
		fmt.Println("listentcp err", err)
		return
	}
	//defer TCPListener.Close()

	//阻塞等待客户端发送请求
	go func() {
		for {
			//阻塞等待客户端发送请求
			Conn, err := TCPListener.Accept()
			if err != nil {
				fmt.Println("Accept err", err)
				return
			}

			//此时conn就已经和对段客户端连接
			go func() {
				//defer Conn.Close()
				//4 客户端有数据请求 处理客户端业务（读，写）
				for {
					buf := make([]byte, 4096)
					n, err := Conn.Read(buf)
					if err != nil && err != io.EOF {
						fmt.Println("read buf err", err)
						break
					}
					if n == 0 {
						fmt.Println("client is over", err)
						break
					}

					//回显示功能(属于业务)
					//msg := []byte("okokok")
					_, err = Conn.Write(buf[:n])
					if err != nil {
						fmt.Println("conn write err:", err)
						return
					}
				}
			}()
		}
	}()

}

//停止服务器
func (s *Server) Stop() {

}

//关闭服务器
func (s *Server) Server() {
	//启动server的监听功能
	s.Start() //并不希望他永久阻塞

	//阻塞//告诉CPU不再需要处理的，节省cpu资源
	select {}
}
