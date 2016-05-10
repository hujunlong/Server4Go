package game

import (
	"fmt"
	"strings"
	"testing"
)

func TestCsv(t *testing.T) {
	csv := new(CsvConfig)
	hero := new(SimpleInfo)

	index_value, simp_buff := readConfig_Simple("config/csv/hero.csv", "101")
	hero.index_value = index_value
	hero.simple_info_map = simp_buff

	csv.hero = hero

	if len(csv.hero.simple_info_map) != 58 {
		t.Error("have 59 len")
	}

	if len(csv.hero.simple_info_map["101"]) > 0 {
	} else {
		t.Error("not found")
	}

	if len(csv.hero.simple_info_map["3"]) == 0 {
	} else {
		t.Error("must len = 0")
	}

	//按照 index_value 查找
	index := hero.index_value["107"]

	if !strings.EqualFold(hero.simple_info_map["40003"][index], "") {
		t.Error("get data error")
	}

	//复合查找测试
	drop_stage := new(CompoundInfo)
	index_value2, Compound_buff := readConfig_Compound("config/csv/drop_stage.csv", "102")
	drop_stage.index_value = index_value2
	drop_stage.compound_info_map = Compound_buff

	fmt.Println("len:", len(drop_stage.compound_info_map))
	if len(drop_stage.compound_info_map) != 31 {
		t.Error("compound buff error")
	}

	if len(drop_stage.compound_info_map["1001"]) == 4 {

	} else {
		t.Error("Compound_buff  error")
	}

	key := drop_stage.index_value["103"]

	value := drop_stage.compound_info_map["1001"][3]

	if strings.EqualFold(value[key], "4") {

	} else {
		t.Error("get data error")
	}

}
