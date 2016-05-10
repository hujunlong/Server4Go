//读取策划配置
package game

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

type SimpleInfo struct { //单map
	index_value     map[string]int
	simple_info_map map[string][]string
}

type CompoundInfo struct { //复合类型
	index_value       map[string]int
	compound_info_map map[string][][]string
}

type CsvConfig struct {
	create_role *SimpleInfo
	hero        *SimpleInfo
	map_stage   *SimpleInfo
	property    *SimpleInfo
	role_exp    *SimpleInfo
	hero_exp    *SimpleInfo
	map_guaji   *SimpleInfo
}

//返回一维数组
func readConfig_Simple(path string, key_index_str string) (map[string]int, map[string][]string) {
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
	buf_map := make(map[string][]string)
	for _, v := range strs {
		key := v[key_index]
		buf_map[key] = v
	}
	return index_value, buf_map
}

//返回二维数组
func readConfig_Compound(path string, key_index_string string) (map[string]int, map[string][][]string) {

	rfile, err_path := os.Open(path)
	defer rfile.Close()
	if err_path != nil {
		writeInfo("readConfig_Compound read file err,path=" + path + "not found")
	}

	r := csv.NewReader(rfile)
	strs, err := r.ReadAll()

	if err != nil {
		writeInfo("readConfigCompound config err:" + path + "not found")
		return nil, nil
	}

	//检测key_index_str是否合法
	var is_find bool = false
	for _, v := range strs[0] {
		if strings.EqualFold(v, key_index_string) {
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

	key_index := index_value[key_index_string]

	//具体值
	config_map := make(map[string][][]string)
	var begin int = 0
	key := strs[0][key_index]

	for i, v := range strs {
		if !strings.EqualFold(key, v[key_index]) {
			config_map[key] = strs[begin:i]
			begin = i
			key = v[key_index]
		}
	}
	//添加最后一个值
	str_len := len(strs)
	config_map[key] = strs[begin:str_len]

	return index_value, config_map
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

func (this *CsvConfig) Init() {

	create_role := new(SimpleInfo)
	hero := new(SimpleInfo)
	map_stage := new(SimpleInfo)
	property := new(SimpleInfo)
	role_exp := new(SimpleInfo)
	hero_exp := new(SimpleInfo)
	map_guaji := new(SimpleInfo)
	create_role.index_value, create_role.simple_info_map = readConfig_Simple("config/csv/creat_role.csv", "101")
	if create_role.index_value == nil || create_role.simple_info_map == nil {
		writeInfo("create_role have erro")
	}
	this.create_role = create_role

	hero.index_value, hero.simple_info_map = readConfig_Simple("config/csv/hero.csv", "101")
	if hero.index_value == nil || hero.simple_info_map == nil {
		writeInfo("hero have erro")
	}
	this.hero = hero

	map_stage.index_value, map_stage.simple_info_map = readConfig_Simple("config/csv/map_stage.csv", "101")
	if map_stage.index_value == nil || map_stage.simple_info_map == nil {
		writeInfo("map_stage have erro")
	}
	this.map_stage = map_stage

	property.index_value, property.simple_info_map = readConfig_Simple("config/csv/property.csv", "101")
	if property.index_value == nil || property.simple_info_map == nil {
		writeInfo("property have erro")
	}
	this.property = property

	role_exp.index_value, role_exp.simple_info_map = readConfig_Simple("config/csv/roleexp.csv", "101")
	if role_exp.index_value == nil || role_exp.simple_info_map == nil {
		writeInfo("property have erro")
	}
	this.role_exp = role_exp

	hero_exp.index_value, hero_exp.simple_info_map = readConfig_Simple("config/csv/hero_exp.csv", "101")
	if hero_exp.index_value == nil || hero_exp.simple_info_map == nil {
		writeInfo("property have erro")
	}
	this.hero_exp = hero_exp

	map_guaji.index_value, map_guaji.simple_info_map = readConfig_Simple("config/csv/map_guaji.csv", "101")
	if map_guaji.index_value == nil || map_guaji.simple_info_map == nil {
		writeInfo("map_guaji have erro")
	}
	this.map_guaji = map_guaji

}
