//角色
package game

import (
	"fmt"
	"net"
	"server/share/global"
	"server/share/protocol"
	"strconv"
	"time"

	"github.com/game_engine/cache/redis"
	"github.com/golang/protobuf/proto"
)

type PlayerInfo struct { //角色的基本信息
	Level     int32  //等级
	Exp       int32  //经验
	Hp        int32  //血量
	Nick      string //角色的名字
	Gold      int32  //金钱
	Vip       int32  //vip
	Energy    int32  //当前体力体力
	EnergyMax int32  //最大体力
	Diamond   int32  //元宝
	Power     int32  //战力
	Signature string //签名
	Option    []bool //设置()
}

type Player struct {
	PlayerId      int64                 //账号id + server_id 注意区别
	ProfessionId  int32                 //职业id
	CreateTime    int32                 //创建时间
	LastTime      int32                 //上次在线时间
	Heros         map[int32]*HeroStruct //英雄列表 key:uid
	Info          PlayerInfo            //角色基础属性
	Stage         *MapStage             //副本关卡
	Bag_Prop      *BagProp              //道具背包
	Bag_Equip     *BagEquip             //装备背包
	Guaji_Stage   *GuajiMapStage        //挂机关卡
	StageFomation *Fomation             //阵型(挂机跟关卡阵型)
	ArenaFomation *Fomation             //阵型 竞技场
	conn          *net.Conn             //角色的网络连接,不需要保存
}

func LoadPlayer(id_str string) *Player { //从数据库读取玩家信息
	player := new(Player)
	err := redis.Find(id_str, player)
	if err == nil {
		return player
	}
	return nil
}

func (this *Player) Init() {
	//英雄相关
	this.Heros = make(map[int32]*HeroStruct)

	//副本地图
	this.Stage = new(MapStage)
	this.Stage.Init()

	//挂机地图
	this.Guaji_Stage = new(GuajiMapStage)
	this.Guaji_Stage.Init()

	//背包
	this.Bag_Prop = new(BagProp)
	this.Bag_Prop.Init()

	this.Bag_Equip = new(BagEquip)
	this.Bag_Equip.Init()
	//阵型(关卡)
	this.StageFomation = new(Fomation)

	//阵型(竞技场)
	this.ArenaFomation = new(Fomation)
}

func (this *Player) SetConn(conn *net.Conn) {
	this.conn = conn
}

