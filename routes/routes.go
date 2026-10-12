package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/polar-bear-cu/sgt-noti-service/controllers"
	"github.com/polar-bear-cu/sgt-noti-service/middlewares"
)

func Register(r *gin.Engine, inAppNotification *controllers.InAppNotificationController, jwtSecret string) {
	r.GET("/health", controllers.Health)

	notifications := r.Group("/notifications")
	notifications.Use(middlewares.RequireAuth(jwtSecret))
	notifications.GET("", inAppNotification.List)
	notifications.DELETE("/:id", inAppNotification.DeleteOne)
	notifications.DELETE("", inAppNotification.DeleteAll)
}
