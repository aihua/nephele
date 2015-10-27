package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/ctripcorp/nephele/imgws/models"
)

type ChannelController struct {
	beego.Controller
}

// @router /channel/add/ [get]
func (this *ChannelController) Add() {
	name := this.GetString("name")
	code := this.GetString("code")
	if name == "" || code == "" {
		this.Ctx.WriteString("params isn't empty!")
		return
	}
	channel := models.Channel{Name: name, Code: code}
	err := channel.Insert()
	if err != nil {
		this.Ctx.WriteString(err.Error())
	} else {
		this.Ctx.WriteString("success")
	}
}

// @router /channel/get/ [get]
func (this *ChannelController) Get() {
	ch := models.Channel{}
	channels, err := ch.GetAll()
	if err.Err != nil {
		this.Ctx.WriteString((err.Err.(error)).Error())
	} else {
		bts, _ := json.Marshal(channels)
		this.Ctx.WriteString(string(bts))
	}
}

// @router /channel/update/ [get]
func (this *ChannelController) Update() {
	name := this.GetString("name")
	code := this.GetString("code")
	if name == "" || code == "" {
		this.Ctx.WriteString("params is't empty")
	}
	channel := models.Channel{Name: name, Code: code}
	err := channel.Update()
	if err != nil {
		this.Ctx.WriteString(err.Error())
	} else {
		this.Ctx.WriteString("success")
	}
}
