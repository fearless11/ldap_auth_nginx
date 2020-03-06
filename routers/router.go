package routers

import (
	"ldap_auth_nginx/controllers"

	"github.com/astaxie/beego"
)

func init() {
	// 固定路由
	beego.Router("/", &controllers.MainController{})
	// 自动路由
	beego.Router("/auth", &controllers.AuthController{}, "get:Auth")
	beego.Router("/login", &controllers.LoginController{}, "get:LoginGet;post:LoginPost")
}
