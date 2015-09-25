package global

import (
	"github.com/game_engine/data/redis"
	"github.com/game_engine/logs"
)

//用于全局
var Redis *redis.Client
var Log *logs.BeeLogger

func Init() {
	Redis = new(redis.Client)   //内存数据
	Log = logs.NewLogger(10000) //日志
	Log.EnableFuncCallDepth(true)
	Log.SetLogger("file", `{"filename":"game.log"}`)
}

const (
	REGISTERSUCCESS = 0 //注册成功
	REGISTERERROR   = 1 //注册错误
	LOGINERROR      = 2 //登陆错误
	PASSWDERROR     = 3 //密码错误
	SAMENICK        = 4 //注册名字相同
	LOGINSUCCESS    = 5 //登陆成功
)
