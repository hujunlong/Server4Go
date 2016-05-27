//装备相关
package game

type Equip struct {
	Equip_id         int32 //装备id
	Equip_uid        int32 //装备唯一id
	Pos              int32 //装备位置（-1 表示未装备）
	Quality          int32 //装备品质
	Equip_level      int32 //装备等级
	Strengthen_count int32 //强化次数
	Refine_count     int32 //精炼次数
}
