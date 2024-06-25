package main

import (
	_ "testhub-spec-uni/routers"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run()
}
