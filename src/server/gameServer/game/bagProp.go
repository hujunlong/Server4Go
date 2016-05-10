package game

type BagProp struct {
	Max       int32  //掉落类型
	OpenIndex int32  //开启的个数
	UseCount  int32  //使用个数
	Props     []Prop //装备
}
