package player

import (
	"net"
	"github.com/game_engine/cache/redis"
)

type PlayerInfo struct {
	ID   string
	Name string
	Age  int32
	Sex  int32
}

type Player struct {
	Info     PlayerInfo
	Money    int32
	Exp      int32
	Password string
	Conn     *net.Conn
}

func LoadPlayer(conn *net.Conn, id string) *Player { //读取玩家数据
	playerInfo := new(Player)
	err := redis.Find(id,playerInfo)
	if err == nil {
		return playerInfo
	}
	return nil
}
