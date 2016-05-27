package main

import (
	"server/share/global"
	"server/share/protocol"

	"github.com/game_engine/logs"
	"github.com/golang/protobuf/proto"
)

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

var nick string = "ttt"
var pwd string = "ttt"

var log *logs.BeeLogger

func init() {
	log = logs.NewLogger(100000) //日志
	log.EnableFuncCallDepth(true)
	log.SetLogger("file", `{"filename":"log/testLoginServer.log"}`)
}

func CheckError(err error) bool {
	if err != nil {
		fmt.Println("err:", err)
		return false
	}
	return true
}

func SendPackage(conn *net.Conn, pid int, body []byte) {
	var pid_32 int32 = int32(pid)

	len := 8 + len(body)
	var len_32 = int32(len)

	len_buf := bytes.NewBuffer([]byte{})
	binary.Write(len_buf, binary.BigEndian, &len_32)

	pid_buf := bytes.NewBuffer([]byte{})
	binary.Write(pid_buf, binary.BigEndian, &pid_32)

	msg := append(len_buf.Bytes(), pid_buf.Bytes()...)
	msg2 := append(msg, body...)
	(*conn).Write(msg2)
	fmt.Println(msg2)
}

func GetHead(buf []byte) (int32, int32) {
	if len(buf) < 8 {
		return 0, 0
	}

	var head_len int32 = 0
	var head_pid int32 = 0
	buffer_len := bytes.NewBuffer(buf[0:4])
	buffer_pid := bytes.NewBuffer(buf[4:8])
	binary.Read(buffer_len, binary.BigEndian, &head_len)
	binary.Read(buffer_pid, binary.BigEndian, &head_pid)

	return head_len, head_pid
}
func SendMsgRegister(conn net.Conn, nick string, pwd string) {
	register := &protocol.Account_RegisterPlayer{
		Playername: proto.String(nick),
		Passworld:  proto.String(pwd),
	}

	encObj, err := proto.Marshal(register)
	is_ok := CheckError(err)
	if is_ok {
		SendPackage(&conn, 1, encObj)
	}
}

func SendServerList(conn net.Conn) {
	serverList := &protocol.Account_ServerListResult{}
	encObj, err := proto.Marshal(serverList)
	is_ok := CheckError(err)
	if is_ok {
		SendPackage(&conn, 3, encObj)
	}
}

func SenMsgLogin(conn net.Conn) {
	//登陆相关
	loginInfo := &protocol.Account_LoginInfo{
		Playername: proto.String(nick),
		Passworld:  proto.String(pwd),
	}

	encObj, err := proto.Marshal(loginInfo)
	is_ok := CheckError(err)
	if is_ok {
		SendPackage(&conn, 2, encObj)
	}
}

func ReciveResult(conn net.Conn) {
	const MAXLEN = 1024
	buf := make([]byte, MAXLEN)

	for true {
		n, _ := conn.Read(buf) //接收具体消息
		_, head_pid := GetHead(buf)

		switch head_pid {
		case 2:
			result := new(protocol.Account_LoginResult)
			if err := proto.Unmarshal(buf[8:n], result); err == nil {
				switch result.GetResult() {
				case global.LOGINSUCCESS:
					log.Info("login sucessfull and player id=%d gameserver = %s", result.GetPlayerId(), result.GetGameserver())
				default:
					log.Error("login error")
				}
				conn.Close()
				return
			} else {
				fmt.Println(err)
			}
		case 1:
			result := new(protocol.Account_RegisterResult)
			if err := proto.Unmarshal(buf[8:n], result); err == nil {
				switch result.GetResult() {
				case global.REGISTERSUCCESS:
					log.Trace("register sccessfull!")
				default:
					log.Error("register error")
				}
				//注册后登陆
				SenMsgLogin(conn)
			} else {
				fmt.Println("err:", err)
			}

		case 3:
			result := new(protocol.Account_ServerListResult)
			if err := proto.Unmarshal(buf[8:n], result); err == nil {
				fmt.Println("server list:", len(result.ServerInfo))
			}
		}

		time.Sleep(5 * time.Millisecond)
	}
}

