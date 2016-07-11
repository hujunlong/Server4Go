//任务系统
package game

import (
	"fmt"
	"net"
	"server/share/protocol"
	"time"

	"github.com/golang/protobuf/proto"
)

//任务系统
type TaskInfo struct {
	Sub_id         int32 //任务子id
	Progress       int32 //任务进度
	Total_Progress int32 //总进度
	Par            int32 //对比附加参数
}

type HandleTaskStruct struct {
	Type_id     int32       //任务类型(1:主线 2:悬赏任务 3:奇遇)
	Quality_ref int32       //品质组
	Stage       int32       //状态 1:可以领取 2:已领取 3:完成未领取奖励 4:完成领取奖励
	Task_info   []*TaskInfo //子项任务
}

//悬赏相关
type XuanShangStruct struct {
	Ref_Total_num int32 //悬赏元宝刷新总次数
	Last_time     int32 //上次免费刷新时间
}

type TaskStruct struct {
	Tasks       map[int32]*HandleTaskStruct //任务列表
	XuanShang   *XuanShangStruct            //悬赏相关
	Achievement *AchievementStruct          //成就系统
	player      *Player                     //player引用
}

func (this *TaskStruct) Init(player *Player) {
	this.Tasks = make(map[int32]*HandleTaskStruct)
	this.XuanShang = new(XuanShangStruct)
	this.Achievement = new(AchievementStruct)
	this.player = player
	this.Achievement.Init() //成就初始化
}

//模板1
func (this *TaskStruct) Temple1(task *TaskInfo, par1 int32, par2 int32) bool {
	if task.Par == par2 {
		task.Progress += par1
	}

	if task.Progress >= task.Total_Progress {
		return true
	}
	return false
}

//模板2
func (this *TaskStruct) Temple2(task *TaskInfo, num int32) bool {
	task.Progress = num
	if task.Progress >= task.Total_Progress {
		return true
	}
	return false
}

//模板3
func (this *TaskStruct) Temple3(task *TaskInfo, num int32) bool {
	task.Progress += num
	if task.Progress >= task.Total_Progress {
		return true
	}
	return false
}

//购买道具
func (this *TaskStruct) BuyPropsTask(task *TaskInfo, count int32) bool {
	return true
}

//获得X个X道具
func (this *TaskStruct) GetPropsTask(task *TaskInfo, count int32, par2 int32) bool {
	return this.Temple1(task, count, par2)
}

//主角达到xx等级
func (this *TaskStruct) UpdateLevelTask(task *TaskInfo, level int32) bool {
	return this.Temple2(task, level)
}

//获取xx经验
func (this *TaskStruct) GetExpTask(task *TaskInfo, exp int32) bool {
	return this.Temple3(task, exp)
}

//获取xx铜钱
func (this *TaskStruct) GetGoldTask(task *TaskInfo, gold int32) bool {
	return this.Temple3(task, gold)
}

//拥有多少个英雄(必须得出已经拥有的英雄)
func (this *TaskStruct) GetHerosTask(task *TaskInfo, hero_num int32) bool {
	return this.Temple3(task, hero_num)
}

//拥有x个x阶的多少个英雄(必须得出已经拥有的英雄)
func (this *TaskStruct) GetHeroStepTask(task *TaskInfo, num int32, step int32) bool {
	return this.Temple1(task, num, step)
}

//拥有x个x星的多少个英雄(必须得出已经拥有的英雄)
func (this *TaskStruct) GetHeroStarTask(task *TaskInfo, num int32, star int32) bool {
	return this.Temple1(task, num, star)
}

//获得x件装备
func (this *TaskStruct) GetEquipTask(task *TaskInfo, equip_num int32) bool {
	return this.Temple3(task, equip_num)
}

//拥有x件x品质的装备
func (this *TaskStruct) GetEquipQualityTask(task *TaskInfo, equip_num int32, par2 int32) bool {
	return this.Temple1(task, equip_num, par2)
}

