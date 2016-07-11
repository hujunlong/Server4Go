package main

import (
	"server/share/protocol"

	"github.com/golang/protobuf/proto"
)

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

//初始化
func init() {

}

//字符串转int32
func Str2Int32(str string) int32 {
	str = strings.TrimSpace(str)
	data_int, error := strconv.Atoi(str)
	if error != nil {
		return 0
	}
	return int32(data_int)
}

//错误检查
func CheckError(err error) bool {
	if err != nil {
		fmt.Println("错误消息:", err)
		return false
	}
	return true
}

func Str2int32(str string) int32 {
	str = strings.TrimSpace(str)

	if len(str) <= 0 {
		return 0
	}

	a2_i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return int32(a2_i)
}

//数据发包
func SendPackage(conn net.Conn, pid int, body []byte) {
	var pid_32 int32 = int32(pid)

	len := 8 + len(body)
	var len_32 = int32(len)

	len_buf := bytes.NewBuffer([]byte{})
	binary.Write(len_buf, binary.BigEndian, &len_32)

	pid_buf := bytes.NewBuffer([]byte{})
	binary.Write(pid_buf, binary.BigEndian, &pid_32)

	msg := append(len_buf.Bytes(), pid_buf.Bytes()...)
	msg2 := append(msg, body...)
	conn.Write(msg2)
	fmt.Println("发送数据 pid=", pid)
}

//获取包头
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

//账号服务器接收消息
func ReciveAccountResult(conn net.Conn) {
	const MAXLEN = 1024
	buf := make([]byte, MAXLEN)

	for true {
		n, _ := conn.Read(buf)
		_, head_pid := GetHead(buf)

		if n > MAXLEN || n < 8 {
			//fmt.Println("recive error n> MAXLEN")
			return
		}

		switch head_pid {
		case 1:
			result := new(protocol.Account_RegisterResult)
			if err := proto.Unmarshal(buf[8:n], result); err == nil {
				switch result.GetResult() {
				case 0:
					fmt.Println("注册成功\n")
				case 1:
					fmt.Println("注册失败\n")
				case 4:
					fmt.Println("注册名被占用\n")
				default:
				}
			}

		case 2:
			result := new(protocol.Account_LoginResult)
			//fmt.Println("buf:", buf, "n:", n)
			if err := proto.Unmarshal(buf[8:n], result); err == nil {
				switch result.GetResult() {
				case 5:
					fmt.Println("登录成功\t", "player id:", result.GetPlayerId(), "游戏地址服务器:", result.GetGameserver(), "\n")
					result.GetPlayerId()
					fout, _ := os.Create("account.txt")
					defer fout.Close()
					fout.WriteString(strconv.Itoa(int(result.GetPlayerId())))
					conn.Close()
				case 2:
					fmt.Println("登录失败\n")
				case 6:
					fmt.Println("禁止登录\n")
				default:
				}

			}

		case 3:
			result := new(protocol.Account_ServerListResult)
			if err := proto.Unmarshal(buf[8:n], result); err == nil {
				fmt.Println("服务器列表", result.GetServerInfo())
			}
		}

		time.Sleep(5 * time.Millisecond)
	}
}

