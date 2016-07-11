package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"server/loginServer/account"
	"server/share/global"
	"server/share/protocol"

	"github.com/game_engine/cache/redis"
	"github.com/golang/protobuf/proto"
)

type Deal4C struct {
	account_info *account.AccountInfo
}

func (this *Deal4C) Init() {
	this.account_info = new(account.AccountInfo)
	this.account_info.Init()
}

func (this *Deal4C) Deal4Client(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		fmt.Println(conn.RemoteAddr(), "connected")
		if CheckError(err) {
			go this.Handler4C(conn)
		}
	}
}

func (this *Deal4C) getNewAddress() (int32, string, error) {
	var address string = ""
	var max_count int32 = 999999
	var address_id int32 = 0

	if len(config.NewServerAddress) == 0 {
		account.Log.Error("config.newServerAddress len = 0")
		return 0, "", errors.New("config.newServerAddress len = 0")
	}

	for key, _ := range config.NewServerAddress {
		if deal_4g.GameConnects[key].Count <= max_count {
			address_id = deal_4g.GameConnects[key].ServerId
			address = deal_4g.GameConnects[key].Address
			max_count = deal_4g.GameConnects[key].Count
		}
	}
	return address_id, address, nil
}

func (this *Deal4C) NoteGame(player_id int32, game_id int32) error {
	fmt.Println("发送通知游戏服务器 函数进入")
	result4G := &protocol.Account_NoteGame{
		PlayerId: proto.Int32(int32(player_id)),
	}

	encObj, _ := proto.Marshal(result4G)
	if _, ok := deal_4g.GameConnects[game_id]; ok {
		SendPackage(deal_4g.GameConnects[game_id].Conn, 102, encObj)
		fmt.Println("通知游戏服务器send to game", encObj)
		return nil
	}
	return errors.New("game connect error")
}

func (this *Deal4C) Handler4C(conn net.Conn) {
	defer conn.Close()
	const MAXLEN = 2048
	buf := make([]byte, MAXLEN)

	defer func() {
		fmt.Println("socket is close")
		conn.Close()
	}()

	for {
		n, err := conn.Read(buf) //接收具体消息
		if err != nil {
			return
		}

		if n > MAXLEN || n < 8 {
			account.Log.Error("recive error n> MAXLEN")
			return
		}

		//接收包头
		var head_len int32 = 0
		var head_pid int32 = 0
		buffer_len := bytes.NewBuffer(buf[0:4])
		buffer_pid := bytes.NewBuffer(buf[4:8])
		binary.Read(buffer_len, binary.BigEndian, &head_len)
		binary.Read(buffer_pid, binary.BigEndian, &head_pid)

		//接收包体
		switch head_pid {

		case 1:
			var db_count_max int32 = 0
			redis.Find("PlayerCount", db_count_max)
			fmt.Println("获取playerCount:", db_count_max)
			//注册
			register := new(protocol.Account_RegisterPlayer)
			if err := proto.Unmarshal(buf[8:n], register); err == nil {
				game_id, _, _ := this.getNewAddress()
				fmt.Println("注册:", register.GetPlayername(), register.GetPassworld(), game_id)
				result, player_id := this.account_info.Register(register.GetPlayername(), register.GetPassworld(), game_id)

				result4C := &protocol.Account_RegisterResult{
					Result: proto.Int32(int32(result)),
				}

				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 1, encObj)

				//通知game注册成功
				if global.REGISTERSUCCESS == result {
					err := this.NoteGame(player_id, game_id)
					if err != nil {
						account.Log.Error("err1:", err)
					}
					account.Log.Info("player_id = %d game_id=%d", player_id, game_id)
				}
			}

		case 2:
			//登陆
			login := new(protocol.Account_LoginInfo)
			if err := proto.Unmarshal(buf[8:n], login); err == nil {
				result, player_id, server_address := this.account_info.VerifyLogin(login.GetPlayername(), login.GetPassworld(), config.AllServerAddress)
				//发送登陆并断开连接
				result4C := &protocol.Account_LoginResult{
					Result:     proto.Int32(int32(result)),
					PlayerId:   proto.Int32(int32(player_id)),
					Gameserver: proto.String(server_address),
				}

				encObj, _ := proto.Marshal(result4C)
				SendPackage(conn, 2, encObj)

				if result == global.LOGINSUCCESS { //登录成功断开连接
					conn.Close()
				}

			} else {
				fmt.Println(err)
			}

		case 3: //服务器列表
			var serverInfo_list []*protocol.Account_ServerInfo
			for k, v := range config.AllServerAddress {
				serverInfo := new(protocol.Account_ServerInfo)
				type_ := deal_4g.BusyLevel(k)
				serverInfo.Type = &type_
				serverInfo.ServerId = &k
				serverInfo.ServerAddress = &v
				serverInfo_list = append(serverInfo_list, serverInfo)
			}

			result4C := &protocol.Account_ServerListResult{
				ServerInfo: serverInfo_list,
			}
			encObj, _ := proto.Marshal(result4C)
			SendPackage(conn, 3, encObj)

		case 4: //返回自己登录过得服务器
			var serverInfo_list []*protocol.Account_ServerInfo
			server_id, servers := this.account_info.GetServers()

			for _, v := range servers {
				serverInfo := new(protocol.Account_ServerInfo)
				type_ := deal_4g.BusyLevel(v)
				serverInfo.Type = &type_
				serverInfo.ServerId = &v
				server_address, _ := config.AllServerAddress[v]
				serverInfo.ServerAddress = &server_address
				serverInfo_list = append(serverInfo_list, serverInfo)
			}

			result4C := &protocol.Account_MyServerListResult{
				LastServerId: &server_id,
				MyServerList: serverInfo_list,
			}
			encObj, _ := proto.Marshal(result4C)
			SendPackage(conn, 4, encObj)
		default:
		}
	}
}