func (this *Player) Login(player_id int64, conn *net.Conn) (bool, *Player) {

	player_id_str := strconv.FormatInt(player_id, 10)
	var is_login bool = false
	this = LoadPlayer(player_id_str)

	if this != nil {
		this.SetConn(conn)
		if this.CreateTime <= 0 {
			//如果没创建 直接返回 false
			var is_create bool = false
			result4C := &protocol.Game_RoleInfoResult{
				IsCreate: &is_create,
			}
			encObj, _ := proto.Marshal(result4C)
			SendPackage(*conn, 1002, encObj)
			return is_login, nil
		}
	} else {
		return is_login, nil
	}

	//先判断是否创建角色
	var is_create bool = true
	//英雄列表
	/*
		var HeroStruct_ []*protocol.Game_HeroStruct
		for i := 0; i < len(this.Heros); i++ {
			hero_struct := &protocol.Game_HeroStruct{
				HeroId:  &this.Heros[i].Hero_id,
				HeroUid: &this.Heros[i].Hero_Info.Hero_uid,

				HeroInfo: &protocol.Game_HeroInfo{
					Level:     &this.Heros[i].Hero_Info.Level,
					Hp:        &this.Heros[i].Hero_Info.Hp,
					Power:     &this.Heros[i].Hero_Info.Power,
					StarLevel: &this.Heros[i].Hero_Info.Star_level,
					StepLevel: &this.Heros[i].Hero_Info.Step_level,
				},
			}
			HeroStruct_ = append(HeroStruct_, hero_struct)
			fmt.Println("player.Heros[i].Hero_Info.Hero_id:", this.Heros[i].Hero_Info.Hero_id)
		}
	*/
	//道具

	//装备

	//副本开启列表
	var type_ int32 = 1
	var copy_levels []*protocol.Game_Stage
	for buff_id, buff_v := range this.Stage.Map_stage_pass {
		id := buff_id
		v := buff_v
		copy_level := new(protocol.Game_Stage)
		copy_level.Type = &type_
		copy_level.State = &v
		copy_level.StageId = &id
		copy_levels = append(copy_levels, copy_level)
	}
	//挂机开启列表
	var type_guaji int32 = 2
	var guaji_stages []*protocol.Game_Stage
	for buff_id, buff_v := range this.Guaji_Stage.Guaji_Map_stage_pass {
		v := buff_v
		id := buff_id
		guaji_stage := new(protocol.Game_Stage)
		guaji_stage.Type = &type_guaji
		guaji_stage.State = &v
		stage_id_32 := &id
		guaji_stage.StageId = stage_id_32
		guaji_stages = append(guaji_stages, guaji_stage)
	}

	result4C := &protocol.Game_RoleInfoResult{
		IsCreate: &is_create,
		PlayerInfo: &protocol.Game_PlayerInfo{
			Level:     &this.Info.Level,
			Exp:       &this.Info.Exp,
			Hp:        &this.Info.Hp,
			Energy:    &this.Info.Energy,
			EnergyMax: &this.Info.EnergyMax,
			Vip:       &this.Info.Vip,
			Gold:      &this.Info.Gold,
			Diamond:   &this.Info.Diamond,
			Power:     &this.Info.Power,
			Nick:      &this.Info.Nick,
			Signature: &this.Info.Signature,
			Option:    this.Info.Option,
			RoleId:    &player_id,
		},

		//HeroStruct:   HeroStruct_,
		CopyLevels:   copy_levels,
		HangupLevels: guaji_stages,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*conn, 1002, encObj)
	is_login = true

	//发送离线收益
	this.OffNotice2CGuaji()

	//添加到挂机全局中
	global_guaji_players.Enter(this.Guaji_Stage.Now_Guaji_id, this.PlayerId)

	//添加到world中
	word.EnterWorld(this)

	return is_login, this
}

func (this *Player) RegisterRole(player_id int64, nick string, HeroId int32, conn *net.Conn) int32 { //id role id

	//检查是否已经拥有了key
	key := strconv.Itoa(int(player_id))
	is_exists, _ := redis.Exists(key)
	if !is_exists {
		return global.REGISTERROLEERROR
	}

	if this.CreateTime > 0 {
		return global.ALREADYHAVE
	}

	//检测heroid是否在配置中 读取hero_create表
	create_role_map := Csv.create_role.simple_info_map[HeroId]
	if create_role_map == nil {
		return global.CSVHEROIDEERROR
	}

	//体力
	index := Csv.property.index_value["102"]
	int32_energy := int32(Csv.property.simple_info_map[2056][index])

	//英雄相关
	var hero HeroStruct
	new_hero := hero.CreateHero(HeroId)
	new_hero.SetHeroPos(2, 1) //默认设置在第一个位置
	this.Heros[new_hero.Hero_Info.Hero_uid] = new_hero

	//将英雄添加到阵型中
	this.StageFomation.OnFomation(1, new_hero.Hero_Info.Hero_id, new_hero.Hero_Info.Hero_uid) //添加到1号位置

	this.PlayerId = player_id
	this.ProfessionId = HeroId
	this.CreateTime = int32(time.Now().Unix())
	this.LastTime = this.CreateTime
	this.Info.Level = 1
	this.Info.Exp = 0
	this.Info.Hp = int32(new_hero.Hero_Info.Hp)
	this.Info.Nick = nick
	this.Info.Gold = 0
	this.Info.Vip = 0
	this.Info.Energy = int32_energy
	this.Info.EnergyMax = this.Info.Energy
	this.Info.Diamond = 0
	this.Info.Power = 0
	this.Info.Signature = ""
	this.conn = conn

	//初始化副本关卡
	fmt.Println(this.Stage.Map_stage_pass)
	this.Stage.Map_stage_pass[10101] = 1

	//初始化挂机关卡
	this.Guaji_Stage.Now_Guaji_id = 20101
	this.Guaji_Stage.Guaji_Map_stage_pass[20101] = 0 //解锁未通关
	this.Guaji_Stage.SetCurrentStage(this.Guaji_Stage.Now_Guaji_id)

	//初始化背包开启个数
	this.Bag_Equip.OpenCount = int32(Csv.property.simple_info_map[2018][index])
	this.Bag_Equip.Max = int32(Csv.property.simple_info_map[2019][index])
	this.Bag_Equip.UseCount = 0

	this.Bag_Prop.OpenCount = int32(Csv.property.simple_info_map[2018][index])
	this.Bag_Prop.Max = int32(Csv.property.simple_info_map[2019][index])
	this.Bag_Prop.UseCount = 0

	//写内存数据库
	this.Save()

	Log.Info("HeroId = %d int32_energy = %d this.Bag_Equip.OpenIndex = %d this.Bag_Equip.Max = %d this.Bag_Prop.OpenIndex=%d", HeroId, int32_energy, this.Bag_Equip.OpenCount, this.Bag_Equip.Max, this.Bag_Prop.OpenCount)
	return global.REGISTERROLESUCCESS
}

