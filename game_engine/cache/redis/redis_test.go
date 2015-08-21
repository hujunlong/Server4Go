package redis_test

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"game_engine/redis"
)

type PlayerInfo struct { //玩家的基本信息，正常情况不会改变
	ID        string //玩家的唯一ID
	Name      string //玩家的名字
	Gender    bool   //玩家的性别
	WorldName string //玩家所在世界的名字
	Gold      int
}

func (this *PlayerInfo) Save(redis *redis.Client) error {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(this)
	if err == nil {
		err = redis.Set(this.ID, buf.Bytes())
	}
	return err
}

func LoadPlayer(redis *redis.Client, id string) *PlayerInfo {
	data, err := redis.Get(id)
	if err == nil {
		playerInfo := new(PlayerInfo)
		buf := bytes.NewBuffer(data)
		dec := gob.NewDecoder(buf)
		dec.Decode(playerInfo)
		return playerInfo
	}

	return nil
}

func DeletePlayerInfo(redis *redis.Client, id string) (bool, error) {
	return redis.Del(id)
}

//序列化测试
func (this *PlayerInfo) ZAdd(redis *redis.Client, key string, i float64) error {

	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(this)
	if err == nil {
		_, err := redis.Zadd(key, buf.Bytes(), i)
		return err
	}
	return nil

}

func DecodeItem(data []byte, item *PlayerInfo) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	dec.Decode(item)
}

//排序
func ZRange(redis *redis.Client, key string) {
	//items, _ := redis.Zrangebyscore(key, 0, 100) //取出所有拍卖的物品
	//items, _ := redis.Zrevrange(key, 0, -1)
	items, _ := redis.ZMyGet(key, "0", "100", "2")

	fmt.Println("items.len = ", len(items), items)
	player_info := make([]PlayerInfo, len(items))
	for k, data := range items {
		DecodeItem(data, &player_info[k])
		fmt.Println("player_info[k].Gold = ", player_info[k].Gold, "player_info[k].ID = ", player_info[k].ID)
	}

}

/*
func main() {
	redis := new(redis.Client)
	ZRange(redis, "Market:super")
	//pl := PlayerInfo{"114", "hjl", true, "worldName", 5}
	//result := pl.ZAdd(redis, "Market:super", float64(pl.Gold))
	//if result != nil {
	//fmt.Println("result = ", result)
	//}

	//pl := LoadPlayer(redis, "111")
	//fmt.Println("pl.ID = ", pl.ID, "pl.WorldName = ", pl.WorldName, "pl.Gold = ", pl.Gold, "pl.Exp = ", pl.Exp, "pl.GroupName = ", pl.GroupName)
	//result, _ := DeletePlayerInfo(redis, "112")
	//if result {
	//fmt.Println("delete successfull")
	//}
}
*/
