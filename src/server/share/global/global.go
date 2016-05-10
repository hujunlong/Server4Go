package global

//用于全局

func init() {

}

//账号服务器消息
const (
	REGISTERSUCCESS = 0 //注册成功
	REGISTERERROR   = 1 //注册错误
	LOGINERROR      = 2 //登陆错误
	PASSWDERROR     = 3 //密码错误
	SAMENICK        = 4 //注册名字相同
	LOGINSUCCESS    = 5 //登陆成功
	FORBIDLOGIN     = 6 //禁止登陆
)

//游戏服务器
const (
	REGISTERROLESUCCESS = 100 //注册角色成功
	REGISTERROLEERROR   = 101 //注册角色失败
	CSVHEROIDEERROR     = 102 //CSV读取错误
)

const (
	Type_prop         = 1 //道具
	Type_hero         = 2 //英雄
	Type_equip        = 3 //装备
	Type_resource     = 4 //资源
	Type_gem          = 5 //英雄
	Type_jewelry      = 6 //首饰
	Type_magic_weapon = 7 //法宝
)
