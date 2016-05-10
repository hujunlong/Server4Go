//关卡相关
package game

import (
	"fmt"
	"server/share/global"
	"server/share/protocol"
	"strconv"
	"time"
)

//奖励
type GuajiMapStage struct {
	Player_exp           int32            //角色经验
	Player_gold          int32            //战斗金币
	Guaji_time           int32            //挂机时间
	Kill_npc_num         int32            //杀死挂机npc
	Now_Guaji_id         int32            //现在挂机地图id
	Guaji_Map_stage_pass map[string]int32 //通过的关卡(关卡通关的状态) 状态 (-1 未解锁  0解锁未通关 1 一星级通关 2二星通关 3三星通关)
	props                []Prop           //道具列表
	equips               []Equip          //装备列表
}

func (this *GuajiMapStage) Init() {
	this.Guaji_Map_stage_pass = make(map[string]int32)
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

//离线收益
func (this *GuajiMapStage) OffNotice2CGuaji(last_time int32) (int32, int32, int32) {

	buff_Player_exp := this.Player_exp
	buff_Player_gold := this.Player_gold
	buff_guaji_time := this.Guaji_time

	this.Guaji_time = int32(time.Now().Unix()) - last_time
	result4C := this.GuajiInfoResult(this.Now_Guaji_id)

	var kill_npc_total int32 = 0
	var guaji_time_total int32 = 0
	var gold_total int32 = 0
	var exp_total int32 = 0

	for _, k := range result4C.Conditions {
		if k.GetType() == 101 { //击杀怪物
			kill_npc_total = k.GetCount()
		}

		if k.GetType() == 102 { //挂机秒
			guaji_time_total = k.GetCount()
		}

		if k.GetType() == 103 { //金钱
			gold_total = k.GetCount()
		}

		if k.GetType() == 104 { //exp
			exp_total = k.GetCount()
		}
	}

	now_guaji_str := strconv.Itoa(int(this.Now_Guaji_id))
	index := Csv.map_guaji.index_value["106"]
	time_ones := Str2Int32(Csv.map_guaji.simple_info_map[now_guaji_str][index]) //多久遍历一次
	total := last_time / time_ones                                              //遍历总和

	//遍历循环产生对事件
	event := Json_config.guaji_event[now_guaji_str].Item0
	var percent_list_ []int32
	var percent_list_value []int32
	for _, k := range event {
		percent_list_ = append(percent_list_, k.Per)
		percent_list_value = append(percent_list_value, k.Event_type)
	}

	//产生宝箱事件
	guaji_event_box := Json_config.guaji_event_box[now_guaji_str].Item0
	var guaji_event_box_list []int32
	for _, k := range guaji_event_box {
		guaji_event_box_list = append(guaji_event_box_list, k.Per)
	}

	var i int32
	for i = 0; i < total; i++ {

		index := getRandomIndex(percent_list_)
		if percent_list_value[index] == 1 { //怪物事件
			if this.Kill_npc_num == kill_npc_total {
				continue
			} else {
				this.Kill_npc_num += 1
			}
		}

		if percent_list_value[index] == 2 { //宝箱事件
			index_box := getRandomIndex(guaji_event_box_list)
			if guaji_event_box[index_box].ItemType == 1 { //铜钱
				if this.Player_gold > gold_total {
					this.Player_gold = gold_total
				} else {
					this.Player_gold += randGold(guaji_event_box[index_box].Num)
				}
			} else { //挂机的其他物品（待定）

			}
		}

		//增加主角经验
		if this.Player_exp > exp_total {
			this.Player_exp = exp_total
		} else {
			index_map_guaji := Csv.map_guaji.index_value["110"]
			simple_info_map_value := Csv.map_guaji.simple_info_map[now_guaji_str][index_map_guaji]
			this.Player_exp += randStr2int32(simple_info_map_value)
		}

	}

	//可以恢复的体力
	index_property := Csv.property.index_value["102"]
	distance_time := Str2Int32(Csv.property.simple_info_map[now_guaji_str][index_property])
	can_add_power := (this.Guaji_time - buff_guaji_time) / distance_time

	//挂机时间
	if this.Guaji_time > guaji_time_total {
		this.Guaji_time = guaji_time_total
	}

	//告诉player需要添加的物品
	return (this.Player_gold - buff_Player_gold), (this.Player_exp - buff_Player_exp), can_add_power

}

//在线挂机收益
func (this *GuajiMapStage) OnNotice2CGuaji(last_time int32) (int32, int32) {
	now_guaji_str := strconv.Itoa(int(this.Now_Guaji_id))

	fmt.Println("now_guaji_str:", now_guaji_str)
	buff_Player_exp := this.Player_exp
	buff_Player_gold := this.Player_gold

	this.Guaji_time = int32(time.Now().Unix()) - last_time

	fmt.Println("this.Now_Guaji_id:", this.Now_Guaji_id)

	result4C := this.GuajiInfoResult(this.Now_Guaji_id)

	var kill_npc_total int32 = 0
	var guaji_time_total int32 = 0
	var gold_total int32 = 0
	var exp_total int32 = 0

	for _, k := range result4C.Conditions {
		if k.GetType() == 101 { //击杀怪物
			kill_npc_total = k.GetCount()
		}

		if k.GetType() == 102 { //挂机秒
			guaji_time_total = k.GetCount()
		}

		if k.GetType() == 103 { //金钱
			gold_total = k.GetCount()
		}

		if k.GetType() == 104 { //exp
			exp_total = k.GetCount()
		}
	}

	//遍历循环产生对事件
	event := Json_config.guaji_event[now_guaji_str].Item0
	var percent_list_ []int32
	var percent_list_value []int32
	for _, k := range event {
		percent_list_ = append(percent_list_, k.Per)
		percent_list_value = append(percent_list_value, k.Event_type)
	}

	//产生宝箱事件
	guaji_event_box := Json_config.guaji_event_box[now_guaji_str].Item0
	var guaji_event_box_list []int32
	for _, k := range guaji_event_box {
		guaji_event_box_list = append(guaji_event_box_list, k.Per)
	}

	index := getRandomIndex(percent_list_)
	if percent_list_value[index] == 1 { //怪物事件
		if this.Kill_npc_num < kill_npc_total {
			this.Kill_npc_num += 1
		}
	}

	if percent_list_value[index] == 2 { //宝箱事件
		fmt.Println("guaji_event_box_list:", guaji_event_box_list)
		index_box := getRandomIndex(guaji_event_box_list)
		if guaji_event_box[index_box].ItemType == 1 { //铜钱
			if this.Player_gold > gold_total {
				this.Player_gold = gold_total
			} else {
				fmt.Println("guaji_event_box[index_box].Num:", guaji_event_box[index_box].Num)
				this.Player_gold += randGold(guaji_event_box[index_box].Num)
			}
		} else { //挂机的其他物品（待定）

		}
	}

	//增加主角经验
	if this.Player_exp > exp_total {
		this.Player_exp = exp_total
	} else {
		index_map_guaji := Csv.map_guaji.index_value["110"]
		simple_info_map_value := Csv.map_guaji.simple_info_map[now_guaji_str][index_map_guaji]
		this.Player_exp += randStr2int32(simple_info_map_value)
	}

	//挂机时间
	if this.Guaji_time > guaji_time_total {
		this.Guaji_time = guaji_time_total
	}

	//告诉player需要添加的物品
	return (this.Player_gold - buff_Player_gold), (this.Player_exp - buff_Player_exp)
}

func (this *GuajiMapStage) GuajiInfoResult(id int32) *protocol.Game_GuajiInfoResult {

	now_guaji_str := strconv.Itoa(int(this.Now_Guaji_id))
	if _, ok := this.Guaji_Map_stage_pass[now_guaji_str]; !ok {
		return nil
	}

	if this.Guaji_Map_stage_pass[now_guaji_str] != 0 {
		return nil
	}

	guaji_killboss_con := Json_config.guaji_kill_boss_con[now_guaji_str].Item0
	fmt.Println("now_guaji_str:", now_guaji_str, "guaji_killboss_con:", guaji_killboss_con)
	var conditions []*protocol.Game_Conditions

	for _, k := range guaji_killboss_con {
		condition := new(protocol.Game_Conditions)
		condition.Type = &k.Con
		condition.Count = &k.Par
		if k.Con == 101 { //怪物
			condition.CurCount = &this.Kill_npc_num
		}

		if k.Con == 102 { //修炼时间
			condition.CurCount = &this.Guaji_time
		}

		if k.Con == 103 { //金币
			condition.CurCount = &this.Player_gold
		}

		if k.Con == 104 { //exp
			condition.CurCount = &this.Player_exp
		}

		fmt.Println("type:", k.Con, "count:", k.Par, "condition.CurCount:", condition.CurCount)
		conditions = append(conditions, condition)
	}

	result4C := &protocol.Game_GuajiInfoResult{
		Conditions: conditions,
	}
	return result4C
}

func (this *GuajiMapStage) C2SChallengeResult(state int32, stage_id int32) ([]Prop, []Equip) {
	this.props = nil
	this.equips = nil

	var prop Prop
	var equip Equip

	now_guaji_str := strconv.Itoa(int(this.Now_Guaji_id))
	item0 := Json_config.guaji_reward[now_guaji_str].Item0
	item1 := Json_config.guaji_reward[now_guaji_str].Item1

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
			equip.equip_id = v.ItemID
			equip.equip_uid = GetUid()
			equip.pos = -1
			equip.quality = item1[index].Quality
			this.equips = append(this.equips, equip)
		}
	}

	//过关
	stage_id_str := strconv.Itoa(int(stage_id))
	this.Guaji_Map_stage_pass[stage_id_str] = state

	return this.props, this.equips
}
