package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "https://ent.testhub.kz"
	c.Data["Email"] = "superalibek123@gmail.com"
	c.TplName = "index.tpl"
}
