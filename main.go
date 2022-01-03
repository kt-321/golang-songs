package main

import (
	"golang-songs/infrastructure"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	// RDBに接続.
	db, err := gorm.Open("mysql", os.Getenv("mysqlConfig"))
	if err != nil {
		log.Println("RDBの接続失敗:", errors.WithStack(err))
	} else {
		log.Println("RDBの接続成功")
	}

	db.DB().SetMaxIdleConns(10)

	defer db.Close()

	// リモートのRedisのコネクションプールの設定.
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

	// サイドカーのRedisのコネクションプール設定.
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

	// リモートRedis.
	// コネクションの取得.
	redis := pool.Get()
	if redis.Err() != nil {
		log.Println("リモートのRedisに接続失敗: ", redis.Err())
	} else {
		log.Println("リモートのRedisに接続成功")
	}
	defer redis.Close()

	// サイドカーRedis.
	// コネクションの取得.
	sidecarRedis := pool2.Get()
	if sidecarRedis.Err() != nil {
		log.Println("サイドカーのRedisに接続失敗: ", sidecarRedis.Err())
	} else {
		log.Println("サイドカーのRedisに接続成功")
	}
	defer sidecarRedis.Close()

	// 予めtime.Localにタイムゾーンの設定情報を入れておく.
	time.Local = time.FixedZone("Local", 9*60*60)

	infrastructure.Dispatch(db, redis, sidecarRedis)
}
