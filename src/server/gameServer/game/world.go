//游戏世界
package game

import (
	//"fmt"
	//"strconv"
	"sync"

	//"github.com/game_engine/timer"
)

type World struct {
	worldMutex *sync.RWMutex
	players    map[int64]*Player //世界内玩家指针 map
}

func (this *World) Init() {
	this.worldMutex = new(sync.RWMutex)
	this.players = make(map[int64]*Player)
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
