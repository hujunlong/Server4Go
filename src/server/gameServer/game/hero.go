//装备相关
package game

import (
	"errors"
)

import (
	"fmt"
	"net"
	"server/share/protocol"

	"github.com/golang/protobuf/proto"
)

type StepNeed struct { //升阶需要的东西
	goods_id int32
	count    int32
}

type Attribute struct {
	Group int32
	Key   int32
	Value float32
}

type HeroInfo struct {
	//hero基础属性
	Hero_id    int32                 //英雄id
	Hero_uid   int32                 //唯一uid
	Type       int32                 //主角 or 英雄 1：英雄 100主角
	Pos_stage  int32                 //位置（-1 表示未上阵）挂机 关卡阵型
	Pos_Arena  int32                 //位置（-1 表示未上阵）竞技场位置
	Level      int32                 //等级
	Exp        int32                 //经验
	Power      int32                 //战力
	Star_level int32                 //星级
	Step_level int32                 //阶级
	Skill      []int32               //技能
	Gift_group int32                 //天赋组id
	Attr       map[int32][]Attribute //属性 group=0总和 group=1基础属性 group=2升星 group=3升阶
}

type HeroStruct struct {
	Hero_Info HeroInfo //基础属性
}

func (this *HeroStruct) Init() {
	this.Hero_Info.Attr = make(map[int32][]Attribute)
}

func (this *HeroStruct) GetBase() { //对应等级血量

	id := this.Hero_Info.Hero_id
	level := this.Hero_Info.Level

	base_hp := Csv.hero[id].Id_130 //基础血
	hp_dev := Csv.hero[id].Id_109
	Base_Hp := base_hp + float32(level-1)*hp_dev

	physical_attack_dev := Csv.hero[id].Id_110 //物理攻击
	Base_Physical_attack := Csv.hero[id].Id_131 + float32(level-1)*physical_attack_dev

	magic_attack_dev := Csv.hero[id].Id_111 //魔法攻击
	Base_Magic_attack := Csv.hero[id].Id_132 + float32(level-1)*magic_attack_dev

	physical_def_dev := Csv.hero[id].Id_112 //物理防御
	Base_Physical_defense := Csv.hero[id].Id_133 + float32(level-1)*physical_def_dev

	magic_def_dev := Csv.hero[id].Id_113 //魔法防御
	Base_Magic_defense := Csv.hero[id].Id_134 + float32(level-1)*magic_def_dev

	this.AttrAdd(1, Base_Hp, Base_Physical_attack, Base_Magic_attack, Base_Physical_defense, Base_Magic_defense)
}

//添加属性
func (this *HeroStruct) AttrAdd(group int32, Base_Hp float32, Base_Physical_attack float32, Base_Magic_attack float32, Base_Physical_defense float32, Base_Magic_defense float32) {
	var attrs []Attribute

	var attr Attribute
	attr.Group = group

	attr.Key = 107 //生命
	attr.Value = Base_Hp
	attrs = append(attrs, attr)

	attr.Key = 108 //物理攻击
	attr.Value = Base_Physical_attack
	attrs = append(attrs, attr)

	attr.Key = 109 //魔法攻击
	attr.Value = Base_Magic_attack
	attrs = append(attrs, attr)

	attr.Key = 110 //物理防御
	attr.Value = Base_Physical_defense
	attrs = append(attrs, attr)

	attr.Key = 111 //法术防御
	attr.Value = Base_Magic_defense
	attrs = append(attrs, attr)

	if this.Hero_Info.Attr == nil {
		this.Init()
	}
	this.Hero_Info.Attr[group] = attrs
}

func (this *HeroStruct) StarAdd(star_level int32) { //升星每级增加比例
	pre := Csv.hero_star[star_level].Id_103 / 10000
	delete(this.Hero_Info.Attr, 2)
	this.AttrAdd(2, this.Hero_Info.Attr[1][0].Value, this.Hero_Info.Attr[1][1].Value*pre, this.Hero_Info.Attr[1][2].Value*pre, this.Hero_Info.Attr[1][3].Value*pre, this.Hero_Info.Attr[1][4].Value*pre)
}

