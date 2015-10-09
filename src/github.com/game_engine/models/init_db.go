package models

import (
	"github.com/astaxie/beego/orm"
)

var o orm.Ormer

func init() {
	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", "root:game9z@tcp(192.168.1.207:3306)/orm_test?charset=utf8")
	o = orm.NewOrm()
}
