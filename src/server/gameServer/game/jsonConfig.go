//读取策划配置
package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
)

func dealFileJson(inter interface{}, file_name string) {
	bytes, err := ioutil.ReadFile(file_name)
	if err != nil {
		fmt.Println("ReadFile: " + err.Error())
	}

	if err := json.Unmarshal(bytes, &inter); err != nil {
		fmt.Println("Unmarshal: " + err.Error())
	}

}

//挂机
type GuajiEvent_UnMarshal struct {
	Event_type int32
	Per        int32
}

type GuajiEvent struct {
	Item0 []GuajiEvent_UnMarshal
}

type GuajiEventMonster_UnMarshal struct {
	MonModelID int32
	Percent    int32
}

type GuajiEventMonster struct {
	Item0   []GuajiEventMonster_UnMarshal
	Exp_Min int32
	Exp_Max int32
}

type GuajiQualityEquip_UnMarshal struct {
	Quality int32
	Per     int32
}

type GuajiQualityEquip struct {
	Item0 []GuajiQualityEquip_UnMarshal
}

type GuajiPercentEquip_UnMarshal struct {
	EquipID int32
	Per     int32
}

type GuajiPercentEquip struct {
	Item0 []GuajiPercentEquip_UnMarshal
}

type GuajiEventBox_UnMarshal struct {
	ItemType int32
	ItemID   int32
	Num      int32
	Per      int32
}

type GuajiEventBox struct {
	Item0 []GuajiEventBox_UnMarshal
}

type GuajiEventQiyu struct {
	EventBoss         int32
	PercentBoss       int32
	Par_boss          int32
	Type_reward       int32
	Reward_boss       string
	Event_sale        int32
	Percent_sale      int32
	Par_sale          int32
	Event_personQuest int32
	Par_personQuest   int32
	Event_worldQuest  int32
	Par_worldQuest    int32
}

type GuajiBossInfo_UnMarshal struct {
	PosID     int32
	MonsterID int32
}

type GuajiBossInfo struct {
	Item0 []GuajiBossInfo_UnMarshal
}

type GuajiKillBossCon_UnMarshal struct {
	Con int32
	Par int32
}

type GuajiKillBossCon struct {
	Item0 []GuajiKillBossCon_UnMarshal
}

type GuajiReward_UnMarshal_1 struct {
	ItemType int32
	ItemID   int32
	Num      int32
}

type GuajiReward_UnMarshal_2 struct {
	Quality int32
	Per     int32
}

type GuajiReward struct {
	Item0 []GuajiReward_UnMarshal_1
	Item1 []GuajiReward_UnMarshal_2
}

type GuajiEventPlayer_UnMarshal struct {
	HudongType      int32
	Num_qiyu_ziji   int32
	Num_qiyu_difang int32
	ItemType        int32
	ItemID          int32
	Num_ziji        int32
	Num_difang      int32
}

type GuajiEventPlayer struct {
	Item0 []GuajiEventPlayer_UnMarshal
}

//基础杂项
type StageMonsterInfo_UnMarshal struct {
	Pos1_monID int32
	MonsterLv  int32
}
type StageMonsterInfo struct {
	Item0 []StageMonsterInfo_UnMarshal
}

type StageStdReward_UnMarshal struct {
	ItemID int32
	Num    int32
}

type StageStdReward struct {
	Item0 []StageStdReward_UnMarshal
}

type StageRandReward_UnMarshal struct {
	ItemID  int32
	Group   int32
	Num_min int32
	Num_Max int32
	Percent int32
}

type StageRandReward struct {
	Item0 []StageRandReward_UnMarshal
}

type StageEquipReward_UnMarshal struct {
	EquipID int32
	Num     int32
	Percent int32
}

type StageEquipReward struct {
	Item0 []StageEquipReward_UnMarshal
}

type StageEquipQuality_UnMarshal struct {
	Quality int32
	Percent int32
}

type StageEquipQuality struct {
	Item0 []StageEquipQuality_UnMarshal
}

//奇遇
type QiyuBossJiSha_UnMarshal struct {
	ItemType int32
	ItemID   int32
	Num      int32
}

type QiyuBossJiSha struct {
	Item0 []QiyuBossJiSha_UnMarshal
}

