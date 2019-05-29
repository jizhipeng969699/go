/**
* @Author: Aceld(刘丹冰)
* @Date: 2019/5/25 9:35
* @Mail: danbing.at@gmail.com
*/
package ziface

type IMessage interface {
	//getter
	GetMsgId() uint32
	GetMsgLen() uint32
	GetMsgData() []byte

	//setter
	SetMsgId(uint32)
	SetData([]byte)
	SetDatalen(uint32)
}
