package core

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"mmo_game_server/pb"
	"sync"
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
		V:    0,                            //角度没有实现
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

//将玩家的聊天数据发送给所有在线玩家
func (p *Player) SendTalkMsgToAll(content string) {
	/*
	//返回给上线玩家初始的坐标
	message BroadCast{
		int32 Pid=1;
		int32 Tp=2; //Tp: 1 世界聊天, 2 坐标, 3 动作, 4 移动之后坐标信息更新
		oneof Data {
			string Content=3;
			Position P=4;
			int32 ActionData=5;
		}
	}
	*/

	//1 定义一个广播的proto 消息数据类型
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1,
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	//2 获取全部的在线玩家
	players := WorldMgrObj.GetAllplayer()

	//向全部的玩家  进行广播
	for _, player := range players {
		err := player.Sendmsg(200, proto_msg)
		if err != nil {
			fmt.Println(player.Pid, "msg send error:", err)
			continue
		}
	}
}

//得到玩家周边的玩家有哪些
func (p *Player) GetSurroundingPlayers() []*Player {
	pids := WorldMgrObj.AoiMgr.GetSurroundPIDsByPos(p.X, p.Z)
	fmt.Println("Surrounding Plyers  = ", pids)

	players := make([]*Player, 0, )
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}

	return players
}

//将自己的消息同步给周边的玩家
func (p *Player) SyncSurrounding() {
	//获取当前玩家周边的玩家有哪些
	players := p.GetSurroundingPlayers()

	//构建一个广播消息200 循环全部players 分别给player对应的客户端发送200消息
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	//将消息发送个给周边的
	for _, player := range players {
		err := player.Sendmsg(200, proto_msg)
		if err != nil {
			fmt.Println("player send 200 err:", err)
			continue
		}
	}

	//将其他玩家告诉 当前玩家 （让当前玩家 能 看到周边的玩家
	//1 构建202消息 players 的信息， 告知当前玩家  p.send (202...)
	//2 得到周边玩家的集合 message player

	players_proto_msg := make([]*pb.Player, 0)
	for _, player := range players {
		//制作一个message player 消息
		pp := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		players_proto_msg = append(players_proto_msg, pp)
	}

	//创价一个 message syncplayers
	syncPlayers_proto_msg := &pb.SyncPlayers{
		Ps: players_proto_msg[:],
	}

	//将当前的周边全部玩家信息 发送到客户端
	err := p.Sendmsg(202, syncPlayers_proto_msg)
	if err != nil {
		fmt.Println("send msg error : ", err)
		return
	}
}

//更新广播当前玩家的最新位置  发送广播
func (p *Player) UpdatePosition(x, y, z, v float32) {
	//计算以下当前的玩家是否已经跨越格子了？
	//旧的格子ID
	oldGrid := WorldMgrObj.AoiMgr.GetGidByPos(p.X, p.Z)
	//新的格子ID
	newGrid := WorldMgrObj.AoiMgr.GetGidByPos(x, z)

	if oldGrid != newGrid {
		//触发grid格子切换

		//把pid从旧的aoi格子中删除
		WorldMgrObj.AoiMgr.RemovePidFromGrid(int(p.Pid), oldGrid)
		//将pid添加到新的aoi格子中去
		WorldMgrObj.AoiMgr.AddPidToGrid(int(p.Pid), newGrid)

		//视野消失的业务
		p.OnExchangeAoiGrid(oldGrid, newGrid)

	}

	//将最新的玩家坐标 ，更新给当前玩家
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	//组建广播protobuf协议 ，msgid 200 tp -4
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  4,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	//过去当前玩家周边的Aoi 九宫个之内的玩家player
	players := p.GetSurroundingPlayers()
	//依次调用player对象，send方法将200法国去
	for _, player := range players {
		err := player.Sendmsg(200, proto_msg) //每个玩家都会给格子的client客户端发送200 消息
		if err != nil {
			fmt.Println("player sendmsg error :", player.Pid, err)
			continue
		}
	}
}