type QiyuBossShanghai_UnMarshal struct {
	ItemType int32
	ItemID   int32
	Num      int32
}

type QiyuBossShanghai struct {
	Item0 []QiyuBossShanghai_UnMarshal
}

type QiyuBossWanCheng_UnMarshal struct {
	ItemType int32
	ItemID   int32
	Num      int32
}

type QiyuBossWanCheng struct {
	Item0 []QiyuBossWanCheng_UnMarshal
}

type JsonConfig struct {
	//挂机
	guaji_event         map[int32]GuajiEvent
	guaji_event_monster map[int32]GuajiEventMonster
	guaji_quality_equip map[int32]GuajiQualityEquip
	guaji_percent_equip map[int32]GuajiPercentEquip
	guaji_event_box     map[int32]GuajiEventBox
	guaji_boss_info     map[int32]GuajiBossInfo
	guaji_kill_boss_con map[int32]GuajiKillBossCon
	guaji_reward        map[int32]GuajiReward
	guaji_event_qiyu    map[int32]GuajiEventQiyu
	guaji_event_player  map[int32]GuajiEventPlayer

	//基础杂项
	stage_monster_info  map[int32]StageMonsterInfo
	stage_std_reward    map[int32]StageStdReward
	stage_rand_reward   map[int32]StageRandReward
	stage_equip_reward  map[int32]StageEquipReward
	stage_equip_quality map[int32]StageEquipQuality

	//奇遇
	qiyu_boss_jisha    map[int32]QiyuBossJiSha
	qiyu_boss_shanghai map[int32]QiyuBossShanghai
	qiyu_boss_wancheng map[int32]QiyuBossWanCheng
}

