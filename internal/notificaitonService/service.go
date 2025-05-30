package notificaitonService

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/DosyaKitarov/notification-service/pkg/email"
	"go.uber.org/zap"
)

type NotificationRepository interface {
	SaveEmailNotificationWithTx(ctx context.Context, tx *sql.Tx, n Notification) error
	SaveWebNotificationWithTx(ctx context.Context, tx *sql.Tx, n Notification) error
	GetUnreadWebNotifications(ctx context.Context, userID uint64) ([]Notification, error)
	BeginTransaction(ctx context.Context) (*sql.Tx, error)
	MarkNotificationAsRead(ctx context.Context, tx *sql.Tx, notificationID uint64) error
	GetEmailNotifications(ctx context.Context, request GetNotificationsRequest) (GetEmailNotifications, error)
	GetWebNotifications(ctx context.Context, request GetNotificationsRequest) (GetWebNotifications, error)
}

type WebNotifier interface {
	SendToUser(userID uint64, data []byte) error
}

type NotificationService struct {
	repo        NotificationRepository
	emailSender email.EmailSender
	webNotifier WebNotifier
	logger      *zap.Logger
}

func NewNotificationService(repo NotificationRepository, emailSender email.EmailSender, logger *zap.Logger) *NotificationService {
	return &NotificationService{
		repo:        repo,
		emailSender: emailSender,
		logger:      logger,
	}
}

func (s *NotificationService) SetWebNotifier(webNotifier WebNotifier) {
	s.webNotifier = webNotifier
}

func (s *NotificationService) RegistrationNotification(ctx context.Context, notification AuthNotificationRequestDTO) error {
	s.logger.Info("Starting RegistrationNotification process")

	tx, err := s.repo.BeginTransaction(ctx)
	if err != nil {
		s.logger.Error("Error starting transaction", zap.Error(err))
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				s.logger.Error("Error committing transaction", zap.Error(commitErr))
				err = commitErr
			}
		}
	}()

	err = s.repo.SaveEmailNotificationWithTx(ctx, tx, notification.ToModel())
	if err != nil {
		s.logger.Error("Error saving notification", zap.Error(err))
		return err
	}
	s.logger.Info("Notification saved successfully")

	err = s.emailSender.SendAuthEmail(EmailTemplateRegistration, notification.Email, notification.Name)
	if err != nil {
		s.logger.Error("Error sending email", zap.Error(err))
		return err
	}
	s.logger.Info("Email sent successfully")

	return nil
}

func (s *NotificationService) LoginNotification(ctx context.Context, notification AuthNotificationRequestDTO) error {
	s.logger.Info("Starting LoginNotification process")

	tx, err := s.repo.BeginTransaction(ctx)
	if err != nil {
		s.logger.Error("Error starting transaction", zap.Error(err))
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				s.logger.Error("Error committing transaction", zap.Error(commitErr))
				err = commitErr
			}
		}
	}()

	err = s.repo.SaveEmailNotificationWithTx(ctx, tx, notification.ToModel())
	if err != nil {
		s.logger.Error("Error saving notification", zap.Error(err))
		return err
	}
	s.logger.Info("Notification saved successfully")

	err = s.emailSender.SendAuthEmail(EmailTemplateLogin, notification.Email, notification.Name)
	if err != nil {
		s.logger.Error("Error sending email", zap.Error(err))
		return err
	}

	return nil
}

func (s *NotificationService) UserNotification(ctx context.Context, notification UserNotificationRequestDTO) error {
	s.logger.Info("Starting UserNotification process")

	tx, err := s.repo.BeginTransaction(ctx)
	if err != nil {
		s.logger.Error("Error starting transaction", zap.Error(err))
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				s.logger.Error("Error committing transaction", zap.Error(commitErr))
				err = commitErr
			}
		}
	}()

	for _, channel := range notification.Channels {
		if channel == nCtoString(NotificationChannelEmail) {
			err = s.emailNotification(ctx, tx, s.emailSender, notification)
			if err != nil {
				s.logger.Error("Error sending email notification", zap.Error(err))
				return err
			}
		} else if channel == nCtoString(NotificationChannelWeb) {
			err = s.webNotification(ctx, tx, notification)
			if err != nil {
				s.logger.Error("Error sending web notification", zap.Error(err))
				return err
			}
		} else {
			s.logger.Error("Unknown notification channel", zap.String("channel", channel))
			return nil
		}
	}
	return nil
}

