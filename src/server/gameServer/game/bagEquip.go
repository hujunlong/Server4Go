package game

import (
	"server/share/protocol"

	"net"

	"github.com/golang/protobuf/proto"
)

type BagEquip struct {
	Max       int32           //掉落类型
	OpenCount int32           //开启的个数
	UseCount  int32           //使用个数
	BagEquip  map[int32]Equip //装备key uid
}

func (this *BagEquip) Init() {
	this.BagEquip = make(map[int32]Equip)
}

func (this *BagEquip) Add(equip Equip) int32 { //0：ok 1:背包无可用格子
	if this.UseCount >= this.OpenCount {
		return 1
	}
	id := equip.Equip_uid
	this.BagEquip[id] = equip
	this.UseCount += 1
	return 0
}

func (this *BagEquip) Adds(equps []Equip, conn *net.Conn) (int, int32) { //返回参数1：index     参数2:(0：ok 1:背包无可用格子)
	for i, v := range equps {
		result := this.Add(v)
		if result != 0 {
			return i, result
		} else if result == 1 {
			var result int32 = 1
			result4C := &protocol.Game_Notice2CMsg{
				Msg: &result,
			}
			encObj, _ := proto.Marshal(result4C)
			SendPackage(*conn, 1028, encObj)
		}
	}
	return len(equps), 0
}

func (this *BagEquip) Del(uid int32) int32 { //0:ok 1:未找到道具uid
	if _, ok := this.BagEquip[uid]; !ok {
		return 1
	}
	delete(this.BagEquip, uid)
	this.UseCount -= 1
	return 0
}
