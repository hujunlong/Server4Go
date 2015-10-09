package account

import (
	"errors"
	"fmt"
	"github.com/game_engine/cache/redis"
	"server/global"
	"server/player"
	"server/protocol"
	"strings"
)

func Login(logininfo_proto *protocol.S2SSystem_LoginInfo) (int32, *player.Player) {
	player := new(player.Player)
	err := redis.Find("player:"+logininfo_proto.GetName(), player)

	if err == nil { //查到改数据
		if strings.EqualFold(logininfo_proto.GetPassworld(), player.Password) {
			return global.LOGINSUCCESS, player
		} else {
			return global.PASSWDERROR, nil
		}
	}
	return global.LOGINERROR, nil
}

func Register(registerInfo_proto *protocol.S2SSystem_RegisterInfo) error {
	player := new(player.Player)

	//查询玩家名字是否被占用
	err := redis.Find("player:"+registerInfo_proto.GetName(), player)
	if err == nil {
		return errors.New("same nick")
	}
	//玩家id 增加1 读取现在最大玩家数据
	id, err := redis.Incr("Register:MaxId")
	if err != nil {
		return err
	}

	player.Info.ID = fmt.Sprintf("%s", id)
	player.Info.Name = registerInfo_proto.GetName()
	player.Info.Age = registerInfo_proto.GetAge()
	player.Info.Sex = registerInfo_proto.GetSex()

	player.Money = 0
	player.Exp = 0
	player.Password = registerInfo_proto.GetPassworld()
	player.Conn = nil

	//存储现在的玩家数据
	redis.Add("player:"+registerInfo_proto.GetName(), player)
	return err
}
