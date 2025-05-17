package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"

	"github.com/DosyaKitarov/notification-service/internal/handler"
	"github.com/DosyaKitarov/notification-service/internal/syncNotificaiton"
	"github.com/DosyaKitarov/notification-service/pkg/config"
	"github.com/DosyaKitarov/notification-service/pkg/database"
	"github.com/DosyaKitarov/notification-service/pkg/email"
	pb "github.com/DosyaKitarov/notification-service/pkg/grpc"
)

func main() {
	grpcPort := flag.Int("grpc-port", 8080, "gRPC server port")
	wsPort := flag.Int("ws-port", 6969, "WebSocket server port")
	flag.Parse()

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
	addr := fmt.Sprintf(":%d", *grpcPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("Failed to listen on port", zap.Int("port", *grpcPort), zap.Error(err))
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
	repo := syncNotificaiton.NewRepository(db, logger)
	service := syncNotificaiton.NewNotificationService(repo, *emailSender, logger)
	wsHandler := handler.NewWSHandler(service, logger)
	handler := handler.NewNotificationServiceHandler(db, service, logger)
	service.SetWebNotifier(wsHandler)

	go func() {
		http.HandleFunc("/ws", wsHandler.ServeWS)
		wsAddr := fmt.Sprintf(":%d", *wsPort)
		logger.Info("Starting WebSocket server", zap.Int("port", *wsPort))
		if err := http.ListenAndServe(wsAddr, nil); err != nil {
			logger.Fatal("Failed to start WebSocket server", zap.Error(err))
		}
	}()

	pb.RegisterNotificationServiceServer(grpcServer, handler)

	logger.Info("Starting gRPC server", zap.Int("port", *grpcPort))
	if err := grpcServer.Serve(listener); err != nil {
		logger.Fatal("Failed to serve gRPC server", zap.Error(err))
	}

}
