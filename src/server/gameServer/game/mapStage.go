//关卡相关
package game

import (
	"net"
	"server/share/global"
	"server/share/protocol"

	"github.com/golang/protobuf/proto"
)

//奖励
type MapStage struct {
	consumption_energy int32           //消耗体力
	player_exp         int32           //角色经验
	player_gold        int32           //战斗金币
	hero_exp           int32           //战斗英雄经验
	props              []Prop          //道具列表
	equips             []Equip         //装备列表
	Map_stage_pass     map[int32]int32 //通过的关卡(关卡通关的状态) 状态 (0解锁未通关 1 一星级通关 2二星通关 3三星通关)
}

func (this *MapStage) Init() {
	this.Map_stage_pass = make(map[int32]int32)
}

//关卡编号 玩家体力 type_(1战斗关卡 2：扫荡关卡) 扫荡次数
//返回 (0：允许 1：前置地图未完成 2：体力不够)
func (this *MapStage) IsCanThroughMap(stage_id int32, player_consumption_energy int32, type_ int32, sweep_count int32) int32 {
	//是否解锁
	if _, ok := this.Map_stage_pass[stage_id]; !ok {
		return 1
	}

	//体力验证
	consumption_energy := Csv.map_stage[stage_id].Id_109 //消耗体力ID
	if consumption_energy*sweep_count > player_consumption_energy {
		return 2
	}

	if type_ == 2 && this.Map_stage_pass[stage_id] <= 0 { //扫荡
		return 1
	}

	return 0
}

//清除相关数据
func (this *MapStage) Clean() {
	this.consumption_energy = 0
	this.player_exp = 0
	this.player_gold = 0
	this.hero_exp = 0
	this.props = nil
	this.equips = nil
}

func (this *MapStage) Reward(stage_id_int32 int32) { //能否获取奖励
	//清空slice
	this.Clean()
	this.consumption_energy = Csv.map_stage[stage_id_int32].Id_109 //消耗体力

	//调用json
	stage_stdReward := Json_config.stage_std_reward[stage_id_int32].Item0
	stage_randReward := Json_config.stage_rand_reward[stage_id_int32].Item0
	stage_equip_reward := Json_config.stage_equip_reward[stage_id_int32].Item0

	if len(stage_stdReward) < 3 || len(stage_randReward) < 1 || len(stage_equip_reward) < 1 {
		Log.Error("get stage_stdReward.json stage_randReward.json stage_equip_reward.json stage_equip_quality.json error")
		return
	}

	//固定奖励
	this.player_exp = stage_stdReward[0].Num
	this.player_gold = stage_stdReward[1].Num
	this.hero_exp = stage_stdReward[2].Num

	var index int32 = 0

	//动态奖励
	var percent_list []int32
	for _, v := range stage_randReward {
		percent_list = append(percent_list, v.Percent)
	}
	index = GetRandomIndex(percent_list)
	if stage_randReward[index].Group == global.Type_prop { //道具
		prop_id := stage_randReward[index].ItemID
		prop_count := RandNum(stage_randReward[index].Num_min, stage_randReward[index].Num_Max)

		for i := 0; i < int(prop_count); i++ {
			this.props = append(this.props, Prop{prop_id, 0, 1})
		}

	} else if stage_randReward[index].Group == global.Type_equip { //装备
		//equip_id := stage_randReward[index].ItemID
		equip_count := RandNum(stage_randReward[index].Num_min, stage_randReward[index].Num_Max)
		for i := 0; i < int(equip_count); i++ {
			//uid := GetUid()
			//this.equips = append(this.equips, Equip{equip_id, uid, -1, quality, 1, 0, 1})
			//var object Equip
			//object.Create(equip_id)
			//equip_id
		}
	}

	/*
		//装备奖励
		var percent_list_equip []int32
		for _, v := range stage_equip_reward {
			percent_list_equip = append(percent_list_equip, v.Percent)
		}
		index = GetRandomIndex(percent_list_equip)

		for i := 0; i < int(stage_equip_reward[index].Num); i++ {
			uid := GetUid()
			this.equips = append(this.equips, Equip{stage_equip_reward[index].EquipID, uid, -1, quality, 1, 0, 1})
		}

		Log.Info("map stage reward player_exp = %d player_gold = %d hero_exp = %d props = %d equips = %d", this.player_exp, this.player_gold, this.hero_exp, this.props, this.equips)
	*/
}

