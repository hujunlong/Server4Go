package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
)

//需要写入的表格结构

var f *os.File
var distance_ints []int //切分字符串处理

type GetData struct {
	can_write bool   //该行能否写入
	tip       string //标题
}

func CreatePath(name string) {
	path := getCurrentDirectory()

	var symbol string = "/"
	if os.IsPathSeparator('\\') {
		symbol = "\\"
	}

	os.Mkdir(path+symbol+name, os.ModePerm)
}

//读取目录
func ListDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			if !strings.Contains(fi.Name(), "~") {
				files = append(files, dirPth+PthSep+fi.Name())
			}

		}
	}
	return files, nil
}

//当前路径
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//切分多个数组与单字符串
func getDistance(get_data_list []GetData, Cells []*xlsx.Cell) {

	var same_count int = 0
	var my_distance_int int = 0
	var last_id int = 0

	for i, v := range get_data_list {
		if strings.EqualFold(get_data_list[0].tip, v.tip) {
			same_count += 1
			last_id = i
		}
	}

	if same_count == 1 {
		distance_ints = append(distance_ints, len(get_data_list))
	} else {
		my_distance_int = (last_id / (same_count - 1)) * same_count
		distance_ints = append(distance_ints, my_distance_int)
		if my_distance_int < len(get_data_list) {
			getDistance(get_data_list[my_distance_int:], Cells[my_distance_int:])
		}
	}
}

//数组拼装
func assemblyArray(index_key int, same_count int, get_data_list []GetData, Cells []*xlsx.Cell) string {
	str := strconv.Itoa(index_key)
	var index_key_str string = ""
	index_key_str = `"` + "item" + str + `":`

	var begin_str string = index_key_str + `[`
	var end_str string = `]`

	distance := len(get_data_list) / same_count

	for i := 0; i < same_count; i++ {
		str = dealsingle(get_data_list[i*distance:(i+1)*distance], Cells[i*distance:(i+1)*distance])
		if i == 0 {
			begin_str += str
		} else {
			begin_str = begin_str + "," + str
		}
	}
	return begin_str + end_str
}

//单字符拼装
func dealsingle(get_data_list []GetData, Cells []*xlsx.Cell) string {
	var begin_str string = `{`
	var end_str string = `}`

	for i, v := range get_data_list {
		cell_str, _ := Cells[i].String()
		if i == 0 {
			begin_str = begin_str + `"` + v.tip + `":` + cell_str
		} else {
			begin_str = begin_str + `,"` + v.tip + `":` + cell_str
		}
	}
	return begin_str + end_str
}

//单字符拼装
func dealLast(get_data_list []GetData, Cells []*xlsx.Cell) string {
	var begin_str string = ""
	for i, v := range get_data_list {
		cell_str, _ := Cells[i].String()
		if i == 0 {
			begin_str = begin_str + `"` + v.tip + `":` + cell_str
		} else {
			begin_str = begin_str + `,"` + v.tip + `":` + cell_str
		}
	}
	return begin_str
}

//查询是否可读取
func findCanWrite(Rows []*xlsx.Row) []GetData {
	var get_data_list []GetData
	//第二行是否为空
	row := Rows[1]
	for _, cell := range row.Cells {
		var get_data GetData
		str, _ := cell.String()
		if len(str) > 0 {
			get_data.can_write = true
			get_data.tip = str
			get_data_list = append(get_data_list, get_data)
		}
	}
	return get_data_list
}

//判断有效行
func EffectiveNum(Cells []*xlsx.Row) int {
	var count int = 0
	for i, row := range Cells {
		if i < 2 { //前面两行不需要
			continue
		}

		//判断该行是否全部为空
		if isAllFull(row.Cells) {
			count += 1
		}
	}
	return count
}

//判断该行是否为空
func isAllFull(Cells []*xlsx.Cell) bool {
	if len(Cells) == 0 {
		return false
	}
	str, _ := Cells[0].String()
	if len(str) > 0 {
		return true
	}
	return false
}

func writeRow(f *os.File, json_str string) {
	f.WriteString(json_str)
}

func getSameCount(buff_get_data_list []GetData) int { //获取相同的
	var same_count int = 0
	for i := 0; i < len(buff_get_data_list); i++ {
		if strings.EqualFold(buff_get_data_list[0].tip, buff_get_data_list[i].tip) {
			same_count += 1
		}
	}
	return same_count
}

func dealjson(file_name string, save_name string) {
	//读取文件
	xlFile, _ := xlsx.OpenFile(file_name)

	//遍历循环
	for _, sheet := range xlFile.Sheets {
		get_data_list := findCanWrite(sheet.Rows)

		if len(get_data_list) > 0 {
			f, _ = os.Create("result/" + sheet.Name + ".json") //创建文件
			f.WriteString("{\n")
			fmt.Println("sheet.Name:", sheet.Name)
			defer f.Close()
		}

		count := EffectiveNum(sheet.Rows)
		distance_ints = nil
		getDistance(get_data_list[1:], sheet.Rows[1].Cells[1:])
		var total_str string = ""
		for i, row := range sheet.Rows[:count+2] {
			if i < 2 { //前面两行不需要
				continue
			}

			if i == count+2 {
				continue
			}

			//除去无用数据
			buff_get_data_list := get_data_list[1:]
			key, _ := row.Cells[0].String()
			row.Cells = row.Cells[1:]

			//单值
			if getSameCount(buff_get_data_list) == 1 {
				total_str = `"` + key + `":`
				if i != count+1 {
					total_str += dealsingle(buff_get_data_list, row.Cells) + ",\n"
				} else {
					total_str += dealsingle(buff_get_data_list, row.Cells) + "\n}" //最后一行不需要加,
				}
			} else {

				total_str = `"` + key + `":{`

				//获取遍历数组
				for j, v := range distance_ints {
					same_count := getSameCount(buff_get_data_list[:v])
					if same_count == 1 { //最后为单值
						if i != count+1 {
							total_str += dealLast(buff_get_data_list, row.Cells) + "},\n"
						} else {
							total_str += dealLast(buff_get_data_list, row.Cells) + "}\n}"
						}
					} else { //数组
						if j != len(distance_ints)-1 {
							total_str += assemblyArray(j, same_count, buff_get_data_list[:v], row.Cells[:v]) + ","
						} else { //最后一个值
							if i != count+1 {
								total_str += assemblyArray(j, same_count, buff_get_data_list[:v], row.Cells[:v]) + "},\n"
							} else {
								total_str += assemblyArray(j, same_count, buff_get_data_list[:v], row.Cells[:v]) + "}\n}"
							}
						}
						buff_get_data_list = buff_get_data_list[v:]
						row.Cells = row.Cells[v:]
					}
				}
			}
			writeRow(f, total_str)

		}

	}
}

func main() {

	CreatePath("result")
	//当前路径
	path := getCurrentDirectory()
	fires, _ := ListDir(path, "xlsx")

	for _, v := range fires {
		file_name := strings.TrimLeft(v, path)
		fmt.Println(v, file_name[1:])
		dealjson(v, file_name[1:])
	}

}
