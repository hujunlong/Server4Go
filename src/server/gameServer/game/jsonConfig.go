//读取策划配置
package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	PosID     int
	MonsterID int
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
	guaji_event         map[string]GuajiEvent
	guaji_event_monster map[string]GuajiEventMonster
	guaji_quality_equip map[string]GuajiQualityEquip
	guaji_percent_equip map[string]GuajiPercentEquip
	guaji_event_box     map[string]GuajiEventBox
	guaji_boss_info     map[string]GuajiBossInfo
	guaji_kill_boss_con map[string]GuajiKillBossCon
	guaji_reward        map[string]GuajiReward
	guaji_event_qiyu    map[string]GuajiEventQiyu
	guaji_event_player  map[string]GuajiEventPlayer

	//基础杂项
	stage_monster_info  map[string]StageMonsterInfo
	stage_std_reward    map[string]StageStdReward
	stage_rand_reward   map[string]StageRandReward
	stage_equip_reward  map[string]StageEquipReward
	stage_equip_quality map[string]StageEquipQuality

	//奇遇
	qiyu_boss_jisha    map[string]QiyuBossJiSha
	qiyu_boss_shanghai map[string]QiyuBossShanghai
	qiyu_boss_wancheng map[string]QiyuBossWanCheng
}

func (this *JsonConfig) Init() {

	//挂机
	var name_json string = "config/json/guaji_event.json"
	this.guaji_event = make(map[string]GuajiEvent)
	dealFileJson(&this.guaji_event, name_json)

	name_json = "config/json/guaji_event_monster.json"
	this.guaji_event_monster = make(map[string]GuajiEventMonster)
	dealFileJson(&this.guaji_event_monster, name_json)

	name_json = "config/json/guaji_quality_equip.json"
	this.guaji_quality_equip = make(map[string]GuajiQualityEquip)
	dealFileJson(&this.guaji_quality_equip, name_json)

	name_json = "config/json/guaji_percent_equip.json"
	this.guaji_percent_equip = make(map[string]GuajiPercentEquip)
	dealFileJson(&this.guaji_percent_equip, name_json)

	name_json = "config/json/guaji_event_box.json"
	this.guaji_event_box = make(map[string]GuajiEventBox)
	dealFileJson(&this.guaji_event_box, name_json)

	name_json = "config/json/guaji_boss_info.json"
	this.guaji_boss_info = make(map[string]GuajiBossInfo)
	dealFileJson(&this.guaji_boss_info, name_json)

	name_json = "config/json/guaji_killboss_con.json"
	this.guaji_kill_boss_con = make(map[string]GuajiKillBossCon)
	dealFileJson(&this.guaji_kill_boss_con, name_json)

	name_json = "config/json/guaji_reward.json"
	this.guaji_reward = make(map[string]GuajiReward)
	dealFileJson(&this.guaji_reward, name_json)

	name_json = "config/json/guaji_event_qiyu.json"
	this.guaji_event_qiyu = make(map[string]GuajiEventQiyu)
	dealFileJson(&this.guaji_event_qiyu, name_json)

	name_json = "config/json/guaji_event_player.json"
	this.guaji_event_player = make(map[string]GuajiEventPlayer)
	dealFileJson(&this.guaji_event_player, name_json)

	//基础
	name_json = "config/json/stage_monsterInfo.json"
	this.stage_monster_info = make(map[string]StageMonsterInfo)
	dealFileJson(&this.stage_monster_info, name_json)

	name_json = "config/json/stage_stdReward.json"
	this.stage_std_reward = make(map[string]StageStdReward)
	dealFileJson(&this.stage_std_reward, name_json)

	name_json = "config/json/stage_randReward.json"
	this.stage_rand_reward = make(map[string]StageRandReward)
	dealFileJson(&this.stage_rand_reward, name_json)

	name_json = "config/json/stage_equip_reward.json"
	this.stage_equip_reward = make(map[string]StageEquipReward)
	dealFileJson(&this.stage_equip_reward, name_json)

	name_json = "config/json/stage_equip_quality.json"
	this.stage_equip_quality = make(map[string]StageEquipQuality)
	dealFileJson(&this.stage_equip_quality, name_json)

	//奇遇
	name_json = "config/json/qiyu_boss_jisha.json"
	this.qiyu_boss_jisha = make(map[string]QiyuBossJiSha)
	dealFileJson(&this.qiyu_boss_jisha, name_json)

	name_json = "config/json/qiyu_boss_shanghai.json"
	this.qiyu_boss_shanghai = make(map[string]QiyuBossShanghai)
	dealFileJson(&this.qiyu_boss_shanghai, name_json)

	name_json = "config/json/qiyu_boss_wancheng.json"
	this.qiyu_boss_wancheng = make(map[string]QiyuBossWanCheng)
	dealFileJson(&this.qiyu_boss_wancheng, name_json)

}
