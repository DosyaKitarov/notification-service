package handler

import (
	"context"
	"database/sql"
	"fmt"

	sync "github.com/DosyaKitarov/notification-service/internal/sync"
	pb "github.com/DosyaKitarov/notification-service/pkg/grpc"
	"github.com/DosyaKitarov/notification-service/pkg/validator"
	"go.uber.org/zap"
)

type NotificationServiceHandler struct { // Add logic to handle registration notification
	SyncService
	pb.UnimplementedNotificationServiceServer
	db     *sql.DB
	logger *zap.Logger
}

type SyncService interface {
	RegistrationNotification(ctx context.Context, notification sync.AuthNotificationRequestDTO) error
	LoginNotification(ctx context.Context, notification sync.AuthNotificationRequestDTO) error
	UserNotification(ctx context.Context, notification sync.UserNotificationRequestDTO) error
}

func NewNotificationServiceHandler(db *sql.DB, syncService SyncService, logger *zap.Logger) *NotificationServiceHandler {
	return &NotificationServiceHandler{
		db:          db,
		SyncService: syncService,
		logger:      logger,
	}
}

func (h *NotificationServiceHandler) SendRegistrationNotification(ctx context.Context, req *pb.AuthNotificationRequest) (*pb.SendNotificationResponse, error) {
	h.logger.Info("Received SendRegistrationNotification request", zap.Any("request", req))

	var (
		userID  = req.GetUserId()
		email   = req.GetEmail()
		channel = sync.ToNotificationChannel(1)
	)

	notification := &sync.AuthNotificationRequest{
		UserID:              userID,
		Email:               email,
		NotificationChannel: channel,
	}

	if err := validator.ValidateAuthNotificationRequest(*notification); err != nil {
		h.logger.Error("Validation failed for SendRegistrationNotification", zap.Error(err))
		return &pb.SendNotificationResponse{Success: false, Error: err.Error()}, nil
	}

	if err := h.SyncService.RegistrationNotification(ctx, notification.ToDTO(string(sync.NotificationTypeRegistration))); err != nil {
		h.logger.Error("Failed to process SendRegistrationNotification", zap.Error(err))
		return &pb.SendNotificationResponse{Success: false, Error: err.Error()}, nil
	}

	h.logger.Info("Successfully processed SendRegistrationNotification",
		zap.Uint64("user_id", userID),
	)
	return &pb.SendNotificationResponse{Success: true}, nil
}

func (h *NotificationServiceHandler) SendLoginNotification(ctx context.Context, req *pb.AuthNotificationRequest) (*pb.SendNotificationResponse, error) {
	h.logger.Info("Received SendLoginNotification request", zap.Any("request", req))

	fmt.Printf("\n\n%+v\n\n", req)

	var (
		userID  = req.GetUserId()
		email   = req.GetEmail()
		name    = req.GetName()
		channel = sync.ToNotificationChannel(1)
	)

	notification := &sync.AuthNotificationRequest{
		UserID:              userID,
		Email:               email,
		Name:                name,
		NotificationChannel: channel,
	}

	if err := validator.ValidateAuthNotificationRequest(*notification); err != nil {
		h.logger.Error("Validation failed for SendLoginNotification", zap.Error(err))
		return &pb.SendNotificationResponse{Success: false, Error: err.Error()}, nil
	}

	if err := h.SyncService.LoginNotification(ctx, notification.ToDTO(string(sync.NotificationTypeLogin))); err != nil {
		h.logger.Error("Failed to process SendLoginNotification", zap.Error(err))
		return &pb.SendNotificationResponse{Success: false, Error: err.Error()}, nil
	}

	h.logger.Info("Successfully processed SendLoginNotification",
		zap.Uint64("user_id", userID),
	)

	return &pb.SendNotificationResponse{Success: true}, nil
}

func (h *NotificationServiceHandler) SendUserNotification(ctx context.Context, req *pb.UserNotificationRequest) (*pb.SendNotificationResponse, error) {
	h.logger.Info("Received SendUserNotification request", zap.Any("request", req))
	var (
		userID           = req.GetUserId()
		email            = req.GetEmail()
		name             = req.GetName()
		notificationType = req.GetType()
		channels         = req.GetChannels()
		metadata         = req.GetMetadata()
	)

	notification := &sync.UserNotificationRequest{
		UserID:   userID,
		Email:    email,
		Name:     name,
		Type:     sync.ToNotificationType(notificationType),
		Channels: sync.ToNotificationChannels(channels),
		Metadata: metadata,
	}

	if err := validator.ValidateUserNotificationRequest(*notification); err != nil {
		h.logger.Error("Validation failed for SendUserNotification", zap.Error(err))
		return &pb.SendNotificationResponse{Success: false, Error: err.Error()}, nil
	}

	if err := h.SyncService.UserNotification(ctx, notification.ToDTO()); err != nil {
		h.logger.Error("Failed to process SendUserNotification", zap.Error(err))
		return &pb.SendNotificationResponse{Success: false, Error: err.Error()}, nil
	}

	return &pb.SendNotificationResponse{Success: true}, nil
}
