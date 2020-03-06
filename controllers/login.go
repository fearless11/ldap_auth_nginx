package controllers

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) LoginGet() {
	// 默认情况下 ReadFromRequest 函数已经实现了读取的数据赋值给 flash
	// 模板读取数据 {{.flash.warning}}  {{.flash.notice}}
	flash := beego.ReadFromRequest(&c.Controller)
	if n, ok := flash.Data["notice"]; ok {
		c.Data["msg"] = n
	}

	target := c.Ctx.Input.Header("X-Target")
	getTarget := c.GetString("target")
	if target == "" && getTarget == "" {
		target = "/"
	}
	if getTarget != "" {
		target = getTarget
	}

	if len(target) == 0 {
		c.Ctx.Abort(500, "500")
	}

	c.Data["target"] = target
	c.TplName = "login.html"
}

func (c *LoginController) LoginPost() {
	flash := beego.NewFlash()
	c.Ctx.Request.ParseForm()
	username := strings.TrimSpace(c.Ctx.Request.Form.Get("username"))
	password := strings.TrimSpace(c.Ctx.Request.Form.Get("password"))
	target := c.Ctx.Request.Form.Get("target")

	if len(username) == 0 || len(password) == 0 {
		flash.Notice("用户密码不能为空")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, fmt.Sprintf("/login?target=%v", target))
	}

	if _, ok := beego.BConfig.WhiteMap[username]; !ok {
		flash.Notice("无权限")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, fmt.Sprintf("/login?target=%v", target))
	}

	auth := fmt.Sprintf("%v:%v", username, password)
	token := base64.StdEncoding.EncodeToString([]byte(auth))
	c.Ctx.SetCookie("nginxauth", token)
	c.Ctx.Redirect(302, target)
}
