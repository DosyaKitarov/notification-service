package handler

import (
	"context"

	pb "github.com/DosyaKitarov/notification-service/pkg/grpc"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type KafkaHandler struct {
	GrpcNotificationServiceHandler
	logger *zap.Logger
}

func NewKafkaHandler(notificationServiceHandler *GrpcNotificationServiceHandler, logger *zap.Logger) *KafkaHandler {
	return &KafkaHandler{
		GrpcNotificationServiceHandler: *notificationServiceHandler,
		logger:                         logger,
	}
}

func (kh *KafkaHandler) ListenAndServe(ctx context.Context, brokers []string, topic string) error {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		kh.logger.Error("Couldn't create consumer", zap.Error(err))
		return err
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		kh.logger.Error("Error subscribing to topic", zap.Error(err))
		return err
	}
	defer partitionConsumer.Close()

	kh.logger.Info("Waiting for messages from producer...")

	for msg := range partitionConsumer.Messages() {
		kh.logger.Info("Received message", zap.String("message", string(msg.Value)))
		var notificaion pb.UserNotificationRequest

		err := proto.Unmarshal(msg.Value, &notificaion)
		if err != nil {
			kh.logger.Error("Proto unmarshall error: %v", zap.Error(err))
			continue
		}
		_, err = kh.GrpcNotificationServiceHandler.SendUserNotification(ctx, &notificaion)
		if err != nil {
			kh.logger.Error("Error occured on send user notification", zap.Error(err))
			continue
		}
	}

	return nil
}