func (this *JsonConfig) changeKey(inter interface{}) {
	type_ := reflect.TypeOf(inter)
	switch type_.String() {
	case "map[string]game.GuajiEvent":
		this.guaji_event = make(map[int32]GuajiEvent)
		for key, v := range inter.(map[string]GuajiEvent) {
			key_int, _ := strconv.Atoi(key)
			this.guaji_event[int32(key_int)] = v
		}
	case "map[string]game.GuajiEventMonster":
		this.guaji_event_monster = make(map[int32]GuajiEventMonster)
		for key, v := range inter.(map[string]GuajiEventMonster) {
			key_int, _ := strconv.Atoi(key)
			this.guaji_event_monster[int32(key_int)] = v
		}
	case "map[string]game.GuajiQualityEquip":
		this.guaji_quality_equip = make(map[int32]GuajiQualityEquip)
		for key, v := range inter.(map[string]GuajiQualityEquip) {
			key_int, _ := strconv.Atoi(key)
			this.guaji_quality_equip[int32(key_int)] = v
		}
	case "map[string]game.GuajiPercentEquip":
		this.guaji_percent_equip = make(map[int32]GuajiPercentEquip)
		for key, v := range inter.(map[string]GuajiPercentEquip) {
			key_int, _ := strconv.Atoi(key)
			this.guaji_percent_equip[int32(key_int)] = v
		}
	case "map[string]game.GuajiEventBox":
		this.guaji_event_box = make(map[int32]GuajiEventBox)
		for key, v := range inter.(map[string]GuajiEventBox) {
			key_int, _ := strconv.Atoi(key)
			this.guaji_event_box[int32(key_int)] = v
		}
	case "map[string]game.GuajiBossInfo":
		this.guaji_boss_info = make(map[int32]GuajiBossInfo)
		for key, v := range inter.(map[string]GuajiBossInfo) {
			key_int, _ := strconv.Atoi(key)
			this.guaji_boss_info[int32(key_int)] = v
		}
	case "map[string]game.GuajiKillBossCon":
		this.guaji_kill_boss_con = make(map[int32]GuajiKillBossCon)
		for key, v := range inter.(map[string]GuajiKillBossCon) {
			key_int, _ := strconv.Atoi(key)
			this.guaji_kill_boss_con[int32(key_int)] = v
		}
	case "map[string]game.GuajiReward":
		this.guaji_reward = make(map[int32]GuajiReward)
		for key, v := range inter.(map[string]GuajiReward) {
			key_int, _ := strconv.Atoi(key)
			this.guaji_reward[int32(key_int)] = v
		}
	case "map[string]game.GuajiEventQiyu":
		this.guaji_event_qiyu = make(map[int32]GuajiEventQiyu)
		for key, v := range inter.(map[string]GuajiEventQiyu) {
			key_int, _ := strconv.Atoi(key)
			this.guaji_event_qiyu[int32(key_int)] = v
		}
	case "map[string]game.GuajiEventPlayer":
		this.guaji_event_player = make(map[int32]GuajiEventPlayer)
		for key, v := range inter.(map[string]GuajiEventPlayer) {
			key_int, _ := strconv.Atoi(key)
			this.guaji_event_player[int32(key_int)] = v
		}
	case "map[string]game.StageMonsterInfo":
		this.stage_monster_info = make(map[int32]StageMonsterInfo)
		for key, v := range inter.(map[string]StageMonsterInfo) {
			key_int, _ := strconv.Atoi(key)
			this.stage_monster_info[int32(key_int)] = v
		}
	case "map[string]game.StageStdReward":
		this.stage_std_reward = make(map[int32]StageStdReward)
		for key, v := range inter.(map[string]StageStdReward) {
			key_int, _ := strconv.Atoi(key)
			this.stage_std_reward[int32(key_int)] = v
		}
	case "map[string]game.StageRandReward":
		this.stage_rand_reward = make(map[int32]StageRandReward)
		for key, v := range inter.(map[string]StageRandReward) {
			key_int, _ := strconv.Atoi(key)
			this.stage_rand_reward[int32(key_int)] = v
		}
	case "map[string]game.StageEquipReward":
		this.stage_equip_reward = make(map[int32]StageEquipReward)
		for key, v := range inter.(map[string]StageEquipReward) {
			key_int, _ := strconv.Atoi(key)
			this.stage_equip_reward[int32(key_int)] = v
		}
	case "map[string]game.StageEquipQuality":
		this.stage_equip_quality = make(map[int32]StageEquipQuality)
		for key, v := range inter.(map[string]StageEquipQuality) {
			key_int, _ := strconv.Atoi(key)
			this.stage_equip_quality[int32(key_int)] = v
		}
	case "map[string]game.QiyuBossJiSha":
		this.qiyu_boss_jisha = make(map[int32]QiyuBossJiSha)
		for key, v := range inter.(map[string]QiyuBossJiSha) {
			key_int, _ := strconv.Atoi(key)
			this.qiyu_boss_jisha[int32(key_int)] = v
		}
	case "map[string]game.QiyuBossShanghai":
		this.qiyu_boss_shanghai = make(map[int32]QiyuBossShanghai)
		for key, v := range inter.(map[string]QiyuBossShanghai) {
			key_int, _ := strconv.Atoi(key)
			this.qiyu_boss_shanghai[int32(key_int)] = v
		}
	case "map[string]game.QiyuBossWanCheng":
		this.qiyu_boss_wancheng = make(map[int32]QiyuBossWanCheng)
		for key, v := range inter.(map[string]QiyuBossWanCheng) {
			key_int, _ := strconv.Atoi(key)
			this.qiyu_boss_wancheng[int32(key_int)] = v
		}
	default:
	}

}

