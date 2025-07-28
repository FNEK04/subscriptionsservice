package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          int
	ServiceName string
	Price       int
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
}

type TotalCostRequest struct {
	PeriodStart string  `json:"period_start"`
	PeriodEnd   string  `json:"period_end"`
	UserID      *string `json:"user_id,omitempty"`
	ServiceName *string `json:"service_name,omitempty"`
}
