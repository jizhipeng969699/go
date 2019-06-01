package main

import (
	"fmt"
	"mmo_game_server/apis"
	"mmo_game_server/core"
	"zinx/net"
	"zinx/ziface"
)

//客户端登陆出发的事件业务
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

	//给conn添加一个属性  pid属性
	conn.SetProperty("pid", p.Pid)

	//同步周边的玩家 ， 告知他们当前玩家已经上线，广播当前的玩家的位置信息
	p.SyncSurrounding()

	fmt.Println("------>player Id = ", p.Pid, "Online...", "Player num =", len(core.WorldMgrObj.Players))
}

//客户退出的时候发出的事件业务
func LostFunc(conn ziface.IConnection) {
	//客户端已经关闭

	//得到当前下限的是那个玩家
	pid, err := conn.GetProperty("pid")
	if err != nil {
		fmt.Println("conn .getroperty pid error:", err)
		return
	}

	//通过pid找到用户  player

	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//玩家的下线业务
	player.Offline()
}

func main() {
	s := net.NewServer("mmo_game_server")

	//添加 hook 函数  处理上线业务
	s.AddOnConnStart(OnlineFunc)
	//添加hook函数， 处理下限业务
	s.AddOnConnStop(LostFunc)

	//针对msgID：2 作一个业务 上线聊天
	s.AddRouter(2, &apis.WorldChar{})
	//针对msgid3 做一个移动业务
	s.AddRouter(3, &apis.Move{})

	//启动游戏服务器
	s.Serve()

	return
}
