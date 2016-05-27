//关卡相关
package game

import (
	"fmt"

	"server/share/global"
	"server/share/protocol"
	"time"
)

//离线产生物品概率
type OffPrecent struct {
	type_id int32 //事件id 1怪物事件 2开宝箱事件 3装备事件 4玩家事件 5奇遇事件
	Precent int32 //概率
}

//奖励
type GuajiMapStage struct {
	Player_exp           int32           //角色经验
	Player_gold          int32           //战斗金币
	Guaji_time           int32           //挂机时间
	Kill_npc_num         int32           //杀死挂机npc
	Now_Guaji_id         int32           //现在挂机地图id
	Guaji_Map_stage_pass map[int32]int32 //通过的关卡(关卡通关的状态) 状态 (-1 未解锁  0解锁未通关 1 一星级通关 2二星通关 3三星通关)
	Guaji_PK             *GuajiPK        //挂机pk相关
	props                []Prop          //道具列表
	equips               []Equip         //装备列表

}

func (this *GuajiMapStage) Init() {
	this.Guaji_Map_stage_pass = make(map[int32]int32)
	this.Guaji_PK = new(GuajiPK)
	this.Guaji_PK.Init()
}

func (this *GuajiMapStage) SetCurrentStage(id int32) {
	this.Now_Guaji_id = id
	this.Player_exp = 0
	this.Player_gold = 0
	this.Guaji_time = 0
	this.Kill_npc_num = 0
	this.props = nil
	this.equips = nil
}

//切换关卡
func (this *GuajiMapStage) ChangeStage(id int32, role_id int64) bool {

	if _, ok := this.Guaji_Map_stage_pass[id]; ok {
		//添加到当前关卡中
		if this.Now_Guaji_id > 0 {
			global_guaji_players.Exit(this.Now_Guaji_id, role_id)
		}

		this.SetCurrentStage(id)
		return true
	}
	return false
}

//处理离线 宝箱产生物品
func (this *GuajiMapStage) DealGoods(now_guaji_id int32, num float32) (int32, []Prop) {
	var props []Prop

	var gold_total int32 = 0
	var box_total_precent int32 = 0

	guaji_event_box := Json_config.guaji_event_box[now_guaji_id].Item0
	for _, k := range guaji_event_box {
		box_total_precent += k.Per
	}

	for _, k := range guaji_event_box { //1道具 2英雄 3装备 4资源
		goods_num := num * (float32(k.Per) / float32(box_total_precent))
		switch k.ItemType {
		case 1:
			prop_uid := GetUid()
			props = append(props, Prop{k.ItemID, prop_uid, k.Num})
		case 4:
			gold_total += int32(goods_num) * randGoodsNum(k.Num, k.Num) //1.2倍
		}
	}

	Log.Info("gold_total = %d,props = %d", gold_total, props)
	return gold_total, props
}

//处理在线 宝箱产生物品
func (this *GuajiMapStage) DealOnlineGoods(now_guaji_id int32) (int32, []Prop) {
	var props []Prop

	var gold_total int32 = 0
	var box_total_precent int32 = 0

	guaji_event_box := Json_config.guaji_event_box[now_guaji_id].Item0
	for _, k := range guaji_event_box {
		box_total_precent += k.Per
	}

	for _, k := range guaji_event_box { //1道具 2英雄 3装备 4资源

		switch k.ItemType {
		case 1:
			prop_uid := GetUid()
			props = append(props, Prop{k.ItemID, prop_uid, k.Num})
		case 4:
			gold_total += randGoodsNum(k.Num, k.Num*6/5)
		}
	}

	Log.Info("gold_total = %d props = %d", gold_total, props)
	return gold_total, props
}

