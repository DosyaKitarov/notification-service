package validator

import (
	"errors"
	"strconv"

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

func ParsePageNumParam(page string) uint32 {
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		return notificaitonService.DefaultPageNumParam
	}
	return uint32(pageNum)
}

func ParsePerPageParam(perPage string) uint32 {
	perPageNum, err := strconv.Atoi(perPage)
	if err != nil || perPageNum <= 0 || perPageNum > notificaitonService.PerPageLimit {
		return notificaitonService.DefaultPerPageParam
	}
	return uint32(perPageNum)
}
