//读取策划配置
package game

import (
	"encoding/csv"
	//"fmt"
	"os"
	"strconv"
	"strings"
)

type Data struct {
	Id  int32
	Num int32
}

//任务系统
type TaskData struct {
	Sub_id int32 //任务子id
	Par1   int32 //参数1
	Par2   int32 //参数2
}

type RewardData struct {
	Type        int32 //物品类型
	Data_reward Data  //物品
}

//任务系统中的物品奖励
func DealGoodsData(str string) []RewardData {
	str = strings.TrimSpace(str)

	var datas []RewardData
	str2 := strings.Split(str, " ")

	for _, str3 := range str2 {
		str4 := strings.Split(str3, "|")
		var data RewardData
		len_ := len(str4)

		if len_ > 0 {
			if len_ == 3 {
				data.Type = Str2Int32(str4[0])
				data.Data_reward.Id = Str2Int32(str4[1])
				data.Data_reward.Num = Str2Int32(str4[2])
			}
		}
		if data.Type > 0 {
			datas = append(datas, data)
		}
	}
	return datas
}

//处理任务系统
func DealTaskData(str string) []TaskData {
	str = strings.TrimSpace(str)

	var datas []TaskData
	str2 := strings.Split(str, " ")

	for _, str3 := range str2 {
		str4 := strings.Split(str3, "|")
		var data TaskData
		len_ := len(str4)

		if len_ > 0 {
			if len_ == 1 {
				data.Sub_id = Str2Int32(str4[0])
				data.Par1 = 0
				data.Par2 = 0
			}

			if len_ == 2 {
				data.Sub_id = Str2Int32(str4[0])
				data.Par1 = Str2Int32(str4[1]) //数量
				data.Par2 = 0                  //附加参数
			}

			if len_ == 3 {
				data.Sub_id = Str2Int32(str4[0])
				data.Par1 = Str2Int32(str4[1])
				data.Par2 = Str2Int32(str4[2])
			}
		}
		if data.Sub_id > 0 {
			datas = append(datas, data)
		}
	}
	return datas
}

func DealStr2Array(str string) []Data { //处理类型 (30002|1 3023|2) //[]Data
	str = strings.TrimSpace(str)
	var datas []Data
	str2 := strings.Split(str, " ")

	for _, str3 := range str2 {
		str4 := strings.Split(str3, "|")
		var data Data
		if len(str4) == 2 {
			data.Id = Str2Int32(str4[0])
			data.Num = Str2Int32(str4[1])
		}
		datas = append(datas, data)
	}
	return datas
}

//处理单成
func DealStr2Int32Array(str string) []int32 { //处理类型 400|100
	str = strings.TrimSpace(str)
	str2 := strings.Split(str, "|")
	var datas []int32

	for _, str3 := range str2 {
		data := Str2Int32(str3)
		datas = append(datas, data)
	}

	return datas
}

func DealStr2Float32Array(str string) []float32 { //处理类型 400|100
	str = strings.TrimSpace(str)
	str2 := strings.Split(str, "|")
	var datas []float32

	for _, str3 := range str2 {
		data := Str2float32(str3)
		datas = append(datas, data)
	}

	return datas
}

var nums []int32 //品质索引
func sortRank(num int32) {
	if len(nums) == 0 {
		nums = append(nums, num)
		return
	}

	for i, v := range nums {
		if num == v {
			return
		}

		if num < v {
			result := []int32{}
			result = append(result, nums[:i]...)
			result = append(result, num)
			result = append(result, nums[i:]...)
			nums = result
			return
		}
	}
	nums = append(nums, num)
}

//单map
type SimpleInfoStr struct {
	index_value     map[string]int
	simple_info_map map[int32][]string
}

//返回一维数组
func readConfig_Simple(path string, key_index_str string) (map[string]int, map[int32][]string) {
	rfile, err_path := os.Open(path)
	defer rfile.Close()
	if err_path != nil {
		writeInfo("readConfig_Simple read file err,path=" + path)
	}

	r := csv.NewReader(rfile)
	strs, err := r.ReadAll()

	if err != nil {
		writeInfo("readConfigSimple config err:" + path + " not found")
		return nil, nil
	}

	//检测key_index_str是否合法
	var is_find bool = false
	for _, v := range strs[0] {
		if strings.EqualFold(v, key_index_str) {
			is_find = true
			continue
		}
	}
	if !is_find {
		return nil, nil
	}

	//第一行存入 key index 方便找位置
	index_value := make(map[string]int)
	for i, v := range strs[0] {
		index_value[v] = i
	}

	//存入内容
	key_index := index_value[key_index_str]
	buf_map := make(map[int32][]string)
	for _, v := range strs[1:] {
		key := v[key_index]
		key_int, _ := strconv.Atoi(key)
		buf_map[int32(key_int)] = v
	}

	return index_value, buf_map
}

//创建角色
type Create_Role_Struct struct {
	Id_104 string
}

func Get_Create_Role(simple *SimpleInfoStr) map[int32]Create_Role_Struct {
	create_roles := make(map[int32]Create_Role_Struct)
	for key, v := range simple.simple_info_map {
		var create_role Create_Role_Struct
		index := simple.index_value["104"]
		create_role.Id_104 = v[index]
		create_roles[key] = create_role
	}
	return create_roles
}

//英雄表
type Hero_Struct struct {
	Id_102   int32
	Id_103   string
	Id_105   int32
	Id_106   int32
	Id_107   int32
	Id_108   string
	Id_130   float32
	Id_131   float32
	Id_132   float32
	Id_133   float32
	Id_134   float32
	Id_109   float32
	Id_110   float32
	Id_111   float32
	Id_112   float32
	Id_113   float32
	Id_skill []int32
	Id_122   int32
	Id_119   int32
}

