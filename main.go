package main

import (
	"golang-songs/infrastructure"
	"io/ioutil"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/yaml.v2"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".envファイルの読み込み失敗")
	}

	yml, err := ioutil.ReadFile("conf/db.yml")
	if err != nil {
		log.Println("conf/db.ymlの読み込み失敗")
	}

	t := make(map[interface{}]interface{})

	_ = yaml.Unmarshal([]byte(yml), &t)

	//環境を取得
	conn := t[os.Getenv("GOJIENV")].(map[interface{}]interface{})

	db, err := gorm.Open("mysql", conn["user"].(string)+conn["password"].(string)+"@"+conn["rds"].(string)+"/"+conn["db"].(string)+"?charset=utf8&parseTime=True")
	if err != nil {
		log.Println(err)
	}

	infrastructure.Dispatch(db)

	db.DB().SetMaxIdleConns(10)
	defer db.Close()
}
