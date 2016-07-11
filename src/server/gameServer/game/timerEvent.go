//游戏世界
package game

import (
	"fmt"

	"github.com/game_engine/timer"
)

type TimerManager struct {
}

func (this *TimerManager) Init() {

	//体力恢复刷新
	//distance_time := int32(Csv.property[2057].Id_102)
	timer.CreateTimer(100, true, this.TimerDealEnergy)

	//悬赏任务刷新
	//distance_time2 := int32(Csv.property[2013].Id_102)
	//timer.CreateTimer(30, true, this.TimerDealXuanShang)
}

//在线推送体力恢复
func (this *TimerManager) TimerDealEnergy() {
	for _, v := range word.players {
		if v.Info.Energy < v.Info.EnergyMax {
			v.ModifyEnergy(100)
		}
	}
}

//在线悬赏
func (this *TimerManager) TimerDealXuanShang() {
	fmt.Println("在线悬赏 定时刷新")
	for _, v := range word.players {
		v.Task.TaskXuanshangTimer(true)
	}
}
