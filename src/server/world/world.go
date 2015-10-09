package world

import (
	_ "server/chat"
	"server/player"
)

type WorldMsg struct {
	Players    map[string]*player.Player //世界内玩家指针 map
	Money_rank []*player.Player          //财富排名
}

var World *WorldMsg //全局世界的数据

func init() { //初始化world相关数据
	World = &WorldMsg{make(map[string]*player.Player), make([]*player.Player, 100)}
}
