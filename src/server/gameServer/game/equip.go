//装备相关
package game

type Equip struct {
	equip_id         int32 //装备id
	equip_uid        int32 //装备唯一id
	pos              int32 //装备位置（-1 表示未装备）
	quality          int32 //装备品质
	equip_level      int32 //装备等级
	strengthen_count int32 //强化次数
	refine_count     int32 //精炼次数
}
