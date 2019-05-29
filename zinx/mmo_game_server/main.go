package main

import (
	"fmt"
	"zinx/mmo_game_server/core"
	"zinx/net"
	"zinx/ziface"
)

//客户端登陆出发的shijian
func OnlineFunc(conn ziface.IConnection) {
	fmt.Println("player on line")

	//创建一个玩家，将连接和玩家模块绑定
	p := core.NewPlayer(conn)

	//给客户端发送一个msgid ： 1
	p.ReturnPid()

	//给客户端发送一个msgID:200
	p.ReturnPlayerPosition()

	//上线成功了
	//将玩家对象添加到世界管理器中
	core.WorldMgrObj.Addplayer(p)

	fmt.Println("------>player Id = ",p.Pid,"Online...","Player num =",len(core.WorldMgrObj.Players))
}

func main() {
	s := net.NewServer("mmo_game_server")

	s.AddOnConnStart(OnlineFunc)

	s.Serve()
}