//拥有x件x星装备
func (this *TaskStruct) GetEquipStarTask(task *TaskInfo, equip_num int32, par2 int32) bool {
	return this.Temple1(task, equip_num, par2)
}

//进阶英雄x次
func (this *TaskStruct) GetHeroStep(task *TaskInfo, add_num int32) bool {
	return this.Temple3(task, add_num)
}

//强化装备X次
func (this *TaskStruct) GetStrengthen(task *TaskInfo, add_num int32) bool {
	return this.Temple3(task, add_num)
}

//精炼装备X次
func (this *TaskStruct) GetJinglian(task *TaskInfo, add_num int32) bool {
	return this.Temple3(task, add_num)
}

//击杀x个个怪物
func (this *TaskStruct) GetKillMonster(task *TaskInfo, add_num int32) bool {
	return this.Temple3(task, add_num)
}

//通关某关卡
func (this *TaskStruct) GetPassLevel(task *TaskInfo, add_num int32, par2 int32) bool {
	return this.Temple1(task, add_num, par2)
}

//进行X次单抽
func (this *TaskStruct) GetOneNiudan(task *TaskInfo, add_num int32) bool {
	return this.Temple3(task, add_num)
}

//进行X次10连抽
func (this *TaskStruct) GetTenNiudan(task *TaskInfo, add_num int32) bool {
	return this.Temple3(task, add_num)
}

//添加X个好友
func (this *TaskStruct) GetFriend(task *TaskInfo, add_num int32) bool {
	return this.Temple3(task, add_num)
}

//进行X次签到
func (this *TaskStruct) GetSign(task *TaskInfo, add_num int32) bool {
	return this.Temple3(task, add_num)
}

//参加X次竞技场
func (this *TaskStruct) GetArena(task *TaskInfo, add_num int32) bool {
	return this.Temple3(task, add_num)
}

//竞技场获胜X次
func (this *TaskStruct) GetArenaSucess(task *TaskInfo, add_num int32) bool {
	return this.Temple3(task, add_num)
}

//竞技场达到前X名
func (this *TaskStruct) GetArenaRank(task *TaskInfo, rank_num int32) bool {
	if rank_num < task.Total_Progress {
		return true
	}
	return false
}

//战力达到x
func (this *TaskStruct) GetPower(task *TaskInfo, add_num int32) bool {
	return this.Temple2(task, add_num)
}

//穿戴X件X部位的装备
func (this *TaskStruct) GetwearEquip(task *TaskInfo, add_num int32) bool {
	return false
}

//上阵x个英雄
func (this *TaskStruct) GetHeroOnFormation(task *TaskInfo, add_num int32) bool {
	return this.Temple2(task, add_num)
}

//使用道具
func (this *TaskStruct) UseProp(task *TaskInfo, add_num int32) bool {
	return this.Temple3(task, add_num)
}

//扣除某个道具完成某任务
func (this *TaskStruct) DeductionProps(task *TaskInfo, num int32, par int32) bool {
	return this.Temple1(task, num, par)
}

