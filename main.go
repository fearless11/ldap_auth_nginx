package main

import (
	_ "ldap_auth_nginx/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
