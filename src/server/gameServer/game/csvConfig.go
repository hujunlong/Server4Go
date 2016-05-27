//读取策划配置
package game

import (
	"encoding/csv"
	//"fmt"
	"os"
	"strconv"
	"strings"
)

type SimpleInfoStr struct { //单map
	index_value     map[string]int
	simple_info_map map[int32][]string
}

type SimpleInfoInt32 struct { //int32
	index_value     map[string]int
	simple_info_map map[int32][]int32
}

type SimpleInfoFloat32 struct { //float32
	index_value     map[string]int
	simple_info_map map[int32][]float32
}

type CsvConfig struct {
	create_role    *SimpleInfoStr
	hero           *SimpleInfoFloat32
	hero_str       *SimpleInfoStr
	property       *SimpleInfoFloat32
	role_exp       *SimpleInfoInt32
	hero_exp       *SimpleInfoInt32
	map_stage      *SimpleInfoInt32
	map_guaji      *SimpleInfoInt32
	item           *SimpleInfoInt32
	monster        *SimpleInfoInt32
	monster_str    *SimpleInfoStr
	monsterability *SimpleInfoInt32
	hero_jinhua    *SimpleInfoInt32
	hero_star      *SimpleInfoInt32
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
	for _, v := range strs {
		key := v[key_index]
		key_int, _ := strconv.Atoi(key)
		buf_map[int32(key_int)] = v
	}

	return index_value, buf_map
}

func Str2int32(str string) int32 {
	if len(str) <= 0 {
		return 0
	}

	a2_i, err := strconv.Atoi(str)
	if err != nil {
		writeInfo("can't to number %s=" + str)
		return 0
	}
	return int32(a2_i)
}

func Str2float32(str string) float32 {
	if len(str) <= 0 {
		return 0
	}

	f, err := strconv.ParseFloat(str, 32)
	if err != nil {
		writeInfo("can't to number %s=" + str)
		return 0
	}
	return float32(f)
}

func readConfig_Simple_int32(path string, key_index_str string) (map[string]int, map[int32][]int32) {
	rfile, _ := os.Open(path)
	defer rfile.Close()

	r := csv.NewReader(rfile)
	strs, err := r.ReadAll()

	if err != nil {
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
	buf_map := make(map[int32][]int32)
	for _, v := range strs {
		key := Str2int32(v[key_index])
		var int32_list []int32

		for _, k := range v {
			k_int32 := Str2int32(k)
			int32_list = append(int32_list, k_int32)
		}
		buf_map[key] = int32_list
	}
	return index_value, buf_map
}

func readConfig_Simple_float32(path string, key_index_str string) (map[string]int, map[int32][]float32) {
	rfile, _ := os.Open(path)
	defer rfile.Close()

	r := csv.NewReader(rfile)
	strs, err := r.ReadAll()

	if err != nil {
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
	buf_map := make(map[int32][]float32)
	for _, v := range strs {
		key := Str2int32(v[key_index])
		var float32_list []float32

		for _, k := range v {
			k_float32 := Str2float32(k)
			float32_list = append(float32_list, k_float32)
		}
		buf_map[key] = float32_list
	}
	return index_value, buf_map
}

func (this *CsvConfig) Init() {
	create_role := new(SimpleInfoStr)
	hero := new(SimpleInfoFloat32)
	hero_str := new(SimpleInfoStr)
	property := new(SimpleInfoFloat32)
	map_stage := new(SimpleInfoInt32)
	role_exp := new(SimpleInfoInt32)
	hero_exp := new(SimpleInfoInt32)
	map_guaji := new(SimpleInfoInt32)
	item := new(SimpleInfoInt32)
	monster := new(SimpleInfoInt32)
	monster_str := new(SimpleInfoStr)
	monsterability := new(SimpleInfoInt32)
	hero_jinhua := new(SimpleInfoInt32)
	hero_star := new(SimpleInfoInt32)
	create_role.index_value, create_role.simple_info_map = readConfig_Simple("config/csv/creat_role.csv", "101")
	if create_role.index_value == nil || create_role.simple_info_map == nil {
		writeInfo("create_role have erro")
	}
	this.create_role = create_role

	hero.index_value, hero.simple_info_map = readConfig_Simple_float32("config/csv/hero.csv", "101")
	if hero.index_value == nil || hero.simple_info_map == nil {
		writeInfo("hero have erro")
	}
	this.hero = hero

	hero_str.index_value, hero_str.simple_info_map = readConfig_Simple("config/csv/hero.csv", "101")
	if hero_str.index_value == nil || hero_str.simple_info_map == nil {
		writeInfo("hero_str have erro")
	}
	this.hero_str = hero_str

	map_stage.index_value, map_stage.simple_info_map = readConfig_Simple_int32("config/csv/map_stage.csv", "101")
	if map_stage.index_value == nil || map_stage.simple_info_map == nil {
		writeInfo("map_stage have erro")
	}
	this.map_stage = map_stage

	property.index_value, property.simple_info_map = readConfig_Simple_float32("config/csv/property.csv", "101")
	if property.index_value == nil || property.simple_info_map == nil {
		writeInfo("property have erro")
	}
	this.property = property

	role_exp.index_value, role_exp.simple_info_map = readConfig_Simple_int32("config/csv/roleexp.csv", "101")
	if role_exp.index_value == nil || role_exp.simple_info_map == nil {
		writeInfo("property have erro")
	}
	this.role_exp = role_exp

	hero_exp.index_value, hero_exp.simple_info_map = readConfig_Simple_int32("config/csv/hero_exp.csv", "101")
	if hero_exp.index_value == nil || hero_exp.simple_info_map == nil {
		writeInfo("property have erro")
	}
	this.hero_exp = hero_exp

	map_guaji.index_value, map_guaji.simple_info_map = readConfig_Simple_int32("config/csv/map_guaji.csv", "101")
	if map_guaji.index_value == nil || map_guaji.simple_info_map == nil {
		writeInfo("map_guaji have erro")
	}
	this.map_guaji = map_guaji

	item.index_value, item.simple_info_map = readConfig_Simple_int32("config/csv/item.csv", "101")
	if item.index_value == nil || item.simple_info_map == nil {
		writeInfo("item have erro")
	}
	this.item = item

	monster.index_value, monster.simple_info_map = readConfig_Simple_int32("config/csv/monster.csv", "101")
	if monster.index_value == nil || monster.simple_info_map == nil {
		writeInfo("monster have erro")
	}
	this.monster = monster

	monster_str.index_value, monster_str.simple_info_map = readConfig_Simple("config/csv/monster.csv", "101")
	if monster_str.index_value == nil || monster_str.simple_info_map == nil {
		writeInfo("monster_str have erro")
	}
	this.monster_str = monster_str

	monsterability.index_value, monsterability.simple_info_map = readConfig_Simple_int32("config/csv/monsterability.csv", "101")
	if monsterability.index_value == nil || monsterability.simple_info_map == nil {
		writeInfo("monsterability have erro")
	}
	this.monsterability = monsterability

	hero_jinhua.index_value, hero_jinhua.simple_info_map = readConfig_Simple_int32("config/csv/hero_jinhua.csv", "101")
	if hero_jinhua.index_value == nil || hero_jinhua.simple_info_map == nil {
		writeInfo("hero_jinhua have erro")
	}
	this.hero_jinhua = hero_jinhua

	hero_star.index_value, hero_star.simple_info_map = readConfig_Simple_int32("config/csv/hero_star.csv", "101")
	if hero_star.index_value == nil || hero_star.simple_info_map == nil {
		writeInfo("hero_star have erro")
	}
	this.hero_star = hero_star

}
