//扭蛋相关
package game

import (
	"fmt"
	"server/share/protocol"

	"github.com/golang/protobuf/proto"
)

type ChouKaAddr struct {
	Total   int32 //抽卡次数
	Now_num int32 //现在这轮抽卡次数
	IsGet   bool  //这轮是否已经抽取
}

type NiuDan struct {
	ChouKa map[int32]ChouKaAddr
}

func (this *NiuDan) Init() {
	this.ChouKa = make(map[int32]ChouKaAddr)
}

//产生物品
func (this *NiuDan) CreateProp(type_group int32, player *Player) (Prop, int32, int32) { //道具 钱 hero_uid
	var prop Prop
	index := GetRandomIndex(Csv.niudan_reward[type_group].QuanZhi)

	if Csv.niudan_reward[type_group].Info[index].Id_103 == 1 { //道具或钱
		if Csv.niudan_reward[type_group].Info[index].Id_104 == 30002 { //加钱
			return prop, Csv.niudan_reward[type_group].Info[index].Id_105, 0
		} else {
			prop.Prop_id = Csv.niudan_reward[type_group].Info[index].Id_104
			prop.Count = Csv.niudan_reward[type_group].Info[index].Id_105
			prop.Prop_uid = GetUid()
			return prop, 0, 0
		}
	}

	if Csv.niudan_reward[type_group].Info[index].Id_103 == 2 { //英雄
		hero_id := Csv.niudan_reward[type_group].Info[index].Id_104
		hero := new(HeroStruct)
		hero_uid, _ := hero.CreateHero(hero_id, player)
		player.Heros[hero_uid] = hero
		fmt.Println("英雄相关:hero_uid=", hero_uid)
		return prop, 0, hero_uid
	}

	return prop, 0, 0
}

//是否触发保底
func (this *NiuDan) IsTriggerEvent(type_ int32, add_num int32, baodi_num int32) bool { //type_:1地仙 2天仙 add_num：添加次数 baodi_num：保底次数

	if baodi_num == 0 { //保底次数配置为空
		return false
	}

	if v, ok := this.ChouKa[type_]; ok {
		v.Total += add_num
		v.Now_num += add_num

		if v.Now_num >= baodi_num { //最后一次发送保底
			if !v.IsGet { //未领取
				return true
			} else {
				v.IsGet = false
				v.Now_num -= baodi_num
			}
		}

		if v.Now_num < baodi_num { //随机发送
			if !v.IsGet { //未领取
				rand_num := RandNum(1, baodi_num)
				if this.ChouKa[type_].Now_num > rand_num {
					v.IsGet = true
					return true
				}
			}
		}

		this.ChouKa[type_] = v
	}
	return false
}

//背包检查并扣除相关物品
func (this *NiuDan) DeleteGoods(type_ int32, type_group int32, player *Player) bool {
	var is_goods_enough bool = false //背包道具是否足够

	if type_ == 1 {
		is_goods_enough = player.Bag_Prop.PropIsenough(Csv.niudan[type_group].One_Need)

	} else {
		is_goods_enough = player.Bag_Prop.PropIsenough(Csv.niudan[type_group].Ten_Need)
	}

	//背包道具检查
	if !is_goods_enough {
		return false
	}

	//扣除单抽物品
	if type_ == 1 {
		for _, v := range Csv.niudan[type_group].One_Need {
			player.Bag_Prop.DeleteItemById(v.Id, v.Num, player.conn)
		}
	}

	//扣除十连抽
	if type_ == 2 {
		for _, v := range Csv.niudan[type_group].Ten_Need {
			player.Bag_Prop.DeleteItemById(v.Id, v.Num, player.conn)
		}
	}

	return true
}

//单抽奖励
func (this *NiuDan) DanChou(type_group int32, is_trigger bool, baodi_group int32, player *Player) (Prop, int32, int32) { //道具 金币 hero_uid
	player.Task.TriggerEvent(17, 1, 0)

	if is_trigger {
		return this.CreateProp(baodi_group, player)
	} else {
		return this.CreateProp(type_group, player)
	}

}

