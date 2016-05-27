//阵型
package game

//阵型规则
//7 4 1
//8 5 2
//9 6 3
type FomationStruct struct {
	Pos_id   int32
	Hero_id  int32
	Hero_uid int32
}

type Fomation struct {
	Hero_fomations []FomationStruct
}

func (this *Fomation) OnFomation(pos_id int32, hero_id int32, hero_uid int32) { //上阵
	if pos_id < 0 || pos_id > 9 {
		return
	}

	//先下掉老位置
	this.OffFomation(pos_id)

	var fomation FomationStruct
	fomation.Pos_id = pos_id
	fomation.Hero_id = hero_id
	fomation.Hero_uid = hero_uid
	this.Hero_fomations = append(this.Hero_fomations, fomation)

}

func (this *Fomation) OffFomation(pos_id int32) { //下阵
	if pos_id < 0 || pos_id > 9 {
		return
	}

	for i, v := range this.Hero_fomations {
		if v.Pos_id == pos_id {
			buff := this.Hero_fomations[:i]
			buff = append(buff, this.Hero_fomations[i+1:]...)
			this.Hero_fomations = buff
		}
	}
}

func (this *Fomation) ChangePos(pos_id_1 int32, pos_id_2 int32) { //交回位置
	if pos_id_1 < 0 || pos_id_1 > 9 || pos_id_2 < 0 || pos_id_2 > 9 {
		return
	}

	var pos_1 int = 0
	var pos_2 int = 0
	for i, v := range this.Hero_fomations {
		if v.Pos_id == pos_id_1 {
			pos_1 = i
		}
		if v.Pos_id == pos_id_2 {
			pos_2 = i
		}
	}

	if pos_1 == 0 || pos_2 == 0 {
		return
	}

	//交回数据
	fomation1 := this.Hero_fomations[pos_1]
	this.Hero_fomations[pos_1] = this.Hero_fomations[pos_2]
	this.Hero_fomations[pos_2] = fomation1
}
