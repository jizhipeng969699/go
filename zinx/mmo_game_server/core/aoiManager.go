package core

import "fmt"

type AOIManager struct {
	//区域的左边界
	MinX int
	//区域的右边边界
	MaxX int
	//x轴方向的格子数量
	CntsX int
	//区域的上边界
	MinY int
	//区域的下边边界
	MaxY int
	//y轴方向的格子数量
	CntsY int
	//整体区域(地图中)拥有哪些格子map:key 格子ID, value:格子对象
	grids map[int]*Grid
}

//初始化一个地图
func NewAOIManager(MinX int, MaxX int, CntsX int, MinY int, MaxY int, CntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  MinX,
		MaxX:  MaxX,
		CntsX: CntsX,
		MinY:  MinY,
		MaxY:  MaxY,
		CntsY: CntsY,
		grids: make(map[int]*Grid),
	}

	//隶属于当前地图的全部格子， 也一并进行初始化
	for y := 0; y < CntsY; y++ {
		for x := 0; x < CntsX; x++ {
			//初始化一个格子
			//格子ID := cntsX * y + x
			gid := y*CntsX + x

			//给aoiManager添加一个格子
			aoiMgr.grids[gid] = NewGrid(gid,
				aoiMgr.MinX+x*aoiMgr.GridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.GridWidth(),
				aoiMgr.MinY+y*aoiMgr.GridHeight(),
				aoiMgr.MinY+(y+1)*aoiMgr.GridHeight())
		}
	}

	return aoiMgr

}

//得到每个格子在x轴方向的宽度
func (a *AOIManager) GridWidth() int {
	return (a.MaxX - a.MinX) / a.CntsX
}

//得到每个格子在Y轴方向的高度
func (a *AOIManager) GridHeight() int {
	return (a.MaxY - a.MinY) / a.CntsY
}

//打印当前的地图信息
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager : \n MinX:%d,MaxX:%d,cntsX:%d, minY:%d, maxY:%d,cntsY:%d, Grids inManager:\n",
		m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)

	//打印全部的格子
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}

	return s
}

//添加一个playerID到一个AOI格子中
func (m *AOIManager) AddPidToGrid(pId, gId int) {
	m.grids[gId].Add(pId, nil)
}

//移除一个playerId  从一个AOI 区域中
func (m *AOIManager) RemovePidFromGrid(pId, gId int) {
	m.grids[gId].Remove(pId)
}

//通过格子Id 获取当前格子的全部Playerid
func (m *AOIManager) GetPidsByGid(gId int) (playerIds []int) {
	return m.grids[gId].GetPlayerIds()
}

//通过一个格子Id得到当前格子的周边九宫个的格子Id集合
func (m *AOIManager) GetSurroundGridsByGid(gId int) (grids []*Grid) {
	//判断Gid 是否在Aoi中
	if _, ok := m.grids[gId]; !ok {
		return
	}

	//将当前中心GID放入九宫个切片中
	grids = append(grids, m.grids[gId])

	//判断GId的左边是否有各自？右边是否有格子
	//通过格子id得到x轴编号   idx = gid%cntsx
	idx := gId % m.CntsX

	//判断GId的左边是否有各自
	if idx > 0 {
		//将左边的格子加入到 grids 切片中
		grids = append(grids, m.grids[gId-1])
	}

	//判断GId的右边是否有各自
	if idx < m.CntsX-1 {
		//将左边的格子加入到 grids 切片中
		grids = append(grids, m.grids[gId+1])
	}

	// ===> 得到一个x轴的格子集合，遍历这个格子集合
	// for ... 依次判断  格子的上面是否有格子？下面是否有格子

	//将X轴全部的Grid ID 放到一个slice中 ，遍历整个slice
	gidsX := make([]int, 0, len(grids))

	for _, v := range grids {
		gidsX = append(gidsX, v.GID)
	}

	for _, gId := range gidsX {
		//10,11,12
		//通过Gid得到当前Gid的Y轴编号
		//idy = gID / cntsX
		idy := gId / m.CntsX

		//上方是否还有格子
		if idy > 0 {
			grids = append(grids, m.grids[gId-m.CntsX])
		}
		//下方是否还有格子
		if idy < m.CntsY-1 {
			grids = append(grids, m.grids[gId+m.CntsX])
		}
	}

	return
}

//通过x,y坐标得到对应的格子Id
func (m *AOIManager) GetGidByPos(x, y float32) int {
	if x < 0 || int(x) >= m.MaxX {
		return -1
	}
	if y < 0 || int(y) >= m.MaxY {
		return -1
	}

	//根据坐标，得到当前玩家所在格子的 Id
	idx := (int(x) - m.MinX) / m.GridWidth()
	idy := (int(y) - m.MinY) / m.GridHeight()

	//gid = idy*m.cntsX + idx
	gid := idy*m.CntsX + idx

	return gid
}

//根据一个坐标，得到周边九宫个之内的全不的玩家id集合
func (m *AOIManager) GetSurroundPIDsByPos(x, y float32) (playerids []int) {
	//通过x，y坐标得到一个格子对应的id
	gid := m.GetGidByPos(x, y)

	//通过格子Id 得到周边九宫个集合
	grids := m.GetSurroundGridsByGid(gid)

	fmt.Println("gids = ", gid)

	//将分别将九宫内的全部的玩家，放在playerds
	for _, grid := range grids {
		playerids = append(playerids, grid.GetPlayerIds()...)
	}

	return
}

//通过坐标 将pid 加入到一个格子中
func (m *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x,y)
	//取出当前的格子
	grid := m.grids[gID]
	//给格子添加玩家
	grid.Add(pID, nil)
}

//通过坐标 把一个player从一个格子中删除
func (m *AOIManager) RemoteFromGridbyPos(pID int , x, y float32) {
	gID := m.GetGidByPos(x,y)

	grid := m.grids[gID]

	grid.Remove(pID)

}