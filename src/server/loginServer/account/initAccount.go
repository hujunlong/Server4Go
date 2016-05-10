package account

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/game_engine/cache/redis"
	"github.com/game_engine/i18n"
	"github.com/game_engine/logs"
)

//全局性
var Log *logs.BeeLogger
var db_count_max int32 //数据库最大注册人数
type Config struct {
	count              int32 //现在最大的playerid
	account_log_max    int64
	Listen4CAddress    string
	Listen4GameAddress string
	NewServerAddress   map[int32]string
	AllServerAddress   map[int32]string
}

func (this *Config) Init() {
	this.NewServerAddress = make(map[int32]string)
	this.AllServerAddress = make(map[int32]string)
	this.setLog()
	this.readConfig()
	this.openNewServerConfig()
	this.openAllServerConfig()
}

func (this *Config) setLog() {
	Log = logs.NewLogger(this.account_log_max) //日志
	Log.EnableFuncCallDepth(true)
	err := Log.SetLogger("file", `{"filename":"log/account.log"}`)
	if err != nil {
		fmt.Println(err)
	}
}

func (this *Config) readConfig() {
	err := il8n.GetInit("config/account_cfg.ini")
	if err == nil {
		this.account_log_max, _ = strconv.ParseInt(il8n.Data["account_log_max"].(string), 10, 64)
		this.Listen4CAddress = il8n.Data["login_listen_4c_ip"].(string)
		this.Listen4GameAddress = il8n.Data["login_listen_4game_ip"].(string)
	} else {
		Log.Error(err.Error())
	}
}

func (this *Config) openNewServerConfig() {
	for k, v := range il8n.Data {
		if strings.Contains(k.(string), "new_") {
			key_str := strings.TrimLeft(k.(string), "new_")
			key_int, _ := strconv.Atoi(key_str)
			this.NewServerAddress[int32(key_int)] = v.(string)
		}
	}

	if len(this.NewServerAddress) == 0 {
		Log.Error("new player can't connect,config can't find new server id")
	}
}

func (this *Config) openAllServerConfig() {
	for k, v := range il8n.Data {
		if strings.Contains(k.(string), "server_") {
			key_str := strings.TrimLeft(k.(string), "server_")
			key_int, _ := strconv.Atoi(key_str)
			this.AllServerAddress[int32(key_int)] = v.(string)
		}
	}

	if len(this.AllServerAddress) == 0 {
		Log.Error("new player can't connect,config can't find new server id")
	}
}

func getMaxId() int32 {
	var count int32 = 0
	err := redis.Find("PlayerCount", count)
	if err != nil {
		return 0
	}
	return count
}

func init() {
	db_count_max = getMaxId()
}
