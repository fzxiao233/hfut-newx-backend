package db

import (
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/xormplus/xorm"
)

var DB *xorm.Engine

func init() {
	//open a DB connection
	var err error
	DB, err = xorm.NewEngine("mysql", "root:root@tcp(127.0.0.1:8889)/newx?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Printf("%s", err)
		panic("failed to connect database")
	}
}
