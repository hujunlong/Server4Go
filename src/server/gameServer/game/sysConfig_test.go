package game

import (
	"strings"
	"testing"
)

func TestSysConfig(t *testing.T) {
	sys_cfg := new(SysConfig)
	sys_cfg.Init()

	sys_cfg.readConfig()
	if sys_cfg.GameId != 1 {
		t.Error("game id error")
	}

	if strings.EqualFold(sys_cfg.Server2AccountAddress, "127.0.0.1:8081") && strings.EqualFold(sys_cfg.ServerAddress, "0.0.0.0:8082") && sys_cfg.DistanceTime == 20 {

	} else {
		t.Error("sys_cfg server config error")
	}

}
