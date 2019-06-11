package main

import (
	"dtools/models/admin"
	_ "dtools/routers"
	"dtools/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
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

	orm.RegisterDataBase("default", "mysql", dsn)

	orm.RegisterModel(
		new(admin.Admin),
		//new(Role),
		//		//new(RoleAuth),
		//		//new(Admin),
		//		//new(Group),
		//		//new(Env),
		//		//new(Code),
		//		//new(ApiSource),
		//		//new(ApiDetail),
		//		//new(ApiPublic),
		//		//new(Template),
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