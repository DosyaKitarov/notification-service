package service

import (
	"context"
	"log"
)

type NotificationRepository interface {
	SaveNotification(ctx context.Context, n Notification) error
}

type NotificationService struct {
	repo NotificationRepository
}

func NewNotificationService(repo NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) RegistrationNotification(ctx context.Context, notification AuthNotificationRequestDTO) error {
	log.Println("Starting RegistrationNotification process")
	err := s.repo.SaveNotification(ctx, notification.ToModel())
	if err != nil {
		log.Printf("Error saving notification: %v", err)
		return err
	}
	log.Println("Notification saved successfully")
	// Add logic to handle registration notification
	// Send email/front
	return nil
}
