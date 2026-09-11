package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/polar-bear-cu/sgt-noti-service/config"
	"github.com/polar-bear-cu/sgt-noti-service/consumer"
	"github.com/polar-bear-cu/sgt-noti-service/mailer"
	"github.com/polar-bear-cu/sgt-noti-service/repositories"
	"github.com/polar-bear-cu/sgt-noti-service/routes"
	"github.com/polar-bear-cu/sgt-noti-service/usecases"
)

func main() {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mongoClient, err := config.ConnectMongo(cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = mongoClient.Disconnect(context.Background()) }()

	rmqConn, err := config.ConnectRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = rmqConn.Close() }()

	logs := repositories.NewNotificationLogMongo(mongoClient.Database(cfg.MongoDBName))
	mail := mailer.NewSMTPMailer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPFrom, cfg.SMTPUsername, cfg.SMTPPassword)
	uc := usecases.NewNotification(mail, logs)
	cons := consumer.NewRabbitMQConsumer(rmqConn, uc)

	go func() {
		if err := cons.Run(ctx); err != nil {
			log.Printf("consumer stopped: %v", err)
		}
	}()

	r := gin.Default()
	routes.Register(r)

	log.Println("listening :" + cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
