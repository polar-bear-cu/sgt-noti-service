package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/polar-bear-cu/sgt-noti-service/config"
	"github.com/polar-bear-cu/sgt-noti-service/routes"
)

func main() {
	cfg := config.Load()

	r := gin.Default()
	routes.Register(r)

	log.Println("listening :" + cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
