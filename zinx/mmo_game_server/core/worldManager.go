package core

import "sync"

//当前世界第如的边界参数
const (
	AOI_MIN_X  int = 85
	AOI_MAX_X  int = 410
	AOI_CNTX_X int = 10
	AOI_MIN_Y  int = 75
	AOI_MAX_Y  int = 400
	AOI_CNTY_Y int = 20
)

//当前场景的世界管理模块

type WorldManager struct {
	//当前全部在县的player集合
	Players map[int32]*Player
	//保护player集合的护持所
	pLock sync.RWMutex
	//AOIManager 当前的地图的管理其
	AoiMgr *AOIManager
}

//对外提供一个全局的世界管路模块
var WorldMgrObj *WorldManager

func init() {
	WorldMgrObj = NewWorldManager()
}

//初始化方法
func NewWorldManager() *WorldManager {
	return &WorldManager{
		Players: make(map[int32]*Player),
		AoiMgr:  NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTX_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTY_Y),
	}
}

//添加一个玩家
func (w *WorldManager) Addplayer(player *Player) {
	//加入到世界管理中
	w.pLock.Lock()
	w.Players[player.Pid] = player
	w.pLock.Unlock()

	//加入到世界地图中
	w.AoiMgr.AddToGridByPos(int(player.Pid), player.X, player.Z)
}

//删除一个玩家
func (w *WorldManager) RemovePlayerByPid(pid int32) {
	//从世界管理其删除
	w.pLock.Lock()

	player := w.Players[pid]
	//从世界地图中shanchu
	w.AoiMgr.RemoteFromGridbyPos(int(pid), player.X, player.Z)

	delete(w.Players, pid)

	w.pLock.Unlock()
}

//通过玩家id得到一个玩家对象
func (w *WorldManager) GetPlayerByPid(pid int32) *Player {
	w.pLock.Lock()
	defer w.pLock.Unlock()

	return w.Players[pid]
}

//获取全部的在线玩家集合
func (w *WorldManager) GetAllplayer() []*Player {
	w.pLock.Lock()
	defer w.pLock.Unlock()

	players := make([]*Player, 0)

	//将世界管理其的所有对象全部加入奥players中去
	for _, player := range w.Players {
		players = append(players, player)
	}

	return players

}
