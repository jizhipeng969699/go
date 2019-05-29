/**
* @Author: Aceld(刘丹冰)
* @Date: 2019/5/25 9:35
* @Mail: danbing.at@gmail.com
*/
package net

import "zinx/ziface"

type Message struct {
	Id  uint32
	Datalen uint32
	Data []byte
}

//提供一个创建Message的方法
func NewMsgPackage(id uint32, data []byte) ziface.IMessage {
	return &Message{
		Id:id,
		Datalen:uint32(len(data)),
		Data:data,
	}
}

//getter
func (m *Message) GetMsgId() uint32 {
	return m.Id
}
func (m *Message) GetMsgLen() uint32 {
	return m.Datalen
}
func (m *Message) GetMsgData() []byte {
	return m.Data
}

//setter
func (m *Message) SetMsgId(id uint32) {
	m.Id =id
}
func (m *Message) SetData(data []byte) {
	m.Data = data
}
func (m *Message) SetDatalen(len uint32) {
	m.Datalen = len
}
