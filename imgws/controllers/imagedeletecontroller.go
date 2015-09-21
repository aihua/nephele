package controllers

import (
	"errors"

	"github.com/astaxie/beego"

	//"github.com/ctripcorp/nephele/fdfs"
	"github.com/ctripcorp/nephele/imgws/models"
)

var (
	ErrIllegalStorageType = errors.New("illegal storage type")
)

type ImageDeleteController struct {
	beego.Controller
}

// @router /image/delete/:uri [get]
func (c *ImageDeleteController) Get() {
	var (
		st string
		sp string
	)
	uri := c.GetString(":uri")
	var ii models.ImageIndex = models.ImageIndex{}
	err := ii.Parse(uri)
	if err.Type != "" {
		return
	}

	st, sp = ii.StorageType, ii.StoragePath
	if st == EmptyString || sp == EmptyString {
		return
	}

	switch st {
	case "fdfs":

	case "nfs":

	default:
		//ErrIllegalStorageType
		return
	}
	println("delete record from MySQL")
}
