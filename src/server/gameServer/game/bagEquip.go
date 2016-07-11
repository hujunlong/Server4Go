package game

import (
	"fmt"
	"net"
	"server/share/protocol"

	"github.com/golang/protobuf/proto"
)

type BagEquip struct {
	Max          int32            //掉落类型
	OpenCount    int32            //开启的个数
	UseCount     int32            //使用个数
	BagEquip     map[int32]*Equip //装备key uid
	BagEquipById map[int32]int32  //装备id对应的数量

}

func (this *BagEquip) Init() {
	this.BagEquip = make(map[int32]*Equip)
	this.BagEquipById = make(map[int32]int32)

}

//添加道具不发送推送消息
func (this *BagEquip) Add(equip Equip) {

	id := equip.Equip_uid
	this.BagEquip[id] = &equip
	this.UseCount += 1

	if v, ok := this.BagEquipById[equip.Equip_id]; ok {
		this.BagEquipById[equip.Equip_id] = v + 1
	} else {
		this.BagEquipById[equip.Equip_id] = 1
	}

}

//背包是否足够
func (this *BagEquip) BagIsEnough(equips []Equip) (bool, int32) {
	need_add := int32(len(equips))
	if this.UseCount+need_add >= this.OpenCount {
		return false, (this.OpenCount - this.UseCount)
	}
	return true, (this.OpenCount - this.UseCount)
}

//添加道具并发送推送消息
func (this *BagEquip) AddAndNotice(equip Equip, conn *net.Conn) {

	if this.OpenCount <= this.UseCount {
		this.Notice2CBagWeek(conn)
	}

	id := equip.Equip_uid
	this.BagEquip[id] = &equip
	this.UseCount += 1

	if v, ok := this.BagEquipById[equip.Equip_id]; ok {
		this.BagEquipById[equip.Equip_id] = v + 1
	} else {
		this.BagEquipById[equip.Equip_id] = 1
	}

	var equips []Equip
	equips = append(equips, equip)
	this.Notice2CEquip(equips, conn)
}

//添加道具
func (this *BagEquip) Adds(equips []Equip, conn *net.Conn) {

	is_enough, last_num := this.BagIsEnough(equips)
	if !is_enough {
		this.Notice2CBagWeek(conn)
		equips = equips[:last_num-1]
	}

	for _, v := range equips {
		this.Add(v)
	}

	this.Notice2CEquip(equips, conn)

}

//删除道具
func (this *BagEquip) Del(uid int32, player *Player) int32 { //0:ok 1:未找到道具uid
	if _, ok := this.BagEquip[uid]; !ok {
		return 1
	}

	//key = equip id
	id := this.BagEquip[uid].Equip_id
	if v, ok := this.BagEquipById[id]; ok {
		if v > 1 {
			this.BagEquipById[id] = v - 1
		} else {
			delete(this.BagEquipById, id)
		}
	}

	//key = equip uid
	delete(this.BagEquip, uid)
	this.UseCount -= 1

	//移除推送
	var type_ int32 = 1
	var uids []int32
	uids = append(uids, uid)
	result4C := &protocol.NoticeMsg_Notice2CRemove{
		Type: &type_,
		Uid:  uids,
	}
	encObj, _ := proto.Marshal(result4C)
	SendPackage(*player.conn, 1210, encObj)

	return 0
}

//删除多个装备
func (this *BagEquip) Dels(uids []int32, player *Player) int32 { //0:ok 1:未找到道具uid

	for _, uid := range uids {
		//检查uid
		if _, ok := this.BagEquip[uid]; !ok {
			return 1
		}

		//检查id
		id := this.BagEquip[uid].Equip_id
		if v, ok := this.BagEquipById[id]; ok {
			if v > 1 {
				this.BagEquipById[id] = v - 1
			} else {
				delete(this.BagEquipById, id)
			}
		}

		delete(this.BagEquip, uid)
		this.UseCount -= 1

	}

	//移除推送
	this.Notice2CRemove(uids, player.conn)

	return 0
}

//使用道具
func (this *BagEquip) UseEquip(equip_uid int32, pos int32, type_ int32, player *Player) int32 {
	if v, ok := this.BagEquip[equip_uid]; ok {
		equip_id := v.Equip_id
		hold_id := Csv.equip[equip_id].Id_104 //穿的位置

		//等级判断
		if player.Info.Level < Csv.equip[equip_id].Id_105 {
			return 1
		}

		if type_ == 1 { //穿戴

			v.Pos = pos
			this.BagEquip[equip_uid] = v

			//推送装备数据
			var equips []Equip
			equips = append(equips, *v)
			this.Notice2CEquip(equips, player.conn)

			//阵型
			if player.StageFomation.Hero_fomations[pos].Equips == nil {
				data := player.StageFomation.Hero_fomations[pos]
				data.Equips = make(map[int32]int32)
				player.StageFomation.Hero_fomations[pos] = data
			}
			player.StageFomation.Hero_fomations[pos].Equips[hold_id] = equip_uid
			this.UseCount -= 1

		} else if type_ == 2 { //卸载

			v.Pos = -1
			this.BagEquip[equip_uid] = v
			//推送装备数据
			var equips []Equip
			equips = append(equips, *v)
			this.Notice2CEquip(equips, player.conn)

			//移除孔位置
			delete(player.StageFomation.Hero_fomations[pos].Equips, hold_id)
			this.UseCount += 1
		}

	}

	return 0
}