//阶增加英雄属性
//阶增加对应属性比例(裸妆的最大值的万分比)
//计算裸妆最大值时候的值
func (this *HeroStruct) StepAdd(step_level int32) {
	extra_star_data_pre := Csv.hero_jinhua[step_level].Id_107 / 10000
	id := this.Hero_Info.Hero_id

	var max_level float32 = 100
	base_hp := Csv.hero[id].Id_130 //基础血
	hp_dev := Csv.hero[id].Id_109
	physical_attack_dev := Csv.hero[id].Id_110 //物理攻击
	magic_attack_dev := Csv.hero[id].Id_111    //魔法攻击
	physical_def_dev := Csv.hero[id].Id_112    //物理防御
	magic_def_dev := Csv.hero[id].Id_113       //魔法防御

	max_level_hp := base_hp + (max_level-1)*hp_dev
	max_level_physical_attack := max_level * physical_attack_dev
	max_level_magic_attack := max_level * magic_attack_dev
	max_level_physical_defense := max_level * physical_def_dev
	max_level_magic_defense := max_level * magic_def_dev

	delete(this.Hero_Info.Attr, 3)
	this.AttrAdd(3, max_level_hp*extra_star_data_pre, max_level_physical_attack*extra_star_data_pre, max_level_magic_attack*extra_star_data_pre, max_level_physical_defense*extra_star_data_pre, max_level_magic_defense*extra_star_data_pre)
}

//设置英雄位置
func (this *HeroStruct) SetHeroPos(type_ int32, pos int32) { //type = 1 竞技场 type=2关卡
	if type_ == 1 {
		this.Hero_Info.Pos_Arena = pos
	}

	if type_ == 2 {
		this.Hero_Info.Pos_stage = pos
	}
}

//创建英雄
func (this *HeroStruct) CreateHero(hero_id int32, player *Player) (int32, error) {

	if _, ok := Csv.hero[hero_id]; !ok {
		Log.Error("%d %s", "hero_id,input error", hero_id)
		return 0, errors.New("hero_id not found")
	}

	this.Hero_Info.Hero_id = hero_id   //id
	this.Hero_Info.Hero_uid = GetUid() //uid
	this.Hero_Info.Level = 1           //level
	this.Hero_Info.Exp = 0
	this.Hero_Info.Pos_Arena = -1
	this.Hero_Info.Pos_stage = -1
	this.Hero_Info.Type = Csv.hero[hero_id].Id_102 //角色类型
	this.Hero_Info.Gift_group = Csv.hero[hero_id].Id_122
	this.Hero_Info.Star_level = Csv.hero[hero_id].Id_106

	//基础属性
	this.GetBase()

	Log.Info("%d", this)

	player.Task.TriggerEvent(6, 1, 0) //任务
	return this.Hero_Info.Hero_uid, nil
}

//英雄加经验
func (this *HeroStruct) heroAddExp(exp int32, player *Player) {
	//找出需要加经验的hero

	csv_exp_int32 := Csv.hero_exp[this.Hero_Info.Level].Id_102

	for true {
		if this.Hero_Info.Level+exp > csv_exp_int32 {

			exp = this.Hero_Info.Exp + exp - csv_exp_int32
			this.Hero_Info.Level += 1

			if _, ok := Csv.role_exp[this.Hero_Info.Level]; !ok {
				this.Hero_Info.Level -= 1

				this.Hero_Info.Exp = Csv.hero_exp[this.Hero_Info.Level].Id_102
				break
			}

			//下一级需要exp
			csv_exp_int32 = Csv.hero_exp[this.Hero_Info.Level].Id_102
		} else {
			this.Hero_Info.Exp += exp
			break
		}
	}
	Log.Info("level =%d exp = %d", this.Hero_Info.Level, this.Hero_Info.Exp)

	this.Note2CHeroChange(player.conn)
}

func (this *HeroStruct) TaskGoods(prop int32, num int32) int32 { //0成功 1:道具不足
	return 0
}

