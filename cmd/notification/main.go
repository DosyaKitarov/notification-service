package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"

	"github.com/DosyaKitarov/notification-service/internal/handler"
	"github.com/DosyaKitarov/notification-service/internal/notificaitonService"
	"github.com/DosyaKitarov/notification-service/pkg/config"
	"github.com/DosyaKitarov/notification-service/pkg/database"
	"github.com/DosyaKitarov/notification-service/pkg/email"
	pb "github.com/DosyaKitarov/notification-service/pkg/grpc"
)

func main() {

	logger := initLogger()
	defer logger.Sync()

	ctx := context.Background()
	cfg := loadConfig(logger)
	db := initDB(cfg, logger)
	defer db.Close()

	grpcPort := cfg.Ports.GRPC
	wsPort := cfg.Ports.WS
	httpPort := cfg.Ports.HTTP

	emailSender := initEmailSender(cfg, logger)
	repo := notificaitonService.NewRepository(db, logger)
	service := notificaitonService.NewNotificationService(repo, *emailSender, logger)
	wsHandler := handler.NewWSHandler(service, logger)
	grpcHandler := handler.NewNotificationServiceHandler(db, service, logger)
	service.SetWebNotifier(wsHandler)
	kafkaHandler := handler.NewKafkaHandler(grpcHandler, logger)
	restHandler := handler.NewRestNotificationServiceHandler(db, service, logger)
	router := handler.NewRouter(restHandler)
	go startHTTPServer(httpPort, router, logger)
	go startWebSocketServer(wsPort, wsHandler, logger)
	go startGRPCServer(grpcPort, grpcHandler, logger)
	startConsumer(ctx, kafkaHandler, cfg, logger)
}

func initLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	return logger
}

func loadConfig(logger *zap.Logger) config.Config {
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		logger.Fatal("Failed to read config", zap.Error(err))
	}
	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		logger.Fatal("Failed to unmarshal config", zap.Error(err))
	}
	return cfg
}

func initDB(cfg config.Config, logger *zap.Logger) *sql.DB {
	db, err := database.ConnectToDB(cfg)
	if err != nil {
		logger.Fatal("DB connect error", zap.Error(err))
	}
	logger.Info("DB connected",
		zap.String("host", cfg.Database.Host),
		zap.Int("port", cfg.Database.Port),
		zap.String("user", cfg.Database.User),
		zap.String("dbname", cfg.Database.DBName),
	)
	return db
}

func initEmailSender(cfg config.Config, logger *zap.Logger) *email.EmailSender {
	emailSender := &email.EmailSender{
		Sender:   cfg.Smtp.Sender,
		Username: cfg.Smtp.Username,
		Password: cfg.Smtp.Password,
	}
	if err := emailSender.LoadTemplates("pkg/email/templates.json"); err != nil {
		logger.Fatal("Failed to load templates", zap.Error(err))
	}
	return emailSender
}

func startWebSocketServer(port int, wsHandler *handler.WSHandler, logger *zap.Logger) {
	http.HandleFunc("/ws", wsHandler.ServeWS)
	addr := fmt.Sprintf(":%d", port)
	logger.Info("Starting WebSocket server", zap.Int("port", port))
	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Fatal("Failed to start WebSocket server", zap.Error(err))
	}
}

func startGRPCServer(port int, grpcHandler pb.NotificationServiceServer, logger *zap.Logger) {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("Failed to listen on port", zap.Int("port", port), zap.Error(err))
	}
	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, grpcHandler)
	logger.Info("Starting gRPC server", zap.Int("port", port))
	if err := grpcServer.Serve(listener); err != nil {
		logger.Fatal("Failed to serve gRPC server", zap.Error(err))
	}
}

func startConsumer(ctx context.Context, kh *handler.KafkaHandler, cfg config.Config, logger *zap.Logger) {
	err := kh.ListenAndServe(ctx, cfg.Ports.Kafka, cfg.Kafka.Topic)
	if err != nil {
		logger.Fatal("Kafka consumer ended with error", zap.Error(err))
	}
}

func startHTTPServer(port int, router http.Handler, logger *zap.Logger) {
	addr := fmt.Sprintf(":%d", port)
	logger.Info("Starting HTTP server", zap.Int("port", port))
	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("Failed to start HTTP server", zap.Error(err))
	}
}
