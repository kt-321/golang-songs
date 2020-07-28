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

	infrastructure.Dispatch(db)

	db.DB().SetMaxIdleConns(10)
	defer db.Close()
}