//处理整个逻辑
func (this *TaskStruct) DealEvent(task *TaskInfo, par1 int32, par2 int32) bool { //par1:数量 par2:比较类型
	var result bool = false
	switch task.Sub_id {
	case 1:
		result = this.BuyPropsTask(task, par1) //购买道具
	case 2:
		result = this.GetPropsTask(task, par1, par2) //获得X个X道具
	case 3:
		result = this.UpdateLevelTask(task, par1) //主角达到X级
	case 4:
		result = this.GetExpTask(task, par1) //获得XX经验
	case 5:
		result = this.GetGoldTask(task, par1) //获得XX铜钱
	case 6:
		result = this.GetHerosTask(task, par1) //拥有X个英雄
	case 7:
		result = this.GetHeroStepTask(task, par1, par2) //拥有X个X阶英雄
	case 8:
		result = this.GetHeroStarTask(task, par1, par2) //拥有X个X星英雄
	case 9:
		result = this.GetEquipTask(task, par1) //获得X件装备
	case 10:
		result = this.GetEquipQualityTask(task, par1, par2) //拥有X个X品质装备
	case 11:
		result = this.GetEquipStarTask(task, par1, par2) //拥有强化次数
	case 12:
		result = this.GetHeroStep(task, par1) //进阶英雄X次
	case 13:
		result = this.GetStrengthen(task, par1) //强化装备X次
	case 14:
		result = this.GetJinglian(task, par1) //精炼装备X次
	case 15:
		result = this.GetKillMonster(task, par1) //击杀X个怪物
	case 16:
		result = this.GetPassLevel(task, par1, par2) //通关某关卡
	case 17:
		result = this.GetOneNiudan(task, par1) //进行X次单抽
	case 18:
		result = this.GetTenNiudan(task, par1) //进行X次10连抽
	case 19:
		result = this.GetFriend(task, par1) //添加X个好友
	case 20:
		result = this.GetSign(task, par1) //进行X次签到
	case 21:
		result = this.GetArena(task, par1) //参加X次竞技场
	case 22:
		result = this.GetArenaSucess(task, par1) //竞技场获胜X次
	case 23:
		result = this.GetArenaRank(task, par1) //竞技场达到前X名
	case 24:
		result = this.GetPower(task, par1) //战斗力达到X
	case 25:
		//result = this.GetwearEquip(task, par1) //穿戴X件X部位的装备
	case 26:
		result = this.GetHeroOnFormation(task, par1) //上阵X个英雄
	case 27:
		result = this.UseProp(task, par1) //使用道具
	case 28:
		result = this.DeductionProps(task, par1, par2) //使用xx个道具
	case 100:
		result = true
	}
	return result
}

//获取当前状态的数据(用于初始化创建)
func (this *TaskStruct) GetProgress(task_id int32, par int32) int32 { //任务id 附加参数
	var num int32 = 0
	switch task_id {
	case 6: //拥有多少个英雄
		num = int32(len(this.player.Heros))
	case 7: //拥有X个X阶英雄
		for _, v := range this.player.Heros {
			if v.Hero_Info.Step_level == par {
				num += 1
			}
		}
	case 8: //拥有X个X星英雄
		for _, v := range this.player.Heros {
			if v.Hero_Info.Star_level == par {
				num += 1
			}
		}
	case 10: //用有多少个品质装备
		for _, v := range this.player.Bag_Equip.BagEquip {
			if v.Quality == par {
				num += 1
			}
		}
	case 11: //拥有X个X星装备
		for _, v := range this.player.Bag_Equip.BagEquip {
			if v.Strengthen_level == par {
				num += 1
			}
		}
	}
	return num
}

//外部调用
func (this *TaskStruct) TriggerEvent(event_type int32, par1 int32, par2 int32) { //event_type:事件类型 par1:数量 par2:附加参数
	//调用成就
	fmt.Println("成就调用:", event_type, par1, par2)
	this.Achievement.TriggerAchievementEvent(event_type, this.player)

	for key, v1 := range this.Tasks { //添加对应事件
		for _, v2 := range v1.Task_info {
			if event_type == v2.Sub_id && v2.Par == par2 {
				this.DealEvent(v2, par1, par2)
				this.Notice2CUpdateProgress(key, v2.Sub_id)
			}
		}
	}

	for key, v1 := range this.Tasks { //判断能否领取奖励
		for _, v2 := range v1.Task_info {
			if v2.Progress < v2.Total_Progress {
				return
			}
		}

		if Csv.quest[key].Id_116 == 1 { //自动领取奖励
			this.NoticeReward(key)
		} else {
			this.Notice2CHandleTask(v1.Type_id, key) //非自动的则提醒客户端手动接取任务
		}
	}
}

//获取任务事件
func (this *TaskStruct) TaskEvent(Task_event []TaskData, Type_id int32) []*TaskInfo {
	var task_infos []*TaskInfo
	for _, v := range Task_event {
		task_info := new(TaskInfo)
		task_info.Sub_id = v.Sub_id                                     //事件id
		task_info.Progress = this.GetProgress(task_info.Sub_id, v.Par2) //当前进度
		task_info.Total_Progress = v.Par1                               //总进度
		task_info.Par = v.Par2                                          //附加参数 用来条件判断
		task_infos = append(task_infos, task_info)
	}
	return task_infos
}