//十连抽奖励
func (this *NiuDan) TenChou(type_group int32, is_trigger bool, baodi_group int32, ten_will_out int32, player *Player) ([]Prop, int32, []int32, []int32) { //道具 金币

	var last_num int = 10
	var gold_total int32 = 0
	var golds []int32
	var props []Prop
	var hero_uids []int32

	if is_trigger { //产生保底
		last_num -= 1
		prop, gold, hero_uid := this.CreateProp(baodi_group, player)
		if gold > 0 {
			gold_total += gold
			golds = append(golds, gold)
		} else if hero_uid > 0 {
			hero_uids = append(hero_uids, hero_uid)
		} else {
			props = append(props, prop)
		}
		fmt.Println("产生一次保底\n")
	}

	//十连必出
	if ten_will_out > 0 {
		prop, gold, hero_uid := this.CreateProp(ten_will_out, player)
		if gold > 0 {
			gold_total += gold
			golds = append(golds, gold)
		} else if hero_uid > 0 {
			hero_uids = append(hero_uids, hero_uid)
		} else {
			props = append(props, prop)
		}

		last_num -= 1
		fmt.Println("十连必出\n")
	}

	//产生剩余次数
	for i := 0; i < last_num; i++ {
		prop, gold, hero_uid := this.CreateProp(type_group, player)

		if gold > 0 {
			gold_total += gold
			golds = append(golds, gold)
		} else if hero_uid > 0 {
			hero_uids = append(hero_uids, hero_uid)
		} else {
			props = append(props, prop)
		}
		fmt.Println("剩余产生\n")
	}

	player.Task.TriggerEvent(18, 1, 0)
	return props, gold_total, hero_uids, golds
}

//input: type_ 1：单抽 2:10连抽  type_group:地仙阁，类型2，天仙阁
//output: 0:ok 1:等级不足 2:消耗道具不足 3:抽卡地址传入配置表未找到
func (this *NiuDan) NiuDanMsg(type_ int32, type_group int32, player *Player) (int32, []Prop, []int32, []int32) { //错误码 道具 金币 hero_uid

	var props []Prop
	var golds []int32
	var hero_uids []int32
	var total_gold int32 = 0
	//非法数据
	if _, ok := Csv.niudan[type_group]; !ok {
		fmt.Println("11111 type_group:", type_group)
		return 3, props, golds, hero_uids
	}

	if type_ != 1 && type_ != 2 {
		fmt.Println("type_ != 1 || type_ != 2 ", type_)
		return 3, props, golds, hero_uids
	}

	ten_will_out := Csv.niudan[type_group].Id_110 //十连必出
	baodi_num := Csv.niudan[type_group].Id_111    //保底次数
	baodi_group := Csv.niudan[type_group].Id_112  //保底组

	//检查等级是否足够
	if Csv.niudan[type_group].Id_109 > player.Info.Level {
		return 1, props, golds, hero_uids
	}

	//检查物品并扣除
	is_delete_goods := this.DeleteGoods(type_, type_group, player)
	if !is_delete_goods {
		return 2, props, golds, hero_uids
	}

	var is_trigger bool = false //1：是否触发保底

	if type_ == 1 { //单抽
		is_trigger = this.IsTriggerEvent(type_group, 1, baodi_num)
		prop, gold, hero_uid := this.DanChou(type_group, is_trigger, baodi_group, player)
		if gold > 0 {
			golds = append(golds, gold)
			total_gold += gold
		} else if hero_uid > 0 {
			hero_uids = append(hero_uids, hero_uid)
		} else {
			props = append(props, prop)
		}
	}

	if type_ == 2 { //十连抽
		is_trigger = this.IsTriggerEvent(type_group, 10, baodi_num)
		props_, gold, hero_uids_, golds_ := this.TenChou(type_group, is_trigger, baodi_group, ten_will_out, player)
		total_gold += gold
		props = append(props, props_...)
		golds = append(golds, golds_...)
		hero_uids = append(hero_uids, hero_uids_...)
	}

	//背包添加物品
	player.Bag_Prop.Adds(props, player.conn)

	//添加金币
	player.ModifyGold(total_gold)

	//推送获取英雄
	if len(hero_uids) > 0 {
		heros := make(map[int32]*HeroStruct)
		for _, v := range hero_uids {
			heros[v] = player.Heros[v]
		}
		heros_FightingAttr := player.GetHeroStruct(heros)
		result4C := &protocol.NoticeMsg_NoticeGetHeros{
			Heros: heros_FightingAttr,
		}
		encObj, _ := proto.Marshal(result4C)
		SendPackage(*player.conn, 1212, encObj)
	}

	//返回数据
	return 0, props, golds, hero_uids
}
