//装备相关
package game

//"strconv"

type HeroInfo struct { //hero基础属性
	Hero_id    int32 //英雄id
	Hero_uid   int32 //唯一uid
	Level      int32 //等级
	Exp        int32 //经验
	Hp         int32 //血量
	Power      int32 //战力
	Star_level int32 //星级
	Step_level int32 //阶级
}

type HeroAttibute struct { //英雄动态属性
	Key   int32
	Value int32
	Group int32
}

type HeroStruct struct {
	Hero_Info     HeroInfo
	Hero_Attibute HeroAttibute
}

//type HeroStructList struct {
//OnHeroStruct  []HeroStruct
//OffHeroStruct []HeroStruct
//}

//func (this *HeroStructList) addExp(exp int32) {
//var add_level int32 = 0
//var new_exp int32 = 0
//role_level_str := strconv.Itoa(int(this.Info.Level))

//}