//具体任务消息格式
func (this *TaskStruct) ProtocolNewTaskInfo(id int32, tasks []*TaskInfo) *protocol.TaskInfo {
	//任务类型
	typeInfo := new(protocol.TaskType)
	typeInfo.Type = &this.Tasks[id].Type_id
	typeInfo.Id = &id
	typeInfo.QualityRef = &this.Tasks[id].Quality_ref

	//任务下面的具体子项任务
	var subitems []*protocol.SubItem
	for _, v_buff := range tasks {
		v := v_buff
		subitem := new(protocol.SubItem)
		subitem.TaskId = &v.Sub_id
		subitem.CompleteProgress = &v.Progress
		subitem.CompleteTotal = &v.Total_Progress
		subitems = append(subitems, subitem)
	}

	//任务
	task_info := new(protocol.TaskInfo)
	task_info.TypeInfo = typeInfo
	task_info.Subitems = subitems

	return task_info
}

//推送接取新任务消息
func (this *TaskStruct) Notice2CNewTask(id int32) {
	protocol_task_info := this.ProtocolNewTaskInfo(id, this.Tasks[id].Task_info)

	result4C := &protocol.NoticeMsg_Notice2CNewTask{
		Task: protocol_task_info,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.player.conn, 1214, encObj)
}

//创建任务
func (this *TaskStruct) CreateNewTask(Type_id int32, id int32) { //Type_id:1 主线 Type_id:悬赏
	//配置检查
	if Type_id == 1 {
		if _, ok := Csv.quest[id]; !ok { //主线任务
			return
		}
	}

	if Type_id == 2 {
		if _, ok := Csv.quest_xuanshang_array[id]; !ok { //悬赏任务
			return
		}
	}

	if Type_id == 1 {
		this.AddNewTask(1, id, Csv.quest[id].Task_event, 2, 0)
	} else if Type_id == 2 {
		index := this.GetQualityId(Csv.quest_xuanshang_array[id])
		this.AddNewTask(2, id, Csv.quest_xuanshang_array[id].Task_event, 2, index)
	}
	this.Notice2CNewTask(id)
}

//任务变化推送
func (this *TaskStruct) Notice2CUpdateProgress(id int32, submit_id int32) {
	if _, ok := this.Tasks[id]; !ok {
		return
	}

	typeInfo := new(protocol.TaskType) //任务类型
	typeInfo.Type = &this.Tasks[id].Type_id
	typeInfo.Id = &id

	subitem := new(protocol.SubItem) //任务下面的具体子项任务
	for _, v := range this.Tasks[id].Task_info {
		if submit_id == v.Sub_id {
			subitem.TaskId = &v.Sub_id
			subitem.CompleteProgress = &v.Progress
			subitem.CompleteTotal = &v.Total_Progress
			break
		}
	}

	result4C := &protocol.NoticeMsg_Notice2CUpdateProgress{
		TypeInfo: typeInfo,
		Subitem:  subitem,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.player.conn, 1216, encObj)
}

