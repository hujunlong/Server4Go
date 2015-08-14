package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"server/protocol"
	"time"
)

func CheckError(err error) bool {
	if err != nil {
		fmt.Println("err:", err)
		return true
	}
	return false
}

func main() {

	buf := make([]byte, 1024)
	/*
		loginInfo := &protocol.S2SSystem_LoginInfo{
			Name:      proto.String("胡俊龙"),
			Passworld: proto.String("1234567"),
		}
	*/

	registerInfo := &protocol.S2SSystem_RegisterInfo{
		Name:      proto.String("胡俊龙"),
		Passworld: proto.String("1234567"),
		Age:       proto.Int32(12),
		Sex:       proto.Int32(1),
	}

	//发送数据包
	encObj, err := proto.Marshal(registerInfo)
	CheckError(err)
	conn, _ := net.Dial("tcp", "127.0.0.1:8087")
	conn.Write(encObj)

	//接收数据包
	n, _ := conn.Read(buf) //接收具体消息
	result := new(protocol.S2SSystem_ResultInfo)
	if err = proto.Unmarshal(buf[0:n], result); err == nil {

		switch result.GetResult() {
		case 0:
			fmt.Println("login sccessfull!")
			break
		case 1:
			fmt.Println("register false!")
			break
		case 3:
			fmt.Println("login false!")
		}

	}

	conn.Close()
	time.Sleep(10 * time.Second)

}
