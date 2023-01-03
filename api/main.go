package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/ponyo877/folks/api/handler"
	"github.com/ponyo877/folks/config"
	"github.com/ponyo877/folks/repository"
	"github.com/ponyo877/folks/usecase/message"

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
	mq, err := nats.Connect("nats://" + natsConfig.MQHost + ":" + natsConfig.MQPort)
	if err != nil {
		log.Panicf("NATSの接続に失敗しました: %v", err)
	}

	messageRepository := repository.NewMessageNats(mq)
	messageService := message.NewService(messageRepository)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	handler.MakeRoomHandlers(e)
	handler.MakeMessageHandlers(e, messageService)
	e.Logger.Fatal(e.Start(":" + appConfig.APPort))
}