type ConnStruct struct {
	conn net.Conn
}

//测试账号服务器
func testAccount(nick string, pwd string, conn net.Conn) {
	go ReciveResult(conn)
	SendServerList(conn)             //获取服务器列表
	SendMsgRegister(conn, nick, pwd) //注册登录
}

func ReciveResult4Game(conn *net.Conn) {
	const MAXLEN = 2048
	buff := make([]byte, MAXLEN)

	for true {
		n, _ := (*conn).Read(buff) //接收具体消息

		fmt.Println("receve:", buff[0:n], "n:", n)
		if n < 8 {
			return
		}

		for n >= 8 {
			body_len, head_pid := GetHead(buff[:8])
			buf := buff[:body_len]
			buff = buff[body_len:]
			n = n - int(body_len)

			switch head_pid {
			case 1001:
				result := new(protocol.Game_RegisterRoleResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("result:", result.GetResult())
				}

			case 1002:
				result := new(protocol.Game_RoleInfoResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("result.HangupLevels:", result.HangupLevels, "copy_levels:", result.CopyLevels, "result.GetIsCreate():", result.GetIsCreate())
					if !result.GetIsCreate() {
						testRegisterGame(conn)
					} else {
						/*
								var id int32 = 10101 //副本关卡
								base := new(protocol.Game_WarMapStage)
								base.StageId = &id
								encObj, _ := proto.Marshal(base)
								SendPackage(conn, 1003, encObj)

							//发送扫荡
							var stage_id int32 = 10101
							var count int32 = 10

							base := new(protocol.Game_SweepMapStage)
							base.StageId = &stage_id
							base.Count = &count

							encObj, err := proto.Marshal(base)
							is_ok := CheckError(err)
							if is_ok {
								fmt.Println("send 战斗结果")
								SendPackage(conn, 1005, encObj)
							}*/

						//发送请求挂机
						/*
							var id int32 = 20101
							base := new(protocol.GameGetGuajiInfo)
							base.Id = &id
							encObj, _ := proto.Marshal(base)
							SendPackage(conn, 1016, encObj)
						*/

						//挑战boss
						var id int32 = 20101
						base := new(protocol.Game_ChallengeBoss)
						base.Id = &id
						encObj, _ := proto.Marshal(base)
						SendPackage(conn, 1017, encObj)

						/*
							var id int32 = 20101
							var state int32 = 1
							var type_ int32 = 2

							stage := new(protocol.Game_Stage)
							stage.Type = &type_
							stage.State = &state
							stage.StageId = &id

							base := new(protocol.Game_C2SChallenge)
							base.Stage = stage
							encObj, _ := proto.Marshal(base)
							fmt.Println(encObj)
							SendPackage(conn, 1018, encObj)
						*/

						/*
							//切换关卡
							var id int32 = 20102
							base := new(protocol.Game_ChangeGuajiInfo)
							base.Id = &id
							encObj, _ := proto.Marshal(base)
							SendPackage(conn, 1020, encObj)
						*/
						/*
									//快速战斗
									var type_ int32 = 2
									var state int32 = 1
									var stage_id int32 = 20102

									stage := new(protocol.Game_Stage)
									stage.Type = &type_
									stage.State = &state
									stage.StageId = &stage_id

									base := new(protocol.Game_FastWar)
									base.Stage = stage
									encObj, _ := proto.Marshal(base)
									SendPackage(conn, 1021, encObj)


								//发送扫荡

								var stage_id int32 = 10101
								var count int32 = 10

								base := new(protocol.Game_SweepMapStage)
								base.StageId = &stage_id
								base.Count = &count

								encObj, err := proto.Marshal(base)
								is_ok := CheckError(err)
								if is_ok {
									fmt.Println("send 战斗结果")
									SendPackage(conn, 1005, encObj)
								}
								/*
									//获取该关卡的挂机列表
									base := new(protocol.Game_GetGuajiRoleList)
									encObj, _ := proto.Marshal(base)
									SendPackage(conn, 1022, encObj)

									//请求阵型
									base := new(protocol.Game_GetGuajiRoleFormation)
									var role_id int64 = 1000003
									var type_ int32 = 2
									base.RoleId = &role_id
									base.Type = &type_

									encObj, _ := proto.Marshal(base)
									SendPackage(conn, 1023, encObj)
									/*
											//上下阵
											base := new(protocol.Game_HerosFormation)
											var type_ int32 = 2
											var is_on bool = true
											var pos_id int32 = 1
											var hero_uid int32 = 301648466

											base.Type = &type_
											base.IsOn = &is_on
											base.PosId = &pos_id
											base.HeroUid = &hero_uid

											encObj, _ := proto.Marshal(base)
											SendPackage(conn, 1024, encObj)

										//交换阵型
										base := new(protocol.Game_ChangeHerosFormation)
										var type_ int32 = 2
										var pos_id_1 int32 = 1
										var pos_id_2 int32 = 2

										base.Type = &type_
										base.PosId_1 = &pos_id_1
										base.PosId_2 = &pos_id_2

										encObj, _ := proto.Marshal(base)
										SendPackage(conn, 1025, encObj)

							var uid int32 = 578459847
							var count int32 = 21
							base := new(protocol.Game_UseProp)
							base.Uid = &uid
							base.Count = &count
							encObj, _ := proto.Marshal(base)
							SendPackage(conn, 1027, encObj)

							//查看道具背包
							base2 := new(protocol.Game_CheckPropBag)
							encObj2, _ := proto.Marshal(base2)
							SendPackage(conn, 1029, encObj2)

							//在线收益
							/*base := new(protocol.Game_C2SOnlineGuaji)
							encObj, _ := proto.Marshal(base)
							SendPackage(conn, 1015, encObj)
						*/

						/*
							base := new(protocol.Game_CheckEquipBag)
							encObj, _ := proto.Marshal(base)
							SendPackage(conn, 1030, encObj)
						*/
					}
				}
			case 1003:
				result := new(protocol.Game_MapStageResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("1003:", result.GetResult())
					if result.GetResult() == 0 { //允许通关

						var id int32 = 10101
						var state int32 = 1
						var type_ int32 = 1

						stage := new(protocol.Game_Stage)
						stage.StageId = &id
						stage.Type = &type_
						stage.State = &state

						base := new(protocol.Game_WarMapNoteServer)
						base.Stage = stage
						encObj, _ := proto.Marshal(base)
						SendPackage(conn, 1004, encObj) //发送战斗结构

					}
				}
			case 1004: //战斗结果
				result := new(protocol.Game_WarMapNoteServerResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("1004", result.GetReward().GetPlayerGold(), result.GetReward().GetHeroExp(), result.GetReward().GetEquipUids(), result.GetReward().GetPropUids())
				}

				//发送扫荡
				var stage_id int32 = 10101
				var count int32 = 1

				base := new(protocol.Game_SweepMapStage)
				base.StageId = &stage_id
				base.Count = &count

				encObj, err := proto.Marshal(base)
				is_ok := CheckError(err)
				if is_ok {
					fmt.Println("send 战斗结果")
					SendPackage(conn, 1005, encObj)
				}

			case 1005: //扫荡结果
				result := new(protocol.Game_SweepMapStageResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("1005", result.GetResult(), result.GetReward())
				}

			case 1006: //道具变化
				result := new(protocol.Game_Notice2CProp)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("1006", result.GetType(), result.GetProp())
				}
			case 1008: //装备
				result := new(protocol.Game_Notice2CEquip)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("1008", result.GetType(), result.GetEquip())
				}
			case 1010: //得到的物品
				result := new(protocol.Game_Notice2CRoleInfo)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("1010", result.GetLevel(), result.GetHp(), result.GetPower(), result.GetHp())
				}

			case 1011: //得到的物品
				result := new(protocol.Game_Notice2CMoney)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("1011", result.GetGold(), result.GetDiamond())
				}

			case 1012: //体力变化
				result := new(protocol.Game_Notice2CEnergy)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("1012", result.GetEnergy(), result.GetEnergyMax())
				}

			case 1013: //关卡变化
				result := new(protocol.Game_Notice2CheckPoint)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("1013", result.GetLatestCheckpoint().GetType(), result.GetLatestCheckpoint().GetState(), result.GetLatestCheckpoint().GetStageId())
				}

			case 1014: //离线收益
				result := new(protocol.Game_OffNotice2CGuaji)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("1014", result.GetPointId(), result.GetGold(), result.GetExp(), result.GetGuajiTime(), result.GetKillNpcNum())
				}

			case 1015: //在线收益
				result := new(protocol.Game_OnNotice2CGuaji)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("1015", result.GetGuajiType(), result.GetNpcId(), result.GetGold(), result.GetExp())
				}

			case 1016: //请求挂机返回
				result := new(protocol.Game_GuajiInfoResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("1016", result.GetConditions())
				}
			case 1017: //可否挑战boss返回
				result := new(protocol.Game_ChallengeBossResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("1017", result.GetIsCanChange(), result.GetTeam_2())
				}
			case 1018: //挑战挂机boss
				result := new(protocol.Game_C2SChallengeResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("1018", result.GetEquipUids(), result.GetPropUids())
				}
			case 1020: //关卡切换
				result := new(protocol.Game_ChangeGuajiInfoResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("1020", result.GetIsOk())
				}

			case 1021: //快速战斗
				result := new(protocol.Game_FastWarResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("1021", result.GetResult(), result.GetReward().GetPlayerGold())
				}
			case 1022: //挂机列表返回
				result := new(protocol.Game_GetGuajiRoleListResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("1022", result.GetGuajiRoleInfos())
				}

			case 1023: //获取某玩家阵型
				result := new(protocol.Game_GetGuajiRoleFormationResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("1023", result.GetResult(), result.GetFormations())
				}
			case 1024: //英雄上下阵
				result := new(protocol.Game_HerosFormationResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("1024", result.GetResult())
				}
			case 1025: //交换阵型
				result := new(protocol.Game_ChangeHerosFormationResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("1025", result.GetResult())
				}
			case 1027: //使用某个道具
				result := new(protocol.Game_UsePropResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("1027", result.GetResult())
				}
			case 1029: //查看道具背包
				result := new(protocol.Game_CheckPropBagResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("1029", result.GetPropStruct())
				}
			case 1030: //查看装备背包
				result := new(protocol.Game_CheckEquipBagResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("1030", result.GetEquipStruct())
				}
			default:
			}
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func testRegisterGame(conn *net.Conn) {

	var player_id int32 = 1
	var hero_id int32 = 901
	var nick string = "ttt"

	base := new(protocol.Game_RegisterRole)
	base.PlayerId = &player_id
	base.Nick = &nick
	base.HeroId = &hero_id

	encObj, err := proto.Marshal(base)
	is_ok := CheckError(err)
	if is_ok {
		fmt.Println("send RegisterGame")
		SendPackage(conn, 1001, encObj)
	}
}

func testGetInfoGame(conn net.Conn) {
	go ReciveResult4Game(&conn)

	role_info := new(protocol.Game_GetRoleInfo)
	var playerid int32 = 1
	role_info.PlayerId = &playerid
	encObj, err := proto.Marshal(role_info)
	fmt.Println(encObj)
	is_ok := CheckError(err)
	if is_ok {
		SendPackage(&conn, 1002, encObj)
	}
}

func main() {
	/*
		conn, err := net.Dial("tcp", "127.0.0.1:8080") //121.52.235.141:8080
		if err != nil {
			fmt.Println(err)
		}
		go testAccount(nick, pwd, conn)
	*/

	conn, err := net.Dial("tcp", "127.0.0.1:8082")
	if err != nil {
		fmt.Println(err)
	}
	testGetInfoGame(conn)

	time.Sleep(1000000 * time.Millisecond)
}
