//关卡相关
package game

import (
	"fmt"
	"server/share/global"
	"strconv"
)

//奖励
type MapStage struct {
	consumption_energy int32            //消耗体力
	player_exp         int32            //角色经验
	player_gold        int32            //战斗金币
	hero_exp           int32            //战斗英雄经验
	props              []Prop           //道具列表
	equips             []Equip          //装备列表
	Map_stage_pass     map[string]int32 //通过的关卡(关卡通关的状态) 状态 (-1 未通关  0解锁未通关 1 一星级通关 2二星通关 3三星通关)
}

func (this *MapStage) Init() {
	fmt.Println("MapStage init")

	this.Map_stage_pass = make(map[string]int32)
}

func (this *MapStage) IsCanThroughMap(stage_id_int32 int32, player_consumption_energy int32, sweep_count int32) int32 { //关卡编号 玩家体力 扫荡次数

	stage_id_str := strconv.Itoa(int(stage_id_int32))
	//前置任务
	pre_stage_id_index := Csv.map_stage.index_value["104"] //前置任务ID
	pre_stage_id_str := Csv.map_stage.simple_info_map[stage_id_str][pre_stage_id_index]
	if _, ok := this.Map_stage_pass[pre_stage_id_str]; ok {
		return 1
	}

	//体力验证
	consumption_energy_index := Csv.map_stage.index_value["109"] //消耗体力ID
	consumption_energy, _ := strconv.Atoi(Csv.map_stage.simple_info_map[stage_id_str][consumption_energy_index])
	if int32(consumption_energy)*sweep_count < player_consumption_energy {
		return 2
	}

	return 0
}

func (this *MapStage) Reward(stage_id_int32 int32) { //能否获取奖励
	//清空slice
	this.props = nil
	this.equips = nil
	stage_id_str := strconv.Itoa(int(stage_id_int32))
	consumption_energy_index := Csv.map_stage.index_value["109"]
	consumption_energy, _ := strconv.Atoi(Csv.map_stage.simple_info_map[stage_id_str][consumption_energy_index]) //消耗体力
	this.consumption_energy = int32(consumption_energy)

	//调用json
	stage_stdReward := Json_config.stage_std_reward[stage_id_str].Item0
	stage_randReward := Json_config.stage_rand_reward[stage_id_str].Item0
	stage_equip_reward := Json_config.stage_equip_reward[stage_id_str].Item0
	stage_equip_quality := Json_config.stage_equip_quality[stage_id_str].Item0

	if len(stage_stdReward) < 3 || len(stage_randReward) < 1 || len(stage_equip_reward) < 1 || len(stage_equip_quality) < 1 {
		Log.Error("get stage_stdReward.json stage_randReward.json stage_equip_reward.json stage_equip_quality.json error")
		return
	}

	//固定奖励
	this.player_exp = stage_stdReward[0].Num
	this.player_gold = stage_stdReward[1].Num
	this.hero_exp = stage_stdReward[2].Num

	var index int = 0
	//装备品质
	var percent_list_quality []int32
	for _, v := range stage_equip_quality {
		percent_list_quality = append(percent_list_quality, v.Percent)
	}
	index = getRandomIndex(percent_list_quality)
	quality := stage_equip_quality[index].Quality

	//动态奖励
	var percent_list []int32
	for _, v := range stage_randReward {
		percent_list = append(percent_list, v.Percent)
	}
	index = getRandomIndex(percent_list)

	if stage_randReward[index].Group == global.Type_prop { //道具
		prop_id := stage_randReward[index].ItemID
		prop_count := randGoodsNum(stage_randReward[index].Num_min, stage_randReward[index].Num_Max)

		for i := 0; i < int(prop_count); i++ {
			uid := GetUid()
			this.props = append(this.props, Prop{prop_id, uid, 1})
		}

	} else if stage_randReward[index].Group == global.Type_equip { //装备
		equip_id := stage_randReward[index].ItemID
		equip_count := randGoodsNum(stage_randReward[index].Num_min, stage_randReward[index].Num_Max)
		for i := 0; i < int(equip_count); i++ {
			uid := GetUid()
			this.equips = append(this.equips, Equip{equip_id, uid, -1, quality, 1, 0, 1})
		}
	}

	//装备奖励
	var percent_list_equip []int32
	for _, v := range stage_equip_reward {
		percent_list_equip = append(percent_list_equip, v.Percent)
	}
	index = getRandomIndex(percent_list_equip)

	for i := 0; i < int(stage_equip_reward[index].Num); i++ {
		uid := GetUid()
		this.equips = append(this.equips, Equip{stage_equip_reward[index].EquipID, uid, -1, quality, 1, 0, 1})
	}

}
