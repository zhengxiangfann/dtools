package routers

import (
	"dtools/controllers/auth"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &auth.AuthController{}, "*:Index")
}