func (s *NotificationService) emailNotification(ctx context.Context, tx *sql.Tx, emailSender email.EmailSender, notification UserNotificationRequestDTO) error {
	s.logger.Info("Starting emailNotification process")

	err := s.repo.SaveEmailNotificationWithTx(ctx, tx, notification.ToModel())
	if err != nil {
		s.logger.Error("Error saving notification", zap.Error(err))
		return err
	}
	s.logger.Info("Notification saved successfully")

	err = emailSender.SendUserEmail(notification.Email, notification.Metadata)
	if err != nil {
		s.logger.Error("Error sending email", zap.Error(err))
		return err
	}
	s.logger.Info("Email sent successfully")
	return nil
}

func (s *NotificationService) webNotification(ctx context.Context, tx *sql.Tx, notification UserNotificationRequestDTO) error {
	s.logger.Info("Starting webNotification process")

	err := s.repo.SaveWebNotificationWithTx(ctx, tx, notification.ToModel())
	if err != nil {
		s.logger.Error("Error saving notification", zap.Error(err))
		return err
	}
	s.logger.Info("Notification saved successfully")

	// Отправка по websocket
	if s.webNotifier != nil {
		webNotif := notification.toWebNotification()
		data, _ := json.Marshal(webNotif)
		sendErr := s.webNotifier.SendToUser(notification.UserID, data)
		if sendErr != nil {
			s.logger.Error("Error sending websocket notification", zap.Error(sendErr))
			// не возвращаем ошибку, чтобы не откатывать транзакцию
		}
	}

	return nil
}
func (s *NotificationService) GetUnreadWebNotifications(ctx context.Context, userID uint64) ([]WebNotification, error) {
	s.logger.Info("Starting GetUnreadWebNotifications process")

	notifications, err := s.repo.GetUnreadWebNotifications(ctx, userID)
	if err != nil {
		s.logger.Error("Error getting unread web notifications", zap.Error(err))
		return nil, err
	}
	s.logger.Info("Unread web notifications retrieved successfully")

	var webNotifications []WebNotification
	for _, n := range notifications {
		webNotifications = append(webNotifications, n.toWebNotification())
	}

	return webNotifications, nil
}

func (s *NotificationService) MarkNotificationAsRead(ctx context.Context, notificationID uint64) error {
	s.logger.Info("Starting MarkNotificationAsRead process")

	tx, err := s.repo.BeginTransaction(ctx)
	if err != nil {
		s.logger.Error("Error starting transaction", zap.Error(err))
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				s.logger.Error("Error committing transaction", zap.Error(commitErr))
				err = commitErr
			}
		}
	}()

	err = s.repo.MarkNotificationAsRead(ctx, tx, notificationID)
	if err != nil {
		s.logger.Error("Error marking notification as read", zap.Error(err))
		return err
	}
	s.logger.Info("Notification marked as read successfully")

	return nil
}

func (s *NotificationService) GetEmailNotifications(ctx context.Context, request GetNotificationsRequest) (GetEmailNotifications, error) {
	s.logger.Info("Starting GetEmailNotifications process")

	response, err := s.repo.GetEmailNotifications(ctx, request)
	if err != nil {
		s.logger.Error("GetEmailNotifications response error", zap.Error(err))
		return GetEmailNotifications{}, err
	}

	return response, nil
}

func (s *NotificationService) GetWebNotifications(ctx context.Context, request GetNotificationsRequest) (GetWebNotifications, error) {
	s.logger.Info("Starting GetWebNotifications process")

	response, err := s.repo.GetWebNotifications(ctx, request)
	if err != nil {
		s.logger.Error("GetWebNotifications response error", zap.Error(err))
		return GetWebNotifications{}, err
	}

	return response, nil
}
