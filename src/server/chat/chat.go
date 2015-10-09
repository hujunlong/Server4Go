package chat

import ()

type ChatMsg struct {
	Playername string
	Msg        string
}

const (
	MAXLEN = 30 //最大消息条数
)

var ChatList []ChatMsg

func init() {
	ChatList = make([]ChatMsg, 0)
}

func AddMsg(name string, msg string) int {
	if len(ChatList) >= MAXLEN {
		ChatList = append(ChatList[1:], ChatMsg{name, msg})
	} else {
		ChatList = append(ChatList, ChatMsg{name, msg})
	}
	return len(ChatList)
}
