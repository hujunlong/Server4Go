//阵型
package game

import (
	"github.com/golang/protobuf/proto"
	"server/share/protocol"
)
//阵型规则
//7 4 1
//8 5 2
//9 6 3
type FomationStruct struct {
	Pos_id   int32
	Hero_uid int32
	Equips   map[int32]int32 //key=孔位置 value=equip_uid
}

type Fomation struct {
	Hero_fomations map[int32]FomationStruct //key=pos_id
}

func (this *Fomation) Init() {
	this.Hero_fomations = make(map[int32]FomationStruct)
}

func (this *Fomation) OnFomation(pos_id int32, hero_uid int32) { //上阵
	if pos_id < 0 || pos_id > 9 {
		return
	}

	//下阵hero
	this.OffFomation(pos_id)

	//上阵新hero
	if v, ok := this.Hero_fomations[pos_id]; !ok {
		var fomation FomationStruct
		fomation.Pos_id = pos_id
		fomation.Hero_uid = hero_uid
		this.Hero_fomations[pos_id] = fomation
	} else {
		v.Hero_uid = hero_uid
		this.Hero_fomations[pos_id] = v
	}
}

func (this *Fomation) OffFomation(pos_id int32) { //下阵
	if pos_id < 0 || pos_id > 9 {
		return
	}

	if value, ok := this.Hero_fomations[pos_id]; ok {
		value.Hero_uid = 0
		this.Hero_fomations[pos_id] = value
	}

}

func (this *Fomation) ChangePos(pos_id_1 int32, pos_id_2 int32) { //交回位置
	if pos_id_1 < 0 || pos_id_1 > 9 || pos_id_2 < 0 || pos_id_2 > 9 {
		return
	}

	v1, pos1_ok := this.Hero_fomations[pos_id_1]
	v2, pos2_ok := this.Hero_fomations[pos_id_2]

	//交回数据
	if pos1_ok && pos2_ok { //两个位置都有英雄
		buff := v1

		v1.Hero_uid = v2.Hero_uid
		this.Hero_fomations[pos_id_1] = v1

		v2.Hero_uid = buff.Hero_uid
		this.Hero_fomations[pos_id_2] = v2
	}

	if pos1_ok && !pos2_ok { //pos2无英雄
		var fomation FomationStruct
		fomation.Pos_id = pos_id_2
		fomation.Hero_uid = v1.Hero_uid
		this.Hero_fomations[pos_id_2] = fomation
	}
}

func (this *Fomation) AddEquip(pos_id int32, hold_id int32, equip_uid int32) { //装备阵型位置 孔位置 装备uid
	fomation := this.Hero_fomations[pos_id]

	if fomation.Equips == nil {
		fomation.Equips = make(map[int32]int32)
		this.Hero_fomations[pos_id] = fomation
	}
	this.Hero_fomations[pos_id].Equips[hold_id] = equip_uid
}

//英雄上下阵
func (this *Fomation) HerosFormation(is_on bool, pos_id int32, hero_uid int32,player *Player) {
	var result int32 = 2

	if _,ok := player.Heros[hero_uid];ok {
		if is_on { //上下阵
				this.OnFomation(pos_id, hero_uid)
			} else {
				this.OffFomation(pos_id)
			}
		result = 0
	}

	result4C := &protocol.Formation_HerosFormationResult{
		Result: &result,
	}
	encObj, _ := proto.Marshal(result4C)
	SendPackage(*player.conn, 1302, encObj)
}

//英雄交换阵型
func (this *Fomation) ChangeHerosFormation(pos_id_1 int32, pos_id_2 int32,player *Player) {
	var result int32 = 0

	this.ChangePos(pos_id_1, pos_id_2)

	result4C := &protocol.Formation_ChangeHerosFormationResult{
		Result: &result,
	}
	encObj, _ := proto.Marshal(result4C)
	SendPackage(*player.conn, 1303, encObj)
}