func GetHeroStruct(simple *SimpleInfoStr) map[int32]Hero_Struct {
	heros := make(map[int32]Hero_Struct)
	for key, v := range simple.simple_info_map {
		var hero Hero_Struct
		index := simple.index_value["102"]
		hero.Id_102 = Str2Int32(v[index])

		index = simple.index_value["103"]
		hero.Id_103 = v[index]

		index = simple.index_value["105"]
		hero.Id_105 = Str2Int32(v[index])

		index = simple.index_value["106"]
		hero.Id_106 = Str2Int32(v[index])

		index = simple.index_value["107"]
		hero.Id_107 = Str2Int32(v[index])

		index = simple.index_value["108"]
		hero.Id_108 = v[index]

		index = simple.index_value["130"]
		hero.Id_130 = Str2float32(v[index])

		index = simple.index_value["131"]
		hero.Id_131 = Str2float32(v[index])

		index = simple.index_value["132"]
		hero.Id_132 = Str2float32(v[index])

		index = simple.index_value["133"]
		hero.Id_133 = Str2float32(v[index])

		index = simple.index_value["134"]
		hero.Id_134 = Str2float32(v[index])

		index = simple.index_value["109"]
		hero.Id_109 = Str2float32(v[index])

		index = simple.index_value["110"]
		hero.Id_110 = Str2float32(v[index])

		index = simple.index_value["111"]
		hero.Id_111 = Str2float32(v[index])

		index = simple.index_value["112"]
		hero.Id_112 = Str2float32(v[index])

		index = simple.index_value["113"]
		hero.Id_113 = Str2float32(v[index])

		//skill
		index = simple.index_value["114"]
		hero.Id_skill = append(hero.Id_skill, Str2Int32(v[index]))

		index = simple.index_value["115"]
		hero.Id_skill = append(hero.Id_skill, Str2Int32(v[index]))

		index = simple.index_value["116"]
		hero.Id_skill = append(hero.Id_skill, Str2Int32(v[index]))

		index = simple.index_value["117"]
		hero.Id_skill = append(hero.Id_skill, Str2Int32(v[index]))

		index = simple.index_value["118"]
		hero.Id_skill = append(hero.Id_skill, Str2Int32(v[index]))

		index = simple.index_value["122"]
		hero.Id_122 = Str2Int32(v[index])

		index = simple.index_value["119"]
		hero.Id_119 = Str2Int32(v[index])

		heros[key] = hero
	}
	return heros
}

//杂项表
type Property_Struct struct {
	Id_102 float32
}

func GetPropertyStruct(simple *SimpleInfoStr) map[int32]Property_Struct {
	propertys := make(map[int32]Property_Struct)
	for key, v := range simple.simple_info_map {
		var property Property_Struct
		index := simple.index_value["102"]
		property.Id_102 = Str2float32(v[index])
		propertys[key] = property
	}
	return propertys
}

//关卡
type Map_Stage_Struct struct {
	Id_109 int32
	Id_114 int32
	Id_117 int32
}

func GetMapStageStruct(simple *SimpleInfoStr) map[int32]Map_Stage_Struct {
	map_stages := make(map[int32]Map_Stage_Struct)
	for key, v := range simple.simple_info_map {
		var map_stage Map_Stage_Struct
		index := simple.index_value["109"]
		map_stage.Id_109 = Str2Int32(v[index])

		index = simple.index_value["114"]
		map_stage.Id_114 = Str2Int32(v[index])

		index = simple.index_value["117"]
		map_stage.Id_117 = Str2Int32(v[index])

		map_stages[key] = map_stage
	}
	return map_stages
}

//玩家经验
type Role_Exp_Struct struct {
	Id_102 int32
}

func GetRoleExpStruct(simple *SimpleInfoStr) map[int32]Role_Exp_Struct {
	role_exps := make(map[int32]Role_Exp_Struct)
	for key, v := range simple.simple_info_map {
		var role_exp Role_Exp_Struct
		index := simple.index_value["102"]
		role_exp.Id_102 = Str2Int32(v[index])
		role_exps[key] = role_exp
	}
	return role_exps
}

//英雄经验
type Hero_Exp_Struct struct {
	Id_102 int32
}

func GetHeroExpStruct(simple *SimpleInfoStr) map[int32]Hero_Exp_Struct {
	hero_exps := make(map[int32]Hero_Exp_Struct)
	for key, v := range simple.simple_info_map {
		var role_exp Hero_Exp_Struct
		index := simple.index_value["102"]
		role_exp.Id_102 = Str2Int32(v[index])
		hero_exps[key] = role_exp
	}
	return hero_exps
}

//挂机关卡
type Map_Guaji_Struct struct {
	Id_102 int32
	Id_105 int32
	Id_106 int32
	Id_116 int32
}

func GetMapGuajiStruct(simple *SimpleInfoStr) map[int32]Map_Guaji_Struct {
	map_guajis := make(map[int32]Map_Guaji_Struct)
	for key, v := range simple.simple_info_map {
		var map_guaji Map_Guaji_Struct
		index := simple.index_value["102"]
		map_guaji.Id_102 = Str2Int32(v[index])

		index = simple.index_value["105"]
		map_guaji.Id_105 = Str2Int32(v[index])

		index = simple.index_value["106"]
		map_guaji.Id_106 = Str2Int32(v[index])

		index = simple.index_value["116"]
		map_guaji.Id_116 = Str2Int32(v[index])

		map_guajis[key] = map_guaji
	}
	return map_guajis
}

//iteam
type Item_Struct struct {
	Id_102 int32
	Id_104 int32
	Id_105 int32
	Id_106 int32
	Id_107 int32
	Id_108 int32
	Id_109 int32
	Id_110 int32
	Id_111 int32
	Id_116 int32
}

func GetItemStruct(simple *SimpleInfoStr) map[int32]Item_Struct {
	items := make(map[int32]Item_Struct)
	for key, v := range simple.simple_info_map {
		var item Item_Struct
		index := simple.index_value["102"]
		item.Id_102 = Str2Int32(v[index])

		index = simple.index_value["104"]
		item.Id_104 = Str2Int32(v[index])

		index = simple.index_value["105"]
		item.Id_105 = Str2Int32(v[index])

		index = simple.index_value["106"]
		item.Id_106 = Str2Int32(v[index])

		index = simple.index_value["107"]
		item.Id_107 = Str2Int32(v[index])

		index = simple.index_value["108"]
		item.Id_108 = Str2Int32(v[index])

		index = simple.index_value["109"]
		item.Id_109 = Str2Int32(v[index])

		index = simple.index_value["110"]
		item.Id_110 = Str2Int32(v[index])

		index = simple.index_value["111"]
		item.Id_111 = Str2Int32(v[index])

		index = simple.index_value["116"]
		item.Id_116 = Str2Int32(v[index])

		items[key] = item
	}
	return items
}