func (this *Player) Save() {
	err := redis.Modify(strconv.Itoa(int(this.PlayerId)), this)
	if err != nil {
		Log.Error("Save Database error:%s", err.Error())
	}

}

//主角加经验
func (this *Player) addRoleExp(exp int32) {
	exp_index := Csv.role_exp.index_value["102"]
	csv_exp_int32 := Csv.role_exp.simple_info_map[this.Info.Level][exp_index]

	for true {
		if this.Info.Exp+exp > csv_exp_int32 {
			exp -= csv_exp_int32
			this.Info.Level += 1
			if Csv.role_exp.simple_info_map[this.Info.Level] == nil {
				//玩家满级判断
				this.Info.Level -= 1
				this.Info.Exp = Csv.role_exp.simple_info_map[this.Info.Level][exp_index]
				break
			}

			//下一级需要exp
			csv_exp_int32 = Csv.role_exp.simple_info_map[this.Info.Level][exp_index]
		} else {
			this.Info.Exp = this.Info.Exp + exp
			break
		}
	}
	Log.Info("level =%d exp = %d", this.Info.Level, this.Info.Exp)
}

//英雄加经验
func (this *Player) heroAddExp(exp int32, hero_uid int32) {
	//找出需要加经验的hero
	if _, ok := this.Heros[hero_uid]; !ok {
		return
	}

	hero := this.Heros[hero_uid] //用来存储改变的值
	exp_index := Csv.hero_exp.index_value["102"]
	hero_level := this.Heros[hero_uid].Hero_Info.Level
	csv_exp_int32 := Csv.hero_exp.simple_info_map[hero_level][exp_index]

	for true {
		if hero.Hero_Info.Exp+exp > csv_exp_int32 {
			exp -= csv_exp_int32
			hero.Hero_Info.Level += 1
			if Csv.role_exp.simple_info_map[this.Heros[hero_uid].Hero_Info.Level] == nil {
				//玩家满级判断
				hero.Hero_Info.Level -= 1
				hero.Hero_Info.Exp = Csv.role_exp.simple_info_map[hero.Hero_Info.Level][exp_index]
				break
			}

			//下一级需要exp
			csv_exp_int32 = Csv.role_exp.simple_info_map[this.Heros[hero_uid].Hero_Info.Level][exp_index]
		} else {
			hero.Hero_Info.Exp = hero.Hero_Info.Exp + exp
			break
		}
	}
	this.Heros[hero_uid] = hero
	Log.Info("level =%d exp = %d", this.Heros[hero_uid].Hero_Info.Level, this.Heros[hero_uid].Hero_Info.Exp)
}

//上阵英雄每个都加exp type_ 1:竞技场 2：挂机&关卡
func (this *Player) addExpOnFormation(exp int32, type_ int32) {
	if type_ == 1 { //竞技场
		for _, v := range this.ArenaFomation.Hero_fomations {
			this.heroAddExp(exp, v.Hero_uid)
		}
	}

	if type_ == 2 { //挂机
		for _, v := range this.StageFomation.Hero_fomations {
			this.heroAddExp(exp, v.Hero_uid)
		}
	}
}