//游戏服务器接收数据消息
func ReciveResult4Game(conn *net.Conn) {
	const MAXLEN = 20480
	buff := make([]byte, MAXLEN)

	for true {
		n, _ := (*conn).Read(buff)

		if n < 8 {
			return
		}

		for n >= 8 {
			body_len, head_pid := GetHead(buff[:8])
			buf := buff[:body_len]
			buff = buff[body_len:]
			n = n - int(body_len)

			fmt.Println("接收消息id=", head_pid)
			switch head_pid {
			case 1001:
				result := new(protocol.PlayerBase_RegisterRoleResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					if result.GetResult() == 100 {
						fmt.Println("result 1001 角色注册成功\n")
					} else {
						fmt.Println("result 1001 角色注册失败:", result.GetResult(), "\n")
					}
				}

			case 1002:
				result := new(protocol.PlayerBase_RoleInfoResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("是否创建角色", result.GetIsCreate(), "\n基础信息:", result.GetPlayerInfo(), "\n英雄列表:", result.GetHeros(), "\n装备列表:", result.GetEquips(), "\n道具列表:", result.GetProps(), "\n副本列表:", result.GetCopyLevels(), "\n挂机列表:", result.GetHangupLevels(), "\n职业id", result.GetProfessionId())
				}

			case 1203:
				result := new(protocol.NoticeMsg_Notice2CRoleInfo)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("玩家基础属性变化 level:", result.GetLevel(), "exp:", result.GetExp(), "power:", result.GetPower(), "hp:", result.GetHp())
				}

			case 1204:
				result := new(protocol.NoticeMsg_Notice2CMoney)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("钱变化 gold:", result.GetGold(), " diamond:", result.GetDiamond())
				}
			case 1205:
				result := new(protocol.NoticeMsg_Notice2CEnergy)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("体力变化 energy:", result.GetEnergy(), " EnergyMax:", result.GetEnergyMax())
				}
			case 1506: //道具出售
				result := new(protocol.Goods_SalePropResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("道具出售:", result.GetResult())
				}
			case 1601:
				result := new(protocol.GM_MsgResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("GM返回:", result.GetResult())
				}
			case 1701:
				result := new(protocol.Activity_NiuDanResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("抽卡返回:", result.GetResult(), result.GetProps(), result.GetGold(), result.GetHeroUids())
				}
			case 1212: //推送新添加英雄
				result := new(protocol.NoticeMsg_NoticeGetHeros)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("收到推送新英雄:", result.GetHeros())
				}
			case 1801: //当前满足任务
				result := new(protocol.Task_AllTaskResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("当前所有满足任务:", result.GetTasks())
				}
			case 1214: //推送新任务
				result := new(protocol.NoticeMsg_Notice2CNewTask)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("推送新任务:", result.GetTask())
				}
			case 1215: //推送当前可以接受的手动接取任务id
				result := new(protocol.NoticeMsg_Notice2CCanAccept)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("推送当前可以接受的手动接取任务id:", result.GetCanAcceptTasks())
				}
			case 1216: //当前任务变化(用来更新任务进度)
				result := new(protocol.NoticeMsg_Notice2CUpdateProgress)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("当前任务变化(用来更新任务进度):", result.GetTypeInfo(), result.GetSubitem())
				}
			case 1217:
				result := new(protocol.NoticeMsg_SubmitTask)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("推送自动提交任务所获取的奖励:", result.GetTypeInfo(), result.GetGold(), result.GetRoleExp(), result.GetProps(), result.GetEquipUids())
				}
			case 1218:
				result := new(protocol.NoticeMsg_Notice2CHandleTask)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("客户端需手动提交任务获取奖励:", result.GetTypeInfo())
				}
			case 1219:
				result := new(protocol.NoticeMsg_Notice2CAchievementChange)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("成就进度变化:", result.GetInfos())
				}

			case 1803:
				result := new(protocol.Task_SubmitTaskResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("手动提交任务获取奖励:", result.GetResult(), result.GetGold(), result.GetRoleExp(), result.GetProps(), result.GetEquipUids(), result.GetHeroUids())
				}
			case 1821:
				result := new(protocol.Task_GetXuanShangInfoResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("获取悬赏任务相关:", result.GetLastNum(), result.GetFreeLastTime())
				}
			case 1822:
				result := new(protocol.Task_XuanShangDiamondRefResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("获取悬赏任务相关:", result.GetResult())
				}
			case 1824:
				result := new(protocol.Task_GiveUpTaskResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("任务放弃:", result.GetResult())
				}
			case 1851: //成就
				result := new(protocol.Task_GetAchievementResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("成就:", result.GetInfos())
				}
			case 1852: //请求领取成就奖励
				result := new(protocol.Task_GetAchievementRewardResult)
				if err := proto.Unmarshal(buf[8:body_len], result); err == nil {
					fmt.Println("成就奖励:", result.GetResult(), result.GetProps())
				}
			}
		}
	}
}

//账号注册消息
func SendMsgRegister(nick string, pwd string, conn net.Conn) {
	register := &protocol.Account_RegisterPlayer{
		Playername: proto.String(nick),
		Passworld:  proto.String(pwd),
	}

	encObj, err := proto.Marshal(register)
	is_ok := CheckError(err)
	if is_ok {
		SendPackage(conn, 1, encObj)
	}
}

//账号登录
func SenMsgLogin(nick string, pwd string, conn net.Conn) {
	loginInfo := &protocol.Account_LoginInfo{
		Playername: proto.String(nick),
		Passworld:  proto.String(pwd),
	}

	encObj, err := proto.Marshal(loginInfo)
	is_ok := CheckError(err)
	if is_ok {
		SendPackage(conn, 2, encObj)
	}
}

