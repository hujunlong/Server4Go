//怪物相关
package game

import (
	"fmt"
	"server/share/protocol"
)

type Monster struct { //hero基础属性
	id               int32             //id
	uid              int32             //uid
	pos              int32             //位置
	level            int32             //等级
	hp               int32             //血量
	step_level       int32             //怪物阶数
	star_level       int32             //怪物星级
	speed            int32             //速度
	zodiac           int32             //属相 金木水火土
	sex              int32             //性别
	physical_attack  float32           //物理攻击
	magic_attack     float32           //法术攻击
	physical_defense float32           //物理防御
	magic_defense    float32           //魔法防御
	feature          string            //特征
	name             string            //名字
	attrs            map[int32]float32 //额外属性
}

type Monsters struct {
	monsters []Monster
}

//获取一个关卡的怪物组数据
func (this *Monsters) GetMonsters(stage_id int32) int32 { //0数据正常 1该关卡不存在
	if _, ok := Json_config.guaji_boss_info[stage_id]; !ok {
		return 1
	}

	for _, v := range Json_config.guaji_boss_info[stage_id].Item0 {
		if v.MonsterID > 0 && v.PosID > 0 {
			monster_ := this.CreateMonster(v.MonsterID, v.PosID)
			this.monsters = append(this.monsters, monster_)
		}
	}
	return 0
}

func (this *Monsters) dealMonster2Protocol() []*protocol.Game_MonsterAttr {
	var MonsterAttrs []*protocol.Game_MonsterAttr
	for _, v_buff := range this.monsters {
		v := v_buff
		monster_attr := new(protocol.Game_MonsterAttr)
		monster_attr.Id = &v.id
		monster_attr.Uid = &v.uid
		monster_attr.Pos = &v.pos
		monster_attr.Level = &v.level
		monster_attr.Hp = &v.hp
		monster_attr.PhysicalAttack = &v.physical_attack
		monster_attr.MagicAttack = &v.magic_attack
		monster_attr.PhysicalDefense = &v.physical_defense
		monster_attr.MagicDefense = &v.magic_defense
		monster_attr.Speed = &v.speed
		monster_attr.StepLevel = &v.step_level
		monster_attr.StarLevel = &v.star_level
		monster_attr.Name = &v.name
		monster_attr.Zodiac = &v.zodiac
		monster_attr.Sex = &v.sex
		monster_attr.Feature = &v.feature

		var Attributes []*protocol.Game_Attribute
		for buff_key, buff_v := range v.attrs {
			key := buff_key
			v := buff_v
			attribute := new(protocol.Game_Attribute)
			attribute.Key = &key
			attribute.Value = &v
			Attributes = append(Attributes, attribute)
		}
		monster_attr.Attribute = Attributes
		MonsterAttrs = append(MonsterAttrs, monster_attr)
	}

	fmt.Println(MonsterAttrs)
	Log.Info("%s", MonsterAttrs)
	return MonsterAttrs
}

