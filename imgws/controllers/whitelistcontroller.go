package controllers

import (
	"errors"
	"regexp"

	"github.com/astaxie/beego"

	"github.com/ctripcorp/nephele/imgws/models"
)

var (
	ErrExistedSize    = errors.New("size already exists in whitelist")
	ErrConflict       = errors.New("data updates conflict")
	ErrIllegalChannel = errors.New("illegal channel")
	ErrIllegalSize    = errors.New("illegal size")
)

var (
	TgPattern   = "^[a-zA-Z]+$"
	SizePattern = "^[1-9]\\d*X[1-9]\\d*$"
)

type WhitelistController struct {
	beego.Controller
}

// @router /whitelist [get]
func (c *WhitelistController) Get() {
	c.Ctx.WriteString(`Whitelist API Reference:
{
	"Get All Sizes" :
	{
		"path"    : "/whitelist/getsizes",
		"pattern" : "/whitelist/getsizes?channel={channel name}"
		"example" : "/whitelist/getsizes?channel=tg"
	},
	"Add Size to Whitelist" :
	{
		"path"    : "/whitelist/add",
		"pattern" : "/whitelist/add?channel={channel name}&size={new size}"
		"example" : "/whitelist/add?channel=tg&size=200X100"
	}
}`)
}

// @router /whitelist/getsizes/ [get]
func (c *WhitelistController) GetSizes() {
	channel := c.GetString("channel")
	isMatch, err := regexp.Match(TgPattern, []byte(channel))
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}
	if !isMatch {
		//ErrIllegalChannel
		c.Ctx.WriteString(ErrIllegalChannel.Error())
		return
	}
	var conf models.Config = models.Config{
		ChannelCode: channel,
		Key:         "sizes",
	}
	sizes, err := conf.GetSizes()
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}
	c.Ctx.WriteString(sizes)
}

// @router /whitelist/add/ [get]
func (c *WhitelistController) Add() {
	channel := c.GetString("channel")
	isMatch, err := regexp.Match(TgPattern, []byte(channel))
	if err != nil {
		//err
		c.Ctx.WriteString(err.Error())
		return
	}
	if !isMatch {
		//ErrIllegalChannel
		c.Ctx.WriteString(ErrIllegalChannel.Error())
		return
	}
	size := c.GetString("size")
	isMatch, err = regexp.Match(SizePattern, []byte(size))
	if err != nil {
		//err
		c.Ctx.WriteString(err.Error())
		return
	}
	if !isMatch {
		//ErrIllegalSize
		c.Ctx.WriteString(ErrIllegalSize.Error())
		return
	}

	var conf models.Config = models.Config{
		ChannelCode: channel,
		Key:         "sizes",
	}
	err = conf.AddSize(size)
	if err != nil {
		//err
		c.Ctx.WriteString(conf.Value)
		c.Ctx.WriteString(err.Error())
		return
	}
	c.Ctx.WriteString(conf.Value)
}
