package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

var isinit bool = false

func InitDB() {
	if !isinit {
		orm.RegisterDriver("mysql", orm.DR_MySQL)
		orm.RegisterDataBase("default", "mysql", "root:@/imagedb?charset=utf8")
		orm.RegisterModel(new(Config))
		isinit = true
	}
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
