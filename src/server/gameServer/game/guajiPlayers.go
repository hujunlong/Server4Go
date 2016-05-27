//全局关卡玩家相关数据 不需要保存数据
package game

import (
	"sync"
)

type PlayerGuajiPk struct {
	role_id  int64
	guaji_pk *GuajiPK
}

//各个关卡挂机玩家
type GuajiPlayers struct {
	player_list map[int32][]int64 //key:地图id
	guaji_Mutex *sync.RWMutex
}

func (this *GuajiPlayers) Init() {
	this.player_list = make(map[int32][]int64)
	this.guaji_Mutex = new(sync.RWMutex)
}

func (this *GuajiPlayers) Enter(stage_id int32, role_id int64) {
	this.guaji_Mutex.Lock()
	defer this.guaji_Mutex.Unlock()

	this.player_list[stage_id] = append(this.player_list[stage_id], role_id)
	Log.Info("stage_id = %d PlayerGuajiPk = %d", stage_id, role_id)
}

func (this *GuajiPlayers) Exit(stage_id int32, role_id int64) {
	this.guaji_Mutex.Lock()
	defer this.guaji_Mutex.Unlock()

	roles := this.player_list[stage_id]
	for i, v := range roles {
		if v == role_id {
			ss := roles[:i]
			ss = append(ss, roles[i+1:]...)
			this.player_list[stage_id] = ss
		}
	}
	Log.Info("role_id = %d exit stage id %d", role_id, stage_id)
}
