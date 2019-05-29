/**
* @Author: Aceld(刘丹冰)
* @Date: 2019/5/25 9:59
* @Mail: danbing.at@gmail.com
*/
/*
  解决TCP粘包问题的 封包 拆包模块
 */
package net

import (
	"bytes"
	"encoding/binary"
	"zinx/ziface"
)

type DataPack struct {

}

//初始化一个DataPack对象
func NewDataPack() ziface.IDataPack {
	return &DataPack{}
}

//获取二进制包的头部长度  固定返回8
func (dp *DataPack) GetHeadLen() uint32 {
	//Datalen uint32（4字节) + ID uint32（4字节)
	return 8
}

//封包方法  ---- 将 Message  打包成 |datalen|dataID|data|\
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放二进制的字节缓冲
	dataBuffer := bytes.NewBuffer([]byte{})

	//"8[4]"
	//    |
	//   databuff
	//将datalen 写进buffer中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgLen()) ; err != nil {
		return nil, err
	}


	//"8[4]|2[4]"
	//          |
	//       databuff
	//将 dataID写进buffer中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgId()) ; err != nil {
		return nil, err
	}


	//"8[4]|2[4]|12345678"
	//                   |
	//                databuff
	//将 data写进buffer中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgData()) ; err != nil {
		return nil, err
	}

	//返回这个缓冲
	return dataBuffer.Bytes(), nil
}

//拆包方法  ---  将|datalen|dataID|data|   拆解到 Message 结构体中
func (dp *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	//解包的时候 是分2次解压，  第一次读取固定的长度8字节， 第二次是根据len 再次进行read
	msgHead := &Message{} //msgHead.Datalen, msgHead.dataID

	//创建一个 读取二进制数据流的io.Reader
	dataBuff := bytes.NewReader(binaryData)

	// |8[4]|2[4]|"12345678"
	//      |
	//    databuf
	//binaryData
	//将二进制流 先读datalen 放在msg的DataLen属性中
	if err := binary.Read(dataBuff, binary.LittleEndian, &msgHead.Datalen); err != nil {
		return nil, err
	}

	// |8[4]|2[4]|"12345678"
	//           |
	//         databuf

	//将二进制流的  DataID 方在Msg的DataID属性中
	if err := binary.Read(dataBuff, binary.LittleEndian, &msgHead.Id); err != nil {
		return nil, err
	}

	return msgHead, nil
}
