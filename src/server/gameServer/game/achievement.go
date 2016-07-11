//成就系统
package game

import (
	"fmt"
	"net"
	"server/share/protocol"

	"github.com/golang/protobuf/proto"
)

type AchievementInfo struct {
	Id             int32 //成就id
	Progress       int32 //进度
	Is_complete    bool  //是否完成
	Is_reward      bool  //是否已经领取了奖励
	Total_Progress int32 //总进度
}

type AchievementStruct struct {
	Achievement_Map map[int32][]*AchievementInfo //key:成就id bool:是否领取奖励
}

func (this *AchievementStruct) Init() {
	this.Achievement_Map = make(map[int32][]*AchievementInfo)
}

func (this *AchievementStruct) GetNowData(type_id int32, player *Player) int32 {
	switch type_id {
	case 1: //等级
		return player.Info.Level
	case 2: //登录游戏
		return player.Info.LoginDays
	case 4: //铜钱
		return player.Info.Gold
	default:
	}
	return 0
}

//对外接口
func (this *AchievementStruct) TriggerAchievementEvent(event_id int32, player *Player) { //成就组别

	fmt.Println("成就event_id = ", event_id)
	var type_id int32 = 0
	switch event_id {
	case 3:
		type_id = 1
	}

	if _, ok := Csv.achievement_type[type_id]; !ok {
		return
	}

	if _, ok := this.Achievement_Map[type_id]; !ok {
		fmt.Println("Csv.achievement_type len=", len(Csv.achievement_type[type_id]))
		//添加对应组别的全部内容
		var buff_achievements []*AchievementInfo
		for _, v := range Csv.achievement_type[type_id] {
			buff_achievement := new(AchievementInfo)
			buff_achievement.Is_reward = false
			buff_achievement.Is_complete = false
			buff_achievement.Progress = this.GetNowData(type_id, player)
			buff_achievement.Total_Progress = v.Id_107
			buff_achievement.Id = v.Id_101
			buff_achievements = append(buff_achievements, buff_achievement)
		}

		this.Achievement_Map[type_id] = buff_achievements
	} else {
		switch type_id {
		case 1: //升级
			this.RoleUpLevel(player) //对应等级事件
		case 2: //连续登陆
			this.LoginDays(player) //对应等级事件
		case 4: //铜钱事件
			this.GoldAchievement(player) //对应等级事件
		default:
		}
	}

	if type_id != 0 {
		this.Notice2CAchievementChange(type_id, player.conn)
	}
}

//角色升级
func (this *AchievementStruct) RoleUpLevel(player *Player) {
	for i, _ := range this.Achievement_Map[1] {
		if !this.Achievement_Map[1][i].Is_reward {
			this.Achievement_Map[1][i].Progress = player.Info.Level
			if !this.Achievement_Map[1][i].Is_complete {
				if this.Achievement_Map[1][i].Progress >= this.Achievement_Map[1][i].Total_Progress {
					this.Achievement_Map[1][i].Is_complete = true
				}
			}
		}

	}
}

//登陆游戏次数
func (this *AchievementStruct) LoginDays(player *Player) {
	for i, _ := range this.Achievement_Map[2] {
		this.Achievement_Map[2][i].Progress = player.Info.LoginDays

		if !this.Achievement_Map[2][i].Is_complete {
			if this.Achievement_Map[2][i].Progress >= this.Achievement_Map[2][i].Total_Progress {
				this.Achievement_Map[2][i].Is_complete = true
			}
		}
	}
}

//铜钱成就
func (this *AchievementStruct) GoldAchievement(player *Player) {
	for i, _ := range this.Achievement_Map[4] {
		this.Achievement_Map[4][i].Progress = player.Info.LoginDays

		if !this.Achievement_Map[4][i].Is_complete {
			if this.Achievement_Map[4][i].Progress >= this.Achievement_Map[4][i].Total_Progress {
				this.Achievement_Map[4][i].Is_complete = true
			}
		}
	}
}

//pid 1219 推送成就变化
func (this *AchievementStruct) Notice2CAchievementChange(type_id int32, conn *net.Conn) { //对应
	if _, ok := this.Achievement_Map[type_id]; !ok {
		return
	}

	var infos []*protocol.AchievementInfo
	for _, v_buff := range this.Achievement_Map[type_id] {
		v := v_buff
		info := new(protocol.AchievementInfo)
		info.Id = &v.Id
		info.Progress = &v.Progress
		info.IsComplete = &v.Is_complete
		info.IsReward = &v.Is_reward

		infos = append(infos, info)
	}

	result4C := &protocol.NoticeMsg_Notice2CAchievementChange{
		Infos: infos,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*conn, 1219, encObj)

}

//请求获取完成成就列表
func (this *AchievementStruct) GetAchievementResult(conn *net.Conn) {

	var infos []*protocol.AchievementInfo
	for _, v1 := range this.Achievement_Map {
		for _, v_buff := range v1 {
			v2 := v_buff
			info := new(protocol.AchievementInfo)
			info.Id = &v2.Id
			info.Progress = &v2.Progress
			info.IsComplete = &v2.Is_complete
			info.IsReward = &v2.Is_reward
			infos = append(infos, info)
		}
	}

	result4C := &protocol.Task_GetAchievementResult{
		Infos: infos,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*conn, 1851, encObj)
}

//请求奖励条件检查
func (this *AchievementStruct) GetConditionCheck(ach *AchievementInfo, player *Player) (int32, []*protocol.RwardProp) {
	//0:ok 1:该任务未完成 2:已经领取 3:未查询到该任务id
	var props []*protocol.RwardProp

	if !ach.Is_complete {
		return 1, props
	}

	if ach.Is_reward {
		return 2, props
	}

	for _, v := range Csv.achievement[ach.Id].Reward {
		var prop Prop
		prop.Prop_id = v.Data_reward.Id
		prop.Count = v.Data_reward.Num
		prop.Prop_uid = GetUid()
		player.Bag_Prop.AddAndNotice(prop, player.conn)

		proto_prop := new(protocol.RwardProp)
		proto_prop.PropUid = &prop.Prop_uid
		proto_prop.Num = &prop.Count

		props = append(props, proto_prop)
	}

	ach.Is_reward = true
	return 0, props
}

//请求领取成就奖励
func (this *AchievementStruct) GetAchievementReward(id int32, player *Player) {
	var result int32 = 3
	var reward_props []*protocol.RwardProp

	for _, v1 := range this.Achievement_Map {
		for _, v2 := range v1 {
			if v2.Id == id {
				result, reward_props = this.GetConditionCheck(v2, player)
			}
		}
	}

	result4C := &protocol.Task_GetAchievementRewardResult{
		Result: &result,
		Props:  reward_props,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*player.conn, 1852, encObj)
}
