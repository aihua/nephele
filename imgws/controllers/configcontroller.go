package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/ctripcorp/nephele/imgws/models"
)

type ConfigController struct {
	beego.Controller
}

// @router /config/add/ [get]
func (this *ConfigController) Add() {
	channel := this.GetString("channel")
	key := this.GetString("key")
	value := this.GetString("value")
	config := models.Config{Channel: channel, Key: key, Value: value}
	err := config.Insert()
	if err != nil {
		this.Ctx.WriteString(err.Error())
	} else {
		this.Ctx.WriteString("sucess")
	}
}

// @router /config/get/ [get]
func (this *ConfigController) Get() {
	channel := this.GetString("channel")
	key := this.GetString("key")

	m, e := models.GetConfigs()
	if e.Err != nil {
		this.Ctx.WriteString(fmt.Sprintf("%v", e.Err))
	}
	var result interface{}
	if channel != "" {
		m1, exists := m[channel]
		if !exists {
			this.Ctx.WriteString("No record")
			return
		}
		if key != "" {
			m2, exists := m1[key]
			if !exists {
				this.Ctx.WriteString("No record")
				return
			}
			result = m2
		} else {
			result = m1
		}
	} else {
		result = m
	}
	bts, err := json.Marshal(result)
	if err != nil {
		this.Ctx.WriteString(err.Error())
	} else {
		this.Ctx.WriteString(string(bts))
	}
}

// @router /config/update/ [get]
func (this *ConfigController) Update() {
	channel := this.GetString("channel")
	key := this.GetString("key")
	value := this.GetString("value")
	if channel == "" || key == "" {
		this.Ctx.WriteString("params is invalid")
		return
	}
	config := models.Config{Channel: channel, Key: key, Value: value}
	err := config.UpdateValue()
	if err != nil {
		this.Ctx.WriteString(err.Error())
	} else {
		this.Ctx.WriteString("sucess")
	}
}
