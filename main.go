package main

import (
	"golang-songs/infrastructure"
	"log"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	//RDBに接続
	db, err := gorm.Open("mysql", os.Getenv("mysqlConfig"))
	if err != nil {
		log.Println(err)
	}
	db.DB().SetMaxIdleConns(10)
	defer db.Close()

	//リモートのRedisのコネクションプールの設定
	pool := &redis.Pool{
		MaxIdle:     3,
		MaxActive:   6,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			rc, err := redis.Dial("tcp", os.Getenv("REDIS_URL"))
			if err != nil {
				return nil, err
			}
			return rc, nil
		},
	}

	//サイドカーのRedisのコネクションプール設定
	pool2 := &redis.Pool{
		MaxIdle:     3,
		MaxActive:   6,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			rc, err := redis.Dial("tcp", os.Getenv("SIDECAR_REDIS_ADDRESS"))
			if err != nil {
				return nil, err
			}
			return rc, nil
		},
	}

	//リモートRedis
	// コネクションの取得
	redis := pool.Get()
	if redis.Err() != nil {
		log.Println(redis.Err())
	} else {
		log.Println("リモートのRedisに接続成功")
	}
	defer redis.Close()

	//サイドカーRedis
	// コネクションの取得
	sidecar_redis := pool2.Get()
	if sidecar_redis.Err() != nil {
		log.Println(sidecar_redis.Err())
	} else {
		log.Println("サイドカーのRedisに接続成功")
	}
	defer sidecar_redis.Close()

	infrastructure.Dispatch(db, redis, sidecar_redis)
}
