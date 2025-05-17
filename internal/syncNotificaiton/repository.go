package syncNotificaiton

import (
	"context"
	"database/sql"
	"encoding/json"

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

func (r *Repository) SaveEmailNotificationWithTx(ctx context.Context, tx *sql.Tx, n Notification) error {
	metadataJSON, err := json.Marshal(n.Metadata)
	if err != nil {
		return err
	}

	query := `
        INSERT INTO email_notifications (user_id,email, name, type, metadata)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err = tx.ExecContext(ctx, query, n.UserID, n.Email, n.Name, n.NotificationType, metadataJSON)
	if err != nil {
		r.logger.Error("Failed to save notification to database", zap.Error(err))
		return err
	}
	r.logger.Info("Notification saved to database", zap.Uint64("user_id", n.UserID))
	return nil
}

func (r *Repository) SaveWebNotificationWithTx(ctx context.Context, tx *sql.Tx, n Notification) error {
	metadataJSON, err := json.Marshal(n.Metadata)
	if err != nil {
		return err
	}

	query := `
        INSERT INTO web_notifications (user_id,email, name, type, metadata)
        VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.ExecContext(ctx, query, n.UserID, n.Email, n.Name, n.NotificationType, metadataJSON)
	if err != nil {
		r.logger.Error("Failed to save notification to database", zap.Error(err))
		return err
	}
	r.logger.Info("Notification saved to database", zap.Uint64("user_id", n.UserID))
	return nil
}

func (r *Repository) GetUnreadWebNotifications(ctx context.Context, userID uint64) ([]Notification, error) {
	query := `
		SELECT id, user_id,email, name, type, metadata
		FROM web_notifications
		WHERE user_id = $1 AND is_read = false
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		r.logger.Error("Failed to get unread notifications from database", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var n Notification
		var metadataBytes []byte
		if err := rows.Scan(&n.ID, &n.UserID, &n.Email, &n.Name, &n.NotificationType, &metadataBytes); err != nil {
			r.logger.Error("Failed to scan notification", zap.Error(err))
			return nil, err
		}
		if err := json.Unmarshal(metadataBytes, &n.Metadata); err != nil {
			r.logger.Error("Failed to unmarshal metadata", zap.Error(err))
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

func (r *Repository) MarkNotificationAsRead(ctx context.Context, tx *sql.Tx, notificationID uint64) error {
	query := `
		UPDATE web_notifications
		SET is_read = true
		WHERE id = $1
	`
	_, err := tx.ExecContext(ctx, query, notificationID)
	if err != nil {
		r.logger.Error("Failed to mark notification as read", zap.Error(err))
		return err
	}
	return nil
}
