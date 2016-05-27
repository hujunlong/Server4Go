//游戏世界
package game

import (
	"fmt"
	//"strconv"
	"sync"

	"github.com/game_engine/timer"
)

type World struct {
	worldMutex *sync.RWMutex
	players    map[int64]*Player //世界内玩家指针 map
}

func (this *World) Init() {
	this.worldMutex = new(sync.RWMutex)
	this.players = make(map[int64]*Player)

	//定时器开启
	index_property := Csv.property.index_value["102"]
	distance_time := int(Csv.property.simple_info_map[2057][index_property])
	timer.CreateTimer(1, true, this.TimerDealEnergy)
	fmt.Println(distance_time)
}

func (this *World) SearchPlayer(ID int64) *Player {
	return this.players[ID]
}

func (this *World) EnterWorld(player *Player) {
	this.worldMutex.Lock()
	defer this.worldMutex.Unlock()
	this.players[player.PlayerId] = player
}

func (this *World) ExitWorld(player_id int64) { //退出游戏
	delete(this.players, player_id)
}

func (this *World) TimerDealEnergy() { //在线体力恢复

	for _, v := range this.players {
		//在线推送体力恢复
		if v.Info.EnergyMax > v.Info.Energy {
			v.Info.Energy += 1
			v.Notice2CEnergy()
		}
	}
}