func (this *Player) dealEquipStruct(equips []Equip) []*protocol.Game_EquipStruct {
	//装备
	var equip_infos []*protocol.Game_EquipInfo
	for i, _ := range equips {
		equip_info := new(protocol.Game_EquipInfo)
		equip_info.Id = &equips[i].Equip_id
		equip_info.Uid = &equips[i].Equip_uid
		equip_info.EquipLevel = &equips[i].Equip_level
		equip_info.StrengthenCount = &equips[i].Strengthen_count
		equip_info.Pos = &equips[i].Pos
		equip_info.Quality = &equips[i].Quality

		equip_infos = append(equip_infos, equip_info)
	}

	var equips_struct []*protocol.Game_EquipStruct
	for i, _ := range equip_infos {
		equip_struct := new(protocol.Game_EquipStruct)
		equip_struct.EquipInfo = equip_infos[i]
		equips_struct = append(equips_struct, equip_struct)
	}
	return equips_struct
}

func (this *Player) dealPropStruct(props_ []Prop) []*protocol.Game_PropStruct {
	//道具
	var props []*protocol.Game_PropStruct
	for i, _ := range props_ {
		prop := new(protocol.Game_PropStruct)
		prop.PropId = &props_[i].Prop_id
		prop.PropUid = &props_[i].Prop_uid
		prop.PropCount = &props_[i].Count
		props = append(props, prop)
	}
	return props
}

//道具跟装备(关卡使用 挂机 接口)
func (this *Player) mapStagereward(id int32) ([]int32, []*protocol.Game_RwardProp) {
	//产生通过奖励
	this.Stage.Reward(id)

	//扣除体力
	energy_index := Csv.map_stage.index_value["109"]
	energy_comsumer := Csv.map_stage.simple_info_map[id][energy_index]
	if this.Info.Energy > energy_comsumer {
		this.Info.Energy -= energy_comsumer
	} else {
		this.Info.Energy = 0
		return nil, nil
	}

	//主角加钱 经验
	this.addRoleExp(this.Stage.player_exp)
	this.Info.Gold += this.Stage.player_gold
	this.addExpOnFormation(this.Stage.hero_exp, 2) //英雄加经验

	//装备
	if len(this.Stage.equips) > 0 {
		index, _ := this.Bag_Equip.Adds(this.Stage.equips, this.conn)
		this.Notice2CEquip(1, this.Stage.equips[:index])
	}

	//道具
	if len(this.Stage.props) > 0 {
		props_ := this.Bag_Prop.Adds(this.Stage.props, this.conn)
		this.Notice2CProp(1, props_)
	}

	//推送相关属性变化
	this.Notice2CRoleInfo()
	this.Notice2CMoney()
	this.Notice2CEnergy()

	if len(this.Stage.equips) > 0 {
		index, _ := this.Bag_Equip.Adds(this.Stage.equips, this.conn)
		this.Notice2CEquip(1, this.Stage.equips[:index])
	}

	//装备
	equip_uids := []int32{}
	for _, v := range this.Stage.equips {
		equip_uids = append(equip_uids, v.Equip_uid)
	}

	//道具
	var rward_props []*protocol.Game_RwardProp
	for i, _ := range this.Stage.props {
		reward_ := &protocol.Game_RwardProp{
			PropUid: &this.Stage.props[i].Prop_uid,
			Num:     &this.Stage.props[i].Count,
		}
		rward_props = append(rward_props, reward_)
	}

	return equip_uids, rward_props
}

