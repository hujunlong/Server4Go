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
	LoginDays int32  //连续登陆次数
	Option    []bool //设置()
}

type Player struct {
	PlayerId      int64                 //账号id + server_id 注意区别
	ProfessionId  int32                 //职业id
	CreateTime    int32                 //创建时间
	LastTime      int32                 //上次在线时间
	Info          *PlayerInfo           //角色基础属性
	Stage         *MapStage             //副本关卡
	Bag_Prop      *BagProp              //道具背包
	Bag_Equip     *BagEquip             //装备背包
	Guaji_Stage   *GuajiMapStage        //挂机关卡
	StageFomation *Fomation             //阵型(挂机跟关卡阵型)
	ArenaFomation *Fomation             //阵型 竞技场
	Heros         map[int32]*HeroStruct //英雄列表 key:uid
	Gm            *GM                   //GM
	Task          *TaskStruct           //任务
	Niu_Dan       *NiuDan               //扭蛋抽卡
	conn          *net.Conn             //角色的网络连接,不需要保存
}

func LoadPlayer(id_str string) *Player { //从数据库读取玩家信息
	player := new(Player)
	player.Init()
	err := redis.Find(id_str, player)
	if err == nil {
		return player
	}
	return nil
}

func (this *Player) Init() {
	fmt.Println("player 初始化")
	//英雄相关
	this.Heros = make(map[int32]*HeroStruct)

	//角色基础属性
	this.Info = new(PlayerInfo)

	//副本地图
	this.Stage = new(MapStage)
	this.Stage.Init()

	//挂机地图
	this.Guaji_Stage = new(GuajiMapStage)
	this.Guaji_Stage.Init(this)

	//背包
	this.Bag_Prop = new(BagProp)
	this.Bag_Prop.Init(this)

	this.Bag_Equip = new(BagEquip)
	this.Bag_Equip.Init()

	//阵型(关卡)
	this.StageFomation = new(Fomation)
	this.StageFomation.Init()

	//阵型(竞技场)
	this.ArenaFomation = new(Fomation)
	this.ArenaFomation.Init()

	//Gm
	this.Gm = new(GM)

	//扭蛋
	this.Niu_Dan = new(NiuDan)
	this.Niu_Dan.Init()

	//任务系统
	this.Task = new(TaskStruct)
	this.Task.Init(this)
}

func (this *Player) SetConn(conn *net.Conn) {
	this.conn = conn
}