//怪物
type Monster_Struct struct {
	Id_102 int32
	Id_123 string
	Id_122 int32
	Id_104 int32
	Id_105 int32
	Id_106 float32
	Id_107 float32
	Id_108 float32
	Id_109 float32
	Id_110 float32
	Id_111 int32
	Id_112 float32
	Id_113 float32
	Id_114 float32
	Id_115 float32
	Id_116 float32
	Id_117 float32
	Id_118 float32
	Id_119 float32
	Id_120 float32
	Id_121 float32
	Skill  []int32
}

func GetMonsterStruct(simple *SimpleInfoStr) map[int32]Monster_Struct {
	monsters := make(map[int32]Monster_Struct)
	for key, v := range simple.simple_info_map {
		var monster Monster_Struct
		index := simple.index_value["102"]
		monster.Id_102 = Str2Int32(v[index])

		index = simple.index_value["123"]
		monster.Id_123 = v[index]

		index = simple.index_value["122"]
		monster.Id_122 = Str2Int32(v[index])

		index = simple.index_value["104"]
		monster.Id_104 = Str2Int32(v[index])

		index = simple.index_value["105"]
		monster.Id_105 = Str2Int32(v[index])

		index = simple.index_value["106"]
		monster.Id_106 = Str2float32(v[index])

		index = simple.index_value["107"]
		monster.Id_107 = Str2float32(v[index])

		index = simple.index_value["108"]
		monster.Id_108 = Str2float32(v[index])

		index = simple.index_value["109"]
		monster.Id_109 = Str2float32(v[index])

		index = simple.index_value["110"]
		monster.Id_110 = Str2float32(v[index])

		index = simple.index_value["111"]
		monster.Id_111 = Str2Int32(v[index])

		index = simple.index_value["112"]
		monster.Id_112 = Str2float32(v[index])

		index = simple.index_value["113"]
		monster.Id_113 = Str2float32(v[index])

		index = simple.index_value["114"]
		monster.Id_114 = Str2float32(v[index])

		index = simple.index_value["115"]
		monster.Id_115 = Str2float32(v[index])

		index = simple.index_value["116"]
		monster.Id_116 = Str2float32(v[index])

		index = simple.index_value["117"]
		monster.Id_117 = Str2float32(v[index])

		index = simple.index_value["118"]
		monster.Id_118 = Str2float32(v[index])

		index = simple.index_value["119"]
		monster.Id_119 = Str2float32(v[index])

		index = simple.index_value["120"]
		monster.Id_120 = Str2float32(v[index])

		index = simple.index_value["121"]
		monster.Id_121 = Str2float32(v[index])

		monsters[key] = monster
	}
	return monsters
}

//怪物的ability
type Ability_Struct struct {
	Id_102 float32
	Id_103 float32
	Id_104 float32
	Id_105 float32
	Id_106 float32
}

func GetAbilityStruct(simple *SimpleInfoStr) map[int32]Ability_Struct {
	abilitys := make(map[int32]Ability_Struct)
	for key, v := range simple.simple_info_map {
		var ability Ability_Struct
		index := simple.index_value["102"]
		ability.Id_102 = Str2float32(v[index])

		index = simple.index_value["103"]
		ability.Id_103 = Str2float32(v[index])

		index = simple.index_value["104"]
		ability.Id_104 = Str2float32(v[index])

		index = simple.index_value["105"]
		ability.Id_105 = Str2float32(v[index])

		index = simple.index_value["106"]
		ability.Id_106 = Str2float32(v[index])

		abilitys[key] = ability
	}
	return abilitys
}

//英雄进化
type Jinhua_Struct struct {
	Id_102    int32
	Id_104    int32
	NeedEquip []Data
	Id_106    int32
	Id_107    float32
}

func GetJinhuaStruct(simple *SimpleInfoStr) (map[int32]Jinhua_Struct, map[int32]Jinhua_Struct) {
	jinhuas_hero := make(map[int32]Jinhua_Struct)
	jinhuas_guaiwu := make(map[int32]Jinhua_Struct)

	for _, v := range simple.simple_info_map {
		var jinhua Jinhua_Struct
		index := simple.index_value["102"]
		jinhua.Id_102 = Str2Int32(v[index])

		index = simple.index_value["104"]
		jinhua.Id_104 = Str2Int32(v[index])

		index = simple.index_value["105"]
		data := DealStr2Array(v[index])
		jinhua.NeedEquip = append(jinhua.NeedEquip, data...)

		index = simple.index_value["106"]
		jinhua.Id_106 = Str2Int32(v[index])

		index = simple.index_value["107"]
		jinhua.Id_107 = Str2float32(v[index])

		index = simple.index_value["201"]
		if Str2Int32(v[index]) == 1 {
			jinhuas_hero[jinhua.Id_102] = jinhua
		} else {
			jinhuas_guaiwu[jinhua.Id_102] = jinhua
		}
	}
	return jinhuas_hero, jinhuas_guaiwu
}

//英雄升星
type Hero_Star_Struct struct {
	Id_102 int32
	Id_103 float32
}

func GetHeroStarStruct(simple *SimpleInfoStr) map[int32]Hero_Star_Struct {
	hero_stars := make(map[int32]Hero_Star_Struct)
	for key, v := range simple.simple_info_map {
		var hero_star Hero_Star_Struct

		index := simple.index_value["102"]
		hero_star.Id_102 = Str2Int32(v[index])

		index = simple.index_value["103"]
		hero_star.Id_103 = Str2float32(v[index])

		hero_stars[key] = hero_star
	}
	return hero_stars
}

//天赋
type Role_Gift_Struct struct {
	Id_102 int32
	Id_103 int32
	Id_104 int32
	Id_105 int32
}

func GetRoleGiftStruct(simple *SimpleInfoStr) map[int32]Role_Gift_Struct {
	role_gifts := make(map[int32]Role_Gift_Struct)
	for key, v := range simple.simple_info_map {
		var role_gift Role_Gift_Struct

		index := simple.index_value["102"]
		role_gift.Id_102 = Str2Int32(v[index])

		index = simple.index_value["103"]
		role_gift.Id_103 = Str2Int32(v[index])

		index = simple.index_value["104"]
		role_gift.Id_104 = Str2Int32(v[index])

		index = simple.index_value["105"]
		role_gift.Id_105 = Str2Int32(v[index])

		role_gifts[key] = role_gift
	}
	return role_gifts
}

