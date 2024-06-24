package main

import (
	"fmt"
	"log"
	_ "testhub-spec-uni/routers"

	"github.com/astaxie/beego/orm"
	"github.com/beego/beego/v2/server/web"
	beego "github.com/beego/beego/v2/server/web"
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

	dataSource := "user=" + user + " password=" + password + " host=" + host + " port=" + port + " dbname=" + dbName + " sslmode=disable"

	orm.RegisterDataBase("default", driverName, dataSource)
	orm.RunSyncdb("default", false, true)
	fmt.Println(driverName)

}

func main() {
	beego.Run()
}