//input:奖励物品 装备品质组
//output: 金币 主角经验 道具 装备
func (this *TaskStruct) Reward(reward []RewardData, quality_item int32) (int32, int32, []*protocol.RwardProp, []int32) {
	var gold int32 = 0
	var role_exp int32 = 0
	var reward_props []*protocol.RwardProp
	var equip_uids []int32

	for _, v := range reward {
		if v.Type == 1 { //道具
			var prop Prop
			prop.Prop_id = v.Data_reward.Id
			prop.Count = v.Data_reward.Num
			prop.Prop_uid = GetUid()
			this.player.Bag_Prop.AddAndNotice(prop, this.player.conn)

			reward_prop := new(protocol.RwardProp)
			reward_prop.Num = &v.Data_reward.Num
			reward_prop.PropUid = &prop.Prop_uid
			reward_props = append(reward_props, reward_prop)
		}

		if v.Type == 4 { //资源类型
			if v.Data_reward.Id == 30002 { //铜钱
				gold += v.Data_reward.Num
			}

			if v.Data_reward.Id == 30003 { //主角经验
				role_exp += v.Data_reward.Num
			}
		}

		if v.Type == 3 { //装备类型
			var equips []Equip
			var equip Equip

			for i := 0; i <= int(v.Data_reward.Num); i++ {
				my_equip := equip.Create(v.Data_reward.Id, quality_item, this.player)
				equips = append(equips, *my_equip)
				this.player.Bag_Equip.Adds(equips, this.player.conn)
				equip_uids = append(equip_uids, my_equip.Equip_uid)
			}
		}
	}

	return gold, role_exp, reward_props, equip_uids
}

//任务奖励消息格式
func (this *TaskStruct) ProtocolRewardInfo(id int32) (int32, int32, []*protocol.RwardProp, []int32) { //金币 主角经验 道具 装备
	var gold int32 = 0
	var role_exp int32 = 0
	var reward_props []*protocol.RwardProp
	var equip_uids []int32

	if this.Tasks[id].Type_id == 1 { //主线奖励
		gold, role_exp, reward_props, equip_uids = this.Reward(Csv.quest[id].Reward, Csv.quest[id].Id_117)
		this.CreateNewTask(this.Tasks[id].Type_id, Csv.quest[id].Id_108) //产生新任务
	} else if this.Tasks[id].Type_id == 2 { //悬赏任务
		//算出任务品质index
		ref_num := this.XuanShang.Ref_Total_num //刷新次数
		add_quanzhi := Csv.quest_xuanshang_array[id].AddQuanzhi
		task_quality := Csv.quest_xuanshang_array[id].Task_quality
		task_quality[3].Num += add_quanzhi[0] * ref_num //新品质4增加权值
		task_quality[4].Num += add_quanzhi[1] * ref_num //新品质5增加权值
		var quality_list []int32
		for _, v := range task_quality {
			quality_list = append(quality_list, v.Num)
		}
		index := GetRandomIndex(quality_list)

		//发放物品
		gold_buff, role_exp_buff, _, _ := this.Reward(Csv.quest[id].Reward, 0) //悬赏任务只奖励铜钱跟exp

		//进行加倍计算
		gold = int32(Csv.quest_xuanshang_array[id].Multiple[index] * float32(gold_buff))
		role_exp = int32(Csv.quest_xuanshang_array[id].Multiple[index] * float32(role_exp_buff))
	}

	//移除完成任务
	delete(this.Tasks, id)

	return gold, role_exp, reward_props, equip_uids
}

//推送自动领取奖励
func (this *TaskStruct) NoticeReward(id int32) { //Type_id:(1:主线 2:悬赏任务 3:奇遇) id:csv中ID
	TaskInfo := new(protocol.TaskType)
	TaskInfo.Type = &this.Tasks[id].Type_id
	TaskInfo.Id = &id

	if this.Tasks[id].Type_id == 1 {
		gold, role_exp, reward_props, equip_uids := this.ProtocolRewardInfo(id)
		result4C := &protocol.NoticeMsg_SubmitTask{
			TypeInfo:  TaskInfo,
			Gold:      &gold,
			RoleExp:   &role_exp,
			Props:     reward_props,
			EquipUids: equip_uids,
		}

		encObj, _ := proto.Marshal(result4C)
		SendPackage(*this.player.conn, 1217, encObj)
	}
}

//推送非自动提交任务，客户端需手动提交任务获取奖励
func (this *TaskStruct) Notice2CHandleTask(Type_id int32, id int32) {
	TaskInfo := new(protocol.TaskType)
	TaskInfo.Type = &Type_id
	TaskInfo.Id = &id

	result4C := &protocol.NoticeMsg_Notice2CHandleTask{
		TypeInfo: TaskInfo,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.player.conn, 1218, encObj)
}

