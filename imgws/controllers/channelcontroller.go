package controllers

import (
	"github.com/astaxie/beego"
	_ "github.com/ctripcorp/nephele/imgws/models"
)

type ChannelController struct {
	beego.Controller
}

func (this *ChannelController) Add() {
	name := this.GetString("name")
	code := this.GetString("code")
	if name == "" || code == "" {
		this.Ctx.WriteString("params isn't empty!")
		return
	}
	//channel := models.Channel{}

}

func (this *ChannelController) Get() {

}
func (this *ChannelController) Update() {

}