//创建一个怪物
func (this *Monsters) CreateMonster(monster_id int32, pos int32) Monster {
	var monster Monster
	if _, ok := Csv.monster.simple_info_map[monster_id]; !ok { //检查配置
		return monster
	}

	monster.attrs = make(map[int32]float32)
	monster.id = monster_id //怪物id
	monster.pos = pos       //怪物位置
	monster.uid = GetUid()  //怪物uid

	key := Csv.monster.index_value["122"]
	monster.level = Csv.monster.simple_info_map[monster_id][key] //等级

	key = Csv.monster.index_value["104"]
	monster.step_level = Csv.monster.simple_info_map[monster_id][key] //怪物阶数

	key = Csv.monster.index_value["105"]
	monster.star_level = Csv.monster.simple_info_map[monster_id][key] //怪物星级

	key = Csv.monster.index_value["111"]
	monster.speed = Csv.monster.simple_info_map[monster_id][key] //速度值

	key = Csv.monster.index_value["123"]
	monster.name = Csv.monster_str.simple_info_map[monster_id][key] //怪物名

	key = Csv.monster.index_value["112"]
	data := Csv.property.simple_info_map[1013][1]
	monster.attrs[1] = float32(Csv.monster.simple_info_map[monster_id][key]) / data //暴击值

	key = Csv.monster.index_value["113"]
	data = Csv.property.simple_info_map[1018][1]
	monster.attrs[2] = float32(Csv.monster.simple_info_map[monster_id][key]) / data //暴伤

	key = Csv.monster.index_value["114"]
	data = Csv.property.simple_info_map[1019][1]
	monster.attrs[3] = float32(Csv.monster.simple_info_map[monster_id][key]) / data //连击值

	key = Csv.monster.index_value["115"]
	data = Csv.property.simple_info_map[1020][1]
	monster.attrs[4] = float32(Csv.monster.simple_info_map[monster_id][key]) / data //抵抗值

	key = Csv.monster.index_value["116"]
	monster.attrs[5] = float32(Csv.monster.simple_info_map[monster_id][key]) / 10000 //物理吸血百分比

	key = Csv.monster.index_value["117"]
	monster.attrs[6] = float32(Csv.monster.simple_info_map[monster_id][key]) / 10000 //法术吸血百分比

	key = Csv.monster.index_value["118"]
	monster.attrs[7] = float32(Csv.monster.simple_info_map[monster_id][key]) //物理免疫

	key = Csv.monster.index_value["119"]
	monster.attrs[8] = float32(Csv.monster.simple_info_map[monster_id][key]) //法术免疫

	key = Csv.monster.index_value["120"]
	monster.attrs[9] = float32(Csv.monster.simple_info_map[monster_id][key]) //霸体

	key = Csv.monster.index_value["121"]
	monster.attrs[10] = float32(Csv.monster.simple_info_map[monster_id][key]) //负面状态免疫

	//最后这几个值用来
	key = Csv.monster.index_value["106"]
	hp_pre := float32(Csv.monster.simple_info_map[monster_id][key]) //血量比例

	key = Csv.monster.index_value["107"]
	physical_attack_pre := float32(Csv.monster.simple_info_map[monster_id][key]) //物理攻击比例

	key = Csv.monster.index_value["108"]
	magic_attack_pre := float32(Csv.monster.simple_info_map[monster_id][key]) //法术攻击比例

	key = Csv.monster.index_value["109"]
	physical_defense_pre := float32(Csv.monster.simple_info_map[monster_id][key]) //物理防御比例

	key = Csv.monster.index_value["110"]
	magic_defense_pre := float32(Csv.monster.simple_info_map[monster_id][key]) //法术防御比例

	//读取monsterbility文件
	key = Csv.monsterability.index_value["102"] //血量
	monster.hp = int32(float32(Csv.monsterability.simple_info_map[monster.level][key]) * (hp_pre / 10000))

	key = Csv.monsterability.index_value["103"] //物攻
	monster.physical_attack = float32(Csv.monsterability.simple_info_map[monster.level][key]) * (physical_attack_pre / 10000)

	key = Csv.monsterability.index_value["104"] //魔攻
	monster.magic_attack = float32(Csv.monsterability.simple_info_map[monster.level][key]) * (magic_attack_pre / 10000)

	key = Csv.monsterability.index_value["105"] //物防
	monster.physical_defense = float32(Csv.monsterability.simple_info_map[monster.level][key]) * (physical_defense_pre / 10000)

	key = Csv.monsterability.index_value["106"] //法防
	monster.magic_defense = float32(Csv.monsterability.simple_info_map[monster.level][key]) * (magic_defense_pre / 10000)

	fmt.Println("2222222", *this)
	//阶增加 血量 物攻 法攻 物防 法防
	//对应阶数序列号
	index_type := Csv.hero_jinhua.index_value["201"]
	index_jie := Csv.hero_jinhua.index_value["102"]

	var step_level_index int32 = 0
	for k, v := range Csv.hero_jinhua.simple_info_map {
		if v[index_type] == 2 && v[index_jie] == monster.step_level {
			step_level_index = k
			break
		}
	}

	if step_level_index > 0 {
		key = Csv.hero_jinhua.index_value["113"] //加物攻
		monster.hp += monster.hp * int32(float32(Csv.hero_jinhua.simple_info_map[step_level_index][key])/10000)

		key = Csv.hero_jinhua.index_value["113"] //加物攻
		monster.physical_attack += monster.physical_attack * (float32(Csv.hero_jinhua.simple_info_map[step_level_index][key]) / 10000)

		key = Csv.hero_jinhua.index_value["114"] //加法功
		monster.magic_attack += monster.magic_attack * (float32(Csv.hero_jinhua.simple_info_map[step_level_index][key]) / 10000)

		key = Csv.hero_jinhua.index_value["115"] //物防
		monster.physical_defense += monster.physical_defense * (float32(Csv.hero_jinhua.simple_info_map[step_level_index][key]) / 10000)

		key = Csv.hero_jinhua.index_value["116"] //法防
		monster.magic_defense += monster.magic_defense * (float32(Csv.hero_jinhua.simple_info_map[step_level_index][key]) / 10000)
	} else {
		monster.step_level = 0
	}

	//怪物的特征，性别，属相
	//怪物表对应 hero表的id
	key = Csv.monster.index_value["102"]
	monster_2hero := Csv.monster.simple_info_map[monster_id][key]
	if _, ok := Csv.hero.simple_info_map[monster_2hero]; !ok {
		return monster
	}

	key = Csv.hero.index_value["105"] //属相
	monster.zodiac = int32(Csv.hero.simple_info_map[monster_2hero][key])

	key = Csv.hero.index_value["107"] //性别
	monster.sex = int32(Csv.hero.simple_info_map[monster_2hero][key])

	key = Csv.hero.index_value["108"] //特征
	monster.feature = Csv.hero_str.simple_info_map[monster_2hero][key]

	//阶增加 血量 物攻 法攻 物防 法防
	//对应星级
	if monster.star_level > 0 {
		key = Csv.hero_star.index_value["103"]
		extra_star_data := float32(Csv.hero_star.simple_info_map[monster.star_level][key])
		monster.physical_attack += monster.physical_attack * (extra_star_data / 10000)
		monster.magic_attack += monster.magic_attack * (extra_star_data / 10000)
		monster.physical_defense += monster.physical_attack * (extra_star_data / 10000)
		monster.magic_defense += monster.physical_attack * (extra_star_data / 10000)
	}

	//叠加后的比例
	key = Csv.monster.index_value["106"]
	pre_hp := float32(Csv.monster.simple_info_map[monster_id][key]) / 10000
	monster.hp = int32(pre_hp * float32(monster.hp))

	key = Csv.monster.index_value["107"]
	pre_physical := float32(Csv.monster.simple_info_map[monster_id][key]) / 10000 //物攻击比
	monster.physical_attack = monster.physical_attack * pre_physical

	key = Csv.monster.index_value["108"]
	pre_magic := float32(Csv.monster.simple_info_map[monster_id][key]) / 10000 //法术比
	monster.magic_attack = monster.magic_attack * pre_magic

	key = Csv.monster.index_value["109"]
	pre_phsical_def := float32(Csv.monster.simple_info_map[monster_id][key]) / 10000 //物防比
	monster.physical_defense = monster.physical_defense * pre_phsical_def

	key = Csv.monster.index_value["110"]
	pre_magic_def := float32(Csv.monster.simple_info_map[monster_id][key]) / 10000 //魔防比
	monster.magic_defense = monster.magic_defense * pre_magic_def

	fmt.Println(monster)
	Log.Info("%s", monster)
	return monster
}
