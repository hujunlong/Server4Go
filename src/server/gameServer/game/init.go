package game

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/game_engine/cache/redis"
)

var Csv *CsvConfig          //CSV 配置
var Sys_config *SysConfig   //系统配置
var Json_config *JsonConfig //json 配置
var rand_ *rand.Rand        //随机数
var Global_Uid int32 = 0    //uid

var global_guaji_players *GuajiPlayers //全局的玩家进入关卡数据
var word *World                        //全局玩家
var timer_Manager *TimerManager        //全局定时刷新

func createbeginUid() int32 {
	t1 := int(time.Now().Unix())
	fmt.Println(t1)
	str1 := strconv.Itoa(t1)
	str1 = str1[5:]

	rand_ := rand.New(rand.NewSource(time.Now().UnixNano()))
	t2 := int(rand_.Intn(100000000))
	fmt.Println(t2)
	str2 := strconv.Itoa(t2)

	if len(str2) < 4 {
		str2 = str2[:len(str2)]
	} else {
		str2 = str2[:4]
	}

	str3 := str1 + str2
	result, _ := strconv.Atoi(str3)
	return int32(result)
}

func init() {
	//初始化随机数
	rand_ = rand.New(rand.NewSource(time.Now().UnixNano()))

	//csv初始化
	Csv = new(CsvConfig)
	Csv.Init()

	//json初始化
	Json_config = new(JsonConfig)
	Json_config.Init()

	//uid初始化
	err := redis.Find("uid", &Global_Uid)
	if err != nil {
		Global_Uid = createbeginUid()
		fmt.Println("Global_Uid", Global_Uid)
	} else {
		Global_Uid += 200 //解决宕机后Global_uid未存储
	}

	//全局world 初始化
	word = new(World)
	word.Init()

	//全局各个关卡挂机玩家
	global_guaji_players = new(GuajiPlayers)
	global_guaji_players.Init()

	//定时器
	timer_Manager = new(TimerManager)
	timer_Manager.Init()
}

func SendPackage(conn net.Conn, pid int32, body []byte) {

	len := 8 + len(body)
	var len_32 = int32(len)

	len_buf := bytes.NewBuffer([]byte{})
	binary.Write(len_buf, binary.BigEndian, len_32)

	pid_buf := bytes.NewBuffer([]byte{})
	binary.Write(pid_buf, binary.BigEndian, pid)

	msg := append(len_buf.Bytes(), pid_buf.Bytes()...)
	msg2 := append(msg, body...)
	conn.Write(msg2)
	fmt.Println("发送数据：id=", pid)
}

func GetUid() int32 {
	Global_Uid = Global_Uid + 1
	return Global_Uid
}

func Str2Int32(str string) int32 {

	data_int, error := strconv.Atoi(str)
	if error != nil {
		return 0
	}
	return int32(data_int)
}

func Str2Int(str string) int {

	data_int, error := strconv.Atoi(str)
	if error != nil {
		return 0
	}
	return data_int
}

func writeInfo(str string) {
	//fmt.Println("hahah", str)
	//Log.Error(str)
}

//产生随机数
func RandNum(min int32, max int32) int32 {
	if max < min {
		return 0
	}

	if min == max {
		return min
	}
	fmt.Println("RandNum min max", min, max)
	var num int32 = min
	num += int32(rand_.Intn(int(max - min)))
	return num
}

func GetRandomIndex(list []int32) int32 {
	var count int32 = 0
	for _, v := range list {
		count += v
	}
	rand_num_int := rand_.Intn(int(count))
	rand_num := int32(rand_num_int)
	var total_rand int32 = 0
	for i, v := range list {
		total_rand += v
		if rand_num <= total_rand {
			return int32(i)
		}
	}
	return 0
}

func Str2float32(str string) float32 {
	if len(str) <= 0 {
		return 0
	}

	f, err := strconv.ParseFloat(str, 32)
	if err != nil {
		writeInfo("can't to number %s=" + str)
		return 0
	}

	return float32(f)
}
