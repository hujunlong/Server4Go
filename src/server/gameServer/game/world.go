//游戏世界
package game

import (
	"strconv"
	"sync"
)

type World struct {
	worldMutex *sync.RWMutex
	players    map[string]*Player //世界内玩家指针 map
}

func (this *World) Init() {
	this.worldMutex = new(sync.RWMutex)
	this.players = make(map[string]*Player)
}

func (this *World) SearchPlayer(ID string) *Player {
	return this.players[ID]
}

func (this *World) EnterWorld(player *Player) {
	this.worldMutex.Lock()
	defer this.worldMutex.Unlock()
	id_str := strconv.Itoa(int(player.PlayerId))
	this.players[id_str] = player
}

func (this *World) TimerDealOnlineGuaji() {
	for _, v := range this.players {
		if v.CreateTime > 0 {
			v.Guaji_Stage.OnNotice2CGuaji(v.LastTime)
		}
	}
}