//装备品质
type Equip_Quality_Struct struct {
	Id_101  int32
	Id_103  int32
	Id_104  float32
	Id_105  float32
	Id_106  float32
	Id_duan []float32
	Id_hole []float32
	Id_117  int32
}

func GetEquipQualityStruct(simple *SimpleInfoStr) map[int32][]Equip_Quality_Struct {
	nums = nil
	map_equip_qualitys := make(map[int32][]Equip_Quality_Struct)

	for _, v := range simple.simple_info_map { //查找出所有品质索引
		index := simple.index_value["102"]
		sortRank(Str2Int32(v[index]))
	}

	for _, key := range nums { //按品质索引设置map
		var equip_qualitys []Equip_Quality_Struct

		for _, v := range simple.simple_info_map {
			quan_num := []float32{} //权值数组
			hole_num := []float32{} //孔数组

			index := simple.index_value["102"]
			if Str2Int32(v[index]) == key {

				var equip_quality Equip_Quality_Struct

				index = simple.index_value["101"]
				equip_quality.Id_101 = Str2Int32(v[index])

				index = simple.index_value["103"]
				equip_quality.Id_103 = Str2Int32(v[index])

				index = simple.index_value["104"]
				equip_quality.Id_104 = Str2float32(v[index])

				index = simple.index_value["105"]
				equip_quality.Id_105 = Str2float32(v[index])

				index = simple.index_value["106"]
				equip_quality.Id_106 = Str2float32(v[index])

				//权值数组
				index = simple.index_value["107"]
				quan_num = append(quan_num, Str2float32(v[index]))

				index = simple.index_value["108"]
				quan_num = append(quan_num, Str2float32(v[index]))

				index = simple.index_value["109"]
				quan_num = append(quan_num, Str2float32(v[index]))

				index = simple.index_value["110"]
				quan_num = append(quan_num, Str2float32(v[index]))

				index = simple.index_value["111"]
				quan_num = append(quan_num, Str2float32(v[index]))

				equip_quality.Id_duan = append(equip_quality.Id_duan, quan_num...)
				//孔数组
				index = simple.index_value["112"]
				hole_num = append(hole_num, Str2float32(v[index]))

				index = simple.index_value["113"]
				hole_num = append(hole_num, Str2float32(v[index]))

				index = simple.index_value["114"]
				hole_num = append(hole_num, Str2float32(v[index]))

				index = simple.index_value["115"]
				hole_num = append(hole_num, Str2float32(v[index]))

				index = simple.index_value["116"]
				hole_num = append(hole_num, Str2float32(v[index]))

				index = simple.index_value["118"]
				hole_num = append(hole_num, Str2float32(v[index]))
				equip_quality.Id_hole = append(equip_quality.Id_hole, hole_num...)

				index = simple.index_value["117"]
				equip_quality.Id_117 = Str2Int32(v[index])

				equip_qualitys = append(equip_qualitys, equip_quality)

			}
		}

		map_equip_qualitys[key] = equip_qualitys

	}

	return map_equip_qualitys
}

//装备
type Equip_Struct struct {
	Id_104 int32
	Id_105 int32
	Id_106 int32
	Id_107 float32
}

func GetEquipStruct(simple *SimpleInfoStr) map[int32]Equip_Struct {
	equips := make(map[int32]Equip_Struct)

	for key, v := range simple.simple_info_map {
		var equip Equip_Struct

		index := simple.index_value["104"]
		equip.Id_104 = Str2Int32(v[index])

		index = simple.index_value["105"]
		equip.Id_105 = Str2Int32(v[index])

		index = simple.index_value["106"]
		equip.Id_106 = Str2Int32(v[index])

		index = simple.index_value["107"]
		equip.Id_107 = Str2float32(v[index])

		equips[key] = equip
	}
	return equips
}

//装备精炼
type Equip_Jinglian_Struct struct {
	Id_102        int32
	Id_103        int32
	NeedEquip     []Data
	Jinglian_Quan []int32
	Id_106        int32
	Id_107        float32
	Id_108        int32
	Id_109        float32
	Id_110        int32
	Id_111        float32
	Id_112        int32
	Id_113        float32
	Id_114        int32
	Id_115        float32
}

func GetEquipJinglianStruct(simple *SimpleInfoStr) map[int32]Equip_Jinglian_Struct {
	equip_jinglians := make(map[int32]Equip_Jinglian_Struct)

	for key, v := range simple.simple_info_map {
		var equip_jinglian Equip_Jinglian_Struct

		index := simple.index_value["102"]
		equip_jinglian.Id_102 = Str2Int32(v[index])

		index = simple.index_value["103"]
		equip_jinglian.Id_103 = Str2Int32(v[index])

		index = simple.index_value["104"]
		data := DealStr2Array(v[index])
		equip_jinglian.NeedEquip = append(equip_jinglian.NeedEquip, data...)

		//精炼倍数权值
		index = simple.index_value["201"]
		equip_jinglian.Jinglian_Quan = append(equip_jinglian.Jinglian_Quan, Str2Int32(v[index]))

		index = simple.index_value["202"]
		equip_jinglian.Jinglian_Quan = append(equip_jinglian.Jinglian_Quan, Str2Int32(v[index]))

		index = simple.index_value["203"]
		equip_jinglian.Jinglian_Quan = append(equip_jinglian.Jinglian_Quan, Str2Int32(v[index]))

		index = simple.index_value["204"]
		equip_jinglian.Jinglian_Quan = append(equip_jinglian.Jinglian_Quan, Str2Int32(v[index]))

		index = simple.index_value["106"]
		equip_jinglian.Id_106 = Str2Int32(v[index])

		index = simple.index_value["107"]
		equip_jinglian.Id_107 = Str2float32(v[index])

		index = simple.index_value["108"]
		equip_jinglian.Id_106 = Str2Int32(v[index])

		index = simple.index_value["109"]
		equip_jinglian.Id_109 = Str2float32(v[index])

		index = simple.index_value["110"]
		equip_jinglian.Id_106 = Str2Int32(v[index])

		index = simple.index_value["111"]
		equip_jinglian.Id_111 = Str2float32(v[index])

		index = simple.index_value["112"]
		equip_jinglian.Id_112 = Str2Int32(v[index])

		index = simple.index_value["113"]
		equip_jinglian.Id_113 = Str2float32(v[index])

		index = simple.index_value["114"]
		equip_jinglian.Id_114 = Str2Int32(v[index])

		index = simple.index_value["115"]
		equip_jinglian.Id_115 = Str2float32(v[index])

		equip_jinglians[key] = equip_jinglian
	}
	return equip_jinglians
}

