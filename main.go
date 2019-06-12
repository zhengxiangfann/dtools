package main

import (
	"dtools/models/admin"
	"dtools/models/api"
	"dtools/models/auth"
	"dtools/models/code"
	"dtools/models/env"
	"dtools/models/group"
	"dtools/models/role"
	"dtools/models/template"

	_ "dtools/routers"
	"dtools/utils"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/patrickmn/go-cache"
	"net/url"
	"time"
)

func Init() {
	dbhost := beego.AppConfig.String("db.host")
	dbport := beego.AppConfig.String("db.port")
	dbuser := beego.AppConfig.String("db.user")
	dbpassword := beego.AppConfig.String("db.password")
	dbname := beego.AppConfig.String("db.name")
	timezone := beego.AppConfig.String("db.timezone")
	if dbport == "" {
		dbport = "3306"
	}
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"
	// fmt.Println(dsn)

	if timezone != "" {
		dsn = dsn + "&loc=" + url.QueryEscape(timezone)
	}

	fmt.Println(dsn)
	orm.RegisterDataBase("default", "mysql", dsn)

	orm.RegisterModel(
		new(auth.Auth),
		new(role.Role),
		new(role.RoleAuth),
		new(admin.Admin),
		new(group.Group),
		new(env.Env),
		new(code.Code),
		new(api.ApiSource),
		new(api.ApiDetail),
		new(api.ApiPublic),
		new(template.Template),
	)
	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}
}

func main() {
	Init()
	utils.Che = cache.New(60*time.Minute, 120*time.Minute)
	beego.Run()
}