//关卡奖励获取
func (this *Player) WarMapNoteServerResult(state int32, id int32) {
	if state < 1 { //未通关
		return
	}

	//添加通关
	this.Stage.Map_stage_pass[id] = state

	//并开启下一关
	index := Csv.map_stage.index_value["117"]
	if Csv.map_stage.simple_info_map[id] != nil {
		next_id := Csv.map_stage.simple_info_map[id][index]
		this.Stage.Map_stage_pass[next_id] = 0
		this.Notice2CheckPoint(1, 0, next_id) //推送
	}

	//拼装数据
	equips_uids, props := this.mapStagereward(id)

	result4C := &protocol.Game_WarMapNoteServerResult{
		Reward: &protocol.Game_Reward{
			PlayerExp:  &this.Stage.player_exp,
			PlayerGold: &this.Stage.player_gold,
			HeroExp:    &this.Stage.hero_exp,
			EquipUids:  equips_uids,
			PropUids:   props,
		},
	}
	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.conn, 1004, encObj)
}

//扫荡
func (this *Player) SweepMapStageResult(stage_id int32, sweep_count int32) {
	result := this.Stage.IsCanThroughMap(stage_id, this.Info.Energy, 2, sweep_count)
	if result != 0 {
		result4C := &protocol.Game_SweepMapStageResult{
			Result: &result,
		}
		encObj, _ := proto.Marshal(result4C)
		SendPackage(*this.conn, 1005, encObj)
		return
	}

	var Game_Reward_ []*protocol.Game_Reward
	var i int32 = 0
	for ; i < sweep_count; i++ {
		equips_uids, props := this.mapStagereward(stage_id)
		reward_ := &protocol.Game_Reward{
			PlayerExp:  &this.Stage.player_exp,
			PlayerGold: &this.Stage.player_gold,
			HeroExp:    &this.Stage.hero_exp,
			EquipUids:  equips_uids,
			PropUids:   props,
		}
		Game_Reward_ = append(Game_Reward_, reward_)
	}

	result4C := &protocol.Game_SweepMapStageResult{
		Result: &result,
		Reward: Game_Reward_,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.conn, 1005, encObj)
	this.Save()
}

//装备变化
func (this *Player) Notice2CEquip(type_ int32, equips []Equip) { //type_1 添加
	equips_struct := this.dealEquipStruct(equips)

	result4C := &protocol.Game_Notice2CEquip{
		Type:  &type_,
		Equip: equips_struct,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.conn, 1008, encObj)
}

//道具变化
func (this *Player) Notice2CProp(type_ int32, props []Prop) { //type_ 1添加 2删除
	props_struct := this.dealPropStruct(props)

	result4C := &protocol.Game_Notice2CProp{
		Type: &type_,
		Prop: props_struct,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.conn, 1006, encObj)
}

//角色变化
func (this *Player) Notice2CRoleInfo() {
	result4C := &protocol.Game_Notice2CRoleInfo{
		Level: &this.Info.Level,
		Exp:   &this.Info.Exp,
		Power: &this.Info.Power,
		Hp:    &this.Info.Hp,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.conn, 1010, encObj)
}

//经济变化
func (this *Player) Notice2CMoney() {
	result4C := &protocol.Game_Notice2CMoney{
		Gold:    &this.Info.Gold,
		Diamond: &this.Info.Diamond,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.conn, 1011, encObj)
}

//体力变化
func (this *Player) Notice2CEnergy() {
	result4C := &protocol.Game_Notice2CEnergy{
		Energy:    &this.Info.Energy,
		EnergyMax: &this.Info.EnergyMax,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.conn, 1012, encObj)
}

//关卡变化
func (this *Player) Notice2CheckPoint(type_ int32, state int32, id int32) { //状态 (-1 未通关  0解锁未通关 1 一星级通关 2二星通关 3三星通关)
	result4C := &protocol.Game_Notice2CheckPoint{
		LatestCheckpoint: &protocol.Game_Stage{
			Type:    &type_,
			State:   &state,
			StageId: &id,
		},
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.conn, 1013, encObj)
}

//玩家退出
func (this *Player) ExitGame() {
	this.LastTime = int32(time.Now().Unix())
	this.Save()
	word.ExitWorld(this.PlayerId)
	global_guaji_players.Exit(this.Guaji_Stage.Now_Guaji_id, this.PlayerId)
	redis.Modify("uid", Global_Uid)
}

