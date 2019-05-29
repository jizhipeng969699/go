package core

import (
	"fmt"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	//初始化AOIManager
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)
	//打印信息
	fmt.Println(aoiMgr)
}

func TestNewGrid2(t *testing.T) {
	//初始化AOIManager
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)

	//求出每个格子周边的九宫格信息

	for Gid, _ := range aoiMgr.grids {
		grids := aoiMgr.GetSurroundGridsByGid(Gid)
		fmt.Println("gid : ", Gid, " grids num = ", len(grids))
		//当前九宫格的ID集合

		// make（切片,len(slice),cap(slice)）
		gIds := make([]int, 0, len(grids))
		for _, grid := range grids {
			gIds = append(gIds, grid.GID)
		}

		fmt.Println("grids ids are ", gIds)
	}
}
