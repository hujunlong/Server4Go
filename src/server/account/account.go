package account

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"server/global"
	"server/player"
	"server/protocol"
	"strings"
)

func Login(logininfo_proto *protocol.S2SSystem_LoginInfo) int32 {
	player := new(player.Player)
	data, err := global.Redis.Get("player:" + logininfo_proto.GetName())
	if err == nil { //先根据name查询数据 在验证密码
		buf := bytes.NewBuffer(data)
		dec := gob.NewDecoder(buf)
		dec.Decode(player)
		if strings.EqualFold(logininfo_proto.GetPassworld(), player.Password) {
			return global.SUCCESS
		} else {
			return global.PASSWDERROR
		}

	}
	fmt.Println(player.Info.Name, player.Info.Age, player.Money, player.Password)
	return global.LOGINERROR
}

func Register(registerInfo_proto *protocol.S2SSystem_RegisterInfo) error {
	player := new(player.Player)

	//查询玩家名字是否被占用
	data, _ := global.Redis.Get("player:" + registerInfo_proto.GetName())
	if data != nil {
		return errors.New("same nick")
	}
	//玩家id 增加1 读取现在最大玩家数据
	id, err := global.Redis.Incr("Register:MaxId")
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
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err = enc.Encode(player)
	if err == nil {
		err = global.Redis.Set("player:"+registerInfo_proto.GetName(), buf.Bytes())
	}
	return err
}