func (this *JsonConfig) Init() {

	//挂机
	var name_json string = "config/json/guaji_event.json"
	guaji_event := make(map[string]GuajiEvent)
	dealFileJson(&guaji_event, name_json)
	this.changeKey(guaji_event)

	name_json = "config/json/guaji_event_monster.json"
	guaji_event_monster := make(map[string]GuajiEventMonster)
	dealFileJson(&guaji_event_monster, name_json)
	this.changeKey(guaji_event_monster)

	name_json = "config/json/guaji_quality_equip.json"
	guaji_quality_equip := make(map[string]GuajiQualityEquip)
	dealFileJson(&guaji_quality_equip, name_json)
	this.changeKey(guaji_quality_equip)

	name_json = "config/json/guaji_percent_equip.json"
	guaji_percent_equip := make(map[string]GuajiPercentEquip)
	dealFileJson(&guaji_percent_equip, name_json)
	this.changeKey(guaji_percent_equip)

	name_json = "config/json/guaji_event_box.json"
	guaji_event_box := make(map[string]GuajiEventBox)
	dealFileJson(&guaji_event_box, name_json)
	this.changeKey(guaji_event_box)

	name_json = "config/json/guaji_boss_info.json"
	guaji_boss_info := make(map[string]GuajiBossInfo)
	dealFileJson(&guaji_boss_info, name_json)
	this.changeKey(guaji_boss_info)

	name_json = "config/json/guaji_killboss_con.json"
	guaji_kill_boss_con := make(map[string]GuajiKillBossCon)
	dealFileJson(&guaji_kill_boss_con, name_json)
	this.changeKey(guaji_kill_boss_con)

	name_json = "config/json/guaji_reward.json"
	guaji_reward := make(map[string]GuajiReward)
	dealFileJson(&guaji_reward, name_json)
	this.changeKey(guaji_reward)

	name_json = "config/json/guaji_event_qiyu.json"
	guaji_event_qiyu := make(map[string]GuajiEventQiyu)
	dealFileJson(&guaji_event_qiyu, name_json)
	this.changeKey(guaji_event_qiyu)

	name_json = "config/json/guaji_event_player.json"
	guaji_event_player := make(map[string]GuajiEventPlayer)
	dealFileJson(&guaji_event_player, name_json)
	this.changeKey(guaji_event_player)

	//基础
	name_json = "config/json/stage_monsterInfo.json"
	stage_monster_info := make(map[string]StageMonsterInfo)
	dealFileJson(&stage_monster_info, name_json)
	this.changeKey(stage_monster_info)

	name_json = "config/json/stage_stdReward.json"
	stage_std_reward := make(map[string]StageStdReward)
	dealFileJson(&stage_std_reward, name_json)
	this.changeKey(stage_std_reward)

	name_json = "config/json/stage_randReward.json"
	stage_rand_reward := make(map[string]StageRandReward)
	dealFileJson(&stage_rand_reward, name_json)
	this.changeKey(stage_rand_reward)

	name_json = "config/json/stage_equip_reward.json"
	stage_equip_reward := make(map[string]StageEquipReward)
	dealFileJson(&stage_equip_reward, name_json)
	this.changeKey(stage_equip_reward)

	name_json = "config/json/stage_equip_quality.json"
	stage_equip_quality := make(map[string]StageEquipQuality)
	dealFileJson(&stage_equip_quality, name_json)
	this.changeKey(stage_equip_quality)

	//奇遇
	name_json = "config/json/qiyu_boss_jisha.json"
	qiyu_boss_jisha := make(map[string]QiyuBossJiSha)
	dealFileJson(&qiyu_boss_jisha, name_json)
	this.changeKey(qiyu_boss_jisha)

	name_json = "config/json/qiyu_boss_shanghai.json"
	qiyu_boss_shanghai := make(map[string]QiyuBossShanghai)
	dealFileJson(&qiyu_boss_shanghai, name_json)
	this.changeKey(qiyu_boss_shanghai)

	name_json = "config/json/qiyu_boss_wancheng.json"
	qiyu_boss_wancheng := make(map[string]QiyuBossWanCheng)
	dealFileJson(&qiyu_boss_wancheng, name_json)
	this.changeKey(qiyu_boss_wancheng)

	/*
		fmt.Println(this.guaji_event)
		fmt.Println(this.guaji_event_monster)
		fmt.Println(this.guaji_quality_equip)
		fmt.Println(this.guaji_percent_equip)
		fmt.Println(this.guaji_event_box)
		fmt.Println(this.guaji_boss_info)
		fmt.Println(this.guaji_kill_boss_con)
		fmt.Println(this.guaji_reward)
		fmt.Println(this.guaji_event_qiyu)
		fmt.Println(this.guaji_event_player)
		fmt.Println(this.stage_monster_info)
		fmt.Println(this.stage_std_reward)
		fmt.Println(this.stage_rand_reward)
		fmt.Println(this.stage_equip_reward)
		fmt.Println(this.stage_equip_quality)
		fmt.Println(this.qiyu_boss_jisha)
		fmt.Println(this.qiyu_boss_shanghai)
		fmt.Println(this.qiyu_boss_wancheng)
	*/
}
