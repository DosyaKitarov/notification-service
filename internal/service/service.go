package service

import (
	"context"
	"database/sql"
	"log"

	"github.com/DosyaKitarov/notification-service/pkg/email"
)

type NotificationRepository interface {
	SaveNotificationWithTx(ctx context.Context, tx *sql.Tx, n Notification) error
	BeginTransaction(ctx context.Context) (*sql.Tx, error)
}

type NotificationService struct {
	repo        NotificationRepository
	emailSender email.EmailSender
}

func NewNotificationService(repo NotificationRepository, emailSender email.EmailSender) *NotificationService {
	return &NotificationService{
		repo:        repo,
		emailSender: emailSender,
	}
}

func (s *NotificationService) RegistrationNotification(ctx context.Context, notification AuthNotificationRequestDTO) error {
	log.Println("Starting RegistrationNotification process")

	// Начало транзакции
	tx, err := s.repo.BeginTransaction(ctx)
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
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
				log.Printf("Error committing transaction: %v", commitErr)
				err = commitErr
			}
		}
	}()

	// Сохранение уведомления
	err = s.repo.SaveNotificationWithTx(ctx, tx, notification.ToModel())
	if err != nil {
		log.Printf("Error saving notification: %v", err)
		return err
	}
	log.Println("Notification saved successfully")

	// Отправка email (включена в транзакцию)
	err = s.emailSender.SendEmail(notification.Metadata, notification.Email)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}
	log.Println("Email sent successfully")

	return nil
}
