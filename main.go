package main

import (
	"fmt"
	"os"
	"testhub-spec-uni/conf"
	_ "testhub-spec-uni/routers"

	"github.com/astaxie/beego/orm"
	"github.com/beego/beego/v2/server/web"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
)

func init() {
	driverName := os.Getenv("DB_DRIVER")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dataSource := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable search_path=uni_spec",
		user, password, host, port, dbName)

	orm.RegisterDataBase("default", driverName, dataSource)
	orm.RunSyncdb("default", false, true)
	fmt.Println(driverName)

	conf.InitElasticsearch()
}

func main() {
	web.BConfig.WebConfig.DirectoryIndex = true
	web.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"

	beego.BConfig.RouterCaseSensitive = false
	beego.SetStaticPath("/swagger", "swagger")

	//beego.InsertFilter("/api/*", web.BeforeRouter, middleware.AuthMiddleware)
	beego.Run()
}
