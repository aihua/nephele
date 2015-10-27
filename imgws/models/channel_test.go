package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"testing"
	"time"
)

var isinit bool = false

func InitDBForTest() {
	if !isinit {
		orm.RegisterDriver("mysql", orm.DR_MySQL)
		orm.RegisterDataBase("default", "mysql", "root:@/imagedb?charset=utf8")
		//orm.Debug = true
		//orm.RegisterModel(new(models.Config))
		isinit = true
	}
}

var channel Channel = Channel{}

func TestGetChannels(t *testing.T) {
	InitDBForTest()
	m, e := channel.GetChannels()
	if e.Err != nil {
		t.Error(e.Err)
	}
	mcount := len(m)

	//insert
	o := orm.NewOrm()
	o.Raw("INSERT INTO channel(name,code) VALUES(?,?)", "test", "ZZ").Exec()
	m1, e := channel.GetChannels()
	m1count := len(m1)
	if mcount != m1count {
		t.Error("refresh error")
	}

	//refresh time  get
	getChannelsTime = time.Now().Add(-2 * time.Minute)
	m2, e := channel.GetChannels()
	m2count := len(m2)
	if m2count == mcount {
		t.Error("refresh fail")
	}

	//delete
	o.Raw("DELETE from channel WHERE name=?", "test").Exec()
}

func TestIsRefresh(t *testing.T) {
	getChannelsTime = time.Now()
	isrefresh := IsRefresh(getChannelsTime)
	if isrefresh {
		t.Error("fail")
	}

	getChannelsTime = time.Now().Add(-2 * time.Minute)
	isrefresh = IsRefresh(getChannelsTime)
	if !isrefresh {
		t.Error("fail")
	}
}
