package world

import (
	"game_engine/cache/redis"
	"game_engine/logs"
	"server/player"
)

type Msg struct {
	Name string
	Info string
}

type World struct {
	players    map[string]*player.Player //世界内玩家指针 map
	redis      *redis.Client             //redis对象
	money_rank []*player.Player          //财富排名
	chat       []Msg                     //聊天记录
}
