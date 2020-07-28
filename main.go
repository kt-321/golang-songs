package main

import (
	"golang-songs/infrastructure"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	db, err := gorm.Open("mysql", os.Getenv("mysqlConfig"))
	if err != nil {
		log.Println(err)
		log.Println(os.Getenv("mysqlConfig"))
		log.Println("mysqlConfig")
	}

	//一旦直打ち
	//db, err := gorm.Open("mysql", "kt321:route666@tcp(database-2.cvzte0rjvtt7.ap-northeast-1.rds.amazonaws.com:3306)/db_goyoursongs_rds?charset=utf8&parseTime=True")
	//if err != nil {
	//	log.Println(err)
	//}

	infrastructure.Dispatch(db)

	db.DB().SetMaxIdleConns(10)
	defer db.Close()
}
