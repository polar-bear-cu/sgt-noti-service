package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/polar-bear-cu/sgt-noti-service/routes"
)

func main() {
	r := gin.Default()
	routes.Register(r)

	log.Println("listening :8088")
	if err := r.Run(":8088"); err != nil {
		log.Fatal(err)
	}
}
