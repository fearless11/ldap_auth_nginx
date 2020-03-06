package controllers

import (
	"encoding/base64"
	"fmt"
	"strings"

	"ldap_auth_nginx/models"

	"github.com/astaxie/beego"
)

type AuthController struct {
	beego.Controller
}

func (c *AuthController) Auth() {
	clientIP := c.Ctx.Input.IP()
	token := c.Ctx.GetCookie("nginxauth")
	auth, err := base64.StdEncoding.DecodeString(token)
	if err != nil || len(token) == 0 {
		// nginx是否缓存有服务器与源服务器共同决定
		//header设置为Cache-control：no-cache、no-store时 nginx不会缓存
		c.Ctx.Output.Header("Cache-Control", "no-cache")
		c.Ctx.Abort(401, "401")
	}

	info := strings.Split(string(auth), ":")
	username := info[0]
	password := info[1]

	err = models.LDAP_Auth(username, password)
	if err != nil {
		c.Ctx.Output.Header("Cache-Control", "no-cache")
		beego.Error(fmt.Sprintf("%v %v auth fail", username, clientIP))
		c.Ctx.Abort(401, "401")
	}

	beego.Info(fmt.Sprintf("%v %v auth success ", username, clientIP))
	c.Ctx.Output.Body([]byte("ok"))
}
