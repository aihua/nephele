package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"

	cat "github.com/ctripcorp/cat.go"
	//"github.com/ctripcorp/nephele/imgws/models"
	_ "github.com/ctripcorp/nephele/imgws/routers"
)

func ConfigCat() {
	cat.CAT_HOST = cat.UAT
	cat.DOMAIN = "900407"
	cat.TEMPFILE = ".cat"
}

func InitDB() {
	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", "root:@/imagedb?charset=utf8")
	//orm.RegisterModel(new(models.Config))
}

func main() {
	ConfigCat()
	InitDB()
	beego.Run()
}
