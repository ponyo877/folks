package main

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/ponyo877/folks/api/handler"
	"github.com/ponyo877/folks/config"
	"github.com/ponyo877/folks/repository"
	"github.com/ponyo877/folks/usecase/room"

	"github.com/go-redis/redis/v9"
	"github.com/nats-io/nats.go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	appConfig, err := config.LoadAppConfig()
	if err != nil {
		log.Panicf("LoadAppConfigに失敗しました: %v", err)
	}
	natsConfig, err := config.LoadNatsConfig()
	if err != nil {
		log.Panicf("LoadNatsConfigに失敗しました: %v", err)
	}
	mq, err := nats.Connect(
		"nats://"+natsConfig.MQHost+":"+natsConfig.MQPort,
		nats.PingInterval(20*time.Second),
		nats.MaxPingsOutstanding(5),
	)
	if err != nil {
		log.Panicf("NATSの接続に失敗しました: %v", err)
	}
	redisConfig, err := config.LoadRedisConfig()
	if err != nil {
		log.Panicf("LoadRedisConfigに失敗しました: %v", err)
	}
	kvs := redis.NewClient(&redis.Options{
		Addr:     redisConfig.KVSHost + ":" + redisConfig.KVSPort,
		Password: redisConfig.KVSPassword,
		DB:       redisConfig.KVSDatabase,
	})
	mysqlConfig, err := config.LoadMysqlConfig()
	if err != nil {
		log.Panicf("LoadMysqlConfigに失敗しました: %v", err)
	}
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", mysqlConfig.DBUser, mysqlConfig.DBPassword, mysqlConfig.DBHost, mysqlConfig.DBPort, mysqlConfig.DBDatabase)
	gormDB, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		log.Panicf("gormDBクライアント作成に失敗しました: %v", err)
	}
	rdb, err := gormDB.DB()
	if err != nil {
		log.Panicf("SQLクライアント作成に失敗しました: %v", err)
	}
	defer rdb.Close()

	messageRepository := repository.NewMessageRepository(mq, kvs, gormDB)
	messageService := room.NewService(messageRepository)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	handler.MakeRoomHandlers(e, messageService)
	e.Logger.Fatal(e.Start(":" + appConfig.APPort))
}