//离线收益
func (this *GuajiMapStage) OffNotice2CGuaji(last_time int32) (int32, int32, int32) {

	buff_Player_exp := this.Player_exp
	buff_Player_gold := this.Player_gold

	guaji_killboss_con := Json_config.guaji_kill_boss_con[this.Now_Guaji_id].Item0

	var kill_npc_total int32 = 0
	var guaji_time_total int32 = 0
	var gold_total int32 = 0
	var exp_total int32 = 0

	if guaji_killboss_con != nil {
		for _, k := range guaji_killboss_con {
			switch k.Con {
			case 101: //击杀怪物
				kill_npc_total = k.Par
			case 102: //挂机秒
				guaji_time_total = k.Par
			case 103: //金钱
				gold_total = k.Par
			case 104: //exp
				exp_total = k.Par
			}
		}
	}

	fmt.Println(kill_npc_total, guaji_time_total, gold_total, exp_total)
	//需要产生物品的次数
	index := Csv.map_guaji.index_value["106"]
	time_ones := Csv.map_guaji.simple_info_map[this.Now_Guaji_id][index]
	total := (int32(time.Now().Unix()) - last_time) / time_ones

	//产生对事件 怪物事件 开宝箱事件 装备事件 玩家事件 奇遇事件 的概率
	event := Json_config.guaji_event[this.Now_Guaji_id].Item0
	var percent_list []OffPrecent
	var total_precent int32 = 0
	for _, k := range event {
		percent_ := OffPrecent{k.Event_type, k.Per}
		percent_list = append(percent_list, percent_)
		total_precent += k.Per
	}

	//各个事件产生的物品
	for _, v := range percent_list {
		num := (float32(v.Precent) / float32(total_precent)) * float32(total)
		switch v.type_id { //事件id 1怪物事件 2开宝箱事件 3装备事件 4玩家事件 5奇遇事件
		case 1:
			this.Kill_npc_num += int32(num)
			this.Player_exp += randGoodsNum(Json_config.guaji_event_monster[this.Now_Guaji_id].Exp_Min, Json_config.guaji_event_monster[this.Now_Guaji_id].Exp_Max) * int32(num)
		case 2:
			buff_gold_total, buff_props := this.DealGoods(this.Now_Guaji_id, num)
			this.props = append(this.props, buff_props...)
			this.Player_gold += buff_gold_total
		case 3:
		case 4:
		case 5:
		}
	}

	//可以恢复的体力
	index_property := Csv.property.index_value["102"]
	distance_time := Csv.property.simple_info_map[2057][index_property]
	can_add_energy := (int32(time.Now().Unix()) - last_time) / int32(distance_time)
	this.Guaji_time = int32(time.Now().Unix()) - last_time

	//告诉player需要添加的物品
	Log.Info("Player_gold = %d Player_exp = %d can_add_energy = %d props = %d", (this.Player_gold - buff_Player_gold), (this.Player_exp - buff_Player_exp), can_add_energy, this.props)
	return (this.Player_gold - buff_Player_gold), (this.Player_exp - buff_Player_exp), can_add_energy

}

//在线挂机收益
func (this *GuajiMapStage) OnNotice2CGuaji(last_time int32) (int32, int32, int32, int32) {
	//now_guaji_str := strconv.Itoa(int(this.Now_Guaji_id))

	buff_Player_exp := this.Player_exp
	buff_Player_gold := this.Player_gold
	this.props = nil
	this.equips = nil

	//在线挂机
	var npc_id int32 = 0

	//遍历循环产生对事件
	event := Json_config.guaji_event[this.Now_Guaji_id].Item0
	var percent_list_ []int32
	var percent_list_value []int32
	for _, k := range event {
		percent_list_ = append(percent_list_, k.Per)
		percent_list_value = append(percent_list_value, k.Event_type)
	}
	index := getRandomIndex(percent_list_)

	switch percent_list_value[index] {
	case 1: //(1 怪物事件 2开宝箱事件 3装备事件)
		item0 := Json_config.guaji_event_monster[this.Now_Guaji_id].Item0
		var percent_list_item0 []int32
		var percent_list_value_item0 []int32
		for _, k := range item0 {
			percent_list_item0 = append(percent_list_item0, k.Percent)
			percent_list_value_item0 = append(percent_list_value_item0, k.MonModelID)
		}
		index := getRandomIndex(percent_list_item0)
		npc_id = percent_list_value_item0[index]

		buff_exp := randGoodsNum(Json_config.guaji_event_monster[this.Now_Guaji_id].Exp_Min, Json_config.guaji_event_monster[this.Now_Guaji_id].Exp_Max)
		this.Player_exp += buff_exp

	case 2:
		gold, buff_props := this.DealOnlineGoods(this.Now_Guaji_id)
		this.Player_gold += gold
		this.props = append(this.props, buff_props...)
	case 3:
		this.equips = append(this.equips, this.OnlineEquips(this.Now_Guaji_id)...)

	}

	Log.Info("Player_gold = %d buff_Player_gold = %d this.Player_exp = %d buff_Player_exp = %d npc_id =%d %d", this.Player_gold, buff_Player_gold, this.Player_exp, buff_Player_exp, npc_id, int32(index))
	return (this.Player_gold - buff_Player_gold), (this.Player_exp - buff_Player_exp), npc_id, int32(index)
}

