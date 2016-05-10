package main

import (
	"server/share/global"
	"server/share/protocol"

	//"encoding/json"

	"github.com/game_engine/logs"
	"github.com/golang/protobuf/proto"
)

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"time"
)

var log *logs.BeeLogger

const max_client = 1

var end = make(chan int)

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
	fmt.Println(msg2)
}

func GetHead(buf []byte) (int32, int32) {
	var head_len int32 = 0
	var head_pid int32 = 0
	buffer_len := bytes.NewBuffer(buf[0:4])
	buffer_pid := bytes.NewBuffer(buf[4:8])
	binary.Read(buffer_len, binary.BigEndian, &head_len)
	binary.Read(buffer_pid, binary.BigEndian, &head_pid)

	return head_len, head_pid
}
func SendMsgRegister(conn net.Conn, i int) {
	nick := strconv.Itoa(i)
	register := &protocol.Account_RegisterPlayer{
		Playername: proto.String(nick),
		Passworld:  proto.String(nick),
	}

	encObj, err := proto.Marshal(register)
	is_ok := CheckError(err)
	if is_ok {
		SendPackage(conn, 1, encObj)
	}
}

func SendServerList(conn net.Conn) {
	serverList := &protocol.Account_ServerListResult{}
	encObj, err := proto.Marshal(serverList)
	is_ok := CheckError(err)
	if is_ok {
		SendPackage(conn, 3, encObj)
	}
}

func SenMsgLogin(conn net.Conn, i int) {
	nick := strconv.Itoa(i)
	//登陆相关
	loginInfo := &protocol.Account_LoginInfo{
		Playername: proto.String(nick),
		Passworld:  proto.String(nick),
	}

	encObj, err := proto.Marshal(loginInfo)
	is_ok := CheckError(err)
	if is_ok {
		SendPackage(conn, 2, encObj)
	}
}

func ReciveResult(conn net.Conn, i int, recive_result chan int) {
	const MAXLEN = 1024
	buf := make([]byte, MAXLEN)

	for true {
		n, _ := conn.Read(buf) //接收具体消息
		//接收包的type类型用来区分包之间的区别
		_, head_pid := GetHead(buf)

		switch head_pid {
		case 2:
			result := new(protocol.Account_LoginResult)
			if err := proto.Unmarshal(buf[8:n], result); err == nil {
				switch result.GetResult() {
				case global.LOGINSUCCESS:
					SenMsgLogin(conn, i)
					log.Info("login sucessfull and player id=%d gameserver = %s", result.GetPlayerId(), result.GetGameserver())
				default:
					log.Error("login error")
				}

				conn.Close()
				recive_result <- 1
				if i == max_client-1 {
					end <- 1
				}
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
				SenMsgLogin(conn, i)
			} else {
				fmt.Println("err:", err)
			}

		case 3:
			result := new(protocol.Account_ServerListResult)
			if err := proto.Unmarshal(buf[8:n], result); err == nil {
				fmt.Println("result:", len(result.ServerInfo))
			}
		}

		time.Sleep(5 * time.Millisecond)
	}
}

func MessageRun(conn net.Conn, i int) {
	//通信先获取返回数据
	result := make(chan int)
	go ReciveResult(conn, i, result)
	SendMsgRegister(conn, i)
	//SenMsgLogin(conn, i)
	//SendServerList(conn)
	<-result
}

type ConnStruct struct {
	conn net.Conn
}

//测试账号服务器
func testAccount() {
	var arrayConnStruct [max_client]ConnStruct
	var err error
	for i := 0; i < max_client; {

		arrayConnStruct[i].conn, err = net.Dial("tcp", "10.8.2.172:8080") //121.52.235.141:8080
		if err != nil {
			log.Error("connect error %s", err)
			time.Sleep(100 * time.Millisecond)
		} else {
			go MessageRun(arrayConnStruct[i].conn, i)
			time.Sleep(5 * time.Millisecond)
			i++
		}
	}
}

func ReciveResult4Game(conn net.Conn) {
	const MAXLEN = 1024
	buf := make([]byte, MAXLEN)

	for true {
		n, _ := conn.Read(buf) //接收具体消息
		fmt.Println("recive Result:", buf[0:n])
		//接收包的type类型用来区分包之间的区别
		_, head_pid := GetHead(buf)
		fmt.Println(buf[0:n])

		switch head_pid {
		case 1002:
			result := new(protocol.Game_RoleInfoResult)
			if err := proto.Unmarshal(buf[8:n], result); err == nil {
				fmt.Println(result.GetIsCreate(), result.PlayerInfo.GetNick(), result.HeroStruct)
			}

		case 1001:
			result := new(protocol.Game_RegisterRoleResult)
			if err := proto.Unmarshal(buf[8:n], result); err == nil {
				fmt.Println("result:", result.GetResult())
			}
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func testRegisterGame() {
	conn, err := net.Dial("tcp", "10.8.2.172:8082")
	go ReciveResult4Game(conn)

	var player_id int32 = 1
	var hero_id int32 = 901
	var nick string = "aaa"

	base := new(protocol.Game_RegisterRole)
	base.PlayerId = &player_id
	base.Nick = &nick
	base.HeroId = &hero_id

	encObj, err := proto.Marshal(base)
	is_ok := CheckError(err)
	if is_ok {
		SendPackage(conn, 1001, encObj)
	}
}

func testGetInfoGame() {
	conn, err := net.Dial("tcp", "127.0.0.1:8082")
	go ReciveResult4Game(conn)

	role_info := new(protocol.Game_GetRoleInfo)
	var playerid int32 = 2
	role_info.PlayerId = &playerid
	encObj, err := proto.Marshal(role_info)
	fmt.Println(encObj)
	is_ok := CheckError(err)
	if is_ok {
		SendPackage(conn, 1002, encObj)
	}
}

func LoginAccount() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	go ReciveResult4Game(conn)

	var player_id int32 = 3
	var hero_id int32 = 901
	var nick string = "aaa"

	base := new(protocol.Game_RegisterRole)
	base.PlayerId = &player_id
	base.Nick = &nick
	base.HeroId = &hero_id

	encObj, err := proto.Marshal(base)
	is_ok := CheckError(err)
	if is_ok {
		SendPackage(conn, 1001, encObj)
	}
}

type LoginBase struct {
	max int32
	min int32
}

func dojson() {
	/*
		cnf, err := config.NewConfig("json", "test.json")
		ErrorMsg(err)
		fmt.Println(cnf)

		rootArray, _ := cnf.DIY("rootArray")
		rootArrayCasted := rootArray.([]interface{})
		if len(rootArrayCasted) <= 0 {
			fmt.Println("error config")
			return
		}

		for i := 0; i < len(rootArrayCasted); i++ {
			elem := rootArrayCasted[i].(map[string]interface{})

			if elem["name"] == float64(1) {
				fmt.Println("come here")
			}

			fmt.Println(elem["id"])
			fmt.Println(elem["name"])
			fmt.Println(elem["discri"])
			fmt.Println(elem["gold"])
		}
	*/
}

func main() {
	//testAccount()
	//go testRegisterGame()
	go testGetInfoGame()
	//dojson()
	time.Sleep(100000 * time.Millisecond)
}
