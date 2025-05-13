package validator

import (
	"errors"

	"github.com/DosyaKitarov/notification-service/internal/service"
)

type NotificationChannel string

func ValidateAuthNotificationRequest(req service.AuthNotificationRequest) error {
	if req.UserID == 0 {
		return errors.New("UserID cannot be zero")
	}
	if req.NotificationChannel == service.NotificationChannelUnknown {
		return errors.New("NotificationChannel cannot be empty")
	}
	return nil
}

func ValidateUserNotificationRequest(req service.UserNotificationRequest) error {
	if req.UserID == 0 {
		return errors.New("UserID cannot be zero")
	}
	if len(req.Channels) == 0 {
		return errors.New("Channels cannot be empty")
	}
	for _, channel := range req.Channels {
		if channel == service.NotificationChannelUnknown {
			return errors.New("Channels cannot contain an empty NotificationChannel")
		}
	}
	return nil
}
