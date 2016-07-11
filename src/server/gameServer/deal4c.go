package main

import (
	"fmt"
	"net"
	"server/gameServer/game"
	"server/share/protocol"

	"github.com/golang/protobuf/proto"
)

var sys_config *game.SysConfig

type Deal4C struct {
	server_id int32 //游戏服务器具体id编号
}

func (this *Deal4C) Init() {
	sys_config = new(game.SysConfig)
	sys_config.Init()
	this.server_id = sys_config.GameId
}

func (this *Deal4C) Deal4Client(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if CheckError(err) {
			go this.Handler4C(conn)
		}
	}
}

/*
基础属性 	 1001-1100
关卡挂机相关 1101-1200
推送变化     1201-1300
英雄阵型相关 1301-1400
英雄相关     1401-1500
道具装备相关 1501-1600
*/
func (this *Deal4C) Handler4C(conn net.Conn) {
	const MAXLEN = 10240
	buff := make([]byte, MAXLEN)

	var player *game.Player //实例化玩家
	is_login := false

	defer func() {
		if player != nil {
			player.ExitGame()
		}
		fmt.Println("socket is close")
		conn.Close()
	}()

	for {
		n, err := conn.Read(buff) //接收具体消息
		if err != nil {
			fmt.Println("err:", err)
			return
		}

		if n > MAXLEN || n < 8 {
			fmt.Println(" n > MAXLEN || n < 8", n)
			return
		}

		//接收包头
		for n >= 8 {
			body_len, head_pid := GetHead(buff[:8])
			if int(body_len) > n {
				fmt.Println("body_len > n", int(body_len), n)
				return
			}
			buf := buff[:body_len]
			buff = buff[body_len:]
			n -= int(body_len)

			//添加开关
			if (head_pid > 1002) && !is_login {
				fmt.Println("head_pid > 1002", head_pid, is_login)
				return
			}

			fmt.Println("数据接收 pid=", head_pid)

			//基础属性 1001-1100
			switch head_pid {

			case 1001: //注册
				register := new(protocol.PlayerBase_RegisterRole)
				if err := proto.Unmarshal(buf[8:body_len], register); err == nil {
					player = new(game.Player)
					result := player.RegisterRole(int64(register.GetPlayerId()+this.server_id*1000000), register.GetNick(), register.GetHeroId(), &conn)
					result4C := &protocol.PlayerBase_RegisterRoleResult{
						Result: proto.Int32(result),
					}
					encObj, _ := proto.Marshal(result4C)
					SendPackage(conn, 1001, encObj)
				}

			case 1002: //获取player基础属性
				fmt.Println("获取player基础属性")
				get_info := new(protocol.PlayerBase_GetRoleInfo)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				//内存数据库
				is_login, player = player.Login(int64(get_info.GetPlayerId()+this.server_id*1000000), &conn)
				if is_login {
					player.SetConn(&conn)
					player.Save()
				}

			case 1003: //玩家退出
				player.ExitGame()

			//关卡挂机相关 1101-1200
			case 1101: //战斗准备进行关卡
				get_info := new(protocol.StageBase_WarMapStage)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}

				result := player.Stage.IsCanThroughMap(get_info.GetStageId(), player.Info.Energy, 1, 1)
				result4C := &protocol.StageBase_MapStageResult{
					Result: &result,
				}
				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1101, encObj)

			case 1102: //战斗结果客户的通知服务器
				get_info := new(protocol.StageBase_WarMapNoteServer)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				player.Stage.WarMapNoteServerResult(get_info.GetStage().GetState(), get_info.GetStage().GetStageId(), player)

			case 1103: //扫荡
				get_info := new(protocol.StageBase_SweepMapStage)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				player.Stage.SweepMapStageResult(get_info.GetStageId(), get_info.GetCount(), player)

			case 1104: //在线挂机
				player.Guaji_Stage.OnNotice2CGuaji(player)

			case 1105: //挂机事件
				get_info := new(protocol.StageBaseGetGuajiInfo)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				result4C := player.Guaji_Stage.GuajiInfoResult(get_info.GetId())
				if result4C != nil {
					encObj, _ := proto.Marshal(result4C)
					SendPackage(conn, 1105, encObj)
				}

			case 1106: //挑战boss的阵容信息
				get_info := new(protocol.StageBase_ChallengeBoss)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				player.Guaji_Stage.ChallengeBoss(get_info.GetId(), player)

			case 1107: //客户端通知服务器boss挑战结果
				get_info := new(protocol.StageBase_C2SChallenge)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				player.Guaji_Stage.C2SChallengeResult(get_info.Stage.GetState(), get_info.Stage.GetStageId(), player)

			case 1108: //切换挂机地方
				get_info := new(protocol.StageBase_ChangeGuajiInfo)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				is_ok := player.Guaji_Stage.ChangeStage(get_info.GetId(), player.PlayerId)
				result4C := &protocol.StageBase_ChangeGuajiInfoResult{
					IsOk: &is_ok,
				}
				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1108, encObj)

			case 1109: //快速战斗
				get_info := new(protocol.StageBase_FastWar)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				result4C := player.Guaji_Stage.FastWar(get_info.GetStage().GetStageId())
				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1109, encObj)

			case 1110: //获取玩家当前挂机列表
				GuajiRoleInfos_ := player.Guaji_Stage.GetGuajiRoleListResult()
				result4C := &protocol.StageBase_GetGuajiRoleListResult{
					GuajiRoleInfos: GuajiRoleInfos_,
				}
				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1110, encObj)

			case 1111: //挑战玩家
				get_info := new(protocol.StageBase_ChallengePlayer)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				player.ChallengePlayer(get_info.GetType(), get_info.GetRoleId())

			//英雄阵型相关 1301-1400
			case 1301: //请求阵型
				get_info := new(protocol.Formation_GetGuajiRoleFormation)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				player.GetFomation(get_info.GetRoleId(), get_info.GetType())

			case 1302: //英雄上阵 下阵
				get_info := new(protocol.Formation_HerosFormation)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}

				if get_info.GetType() == 2 {
					player.StageFomation.HerosFormation(get_info.GetIsOn(), get_info.GetPosId(), get_info.GetHeroUid(), player)
				}

			case 1303: //位置交换
				get_info := new(protocol.Formation_ChangeHerosFormation)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				player.StageFomation.ChangeHerosFormation(get_info.GetPosId_1(), get_info.GetPosId_2(), player)

			//英雄相关     1401-1500
			case 1401: //获取英雄相关信息
				get_info := new(protocol.Hero_GetHeros)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				player.GetHeroResult(get_info.GetRoleId(), get_info.GetType())

			case 1402: //hero升阶
				get_info := new(protocol.Hero_StepHero)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}

				if _, ok := player.Heros[get_info.GetHeroUid()]; !ok {
					return
				}

				result := player.Heros[get_info.GetHeroUid()].StepUp(player)
				result4C := new(protocol.Hero_StepHeroResult)
				result4C.Result = &result
				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1402, encObj)

				fmt.Println(player.Heros[get_info.GetHeroUid()].Hero_Info.Step_level)
			case 1403: //hero升星
				get_info := new(protocol.Hero_StarHero)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}

				if _, ok := player.Heros[get_info.GetHeroUid()]; !ok {
					return
				}

				result := player.Heros[get_info.GetHeroUid()].StarUp(player)
				result4C := new(protocol.Hero_StarHeroResult)
				result4C.Result = &result
				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1403, encObj)

			case 1404: //天赋开启
				get_info := new(protocol.Hero_OpenHeroGift)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}

				hero_uid := get_info.GetHeroUid()
				if player.Heros[hero_uid] == nil {
					return
				}

				result := player.Heros[hero_uid].OpenHeroGift(get_info.GetId())
				result4C := &protocol.Hero_OpenHeroGiftResult{
					Result: &result,
				}
				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1404, encObj)

			case 1405: //英雄升级
				get_info := new(protocol.Hero_HeroLevelUp)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}

				hero_uid := get_info.GetHeroUid()
				if _, ok := player.Heros[hero_uid]; !ok {
					return
				}
				result := player.Heros[hero_uid].TaskGoods(get_info.GetPropId(), get_info.GetCount())
				result4C := &protocol.Hero_HeroLevelUpResult{
					Result: &result,
				}
				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1405, encObj)

				//道具装备相关 1501-1600
			case 1501: //使用某个道具
				get_info := new(protocol.Goods_UseProp)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}

				_, result := player.Bag_Prop.DeleteItemByUid(get_info.GetUid(), get_info.GetCount())
				result4C := &protocol.Goods_UsePropResult{
					Result: &result,
				}
				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1501, encObj)

			case 1502: //穿戴 或 卸载某装备
				get_info := new(protocol.Goods_UseEquip)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				result := player.Bag_Equip.UseEquip(get_info.GetEquipUid(), get_info.GetPos(), get_info.GetType(), player)
				result4C := &protocol.Goods_UseEquipResult{
					Result: &result,
				}
				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1502, encObj)

			case 1503: //装备强化
				get_info := new(protocol.Goods_EquipStrengthen)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}

				var result int32 = 0
				equip_uid := get_info.GetEquipUid()
				if _, ok := player.Bag_Equip.BagEquip[equip_uid]; ok {
					result = player.Bag_Equip.BagEquip[equip_uid].Strengthen(player)
				} else {
					result = 1
				}

				for _, v := range player.Bag_Equip.BagEquip {
					fmt.Println("equip:", v.Equip_id, v.Equip_uid)
				}

				result4C := &protocol.Goods_EquipStrengthenResult{
					Result: &result,
				}
				player.Save()
				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1503, encObj)

			case 1504: //装备精炼
				get_info := new(protocol.Goods_EquipRefine)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}

				var result int32 = 0
				equip_uid := get_info.GetEquipUid()
				if _, ok := player.Bag_Equip.BagEquip[equip_uid]; ok {
					result = player.Bag_Equip.BagEquip[equip_uid].EquipRefine(player)
				} else {
					result = 1
				}

				result4C := &protocol.Goods_EquipRefineResult{
					Result: &result,
				}
				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1504, encObj)
				player.Save()
			case 1505: //装备分解
				get_info := new(protocol.Goods_EquipDecompose)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}

				var result int32 = 0
				equip_uids := get_info.GetEquipUids()
				result = player.Bag_Equip.EquipDecompose(equip_uids, player)

				result4C := &protocol.Goods_EquipDecomposeResult{
					Result: &result,
				}
				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1505, encObj)
			case 1506: //道具出售
				get_info := new(protocol.Goods_SaleProp)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				player.Bag_Prop.SaleProp(get_info.GetPropUid(), get_info.GetCount())

			//GM
			case 1601:
				get_info := new(protocol.GM_Msg)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				player.Gm.DealMsg(get_info.GetMsg(), player)
				fmt.Println("GM:", get_info.GetMsg())

			//活动 1701-1800
			case 1701:
				get_info := new(protocol.Activity_NiuDan)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				result, props, golds, hero_uids := player.Niu_Dan.NiuDanMsg(get_info.GetType(), get_info.GetTypeGroup(), player)
				fmt.Print("活动：", hero_uids)
				protocl_pros := player.Bag_Prop.DealPropStruct(props)
				result4C := &protocol.Activity_NiuDanResult{
					Result:   &result,
					Props:    protocl_pros,
					Gold:     golds,
					HeroUids: hero_uids,
				}
				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1701, encObj)

			//任务系统 1801-1900

			case 1801: //满足的所有任务
				player.Task.AllTask(&conn)

			//c2s手动接受任务
			case 1802:
				get_info := new(protocol.Task_AcceptTask)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				player.Task.AcceptTask(get_info.GetTaskInfo().GetType(), get_info.GetTaskInfo().GetId())

			//pid 1803 手动提交任务获取奖励
			case 1803:
				get_info := new(protocol.Task_SubmitTask)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				player.Task.SubmitTask(get_info.GetTaskInfo().GetType(), get_info.GetTaskInfo().GetId())

			//获取悬赏任务相关
			case 1821:
				player.Task.GetXuanShangInfo(&conn)
			//悬赏元宝刷新请求
			case 1822:
				get_info := new(protocol.Task_XuanShangDiamondRef)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				player.Task.XuanShangDiamondRef()

			case 1824: //放弃任务
				get_info := new(protocol.Task_GiveUpTask)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				player.Task.GiveUpTask(get_info.GetTask().GetId())

			//请求获取成就列表
			case 1851:
				player.Task.Achievement.GetAchievementResult(&conn)

			//获取成就奖励
			case 1852:
				get_info := new(protocol.Task_GetAchievementReward)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				player.Task.Achievement.GetAchievementReward(get_info.GetId(), player)
			default:
			}
		}
	}
}
