//装备相关
package game

import (
	"fmt"
	//"net"
)

type Equip struct {
	Equip_id         int32                 //装备id
	Equip_uid        int32                 //装备唯一id
	Pos              int32                 //装备阵型位置（-1 表示未装备）
	Level            int32                 //装备所需等级
	Hole_num         int32                 //开孔个数
	Quality          int32                 //装备品质
	Strengthen_level int32                 //强化次数
	Refine_count     int32                 //精炼次数
	Refine_Level     int32                 //当前精炼等级
	Attrs            map[int32][]Attribute //额外属性 group = 0总值 group=1基础  group=2精炼 group=3强化 group=4品质
}

func (this *Equip) Init() {
	this.Attrs = make(map[int32][]Attribute)
}

//获取基础属性
func (this *Equip) GetBase(equip_id int32) {
	//装备等级
	this.Pos = -1
	this.Level = Csv.equip[equip_id].Id_105
	data := Csv.equip[equip_id].Id_107

	var Attrs []Attribute
	var attr Attribute
	attr.Group = 1
	attr.Value = data

	//获取类型
	switch Csv.equip[equip_id].Id_106 {
	case 107: //生命值
		attr.Key = 107
	case 108: //攻击
		attr.Key = 108
	case 109: //法术
		attr.Key = 109
	case 110: //物理防御
		attr.Key = 110
	case 111: //法术防御
		attr.Key = 111
	case 201: //速度值
		attr.Key = 119
	case 202: //暴击率
		attr.Key = 202
	case 203: //基础暴击伤害
		attr.Key = 203
	case 204: //额外暴击伤害
		attr.Key = 204
	case 205: //总暴击伤害
		attr.Key = 205
	case 206: //连击率
		attr.Key = 206
	case 207: //负面状态率
		attr.Key = 207
	default:
	}

	Attrs = append(Attrs, attr)
	this.Attrs[1] = Attrs
}

func (this *Equip) Create(equip_id int32, quality_id int32, player *Player) *Equip { //品质组id
	//初始化
	this.Init()

	this.Equip_id = equip_id
	this.Equip_uid = GetUid()

	//获取基础属性
	this.GetBase(equip_id)

	//获取品
	this.GetEquipDynamic(quality_id)

	player.Task.TriggerEvent(9, 1, 0)
	player.Task.TriggerEvent(10, 1, this.Quality)
	player.Task.TriggerEvent(11, 1, this.Strengthen_level)
	player.Task.TriggerEvent(13, 1, 0)
	return this
}

func (this *Equip) getIndex(slice []float32) int32 {
	var total float32 = 0
	rand_num := float32(RandNum(1, 10000))
	for i, v := range slice {
		if v == 0 {
			continue
		}
		total += v
		if total >= rand_num {
			return (int32(i) + 1)
		}
	}
	return 0
}

//获取品质组编号id
func (this *Equip) GetQualityItem(quality_item int32) int32 { //传入品质组别
	rand_num := RandNum(1, 10000) //随机品质组
	var count float32 = 0         //用来查看掉落在哪个品质组里面

	for key, v := range Csv.equip_quality[quality_item] {
		count += v.Id_106
		if int32(count) > rand_num {
			return int32(key)
		}
	}
	return -1
}

//关卡跟挂机的品质随机
func (this *Equip) GetEquipDynamic(quality_item int32) { //传入品质组别id

	index_key := this.GetQualityItem(quality_item)
	fmt.Println(quality_item, index_key)

	if index_key == -1 {
		return
	}

	this.Quality = Csv.equip_quality[quality_item][index_key].Id_103
	//属性波动差值
	distance := (Csv.equip_quality[quality_item][index_key].Id_105 - Csv.equip_quality[quality_item][index_key].Id_104) / 5

	duanwei := this.getIndex(Csv.equip_quality[quality_item][index_key].Id_duan)
	kong_wei := this.getIndex(Csv.equip_quality[quality_item][index_key].Id_hole)

	dynamic_pre := float32(duanwei) * (float32(distance) / 10000.0)

	//孔数量
	this.Hole_num = kong_wei

	//品质增加属性
	if _, ok := this.Attrs[1]; ok {
		var Attrs []Attribute
		var attr Attribute
		attr.Group = 4

		attr.Key = this.Attrs[1][0].Key
		attr.Value = this.Attrs[1][0].Value * dynamic_pre

		Attrs = append(Attrs, attr)
		this.Attrs[4] = Attrs

		//品质的额外属性
	}
}

