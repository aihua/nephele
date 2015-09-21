package models

import (
	"github.com/astaxie/beego/orm"
	"testing"
)

func InitDB() {
	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", "root:@/imagedb?charset=utf8")
	orm.RegisterModel(new(Config))
}

func TestWhitelistAddSize(t *testing.T) {
	InitDB()
	var conf Config = Config{
		Channel: "tg",
		Key:     "sizes",
	}
	conf.AddSize("320X160")
	conf.AddSize("200X200")
}
