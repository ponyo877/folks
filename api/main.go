package main

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/ponyo877/folks/api/handler"
	"github.com/ponyo877/folks/config"
	"github.com/ponyo877/folks/repository"
	"github.com/ponyo877/folks/usecase/message"

	"github.com/go-redis/redis/v9"
	"github.com/nats-io/nats.go"
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

	messageRepository := repository.NewMessageRepository(mq, kvs)
	messageService := message.NewService(messageRepository)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	handler.MakeRoomHandlers(e)
	handler.MakeMessageHandlers(e, messageService)
	e.Logger.Fatal(e.Start(":" + appConfig.APPort))
}
