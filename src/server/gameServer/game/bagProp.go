package game

import (
	"net"
	"server/share/protocol"

	"github.com/golang/protobuf/proto"
)

type BagProp struct {
	Max       int32          //掉落类型
	OpenCount int32          //开启的个数
	UseCount  int32          //使用个数
	Props     map[int32]Prop //道具uid做key
}

func (this *BagProp) Init() {
	this.Props = make(map[int32]Prop)
}

func (this *BagProp) addItem(prop_id int32, count int32, count_max int32) (int32, int32, Prop) { //参数1（剩余数量） 参数2（0:ok 1:背包不足）
	var last_count int32 = 0
	var result int32 = 0
	var prop_ Prop

	if this.OpenCount > this.UseCount {
		if count > count_max {
			prop_.Count = count_max
			prop_.Prop_id = prop_id
			prop_.Prop_uid = GetUid()
			last_count = count - count_max
		} else {
			prop_.Count = count
			prop_.Prop_id = prop_id
			prop_.Prop_uid = GetUid()
		}

		this.UseCount += 1
		this.Props[prop_.Prop_uid] = prop_
	} else {
		result = 1
	}
	return last_count, result, prop_
}

func (this *BagProp) Add(prop_id int32, count int32) (int32, []Prop) { //（0:ok 1:背包不足 2:不能放入背包）
	var notice_props []Prop
	index := Csv.item.index_value["109"]
	count_max := Csv.item.simple_info_map[prop_id][index]
	if count_max == 0 {
		return 2, nil
	}

	//遍历添加将所有空位置补上
	for key, v := range this.Props {
		if prop_id == v.Prop_id {
			if v.Count+count <= count_max {
				v.Count += count
				count = 0
			} else {
				count = (count + this.Props[key].Count - count_max)
				v.Count = count_max
			}
			this.Props[key] = v
			notice_props = append(notice_props, this.Props[key]) //用来返回通知客户端
		}
	}

	//开启新背包
	for count > 0 {
		last_count, err_id, prop_ := this.addItem(prop_id, count, count_max)
		if err_id > 0 {
			return err_id, notice_props
		}
		count = last_count
		notice_props = append(notice_props, prop_)
	}
	return 0, notice_props
}

func (this *BagProp) Adds(props []Prop, conn *net.Conn) []Prop {

	var notice_props []Prop
	for _, v := range props {
		result, props := this.Add(v.Prop_id, v.Count)
		if result == 0 {
			notice_props = append(notice_props, props...)
		} else {

			var result int32 = 2
			result4C := &protocol.Game_Notice2CMsg{
				Msg: &result,
			}
			encObj, _ := proto.Marshal(result4C)
			SendPackage(*conn, 1028, encObj)

			break
		}
	}

	return notice_props
}

func (this *BagProp) GetTotal(id int32) int32 {
	var total int32
	for _, v := range this.Props {
		if v.Prop_id == id {
			total += v.Count
		}
	}
	return total
}

func (this *BagProp) deleteItem(id int32, count int32) []Prop {
	var notice_props []Prop

	for i, v := range this.Props {
		if v.Prop_id == id {
			if count > v.Count {
				count -= v.Count
				v.Count = 0
				notice_props = append(notice_props, v)
				delete(this.Props, v.Prop_uid)
				this.UseCount -= 1
			} else {
				count = 0
				v.Count = v.Count - count
				this.Props[i] = v
				notice_props = append(notice_props, v)
				break
			}
		}
	}
	return notice_props
}

func (this *BagProp) Use(uid int32, count int32) (int32, []Prop) { //（0:ok 1:不存该道具id 2:道具总量少于请求数量）

	var notice_props []Prop

	if _, ok := this.Props[uid]; !ok {
		return 1, nil
	}

	prop := this.Props[uid]

	if prop.Count > count {
		prop.Count -= count
		this.Props[uid] = prop
		notice_props = append(notice_props, prop)
		return 0, notice_props
	} else {
		total := this.GetTotal(prop.Prop_id)
		if total < count {
			return 2, nil
		} else {
			buff_props := this.deleteItem(prop.Prop_id, count)
			notice_props = append(notice_props, buff_props...)
		}
	}

	return 0, notice_props
}
