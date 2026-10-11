package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/polar-bear-cu/sgt-noti-service/controllers"
)

func Register(r *gin.Engine, inAppNotification *controllers.InAppNotificationController) {
	r.GET("/health", controllers.Health)

	r.GET("/notifications", inAppNotification.List)
	r.DELETE("/notifications/:id", inAppNotification.DeleteOne)
	r.DELETE("/notifications", inAppNotification.DeleteAll)
}