//强化
func (this *Equip) Strengthen(player *Player) int32 { //input 背包的道具
	hold_id := Csv.equip[this.Equip_id].Id_104

	//fmt.Println("come here now", len(Csv.equip_qianghua[hold_id]), int(this.Strengthen_level))
	//if len(Csv.equip_qianghua[hold_id]) <= int(this.Strengthen_level) { //强化满级
	//return 1
	//}

	var data Equip_Qianghua_Struct
	for _, v := range Csv.equip_qianghua[hold_id] {
		if v.Id_103 == this.Strengthen_level+1 { //对应强化等级需求
			data = v
			break
		}
	}

	if player.Info.Level < data.Id_110 { //强化等级不够
		return 1
	}

	//铜钱是否足够

	if player.Info.Gold < data.Id_106 {
		return 2
	}

	//检查道具是否足够
	if !player.Bag_Prop.PropIsenough(data.NeedProp) {
		return 1
	}

	//扣除钱 跟 物品
	player.ModifyGold(-data.Id_106)

	for _, v := range data.NeedProp {
		player.Bag_Prop.DeleteItemById(v.Id, v.Num, player.conn)
	}

	//强化等级+1
	this.Strengthen_level += 1
	player.Task.TriggerEvent(11, 1, this.Strengthen_level)
	//增加物品属性
	dynamic_pre := float32(data.Id_107) / 10000
	fmt.Println("dynamic_re", dynamic_pre, data.Id_107)
	if _, ok := this.Attrs[1]; ok {
		var Attrs []Attribute
		var attr Attribute
		attr.Group = 3

		for _, v := range this.Attrs[1] {
			attr.Key = v.Key
			attr.Value = v.Value * dynamic_pre
			break
		}
		Attrs = append(Attrs, attr)

		if data.Id_108 != 0 { //添加暗属性
			attr.Key = data.Id_108
			attr.Value = data.Id_109
		}
		Attrs = append(Attrs, attr)
		this.Attrs[3] = Attrs
	}

	//推送变化
	var equips []Equip
	equips = append(equips, *this)
	player.Bag_Equip.Notice2CEquip(equips, player.conn)

	return 0
}

//精炼
func (this *Equip) EquipRefine(player *Player) int32 {
	var Attrs []Attribute
	var attr Attribute
	attr.Group = 2

	//检查道具是否足够
	if !player.Bag_Prop.PropIsenough(Csv.equip_jinglian[this.Refine_Level].NeedEquip) {
		return 3
	}

	//扣除物品
	for _, v := range Csv.equip_jinglian[this.Refine_Level].NeedEquip {
		player.Bag_Prop.DeleteItemById(v.Id, v.Num, player.conn)
	}

	//速度
	attr.Key = 119
	attr.Value = Csv.equip_jinglian[this.Refine_Level].Id_107
	Attrs = append(Attrs, attr)

	//计算相关属性
	index := GetRandomIndex(Csv.equip_jinglian[this.Refine_Level].Jinglian_Quan)

	var beishu int32 = 1
	switch index {
	case 0: //1倍权值
		beishu = 1
	case 1: //2倍权值
		beishu = 2
	case 2: //5倍权值
		beishu = 5
	case 3: //10倍权值
		beishu = 10
	}

	//改变精炼次数与等级
	if this.Refine_count+beishu >= Csv.equip_jinglian[this.Refine_Level].Id_103 {
		this.Refine_Level += 1
		this.Refine_count = 0
	} else {
		this.Refine_count += beishu
	}

	//暴击率
	attr.Key = 202
	attr.Value = Csv.equip_jinglian[this.Refine_Level].Id_109
	Attrs = append(Attrs, attr)

	//暴伤率
	attr.Key = 204
	attr.Value = Csv.equip_jinglian[this.Refine_Level].Id_111
	Attrs = append(Attrs, attr)

	//连击率
	attr.Key = 206
	attr.Value = Csv.equip_jinglian[this.Refine_Level].Id_113
	Attrs = append(Attrs, attr)

	//负面状态率
	attr.Key = 207
	attr.Value = Csv.equip_jinglian[this.Refine_Level].Id_115
	Attrs = append(Attrs, attr)

	this.Attrs[2] = Attrs

	var equips []Equip
	equips = append(equips, *this)
	fmt.Println("111", (*this).Refine_count, (*this).Refine_Level, (*this).Equip_uid, (*this).Equip_id)

	player.Bag_Equip.Notice2CEquip(equips, player.conn)

	player.Task.TriggerEvent(14, 1, 0)
	return 0
}

func (this *Equip) EquipDecomposeAddGoods(input_data []Data) []Prop {

	back_goods_pre := Csv.property[1030].Id_102 / 10000
	var datas []Data

	for _, v1 := range input_data {
		var is_add bool = false

		for j, v2 := range datas {
			if v1.Id == v2.Id {
				datas[j].Num += v1.Num
				is_add = true
				continue
			}
		}

		if !is_add {
			datas = append(datas, v1)
		}
	}

	//添加到背包
	var props []Prop
	for _, v := range datas {
		var add_num float32 = float32(v.Num) * back_goods_pre
		var prop Prop
		prop.Count = int32(add_num)
		prop.Prop_id = v.Id
		prop.Prop_uid = GetUid()
		props = append(props, prop)
	}
	return props
}
