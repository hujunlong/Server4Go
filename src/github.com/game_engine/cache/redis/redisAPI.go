package redis

import (
	"bytes"
	"encoding/gob"
)

var redis_ *Client

func init() {
	redis_ = new(Client)
}

func Add(key string, inter interface{}) error {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(inter)
	if err == nil {
		err = redis_.Set(key, buf.Bytes())
	}
	return err
}

func Modify(key string, inter interface{}) error {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(inter)
	if err == nil {
		err = redis_.Set(key, buf.Bytes())
	}
	return err
}

func Find(key string) []byte {
	data, _ := redis_.Get(key)
	return data
}
