package models

import (
	//"database/sql"
	//"github.com/astaxie/beego/orm"
	"testing"
)

func TestWhitelistAddSize(t *testing.T) {
	InitDBForTest()
	var conf Config = Config{
		ChannelCode: "10",
		Key:         "sizes",
	}
	conf.AddSize("320X160")
	conf.AddSize("200X200")
}

func TestGetConfigs(t *testing.T) {
	config := Config{}
	configs, _ := config.GetConfigs()

	_, exists := configs["10"]
	if !exists {
		t.Error("error")
	}
}
