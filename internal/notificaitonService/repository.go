package notificaitonService

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

func (r *Repository) GetEmailNotifications(ctx context.Context, request GetNotificationsRequest) (GetEmailNotifications, error) {
	var (
		result = make([]EmailNotificationRespone, 0, request.getPerPage())
		total  sql.NullInt32

		query = `
        WITH data_CTE AS (
            SELECT 
                user_id,
                email,
                name,
                type,
                metadata,
                created_at
            FROM email_notifications
            ORDER BY created_at DESC
            LIMIT $1
            OFFSET $2
        ),
        total AS (
            SELECT count(1) as total FROM email_notifications
        )
        SELECT 
            t.total,
            d.user_id,
            d.email,
            d.name,
            d.type,
            d.metadata,
            d.created_at
        FROM total t CROSS JOIN data_CTE d;
        `
	)

	rows, err := r.db.QueryContext(
		ctx,
		query,
		request.getPerPage(),
		request.getPage(),
	)
	if err != nil {
		r.logger.Error("Failed to get email notifications from database", zap.Error(err))
		return GetEmailNotifications{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var n EmailNotificationRespone
		var metadataBytes []byte
		if err := rows.Scan(
			&total,
			&n.UserID,
			&n.Email,
			&n.Name,
			&n.NotificationType,
			&metadataBytes,
			&n.CreatedAt,
		); err != nil {
			r.logger.Error("Failed to scan email notification", zap.Error(err))
			return GetEmailNotifications{}, err
		}
		if err := json.Unmarshal(metadataBytes, &n.Metadata); err != nil {
			r.logger.Error("Failed to unmarshal metadata", zap.Error(err))
			return GetEmailNotifications{}, err
		}
		result = append(result, n)
	}
	if rows.Err() != nil {
		return GetEmailNotifications{}, rows.Err()
	}

	return GetEmailNotifications{
		Notifications: result,
		Total:         uint32(total.Int32),
	}, nil
}

func (r *Repository) GetWebNotifications(ctx context.Context, request GetNotificationsRequest) (GetWebNotifications, error) {
	var (
		result = make([]WebNotificationResponse, 0, request.getPerPage())
		total  sql.NullInt32

		query = `
        WITH data_CTE AS (
            SELECT 
                user_id,
                email,
                name,
                type,
                metadata,
				is_read,
                created_at
            FROM web_notifications
            ORDER BY created_at DESC
            LIMIT $1
            OFFSET $2
        ),
        total AS (
            SELECT count(1) as total FROM web_notifications
        )
        SELECT 
            t.total,
            d.user_id,
            d.email,
            d.name,
            d.type,
            d.metadata,
			d.is_read,
            d.created_at
        FROM total t CROSS JOIN data_CTE d;
        `
	)

	rows, err := r.db.QueryContext(
		ctx,
		query,
		request.getPerPage(),
		request.getPage(),
	)
	if err != nil {
		r.logger.Error("Failed to get email notifications from database", zap.Error(err))
		return GetWebNotifications{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var n WebNotificationResponse
		var metadataBytes []byte
		if err := rows.Scan(
			&total,
			&n.UserID,
			&n.Email,
			&n.Name,
			&n.NotificationType,
			&metadataBytes,
			&n.IsRead,
			&n.CreatedAt,
		); err != nil {
			r.logger.Error("Failed to scan email notification", zap.Error(err))
			return GetWebNotifications{}, err
		}
		if err := json.Unmarshal(metadataBytes, &n.Metadata); err != nil {
			r.logger.Error("Failed to unmarshal metadata", zap.Error(err))
			return GetWebNotifications{}, err
		}
		result = append(result, n)
	}
	if rows.Err() != nil {
		return GetWebNotifications{}, rows.Err()
	}

	return GetWebNotifications{
		Notifications: result,
		Total:         uint32(total.Int32),
	}, nil
}
