package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"server/protocol"
	"time"
)

var nick string
var passwd string
var ok bool = false
var isCanGetList bool = false //获取聊天列表
var message string

func CheckError(err error) bool {
	if err != nil {
		fmt.Println("err:", err)
		return true
	}
	return false
}

func ReciveResult(conn net.Conn, cn chan int) {
	const MAXLEN = 1024
	buf := make([]byte, MAXLEN)

	for true {
		n, _ := conn.Read(buf) //接收具体消息
		//接收数据包
		//接收包的type类型用来区分包之间的区别
		typeStruct := new(protocol.S2SSystem_GetType)
		if err := proto.Unmarshal(buf[0:n], typeStruct); err != nil {
			CheckError(err)
			return
		}

		switch *typeStruct.Type {
		case 5: //接收返回结果ResultInfo
			result := new(protocol.S2SSystem_ResultInfo)
			if err := proto.Unmarshal(buf[0:n], result); err == nil {
				switch result.GetResult() {
				case 5:
					fmt.Println("login sccessfull!")
					ok = true

					request := &protocol.S2SSystem_Request{
						Type:   proto.Int32(3),
						Result: proto.Int32(1),
					}
					encObj, _ := proto.Marshal(request)
					conn.Write(encObj)
					break
				}
			}
			break
		case 4: //ResultChatMsg
			result := new(protocol.S2SSystem_ResultChatMsg)
			if err := proto.Unmarshal(buf[0:n], result); err == nil {
				for i := 0; i < len(result.Type); i++ {
					fmt.Println(result.Playername[i], "say:", result.Msg[i])
				}
			}
			break
		}
	}

}

func main() {

	conn, _ := net.Dial("tcp", "127.0.0.1:8087")
	
	ch := make(chan int)
	go ReciveResult(conn, ch)

 
	for true {
		if !ok {
			fmt.Println("please enter nick:")
			fmt.Scanln(&nick)
			fmt.Println("please enter passwd:")
			fmt.Scanln(&passwd)

			//登陆相关
			loginInfo := &protocol.S2SSystem_LoginInfo{
				Type:      proto.Int32(1),
				Name:      proto.String(nick),
				Passworld: proto.String(passwd),
			}

			//发送数据包
			encObj, err := proto.Marshal(loginInfo)
			CheckError(err)
			conn.Write(encObj)
		}

		if len(message) > 0 {
			request := protocol.S2SSystem_ResultChatMsg{Type: make([]int32, 1), Playername: make([]string, 1), Msg: make([]string, 1)}
			var mytype int32 = 4
			request.Type[0] = mytype
			request.Playername[0] = nick
			request.Msg[0] = message
			encObj, _ := proto.Marshal(&request)
			conn.Write(encObj)
			message = ""
		} else {
			fmt.Scanln(&message)
		}

		time.Sleep(100 * time.Millisecond)
	}

	<-ch
 
	/*
	loginInfo := &protocol.S2SSystem_RegisterInfo{
				Type:      proto.Int32(2),
				Name:      proto.String("bbbb"),
				Age :	   proto.Int32(10),
				Passworld: proto.String("12345678"),
				Sex: proto.Int32(2),
			}

			//发送数据包
			encObj, err := proto.Marshal(loginInfo)
			CheckError(err)
			fmt.Println(encObj)
			conn.Write(encObj)
	*/		
	defer conn.Close()

}