//切换格子的视野问题
func (p *Player) OnExchangeAoiGrid(oldgrid, newgrid int) {
	//获取旧的九宫格的成员
	oldGrids := WorldMgrObj.AoiMgr.GetSurroundGridsByGid(oldgrid)
	//旧的九宫格成员建立一个哈希表， 用来快速查找
	oldGridsMap := make(map[int]bool, len(oldGrids))
	for _, grid := range oldGrids {
		oldGridsMap[grid.GID] = true
	}

	//获取新的九宫格的成员
	newGrids := WorldMgrObj.AoiMgr.GetSurroundGridsByGid(newgrid)
	//将新的九宫格成员建立一个哈希表，用来快速查询
	newGridsMap := make(map[int]bool, len(newGrids))
	for _, grid := range newGrids {
		newGridsMap[grid.GID] = true
	}

	// ---- > 处理视野消失 <-----
	//构建一个MsgID:201
	offline_msg := &pb.SyncPid{
		Pid: p.Pid,
	}
	//找到在旧的九宫格中出现，但是在新的九宫格中没有出现的格子
	leavingGrids := make([]*Grid, 0)
	for _, grid := range oldGrids {
		if _, ok := newGridsMap[grid.GID]; !ok {
			leavingGrids = append(leavingGrids, grid)
		}
	}
	//获取leavingGrids中的全部玩家，
	for _, grid := range leavingGrids {
		players := WorldMgrObj.GetPlayerByGid(grid.GID)
		//让自己在其他玩家的客户端中消失

		for _, player := range players {
			err := player.Sendmsg(201, offline_msg)
			if err != nil {
				fmt.Println("player.Sendmsg(201, offline_msg) error:", err, player.Pid)
				continue
			}

			//将其他玩家信息 在自己的客户端中消失
			another_offline_msg := &pb.SyncPid{
				Pid: player.Pid,
			}
			err = p.Sendmsg(201, another_offline_msg)
			if err != nil {
				fmt.Println("p.Sendmsg(201, offline_msg) error:", err, player.Pid)
				continue
			}
		}
	}

	// ---> 处理视野出现 <-----
	onlin_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	//找到在新的九宫格内出现的， 但是没有在旧的九宫格出现的格子
	enteringGrids := make([]*Grid, 0)
	for _, grid := range newGrids {
		if _, ok := oldGridsMap[grid.GID]; !ok {
			enteringGrids = append(enteringGrids, grid)
		}
	}

	//获取enteringGrids中的全部玩家
	for _, grid := range enteringGrids {
		//通过gid获取玩家集合
		players := WorldMgrObj.GetPlayerByGid(grid.GID)

		for _, player := range players {
			//让自己出现在其他人的事业中
			err := player.Sendmsg(200, onlin_msg)
			if err != nil {
				fmt.Println("player.Sendmsg(200, onlin_msg) error:", err, player.Pid)
				continue
			}

			//让其他人出现在自己的事业中
			another_online_msg := &pb.BroadCast{
				Pid: player.Pid,
				Tp:  2,
				Data: &pb.BroadCast_P{
					P: &pb.Position{
						X: player.X,
						Y: player.Y,
						Z: player.Z,
						V: player.V,
					},
				},
			}

			err = p.Sendmsg(200, another_online_msg)
			if err != nil {
				fmt.Println("player.Sendmsg(200, onlin_msg) error:", err, player.Pid)
				continue
			}
		}

	}

}

//用户下线通知
func (p *Player) Offline() {
	//1 得到当前玩家的周边玩家有那些， players
	players := p.GetSurroundingPlayers()

	//2 制作一个msgid 201   前端已经定以好的 msg
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}

	//3 给周边的玩家广播一个消息
	for _, player := range players {
		err := player.Sendmsg(201, proto_msg)
		if err != nil {
			fmt.Println("player.Sendmsg(201, proto_msg) error:", err)
			continue
		}
	}

	//将下线的玩家，从世界管理器移除
	WorldMgrObj.RemovePlayerByPid(p.Pid)

	//将下线玩家从 地图AOImanager 中移除
	WorldMgrObj.AoiMgr.RemoteFromGridbyPos(int(p.Pid), p.X, p.Z)
}
