package models

import (
	"github.com/astaxie/beego/orm"
)

type Profile struct {
	Id  int
	Age int16
}

func init() {
	orm.RegisterModel(new(Profile))
	orm.RunSyncdb("default", false, true) //创建表
}

func Insert() {
	profile := new(Profile)
	profile.Age = 50
	o.Insert(profile)
}
