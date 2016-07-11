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

func (this *Monsters) dealMonster2Protocol() []*protocol.FightingAttr {
	var MonsterAttrs []*protocol.FightingAttr
	for _, v_buff := range this.monsters {
		v := v_buff

		//var physical_attack int32 = int32(v.physical_attack)
		//var magic_attack int32 = int32(v.magic_attack)
		//var physical_defense int32 = int32(v.physical_defense)
		//var magic_defense int32 = int32(v.magic_defense)

		monster_attr := new(protocol.FightingAttr)
		monster_attr.Id = &v.id
		monster_attr.Uid = &v.uid
		monster_attr.Pos = &v.pos
		monster_attr.Level = &v.level
		//monster_attr.Hp = &v.hp
		//monster_attr.PhysicalAttack = &physical_attack
		//monster_attr.MagicAttack = &magic_attack
		//monster_attr.PhysicalDefense = &physical_defense
		//monster_attr.MagicDefense = &magic_defense
		//monster_attr.Speed = &v.speed
		monster_attr.StepLevel = &v.step_level
		monster_attr.StarLevel = &v.star_level

		var Attributes []*protocol.Attribute
		for buff_key, buff_v := range v.attrs {
			key := buff_key
			v := buff_v
			attribute := new(protocol.Attribute)
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
	if _, ok := Csv.monster[monster_id]; !ok { //检查配置
		return monster
	}

	monster.attrs = make(map[int32]float32)
	monster.id = monster_id //怪物id
	monster.pos = pos       //怪物位置
	monster.uid = GetUid()  //怪物uid

	monster.level = Csv.monster[monster_id].Id_122 //等级

	monster.step_level = Csv.monster[monster_id].Id_104 //怪物阶数

	monster.star_level = Csv.monster[monster_id].Id_105 //怪物星级

	monster.speed = Csv.monster[monster_id].Id_111 //速度值

	monster.name = Csv.monster[monster_id].Id_123 //怪物名

	data := Csv.property[1013].Id_102
	monster.attrs[1] = Csv.monster[monster_id].Id_112 / data //暴击值

	data = Csv.property[1018].Id_102
	monster.attrs[2] = Csv.monster[monster_id].Id_113 / data //暴伤

	data = Csv.property[1019].Id_102
	monster.attrs[3] = Csv.monster[monster_id].Id_114 / data //连击值

	data = Csv.property[1020].Id_102
	monster.attrs[4] = Csv.monster[monster_id].Id_115 / data //抵抗值

	monster.attrs[5] = Csv.monster[monster_id].Id_116 / 10000 //物理吸血百分比

	monster.attrs[6] = Csv.monster[monster_id].Id_117 / 10000 //法术吸血百分比

	monster.attrs[7] = Csv.monster[monster_id].Id_118 //物理免疫

	monster.attrs[8] = Csv.monster[monster_id].Id_119 //法术免疫

	monster.attrs[9] = Csv.monster[monster_id].Id_120 //霸体

	monster.attrs[10] = Csv.monster[monster_id].Id_121 //负面状态免疫

	//最后这几个值用来

	hp_pre := Csv.monster[monster_id].Id_106 //血量比例

	physical_attack_pre := Csv.monster[monster_id].Id_107 //物理攻击比例

	magic_attack_pre := Csv.monster[monster_id].Id_108 //法术攻击比例

	physical_defense_pre := Csv.monster[monster_id].Id_109 //物理防御比例

	magic_defense_pre := Csv.monster[monster_id].Id_110 //法术防御比例

	//读取monsterbility文件
	monster.hp = int32(Csv.monsterability[monster.level].Id_102 * (hp_pre / 10000))

	//物攻
	monster.physical_attack = Csv.monsterability[monster.level].Id_103 * (physical_attack_pre / 10000)

	//魔攻
	monster.magic_attack = Csv.monsterability[monster.level].Id_104 * (magic_attack_pre / 10000)

	//物防
	monster.physical_defense = Csv.monsterability[monster.level].Id_105 * (physical_defense_pre / 10000)

	//法防
	monster.magic_defense = Csv.monsterability[monster.level].Id_106 * (magic_defense_pre / 10000)

	//阶增加 血量 物攻 法攻 物防 法防
	//对应阶数序列号

	var step_level_index int32 = monster.step_level
	if step_level_index > 0 {
		//加血
		monster.hp += monster.hp * int32(Csv.hero_jinhua_guaiwu[step_level_index].Id_107/10000)

		//加物攻
		monster.physical_attack += monster.physical_attack * (Csv.hero_jinhua_guaiwu[step_level_index].Id_107 / 10000)

		//加法功
		monster.magic_attack += monster.magic_attack * (Csv.hero_jinhua_guaiwu[step_level_index].Id_107 / 10000)

		//物防
		monster.physical_defense += monster.physical_defense * (Csv.hero_jinhua_guaiwu[step_level_index].Id_107 / 10000)

		//法防
		monster.magic_defense += monster.magic_defense * (Csv.hero_jinhua_guaiwu[step_level_index].Id_107 / 10000)
	} else {
		monster.step_level = 0
	}

	//怪物的特征，性别，属相
	//怪物表对应 hero表的id
	monster_2hero := Csv.monster[monster_id].Id_102
	if _, ok := Csv.hero[monster_2hero]; !ok {
		return monster
	}

	//属相
	monster.zodiac = int32(Csv.hero[monster_2hero].Id_105)

	//性别
	monster.sex = int32(Csv.hero[monster_2hero].Id_107)

	//特征
	monster.feature = Csv.hero[monster_2hero].Id_108

	//阶增加 血量 物攻 法攻 物防 法防
	//对应星级
	if monster.star_level > 0 {

		extra_star_data := Csv.hero_star[monster.star_level].Id_103
		monster.physical_attack += monster.physical_attack * (extra_star_data / 10000)
		monster.magic_attack += monster.magic_attack * (extra_star_data / 10000)
		monster.physical_defense += monster.physical_attack * (extra_star_data / 10000)
		monster.magic_defense += monster.physical_attack * (extra_star_data / 10000)
	}

	//叠加后的比例
	pre_hp := Csv.monster[monster_id].Id_106 / 10000
	monster.hp = int32(pre_hp * float32(monster.hp))

	pre_physical := Csv.monster[monster_id].Id_107 / 10000 //物攻击比
	monster.physical_attack = monster.physical_attack * pre_physical

	pre_magic := Csv.monster[monster_id].Id_108 / 10000 //法术比
	monster.magic_attack = monster.magic_attack * pre_magic

	pre_phsical_def := Csv.monster[monster_id].Id_109 / 10000 //物防比
	monster.physical_defense = monster.physical_defense * pre_phsical_def

	pre_magic_def := Csv.monster[monster_id].Id_110 / 10000 //魔防比
	monster.magic_defense = monster.magic_defense * pre_magic_def

	fmt.Println(monster)
	Log.Info("%s", monster)
	return monster
}
