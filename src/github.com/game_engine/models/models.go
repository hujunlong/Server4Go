package models

import (
	"github.com/astaxie/beego/orm"
)

type Profile struct {
	Id  int
	Age int16
}

func init() {
	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", "root:game9z@tcp(192.168.1.207:3306)/orm_test?charset=utf8")
	orm.RegisterModel(new(Profile))

}

func Insert() {
	o := orm.NewOrm()

	profile := new(Profile)
	profile.Age = 50
	o.Insert(profile)

}