//上阵英雄每个都加exp type_ 1:竞技场 2：挂机&关卡
func (this *MapStage) AddExpOnFormation(exp int32, player *Player) {
	for _, v := range player.StageFomation.Hero_fomations {
		player.Heros[v.Hero_uid].heroAddExp(exp, player)
	}
}

//道具跟装备(关卡使用 挂机 接口)
func (this *MapStage) MapStagereward(id int32, player *Player) ([]int32, []*protocol.RwardProp) {
	//产生通过奖励
	this.Reward(id)

	//扣除体力
	energy_comsumer := Csv.map_stage[id].Id_109
	player.ModifyEnergy(energy_comsumer)

	//主角加钱 经验
	player.AddRoleExp(this.player_exp)
	player.ModifyGold(this.player_gold)

	//英雄加经验
	this.AddExpOnFormation(this.hero_exp, player)

	//添加装备
	if len(this.equips) > 0 {
		player.Bag_Equip.Adds(this.equips, player.conn)
	}

	//添加道具
	if len(this.props) > 0 {
		player.Bag_Prop.Adds(this.props, player.conn)
	}

	//装备
	equip_uids := []int32{}
	for _, v := range this.equips {
		equip_uids = append(equip_uids, v.Equip_uid)
	}

	//道具
	var rward_props []*protocol.RwardProp
	for i, _ := range this.props {
		reward_ := &protocol.RwardProp{
			PropUid: &this.props[i].Prop_uid,
			Num:     &this.props[i].Count,
		}
		rward_props = append(rward_props, reward_)
	}
	return equip_uids, rward_props
}

//关卡奖励获取
func (this *MapStage) WarMapNoteServerResult(state int32, id int32, player *Player) {
	if state < 1 { //未通关
		return
	}

	//添加通关
	this.Map_stage_pass[id] = state

	//任务系统
	player.Task.TriggerEvent(16, 1, id)

	//并开启下一关
	if _, ok := Csv.map_stage[id]; ok {
		next_id := Csv.map_stage[id].Id_117
		this.Map_stage_pass[next_id] = 0
		this.Notice2CheckPoint(1, 0, next_id, player.conn)
	}

	equips_uids, props := this.MapStagereward(id, player)
	result4C := &protocol.StageBase_WarMapNoteServerResult{
		Reward: &protocol.StageBase_Reward{
			PlayerExp:  &this.player_exp,
			PlayerGold: &this.player_gold,
			HeroExp:    &this.hero_exp,
			EquipUids:  equips_uids,
			PropUids:   props,
		},
	}
	encObj, _ := proto.Marshal(result4C)
	SendPackage(*player.conn, 1102, encObj)
}

//扫荡
func (this *MapStage) SweepMapStageResult(stage_id int32, sweep_count int32, player *Player) {
	result := this.IsCanThroughMap(stage_id, player.Info.Energy, 2, sweep_count)
	if result != 0 {
		result4C := &protocol.StageBase_SweepMapStageResult{
			Result: &result,
		}
		encObj, _ := proto.Marshal(result4C)
		SendPackage(*player.conn, 1103, encObj)
		return
	}

	var Game_Reward_ []*protocol.StageBase_Reward
	var i int32 = 0
	for ; i < sweep_count; i++ {
		equips_uids, props := this.MapStagereward(stage_id, player)
		reward_ := &protocol.StageBase_Reward{
			PlayerExp:  &this.player_exp,
			PlayerGold: &this.player_gold,
			HeroExp:    &this.hero_exp,
			EquipUids:  equips_uids,
			PropUids:   props,
		}
		Game_Reward_ = append(Game_Reward_, reward_)
	}

	result4C := &protocol.StageBase_SweepMapStageResult{
		Result: &result,
		Reward: Game_Reward_,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*player.conn, 1103, encObj)

	//扫荡
	player.Task.TriggerEvent(16, sweep_count, stage_id)
}

//关卡变化
func (this *MapStage) Notice2CheckPoint(type_ int32, state int32, id int32, conn *net.Conn) { //状态 (-1 未通关  0解锁未通关 1 一星级通关 2二星通关 3三星通关)
	result4C := &protocol.NoticeMsg_Notice2CheckPoint{
		LatestCheckpoint: &protocol.Stage{
			Type:    &type_,
			State:   &state,
			StageId: &id,
		},
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*conn, 1206, encObj)
}