//离线收益
func (this *Player) OffNotice2CGuaji() {

	gold, exp, can_add_energy := this.Guaji_Stage.OffNotice2CGuaji(this.LastTime)
	this.Info.Gold += gold
	this.Info.Exp += exp
	this.Info.Energy += can_add_energy

	if this.Info.Energy > this.Info.EnergyMax {
		this.Info.Energy = this.Info.EnergyMax
	}

	this.addRoleExp(exp)
	this.Notice2CRoleInfo()
	this.Notice2CMoney()
	this.Notice2CEnergy()

	if len(this.Guaji_Stage.equips) > 0 {
		index, _ := this.Bag_Equip.Adds(this.Guaji_Stage.equips, this.conn)
		this.Notice2CEquip(1, this.Stage.equips[:index])
	}

	if len(this.Guaji_Stage.props) > 0 {
		props_ := this.Bag_Prop.Adds(this.Guaji_Stage.props, this.conn)
		this.Notice2CProp(1, props_)
		fmt.Println("66666666", props_)
	}

	result4C := &protocol.Game_OffNotice2CGuaji{
		PointId:    &this.Guaji_Stage.Now_Guaji_id,
		Gold:       &this.Guaji_Stage.Player_gold,
		Exp:        &this.Guaji_Stage.Player_exp,
		GuajiTime:  &this.Guaji_Stage.Guaji_time,
		KillNpcNum: &this.Guaji_Stage.Kill_npc_num,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.conn, 1014, encObj)

}

//挑战boss奖励
func (this *Player) C2SChallengeResult(state int32, stage_id int32) {
	props, equips := this.Guaji_Stage.C2SChallengeResult(state, stage_id)

	var equips_list []int32
	for _, v := range equips {
		equips_list = append(equips_list, v.Equip_uid)
	}

	var props_list []*protocol.Game_RwardProp
	for _, v := range props {
		var props_ protocol.Game_RwardProp
		props_.PropUid = &v.Prop_uid
		props_.Num = &v.Count
		props_list = append(props_list, &props_)
	}

	result4C := &protocol.Game_C2SChallengeResult{
		PropUids:  props_list,
		EquipUids: equips_list,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.conn, 1018, encObj)

	if len(equips) > 0 {
		index, _ := this.Bag_Equip.Adds(equips, this.conn)
		this.Notice2CEquip(1, this.Stage.equips[:index])
	}
	if len(props) > 0 {
		props_ := this.Bag_Prop.Adds(props, this.conn)
		this.Notice2CProp(1, props_)
	}

	//推送下个挂机关卡
	next_stage_id_index := Csv.map_guaji.index_value["102"]
	next_stage_id_int32 := Csv.map_guaji.simple_info_map[stage_id][next_stage_id_index]
	this.Notice2CheckPoint(2, 0, next_stage_id_int32)

	this.Guaji_Stage.Guaji_Map_stage_pass[next_stage_id_int32] = 0
	this.Guaji_Stage.ChangeStage(next_stage_id_int32, this.PlayerId)
	this.Save()
}

//在线挂机
func (this *Player) OnNotice2CGuaji() {
	gold, exp, npc_id, guiji_type := this.Guaji_Stage.OnNotice2CGuaji(this.LastTime)

	this.Info.Gold += gold
	this.Info.Exp += exp
	this.addRoleExp(exp)

	Log.Info("%d %d %d %d", gold, exp, this.Info.Gold, this.Info.Exp)
	//推送
	this.Notice2CRoleInfo()
	this.Notice2CMoney()
	if len(this.Guaji_Stage.equips) > 0 {
		index, _ := this.Bag_Equip.Adds(this.Guaji_Stage.equips, this.conn)
		fmt.Println("3333333333", index, len(this.Guaji_Stage.equips), this.Guaji_Stage.equips[:index])
		this.Notice2CEquip(1, this.Guaji_Stage.equips[:index])
	}

	if len(this.Guaji_Stage.props) > 0 {
		props_ := this.Bag_Prop.Adds(this.Guaji_Stage.props, this.conn)
		fmt.Println("555555555", props_)
		this.Notice2CProp(1, props_)
	}

	//发送在线挂机
	var Equip_Uids []int32
	for _, v := range this.Guaji_Stage.equips {
		Equip_Uids = append(Equip_Uids, v.Equip_uid)
	}

	var Prop_Uids []*protocol.Game_RwardProp
	for _, v := range this.Guaji_Stage.props {
		prop_uid := new(protocol.Game_RwardProp)
		prop_uid.PropUid = &v.Prop_uid
		prop_uid.Num = &v.Count
		Prop_Uids = append(Prop_Uids, prop_uid)
	}

	result4C := &protocol.Game_OnNotice2CGuaji{
		GuajiType: &guiji_type,
		NpcId:     &npc_id,
		Gold:      &gold,
		Exp:       &exp,
		EquipUids: Equip_Uids,
		PropUids:  Prop_Uids,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.conn, 1015, encObj)
	this.Save()
}