//装备强化
type Equip_Qianghua_Struct struct {
	Id_110   int32
	Id_103   int32
	NeedProp []Data
	Id_106   int32
	Id_107   float32
	Id_108   int32
	Id_109   float32
}

func GetEquipQianghuaStruct(simple *SimpleInfoStr) map[int32][]Equip_Qianghua_Struct {
	nums = nil
	map_equip_qianghuas := make(map[int32][]Equip_Qianghua_Struct)

	for _, v := range simple.simple_info_map { //查找出所有品质索引
		index := simple.index_value["102"]
		sortRank(Str2Int32(v[index]))
	}

	for _, key := range nums {
		var equip_qianghuas []Equip_Qianghua_Struct

		for _, v := range simple.simple_info_map {

			var equip_qianghua Equip_Qianghua_Struct

			index := simple.index_value["102"]
			if Str2Int32(v[index]) != key {
				continue
			}

			index = simple.index_value["110"]
			equip_qianghua.Id_110 = Str2Int32(v[index])

			index = simple.index_value["103"]
			equip_qianghua.Id_103 = Str2Int32(v[index])

			index = simple.index_value["104"]
			data := DealStr2Array(v[index])
			equip_qianghua.NeedProp = append(equip_qianghua.NeedProp, data...)

			index = simple.index_value["106"]
			equip_qianghua.Id_106 = Str2Int32(v[index])

			index = simple.index_value["107"]
			equip_qianghua.Id_107 = Str2float32(v[index])

			index = simple.index_value["108"]
			equip_qianghua.Id_108 = Str2Int32(v[index])

			index = simple.index_value["109"]
			equip_qianghua.Id_109 = Str2float32(v[index])

			equip_qianghuas = append(equip_qianghuas, equip_qianghua)
		}

		map_equip_qianghuas[key] = equip_qianghuas
	}

	return map_equip_qianghuas
}

//装备分解
type Equip_Fenjie_Struct struct {
	Id_103 int32
	Id_104 int32
}

func GetEquipFenjieStruct(simple *SimpleInfoStr) map[int32]Equip_Fenjie_Struct {
	equip_fenjies := make(map[int32]Equip_Fenjie_Struct)

	for key, v := range simple.simple_info_map {
		var equip_fenjie Equip_Fenjie_Struct

		index := simple.index_value["103"]
		equip_fenjie.Id_103 = Str2Int32(v[index])

		index = simple.index_value["104"]
		equip_fenjie.Id_104 = Str2Int32(v[index])

		equip_fenjies[key] = equip_fenjie
	}
	return equip_fenjies
}

//抽卡
type NiuDan_Struct struct {
	One_Need []Data
	Ten_Need []Data
	Id_106   int32
	Id_107   int32
	Id_108   int32
	Id_109   int32
	Id_110   int32
	Id_111   int32
	Id_112   int32
}

func GetNiuDanStruct(simple *SimpleInfoStr) map[int32]NiuDan_Struct {
	niudans := make(map[int32]NiuDan_Struct)

	for key, v := range simple.simple_info_map {
		var niudan NiuDan_Struct

		index := simple.index_value["103"]
		data1 := DealStr2Array(v[index])
		niudan.One_Need = append(niudan.One_Need, data1...)

		index = simple.index_value["104"]
		data2 := DealStr2Array(v[index])
		niudan.Ten_Need = append(niudan.Ten_Need, data2...)

		index = simple.index_value["106"]
		niudan.Id_106 = Str2Int32(v[index])

		index = simple.index_value["107"]
		niudan.Id_107 = Str2Int32(v[index])

		index = simple.index_value["108"]
		niudan.Id_108 = Str2Int32(v[index])

		index = simple.index_value["109"]
		niudan.Id_109 = Str2Int32(v[index])

		index = simple.index_value["110"]
		niudan.Id_110 = Str2Int32(v[index])

		index = simple.index_value["111"]
		niudan.Id_111 = Str2Int32(v[index])

		index = simple.index_value["112"]
		niudan.Id_112 = Str2Int32(v[index])

		niudans[key] = niudan
	}
	return niudans
}

//扭蛋奖品
type NiuDanReward_Struct_Info struct {
	Id_103 int32
	Id_104 int32
	Id_105 int32
	Id_106 int32
}

type NiuDanReward_Struct struct {
	Info    []NiuDanReward_Struct_Info
	QuanZhi []int32
	Total   int32
}

func GetNiuDanRewardStruct(simple *SimpleInfoStr) map[int32]NiuDanReward_Struct {
	nums = nil
	niudan_rewards := make(map[int32]NiuDanReward_Struct)

	for _, v := range simple.simple_info_map {
		index := simple.index_value["102"]
		sortRank(Str2Int32(v[index]))
	}

	for _, num := range nums {
		var niudan_reward NiuDanReward_Struct
		var niudan_rewar_infos []NiuDanReward_Struct_Info
		var total int32 = 0 //用来累加权值
		var quanzhi []int32 //权值列表
		for _, v := range simple.simple_info_map {
			index := simple.index_value["102"]

			if Str2Int32(v[index]) == num {
				var niudan_rewar_info_one NiuDanReward_Struct_Info

				index = simple.index_value["103"]
				niudan_rewar_info_one.Id_103 = Str2Int32(v[index])

				index = simple.index_value["104"]
				niudan_rewar_info_one.Id_104 = Str2Int32(v[index])

				index = simple.index_value["105"]
				niudan_rewar_info_one.Id_105 = Str2Int32(v[index])

				index = simple.index_value["106"]
				niudan_rewar_info_one.Id_106 = Str2Int32(v[index])
				quanzhi = append(quanzhi, niudan_rewar_info_one.Id_106)

				total += niudan_rewar_info_one.Id_106

				niudan_rewar_infos = append(niudan_rewar_infos, niudan_rewar_info_one)
			}
		}
		niudan_reward.Info = niudan_rewar_infos
		niudan_reward.Total = total
		niudan_reward.QuanZhi = quanzhi
		niudan_rewards[num] = niudan_reward
	}
	return niudan_rewards
}

