/**
* @Author: Aceld(刘丹冰)
* @Date: 2019/5/22 17:19
* @Mail: danbing.at@gmail.com
*/
package net

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/config"
	"zinx/ziface"
)

//具体的TCP链接模块
type Connection struct {
	//当前链接是属于哪个server创建
	server ziface.IServer

	//当前链接的原生套接字
	Conn *net.TCPConn

	//链接ID
	ConnID uint32

	//当前的链接状态
	isClosed bool

	//消息管理模块 多路由
	MsgHandler ziface.IMsgHandler

	//添加一个 Reader和Writer通信的Channel
	msgChan chan []byte

	//创建一个Channel  用来Reader通知Writer conn已经关闭，需要退出的消息
	writerExitChan chan bool

	//当前连接模块所具备的一些属性集合
	property map[string]interface{}

	//保护当前的property的锁
	propertyLock sync.RWMutex
}

/*
初始化链接方法
 */
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, handler ziface.IMsgHandler) ziface.IConnection {
	c := &Connection{
		server: server,
		Conn:   conn,
		ConnID: connID,
		//handleAPI:callback_api,
		MsgHandler:     handler,
		isClosed:       false,
		msgChan:        make(chan []byte), //初始化Reader Writer通信的Channel
		writerExitChan: make(chan bool),

		property: make(map[string]interface{}),
	}

	//当已经成功创建一个链接的时候，添加到链接管理器中
	c.server.GetConnMgr().Add(c)

	return c
}

//针对链接读业务的方法
func (c *Connection) StartReader() {
	//从对端读数据
	fmt.Println("[Reader Goroutine isStarted]...")
	defer fmt.Println("[Reader Goroutine Stop...] connID = ", c.ConnID, "Reader is exit, remote addr is = ", c.GetRemoteAddr().String())
	defer c.Stop()

	for {
		/*
		buf := make([]byte, config.GlobalObject.MaxPackageSize)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err", err)
			break
		}
		*/
		//创建拆包封包的对象
		dp := NewDataPack()

		//读取客户端消息的头部
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.Conn, headData); err != nil {
			fmt.Println("read msg head error", err)
			break
		}

		//根据头部 获取数据的长度，进行第二次读取
		msg, err := dp.UnPack(headData) //将msg 头部信息填充满
		if err != nil {
			fmt.Println("unpack error ", err)
			break
		}

		//根据长度 再次读取
		var data []byte
		if msg.GetMsgLen() > 0 {
			//有内容
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.Conn, data); err != nil {
				fmt.Println("read msg data error  ", err)
				break
			}

		}
		msg.SetData(data)

		//将读出来的msg 组装一个request
		//将当前一次性得到的对端客户端请求的数据 封装成一个Request
		req := NewReqeust(c, msg)

		//将req交给worker工作池来处理
		if config.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			go c.MsgHandler.DoMsgHandler(req)
		}

	}

}

/*
 写消息的Goroutine 专门负责给客户端发送消息
 */
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine isStarted]...")
	defer fmt.Println("[Writer Goroutine Stop...]")
	//IO多路复用
	for {
		select {
		case data := <-c.msgChan:
			//有数据需要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error ", err)
				return
			}
		case <-c.writerExitChan:
			//代表reader已经退出了，writer也要退出
			return
		}

	}
}

//启动链接
func (c *Connection) Start() {
	fmt.Println("Conn Start（）  ... id = ", c.ConnID)
	//先进行读业务
	go c.StartReader()

	//进行写业务
	go c.StartWriter()

	//调用创建链接之后  用户自定义的Hook业务
	c.server.CallOnConnStart(c)
}

//停止链接
func (c *Connection) Stop() {
	fmt.Println("c. Stop() ... ConnId = ", c.ConnID)

	//调用销毁链接之前用户自定义的Hook函数
	c.server.CallOnConnStop(c)

	//回收工作
	if c.isClosed == true {
		return
	}

	c.isClosed = true

	//告知writer 链接已经关闭了
	c.writerExitChan <- true

	//关闭原生套接字
	_ = c.Conn.Close()
	//将当前链接从链接管理模块删除
	c.server.GetConnMgr().Remove(c.ConnID)

	//释放channel资源
	close(c.msgChan)
	close(c.writerExitChan)
}

//获取链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取conn的原生socket套接字
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取远程客户端的ip地址
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//发送数据给对方客户端
func (c *Connection) Send(msgId uint32, msgData []byte) error {

	if c.isClosed == true {
		return errors.New("Connection closed ..send Msg ")
	}
	//封装成msg
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, msgData))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return err
	}

	//将要发送的打包好的二进制数发送channel 让writer去写
	c.msgChan <- binaryMsg

	return nil
}

//设置属性
func (c *Connection) SetProperty(key string, value interface{}) {
	//加锁
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	//添加一个连接属性
	c.property[key] = value
}

//获取属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	//加读锁
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	//读取属性
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("not found property :" + key)
	}
}

//删除属性
func (c *Connection) RemoveProperty(key string) {
	//加锁
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	//删除属性
	delete(c.property, key)
}
