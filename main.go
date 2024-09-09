package main

import (
	"fmt"
	"log"
	"testhub-spec-uni/controllers"
	"testhub-spec-uni/middleware"
	_ "testhub-spec-uni/routers"

	"github.com/astaxie/beego/orm"
	"github.com/beego/beego/v2/server/web"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
	_ "github.com/lib/pq"
)

func init() {
	driverName, err := web.AppConfig.String("db_driver")
	if err != nil {
		log.Fatalf("Failed to get 'db_driver': %v", err)
	}

	user, err := web.AppConfig.String("db_user")
	if err != nil {
		log.Fatalf("Failed to get 'db_user': %v", err)
	}

	password, err := web.AppConfig.String("db_password")
	if err != nil {
		log.Fatalf("Failed to get 'db_password': %v", err)
	}

	host, err := web.AppConfig.String("db_host")
	if err != nil {
		log.Fatalf("Failed to get 'db_host': %v", err)
	}

	port, err := web.AppConfig.String("db_port")
	if err != nil {
		log.Fatalf("Failed to get 'db_port': %v", err)
	}

	dbName, err := web.AppConfig.String("db_name")
	if err != nil {
		log.Fatalf("Failed to get 'db_name': %v", err)
	}

	schema := "uni_spec,accounts"
	dataSource := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable search_path=%s",
		user, password, host, port, dbName, schema)

	orm.RegisterDataBase("default", driverName, dataSource)

	fmt.Println(driverName)
}

func main() {
	web.BConfig.WebConfig.DirectoryIndex = true
	web.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"

	beego.BConfig.RouterCaseSensitive = false
	beego.SetStaticPath("/swagger", "swagger")

	beego.InsertFilter("/api/*", beego.BeforeRouter, middleware.AuthMiddleware)
	beego.InsertFilter("/user/universities/*", beego.BeforeRouter, middleware.AuthMiddleware)

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: []string{"http://localhost:3000", "https://admin-course.testhub.kz",
			"https://ent.testhub.kz", "https://console.ps.kz", "https://api-dev.testhub.kz",
			"https://dev-front.testhub.kz"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "lang"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	beego.Include(&controllers.SubjectController{})
	beego.Include(&controllers.SpecialityController{})
	beego.Include(&controllers.UniversityController{})
	beego.Include(&controllers.CityController{})
	beego.Include(&controllers.QuotaController{})

	beego.Run()

}
