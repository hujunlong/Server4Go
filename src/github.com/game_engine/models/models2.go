package models

import (
	"github.com/astaxie/beego/orm"
)

type User struct {
	Id   int
	Name string
}

func init() {
	orm.RegisterModel(new(User))
	orm.RunSyncdb("default", false, true)
}

func Add() {
	user := new(User)
	user.Name = "2010"
	o.Insert(user)
}