type Quest_Struct struct {
	Id_101     int32
	Id_102     int32
	Id_104     int32
	Id_108     int32
	Id_109     int32
	Id_110     int32
	Task_event []TaskData
	Id_112     int32
	Id_113     int32
	Id_114     int32
	Id_116     int32
	Id_117     int32
	Reward     []RewardData
}

func GetQust(simple *SimpleInfoStr) map[int32]Quest_Struct {

	qust_structs := make(map[int32]Quest_Struct)
	for key, v := range simple.simple_info_map {
		var quest Quest_Struct
		index := simple.index_value["101"]
		quest.Id_101 = Str2Int32(v[index])

		index = simple.index_value["102"]
		quest.Id_101 = Str2Int32(v[index])

		index = simple.index_value["104"]
		quest.Id_104 = Str2Int32(v[index])

		index = simple.index_value["108"]
		quest.Id_108 = Str2Int32(v[index])

		index = simple.index_value["109"]
		quest.Id_109 = Str2Int32(v[index])

		index = simple.index_value["110"]
		quest.Id_110 = Str2Int32(v[index])

		index = simple.index_value["111"]
		task_three_Data := DealTaskData(v[index])
		quest.Task_event = append(quest.Task_event, task_three_Data...)

		index = simple.index_value["112"]
		quest.Id_112 = Str2Int32(v[index])

		index = simple.index_value["113"]
		quest.Id_113 = Str2Int32(v[index])

		index = simple.index_value["114"]
		quest.Id_114 = Str2Int32(v[index])

		index = simple.index_value["116"]
		quest.Id_116 = Str2Int32(v[index])

		index = simple.index_value["117"]
		quest.Id_117 = Str2Int32(v[index])

		index = simple.index_value["115"]
		three_Data := DealGoodsData(v[index])
		quest.Reward = append(quest.Reward, three_Data...)

		qust_structs[key] = quest
	}
	return qust_structs
}

type Quest_xuanshang_info struct {
	Id_101       int32
	Task_quality []Data
	AddQuanzhi   []int32
	Id_106       int32
	Id_110       int32
	Id_111       int32
	Id_112       int32
	Id_113       int32
	Task_event   []TaskData
	Multiple     []float32
	Reward       []RewardData
}

type Quest_xuanshang_struct struct {
	Min_level int32
	Max_level int32
	Xuanshang []Quest_xuanshang_info
}

func GetQustXuanshang(simple *SimpleInfoStr) map[int32]Quest_xuanshang_struct {
	nums = nil
	quest_xuanshang_structs := make(map[int32]Quest_xuanshang_struct)
	for _, v := range simple.simple_info_map { //查找出所有品质索引
		index := simple.index_value["105"]
		sortRank(Str2Int32(v[index]))
	}

	for _, key := range nums {
		var quest_xuanshang_struct Quest_xuanshang_struct
		var quest_xuanshangs []Quest_xuanshang_info
		var min_level int32 = 0
		var max_level int32 = 0
		for _, v := range simple.simple_info_map {
			index := simple.index_value["105"]
			if Str2Int32(v[index]) == key {
				var quest_xuanshang Quest_xuanshang_info //具体任务项

				index = simple.index_value["101"]
				quest_xuanshang.Id_101 = Str2Int32(v[index])

				index = simple.index_value["103"]
				data := DealStr2Array(v[index])
				quest_xuanshang.Task_quality = append(quest_xuanshang.Task_quality, data...) //任务品质和权值

				index = simple.index_value["104"]
				data2 := DealStr2Int32Array(v[index])
				quest_xuanshang.AddQuanzhi = append(quest_xuanshang.AddQuanzhi, data2...) //任务品质和权值

				index = simple.index_value["106"]
				quest_xuanshang.Id_106 = Str2Int32(v[index])

				index = simple.index_value["109"]
				data3 := DealStr2Int32Array(v[index])
				min_level = data3[0]
				max_level = data3[1]

				index = simple.index_value["110"]
				quest_xuanshang.Id_110 = Str2Int32(v[index])

				index = simple.index_value["111"]
				quest_xuanshang.Id_111 = Str2Int32(v[index])

				index = simple.index_value["112"]
				quest_xuanshang.Id_112 = Str2Int32(v[index])

				index = simple.index_value["113"]
				quest_xuanshang.Id_113 = Str2Int32(v[index])

				index = simple.index_value["114"]
				quest_xuanshang.Task_event = DealTaskData(v[index])

				index = simple.index_value["115"]
				quest_xuanshang.Multiple = DealStr2Float32Array(v[index])

				index = simple.index_value["116"]
				quest_xuanshang.Reward = DealGoodsData(v[index])

				quest_xuanshangs = append(quest_xuanshangs, quest_xuanshang)
			}
		}
		quest_xuanshang_struct.Min_level = min_level
		quest_xuanshang_struct.Max_level = max_level
		quest_xuanshang_struct.Xuanshang = append(quest_xuanshang_struct.Xuanshang, quest_xuanshangs...)

		quest_xuanshang_structs[key] = quest_xuanshang_struct
	}
	return quest_xuanshang_structs
}

func GetQustXuanshangArray(simple *SimpleInfoStr) map[int32]Quest_xuanshang_info {
	quest_xuanshang_infos := make(map[int32]Quest_xuanshang_info)

	for key, v := range simple.simple_info_map {
		var quest_xuanshang Quest_xuanshang_info //具体任务项

		index := simple.index_value["103"]
		data := DealStr2Array(v[index])
		quest_xuanshang.Task_quality = append(quest_xuanshang.Task_quality, data...) //任务品质和权值

		index = simple.index_value["104"]
		data2 := DealStr2Int32Array(v[index])
		quest_xuanshang.AddQuanzhi = append(quest_xuanshang.AddQuanzhi, data2...) //任务品质和权值

		index = simple.index_value["106"]
		quest_xuanshang.Id_106 = Str2Int32(v[index])

		index = simple.index_value["110"]
		quest_xuanshang.Id_110 = Str2Int32(v[index])

		index = simple.index_value["111"]
		quest_xuanshang.Id_111 = Str2Int32(v[index])

		index = simple.index_value["112"]
		quest_xuanshang.Id_112 = Str2Int32(v[index])

		index = simple.index_value["113"]
		quest_xuanshang.Id_113 = Str2Int32(v[index])

		index = simple.index_value["114"]
		quest_xuanshang.Task_event = DealTaskData(v[index])

		index = simple.index_value["115"]
		quest_xuanshang.Multiple = DealStr2Float32Array(v[index])

		index = simple.index_value["116"]
		quest_xuanshang.Reward = DealGoodsData(v[index])

		quest_xuanshang_infos[key] = quest_xuanshang
	}
	return quest_xuanshang_infos
}

