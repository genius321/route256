package postgres

import (
	"context"
	"fmt"
	"route256/notifications/internal/pkg/notifications"
	"route256/notifications/internal/repository/postgres/tx"
	"time"

	"github.com/georgysavva/scany/pgxscan"
)

type Repository struct {
	provider tx.DBProvider
}

func New(provider tx.DBProvider) *Repository {
	return &Repository{provider: provider}
}

func (r *Repository) GetHistory(ctx context.Context, req *notifications.GetHistoryRequest) (*notifications.GetHistoryResponse, error) {
	db := r.provider.GetDB(ctx)
	query := `
	SELECT order_id, status_name, created_at
	FROM notifications
	WHERE user_id=$1
		AND created_at >= $2
		AND created_at <= $3;
	`

	var result []struct {
		OrderId int64     `db:"order_id"`
		Status  string    `db:"status_name"`
		Time    time.Time `db:"created_at"`
	}

	err := pgxscan.Select(ctx, db, &result, query, req.UserId, req.StartTime, req.EndTime)
	if err != nil {
		return nil, fmt.Errorf("exec select GetHistory: %w", err)
	}

	history := make([]*notifications.Entry, len(result))
	for i, v := range result {
		history[i] = &notifications.Entry{
			OrderId:   v.OrderId,
			Status:    v.Status,
			CreatedAt: v.Time.Format("2006-01-02 15:04:05.000000-07"),
		}
	}

	return &notifications.GetHistoryResponse{Entries: history}, nil
}

func (r *Repository) SaveNotification(ctx context.Context, orderId, userId int64, status string) error {
	db := r.provider.GetDB(ctx)
	query := `
	INSERT INTO notifications (order_id, user_id, status_name)
	VALUES ($1, $2, $3)
	`
	_, err := db.Exec(ctx, query, orderId, userId, status)
	if err != nil {
		return fmt.Errorf("exec insert SaveNotification: %w", err)
	}

	return nil
}