//获取指定玩家阵型
func (this *Player) GetFomation(role_id int64, type_ int32) {
	var result int32 = 0
	var formations []*protocol.Game_FormationStruct
	result4C := new(protocol.Game_GetGuajiRoleFormationResult)
	if value, ok := word.players[role_id]; ok {
		switch type_ { //(1竞技场，2副本&挂机)
		case 1:
		case 2:
			fomations := value.StageFomation.Hero_fomations
			proto_formation := new(protocol.Game_FormationStruct)
			for _, v := range fomations {
				proto_formation.Pos = &v.Pos_id
				proto_formation.HeroUid = &v.Hero_uid
				proto_formation.HeroId = &v.Hero_id
				formations = append(formations, proto_formation)
			}

		default:
		}

	} else {
		result = 1
		result4C.Result = &result
	}

	result4C.Type = &type_
	result4C.Result = &result
	result4C.Formations = formations
	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.conn, 1023, encObj)
	Log.Info("role_id = %d type_ = %d", role_id, type_)
}

//英雄上下阵
func (this *Player) HerosFormation(type_ int32, is_on bool, pos_id int32, hero_uid int32) {
	var result int32 = 2

	switch type_ { //1竞技场，2副本&挂机
	case 1:
		for _, v := range this.Heros {
			if v.Hero_Info.Hero_uid == hero_uid {
				if is_on { //上下阵
					this.ArenaFomation.OnFomation(pos_id, v.Hero_Info.Hero_id, hero_uid)
				} else {
					this.ArenaFomation.OffFomation(pos_id)
				}
				result = 0
				break
			}
		}
	case 2:
		for _, v := range this.Heros {
			if v.Hero_Info.Hero_uid == hero_uid {
				if is_on { //上下阵
					this.StageFomation.OnFomation(pos_id, v.Hero_Info.Hero_id, hero_uid)
				} else {
					this.StageFomation.OffFomation(pos_id)
				}
				result = 0
				break
			}
		}
	default:
	}
	result4C := &protocol.Game_HerosFormationResult{
		Result: &result,
	}
	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.conn, 1024, encObj)
	this.Save()
}

//英雄交换阵型
func (this *Player) ChangeHerosFormation(pos_id_1 int32, pos_id_2 int32) {
	var result int32 = 0
	for i, v := range this.StageFomation.Hero_fomations {
		if pos_id_1 == v.Pos_id {
			this.StageFomation.Hero_fomations[i].Pos_id = pos_id_2
		}

		if pos_id_2 == v.Pos_id {
			this.StageFomation.Hero_fomations[i].Pos_id = pos_id_1
		}
	}
	result4C := &protocol.Game_ChangeHerosFormationResult{
		Result: &result,
	}
	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.conn, 1025, encObj)
	this.Save()
}