func (this *Player) Login(player_id int64, conn *net.Conn) (bool, *Player) {
	player_id_str := strconv.FormatInt(player_id, 10)
	var is_login bool = false
	this = LoadPlayer(player_id_str)

	if this != nil {
		if this.CreateTime <= 0 {
			fmt.Println("创建时间：", this.CreateTime)
			//如果没创建 直接返回 false
			var is_create bool = false
			result4C := &protocol.PlayerBase_RoleInfoResult{
				IsCreate: &is_create,
			}
			encObj, _ := proto.Marshal(result4C)
			SendPackage(*conn, 1002, encObj)
			return is_login, nil
		}
	} else {
		return is_login, nil
	}

	this.conn = conn
	//先判断是否创建角色
	var is_create bool = true

	//英雄列表
	heros_FightingAttr := this.GetHeroStruct(this.Heros)

	//道具
	props_struct := this.DealPropStruct(this.Bag_Prop.Props)

	//装备
	equips_struct := this.DealEquipStruct(this.Bag_Equip.BagEquip)

	//副本开启列表
	var type_ int32 = 1
	var copy_levels []*protocol.Stage
	for buff_id, buff_v := range this.Stage.Map_stage_pass {
		id := buff_id
		v := buff_v
		copy_level := new(protocol.Stage)
		copy_level.Type = &type_
		copy_level.State = &v
		copy_level.StageId = &id
		copy_levels = append(copy_levels, copy_level)
	}
	//挂机开启列表
	var type_guaji int32 = 2
	var guaji_stages []*protocol.Stage
	for buff_id, buff_v := range this.Guaji_Stage.Guaji_Map_stage_pass {
		v := buff_v
		id := buff_id
		guaji_stage := new(protocol.Stage)
		guaji_stage.Type = &type_guaji
		guaji_stage.State = &v
		stage_id_32 := &id
		guaji_stage.StageId = stage_id_32
		guaji_stages = append(guaji_stages, guaji_stage)
	}

	result4C := &protocol.PlayerBase_RoleInfoResult{
		IsCreate: &is_create,
		PlayerInfo: &protocol.PlayerBase_PlayerInfo{
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
		Heros:        heros_FightingAttr,
		Equips:       equips_struct,
		Props:        props_struct,
		CopyLevels:   copy_levels,
		HangupLevels: guaji_stages,
		ProfessionId: &this.ProfessionId,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*conn, 1002, encObj)
	is_login = true

	//发送离线收益
	//this.Guaji_Stage.OffNotice2CGuaji(this)

	//添加到挂机全局中
	global_guaji_players.Enter(this.Guaji_Stage.Now_Guaji_id, this.PlayerId)

	//添加到world中
	word.EnterWorld(this)

	//添加道具
	//this.Test()
	return is_login, this
}

func (this *Player) RegisterRole(player_id int64, nick string, HeroId int32, conn *net.Conn) int32 { //id role id

	key := strconv.Itoa(int(player_id))
	redis.Modify(player_id, "")
	player := LoadPlayer(key)

	if player == nil {
		redis.Modify(player_id, "")
		//redis.Modify()
		//return global.REGISTERROLEERROR
	}

	if player.CreateTime > 0 {
		return global.ALREADYHAVE
	}

	//检测heroid是否在配置中 读取hero_create表
	if _, ok := Csv.create_role[HeroId]; !ok {
		return global.CSVHEROIDEERROR
	}

	//创建player
	this.Init()

	//体力
	int32_energy := int32(Csv.property[2056].Id_102)

	//英雄相关
	hero := new(HeroStruct)
	hero_uid, err := hero.CreateHero(HeroId, this)
	if err == nil {
		this.Heros[hero_uid] = hero
		this.Heros[hero_uid].SetHeroPos(2, 1) //默认设置在第一个位置
	}

	//将英雄添加到阵型中
	this.StageFomation.OnFomation(1, hero_uid) //添加到1号位置

	this.PlayerId = player_id
	this.ProfessionId = HeroId
	this.CreateTime = int32(time.Now().Unix())
	this.LastTime = this.CreateTime
	this.Info.Level = 1
	this.Info.Exp = 0
	this.Info.Hp = int32(hero_uid)
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
	this.Stage.Map_stage_pass[10101] = 1

	//初始化挂机关卡
	this.Guaji_Stage.Now_Guaji_id = 20101
	this.Guaji_Stage.Guaji_Map_stage_pass[20101] = 0 //解锁未通关
	this.Guaji_Stage.SetCurrentStage(this.Guaji_Stage.Now_Guaji_id)

	//初始化背包开启个数
	this.Bag_Equip.OpenCount = int32(Csv.property[2018].Id_102)
	this.Bag_Equip.Max = int32(Csv.property[2019].Id_102)
	this.Bag_Equip.UseCount = 0

	this.Bag_Prop.OpenCount = int32(Csv.property[2018].Id_102)
	this.Bag_Prop.Max = int32(Csv.property[2019].Id_102)
	this.Bag_Prop.UseCount = 0

	//初始化任务系统
	this.Task.CreateNewTask(1, 2001)
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
func (this *Player) AddRoleExp(exp int32) {
	this.Task.TriggerEvent(4, exp, 0) //任务id = 4 主角获得多少经验

	csv_exp_int32 := Csv.role_exp[this.Info.Level].Id_102
	for true {
		if this.Info.Exp+exp > csv_exp_int32 {
			exp -= csv_exp_int32
			this.Info.Level += 1
			if _, ok := Csv.role_exp[this.Info.Level]; !ok {
				//玩家满级判断
				this.Info.Level -= 1
				this.Info.Exp = Csv.role_exp[this.Info.Level].Id_102
				break
			}

			//下一级需要exp
			csv_exp_int32 = Csv.role_exp[this.Info.Level].Id_102
		} else {
			this.Info.Exp = this.Info.Exp + exp
			break
		}
	}
	this.Notice2CRoleInfo()
	this.Task.TriggerEvent(3, this.Info.Level, 0) //任务id = 3 主角达到x级
	Log.Info("level =%d exp = %d", this.Info.Level, this.Info.Exp)
}

//钱变化
func (this *Player) ModifyGold(num int32) {
	if num == 0 {
		return
	}

	this.Info.Gold += num
	if this.Info.Gold < 0 {
		this.Info.Gold = 0
	}
	this.Notice2CMoney()
	this.Task.TriggerEvent(5, this.Info.Gold, 0) //任务id = 5 获取xx铜钱
}

//元宝变化
func (this *Player) ModifyDiamond(num int32) {
	if num == 0 {
		return
	}

	this.Info.Diamond += num
	if this.Info.Diamond < 0 {
		this.Info.Diamond = 0
	}
	this.Notice2CMoney()
}

//体力
func (this *Player) ModifyEnergy(num int32) {
	this.Info.Energy += num
	if this.Info.Energy < 0 {
		this.Info.Energy = 0
	}

	if this.Info.Energy > this.Info.EnergyMax {
		this.Info.Energy = this.Info.EnergyMax
	}

	this.Notice2CEnergy()
}

//推送数据
func (this *Player) Notice2CRoleInfo() {
	result4C := &protocol.NoticeMsg_Notice2CRoleInfo{
		Level: &this.Info.Level,
		Exp:   &this.Info.Exp,
		Power: &this.Info.Power,
		Hp:    &this.Info.Hp,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.conn, 1203, encObj)
}

//经济变化
func (this *Player) Notice2CMoney() {
	result4C := &protocol.NoticeMsg_Notice2CMoney{
		Gold:    &this.Info.Gold,
		Diamond: &this.Info.Diamond,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.conn, 1204, encObj)
}

//体力变化
func (this *Player) Notice2CEnergy() {
	result4C := &protocol.NoticeMsg_Notice2CEnergy{
		Energy:    &this.Info.Energy,
		EnergyMax: &this.Info.EnergyMax,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.conn, 1205, encObj)
}

func (this *Player) DealEquipStruct(Equips map[int32]*Equip) []*protocol.EquipStruct {
	//装备
	var my_equips []Equip
	for _, v := range Equips {
		my_equips = append(my_equips, *v)
	}

	return this.Bag_Equip.DealEquipStruct(my_equips)

}

func (this *Player) DealPropStruct(Props map[int32]*Prop) []*protocol.PropStruct {
	//道具
	var props []*protocol.PropStruct
	for i, _ := range Props {
		prop := new(protocol.PropStruct)
		prop.PropId = &Props[i].Prop_id
		prop.PropUid = &Props[i].Prop_uid
		prop.PropCount = &Props[i].Count
		props = append(props, prop)
	}
	return props
}

//玩家退出
func (this *Player) ExitGame() {
	this.LastTime = int32(time.Now().Unix())
	this.Save()
	word.ExitWorld(this.PlayerId)
	global_guaji_players.Exit(this.Guaji_Stage.Now_Guaji_id, this.PlayerId)
	redis.Modify("uid", Global_Uid)
}

//获取指定玩家阵型
func (this *Player) GetFomation(role_id int64, type_ int32) {
	var result int32 = 0
	var formations []*protocol.Formation_FormationStruct
	result4C := new(protocol.Formation_GetGuajiRoleFormationResult)
	if value, ok := word.players[role_id]; ok {
		switch type_ { //(1竞技场，2副本&挂机)
		case 1:
		case 2:
			fomations := value.StageFomation.Hero_fomations
			proto_formation := new(protocol.Formation_FormationStruct)
			for _, v := range fomations {
				proto_formation.Pos = &v.Pos_id
				proto_formation.HeroUid = &v.Hero_uid
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
	SendPackage(*this.conn, 1301, encObj)
	Log.Info("role_id = %d type_ = %d", role_id, type_)
}

//获取英雄通信相关
func (this *Player) GetHeroStruct(Heros map[int32]*HeroStruct) []*protocol.FightingAttr {
	var Game_HeroStruct []*protocol.FightingAttr
	for _, v_buff := range Heros { //遍历对应玩家的所有英雄列表
		v := v_buff
		HeroInfo := new(protocol.FightingAttr)
		HeroInfo.Id = &v.Hero_Info.Hero_id
		HeroInfo.Uid = &v.Hero_Info.Hero_uid
		HeroInfo.Pos = &v.Hero_Info.Pos_stage
		HeroInfo.Level = &v.Hero_Info.Level
		HeroInfo.Exp = &v.Hero_Info.Exp
		HeroInfo.StarLevel = &v.Hero_Info.Star_level
		HeroInfo.StepLevel = &v.Hero_Info.Step_level
		HeroInfo.Skill = v.Hero_Info.Skill

		//额外属性
		var attributes []*protocol.Attribute
		for _, v1 := range v.Hero_Info.Attr {
			for _, v2_buff := range v1 {
				v2 := v2_buff
				attribute := new(protocol.Attribute)
				attribute.Group = &v2.Group
				attribute.Key = &v2.Key
				attribute.Value = &v2.Value
				attributes = append(attributes, attribute)
			}
		}
		HeroInfo.Attribute = attributes

		Game_HeroStruct = append(Game_HeroStruct, HeroInfo)
	}
	return Game_HeroStruct
}

//挑战其他玩家
func (this *Player) ChallengePlayer(type_ int32, other_role_id int64) {
	result4C := new(protocol.StageBase_ChallengePlayerResult)
	if type_ == 1 { //挂机挑战其他玩家
		result := this.Guaji_Stage.IsCanPk(this.PlayerId, other_role_id)
		result4C.IsCanChange = &result

		if !result { //不能挑战直接返回
			encObj, _ := proto.Marshal(result4C)
			SendPackage(*this.conn, 1111, encObj)
		}

		Heros := make(map[int32]*HeroStruct)
		if value, ok := word.players[other_role_id]; ok {
			for k, v := range value.Heros { //遍历对应玩家的所有英雄列表
				if v.Hero_Info.Pos_stage > 0 { //关卡上阵英雄
					Heros[k] = v
				}
			}
			Game_HeroStruct := this.GetHeroStruct(Heros)
			result4C.Team_2 = Game_HeroStruct
		}

		encObj, _ := proto.Marshal(result4C)
		SendPackage(*this.conn, 1111, encObj)
	}
}

//查看在线玩家的英雄信息
func (this *Player) GetHeroResult(role_id int64, type_ int32) { //0:ok 1:role_id 未找到或者未上线
	var result int32 = 0
	result4C := new(protocol.Hero_GetHerosResult)

	if _, ok := word.players[role_id]; !ok {
		result = 1
	} else {
		Heros := make(map[int32]*HeroStruct)
		switch type_ { //全部英雄 2:只获取竞技场上的英雄 3:挂机 关卡阵型上的英雄
		case 1:
			Heros = word.players[role_id].Heros
		case 2:
			for key, v := range word.players[role_id].Heros {
				if v.Hero_Info.Pos_Arena > 0 {
					Heros[key] = v
				}
			}
		case 3:
			for key, v := range word.players[role_id].Heros {
				if v.Hero_Info.Pos_stage > 0 {
					Heros[key] = v
				}
			}
		}
		result4C.Heros = this.GetHeroStruct(Heros)
	}

	result4C.Result = &result
	encObj, _ := proto.Marshal(result4C)
	fmt.Println(this.CreateTime, this.Info.Nick, "111", this.Info)
	SendPackage(*this.conn, 1401, encObj)
}
