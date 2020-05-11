package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"api/pkg/setting"
	"fmt"
)

var (
	DB *gorm.DB
    DatabaseSetting = &Database{}
)

type Database struct {
	TYPE string
	USER string
	PASSWORD string
	HOST string
	NAME string
}

type Model struct {
	ID int `gorm:"primary_key" json:"id"`
}


func init() {
	var (
		err error
	)

	err = setting.Cfg.Section("database").MapTo(DatabaseSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo DatabaseSetting err: %v", err)
	}

	//sec, err := setting.Cfg.GetSection("database")
	//if err != nil {
	//	log.Fatal(2, "Fail to get section 'database': %v", err)
	//}
	//
	//dbType = sec.Key("TYPE").String()
	//dbName = sec.Key("NAME").String()
	//user = sec.Key("USER").String()
	//password = sec.Key("PASSWORD").String()
	//host = sec.Key("HOST").String()

	//DB, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
	//	user,
	//	password,
	//	host,
	//	dbName))

	DB, err = gorm.Open(DatabaseSetting.TYPE, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		DatabaseSetting.USER,
		DatabaseSetting.PASSWORD,
		DatabaseSetting.HOST,
		DatabaseSetting.NAME))

	if err != nil {
		log.Println(err)
	}

	gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
		return defaultTableName;
	}

	DB.SingularTable(true)
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(1000)
}

func CloseDB() {
	defer DB.Close()
}
