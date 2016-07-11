package account

import (
	"fmt"
	"server/share/global"
	"strings"
	"sync"

	"github.com/game_engine/cache/redis"
	//_ "github.com/go-sql-driver/mysql"
)

type LoginBase struct {
	PlayerId   int32
	PlayerName string
	PlayerPwd  string
	ServerId   int32
	IsForBid   bool
	Servers    []int32
}

type AccountInfo struct {
	accountMutex     *sync.RWMutex
	redis_login_base *LoginBase
}

func (this *AccountInfo) Init() {
	this.accountMutex = new(sync.RWMutex)
	this.redis_login_base = new(LoginBase)
}

func (this *AccountInfo) Register(name string, pwd string, server_id int32) (int32, int32) {
	this.accountMutex.Lock()
	defer this.accountMutex.Unlock()

	//通过内存数据库 先检查是否username相同
	err := redis.Find("PlayerName:"+name, this.redis_login_base)
	if err != nil { //没有查到数据

		db_count_max += 1
		user := LoginBase{PlayerId: db_count_max, PlayerName: name, PlayerPwd: pwd, ServerId: server_id, IsForBid: false}
		this.redis_login_base.Servers = append(this.redis_login_base.Servers, server_id)
		//内存数据库
		err_redis_player := redis.Modify("PlayerName:"+name, user)
		err_redis_count := redis.Modify("PlayerCount", db_count_max)

		if err_redis_player == nil && err_redis_count == nil {
			return global.REGISTERSUCCESS, db_count_max
		} else {
			fmt.Println("数据写入错误")
		}
	} else {
		Log.Trace("name = %s pwd = %s have same SAMENICK", name, pwd)
	}
	return global.SAMENICK, 0
}

func (this *AccountInfo) VerifyLogin(name string, pwd string, allServerAddress map[int32]string) (result int32, player_id int32, game_address string) {
	err := redis.Find("PlayerName:"+name, this.redis_login_base)
	if err == nil {
		if strings.EqualFold(this.redis_login_base.PlayerName, name) && strings.EqualFold(this.redis_login_base.PlayerPwd, pwd) && !this.redis_login_base.IsForBid {
			serverAddress, _ := allServerAddress[this.redis_login_base.ServerId]
			return global.LOGINSUCCESS, this.redis_login_base.PlayerId, serverAddress
		} else {
			return global.FORBIDLOGIN, 0, ""
		}
	}
	return global.LOGINERROR, 0, ""
}

func (this *AccountInfo) GetServers() (int32, []int32) {
	return this.redis_login_base.ServerId, this.redis_login_base.Servers
}
