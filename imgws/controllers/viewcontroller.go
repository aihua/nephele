package controllers

import (
	"github.com/astaxie/beego"
	cat "github.com/ctripcorp/cat.go"
	"github.com/ctripcorp/nephele/imgws/business"
	"github.com/ctripcorp/nephele/util"
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
	var result = util.Error{}
	Cat := cat.Instance()
	title := "/fdupload/target"
	if loadImgRequest.IsSource {
		title = "/fdupload/source"
	}
	tran := Cat.NewTransaction("URL", title)
	defer func() {
		if p := recover(); p != nil {
			Cat.LogPanic(p)
		}
		if result.Err != nil {
			tran.SetStatus(result.Err)
		} else {
			tran.SetStatus("0")
		}
		tran.Complete()
	}()
	util.LogEvent(Cat, "URL", "URL.Client", map[string]string{
		"clientip": util.GetClientIP(this.Ctx.Request),
		"serverip": util.GetIP(),
		"proto":    this.Ctx.Request.Proto,
		"referer":  this.Ctx.Request.Referer(),
		//"agent":    request.UserAgent(),
	})
	util.LogEvent(Cat, "URL", "URL.Method", map[string]string{
		"Http": this.Ctx.Request.Method + " " + loadImgRequest.FilePath,
	})

	imgRequest := business.ImageRequest{}
	imgRequest.Cat = Cat
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
