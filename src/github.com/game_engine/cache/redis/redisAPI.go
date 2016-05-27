package redis

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
)

var redis_ *Client

func init() {
	redis_ = new(Client)
}

func dealInterface(key interface{}) string {
	var str string = ""

	switch key.(type) {
	case string:
		str = key.(string)
	case int32:
		str = strconv.Itoa(int(key.(int32)))
	case int:
		str = strconv.Itoa(key.(int))
	case int64:
		str = strconv.FormatInt(key.(int64), 10)
	default:
	}

	return str
}

func Modify(key interface{}, inter interface{}) error {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(inter)

	if err == nil {
		str := dealInterface(key)
		redis_.Set(str, buf.Bytes())
	}
	return err
}

func Find(key interface{}, inter interface{}) error {

	str := dealInterface(key)
	data, err := redis_.Get(str)

	if err == nil {
		buf := bytes.NewBuffer(data)
		dec := gob.NewDecoder(buf)
		dec.Decode(inter)
	} else {
		fmt.Println("err:", err)
	}
	return err
}

func Incr(key interface{}) (int64, error) {
	str := dealInterface(key)
	id, err := redis_.Incr(str)
	return id, err
}

func Del(key interface{}) (bool, error) {
	str := dealInterface(key)
	ok, err := redis_.Del(str)
	return ok, err
}

func Exists(key interface{}) (bool, error) {
	str := dealInterface(key)
	result, err := redis_.Exists(str)
	return result, err
}