//挑战boss阵容
func (this *Player) ChallengeBoss(stage_id int32) {
	result4C := new(protocol.Game_ChallengeBossResult)

	var is_can_change bool = true
	if _, ok := this.Guaji_Stage.Guaji_Map_stage_pass[stage_id]; ok {
		if this.Guaji_Stage.Guaji_Map_stage_pass[stage_id] > 0 {
			is_can_change = true
		}
	}
	/*
		if this.Guaji_Stage.Now_Guaji_id == stage_id {
			guaji_killboss_con := Json_config.guaji_kill_boss_con[this.Guaji_Stage.Now_Guaji_id].Item0

			for _, v := range guaji_killboss_con {
				switch v.Con {
				case 101: //怪物
					if v.Par > this.Guaji_Stage.Kill_npc_num {
						is_can_change = false
					}
				case 102: //修炼时间
					if v.Par > this.Guaji_Stage.Guaji_time {
						is_can_change = false
					}
				case 103: //金币
					if v.Par > this.Guaji_Stage.Player_gold {
						is_can_change = false
					}
				case 104: //exp
					if v.Par > this.Guaji_Stage.Player_exp {
						is_can_change = false
					}
				default:
				}
			}
		}
	*/
	result4C.IsCanChange = &is_can_change
	if !is_can_change {
		encObj, _ := proto.Marshal(result4C)
		SendPackage(*this.conn, 1017, encObj)
		return
	}

	//怪物阵型
	team2 := new(protocol.Game_MonsterCombatTeam)
	var monster Monsters
	result := monster.GetMonsters(stage_id)
	if result == 0 {
		formations_2 := monster.dealMonster2Protocol()
		team2.MonsterAttrs = formations_2
	}

	result4C.Team_2 = team2

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.conn, 1017, encObj)
}

//挑战其他玩家
func (this *Player) ChallengePlayer(type_ int32, other_role_id int64) {
	result4C := new(protocol.Game_ChallengePlayerResult)
	if type_ == 1 { //挂机挑战其他玩家
		result := this.Guaji_Stage.IsCanPk(this.PlayerId, other_role_id)
		result4C.IsCanChange = &result

		if !result { //不能挑战直接返回
			encObj, _ := proto.Marshal(result4C)
			SendPackage(*this.conn, 1026, encObj)
		}

		if value, ok := word.players[other_role_id]; ok {
			var Game_HeroStruct []*protocol.Game_HeroStruct

			for _, v := range value.Heros { //遍历对应玩家的所有英雄列表
				if v.Hero_Info.Pos_stage > 0 { //关卡上阵英雄

					Hp := int32(v.Hero_Info.Hp)
					PhysicalAttack := int32(v.Hero_Info.Physical_attack)
					MagicAttack := int32(v.Hero_Info.Magic_attack)
					PhysicalDefense := int32(v.Hero_Info.Physical_defense)
					MagicDefense := int32(v.Hero_Info.Magic_defense)

					hero_struct := new(protocol.Game_HeroStruct)
					hero_struct.HeroInfo.Id = &v.Hero_Info.Hero_id
					hero_struct.HeroInfo.Uid = &v.Hero_Info.Hero_uid
					hero_struct.HeroInfo.Type = &v.Hero_Info.Type
					hero_struct.HeroInfo.Sex = &v.Hero_Info.Sex
					hero_struct.HeroInfo.Pos = &v.Hero_Info.Pos_stage
					hero_struct.HeroInfo.Level = &v.Hero_Info.Level
					hero_struct.HeroInfo.Exp = &v.Hero_Info.Exp
					hero_struct.HeroInfo.Hp = &Hp
					hero_struct.HeroInfo.StarLevel = &v.Hero_Info.Star_level
					hero_struct.HeroInfo.StepLevel = &v.Hero_Info.Step_level
					hero_struct.HeroInfo.Speed = &v.Hero_Info.Speed
					hero_struct.HeroInfo.Zodiac = &v.Hero_Info.Zodiac
					hero_struct.HeroInfo.Feature = &v.Hero_Info.Feature
					hero_struct.HeroInfo.Name = &v.Hero_Info.Name
					hero_struct.HeroInfo.PhysicalAttack = &PhysicalAttack
					hero_struct.HeroInfo.MagicAttack = &MagicAttack
					hero_struct.HeroInfo.PhysicalDefense = &PhysicalDefense
					hero_struct.HeroInfo.MagicDefense = &MagicDefense
					hero_struct.HeroInfo.Skill = v.Hero_Info.Skill

					//其他
					Game_HeroStruct = append(Game_HeroStruct, hero_struct)
				}
			}

			result4C.Team_2 = Game_HeroStruct
		}

		encObj, _ := proto.Marshal(result4C)
		SendPackage(*this.conn, 1026, encObj)

	}
}
