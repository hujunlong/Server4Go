package main

import (
	"fmt"
	"net"
	"server/gameServer/game"
	"server/share/global"
	"server/share/protocol"
	"strconv"

	"github.com/game_engine/timer"
	"github.com/golang/protobuf/proto"
)

var sys_config *game.SysConfig

type Deal4C struct {
	word      *game.World
	server_id int32 //游戏服务器具体id编号
}

func (this *Deal4C) Init() {

	sys_config = new(game.SysConfig)
	sys_config.Init()
	this.server_id = sys_config.GameId

	this.word = new(game.World)
	this.word.Init()

	//开启定时器
	timer.CreateTimer(sys_config.DistanceTime, true, this.word.TimerDealOnlineGuaji)
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
	defer conn.Close()

	var player *game.Player //玩家对象的指针

	const MAXLEN = 1024
	buf := make([]byte, MAXLEN)
	for {

		n, err := conn.Read(buf) //接收具体消息
		if err != nil {
			return
		}

		if n > MAXLEN {
			game.Log.Error("recive error n> MAXLEN")
			return
		}

		fmt.Println(buf[0:n])

		//接收包头
		_, head_pid := GetHead(buf)

		//包体
		switch head_pid {
		case 1001: //注册
			register := new(protocol.Game_RegisterRole)
			if err := proto.Unmarshal(buf[8:n], register); err == nil {
				player = new(game.Player)
				player.Init()
				result := player.RegisterRole(int64(register.GetPlayerId()+this.server_id*1000000), register.GetNick(), register.GetHeroId(), &conn)
				result4C := &protocol.Game_RegisterRoleResult{
					Result: proto.Int32(result),
				}
				encObj, _ := proto.Marshal(result4C)
				fmt.Println("register:", encObj)
				SendPackage(conn, 1001, encObj)

				//加载保存于内存中
				if global.REGISTERROLESUCCESS == result {
					this.word.EnterWorld(player)
				}
			}

		case 1002: //获取player基础属性
			get_info := new(protocol.Game_GetRoleInfo)
			if err := proto.Unmarshal(buf[8:n], get_info); err != nil {
				return
			}

			//先查询world中是否存在该玩家 未查询到 就读取内存数据库
			player_id_str := strconv.FormatInt(int64(get_info.GetPlayerId()+this.server_id*1000000), 10)
			player = this.word.SearchPlayer(player_id_str)
			if player == nil {
				player = game.LoadPlayer(player_id_str)
				if player != nil {
					player.Init()
				}

			}

			if player == nil {
				fmt.Println("player_id_str not found id = ", player_id_str)
				game.Log.Error("player_id_str not found id = %s", player_id_str)
				return
			}

			//先判断是否创建角色
			fmt.Println("create time:", player.CreateTime)
			if player.CreateTime <= 0 {
				var is_create bool = false
				result4C := &protocol.Game_RoleInfoResult{
					IsCreate: &is_create,
				}

				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1002, encObj)

			} else {
				var is_create bool = true

				//英雄列表
				var HeroStruct_ []*protocol.Game_HeroStruct
				for i := 0; i < len(player.Heros); i++ {
					hero_struct := &protocol.Game_HeroStruct{
						HeroId:  &player.Heros[i].Hero_Info.Hero_id,
						HeroUid: &player.Heros[i].Hero_Info.Hero_uid,

						HeroInfo: &protocol.Game_HeroInfo{
							Level:     &player.Heros[i].Hero_Info.Level,
							Hp:        &player.Heros[i].Hero_Info.Hp,
							Power:     &player.Heros[i].Hero_Info.Power,
							StarLevel: &player.Heros[i].Hero_Info.Star_level,
							StepLevel: &player.Heros[i].Hero_Info.Step_level,
						},
					}
					HeroStruct_ = append(HeroStruct_, hero_struct)
					fmt.Println("player.Heros[i].Hero_Info.Hero_id:", player.Heros[i].Hero_Info.Hero_id)
				}

				//道具

				//装备

				//副本开启
				var type_ int32 = 1
				var copy_levels []*protocol.Game_Stage
				for id_str, v := range player.Stage.Map_stage_pass {
					copy_level := new(protocol.Game_Stage)
					copy_level.Type = &type_
					copy_level.State = &v
					stage_id_32 := game.Str2Int32(id_str)
					copy_level.StageId = &stage_id_32
					copy_levels = append(copy_levels, copy_level)
				}
				//挂机
				var type_guaji int32 = 2
				var guaji_stages []*protocol.Game_Stage
				for id_str, v := range player.Guaji_Stage.Guaji_Map_stage_pass {
					guaji_stage := new(protocol.Game_Stage)
					guaji_stage.Type = &type_guaji
					guaji_stage.State = &v
					stage_id_32 := game.Str2Int32(id_str)
					guaji_stage.StageId = &stage_id_32
					guaji_stages = append(guaji_stages, guaji_stage)
				}

				result4C := &protocol.Game_RoleInfoResult{
					IsCreate: &is_create,
					PlayerInfo: &protocol.Game_PlayerInfo{
						Level:     &player.Info.Level,
						Exp:       &player.Info.Exp,
						Hp:        &player.Info.Hp,
						Energy:    &player.Info.Energy,
						EnergyMax: &player.Info.EnergyMax,
						Vip:       &player.Info.Vip,
						Gold:      &player.Info.Gold,
						Diamond:   &player.Info.Diamond,
						Power:     &player.Info.Power,
						Nick:      &player.Info.Nick,
						Signature: &player.Info.Signature,
						Option:    player.Info.Option,
					},

					HeroStruct:   HeroStruct_,
					CopyLevels:   copy_levels,
					HangupLevels: guaji_stages,
				}

				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1002, encObj)
			}

		case 1019: //玩家退出
			player.ExitGame()

		case 1003: //战斗准备进行关卡
			get_info := new(protocol.Game_WarMapStage)
			if err := proto.Unmarshal(buf[8:n], get_info); err != nil {
				return
			}

			result := player.Stage.IsCanThroughMap(get_info.GetStageId(), player.Info.Energy, 1)
			result4C := &protocol.Game_MapStageResult{
				Result: &result,
			}
			encObj, _ := proto.Marshal(result4C)
			SendPackage(conn, 1003, encObj)

		case 1004: //战斗结果客户的通知服务器
			get_info := new(protocol.Game_WarMapNoteServer)
			if err := proto.Unmarshal(buf[8:n], get_info); err != nil {
				return
			}
			player.WarMapNoteServerResult(get_info.GetStage().GetState(), get_info.GetStage().GetStageId())

		case 1005: //扫荡
			get_info := new(protocol.Game_SweepMapStage)
			if err := proto.Unmarshal(buf[8:n], get_info); err != nil {
				return
			}
			player.SweepMapStageResult(get_info.GetStageId(), get_info.GetCount())

		case 1016: //挂机事件
			get_info := new(protocol.GameGetGuajiInfo)
			if err := proto.Unmarshal(buf[8:n], get_info); err != nil {
				return
			}
			result4C := player.Guaji_Stage.GuajiInfoResult(get_info.GetId())
			if result4C != nil {
				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1016, encObj)
			}

		case 1017: //挑战boss
			get_info := new(protocol.Game_ChallengeBoss)
			if err := proto.Unmarshal(buf[8:n], get_info); err != nil {
				return
			}

		case 1018: //客户端通知服务器boss挑战结果
			get_info := new(protocol.Game_C2SChallenge)
			if err := proto.Unmarshal(buf[8:n], get_info); err != nil {
				return
			}
			player.C2SChallengeResult(get_info.Stage.GetState(), get_info.Stage.GetStageId())
		default:
		}
	}
}
