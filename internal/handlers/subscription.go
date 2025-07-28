package handlers

import (
	"net/http"
	"strconv"
	"subscriptionsservice/internal/models"
	"subscriptionsservice/internal/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	subscriptionService services.SubscriptionService
}

func NewHandler(subscriptionService services.SubscriptionService) *Handler {
	return &Handler{subscriptionService: subscriptionService}
}

func parseMonthYear(s string) (time.Time, error) {
	return time.Parse("01-2006", "01-"+s)
}

// CreateSubscription godoc
// @Summary Создать новую подписку
// @Description Создает новую подписку с указанными данными
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body models.Subscription true "Данные подписки"
// @Success 201 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Router /subscriptions [post]
func (h *Handler) CreateSubscription(c *gin.Context) {
	var sub models.Subscription
	if err := c.BindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := h.subscriptionService.CreateSubscription(c.Request.Context(), &sub); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sub)
}

// GetSubscription godoc
// @Summary Получить подписку по ID
// @Description Возвращает подписку по указанному ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "ID подписки"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *Handler) GetSubscription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		return
	}
	sub, err := h.subscriptionService.GetSubscription(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{"error": "Subscription not found"})
		return
	}
	c.JSON(http.StatusOK, sub)
}

// UpdateSubscription godoc
// @Summary Обновить подписку
// @Description Обновляет данные подписки по ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "ID подписки"
// @Param subscription body models.Subscription true "Данные подписки"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (h *Handler) UpdateSubscription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		return
	}
	var sub models.Subscription
	if err := c.BindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	sub.ID = id
	if err := h.subscriptionService.UpdateSubscription(c.Request.Context(), &sub); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sub)
}

// DeleteSubscription godoc
// @Summary Удалить подписку
// @Description Удаляет подписку по указанному ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "ID подписки"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *Handler) DeleteSubscription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		return
	}
	if err := h.subscriptionService.DeleteSubscription(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ListSubscriptions godoc
// @Summary Список подписок
// @Description Возвращает список подписок с фильтрацией по user_id и service_name
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "ID пользователя"
// @Param service_name query string false "Название сервиса"
// @Success 200 {array} models.Subscription
// @Failure 400 {object} map[string]string
// @Router /subscriptions [get]
func (h *Handler) ListSubscriptions(c *gin.Context) {
	filters := make(map[string]interface{})
	if userID := c.Query("user_id"); userID != "" {
		filters["user_id"] = userID
	}
	if serviceName := c.Query("service_name"); serviceName != "" {
		filters["service_name"] = serviceName
	}
	subscriptions, err := h.subscriptionService.ListSubscriptions(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subscriptions)
}

// CalculateTotalCost godoc
// @Summary Рассчитать общую стоимость подписок
// @Description Рассчитывает общую стоимость подписок за период с фильтрами
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body models.TotalCostRequest true "Запрос на расчет"
// @Success 200 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Router /subscriptions/total-cost [post]
func (h *Handler) CalculateTotalCost(c *gin.Context) {
	var req models.TotalCostRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	periodStart, err := parseMonthYear(req.PeriodStart)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid period_start"})
		return
	}
	periodEnd, err := parseMonthYear(req.PeriodEnd)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid period_end"})
		return
	}
	var userID *uuid.UUID
	if req.UserID != nil {
		uid, err := uuid.Parse(*req.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user_id"})
			return
		}
		userID = &uid
	}
	totalCost, err := h.subscriptionService.CalculateTotalCost(c.Request.Context(), periodStart, periodEnd, userID, req.ServiceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]int{"total_cost": totalCost})
}
