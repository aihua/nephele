package routers

import (
	"github.com/astaxie/beego"
	"github.com/ctripcorp/nephele/imgws/controllers"
)

func init() {
	ws := &controllers.ImageWS{}
	beego.Handler("/imagews.asmx", ws)
	imageHandler := &controllers.ImageHandler{}
	beego.Handler("/imagehandler.ashx", imageHandler)

	beego.Router("/", &controllers.ImageController{})
	beego.Include(&controllers.WhitelistController{})
	beego.Include(&controllers.ImageDeleteController{})
	beego.Include(&controllers.ConfigController{})
	beego.Include(&controllers.ChannelController{})
}
