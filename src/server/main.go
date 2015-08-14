package main

import (
	"fmt"
	"game_engine/protobuf/github.com/golang/protobuf/proto"
	"io"
	"net"
)

import (
	"server/account"
	"server/global"
	"server/protocol"
)

func CheckError(err error) bool {
	if err != nil {
		fmt.Println("err:", err)
		return true
	}
	return false
}

func Handler(conn net.Conn) {
	const MAXLEN = 1024
	buf := make([]byte, MAXLEN)
	for {
		n, err := conn.Read(buf) //接收具体消息
		if err == io.EOF {
			fmt.Println(conn.RemoteAddr(), " closed")
			conn.Close()
			return
		} else if CheckError(err) {
			return
		}

		if n > MAXLEN {
			global.Log.Error("recive error n> MAXLEN")
		}

		//注册
		register := new(protocol.S2SSystem_RegisterInfo)
		if err := proto.Unmarshal(buf[0:n], register); err == nil {
			err := account.Register(register)
			result := new(protocol.S2SSystem_ResultInfo)

			if err == nil {
				buff := int32(global.SUCCESS)
				result.Result = &buff
				global.Log.Trace("%s register success", register.GetName())

			} else {
				buff := int32(global.REGISTERERROR)
				result.Result = &buff
				global.Log.Trace("%s register faield", register.GetName())
			}
			encObj, _ := proto.Marshal(result)
			conn.Write(encObj)
		}

		//登陆
		login := new(protocol.S2SSystem_LoginInfo)
		if err := proto.Unmarshal(buf[0:n], login); err == nil {
			fmt.Println(login.GetName(), login.GetPassworld())
			login_result := account.Login(login)
			result := new(protocol.S2SSystem_ResultInfo)

			buff := int32(login_result)
			result.Result = &buff
			global.Log.Trace("%s login rsult %d", login.GetName(), login_result)

			encObj, _ := proto.Marshal(result)
			conn.Write(encObj)
		}
	}

}

func main() {

	//初始化数据(redis log)
	global.Init()

	listener, err := net.Listen("tcp", "0.0.0.0:8087")
	if !CheckError(err) {
		for {
			conn, err1 := listener.Accept()
			if !CheckError(err1) {
				fmt.Println(conn.RemoteAddr(), "connected")
				go Handler(conn)
			}
		}
	}
}
