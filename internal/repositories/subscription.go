package repositories

import (
	"context"
	"database/sql"
	"subscriptionsservice/internal/models"
	"time"

	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *models.Subscription) error
	GetByID(ctx context.Context, id int) (*models.Subscription, error)
	Update(ctx context.Context, sub *models.Subscription) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, filters map[string]interface{}) ([]*models.Subscription, error)
	CalculateTotalCost(ctx context.Context, periodStart, periodEnd time.Time, userID *uuid.UUID, serviceName *string) (int, error)
}

type PostgresSubscriptionRepository struct {
	db *sql.DB
}

func NewPostgresSubscriptionRepository(db *sql.DB) *PostgresSubscriptionRepository {
	return &PostgresSubscriptionRepository{db: db}
}

func (r *PostgresSubscriptionRepository) Create(ctx context.Context, sub *models.Subscription) error {
	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRowContext(ctx, query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).Scan(&sub.ID)
}

func (r *PostgresSubscriptionRepository) GetByID(ctx context.Context, id int) (*models.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date 
              FROM subscriptions WHERE id = $1`
	sub := &models.Subscription{}
	var endDate sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &endDate,
	)
	if err != nil {
		return nil, err
	}
	if endDate.Valid {
		sub.EndDate = &endDate.Time
	}
	return sub, nil
}

func (r *PostgresSubscriptionRepository) Update(ctx context.Context, sub *models.Subscription) error {
	query := `UPDATE subscriptions 
              SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5 
              WHERE id = $6`
	_, err := r.db.ExecContext(ctx, query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate, sub.ID)
	return err
}

func (r *PostgresSubscriptionRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresSubscriptionRepository) List(ctx context.Context, filters map[string]interface{}) ([]*models.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date 
              FROM subscriptions WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if userID, ok := filters["user_id"]; ok {
		query += ` AND user_id = $` + string(rune(argIndex))
		args = append(args, userID)
		argIndex++
	}
	if serviceName, ok := filters["service_name"]; ok {
		query += ` AND service_name = $` + string(rune(argIndex))
		args = append(args, serviceName)
		argIndex++
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*models.Subscription
	for rows.Next() {
		sub := &models.Subscription{}
		var endDate sql.NullTime
		if err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &endDate); err != nil {
			return nil, err
		}
		if endDate.Valid {
			sub.EndDate = &endDate.Time
		}
		subscriptions = append(subscriptions, sub)
	}
	return subscriptions, nil
}

func (r *PostgresSubscriptionRepository) CalculateTotalCost(ctx context.Context, periodStart, periodEnd time.Time, userID *uuid.UUID, serviceName *string) (int, error) {
	query := `
        SELECT COALESCE(SUM(price * months_active), 0)
        FROM (
            SELECT s.price,
                   COUNT(*) AS months_active
            FROM subscriptions s
            CROSS JOIN generate_series(
                date_trunc('month', s.start_date),
                COALESCE(date_trunc('month', s.end_date), 'infinity'),
                '1 month'
            ) AS m(month)
            WHERE m.month >= $1 AND m.month <= $2
              AND ($3::uuid IS NULL OR s.user_id = $3)
              AND ($4 IS NULL OR s.service_name = $4)
            GROUP BY s.id, s.price
        ) AS sub`
	var totalCost int
	err := r.db.QueryRowContext(ctx, query, periodStart, periodEnd, userID, serviceName).Scan(&totalCost)
	return totalCost, err
}
