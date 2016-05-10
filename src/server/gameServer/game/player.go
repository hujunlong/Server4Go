//角色
package game

import (
	"fmt"
	"net"
	"server/share/global"
	"strconv"
	"time"

	"server/share/protocol"

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
	PlayerId    int64          //账号id + server_id 注意区别
	CreateTime  int32          //创建时间
	LastTime    int32          //上次在线时间
	Heros       []HeroStruct   //英雄列表
	Info        PlayerInfo     //角色基础属性
	Stage       *MapStage      //副本关卡
	Bag_Prop    *BagProp       //道具背包
	Bag_Equip   *BagEquip      //装备背包
	Guaji_Stage *GuajiMapStage //挂机关卡
	conn        *net.Conn      //角色的网络连接,不需要保存
}

func LoadPlayer(id string) *Player { //从数据库读取玩家信息
	player := new(Player)
	err := redis.Find(id, player)
	if err == nil {
		return player
	}
	return nil
}

func (this *Player) Init() {

	//this.Heros = make([]HeroStruct, 1, 1)
	//副本地图
	this.Stage = new(MapStage)
	this.Stage.Init()

	//挂机地图
	this.Guaji_Stage = new(GuajiMapStage)
	this.Guaji_Stage.Init()

	//背包
	this.Bag_Prop = new(BagProp)
	this.Bag_Equip = new(BagEquip)

	this.CreateTime = 0

}

func (this *Player) RegisterRole(player_id int64, nick string, HeroId int32, conn *net.Conn) int32 { //id role id

	//检查是否已经拥有了key
	key := strconv.Itoa(int(player_id))
	is_exists, _ := redis.Exists(key)

	//检测heroid是否在配置中 读取hero_create表
	HeroId_str := strconv.Itoa(int(HeroId))
	create_role_map := Csv.create_role.simple_info_map[HeroId_str]
	if create_role_map == nil {
		return global.CSVHEROIDEERROR
	}

	//读取hero表
	index_value := Csv.hero.index_value
	hero_map := Csv.hero.simple_info_map[HeroId_str]

	if is_exists {

		//player属性
		var base_mult int32 = 3                              //基础倍数
		hp_growth := Str2int32(hero_map[index_value["109"]]) //生命成长

		//体力
		index := Csv.property.index_value["102"]
		int_energy, _ := strconv.Atoi(Csv.property.simple_info_map["2056"][index])

		this.PlayerId = player_id
		this.CreateTime = int32(time.Now().Unix())
		this.LastTime = this.CreateTime
		this.Info.Level = 1
		this.Info.Exp = 0
		this.Info.Hp = hp_growth * base_mult
		this.Info.Nick = nick
		this.Info.Gold = 0
		this.Info.Vip = 0
		this.Info.Energy = int32(int_energy)
		this.Info.EnergyMax = this.Info.Energy
		this.Info.Diamond = 0
		this.Info.Power = 0
		this.Info.Signature = ""
		this.conn = conn

		//英雄相关(具体数据待完成)
		var hero_struct HeroStruct
		hero_struct.Hero_Info.Hero_id = HeroId
		hero_struct.Hero_Info.Hero_uid = GetUid()
		hero_struct.Hero_Info.Exp = 0
		hero_struct.Hero_Info.Hp = hp_growth * base_mult
		hero_struct.Hero_Info.Power = 0
		hero_struct.Hero_Info.Star_level = 1
		hero_struct.Hero_Info.Step_level = 1
		this.Heros = append(this.Heros, hero_struct)

		//初始化副本关卡
		fmt.Println(this.Stage.Map_stage_pass)
		this.Stage.Map_stage_pass["0"] = 1

		//初始化挂机关卡
		this.Guaji_Stage.Now_Guaji_id = 20101
		this.Guaji_Stage.Guaji_Map_stage_pass["20101"] = 0 //解锁未通关
		this.Guaji_Stage.SetCurrentStage(this.Guaji_Stage.Now_Guaji_id)

		//初始化背包
		begin_open_count, _ := strconv.Atoi(Csv.property.simple_info_map["2018"][index])
		open_max, _ := strconv.Atoi(Csv.property.simple_info_map["2019"][index])

		this.Bag_Equip.OpenIndex = int32(begin_open_count)
		this.Bag_Equip.Max = int32(open_max)
		this.Bag_Equip.UseCount = 0
		this.Bag_Prop.OpenIndex = int32(begin_open_count)
		this.Bag_Prop.Max = int32(open_max)
		this.Bag_Prop.UseCount = 0
		//写内存数据库
		err := this.Save()
		if err != nil {
			fmt.Println("write err:", err)
		}
		return global.REGISTERROLESUCCESS
	}
	return global.REGISTERROLEERROR
}

