package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/DosyaKitarov/notification-service/internal/service"
	"go.uber.org/zap"
)

type Repository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewRepository(db *sql.DB, logger *zap.Logger) *Repository {
	return &Repository{db: db, logger: logger}
}

func (r *Repository) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *Repository) SaveNotificationWithTx(ctx context.Context, tx *sql.Tx, n service.Notification) error {
	metadataJSON, err := json.Marshal(n.Metadata)
	if err != nil {
		return err
	}

	notificationChannelJSON, err := json.Marshal(n.NotificationChannel)
	if err != nil {
		return err
	}

	query := `
        INSERT INTO notifications (user_id,email, name, type, notification_channel, metadata)
        VALUES ($1, $2, $3, $4, $5, $6	)
    `
	_, err = tx.ExecContext(ctx, query, n.UserID, n.Email, n.Name, n.NotificationType, notificationChannelJSON, metadataJSON)
	if err != nil {
		r.logger.Error("Failed to save notification to database", zap.Error(err))
		return err
	}
	r.logger.Info("Notification saved to database", zap.Uint64("user_id", n.UserID))
	return nil
}