//在线产生装备
func (this *GuajiMapStage) OnlineEquips(stage_id int32) []Equip {

	//产生道具id
	item0 := Json_config.guaji_percent_equip[stage_id].Item0
	var percent_list_item0 []int32
	var percent_list_value_item0 []int32
	for _, k := range item0 {
		percent_list_item0 = append(percent_list_item0, k.Per)
		percent_list_value_item0 = append(percent_list_value_item0, k.EquipID)
	}
	index := getRandomIndex(percent_list_item0)
	equip_id := percent_list_value_item0[index]

	//产生品质
	item0_quality := Json_config.guaji_quality_equip[stage_id].Item0
	var percent_list_item0_quality []int32
	var percent_list_value_quality []int32
	for _, k := range item0_quality {
		percent_list_item0_quality = append(percent_list_item0_quality, k.Per)
		percent_list_value_quality = append(percent_list_value_quality, k.Quality)
	}
	index = getRandomIndex(percent_list_item0_quality)
	quality := percent_list_value_quality[index]

	//产生数量
	var number int = 1
	var equips []Equip

	for i := 0; i < number; i++ {
		equips = append(equips, Equip{equip_id, GetUid(), -1, quality, 1, 0, 0})
	}
	Log.Info("%d", equips)
	return equips
}

//挑战boss条件
func (this *GuajiMapStage) GuajiInfoResult(id int32) *protocol.Game_GuajiInfoResult {

	if _, ok := this.Guaji_Map_stage_pass[this.Now_Guaji_id]; !ok {
		return nil
	}

	if this.Guaji_Map_stage_pass[this.Now_Guaji_id] != 0 { //状态 (-1 未解锁  0解锁未通关 1 一星级通关 2二星通关 3三星通关)
		return nil
	}

	guaji_killboss_con := Json_config.guaji_kill_boss_con[this.Now_Guaji_id].Item0
	var conditions []*protocol.Game_Conditions

	for i, _ := range guaji_killboss_con {
		var condition protocol.Game_Conditions
		condition.Type = &guaji_killboss_con[i].Con
		condition.Count = &guaji_killboss_con[i].Par
		if guaji_killboss_con[i].Con == 101 { //怪物
			condition.CurCount = &this.Kill_npc_num
		}

		if guaji_killboss_con[i].Con == 102 { //修炼时间
			condition.CurCount = &this.Guaji_time
		}

		if guaji_killboss_con[i].Con == 103 { //金币
			condition.CurCount = &this.Player_gold
		}

		if guaji_killboss_con[i].Con == 104 { //exp
			condition.CurCount = &this.Player_exp
		}
		conditions = append(conditions, &condition)
	}

	result4C := &protocol.Game_GuajiInfoResult{
		Conditions: conditions,
	}
	return result4C
}