func (this *TaskStruct) GetQualityId(info Quest_xuanshang_info) int32 {

	ref_num := this.XuanShang.Ref_Total_num //刷新次数
	add_quanzhi := info.AddQuanzhi
	task_quality := info.Task_quality
	task_quality[3].Num += add_quanzhi[0] * ref_num //新品质4增加权值
	task_quality[4].Num += add_quanzhi[1] * ref_num //新品质5增加权值
	var quality_list []int32
	for _, v := range task_quality {
		quality_list = append(quality_list, v.Num)
	}
	index := GetRandomIndex(quality_list)
	return (index + 1)
}

func (this *TaskStruct) AddNewTask(type_id int32, id int32, Task_event []TaskData, stage int32, Quality_ref int32) {
	task := new(HandleTaskStruct) //任务列表
	task.Type_id = type_id
	task.Stage = stage
	task.Quality_ref = Quality_ref

	var task_infos []*TaskInfo //子项任务
	for _, v := range Task_event {
		info := new(TaskInfo)
		info.Par = v.Par2
		info.Progress = this.GetProgress(v.Sub_id, v.Par2)
		info.Sub_id = v.Sub_id
		info.Total_Progress = v.Par1
		task_infos = append(task_infos, info)
	}

	task.Task_info = task_infos
	this.Tasks[id] = task
}

//定时器定时刷新悬赏任务
func (this *TaskStruct) TaskXuanshangTimer(is_timer bool) { //是否是定时器驱动
	//四个刷新任务
	if is_timer {
		this.XuanShang.Last_time = int32(time.Now().Unix())
	}

	//删除刷新的未接取任务
	for key, v := range this.Tasks {
		if v.Type_id == 2 {
			if v.Stage == 1 || v.Stage == 4 { //可以领取但未领取
				delete(this.Tasks, key)
			}
		}
	}

	//添加新的悬赏任务
	for _, v := range Csv.quest_xuanshang {
		if this.player.Info.Level >= v.Min_level && this.player.Info.Level <= v.Max_level {
			task_len := len(v.Xuanshang)

			if int32(task_len) >= 4 { //产生4个任务
				rand_num_int := rand_.Intn(task_len - 3)
				rand_num := int32(rand_num_int)

				index := this.GetQualityId(v.Xuanshang[rand_num])
				this.AddNewTask(2, v.Xuanshang[rand_num].Id_101, v.Xuanshang[rand_num].Task_event, 1, index)

				index = this.GetQualityId(v.Xuanshang[rand_num])
				this.AddNewTask(2, v.Xuanshang[rand_num+1].Id_101, v.Xuanshang[rand_num+1].Task_event, 1, index)

				index = this.GetQualityId(v.Xuanshang[rand_num])
				this.AddNewTask(2, v.Xuanshang[rand_num+2].Id_101, v.Xuanshang[rand_num+2].Task_event, 1, index)

				index = this.GetQualityId(v.Xuanshang[rand_num])
				this.AddNewTask(2, v.Xuanshang[rand_num+3].Id_101, v.Xuanshang[rand_num+3].Task_event, 1, index)

				fmt.Println("定时器添加新的悬赏任务:")
				break
			}
		}
	}

	//推送
	this.NoticeCanAccept()
}

//推送可以接取的手动任务
func (this *TaskStruct) NoticeCanAccept() {
	var canAcceptTasks []*protocol.TaskType
	for key_buff, v_buff := range this.Tasks {
		if v_buff.Stage == 1 {
			v := v_buff
			key := key_buff
			canAcceptTask := new(protocol.TaskType)
			canAcceptTask.Type = &v.Type_id
			canAcceptTask.Id = &key
			canAcceptTasks = append(canAcceptTasks, canAcceptTask)
			fmt.Println("推送可以接取的手动任务:", v)
		}
	}

	result4C := &protocol.NoticeMsg_Notice2CCanAccept{
		CanAcceptTasks: canAcceptTasks,
	}
	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.player.conn, 1215, encObj)
}