func (this *Player) Save() error {
	err := redis.Modify(strconv.Itoa(int(this.PlayerId)), this)
	return err
}

//往主角装备背包加东西
func (this *Player) addEquip(equip Equip) {
	if this.Bag_Equip.OpenIndex > this.Bag_Equip.UseCount {
		this.Bag_Equip.BagEquip = append(this.Bag_Equip.BagEquip, equip)
		this.Bag_Equip.UseCount += 1
	}
}

//往主角道具背包加东西
func (this *Player) addProp(prop Prop) {
	if this.Bag_Prop.OpenIndex > this.Bag_Prop.UseCount {
		this.Bag_Prop.Props = append(this.Bag_Prop.Props, prop)
		this.Bag_Equip.UseCount += 1
	}
}

//主角加经验
func (this *Player) addExp(exp int32) {
	var add_level int32 = 0
	role_level_str := strconv.Itoa(int(this.Info.Level))

	exp_index := Csv.role_exp.index_value["102"]
	exp_int, _ := strconv.Atoi(Csv.role_exp.simple_info_map[role_level_str][exp_index])

	for true {
		if this.Info.Exp+exp > int32(exp_int) {
			add_level += 1
			exp -= int32(exp_int)

			role_level_str = strconv.Itoa(int(this.Info.Level + add_level))
			exp_int, _ = strconv.Atoi(Csv.role_exp.simple_info_map[role_level_str][exp_index])
		} else {
			this.Info.Exp = this.Info.Exp + exp
			this.Info.Level += add_level
			break
		}
	}
}

func (this *Player) dealEquipStruct(equips []Equip) []*protocol.Game_EquipStruct {
	//装备
	var equip_infos []*protocol.Game_EquipInfo
	for _, k := range equips {
		equip_info := new(protocol.Game_EquipInfo)
		equip_info.Id = &k.equip_id
		equip_info.Uid = &k.equip_uid
		equip_info.EquipLevel = &k.equip_level
		equip_info.StrengthenCount = &k.strengthen_count
		equip_info.Pos = &k.pos
		equip_info.Quality = &k.quality

		equip_infos = append(equip_infos, equip_info)
	}

	var equips_struct []*protocol.Game_EquipStruct
	for _, k := range equip_infos {
		equip_struct := new(protocol.Game_EquipStruct)
		equip_struct.EquipInfo = k
		equips_struct = append(equips_struct, equip_struct)
	}
	return equips_struct
}

func (this *Player) dealPropStruct(props_ []Prop) []*protocol.Game_PropStruct {
	//道具
	var props []*protocol.Game_PropStruct
	for _, k := range props_ {
		prop := new(protocol.Game_PropStruct)
		prop.PropId = &k.Prop_id
		prop.PropUid = &k.Prop_uid
		props = append(props, prop)
	}
	return props
}

//道具跟装备
func (this *Player) getEquipsAndProps(id int32) ([]*protocol.Game_EquipStruct, []*protocol.Game_PropStruct) {
	//产生通过奖励
	this.Stage.Reward(id)

	//主角加钱 经验
	this.addExp(this.Stage.player_exp)
	this.Info.Gold += this.Stage.player_gold
	this.addExp(this.Stage.hero_exp)

	//装备
	for _, k := range this.Stage.equips {
		//往背包添加道具
		this.addEquip(k)
	}

	//道具
	for _, k := range this.Stage.props {
		//往背包添加道具
		this.addProp(k)
	}

	//推送相关属性变化
	this.Notice2CEquip(1, this.Stage.equips)
	this.Notice2CProp(1, this.Stage.props)
	this.Notice2CRoleInfo()
	this.Notice2CMoney()
	this.Notice2CEnergy()

	equips_struct := this.dealEquipStruct(this.Stage.equips)
	props := this.dealPropStruct(this.Stage.props)
	return equips_struct, props
}

