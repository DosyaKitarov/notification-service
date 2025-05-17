package validator

import (
	"errors"

	"github.com/DosyaKitarov/notification-service/internal/notificaitonService"
)

type NotificationChannel string

func ValidateAuthNotificationRequest(req notificaitonService.AuthNotificationRequest) error {
	if req.UserID == 0 {
		return errors.New("UserID cannot be zero")
	}
	if req.NotificationChannel == notificaitonService.NotificationChannelUnknown {
		return errors.New("NotificationChannel cannot be empty")
	}
	if req.Email == "" {
		return errors.New("Email cannot be empty")
	}
	if req.Name == "" {
		return errors.New("Name cannot be empty")
	}
	return nil
}

func ValidateUserNotificationRequest(req notificaitonService.UserNotificationRequest) error {
	if req.UserID == 0 {
		return errors.New("UserID cannot be zero")
	}
	if len(req.Channels) == 0 {
		return errors.New("Channels cannot be empty")
	}
	if req.Name == "" {
		return errors.New("Name cannot be empty")
	}
	for _, channel := range req.Channels {
		if channel == notificaitonService.NotificationChannelUnknown {
			return errors.New("Channels cannot contain an empty NotificationChannel")
		}
	}
	if req.Email == "" {
		return errors.New("Email cannot be empty")
	}
	if req.Metadata == nil {
		return errors.New("Metadata cannot be nil")
	}
	if len(req.Metadata) == 0 {
		return errors.New("Metadata cannot be empty")
	}
	return nil
}
