package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"math/rand"
	"mmo_game_server/pb"
	"net"
	"os"
	"os/signal"
	"time"
)

type TcpClient struct {
	conn net.Conn
	Pid  int32
	X    float32
	Y    float32
	Z    float32
	V    float32
}

type Message struct {
	Len   uint32
	MsgId uint32
	Data  []byte
}

//永久进行客户端的读写业务
func (this *TcpClient) Start() {
	//循化
	go func() {
		for {
			fmt.Println("deal server msg read and write..")
			//按照框架LTV 先解析头部8字节 ，在得到包体
			headData := make([]byte, 8)
			_, err := io.ReadFull(this.conn, headData)
			if err != nil {
				fmt.Println(err)
				return
			}

			messageHead, err := this.UnPack(headData)
			if err != nil {
				return
			}

			//data
			if messageHead.Len > 0 {
				messageHead.Data = make([]byte, messageHead.Len)
				if _, err := io.ReadFull(this.conn, messageHead.Data); err != nil {
					return
				}
			}

			//得到了一个完整的Message数据包

			//根据不同的MsgID 来处理不同的业务
			this.DoMsg(messageHead)
		}
	}()
}

//拆包方法
func (this *TcpClient) UnPack(headData []byte) (*Message, error) {

	headBufReader := bytes.NewReader(headData)

	head := &Message{}

	//读取Len
	if err := binary.Read(headBufReader, binary.LittleEndian, &head.Len); err != nil {
		return nil, err
	}

	//读取MsgID
	if err := binary.Read(headBufReader, binary.LittleEndian, &head.MsgId); err != nil {

		return nil, err
	}

	return head, nil
}

//封包方法
func (this *TcpClient) Pack(msgId uint32, data []byte) ([]byte, error) {
	outbuff := bytes.NewBuffer([]byte{})

	//写Len
	if err := binary.Write(outbuff, binary.LittleEndian, uint32(len(data))); err != nil {
		fmt.Println("write len error")
		return nil, err
	}

	//写ID
	if err := binary.Write(outbuff, binary.LittleEndian, msgId); err != nil {
		fmt.Println("write id error")
		return nil, err
	}

	//写内容
	if err := binary.Write(outbuff, binary.LittleEndian, data); err != nil {
		fmt.Println("write data error")
		return nil, err
	}

	return outbuff.Bytes(), nil
}

//处理服务器返回的数据业务
func (this *TcpClient) DoMsg(msg *Message) {
	fmt.Println("msgId =", msg.MsgId, ", msgLen = ", msg.Len, ", msg.Data=", msg.Data)
	if msg.MsgId == 1 {
		//服务器回执给客户端 分配ID
		//解析proto协议
		syncpid := &pb.SyncPid{}
		proto.Unmarshal(msg.Data, syncpid)

		//给当前客户端赋值
		this.Pid = syncpid.Pid
	} else if msg.MsgId == 200 {
		//服务器回执给客户端 广播数据

		//解析proto数据
		bdata := &pb.BroadCast{}
		proto.Unmarshal(msg.Data, bdata)

		if bdata.Tp == 2 && bdata.Pid == this.Pid {
			//服务器给自己客户端分配了初始化的位置坐标
			//更新本人的坐标
			this.X = bdata.GetP().X
			this.Y = bdata.GetP().Y
			this.Z = bdata.GetP().Z
			this.V = bdata.GetP().V

			fmt.Printf("Player ID: %d online.. at(%f,%f,%f,%f)\n",this.Pid,this.X,this.Y,this.Z,this.V)


			//玩家就已经上线了
			//触发客户端可以主动请求动作
			go func() {
				for {
					this.AIRobotAction()//自动完成一个AI动作
					time.Sleep(3*time.Second)
				}
			}()

		} else if bdata.Tp == 1 {
			//是一个世界聊天的广播消息
			fmt.Println(fmt.Sprintf("世界聊天: 玩家:%d 说的话是 %s", bdata.Pid, bdata.GetContent()))
		}
	}
}

//简单的AI动作
func (this *TcpClient) AIRobotAction() {
	//聊天 或者 移动

	tp := rand.Intn(2)

	if tp == 0 {
		//自动发送一个聊天信息
		content := fmt.Sprintf("hello 我是 player :%d, 我是你88！？", this.Pid)
		msg := &pb.Talk{
			Content:content,
		}

		//将数据发送给对应的服务端
		this.SendMsg(2, msg)

	} else {
		//自动移动
		x := this.X
		z := this.Z

		randpos := rand.Intn(2)
		if randpos == 0 {
			x -= float32(rand.Intn(10))
			z -= float32(rand.Intn(10))
		} else {
			x += float32(rand.Intn(10))
			z += float32(rand.Intn(10))
		}

		//纠正坐标
		if x > 410 {
			x = 410
		} else if x < 85 {
			x = 85
		}

		if z > 400 {
			z = 400
		} else if z < 75 {
			z = 75
		}

		randv := rand.Intn(2)
		v := this.V
		if randv == 0 {
			v = 25
		} else {
			v = 350
		}

		//打包一个proto结构
		msg := &pb.Position{
			X:x,
			Y:this.Y,
			Z:z,
			V:v,
		}

		fmt.Println("Player Id = ", this.Pid, " walking...")

		this.SendMsg(3, msg)
	}
}

//当前客户端发包的方法
func (this *TcpClient) SendMsg(msgID uint32, data proto.Message) {
	//打包成二进制
	binary_data, err:= proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal error", err)
		return
	}

	//打包LTV
	sendData, err := this.Pack(msgID, binary_data)
	if err == nil {
		//发送给对端
		this.conn.Write(sendData)
	} else {
		fmt.Println(err)
	}
}


func NewTcpClient(ip string, port int) *TcpClient {
	address := fmt.Sprintf("%s:%d", ip, port)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}

	//请求成功
	client := &TcpClient{
		conn: conn,
		Pid:  0,
		X:    0,
		Y:    0,
		Z:    0,
		V:    0,
	}

	return client
}



func main() {

	for i := 0; i < 500; i++ {
		//connection   server

		client := NewTcpClient("127.0.0.1", 8999)

		//connection 读写业务
		client.Start()

		time.Sleep(time.Second)
	}

	//阻塞
	c := make(chan os.Signal, 1)
	//当前对os.Kill 和 os.Interrupt （Ctrl+C）
	signal.Notify(c, os.Kill, os.Interrupt)
	//一旦有os.Kill 和 os.Interrupt信号过来，此时channel就有数据可读，否则就阻塞
	sig := <-c

	fmt.Println("====>recv signal，", sig)

	return
}
