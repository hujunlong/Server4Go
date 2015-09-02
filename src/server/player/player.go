package player

import (
	"bytes"
	"encoding/gob"
	"net"
	"server/global"
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
	data, err := global.Redis.Get(id)
	if err == nil {
		playerInfo := new(Player)
		buf := bytes.NewBuffer(data)
		dec := gob.NewDecoder(buf)
		dec.Decode(playerInfo)
		playerInfo.Conn = conn //内存数据库中没存现在的连接 所以每次load必须获取现在连接
		return playerInfo
	}
	return nil
}
