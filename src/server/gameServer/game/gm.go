//GM相关
package game

import (
	"fmt"
	"server/share/protocol"
	"strings"

	"github.com/golang/protobuf/proto"
)

type GM struct {
}

func (this *GM) DealMsg(str string, player *Player) {
	str = strings.TrimSpace(str)
	strs := strings.Split(str, " ")
	var data []int32

	fmt.Println("strs", strs)

	for i := 1; i <= len(strs)-1; i++ {
		num := Str2Int32(strs[i])
		data = append(data, num)
	}

	if len(data) == 1 {
		switch strs[0] {
		case "gold":
			player.ModifyGold(data[0])
		case "diamond":
			player.ModifyDiamond(data[0])
		case "energy":
			fmt.Println("energy", data[0])
			player.ModifyEnergy(data[0])
		case "exp":
			fmt.Println("exp", data[0])
			player.AddRoleExp(data[0])

		case "hero":
			fmt.Println("hero", data[0])
			hero := new(HeroStruct)
			hero_uid, _ := hero.CreateHero(data[0], player)
			player.Heros[hero_uid] = hero

			//推送
			heros := make(map[int32]*HeroStruct)
			heros[hero_uid] = hero
			heros_FightingAttr := player.GetHeroStruct(heros)
			result4C := &protocol.NoticeMsg_NoticeGetHeros{
				Heros: heros_FightingAttr,
			}
			encObj, _ := proto.Marshal(result4C)
			SendPackage(*player.conn, 1212, encObj)
		}
	}

	if len(data) == 2 {
		switch strs[0] {
		case "prop":
			var prop Prop
			prop.Prop_id = data[0]
			prop.Count = data[1]
			prop.Prop_uid = GetUid()
			player.Bag_Prop.AddAndNotice(prop, player.conn)
		case "equip":
			var equips []Equip
			var equip Equip
			for i := 0; i < int(data[1]); i++ {
				equp_ := equip.Create(data[0], 1, player) //品质组暂时定为1
				equips = append(equips, *equp_)
			}
			player.Bag_Equip.Adds(equips, player.conn)
		}
	}

	//result 返回
	var result int32 = 1
	if strings.EqualFold(strs[0], "gold") || strings.EqualFold(strs[0], "diamond") || strings.EqualFold(strs[0], "energy") || strings.EqualFold(strs[0], "exp") || strings.EqualFold(strs[0], "hero") {
		if len(data) == 1 {
			result = 0
		}
	}

	if strings.EqualFold(strs[0], "prop") || strings.EqualFold(strs[0], "equip") {
		if len(data) == 2 {
			result = 0
		}
	}

	result4C := &protocol.GM_MsgResult{
		Result: &result,
	}
	encObj, _ := proto.Marshal(result4C)
	SendPackage(*player.conn, 1601, encObj)
}