//服务器列表
func GetGameServerList(conn net.Conn) {
	loginInfo := &protocol.Account_GetServerList{}
	encObj, err := proto.Marshal(loginInfo)
	is_ok := CheckError(err)
	if is_ok {
		SendPackage(conn, 3, encObj)
	}
}

func read(path string) (int32, error) {
	fi, err := os.Open(path)
	if err != nil {
		return 0, err
	}

	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	str := string(fd)
	data, _ := strconv.Atoi(str)
	return int32(data), nil
}

//游戏服务器注册
func RegisterGame(conn *net.Conn) {

	player_id, err := read("account.txt")
	if err != nil {
		fmt.Println("游戏注册失败 player id 没有写入account.txt文件中\n")
	}
	var hero_id int32 = 901
	var nick string = "ttt"

	base := new(protocol.PlayerBase_RegisterRole)
	base.PlayerId = &player_id
	base.Nick = &nick
	base.HeroId = &hero_id

	encObj, err := proto.Marshal(base)
	is_ok := CheckError(err)
	if is_ok {
		SendPackage(*conn, 1001, encObj)
	}
}

//游戏服务器登录
func LoginGameServer(conn *net.Conn) {
	player_id, err := read("account.txt")
	if err != nil {
		fmt.Println("player id 未知,不能登录\n")
	}
	base := new(protocol.PlayerBase_GetRoleInfo)
	base.PlayerId = &player_id
	encObj, err := proto.Marshal(base)
	is_ok := CheckError(err)
	if is_ok {
		SendPackage(*conn, 1002, encObj)
	}
}

//GM命令
func GMmessge(str string, conn *net.Conn) {
	base := new(protocol.GM_Msg)
	base.Msg = &str
	encObj, _ := proto.Marshal(base)
	SendPackage(*conn, 1601, encObj)
}

//扭蛋
func NiuDan(type_ int32, type_group int32, conn *net.Conn) {
	base := new(protocol.Activity_NiuDan)
	base.Type = &type_
	base.TypeGroup = &type_group
	encObj, _ := proto.Marshal(base)
	SendPackage(*conn, 1701, encObj)
}

//获取所有任务列表
func AllTask(conn *net.Conn) {
	base := new(protocol.Task_AllTask)
	encObj, _ := proto.Marshal(base)
	SendPackage(*conn, 1801, encObj)
}

//手动获取奖励
func SubmitTask(type_id int32, id int32, conn *net.Conn) {
	taskInfo := new(protocol.TaskType)
	taskInfo.Type = &type_id
	taskInfo.Id = &id

	base := new(protocol.Task_SubmitTask)
	base.TaskInfo = taskInfo

	encObj, _ := proto.Marshal(base)
	SendPackage(*conn, 1803, encObj)
}

//获取悬赏任务相关
func GetXuanShangInfo(conn *net.Conn) {
	base := new(protocol.Task_GetXuanShangInfo)
	encObj, _ := proto.Marshal(base)
	SendPackage(*conn, 1821, encObj)
}

//获取悬赏任务相关
func XuanShangDiamondRef(conn *net.Conn) {
	base := new(protocol.Task_XuanShangDiamondRef)
	encObj, _ := proto.Marshal(base)
	SendPackage(*conn, 1822, encObj)
}

//手动接取任务
func AcceptTask(num1 int32, num2 int32, conn *net.Conn) {
	taskInfo := new(protocol.TaskType)
	taskInfo.Type = &num1
	taskInfo.Id = &num2

	base := new(protocol.Task_AcceptTask)
	base.TaskInfo = taskInfo
	encObj, _ := proto.Marshal(base)
	SendPackage(*conn, 1802, encObj)
}

//任务放弃
func GiveUpTask(num1 int32, num2 int32, conn *net.Conn) {
	taskInfo := new(protocol.TaskType)
	taskInfo.Type = &num1
	taskInfo.Id = &num2

	base := new(protocol.Task_GiveUpTask)
	base.Task = taskInfo
	encObj, _ := proto.Marshal(base)
	SendPackage(*conn, 1824, encObj)
}

//请求获取成就列表
func GetAchievement(conn *net.Conn) {
	base := new(protocol.Task_GetAchievement)
	encObj, _ := proto.Marshal(base)
	SendPackage(*conn, 1851, encObj)
}

