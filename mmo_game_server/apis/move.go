package apis

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"mmo_game_server/core"
	"mmo_game_server/pb"
	"zinx/net"
	"zinx/ziface"
)

type Move struct {
	net.BaseRouter
}

//重写handle函数
func (m *Move) Handle(request ziface.IRequest) {
	//解析客户端法国来的proto协议  msgid 3
	proto_msg := &pb.Position{}
	err := proto.Unmarshal(request.GetMsg().GetMsgData(), proto_msg)
	if err != nil {
		fmt.Println("protobuf unmarshal error:", err)
		return
	}

	//通过连接属性，得到玩家的id
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetConnection().GetProperty error:", err)
		return
	}

	//通过pid找到当前的玩家对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//玩家对象方法，将当前的新坐标位置，发送给全部的周边玩家
	player.UpdatePosition(proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)
}
