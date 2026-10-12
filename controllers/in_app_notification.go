package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/polar-bear-cu/sgt-noti-service/dtos"
	"github.com/polar-bear-cu/sgt-noti-service/models"
	"github.com/polar-bear-cu/sgt-noti-service/repositories"
	"github.com/polar-bear-cu/sgt-noti-service/usecases"
)

const userIDContextKey = "user_id"

type InAppNotificationController struct {
	uc *usecases.InAppNotificationUsecase
}

func NewInAppNotification(uc *usecases.InAppNotificationUsecase) *InAppNotificationController {
	return &InAppNotificationController{uc: uc}
}

func (ctl *InAppNotificationController) List(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}

	notifications, err := ctl.uc.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Error: "failed to list notifications"})
		return
	}

	response := make([]dtos.InAppNotificationResponse, 0, len(notifications))
	for _, notification := range notifications {
		response = append(response, toNotificationResponse(notification))
	}

	c.JSON(http.StatusOK, response)
}

func (ctl *InAppNotificationController) DeleteOne(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}

	err := ctl.uc.DeleteOne(c.Request.Context(), userID, c.Param("id"))
	switch {
	case err == nil:
		c.Status(http.StatusNoContent)
	case errors.Is(err, repositories.ErrInvalidNotificationID):
		c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Error: "invalid notification id"})
	case errors.Is(err, usecases.ErrNotificationNotFound):
		c.JSON(http.StatusNotFound, dtos.ErrorResponse{Error: "notification not found"})
	default:
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Error: "failed to delete notification"})
	}
}

func (ctl *InAppNotificationController) DeleteAll(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}

	deletedCount, err := ctl.uc.DeleteAll(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Error: "failed to clear notifications"})
		return
	}

	c.JSON(http.StatusOK, dtos.DeleteNotificationsResponse{DeletedCount: deletedCount})
}

func (ctl *InAppNotificationController) UnreadCount(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}

	count, err := ctl.uc.UnreadCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Error: "failed to count unread notifications"})
		return
	}

	c.JSON(http.StatusOK, dtos.UnreadCountResponse{UnreadCount: count})
}

func (ctl *InAppNotificationController) MarkAllRead(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}

	updatedCount, err := ctl.uc.MarkAllRead(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Error: "failed to mark notifications as read"})
		return
	}

	c.JSON(http.StatusOK, dtos.MarkReadResponse{UpdatedCount: updatedCount})
}

func authenticatedUserID(c *gin.Context) (string, bool) {
	value, exists := c.Get(userIDContextKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{Error: "unauthorized"})
		return "", false
	}

	userID, ok := value.(string)
	if !ok || userID == "" {
		c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{Error: "unauthorized"})
		return "", false
	}

	return userID, true
}

func toNotificationResponse(notification models.InAppNotification) dtos.InAppNotificationResponse {
	return dtos.InAppNotificationResponse{
		ID:             notification.ID.Hex(),
		UserID:         notification.UserID,
		SubscriptionID: notification.SubscriptionID,
		Title:          notification.Title,
		Content:        notification.Content,
		Status:         notification.Status,
		ReadStatus:     notification.ReadStatus,
		CreatedAt:      notification.CreatedAt,
		SentAt:         notification.SentAt,
	}
}
