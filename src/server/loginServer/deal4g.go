package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"server/loginServer/account"
	"server/share/protocol"

	"github.com/golang/protobuf/proto"
)

type ConnectInfo struct {
	ServerId int32  //服务器id
	Address  string //地址
	Count    int32  //人数
	Conn     net.Conn
}

type Deal4G struct {
	GameConnects map[int32]ConnectInfo
}

func (this *Deal4G) Init() {
	this.GameConnects = make(map[int32]ConnectInfo)
}

func (this *Deal4G) Handler4Game(conn net.Conn) {
	//game与账号服务器断开
	defer func() {
		var key int32 = 0
		for i, v := range this.GameConnects {
			if v.Conn == conn {
				key = i
				break
			}
		}
		delete(this.GameConnects, key)
		conn.Close()
	}()

	const MAXLEN = 2048
	buf := make([]byte, MAXLEN)
	for {
		n, err := conn.Read(buf) //接收具体消息
		if err != nil {
			return
		}

		if n > MAXLEN {
			account.Log.Error("recive error n> MAXLEN")
			return
		}

		var head_len int32 = 0
		var head_pid int32 = 0
		buffer_len := bytes.NewBuffer(buf[0:4])
		buffer_pid := bytes.NewBuffer(buf[4:8])
		binary.Read(buffer_len, binary.BigEndian, &head_len)
		binary.Read(buffer_pid, binary.BigEndian, &head_pid)

		switch head_pid {
		case 101:
			get_account := new(protocol.Account_GameResult)
			if err := proto.Unmarshal(buf[8:n], get_account); err == nil {
				key := get_account.GetGameId()
				this.GameConnects[key] = ConnectInfo{key, get_account.GetGameAddress(), get_account.GetCount(), conn}
				fmt.Println("get_account.GetGameAddress():", get_account.GetGameAddress(), "num:", get_account.GetCount())
			} else {
				fmt.Println(err)
			}
		default:
		}
	}
}

func (this *Deal4G) Deal4GameServer(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		fmt.Println(conn.RemoteAddr().String(), " connet")
		if CheckError(err) {
			this.Handler4Game(conn)
		}
	}
}

func (this *Deal4G) BusyLevel(ServerId int32) int32 { //获取服务器繁忙程度

	var type_ int32 = 3
	_, ok := this.GameConnects[ServerId]
	if !ok {
		return -1
	}
	if this.GameConnects[ServerId].Count < 50 {
		type_ = 1
	} else if this.GameConnects[ServerId].Count >= 50 && this.GameConnects[ServerId].Count <= 200 {
		type_ = 2
	}
	return type_
}
