package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/DosyaKitarov/notification-service/internal/notificaitonService"
	pb "github.com/DosyaKitarov/notification-service/pkg/grpc"
	"github.com/DosyaKitarov/notification-service/pkg/validator"
	"go.uber.org/zap"
)

type RestNotificationServiceHandler struct {
	NotificationService
	pb.UnimplementedNotificationServiceServer
	db     *sql.DB
	logger *zap.Logger
}

func NewRestNotificationServiceHandler(db *sql.DB, NotificationService NotificationService, Logger *zap.Logger) *RestNotificationServiceHandler {
	return &RestNotificationServiceHandler{
		db:                  db,
		NotificationService: NotificationService,
		logger:              Logger,
	}
}

func (h *RestNotificationServiceHandler) GetEmailNotifications(w http.ResponseWriter, r *http.Request) {
	var (
		ctx     = r.Context()
		perPage = r.URL.Query().Get("per_page")
		page    = r.URL.Query().Get("page")
	)

	h.logger.Info("Received GetEmailNotifications request", zap.String("per_page", perPage), zap.String("page", page))

	GetEmailNotification := notificaitonService.GetNotificationsRequest{
		PerPage: validator.ParsePerPageParam(perPage),
		Page:    validator.ParsePageNumParam(page),
	}

	response, err := h.NotificationService.GetEmailNotifications(ctx, GetEmailNotification)
	if err != nil {
		h.logger.Error("Failed to get email notifications", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = statusResponse(w, &R{
		Status:     http.StatusOK,
		StatusCode: http.StatusText(http.StatusOK),
		Data:       response,
	})

}

func (h *RestNotificationServiceHandler) GetWebNotifications(w http.ResponseWriter, r *http.Request) {
	var (
		ctx     = r.Context()
		perPage = r.URL.Query().Get("per_page")
		page    = r.URL.Query().Get("page")
	)

	h.logger.Info("Received GetWebNotifications request", zap.String("per_page", perPage), zap.String("page", page))

	GetWebNotification := notificaitonService.GetNotificationsRequest{
		PerPage: validator.ParsePerPageParam(perPage),
		Page:    validator.ParsePageNumParam(page),
	}

	response, err := h.NotificationService.GetWebNotifications(ctx, GetWebNotification)
	if err != nil {
		h.logger.Error("Failed to get web notifications", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = statusResponse(w, &R{
		Status:     http.StatusOK,
		StatusCode: http.StatusText(http.StatusOK),
		Data:       response,
	})
}

func statusResponse(w http.ResponseWriter, response *R) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	bytes, err := json.MarshalIndent(response, " ", "	")
	if err != nil {
		return err
	}

	if _, err = w.Write(bytes); err != nil {
		return err
	}

	return nil
}

type R struct {
	Status     int         `json:"-"`
	StatusCode string      `json:"status_code"`
	Data       interface{} `json:"data,omitempty"`
	Field      *string     `json:"field,omitempty"`
}
