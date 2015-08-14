package player

import (
	"bytes"
	"encoding/gob"
	"game_engine/cache/redis"
	"net"
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

func LoadPlayer(redis *redis.Client, id string) *Player { //读取玩家数据
	data, err := redis.Get(id)
	if err == nil {
		playerInfo := new(Player)
		buf := bytes.NewBuffer(data)
		dec := gob.NewDecoder(buf)
		dec.Decode(playerInfo)
		return playerInfo
	}
	return nil
}