//当前所有任务
func (this *TaskStruct) AllTask(conn *net.Conn) {
	var task_infos []*protocol.TaskInfo

	for key, v1 := range this.Tasks {
		if v1.Stage > 1 {
			task_info := this.ProtocolNewTaskInfo(key, v1.Task_info)
			task_infos = append(task_infos, task_info)
		}
	}

	result4C := &protocol.Task_AllTaskResult{
		Tasks: task_infos,
	}
	encObj, _ := proto.Marshal(result4C)
	SendPackage(*conn, 1801, encObj)
}

//手动接受任务
func (this *TaskStruct) AcceptTask(Type_id int32, id int32) {
	fmt.Println("手动接受任务")
	var result int32 = 1
	for key, v := range this.Tasks {
		if v.Stage == 1 && v.Type_id == Type_id && key == id {
			this.CreateNewTask(Type_id, id)
			break
		}
	}

	result4C := &protocol.Task_AcceptTaskResult{
		Result: &result,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.player.conn, 1802, encObj)
}

//检查任务是否完成
func (this *TaskStruct) IsComplete(id int32) bool {

	if _, ok := this.Tasks[id]; ok {
		for _, v := range this.Tasks[id].Task_info {
			if v.Progress < v.Total_Progress {
				return false
			}
		}
	}

	return true
}

//手动提交任务获取奖励
func (this *TaskStruct) SubmitTask(Type_id int32, id int32) {

	if _, ok := this.Tasks[id]; !ok {
		return
	}

	//检查是否完成任务
	is_complete := this.IsComplete(id)

	if is_complete {
		var result int32 = 0
		gold, role_exp, reward_props, equip_uids := this.ProtocolRewardInfo(id)

		result4C := &protocol.Task_SubmitTaskResult{
			Result:    &result,
			Gold:      &gold,
			RoleExp:   &role_exp,
			Props:     reward_props,
			EquipUids: equip_uids,
		}
		encObj, _ := proto.Marshal(result4C)
		SendPackage(*this.player.conn, 1803, encObj)
	}

}

//获取悬赏任相关
func (this *TaskStruct) GetXuanShangInfo(conn *net.Conn) {
	//悬赏任务刷新时间/秒
	csv_time := int32(Csv.property[2013].Id_102)
	now_time := int32(time.Now().Unix())
	last_time := csv_time + this.XuanShang.Last_time - now_time
	if last_time < 0 {
		last_time = 0
	}

	fmt.Println("获取悬赏任相关:", csv_time, now_time, last_time)
	var last_num int32 = int32(Csv.property[2015].Id_102) - this.XuanShang.Ref_Total_num
	result4C := &protocol.Task_GetXuanShangInfoResult{
		LastNum:      &last_num,
		FreeLastTime: &last_time,
	}
	encObj, _ := proto.Marshal(result4C)
	SendPackage(*conn, 1821, encObj)
}

//元宝刷新
func (this *TaskStruct) XuanShangDiamondRef() {
	var result int32 = 0

	if this.XuanShang.Ref_Total_num >= int32(Csv.property[2015].Id_102) {
		result = 2
	}

	need_diamond := int32(Csv.property[2014].Id_102) + 5*this.XuanShang.Ref_Total_num
	if this.player.Info.Diamond < need_diamond {
		result = 1
	}

	if result == 0 {
		this.player.ModifyDiamond(-need_diamond)
		this.XuanShang.Ref_Total_num += 1
		this.TaskXuanshangTimer(false)
	}

	result4C := &protocol.Task_XuanShangDiamondRefResult{
		Result: &result,
	}
	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.player.conn, 1822, encObj)
}

//任务放弃
func (this *TaskStruct) GiveUpTask(id int32) {
	var result int32 = 0
	if _, ok := this.Tasks[id]; ok {
		if this.Tasks[id].Type_id == 1 {
			result = 1
		} else {
			delete(this.Tasks, id)
		}
	}

	result4C := &protocol.Task_GiveUpTaskResult{
		Result: &result,
	}
	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.player.conn, 1824, encObj)
}
