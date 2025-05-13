package service

import (
	pb "github.com/DosyaKitarov/notification-service/pkg/grpc"
)

type NotificationChannel int

const (
	NotificationChannelUnknown NotificationChannel = iota
	NotificationChannelEmail
	NotificationChannelWeb
)

type NotificationType int

const (
	NotificationTypeUnknown NotificationType = iota
	NotificationTypeRegistrationConfirmation
	NotificationTypeLoginAlert
	NotificationTypeInvestmentSuccess
	NotificationTypeInvestedInYou
)

type AuthNotificationRequest struct {
	UserID              uint64
	NotificationChannel NotificationChannel
	Metadata            map[string]string
}

type AuthNotificationRequestDTO struct {
	UserID              uint64
	NotificationChannel string
	Metadata            map[string]string
	NotificationType    string
}

func (auth *AuthNotificationRequest) ToDTO(NotificationType string) AuthNotificationRequestDTO {
	return AuthNotificationRequestDTO{
		UserID:              auth.UserID,
		NotificationChannel: nCtoString(auth.NotificationChannel),
		Metadata:            auth.Metadata,
		NotificationType:    NotificationType,
	}
}

func (auth *AuthNotificationRequestDTO) ToModel() Notification {
	return Notification{
		UserID:              auth.UserID,
		NotificationChannel: []string{auth.NotificationChannel},
		Metadata:            auth.Metadata,
		NotificationType:    auth.NotificationType,
	}
}

type UserNotificationRequest struct {
	UserID   uint64
	Channels []NotificationChannel
	Metadata map[string]string
}

type UserNotificationRequestDTO struct {
	UserID           uint64
	Channels         []string
	Metadata         map[string]string
	NotificationType []string
}

func (user *UserNotificationRequest) ToDTO(NotificationType []string) *UserNotificationRequestDTO {
	return &UserNotificationRequestDTO{
		UserID:           user.UserID,
		Channels:         nCtoStringSlice(user.Channels),
		Metadata:         user.Metadata,
		NotificationType: NotificationType,
	}
}

func (user *UserNotificationRequestDTO) ToModel() *Notification {
	return &Notification{
		UserID:              user.UserID,
		NotificationChannel: user.Channels,
		Metadata:            user.Metadata,
	}
}

type Notification struct {
	UserID              uint64            `db:"user_id"`
	NotificationType    string            `db:"notification_type"`
	NotificationChannel []string          `db:"notification_channel"`
	Metadata            map[string]string `db:"metadata"`
}

func ToNotificationChannel(channel pb.NotificationChannel) NotificationChannel {
	switch channel {
	case pb.NotificationChannel_EMAIL:
		return NotificationChannelEmail
	case pb.NotificationChannel_WEB:
		return NotificationChannelWeb
	default:
		return NotificationChannelUnknown
	}
}

func nCtoString(channel NotificationChannel) string {
	switch channel {
	case NotificationChannelEmail:
		return "email"
	case NotificationChannelWeb:
		return "web"
	default:
		return "unknown"
	}
}

func nCtoStringSlice(channels []NotificationChannel) []string {
	var result []string
	for _, channel := range channels {
		result = append(result, nCtoString(channel))
	}
	return result
}