//英雄升阶
func (this *HeroStruct) StepUp(player *Player) int32 { //0:ok 1：等级不足 2：金币不足 3：材料不足 4：非法数据

	//检查 经验 金币 物品是否足够
	step_level := this.Hero_Info.Step_level
	var key int32 = step_level + 1

	if key == 0 {
		return 4
	}

	if this.Hero_Info.Level < Csv.hero_jinhua[key].Id_104 {
		fmt.Println(this.Hero_Info.Level, Csv.hero_jinhua[key].Id_104)
		return 1
	}

	if player.Info.Gold < Csv.hero_jinhua[key].Id_106 {
		fmt.Println(player.Info.Gold, Csv.hero_jinhua[key].Id_106)
		return 2
	}

	for _, v := range Csv.hero_jinhua[key].NeedEquip {
		fmt.Println(player.Bag_Prop.PropsById[v.Id], v.Num)
		if player.Bag_Prop.PropsById[v.Id] < v.Num {
			return 3
		}
	}

	//扣钱
	player.ModifyGold(-Csv.hero_jinhua[key].Id_106)

	//扣除物品
	for _, v := range Csv.hero_jinhua[key].NeedEquip {
		player.Bag_Prop.DeleteItemById(v.Id, v.Num, player.conn)
	}

	//更改英雄属性
	this.Hero_Info.Step_level += 1
	player.Task.TriggerEvent(12, 1, 0)
	//添加属性
	this.StepAdd(this.Hero_Info.Step_level)

	//推送
	this.Note2CHeroChange(player.conn)

	//任务系统
	player.Task.TriggerEvent(7, 1, this.Hero_Info.Step_level)

	player.Save()

	return 0
}

//英雄升星
func (this *HeroStruct) StarUp(player *Player) int32 { //0:ok 1：材料不足 2：非法数据
	hero_id := this.Hero_Info.Hero_id
	star_level := this.Hero_Info.Star_level

	//读取需要碎片id
	goods_id := int32(Csv.hero[hero_id].Id_119)
	need_count := Csv.hero_star[star_level+1].Id_102

	//检查材料是否足够
	if player.Bag_Prop.PropsById[goods_id] < need_count {
		return 1
	}

	//扣除材料
	player.Bag_Prop.DeleteItemById(goods_id, need_count, player.conn)

	this.Hero_Info.Star_level += 1

	this.StarAdd(this.Hero_Info.Star_level)

	//推送
	this.Note2CHeroChange(player.conn)

	player.Task.TriggerEvent(8, 1, this.Hero_Info.Star_level)

	player.Save()
	return 0
}

//天赋开启
func (this *HeroStruct) OpenHeroGift(ids []int32) int32 { //0 ok 1:等级不足 2:主角才能升级天赋 ids暂定为技能id

	if this.Hero_Info.Gift_group == 0 {
		return 2
	}

	for _, id := range ids {
		level := Csv.role_gift[id].Id_105
		if this.Hero_Info.Level < level {
			return 1
		}
	}

	for _, id := range ids {
		this.Hero_Info.Skill = append(this.Hero_Info.Skill, id)
	}

	return 0
}

//英雄对应属性推送
func (this *HeroStruct) Note2CHeroChange(conn *net.Conn) {
	result4C := new(protocol.NoticeMsg_Notice2CHero)
	result4C.HeroUid = &this.Hero_Info.Hero_uid
	result4C.StepLevel = &this.Hero_Info.Step_level
	result4C.StarLevel = &this.Hero_Info.Star_level
	result4C.Level = &this.Hero_Info.Level
	result4C.Exp = &this.Hero_Info.Exp

	var HeroAttrs []*protocol.Attribute
	for _, v1 := range this.Hero_Info.Attr {
		for _, v2_buff := range v1 {
			v2 := v2_buff
			HeroAttr := new(protocol.Attribute)
			HeroAttr.Group = &v2.Group
			HeroAttr.Key = &v2.Key
			HeroAttr.Value = &v2.Value
			HeroAttrs = append(HeroAttrs, HeroAttr)

			fmt.Println("attr:", v2.Group)
			fmt.Println("attr:", v2.Key)
			fmt.Println("attr:", v2.Value)

		}
	}

	result4C.HeroAttr = HeroAttrs
	encObj, _ := proto.Marshal(result4C)
	SendPackage(*conn, 1209, encObj)
}
