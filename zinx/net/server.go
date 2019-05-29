/**
server模块的实现层
*/
package net

import (
	"fmt"
	"net"
	"zinx/config"
	"zinx/ziface"
)

type Server struct {
	//服务器ip
	IPVersion string
	IP string
	//服务器port
	Port int
	//服务器名称
	Name string

	//多路由的消息管理模块
	MsgHandler ziface.IMsgHandler

	//链接管理模块
	connMgr ziface.IConnManager

	//该server创建链接之后自动调用Hook函数
	OnConnStart func(conn ziface.IConnection)
	//该server销毁链接之前自动调用的Hook函数
	OnConnStop func(conn ziface.IConnection)
}


//初始化的New方法
func NewServer(name string) ziface.IServer{
	s := &Server{
		Name:config.GlobalObject.Name,
		IPVersion:"tcp4",
		IP:config.GlobalObject.Host,
		Port:config.GlobalObject.Port,
		MsgHandler:NewMsgHandler(),
		connMgr:NewConnManager(),
	}

	return s
}

//启动服务器
//原生socket 服务器编程
func (s *Server) Start() {
	fmt.Printf("[start] Server Linstenner at IP :%s, Port :%d, is starting..\n", s.IP, s.Port)

	//0 启动worker工作池
	s.MsgHandler.StartWorkerPool()

	//1 创建套接字  ：得到一个TCP的addr
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("resolve tcp addr error:", err)
		return
	}
	//2 监听服务器地址
	listenner, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("listen ", s.IPVersion, " err , ", err)
		return
	}


	//生成id的累加器
	var cid uint32
	cid = 0

	//3 阻塞等待客户端发送请求，
	go func() {
		for {
			//应该是永久存在的。

			//阻塞等待客户端请求,
			conn, err := listenner.AcceptTCP()//只是针对TCP协议
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			//创建一个Connection对象
			//判断当前server链接数量是否已经最大值
			if s.connMgr.Len() >= int(config.GlobalObject.MaxConn) {
				//当前链接已经满了
				fmt.Println("---> Too many Connection MAxConn = ", config.GlobalObject.MaxConn)
				conn.Close()
				continue
			}
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++


			//此时conn就已经和对端客户端连接 //处理与客户端的读写业务
			go dealConn.Start()
		}
	}()

}
//停止服务器
func (s *Server) Stop() {
	//服务器停止  应该清空当前全部的链接
	s.connMgr.ClearConn()
}
//运行服务器
func (s *Server )Serve() {
	//启动server的监听功能
	s.Start()//并不希望他永久的阻塞

	//TODO  做一些其他的扩展
	//阻塞//告诉CPU不再需要处理的，节省cpu资源
	select{}//main函数不退出  //main函数 阻塞在这
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("Add Router SUCC!! msgID = ", msgId)
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.connMgr
}



//注册 创建链接之后 调用的 Hook函数 的方法
func(s *Server) AddOnConnStart(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
}
//注册 销毁链接之前调用的Hook函数 的方法
func (s *Server) AddOnConnStop(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStop = hookFunc
}


//调用 创建链接之后的HOOK函数的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}
//调用 销毁链接之前调用的HOOk函数的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> Call OnConnStop()...")
		s.OnConnStop(conn)
	}
}