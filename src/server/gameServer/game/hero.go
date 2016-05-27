//装备相关
package game

import (
	"fmt"
)

type HeroInfo struct { //hero基础属性
	Hero_id          int32             //英雄id
	Hero_uid         int32             //唯一uid
	Type             int32             //主角 or 英雄 1：英雄 100主角
	Sex              int32             //性别
	Pos_stage        int32             //位置（-1 表示未上阵）挂机 关卡阵型
	Pos_Arena        int32             //位置（-1 表示未上阵）竞技场位置
	Level            int32             //等级
	Exp              int32             //经验
	Hp               float32           //血量
	Power            int32             //战力
	Star_level       int32             //星级
	Step_level       int32             //阶级
	Speed            int32             //速度
	Zodiac           int32             //属相 金木水火土
	Feature          string            //特征
	Name             string            //名称
	Physical_attack  float32           //物理攻击
	Magic_attack     float32           //法术攻击
	Physical_defense float32           //物理防御
	Magic_defense    float32           //魔法防御
	Skill            []int32           //技能
	Attrs            map[int32]float32 //额外属性

}

type HeroStruct struct {
	Hero_Info HeroInfo
}

func (this *HeroStruct) GetBase(id int32, level int32, star int32) { //对应等级血量
	star_pre := this.StarPre(star)

	index := Csv.hero.index_value["109"] //血
	hp_dev := Csv.hero.simple_info_map[id][index]
	this.Hero_Info.Hp = float32(level+2) * hp_dev * (star_pre + 1)

	index = Csv.hero.index_value["110"] //物理攻击
	physical_attack_dev := Csv.hero.simple_info_map[id][index]
	this.Hero_Info.Physical_attack = float32(level+2) * physical_attack_dev * (star_pre + 1)

	index = Csv.hero.index_value["111"] //魔法攻击
	magic_attack_dev := Csv.hero.simple_info_map[id][index]
	this.Hero_Info.Magic_attack = float32(level+2) * magic_attack_dev * (star_pre + 1)

	index = Csv.hero.index_value["112"] //物理防御
	physical_def_dev := Csv.hero.simple_info_map[id][index]
	this.Hero_Info.Physical_defense = float32(level+2) * physical_def_dev * (star_pre + 1)

	index = Csv.hero.index_value["113"] //魔法防御
	magic_def_dev := Csv.hero.simple_info_map[id][index]
	this.Hero_Info.Magic_defense = float32(level+2) * magic_def_dev * (star_pre + 1)
}

func (this *HeroStruct) StarPre(star_level int32) float32 { //升星每级增加比例
	key := Csv.hero_star.index_value["103"]
	extra_star_data_pre := float32(Csv.hero_star.simple_info_map[star_level][key]) / 10000
	return extra_star_data_pre
}

func (this *HeroStruct) StepPre(stp_level int32) { //升阶增加的百分比

	var step_level_index int32 = 0
	index_type := Csv.hero_jinhua.index_value["201"]
	index_jie := Csv.hero_jinhua.index_value["102"]

	for k, v := range Csv.hero_jinhua.simple_info_map {
		if v[index_type] == 1 && v[index_jie] == stp_level {
			step_level_index = k
			break
		}
	}

	if step_level_index == 0 {
		return
	}

	key := Csv.hero_jinhua.index_value["112"] //加血量
	add_hp_pre := (float32(Csv.hero_jinhua.simple_info_map[step_level_index][key]) / 10000)

	key = Csv.hero_jinhua.index_value["113"] //加物理攻击
	add_physic_attack_pre := (float32(Csv.hero_jinhua.simple_info_map[step_level_index][key]) / 10000)

	key = Csv.hero_jinhua.index_value["112"] //魔法攻击
	add_magic_attack_pre := (float32(Csv.hero_jinhua.simple_info_map[step_level_index][key]) / 10000)

	key = Csv.hero_jinhua.index_value["113"] //加物理防御
	add_physic_def_pre := (float32(Csv.hero_jinhua.simple_info_map[step_level_index][key]) / 10000)

	key = Csv.hero_jinhua.index_value["114"] //增加法术防御
	add_magic_def_pre := (float32(Csv.hero_jinhua.simple_info_map[step_level_index][key]) / 10000)

	this.Hero_Info.Hp = this.Hero_Info.Hp * (1 + add_hp_pre)
	this.Hero_Info.Physical_attack = this.Hero_Info.Physical_attack * (1 + add_physic_attack_pre)
	this.Hero_Info.Magic_attack = this.Hero_Info.Magic_attack * (1 + add_magic_attack_pre)
	this.Hero_Info.Physical_defense = this.Hero_Info.Physical_defense * (1 + add_physic_def_pre)
	this.Hero_Info.Magic_defense = this.Hero_Info.Magic_defense * (1 + add_magic_def_pre)
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
func (this *HeroStruct) CreateHero(id int32) *HeroStruct {

	if Csv.hero.simple_info_map[id] == nil {
		return this
	}

	this.Hero_Info.Hero_id = id //id

	this.Hero_Info.Hero_uid = GetUid() //uid

	this.Hero_Info.Level = 1 //leve

	index := Csv.hero.index_value["102"]
	this.Hero_Info.Type = int32(Csv.hero.simple_info_map[id][index]) //角色类型

	index = Csv.hero.index_value["103"]
	this.Hero_Info.Name = Csv.hero_str.simple_info_map[id][index] //英雄名称

	index = Csv.hero.index_value["105"]
	this.Hero_Info.Zodiac = int32(Csv.hero.simple_info_map[id][index]) //属性

	index = Csv.hero.index_value["106"]
	this.Hero_Info.Star_level = int32(Csv.hero.simple_info_map[id][index]) //星级

	index = Csv.hero.index_value["107"]
	this.Hero_Info.Sex = int32(Csv.hero.simple_info_map[id][index]) //sex

	index = Csv.hero.index_value["108"]
	this.Hero_Info.Feature = Csv.hero_str.simple_info_map[id][index] //特征

	//技能
	index = Csv.hero.index_value["114"]
	skill_id := int32(Csv.hero.simple_info_map[id][index])
	this.Hero_Info.Skill = append(this.Hero_Info.Skill, skill_id)

	index = Csv.hero.index_value["115"]
	skill_id = int32(Csv.hero.simple_info_map[id][index])
	this.Hero_Info.Skill = append(this.Hero_Info.Skill, skill_id)

	index = Csv.hero.index_value["116"]
	skill_id = int32(Csv.hero.simple_info_map[id][index])
	this.Hero_Info.Skill = append(this.Hero_Info.Skill, skill_id)

	//基础属性
	this.GetBase(id, 1, this.Hero_Info.Star_level)
	this.Hero_Info.Hp = 0
	this.Hero_Info.Pos_Arena = -1
	this.Hero_Info.Pos_stage = -1

	//阶增加对应属性比例(创建时候都没有)

	fmt.Println(this)
	Log.Info("%d", this)
	return this
}