//关卡奖励获取
func (this *Player) WarMapNoteServerResult(state int32, id int32) {
	if state < 1 { //未通关
		return
	}
	equips_struct, props := this.getEquipsAndProps(id)
	result4C := &protocol.Game_WarMapNoteServerResult{
		Reward: &protocol.Game_Reward{
			PlayerExp:  &this.Stage.player_exp,
			PlayerGold: &this.Stage.player_gold,
			HeroExp:    &this.Stage.hero_exp,
			Equips:     equips_struct,
			Props:      props,
		},
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.conn, 1004, encObj)
}

//扫荡
func (this *Player) SweepMapStageResult(stage_id int32, sweep_count int32) {
	result := this.Stage.IsCanThroughMap(stage_id, this.Info.Energy, sweep_count)
	if result != 0 {
		return
	}

	var Game_Reward_ []*protocol.Game_Reward
	var i int32 = 0
	for ; i < sweep_count; i++ {
		equips_struct, props := this.getEquipsAndProps(stage_id)
		reward_ := &protocol.Game_Reward{
			PlayerExp:  &this.Stage.player_exp,
			PlayerGold: &this.Stage.player_gold,
			HeroExp:    &this.Stage.hero_exp,
			Equips:     equips_struct,
			Props:      props,
		}
		Game_Reward_ = append(Game_Reward_, reward_)
	}

	result4C := &protocol.Game_SweepMapStageResult{
		Result: &result,
		Reward: Game_Reward_,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.conn, 1005, encObj)

}

//装备变化
func (this *Player) Notice2CEquip(type_ int32, equips []Equip) {
	equips_struct := this.dealEquipStruct(equips)

	result4C := &protocol.Game_Notice2CEquip{
		Type:  &type_,
		Equip: equips_struct,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.conn, 1008, encObj)
}

//道具变化
func (this *Player) Notice2CProp(type_ int32, props []Prop) {
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
func (this *Player) Notice2CheckPoint(type_ int32, id int32) { //状态 (-1 未通关  0解锁未通关 1 一星级通关 2二星通关 3三星通关)
	result4C := &protocol.Game_Notice2CheckPoint{
		LatestCheckpoint: &protocol.Game_Stage{
			State:   &type_,
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
}

//离线收益
func (this *Player) OffNotice2CGuaji() {
	gold, exp, can_add_power := this.Guaji_Stage.OffNotice2CGuaji(this.LastTime)
	this.Info.Gold += gold
	this.Info.Exp += exp
	this.Info.Energy += can_add_power

	this.addExp(this.Info.Exp)
	this.Notice2CRoleInfo()
	this.Notice2CMoney()
	this.Notice2CEnergy()

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
	equips_struct := this.dealEquipStruct(equips)
	props_struct := this.dealPropStruct(props)

	result4C := &protocol.Game_C2SChallengeResult{
		PropStruct:  props_struct,
		EquipStruct: equips_struct,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.conn, 1018, encObj)

	this.Notice2CEquip(1, equips)
	this.Notice2CProp(1, props)

	//推送下个挂机关卡
	var type_32 int32 = 2     //状态 挂机
	var state_const int32 = 0 //解锁未通关
	next_stage_id_index := Csv.map_guaji.index_value["102"]
	now_stage_str := strconv.Itoa(int(stage_id))
	next_stage_id_str := Csv.map_guaji.simple_info_map[now_stage_str][next_stage_id_index]
	next_stage_id_int32 := Str2Int32(next_stage_id_str)

	result4C2 := &protocol.Game_Notice2CheckPoint{
		LatestCheckpoint: &protocol.Game_Stage{
			Type:    &type_32,
			State:   &state_const,
			StageId: &next_stage_id_int32,
		},
	}
	encObj2, _ := proto.Marshal(result4C2)
	SendPackage(*this.conn, 1013, encObj2)
}
