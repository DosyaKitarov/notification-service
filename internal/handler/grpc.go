package handler

import (
	"context"
	"database/sql"
	"fmt"

	service "github.com/DosyaKitarov/notification-service/internal/notificaitonService"
	pb "github.com/DosyaKitarov/notification-service/pkg/grpc"
	"github.com/DosyaKitarov/notification-service/pkg/validator"
	"go.uber.org/zap"
)

type GrpcNotificationServiceHandler struct {
	NotificationService
	pb.UnimplementedNotificationServiceServer
	db     *sql.DB
	logger *zap.Logger
}

type NotificationService interface {
	RegistrationNotification(ctx context.Context, notification service.AuthNotificationRequestDTO) error
	LoginNotification(ctx context.Context, notification service.AuthNotificationRequestDTO) error
	UserNotification(ctx context.Context, notification service.UserNotificationRequestDTO) error
	GetEmailNotifications(ctx context.Context, request service.GetNotificationsRequest) (service.GetEmailNotifications, error)
	GetWebNotifications(ctx context.Context, request service.GetNotificationsRequest) (service.GetWebNotifications, error)
}

func NewNotificationServiceHandler(db *sql.DB, syncService NotificationService, logger *zap.Logger) *GrpcNotificationServiceHandler {
	return &GrpcNotificationServiceHandler{
		db:                  db,
		NotificationService: syncService,
		logger:              logger,
	}
}

func (h *GrpcNotificationServiceHandler) SendRegistrationNotification(ctx context.Context, req *pb.AuthNotificationRequest) (*pb.SendNotificationResponse, error) {
	h.logger.Info("Received SendRegistrationNotification request", zap.Any("request", req))

	var (
		userID  = req.GetUserId()
		email   = req.GetEmail()
		channel = service.ToNotificationChannel(1)
		name    = req.GetName()
	)

	notification := &service.AuthNotificationRequest{
		UserID:              userID,
		Email:               email,
		Name:                name,
		NotificationChannel: channel,
	}

	if err := validator.ValidateAuthNotificationRequest(*notification); err != nil {
		h.logger.Error("Validation failed for SendRegistrationNotification", zap.Error(err))
		return &pb.SendNotificationResponse{Success: false, Error: err.Error()}, nil
	}

	if err := h.NotificationService.RegistrationNotification(ctx, notification.ToDTO(string(service.NotificationTypeRegistration))); err != nil {
		h.logger.Error("Failed to process SendRegistrationNotification", zap.Error(err))
		return &pb.SendNotificationResponse{Success: false, Error: err.Error()}, nil
	}

	h.logger.Info("Successfully processed SendRegistrationNotification",
		zap.Uint64("user_id", userID),
	)
	return &pb.SendNotificationResponse{Success: true}, nil
}

func (h *GrpcNotificationServiceHandler) SendLoginNotification(ctx context.Context, req *pb.AuthNotificationRequest) (*pb.SendNotificationResponse, error) {
	h.logger.Info("Received SendLoginNotification request", zap.Any("request", req))

	fmt.Printf("\n\n%+v\n\n", req)

	var (
		userID  = req.GetUserId()
		email   = req.GetEmail()
		name    = req.GetName()
		channel = service.ToNotificationChannel(1)
	)

	notification := &service.AuthNotificationRequest{
		UserID:              userID,
		Email:               email,
		Name:                name,
		NotificationChannel: channel,
	}

	if err := validator.ValidateAuthNotificationRequest(*notification); err != nil {
		h.logger.Error("Validation failed for SendLoginNotification", zap.Error(err))
		return &pb.SendNotificationResponse{Success: false, Error: err.Error()}, nil
	}

	if err := h.NotificationService.LoginNotification(ctx, notification.ToDTO(string(service.NotificationTypeLogin))); err != nil {
		h.logger.Error("Failed to process SendLoginNotification", zap.Error(err))
		return &pb.SendNotificationResponse{Success: false, Error: err.Error()}, nil
	}

	h.logger.Info("Successfully processed SendLoginNotification",
		zap.Uint64("user_id", userID),
	)

	return &pb.SendNotificationResponse{Success: true}, nil
}

func (h *GrpcNotificationServiceHandler) SendUserNotification(ctx context.Context, req *pb.UserNotificationRequest) (*pb.SendNotificationResponse, error) {
	h.logger.Info("Received SendUserNotification request", zap.Any("request", req))
	var (
		userID           = req.GetUserId()
		email            = req.GetEmail()
		name             = req.GetName()
		notificationType = req.GetType()
		channels         = req.GetChannels()
		metadata         = req.GetMetadata()
	)

	notification := &service.UserNotificationRequest{
		UserID:   userID,
		Email:    email,
		Name:     name,
		Type:     service.ToNotificationType(notificationType),
		Channels: service.ToNotificationChannels(channels),
		Metadata: metadata,
	}

	if err := validator.ValidateUserNotificationRequest(*notification); err != nil {
		h.logger.Error("Validation failed for SendUserNotification", zap.Error(err))
		return &pb.SendNotificationResponse{Success: false, Error: err.Error()}, nil
	}

	if err := h.NotificationService.UserNotification(ctx, notification.ToDTO()); err != nil {
		h.logger.Error("Failed to process SendUserNotification", zap.Error(err))
		return &pb.SendNotificationResponse{Success: false, Error: err.Error()}, nil
	}

	return &pb.SendNotificationResponse{Success: true}, nil
}
