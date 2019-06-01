package core

import (
	"fmt"
	"sync"
)

type Grid struct {
	//格子Id
	GID int
	//格子左边边界坐标
	MinX int
	//格子右边边界的坐标
	MaxX int

	//格子的上边边界的坐标
	MinY int
	//格子的上边边界的坐标
	MaxY int

	//当前格子内的 玩家/物体 成员的ID 集合 map[玩家/无土ID]
	playerIds map[int]interface{}
	//保护当前格子内容的map的锁
	plDLock sync.RWMutex
}

//初始化格子
func NewGrid(GID int, MinX int, MaxX int, MinY int, MaxY int) *Grid {
	return &Grid{
		GID:       GID,
		MinX:      MinX,
		MaxX:      MaxX,
		MinY:      MinY,
		MaxY:      MaxY,
		playerIds: make(map[int]interface{}),
	}
}

//给格子添加一个玩家
func (g *Grid) Add(playerId int, player interface{}) {
	g.plDLock.Lock()

	g.playerIds[playerId] = player
	defer g.plDLock.Unlock()
}

//从格子中删除一个玩家
func (g *Grid) Remove(playerId int) {
	g.plDLock.Lock()
	defer g.plDLock.Unlock()

	delete(g.playerIds, playerId)
}

//得到当前格子所有的玩家Id
func (g *Grid) GetPlayerIds() (playerIds []int) {
	g.plDLock.RLock()
	defer g.plDLock.RUnlock()

	//playerIds:= []int{}

	for k, _ := range g.playerIds {
		playerIds = append(playerIds, k)
	}
	return
}

//调试打印格子信息的方法
func (g *Grid) String() string {
	return fmt.Sprintf("Grid id:%d, minX:%d, maxX:%d, minY:%d, maxY:%d, playIds:%v\n", g.GID, g.MinX, g.MaxX, g.MinX, g.MinY, g.playerIds)
}
