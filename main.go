package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

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
	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r}

	go func() {
		log.Println("listening :" + cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	stop()
	log.Println("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("http shutdown: %v", err)
	}
}