//请求领取成就奖励
func GetAchievementReward(id int32, conn *net.Conn) {
	base := new(protocol.Task_GetAchievementReward)
	base.Id = &id
	encObj, _ := proto.Marshal(base)
	SendPackage(*conn, 1852, encObj)
}

//道具出售
func SaleProp(uid int32, count int32, conn *net.Conn) {
	base := new(protocol.Goods_SaleProp)
	base.PropUid = &uid
	base.Count = &count
	encObj, _ := proto.Marshal(base)
	SendPackage(*conn, 1506, encObj)
}

func main() {

	var cmd string = ""

	//账号服务器
	conn1, err1 := net.Dial("tcp", "127.0.0.1:8080") //
	if err1 != nil {
		fmt.Println("账号服务器连接错误:", err1)
	}
	go ReciveAccountResult(conn1)

	//游戏服务器
	conn2, err2 := net.Dial("tcp", "127.0.0.1:8082") //121.52.235.141
	if err2 != nil {
		fmt.Println("游戏服务器连接错误:", err2)
	}
	go ReciveResult4Game(&conn2)

	inputReader := bufio.NewReader(os.Stdin)
	for {
		cmd, _ = inputReader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)

		cmd_id := Str2int32(cmd)

		switch cmd_id {
		case 1: //注册账号
			fmt.Println("请输入用注册户名,密码\n")
			cmd, _ = inputReader.ReadString('\n')
			strs := strings.Split(cmd, " ")
			SendMsgRegister(strs[0], strs[1], conn1)
		case 2: //登录
			fmt.Println("请输入用登录户名,密码\n")
			cmd, _ = inputReader.ReadString('\n')
			strs := strings.Split(cmd, " ")
			SenMsgLogin(strs[0], strs[1], conn1)
		case 3: //服务器列表
			GetGameServerList(conn1)
		case 1001: //游戏服务器注册
			RegisterGame(&conn2)
		case 1002: //登录
			LoginGameServer(&conn2)
		case 1506: //道具出售
			fmt.Println("请输入道具uid 与 count")
			cmd, _ = inputReader.ReadString('\n')
			strs := strings.Split(cmd, " ")
			num1 := Str2Int32(strs[0])
			num2 := Str2Int32(strs[1])
			SaleProp(num1, num2, &conn2)
		case 1601: //GM命令
			fmt.Println("请输入GM命令\n")
			cmd, _ = inputReader.ReadString('\n')
			GMmessge(cmd, &conn2)
		case 1701: //活动抽卡
			fmt.Println("活动抽卡 参数1（1:单抽 2:10连抽) 参数2(抽卡类型)看配置\n")
			cmd, _ = inputReader.ReadString('\n')
			strs := strings.Split(cmd, " ")
			num1 := Str2Int32(strs[0])
			num2 := Str2Int32(strs[1])
			NiuDan(num1, num2, &conn2)
		case 1801:
			AllTask(&conn2)
		case 1803: //手动获取奖励
			fmt.Println("请输入任务类型(1:主线 2:悬赏) id(csv中的配置)\n")
			cmd, _ = inputReader.ReadString('\n')
			strs := strings.Split(cmd, " ")
			num1 := Str2Int32(strs[0])
			num2 := Str2Int32(strs[1])
			SubmitTask(num1, num2, &conn2)
		case 1821: //获取悬赏任务相关
			GetXuanShangInfo(&conn2)
		case 1822: //悬赏元宝刷新请求
			XuanShangDiamondRef(&conn2)
		case 1802: //手动接受任务
			fmt.Println("请输入手动接取 任务类型与任务id")
			cmd, _ = inputReader.ReadString('\n')
			strs := strings.Split(cmd, " ")
			num1 := Str2Int32(strs[0])
			num2 := Str2Int32(strs[1])
			AcceptTask(num1, num2, &conn2)
		case 1824: //放弃任务
			fmt.Println("放弃任务 任务类型与任务id")
			cmd, _ = inputReader.ReadString('\n')
			strs := strings.Split(cmd, " ")
			num1 := Str2Int32(strs[0])
			num2 := Str2Int32(strs[1])
			GiveUpTask(num1, num2, &conn2)
		case 1851: //请求获取成就列表
			GetAchievement(&conn2)
		case 1852: //成就领取奖励
			fmt.Println("输入成就id")
			cmd, _ = inputReader.ReadString('\n')
			num := Str2Int32(cmd)
			GetAchievementReward(num, &conn2)

		default:
		}
	}
}
