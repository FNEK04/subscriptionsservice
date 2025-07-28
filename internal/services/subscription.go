package services

import (
	"context"
	"subscriptionsservice/internal/models"
	"subscriptionsservice/internal/repositories"
	"time"

	"github.com/google/uuid"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, sub *models.Subscription) error
	GetSubscription(ctx context.Context, id int) (*models.Subscription, error)
	UpdateSubscription(ctx context.Context, sub *models.Subscription) error
	DeleteSubscription(ctx context.Context, id int) error
	ListSubscriptions(ctx context.Context, filters map[string]interface{}) ([]*models.Subscription, error)
	CalculateTotalCost(ctx context.Context, periodStart, periodEnd time.Time, userID *uuid.UUID, serviceName *string) (int, error)
}

type SubscriptionServiceImpl struct {
	repo repositories.SubscriptionRepository
}

func NewSubscriptionService(repo repositories.SubscriptionRepository) *SubscriptionServiceImpl {
	return &SubscriptionServiceImpl{repo: repo}
}

func (s *SubscriptionServiceImpl) CreateSubscription(ctx context.Context, sub *models.Subscription) error {
	return s.repo.Create(ctx, sub)
}

func (s *SubscriptionServiceImpl) GetSubscription(ctx context.Context, id int) (*models.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SubscriptionServiceImpl) UpdateSubscription(ctx context.Context, sub *models.Subscription) error {
	return s.repo.Update(ctx, sub)
}

func (s *SubscriptionServiceImpl) DeleteSubscription(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *SubscriptionServiceImpl) ListSubscriptions(ctx context.Context, filters map[string]interface{}) ([]*models.Subscription, error) {
	return s.repo.List(ctx, filters)
}

func (s *SubscriptionServiceImpl) CalculateTotalCost(ctx context.Context, periodStart, periodEnd time.Time, userID *uuid.UUID, serviceName *string) (int, error) {
	return s.repo.CalculateTotalCost(ctx, periodStart, periodEnd, userID, serviceName)
}
