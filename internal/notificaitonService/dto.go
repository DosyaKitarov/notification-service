package notificaitonService

import (
	pb "github.com/DosyaKitarov/notification-service/pkg/grpc"
)

type NotificationChannel string

const (
	NotificationChannelUnknown NotificationChannel = "unknown"
	NotificationChannelEmail   NotificationChannel = "email"
	NotificationChannelWeb     NotificationChannel = "web"
)

type NotificationType string

const (
	NotificationTypeUnknown           NotificationType = "unknown"
	NotificationTypeRegistration      NotificationType = "registration"
	NotificationTypeLogin             NotificationType = "login"
	NotificationTypeInvestmentSuccess NotificationType = "investment_success"
	NotificationTypeInvestedInYou     NotificationType = "invested_in_you"
	NotificationTypeOther             NotificationType = "other"
)

const (
	EmailTemplateRegistration = "registration"
	EmailTemplateLogin        = "login"
)

const (
	DefaultPerPageParam = 10
	DefaultPageNumParam = 1
	PerPageLimit        = 1000
)

type AuthNotificationRequest struct {
	UserID              uint64
	Email               string
	Name                string
	NotificationChannel NotificationChannel
}

type AuthNotificationRequestDTO struct {
	UserID              uint64
	Email               string
	Name                string
	NotificationChannel string
	NotificationType    string
}

func (auth *AuthNotificationRequest) ToDTO(NotificationType string) AuthNotificationRequestDTO {
	return AuthNotificationRequestDTO{
		UserID:              auth.UserID,
		Email:               auth.Email,
		Name:                auth.Name,
		NotificationChannel: nCtoString(auth.NotificationChannel),
		NotificationType:    NotificationType,
	}
}

func (auth *AuthNotificationRequestDTO) ToModel() Notification {
	return Notification{
		UserID:              auth.UserID,
		NotificationChannel: []string{auth.NotificationChannel},
		Name:                auth.Name,
		NotificationType:    auth.NotificationType,
		Email:               auth.Email,
	}
}

type UserNotificationRequest struct {
	UserID   uint64
	Email    string
	Name     string
	Type     NotificationType
	Channels []NotificationChannel
	Metadata map[string]string
}

type UserNotificationRequestDTO struct {
	UserID           uint64
	Email            string
	Name             string
	NotificationType string
	Channels         []string
	Metadata         map[string]string
}

func (user *UserNotificationRequest) ToDTO() UserNotificationRequestDTO {
	return UserNotificationRequestDTO{
		UserID:           user.UserID,
		Email:            user.Email,
		Name:             user.Name,
		Channels:         nCtoStringSlice(user.Channels),
		Metadata:         user.Metadata,
		NotificationType: nTtoString(user.Type),
	}
}

func (user *UserNotificationRequestDTO) ToModel() Notification {
	return Notification{
		UserID:              user.UserID,
		NotificationChannel: user.Channels,
		Name:                user.Name,
		Metadata:            user.Metadata,
		Email:               user.Email,
		NotificationType:    user.NotificationType,
	}
}

func (user *UserNotificationRequestDTO) toWebNotification() WebNotification {
	return WebNotification{
		ID:      user.UserID,
		Subject: parseSubjectFromMetadata(user.Metadata),
		Body:    parseBodyFromMetadata(user.Metadata),
	}
}

type Notification struct {
	ID                  uint64            `db:"id"`
	UserID              uint64            `db:"user_id"`
	Email               string            `db:"email"`
	Name                string            `db:"name"`
	NotificationType    string            `db:"type"`
	NotificationChannel []string          `db:"notification_channel"`
	Metadata            map[string]string `db:"metadata"`
}

type EmailNotificationRespone struct {
	UserID           uint64            `db:"user_id"`
	Email            string            `db:"email"`
	Name             string            `db:"name"`
	NotificationType string            `db:"type"`
	Metadata         map[string]string `db:"metadata"`
	CreatedAt        string            `db:"created_at"`
}

type WebNotificationResponse struct {
	UserID           uint64            `db:"user_id"`
	Email            string            `db:"email"`
	Name             string            `db:"name"`
	NotificationType string            `db:"type"`
	Metadata         map[string]string `db:"metadata"`
	CreatedAt        string            `db:"created_at"`
	IsRead           bool              `db:"is_read"`
}

type WebNotification struct {
	ID      uint64 `json:"id"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (n *Notification) toWebNotification() WebNotification {
	return WebNotification{
		ID:      n.ID,
		Subject: parseSubjectFromMetadata(n.Metadata),
		Body:    parseBodyFromMetadata(n.Metadata),
	}
}

type GetNotificationsRequest struct {
	PerPage uint32
	Page    uint32
}

type GetEmailNotifications struct {
	Notifications []EmailNotificationRespone
	Total         uint32
}
type GetWebNotifications struct {
	Notifications []WebNotificationResponse
	Total         uint32
}

func (g *GetNotificationsRequest) getPage() uint32 {
	if g.Page == 0 {
		g.Page = DefaultPageNumParam
	}
	return (g.Page - 1) * g.getPerPage()
}

func (g *GetNotificationsRequest) getPerPage() uint32 {
	if g.PerPage == 0 {
		return DefaultPerPageParam
	}
	if g.PerPage > PerPageLimit {
		return PerPageLimit
	}
	return g.PerPage
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

func ToNotificationChannels(channels []pb.NotificationChannel) []NotificationChannel {
	var result []NotificationChannel
	for _, channel := range channels {
		result = append(result, ToNotificationChannel(channel))
	}
	return result
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

func ToNotificationType(notificationType pb.NotificationType) NotificationType {
	switch notificationType {
	case pb.NotificationType_REGISTRATION:
		return NotificationTypeRegistration
	case pb.NotificationType_LOGIN:
		return NotificationTypeLogin
	case pb.NotificationType_INVESTMENT_SUCCESS:
		return NotificationTypeInvestmentSuccess
	case pb.NotificationType_INVESTED_IN_YOU:
		return NotificationTypeInvestedInYou
	case pb.NotificationType_OTHER:
		return NotificationTypeOther
	default:
		return NotificationTypeUnknown
	}
}
func nTtoString(notificationType NotificationType) string {
	switch notificationType {
	case NotificationTypeRegistration:
		return "registration"
	case NotificationTypeLogin:
		return "login"
	case NotificationTypeInvestmentSuccess:
		return "investment_success"
	case NotificationTypeInvestedInYou:
		return "invested_in_you"
	case NotificationTypeOther:
		return "other"
	default:
		return "unknown"
	}
}

func parseSubjectFromMetadata(metadata map[string]string) string {
	subject, ok := metadata["subject"]
	if !ok {
		return "No subject"
	}
	return subject
}
func parseBodyFromMetadata(metadata map[string]string) string {
	body, ok := metadata["body"]
	if !ok {
		return "No body"
	}
	return body
}