//装备变化
func (this *BagEquip) DealEquipStructOne(equip Equip) *protocol.EquipStruct {
	equip_info := new(protocol.EquipInfo)
	var Attributes []*protocol.Attribute
	equip_info.Id = &equip.Equip_id
	equip_info.Uid = &equip.Equip_uid
	equip_info.StrengthenCount = &equip.Strengthen_level
	equip_info.Pos = &equip.Pos
	equip_info.Quality = &equip.Quality
	equip_info.RefineCount = &equip.Refine_count
	equip_info.RefineLevel = &equip.Refine_Level
	for _, v1 := range equip.Attrs {
		for _, v2_buff := range v1 {
			v2 := v2_buff
			Attribute := new(protocol.Attribute)
			Attribute.Group = &v2.Group
			Attribute.Key = &v2.Key
			Attribute.Value = &v2.Value
			Attributes = append(Attributes, Attribute)
		}
	}

	equips_struct := new(protocol.EquipStruct)
	equips_struct.Attribute = Attributes
	equips_struct.EquipInfo = equip_info
	return equips_struct
}

func (this *BagEquip) DealEquipStruct(Equips []Equip) []*protocol.EquipStruct {

	var equips_struct []*protocol.EquipStruct
	for _, v := range Equips {
		buff := this.DealEquipStructOne(v)
		equips_struct = append(equips_struct, buff)
	}
	return equips_struct
}

//推送装备变化消息
func (this *BagEquip) Notice2CEquip(equips []Equip, conn *net.Conn) {
	equips_struct := this.DealEquipStruct(equips)

	result4C := &protocol.NoticeMsg_Notice2CEquip{
		Equip: equips_struct,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*conn, 1202, encObj)
}

//推送背包不足
func (this *BagEquip) Notice2CBagWeek(conn *net.Conn) {
	var type_ int32 = 1
	result4C := &protocol.NoticeMsg_Notice2CMsg{
		Msg: &type_,
	}
	encObj, _ := proto.Marshal(result4C)
	SendPackage(*conn, 1208, encObj)
}

//推送移除道具
func (this *BagEquip) Notice2CRemove(uid []int32, conn *net.Conn) {
	var type_ int32 = 1

	result4C := &protocol.NoticeMsg_Notice2CRemove{
		Type: &type_,
		Uid:  uid,
	}
	encObj, _ := proto.Marshal(result4C)
	SendPackage(*conn, 1210, encObj)
}

//分解
func (this *BagEquip) EquipDecompose(equip_uids []int32, player *Player) int32 {

	var props []Prop

	for _, uid := range equip_uids {
		//强化等级
		if _, ok := this.BagEquip[uid]; ok {

			strengthen_level := this.BagEquip[uid].Strengthen_level
			equip_id := this.BagEquip[uid].Equip_id
			csv_pos := Csv.equip[equip_id].Id_104
			quality := this.BagEquip[uid].Quality

			if strengthen_level > 0 {

				if Csv.equip_qianghua_leijia[csv_pos] != nil && Csv.equip_qianghua_leijia[csv_pos][strengthen_level-1] != nil {
					buff_props := this.BagEquip[uid].EquipDecomposeAddGoods(Csv.equip_qianghua_leijia[csv_pos][strengthen_level-1])
					props = append(props, buff_props...)
				}

			}

			//精炼
			refine_Level := this.BagEquip[uid].Refine_Level
			refine_count := this.BagEquip[uid].Refine_count

			if refine_Level > 0 || refine_count > 0 { //精炼
				//添加精炼次数次数的值
				var input_data []Data
				if refine_count > 0 {
					max_level := int32(len(Csv.equip_jinglian_leijia)) - 1
					if refine_count > 0 || (max_level > refine_Level) {
						for _, v := range Csv.equip_jinglian[refine_Level+1].NeedEquip {
							v.Num = v.Num * refine_count
							input_data = append(input_data, v)
						}
					}
				}

				input_data = append(input_data, Csv.equip_jinglian_leijia[refine_Level-1]...)
				buff_props := this.BagEquip[uid].EquipDecomposeAddGoods(input_data)
				fmt.Println("111", Csv.equip_jinglian_leijia[refine_Level-1], buff_props)
				props = append(props, buff_props...)
			}

			//自己物品根据品质发放物品
			if quality > 0 {
				var prop Prop
				prop.Count = Csv.equip_fenjie[quality].Id_104
				prop.Prop_id = Csv.equip_fenjie[quality].Id_103
				prop.Prop_uid = GetUid()
				props = append(props, prop)
			}
		} else {
			return 1
		}
	}

	player.Bag_Prop.Adds(props, player.conn)

	//删除道具
	this.Dels(equip_uids, player)

	return 0
}
