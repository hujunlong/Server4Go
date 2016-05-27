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

func (this *Deal4C) Handler4C(conn net.Conn) {

	var player *game.Player //玩家对象的指针
	var is_login bool = false

	defer func() {
		if player != nil {
			player.ExitGame()
		}
		fmt.Println("socket is close")
		conn.Close()
	}()

	const MAXLEN = 10240
	buff := make([]byte, MAXLEN)

	for {
		n, err := conn.Read(buff) //接收具体消息
		if err != nil {
			return
		}

		if n > MAXLEN || n < 8 {
			return
		}

		//接收包头
		for n >= 8 {
			body_len, head_pid := GetHead(buff[:8])
			if int(body_len) > n {
				return
			}
			buf := buff[:body_len]
			buff = buff[body_len:]
			n = n - int(body_len)

			//添加开关
			if (head_pid > 1002) && !is_login {
				return
			}

			switch head_pid {

			case 1001: //注册
				register := new(protocol.Game_RegisterRole)
				if err := proto.Unmarshal(buf[8:body_len], register); err == nil {

					player = new(game.Player)
					player.Init()

					result := player.RegisterRole(int64(register.GetPlayerId()+this.server_id*1000000), register.GetNick(), register.GetHeroId(), &conn)
					result4C := &protocol.Game_RegisterRoleResult{
						Result: proto.Int32(result),
					}
					encObj, _ := proto.Marshal(result4C)
					SendPackage(conn, 1001, encObj)
				}

			case 1002: //获取player基础属性
				get_info := new(protocol.Game_GetRoleInfo)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				//内存数据库
				is_login, player = player.Login(int64(get_info.GetPlayerId()+this.server_id*1000000), &conn)
			case 1019: //玩家退出
				player.ExitGame()

			case 1003: //战斗准备进行关卡
				get_info := new(protocol.Game_WarMapStage)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}

				result := player.Stage.IsCanThroughMap(get_info.GetStageId(), player.Info.Energy, 1, 1)
				result4C := &protocol.Game_MapStageResult{
					Result: &result,
				}
				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1003, encObj)

			case 1004: //战斗结果客户的通知服务器
				get_info := new(protocol.Game_WarMapNoteServer)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				player.WarMapNoteServerResult(get_info.GetStage().GetState(), get_info.GetStage().GetStageId())

			case 1005: //扫荡
				get_info := new(protocol.Game_SweepMapStage)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}

				player.SweepMapStageResult(get_info.GetStageId(), get_info.GetCount())

			case 1015: //在线挂机
				player.OnNotice2CGuaji()
			case 1016: //挂机事件
				get_info := new(protocol.GameGetGuajiInfo)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				result4C := player.Guaji_Stage.GuajiInfoResult(get_info.GetId())
				if result4C != nil {
					encObj, _ := proto.Marshal(result4C)
					SendPackage(conn, 1016, encObj)
				}

			case 1017: //挑战boss的阵容信息
				get_info := new(protocol.Game_ChallengeBoss)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				player.ChallengeBoss(get_info.GetId())

			case 1018: //客户端通知服务器boss挑战结果
				get_info := new(protocol.Game_C2SChallenge)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				player.C2SChallengeResult(get_info.Stage.GetState(), get_info.Stage.GetStageId())

			case 1020: //切换挂机地方
				get_info := new(protocol.Game_ChangeGuajiInfo)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				is_ok := player.Guaji_Stage.ChangeStage(get_info.GetId(), player.PlayerId)
				result4C := &protocol.Game_ChangeGuajiInfoResult{
					IsOk: &is_ok,
				}
				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1020, encObj)

			case 1021: //快速战斗
				get_info := new(protocol.Game_FastWar)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				result4C := player.Guaji_Stage.FastWar(get_info.GetStage().GetStageId())
				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1021, encObj)

			case 1022: //获取玩家当前挂机列表
				GuajiRoleInfos_ := player.Guaji_Stage.GetGuajiRoleListResult()
				result4C := &protocol.Game_GetGuajiRoleListResult{
					GuajiRoleInfos: GuajiRoleInfos_,
				}
				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1022, encObj)

			case 1023: //请求阵型
				get_info := new(protocol.Game_GetGuajiRoleFormation)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				player.GetFomation(get_info.GetRoleId(), get_info.GetType())

			case 1024: //英雄上阵 下阵
				get_info := new(protocol.Game_HerosFormation)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				player.HerosFormation(get_info.GetType(), get_info.GetIsOn(), get_info.GetPosId(), get_info.GetHeroUid())

			case 1025: //位置交换
				get_info := new(protocol.Game_ChangeHerosFormation)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				player.ChangeHerosFormation(get_info.GetPosId_1(), get_info.GetPosId_2())

			case 1026: //挑战玩家
				get_info := new(protocol.Game_ChallengePlayer)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				player.ChallengePlayer(get_info.GetType(), get_info.GetRoleId())

			case 1027: //使用某个道具
				get_info := new(protocol.Game_UseProp)
				if err := proto.Unmarshal(buf[8:body_len], get_info); err != nil {
					return
				}
				result, props := player.Bag_Prop.Use(get_info.GetUid(), get_info.GetCount())
				result4C := &protocol.Game_UsePropResult{
					Result: &result,
				}
				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1027, encObj)

				if result == 0 { //0:ok 1:不存该道具id 2:道具总量少于请求数量
					player.Notice2CProp(2, props) //道具删除
				}

			case 1029: //查看道具背包
				result4C := new(protocol.Game_CheckPropBagResult)
				var propStruct_s []*protocol.Game_PropStruct

				for _, buff_v := range player.Bag_Prop.Props {
					v := buff_v
					propStruct := new(protocol.Game_PropStruct)
					propStruct.PropUid = &v.Prop_uid
					propStruct.PropId = &v.Prop_id
					propStruct.PropCount = &v.Count
					propStruct_s = append(propStruct_s, propStruct)
				}

				result4C.PropStruct = propStruct_s
				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1029, encObj)

			case 1030: //查看装备背包
				result4C := new(protocol.Game_CheckEquipBagResult)
				var equipStruct_s []*protocol.Game_EquipStruct

				for _, v_buff := range player.Bag_Equip.BagEquip {
					v := v_buff
					equipStruct := new(protocol.Game_EquipStruct)
					EquipInfo := new(protocol.Game_EquipInfo)

					EquipInfo.Id = &v.Equip_id
					EquipInfo.Uid = &v.Equip_uid
					EquipInfo.Pos = &v.Pos
					EquipInfo.Quality = &v.Quality
					EquipInfo.EquipLevel = &v.Equip_level
					EquipInfo.StrengthenCount = &v.Strengthen_count
					EquipInfo.RefineCount = &v.Refine_count

					equipStruct.EquipInfo = EquipInfo
					equipStruct_s = append(equipStruct_s, equipStruct)
				}

				result4C.EquipStruct = equipStruct_s
				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1030, encObj)
			default:
			}
		}
	}
}
