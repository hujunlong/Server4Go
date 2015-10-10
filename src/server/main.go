package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"net"
)

import (
	"server/account"
	"server/chat"
	"server/global"
	"server/protocol"
	"server/world"
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

		//接收包的type类型用来区分包之间的区别
		typeStruct := new(protocol.S2SSystem_GetType)
		if err1 := proto.Unmarshal(buf[0:n], typeStruct); err1 != nil {
			CheckError(err1)
			continue
		}
		switch *typeStruct.Type {
		case 1:
			//登陆
			login := new(protocol.S2SSystem_LoginInfo)
			if err := proto.Unmarshal(buf[0:n], login); err == nil {
				login_result, player := account.Login(login)
				if player != nil {
					//写入到内存
					player.Conn = &conn
					world.World.Players[player.Info.Name] = player
				}

				result := &protocol.S2SSystem_ResultInfo{
					Type:   proto.Int32(protocol.Default_S2SSystem_ResultInfo_Type),
					Result: proto.Int32(login_result),
				}

				global.Log.Trace("%s login rsult %d", login.GetName(), login_result)
				encObj, _ := proto.Marshal(result)
				conn.Write(encObj)
			}
			break
		case 2:
			//注册
			register := new(protocol.S2SSystem_RegisterInfo)
			if err := proto.Unmarshal(buf[0:n], register); err == nil {
				err := account.Register(register)
				result := new(protocol.S2SSystem_ResultInfo)

				if err == nil {
					buff := int32(global.REGISTERSUCCESS)
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
			break
		case 3:
			//聊天历史消息
			request := new(protocol.S2SSystem_Request)
			if err := proto.Unmarshal(buf[0:n], request); err == nil {
				switch *request.Result {
				case 1: //请求获取聊天记录
					len_ := len(chat.ChatList)
					result := protocol.S2SSystem_ResultChatMsg{Type: make([]int32, len_), Playername: make([]string, len_), Msg: make([]string, len_)}
					for i, v := range chat.ChatList {
						result.Type[i] = 4
						result.Playername[i] = v.Playername
						result.Msg[i] = v.Msg
					}
					encObj, _ := proto.Marshal(&result)
					conn.Write(encObj)
					global.Log.Info("get history chat")
					break
				}
			}
			break
		case 4:
			//玩家发言聊天
			Chat := new(protocol.S2SSystem_ResultChatMsg)
			if err := proto.Unmarshal(buf[0:n], Chat); err == nil {
				chat.AddMsg(Chat.GetPlayername()[0], Chat.GetMsg()[0])
				//发送给在线所以玩家
				encObj, _ := proto.Marshal(Chat)
				for _, v := range world.World.Players {
					(*v.Conn).Write(encObj) //直接将接收到的数据发出去
				}
			}
			break
		case 5:
			break
		}
	}

}

func main() {
	/*
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
	*/

}