type ChengjiuStruct struct {
	Id_101 int32
	Id_102 int32
	Id_105 int32
	Id_106 int32
	Id_107 int32
	Id_108 int32
	Id_109 int32
	Reward []RewardData
}

func GetChengjiuStruct(simple *SimpleInfoStr) map[int32]ChengjiuStruct {
	chengjiuStructs := make(map[int32]ChengjiuStruct)

	for key, v := range simple.simple_info_map {
		var chengjiuStruct ChengjiuStruct

		index := simple.index_value["101"]
		data := Str2Int32(v[index])
		chengjiuStruct.Id_101 = data

		index = simple.index_value["102"]
		data = Str2Int32(v[index])
		chengjiuStruct.Id_102 = data

		index = simple.index_value["105"]
		data = Str2Int32(v[index])
		chengjiuStruct.Id_105 = data

		index = simple.index_value["106"]
		data = Str2Int32(v[index])
		chengjiuStruct.Id_106 = data

		index = simple.index_value["107"]
		data = Str2Int32(v[index])
		chengjiuStruct.Id_107 = data

		index = simple.index_value["108"]
		data = Str2Int32(v[index])
		chengjiuStruct.Id_108 = data

		index = simple.index_value["109"]
		data = Str2Int32(v[index])
		chengjiuStruct.Id_109 = data

		index = simple.index_value["110"]
		chengjiuStruct.Reward = DealGoodsData(v[index])

		chengjiuStructs[key] = chengjiuStruct
	}
	return chengjiuStructs
}

func GetChengjiuTypeStruct(simple *SimpleInfoStr) map[int32][]ChengjiuStruct {
	nums = nil
	map_chengjius := make(map[int32][]ChengjiuStruct)

	for _, v := range simple.simple_info_map { //查找出所有品质索引
		index := simple.index_value["105"]
		sortRank(Str2Int32(v[index]))
	}

	for _, key := range nums { //按品质索引设置map
		var chengjius []ChengjiuStruct

		for _, v := range simple.simple_info_map {
			index := simple.index_value["105"]
			if Str2Int32(v[index]) == key {
				var chengjiuStruct ChengjiuStruct

				index := simple.index_value["101"]
				data := Str2Int32(v[index])
				chengjiuStruct.Id_101 = data

				index = simple.index_value["102"]
				data = Str2Int32(v[index])
				chengjiuStruct.Id_102 = data

				index = simple.index_value["105"]
				data = Str2Int32(v[index])
				chengjiuStruct.Id_105 = data

				index = simple.index_value["106"]
				data = Str2Int32(v[index])
				chengjiuStruct.Id_106 = data

				index = simple.index_value["107"]
				data = Str2Int32(v[index])
				chengjiuStruct.Id_107 = data

				index = simple.index_value["108"]
				data = Str2Int32(v[index])
				chengjiuStruct.Id_108 = data

				index = simple.index_value["109"]
				data = Str2Int32(v[index])
				chengjiuStruct.Id_109 = data

				index = simple.index_value["110"]
				chengjiuStruct.Reward = DealGoodsData(v[index])

				chengjius = append(chengjius, chengjiuStruct)
			}
		}
		map_chengjius[key] = chengjius
	}
	return map_chengjius
}

type CsvConfig struct {
	create_role           map[int32]Create_Role_Struct
	hero                  map[int32]Hero_Struct
	property              map[int32]Property_Struct
	role_exp              map[int32]Role_Exp_Struct
	hero_exp              map[int32]Hero_Exp_Struct
	map_stage             map[int32]Map_Stage_Struct
	map_guaji             map[int32]Map_Guaji_Struct
	item                  map[int32]Item_Struct
	monster               map[int32]Monster_Struct
	monsterability        map[int32]Ability_Struct
	hero_jinhua           map[int32]Jinhua_Struct
	hero_jinhua_guaiwu    map[int32]Jinhua_Struct
	hero_star             map[int32]Hero_Star_Struct
	role_gift             map[int32]Role_Gift_Struct
	equip_quality         map[int32][]Equip_Quality_Struct
	equip                 map[int32]Equip_Struct
	equip_jinglian        map[int32]Equip_Jinglian_Struct
	equip_qianghua        map[int32][]Equip_Qianghua_Struct
	equip_fenjie          map[int32]Equip_Fenjie_Struct
	equip_jinglian_leijia map[int32][]Data
	equip_qianghua_leijia map[int32][][]Data
	niudan                map[int32]NiuDan_Struct
	niudan_reward         map[int32]NiuDanReward_Struct
	quest                 map[int32]Quest_Struct
	quest_xuanshang       map[int32]Quest_xuanshang_struct
	quest_xuanshang_array map[int32]Quest_xuanshang_info
	achievement           map[int32]ChengjiuStruct
	achievement_type      map[int32][]ChengjiuStruct
}

func (this *CsvConfig) DealEquipLeijia() {
	this.equip_qianghua_leijia = make(map[int32][][]Data)
	this.equip_jinglian_leijia = make(map[int32][]Data)
	//强化
	var i int32 = 1
	for ; i <= 6; i++ { //类型编号遍历
		v1, _ := this.equip_qianghua[i]
		my_datas := make([][]Data, len(v1))

		for j := 1; j <= len(v1); j++ { //等级从1开始
			var my_data []Data
			for _, v2 := range v1 {
				if int(v2.Id_103) == j {
					if j > 1 {
						my_data = append(my_data, my_datas[j-2]...)
					}
					my_data = append(my_data, v2.NeedProp...)
					break
				}
			}
			my_datas[j-1] = append(my_datas[j-1], my_data...)
		}
		this.equip_qianghua_leijia[i] = my_datas
	}

	//精炼
	var j int = 0
	for ; j < len(this.equip_jinglian); j++ {
		var my_datas []Data
		for k, v := range this.equip_jinglian {
			if int(k) == j {
				if j > 0 {
					my_datas = append(my_datas, this.equip_jinglian_leijia[k-1]...)
				}

				var NeedEquip []Data
				for _, v1 := range v.NeedEquip {
					NeedEquip = append(NeedEquip, Data{v1.Id, v1.Num * v.Id_103})
				}

				my_datas = append(my_datas, NeedEquip...)
				break
			}
		}
		this.equip_jinglian_leijia[int32(j)] = my_datas
	}
}

