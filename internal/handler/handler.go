package handler

import (
	"context"
	"database/sql"

	"github.com/DosyaKitarov/notification-service/internal/service"
	pb "github.com/DosyaKitarov/notification-service/pkg/grpc"
	"github.com/DosyaKitarov/notification-service/pkg/validator"
	"go.uber.org/zap"
)

const (
	NotificationTypeRegistration = "registration_confirmation"
	NotificationTypeLogin        = "login_alert"
	NotificationTypeInvestment   = "investment_success"
	NotificationTypeInvestor     = "invested_in_you"
)

type NotificationServiceHandler struct { // Add logic to handle registration notification
	*service.NotificationService
	pb.UnimplementedNotificationServiceServer
	db     *sql.DB
	logger *zap.Logger
}

type Service interface {
	RegistrationNotification(ctx context.Context, notification service.AuthNotificationRequest) error
}

func NewNotificationServiceHandler(db *sql.DB, service *service.NotificationService, logger *zap.Logger) *NotificationServiceHandler {
	return &NotificationServiceHandler{
		db:                  db,
		NotificationService: service,
		logger:              logger,
	}
}

func (h *NotificationServiceHandler) SendRegistrationNotification(ctx context.Context, req *pb.AuthNotificationRequest) (*pb.SendNotificationResponse, error) {
	h.logger.Info("Received SendRegistrationNotification request",
		zap.Uint64("user_id", req.GetUserId()),
		zap.String("channel", req.GetChannel().String()),
	)

	var (
		userID   = req.GetUserId()
		email    = req.GetEmail()
		channel  = service.ToNotificationChannel(req.GetChannel())
		metaData = req.GetMetadata()
	)

	notification := &service.AuthNotificationRequest{
		UserID:              userID,
		Email:               email,
		NotificationChannel: channel,
		Metadata:            metaData,
	}

	if err := validator.ValidateAuthNotificationRequest(*notification); err != nil {
		h.logger.Error("Validation failed for SendRegistrationNotification", zap.Error(err))
		return &pb.SendNotificationResponse{Success: false, Error: err.Error()}, nil
	}

	if err := h.NotificationService.RegistrationNotification(ctx, notification.ToDTO(NotificationTypeRegistration)); err != nil {
		h.logger.Error("Failed to process SendRegistrationNotification", zap.Error(err))
		return &pb.SendNotificationResponse{Success: false, Error: err.Error()}, nil
	}

	h.logger.Info("Successfully processed SendRegistrationNotification",
		zap.Uint64("user_id", userID),
	)
	return &pb.SendNotificationResponse{Success: true}, nil
}

func (h *NotificationServiceHandler) SendLoginNotification(ctx context.Context, req *pb.AuthNotificationRequest) (*pb.SendNotificationResponse, error) {
	// Add logic to handle login notification
	return &pb.SendNotificationResponse{Success: true}, nil
}

func (h *NotificationServiceHandler) SendInvestmentNotification(ctx context.Context, req *pb.UserNotificationRequest) (*pb.SendNotificationResponse, error) {
	// Add logic to handle investment notification
	return &pb.SendNotificationResponse{Success: true}, nil
}

func (h *NotificationServiceHandler) SendInvestorNotification(ctx context.Context, req *pb.UserNotificationRequest) (*pb.SendNotificationResponse, error) {
	// Add logic to handle investor notification
	return &pb.SendNotificationResponse{Success: true}, nil
}
