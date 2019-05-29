/**
* @Author: Aceld(刘丹冰)
* @Date: 2019/5/25 10:23
* @Mail: danbing.at@gmail.com
*/
package net

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 函数名 Test开头  后面的函数名 自定义
//形参 (t *testing.T)
func TestDataPack(t *testing.T) {
	fmt.Println("test datapack ...")

	/*
	  模拟写一个server
	  收到二进制流 进行解包
	 */
	 // 1 创建listenner
	 listenner, err := net.Listen("tcp", "127.0.0.1:7777")
	 if err != nil {
	 	fmt.Println("server listenner err", err)
	 	return
	 }
	 // 2 AcceptTCP
	 go func() {
	 	for {
	 		conn , err:= listenner.Accept()
	 		if err != nil {
	 			fmt.Println("server accpet err", err)
			}

			//3 读写业务
			go func(conn *net.Conn) {
				//读取客户端的请求
				// ---- 拆包过程 ---
				// |datalen|id|data|
				dp := NewDataPack()
				for {
					//进行第一次从conn读， 把head读出来
					headData := make([]byte, dp.GetHeadLen())//8个字节
					_, err := io.ReadFull(*conn, headData) //直到headData填充满，才会返回，否则阻塞
					if err != nil {
						fmt.Println("read head error")
						break
					}
					//headData ==  > |datalen|id|  （8字节的长度）
					//将headData ---> Message结构体中 填充 len  id两个字段
					msgHead, err := dp.UnPack(headData)
					//msgHead : 已经填充好了 Datalen  id 两个字段，data -->nil
					if err !=nil {
						fmt.Println("server unpack err ", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						//数据区有内容，需要进行第二次读取
						//将msgHead进行向下装换 将IMessage 转换成Message
						msg := msgHead.(*Message)
						//给msg的Data属性开辟 ， 长度就是数据的长度  data|
						msg.Data = make([]byte, msg.GetMsgLen())

						//根据datalen的长度进行第二次read
						_, err := io.ReadFull(*conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data error ", err)
							return
						}
						fmt.Println("---> Recv MsgID = ", msg.Id, " datalen = ", msg.Datalen, "data = ", string(msg.Data))
					}
				}

			}(&conn)
		}
	 }()


	/*
	 模拟写一个client  封包之后再发包
	 */
	 //connection Dail
	 conn, err := net.Dial("tcp", "127.0.0.1:7777")
	 if err != nil {
	 	fmt.Println("client dail err: ", err)
	 	return
	 }
	 //封包
	 //创建dp拆包 封包的工具
	 dp := NewDataPack()

	 //模拟粘包过程发包
	 //封装第一个包
	 msg1 := &Message{
		Id:1,
		Datalen:4,
		Data: []byte{'z','i','n','x'},
	}
	 sendData1, err := dp.Pack(msg1)
	 if err != nil {
	 	fmt.Println("client send data1 error")
	 	return
	 }
	//封装第2个包
	msg2 := &Message{
		Id:2,
		Datalen:5,
		Data: []byte{'h','e','l','l','o'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client send data2 error")
		return
	}

	//将两个包黏在一起
	sendData1 = append(sendData1, sendData2...) //[4][1]zinx[5][2]hello
	//发送
	conn.Write(sendData1)


	//让test不结束
	select{}
}