func (this *CsvConfig) Init() {
	create_role := new(SimpleInfoStr)
	create_role.index_value, create_role.simple_info_map = readConfig_Simple("config/csv/creat_role.csv", "101")
	this.create_role = Get_Create_Role(create_role)

	hero := new(SimpleInfoStr)
	hero.index_value, hero.simple_info_map = readConfig_Simple("config/csv/hero.csv", "101")
	this.hero = GetHeroStruct(hero)

	property := new(SimpleInfoStr)
	property.index_value, property.simple_info_map = readConfig_Simple("config/csv/property.csv", "101")
	this.property = GetPropertyStruct(property)

	map_stage := new(SimpleInfoStr)
	map_stage.index_value, map_stage.simple_info_map = readConfig_Simple("config/csv/map_stage.csv", "101")
	this.map_stage = GetMapStageStruct(map_stage)

	role_exp := new(SimpleInfoStr)
	role_exp.index_value, role_exp.simple_info_map = readConfig_Simple("config/csv/roleexp.csv", "101")
	this.role_exp = GetRoleExpStruct(role_exp)

	hero_exp := new(SimpleInfoStr)
	hero_exp.index_value, hero_exp.simple_info_map = readConfig_Simple("config/csv/hero_exp.csv", "101")
	this.hero_exp = GetHeroExpStruct(hero_exp)

	map_guaji := new(SimpleInfoStr)
	map_guaji.index_value, map_guaji.simple_info_map = readConfig_Simple("config/csv/map_guaji.csv", "101")
	this.map_guaji = GetMapGuajiStruct(map_guaji)

	item := new(SimpleInfoStr)
	item.index_value, item.simple_info_map = readConfig_Simple("config/csv/item.csv", "101")
	this.item = GetItemStruct(item)

	monster := new(SimpleInfoStr)
	monster.index_value, monster.simple_info_map = readConfig_Simple("config/csv/monster.csv", "101")
	this.monster = GetMonsterStruct(monster)

	monsterability := new(SimpleInfoStr)
	monsterability.index_value, monsterability.simple_info_map = readConfig_Simple("config/csv/monsterability.csv", "101")
	this.monsterability = GetAbilityStruct(monsterability)

	hero_jinhua := new(SimpleInfoStr)
	hero_jinhua.index_value, hero_jinhua.simple_info_map = readConfig_Simple("config/csv/hero_jinhua.csv", "101")
	this.hero_jinhua, this.hero_jinhua_guaiwu = GetJinhuaStruct(hero_jinhua)

	hero_star := new(SimpleInfoStr)
	hero_star.index_value, hero_star.simple_info_map = readConfig_Simple("config/csv/hero_star.csv", "101")
	this.hero_star = GetHeroStarStruct(hero_star)

	role_gift := new(SimpleInfoStr)
	role_gift.index_value, role_gift.simple_info_map = readConfig_Simple("config/csv/role_gift.csv", "104")
	this.role_gift = GetRoleGiftStruct(role_gift)

	equip_quality := new(SimpleInfoStr)
	equip_quality.index_value, equip_quality.simple_info_map = readConfig_Simple("config/csv/equip_quality.csv", "101")
	this.equip_quality = GetEquipQualityStruct(equip_quality)

	equip := new(SimpleInfoStr)
	equip.index_value, equip.simple_info_map = readConfig_Simple("config/csv/equip.csv", "101")
	this.equip = GetEquipStruct(equip)

	equip_qianghua := new(SimpleInfoStr)
	equip_qianghua.index_value, equip_qianghua.simple_info_map = readConfig_Simple("config/csv/equip_qianghua.csv", "101")
	this.equip_qianghua = GetEquipQianghuaStruct(equip_qianghua)

	equip_jinglian := new(SimpleInfoStr)
	equip_jinglian.index_value, equip_jinglian.simple_info_map = readConfig_Simple("config/csv/equip_jinglian.csv", "102")
	this.equip_jinglian = GetEquipJinglianStruct(equip_jinglian)

	equip_fenjie := new(SimpleInfoStr)
	equip_fenjie.index_value, equip_fenjie.simple_info_map = readConfig_Simple("config/csv/equip_fenjie.csv", "102")
	this.equip_fenjie = GetEquipFenjieStruct(equip_fenjie)

	this.DealEquipLeijia()

	niudan := new(SimpleInfoStr)
	niudan.index_value, niudan.simple_info_map = readConfig_Simple("config/csv/niudan.csv", "102")
	this.niudan = GetNiuDanStruct(niudan)

	niudan_reward := new(SimpleInfoStr)
	niudan_reward.index_value, niudan_reward.simple_info_map = readConfig_Simple("config/csv/niudan_reward.csv", "101")
	this.niudan_reward = GetNiuDanRewardStruct(niudan_reward)

	quest := new(SimpleInfoStr)
	quest.index_value, quest.simple_info_map = readConfig_Simple("config/csv/quest.csv", "101")
	this.quest = GetQust(quest)

	quest_xuanshang := new(SimpleInfoStr)
	quest_xuanshang.index_value, quest_xuanshang.simple_info_map = readConfig_Simple("config/csv/quest_xuanshang.csv", "101")
	this.quest_xuanshang = GetQustXuanshang(quest_xuanshang)            //分类
	this.quest_xuanshang_array = GetQustXuanshangArray(quest_xuanshang) //数组

	achievement := new(SimpleInfoStr)
	achievement.index_value, achievement.simple_info_map = readConfig_Simple("config/csv/achievement.csv", "101")
	this.achievement = GetChengjiuStruct(achievement)          //数组
	this.achievement_type = GetChengjiuTypeStruct(achievement) //分类

}
