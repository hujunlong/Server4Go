//挂机挑战玩家
package game

import (
	"time"
)

//玩家相关显示属性
type GuajiPK struct {
	PK_type       int32 //(1:能够pk 2:免战牌不能pk 3:受保护不能pk 4:等级不够未开放)
	Last_pk_num   int32 //剩余pk次数
	Last_pked_num int32 //剩余被pk次数
	Kill_num      int32 //杀死玩家数量
	Protect_time  int32 //保护时间产生点 记录当时时间
	Pk_less_time  int32 //减少收益时间点 记录当时时间
	Last_add_time int32 //上次加次数时间 记录当时时间
	Use_buy_num   int32 //使用购买次数

	pk_open_level       int32 //开放等级2062
	pk_num              int32 //pk次数2063
	pked_num            int32 //被pk次数2064
	pk_consumer_yuanbao int32 //pk消耗的元宝2065
	pk_protect_time     int32 //保护时间2066
	pk_cd               int32 //pk次数恢复时间2067
	pk_less_per         int32 //收益减少比例2068
	pk_less_time        int32 //收益减少时间 2069
	pk_other_less_per   int32 //pk成功对手减少比例 2076
	pk_other_less_time  int32 //pk成功对手减少比例时间2070
	pk_fanji_yuanbao    int32 //反击消耗元宝 2071
	pk_buy_num          int32 //购买次数 2075
}

func (this *GuajiPK) Init() {

	this.pk_open_level = int32(Csv.property[2062].Id_102)

	this.pk_num = int32(Csv.property[2063].Id_102)

	this.pked_num = int32(Csv.property[2064].Id_102)

	this.pk_consumer_yuanbao = int32(Csv.property[2065].Id_102)

	this.pk_protect_time = int32(Csv.property[2066].Id_102)

	this.pk_cd = int32(Csv.property[2067].Id_102)

	this.pk_less_per = int32(Csv.property[2068].Id_102)

	this.pk_less_time = int32(Csv.property[2069].Id_102)

	this.pk_other_less_per = int32(Csv.property[2076].Id_102)

	this.pk_other_less_time = int32(Csv.property[2070].Id_102)

	this.pk_fanji_yuanbao = int32(Csv.property[2071].Id_102)

	this.pk_buy_num = int32(Csv.property[2075].Id_102)

	Log.Info("pk_open_level = %d pk_num=%d pked_num =%d pk_consumer_yuanbao = %d pk_protect_time = %d pk_cd = %d pk_less_per=%d pk_less_time = %d pk_other_less_per = %d pk_other_less_time = %d pk_fanji_yuanbao=%d pk_buy_num = %d", this.pk_open_level, this.pk_num, this.pked_num, this.pk_consumer_yuanbao, this.pk_protect_time, this.pk_cd, this.pk_less_per, this.pk_less_time, this.pk_other_less_per, this.pk_other_less_time, this.pk_fanji_yuanbao, this.pk_buy_num)
}

//获取恢复次数与Last_add_time
func (this *GuajiPK) AddPkNum() {
	Log.Info("this.Last_add_time = %d", this.Last_add_time)
	now := int32(time.Now().Unix())
	add_pk_num := (now - this.Last_add_time) / this.pk_cd

	if add_pk_num > 0 {
		//检查是否超过上限
		if (add_pk_num + this.Last_pk_num) >= this.pk_num {
			this.Last_pk_num = this.pk_num
			this.Last_add_time = now
		} else {
			this.Last_pk_num += add_pk_num
			this.Last_add_time = now
		}
	}
}

//计算受保护剩余时间
func (this *GuajiPK) GetLastProtectTime() int32 {
	Log.Info("this.Protect_time = %d", this.Protect_time)
	if this.Protect_time == 0 {
		return 0
	}
	now := int32(time.Now().Unix())
	if now-this.Protect_time > this.pk_protect_time {
		this.Protect_time = 0
		return 0
	} else {
		return (this.pk_protect_time + this.Protect_time - now)
	}
}

//计算减少收益时间
func (this *GuajiPK) GetPKLessTime() int32 {
	Log.Info("%d", this.Pk_less_time)
	if this.Pk_less_time == 0 {
		return 0
	}

	now := int32(time.Now().Unix())
	if now-this.Pk_less_time > this.pk_less_time {
		this.Pk_less_time = 0
		return 0
	} else {
		return (this.pk_less_time + this.Pk_less_time - now)
	}
}

//获取最新状况
func (this *GuajiPK) GetGuajiNowInfo() (int32, int32) {
	this.Init()
	this.AddPkNum()
	protected_last_time := this.GetLastProtectTime()
	less_last_time := this.GetPKLessTime()
	return protected_last_time, less_last_time
}