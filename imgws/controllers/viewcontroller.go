package controllers

import (
	"github.com/astaxie/beego"
	"github.com/ctripcorp/nephele/imgws/business"
	"github.com/ctripcorp/nephele/util/soapparse"
	"github.com/ctripcorp/nephele/util/soapparse/request"
)

type ViewController struct {
	beego.Controller
}

// @router /fdupload/target/* [get]
func (this *ViewController) View() {
	path := "/" + this.GetString(":splat")
	loadImgRequest := request.LoadImgRequest{}
	loadImgRequest.FilePath = path
	loadImgRequest.IsSource = false
	View(this, &loadImgRequest)
}

// @router /fdupload/source/* [get]
func (this *ViewController) ViewSource() {
	path := "/" + this.GetString(":splat")

	loadImgRequest := request.LoadImgRequest{}
	loadImgRequest.FilePath = path
	loadImgRequest.IsSource = true
	View(this, &loadImgRequest)
}

func View(this *ViewController, loadImgRequest *request.LoadImgRequest) {
	imgRequest := business.ImageRequest{}
	resp, e := imgRequest.Download(loadImgRequest)
	if e.Err != nil {
		this.Ctx.WriteString(e.Err.(error).Error())
	} else {
		//this.Ctx.Output.ContentType("image/Jpeg")

		//this.Ctx.Output.Body(resp.FileBytes)
		bts, err := soapparse.B64.DecodeString(string(resp.FileBytes))
		if err != nil {
			this.Ctx.WriteString(err.Error())
		} else {
			this.Ctx.Output.Header("Content-Type", "image/jpeg")
			this.Ctx.Output.Body(bts)
		}
	}
}
