package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"

	cat "github.com/ctripcorp/cat.go"
	_ "github.com/ctripcorp/nephele/imgws/routers"
	"github.com/ctripcorp/nephele/util"
)

func ConfigCat() {
	cat.CAT_HOST = cat.UAT
	cat.DOMAIN = "900408"
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
	func() {
		Reboot := "Reboot"
		Cat := cat.Instance()
		tran := Cat.NewTransaction("System", Reboot)
		defer func() {
			tran.SetStatus("0")
			tran.Complete()
		}()
	}()
	go LogHeartbeat()
	beego.Run()
}

func LogHeartbeat() {
	defer func() {
		if err := recover(); err != nil {
			//log.Error(fmt.Sprintf("%v", err))
			LogHeartbeat()
		}
	}()

	ip := util.GetIP()
	second := time.Now().Second()
	if second < 29 {
		sleep := time.Duration((29 - second) * 1000000000)
		time.Sleep(sleep)
	}

	catinstance := cat.Instance()
	for {
		//log.Debug("send cat heartbeat")
		stats := util.GetStatus()
		tran := catinstance.NewTransaction("System", "Status")
		h := catinstance.NewHeartbeat("HeartBeat", ip)
		for key, value := range stats {
			switch key {
			case "Alloc", "TotalAlloc", "Sys", "Mallocs", "Frees", "OtherSys", "PauseNs":
				h.Set("System", key, value)
			case "HeapAlloc", "HeapSys", "HeapIdle", "HeapInuse", "HeapReleased", "HeapObjects":
				h.Set("HeapUsage", key, value)
			case "NextGC", "LastGC", "NumGC":
				h.Set("GC", key, value)
			}
		}
		h.SetStatus("0")
		h.Complete()
		tran.SetStatus("0")
		tran.Complete()
		second = time.Now().Second()
		sleep := time.Duration((90 - second) * 1000000000)
		time.Sleep(sleep)
	}
}
