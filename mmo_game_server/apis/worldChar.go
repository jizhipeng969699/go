package apis

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"mmo_game_server/core"
	"mmo_game_server/pb"
	"zinx/net"
	"zinx/ziface"
)

type WorldChar struct {
	net.BaseRouter
}

func (w *WorldChar) Handle(request ziface.IRequest) {
	//1 解析客户端传递进来的protobuf 数据
	proto_msg := &pb.Talk{}
	err := proto.Unmarshal(request.GetMsg().GetMsgData(), proto_msg)
	if err != nil {
		fmt.Println("Talk message unmarshal error:", err)
		return
	}

	//2 通过获取连接属性，得到当前的玩家Id
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("get pid error:", err)
		return
	}

	//3 通过pid 来得到对应player对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//把聊天数据 广播给所有 在线玩家
	player.SendTalkMsgToAll(proto_msg.GetContent())
}
