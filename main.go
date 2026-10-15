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
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/polar-bear-cu/sgt-noti-service/config"
	"github.com/polar-bear-cu/sgt-noti-service/consumer"
	"github.com/polar-bear-cu/sgt-noti-service/controllers"
	"github.com/polar-bear-cu/sgt-noti-service/mailer"
	"github.com/polar-bear-cu/sgt-noti-service/repositories"
	"github.com/polar-bear-cu/sgt-noti-service/routes"
	"github.com/polar-bear-cu/sgt-noti-service/usecases"
)

func main() {
	cfg := config.Load()
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET required")
	}

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
	rmqClosed := rmqConn.NotifyClose(make(chan *amqp.Error, 1))

	db := mongoClient.Database(cfg.MongoDBName)
	emailRepo := repositories.NewEmailNotificationMongo(db)
	if err := emailRepo.EnsureIndexes(ctx); err != nil {
		log.Fatal(err)
	}
	inAppRepo := repositories.NewInAppNotificationMongo(db)
	if err := inAppRepo.EnsureIndexes(ctx); err != nil {
		log.Fatal(err)
	}
	mail := mailer.NewSMTPMailer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPFrom, cfg.SMTPUsername, cfg.SMTPPassword)
	emailUC := usecases.NewEmailNotification(mail, emailRepo)
	inAppUC := usecases.NewInAppNotification(inAppRepo)
	inAppController := controllers.NewInAppNotification(inAppUC)
	cons := consumer.NewRabbitMQConsumer(rmqConn, emailUC, inAppUC)

	go func() {
		if err := cons.Run(ctx); err != nil {
			log.Printf("consumer stopped: %v", err)
		}
	}()

	userDeleted := consumer.NewUserDeletedConsumer(rmqConn, emailUC, inAppUC)
	go func() {
		if err := userDeleted.Run(ctx); err != nil {
			log.Printf("user_deleted consumer stopped: %v", err)
		}
	}()

	r := gin.Default()
	routes.Register(r, inAppController, cfg.JWTSecret)
	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r}

	go func() {
		log.Println("listening :" + cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	select {
	case <-ctx.Done():
	case amqpErr := <-rmqClosed:
		log.Fatalf("rabbitmq connection closed: %v", amqpErr)
	}
	stop()
	log.Println("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("http shutdown: %v", err)
	}
}
