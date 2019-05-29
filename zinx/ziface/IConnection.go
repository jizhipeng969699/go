/**
* @Author: Aceld(刘丹冰)
* @Date: 2019/5/22 17:18
* @Mail: danbing.at@gmail.com
*/
package ziface

import "net"

/*
	抽象链接层
 */
 type IConnection interface {
 	//启动链接
	Start()

 	//停止链接
 	Stop()

 	//获取链接ID
 	GetConnID() uint32

 	//获取conn的原生socket套接字
 	GetTCPConnection() *net.TCPConn

 	//获取远程客户端的ip地址
 	GetRemoteAddr() net.Addr

 	//发送数据给对方客户端
 	Send(msgId uint32, msgData []byte) error

	//设置属性
	SetProperty(key string,value interface{})
	//获取属性
	GetProperty(key string)(interface{},error)
	//删除属性
	RemoveProperty(key string)
 }

 //业务处理方法 抽象定义
 //type HandleFunc func(request IRequest) error


