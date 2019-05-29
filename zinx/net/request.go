/**
* @Author: Aceld(刘丹冰)
* @Date: 2019/5/23 11:57
* @Mail: danbing.at@gmail.com
*/
package net

import "zinx/ziface"

type Request struct {
	//链接信息
	conn ziface.IConnection

	//客户端发送的消息
	msg ziface.IMessage
}

func NewReqeust(conn ziface.IConnection, msg ziface.IMessage) ziface.IRequest {
	req := &Request{
		conn:conn,
		msg:msg,
	}

	return req
}
//得到当前的请求的链接
func(r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

//得到链接的数据
func(r *Request) GetMsg() ziface.IMessage {
	return r.msg
}