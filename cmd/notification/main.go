package main

import (
	"net"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"

	"github.com/DosyaKitarov/notification-service/internal/handler"
	"github.com/DosyaKitarov/notification-service/internal/repository"
	"github.com/DosyaKitarov/notification-service/internal/service"
	"github.com/DosyaKitarov/notification-service/pkg/config"
	"github.com/DosyaKitarov/notification-service/pkg/database"
	"github.com/DosyaKitarov/notification-service/pkg/email"
	pb "github.com/DosyaKitarov/notification-service/pkg/grpc"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		logger.Fatal("Failed to read config", zap.Error(err))
	}
	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		logger.Fatal("Failed to unmarshal config", zap.Error(err))
	}

	// DB connection setup
	db, err := database.ConnectToDB(cfg)
	if err != nil {
		logger.Fatal("DB connect error", zap.Error(err))
	}
	defer db.Close()

	logger.Info("DB connected",
		zap.String("host", cfg.Database.Host),
		zap.Int("port", cfg.Database.Port),
		zap.String("user", cfg.Database.User),
		zap.String("dbname", cfg.Database.DBName),
	)

	// gRPC server setup
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.Fatal("Failed to listen on port 8080", zap.Error(err))
	}

	emailSender := &email.EmailSender{
		Sender:   cfg.Smtp.Sender,
		Username: cfg.Smtp.Username,
		Password: cfg.Smtp.Password,
	}

	err = emailSender.LoadTemplates("pkg/email/templates.json")
	if err != nil {
		logger.Fatal("Failed to load templates", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	repo := repository.NewRepository(db, logger)
	service := service.NewNotificationService(repo, *emailSender, logger)
	handler := handler.NewNotificationServiceHandler(db, service, logger)

	pb.RegisterNotificationServiceServer(grpcServer, handler)

	logger.Info("Starting gRPC server on port 8080")
	if err := grpcServer.Serve(listener); err != nil {
		logger.Fatal("Failed to serve gRPC server", zap.Error(err))
	}
}
