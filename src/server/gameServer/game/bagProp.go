package game

import (
	"net"
	"server/share/protocol"

	"fmt"

	"github.com/golang/protobuf/proto"
)

type BagProp struct {
	Max       int32           //掉落类型
	OpenCount int32           //开启的个数
	UseCount  int32           //使用个数
	Props     map[int32]*Prop //道具uid做key
	PropsById map[int32]int32 //道具id 做key
	player    *Player         //道具装备
}

func (this *BagProp) Init(player *Player) {
	this.Props = make(map[int32]*Prop)
	this.PropsById = make(map[int32]int32)
	this.player = player
}

//新位置添加道具
func (this *BagProp) AddNewItem(prop_id int32, count int32, count_max int32) (int32, int32, Prop) { //参数1（剩余数量） 参数2（0:ok 1:背包不足）
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
		this.Props[prop_.Prop_uid] = &prop_
	} else {
		result = 1
	}

	this.PropsById[prop_id] += prop_.Count
	return last_count, result, prop_
}

//添加一个道具 不推送消息
func (this *BagProp) Add(prop_id int32, count int32) (int32, []Prop) { //（0:ok 1:背包不足 2:不能放入背包）
	buff_count := count //该数据用任务系统计算
	var notice_props []Prop
	count_max := Csv.item[prop_id].Id_109
	if count_max == 0 {
		return 2, nil
	}

	//遍历添加将所有空位置补上
	for key, v := range this.Props {
		if prop_id == v.Prop_id {
			if v.Count+count <= count_max {
				v.Count += count
				this.PropsById[prop_id] += count
				count = 0
			} else {
				count = (count + this.Props[key].Count - count_max)
				v.Count = count_max
				this.PropsById[prop_id] += (count_max - this.Props[key].Count)
			}
			this.Props[key] = v
			notice_props = append(notice_props, *this.Props[key]) //用来返回通知客户端
		}
	}

	//开启新背包
	for count > 0 {
		last_count, err_id, prop_ := this.AddNewItem(prop_id, count, count_max)
		if err_id > 0 {
			return err_id, notice_props
		}
		count = last_count
		notice_props = append(notice_props, prop_)
	}

	this.player.Task.TriggerEvent(2, buff_count-count, prop_id) //任务
	return 0, notice_props
}

//添加并通知
func (this *BagProp) AddAndNotice(prop Prop, conn *net.Conn) {
	result, props := this.Add(prop.Prop_id, prop.Count)
	if result != 0 {

		this.BagWeek(conn)
	} else {
		this.Notice2CProp(props, conn)
	}

}

//添加多个道具并推送
func (this *BagProp) Adds(props []Prop, conn *net.Conn) {
	var notice_props []Prop
	for _, v := range props {
		result, props := this.Add(v.Prop_id, v.Count)
		if result == 0 {
			notice_props = append(notice_props, props...)
		} else {
			this.BagWeek(conn)
			break
		}
	}

	this.Notice2CProp(notice_props, conn)
}

//删除该物品by uid
func (this *BagProp) DeleteItemByUid(uid int32, count int32) ([]Prop, int32) { //（0:ok 1:不存该道具id 2:道具总量少于请求数量）
	var notice_props []Prop

	if _, ok := this.Props[uid]; !ok {
		return notice_props, 1
	}

	if this.Props[uid].Count < count {
		return notice_props, 2
	} else {
		buff_prop := Prop{this.Props[uid].Prop_id, uid, this.Props[uid].Count - count}
		this.Props[uid] = &buff_prop
		notice_props = append(notice_props, buff_prop)
		this.PropsById[this.Props[uid].Prop_id] -= count

		//任务
		this.player.Task.TriggerEvent(28, count, this.Props[uid].Prop_id)
		return notice_props, 0
	}
}

//删除该物品by id
func (this *BagProp) DeleteItemById(id int32, count int32, conn *net.Conn) { //（0:ok  2:道具总量少于请求数量）
	var notice_props []Prop

	if this.PropsById[id] < count {
		return
	}

	for i, v := range this.Props {
		if v.Prop_id == id {
			if count > v.Count {
				count -= v.Count
				v.Count = 0
				notice_props = append(notice_props, Prop{id, v.Prop_uid, 0})
				delete(this.Props, v.Prop_uid)
				this.UseCount -= 1
			} else {
				v.Count = v.Count - count
				this.Props[i] = v
				notice_props = append(notice_props, Prop{id, v.Prop_uid, v.Count})
				break
			}
		}
	}

	this.PropsById[id] -= count
	this.Notice2CProp(notice_props, conn)

	//任务系统

	this.player.Task.TriggerEvent(28, count, id)
}

//转换protocol格式
func (this *BagProp) DealPropStruct(Props []Prop) []*protocol.PropStruct {
	//道具
	var props []*protocol.PropStruct
	for i, _ := range Props {
		prop := new(protocol.PropStruct)
		prop.PropId = &Props[i].Prop_id
		prop.PropUid = &Props[i].Prop_uid
		prop.PropCount = &Props[i].Count
		props = append(props, prop)
	}
	return props
}

//推送消息
func (this *BagProp) Notice2CProp(props []Prop, conn *net.Conn) {
	props_struct := this.DealPropStruct(props)

	result4C := &protocol.NoticeMsg_Notice2CProp{
		Prop: props_struct,
	}

	encObj, _ := proto.Marshal(result4C)
	SendPackage(*conn, 1201, encObj)
}

//推送背包不足
func (this *BagProp) BagWeek(conn *net.Conn) {
	var result int32 = 2
	result4C := &protocol.NoticeMsg_Notice2CMsg{
		Msg: &result,
	}
	encObj, _ := proto.Marshal(result4C)
	SendPackage(*conn, 1208, encObj)
}

//背包道具是否充足
func (this *BagProp) PropIsenough(datas []Data) bool {
	for _, v := range datas {
		if _, ok := this.PropsById[v.Id]; !ok {
			return false
		}

		fmt.Println("背包拥有物品数量", this.PropsById[v.Id])
		if this.PropsById[v.Id] < v.Num {
			return false
		}
	}
	return true
}

//扣除道具
func (this *BagProp) SaleProp(uid int32, count int32) {
	props, result := this.DeleteItemByUid(uid, count)
	this.Notice2CProp(props, this.player.conn)

	result4C := &protocol.Goods_SalePropResult{
		Result: &result,
	}
	encObj, _ := proto.Marshal(result4C)
	SendPackage(*this.player.conn, 1506, encObj)
}
