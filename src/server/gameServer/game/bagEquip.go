package game

type BagEquip struct {
	Max       int32   //掉落类型
	OpenIndex int32   //开启的个数
	UseCount  int32   //使用个数
	BagEquip  []Equip //装备
}
