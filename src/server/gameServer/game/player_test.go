package game

import (
	"fmt"
	"server/share/global"
	"strings"
	"testing"
)

func TestSysPlayer(t *testing.T) {

	//验证创建
	player := new(Player)
	result := player.RegisterRole(11118, 11112, 1, "gbkd", 901)
	if result == global.REGISTERROLESUCCESS {
		fmt.Println("create ok")
	}

	if result == global.REGISTERROLEERROR {
		fmt.Println("existing the data")
	}

	//验证获取
	get_player := LoadPlayer("11118")
	fmt.Println("hp:", get_player.Info.Hp, "df:", get_player.Info.Physical_def)
	if get_player.Info.ID == 11118 && strings.EqualFold(get_player.Info.Nick, "gbkd") && get_player.PlayerId == 11112 {
		fmt.Println("ok")
	} else {
		t.Error("error not found the data")
	}

}