//挑战boss发放奖励
func (this *GuajiMapStage) C2SChallengeResult(state int32, stage_id int32) ([]Prop, []Equip) {
	this.props = nil
	this.equips = nil

	var prop Prop
	var equip Equip

	item0 := Json_config.guaji_reward[this.Now_Guaji_id].Item0
	item1 := Json_config.guaji_reward[this.Now_Guaji_id].Item1

	//产生道具的品质
	var percent_list_ []int32
	for _, v := range item1 {
		percent_list_ = append(percent_list_, v.Per)
	}
	index := getRandomIndex(percent_list_)

	//产生道具跟装备
	for _, v := range item0 {
		if v.ItemType == global.Type_prop {
			prop.Prop_id = v.ItemID
			prop.Prop_uid = GetUid()
			prop.Count = v.Num
			this.props = append(this.props, prop)
		}

		if v.ItemType == global.Type_equip {
			equip.Equip_id = v.ItemID
			equip.Equip_uid = GetUid()
			equip.Pos = -1
			equip.Quality = item1[index].Quality
			this.equips = append(this.equips, equip)
		}
	}

	//过关
	this.Guaji_Map_stage_pass[stage_id] = state

	return this.props, this.equips
}

//快速战斗
func (this *GuajiMapStage) FastWar(id int32) *protocol.Game_FastWarResult {
	reward := new(protocol.Game_Reward)
	result4c := new(protocol.Game_FastWarResult)
	var result int32 = 0 //能否快速战斗 0：可以快速战斗 1：该关卡不能快速战斗 2：快速战斗用完
	if id != this.Now_Guaji_id {
		result = 1
		result4c.Result = &result
		return result4c
	}

	var exp int32 = 100
	var gold int32 = 50
	reward.PlayerExp = &exp
	reward.PlayerGold = &gold
	result4c.Result = &result
	result4c.Reward = reward

	return result4c
}

//获取该关卡玩家列表
func (this *GuajiMapStage) GetGuajiRoleListResult() []*protocol.Game_GuajiRoleInfo {
	for _, v := range global_guaji_players.player_list {
		fmt.Println("global_guaji_players_id", v)
	}

	roles_id := global_guaji_players.player_list[this.Now_Guaji_id]
	var GuajiRoleInfos []*protocol.Game_GuajiRoleInfo

	for _, v := range roles_id {

		player := word.players[v]
		guaji_role_info := new(protocol.Game_GuajiRoleInfo)

		//获取现在的状态
		protected_last_time, _ := player.Guaji_Stage.Guaji_PK.GetGuajiNowInfo()

		guaji_role_info.ProtectTime = &protected_last_time
		guaji_role_info.LastPkNum = &player.Guaji_Stage.Guaji_PK.Last_pk_num
		guaji_role_info.KillNum = &player.Guaji_Stage.Guaji_PK.Kill_num

		//玩家基础信息
		guaji_role_info.RoleId = &player.PlayerId
		guaji_role_info.ProfessionId = &player.ProfessionId
		guaji_role_info.Level = &player.Info.Level
		guaji_role_info.Power = &player.Info.Power
		guaji_role_info.Nick = &player.Info.Nick

		var PkType int32 = 1 //1:能够pk 2:免战牌不能pk 3:受保护不能pk 4:等级不够未开放
		if player.Info.Level < player.Guaji_Stage.Guaji_PK.pk_open_level {
			PkType = 4
		}
		if protected_last_time > 0 {
			PkType = 3
		}
		guaji_role_info.PkType = &PkType

		GuajiRoleInfos = append(GuajiRoleInfos, guaji_role_info)
	}
	Log.Info("GuajiRoleInfos = %d", GuajiRoleInfos)
	return GuajiRoleInfos
}

//能否挑战该玩家
func (this *GuajiMapStage) IsCanPk(my_role_id int64, other_role_id int64) bool {
	var count int32 = 0
	for _, v := range global_guaji_players.player_list[this.Now_Guaji_id] {
		if v == my_role_id {
			if this.Guaji_PK.Last_pk_num > 0 { //剩余PK次数>1
				count += 1
			}
		}

		if v == other_role_id {
			if word.players[other_role_id].Guaji_Stage.Guaji_PK.Protect_time == 0 { //对手是否是保护时间
				count += 1
			}
		}
	}

	if count == 2 {
		return true
	} else {
		return false
	}

}
