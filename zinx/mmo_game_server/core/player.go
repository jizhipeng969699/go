package core

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"sync"
	"zinx/mmo_game_server/pb"
	"zinx/ziface"
)

type Player struct {
	Pid  int32              //玩家id
	Conn ziface.IConnection //当前玩家的连接，与对应客户端通信
	X    float32            //平面X轴坐标
	Y    float32            //高度
	Z    float32            //平面y轴坐标
	V    float32            //玩家脸部长相的方向
}

//playerid生成器
var PidGen int32 = 1
var Idlock sync.Mutex // 用于保护id生成器的护持所

//初始化玩家
func NewPlayer(conn ziface.IConnection) *Player {
	//分配一个玩家
	Idlock.Lock()
	id := PidGen
	PidGen++
	Idlock.Unlock()

	//创建一个玩家对象
	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)), //随即生成玩家上线所在的x轴坐标
		Y:    0,                            //高度没有实现
		Z:    float32(160 + rand.Intn(10)), //随即生成玩家上线所在的y轴坐标
		V:   0,                            //角度没有实现
	}

	return p

}

//服务端发送玩家初始化id
func (p *Player) ReturnPid() {
	//定义个 msg：Id 1 的派人 proto数据结构
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}

	//将这个消息发送给客户端
	_ = p.Sendmsg(1, proto_msg)
}

//玩家可以和对段发送消息的方法
func (p *Player) Sendmsg(msgId uint32, proto_struct proto.Message) error {
	//要将proto的结构体  转换成 二进制的数据
	data, err := proto.Marshal(proto_struct)
	if err != nil {
		fmt.Println("proto marshal error:", err)
		return err
	}

	//再调用zinx原生的conn.send 方法
	if err := p.Conn.Send(msgId, data); err != nil {
		fmt.Println("conn send error:", err)
		return err
	}
	return nil
}

//服务器给客户端发送一个玩家的初始化位置信息

func (p *Player) ReturnPlayerPosition() {
	//组建msgid 200 的消息

	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2, //2 -坐标信息
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	//将这个消息 发送给客户端
	_ = p.Sendmsg(200, proto_msg)
}
