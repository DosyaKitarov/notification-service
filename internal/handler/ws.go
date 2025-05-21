package handler

import (
	"encoding/json"
	"net/http"
	"sync"

	syncService "github.com/DosyaKitarov/notification-service/internal/notificaitonService"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type WSHandler struct {
	Service      *syncService.NotificationService
	Logger       *zap.Logger
	connections  map[uint64]*websocket.Conn
	connectionsM sync.RWMutex
}

func NewWSHandler(service *syncService.NotificationService, logger *zap.Logger) *WSHandler {
	return &WSHandler{
		Service:     service,
		Logger:      logger,
		connections: make(map[uint64]*websocket.Conn),
	}
}

type AuthMessage struct {
	UserID uint64 `json:"user_id"`
}

type IsReadMessage struct {
	Id uint64 `json:"is_read"`
}

type NotificationMessage struct {
	Id       uint64            `json:"id"`
	Metadata map[string]string `json:"metadata"`
}

func (h *WSHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.Logger.Error("Failed to upgrade connection", zap.Error(err))
		return
	}
	defer conn.Close()

	// Ожидаем первое сообщение с user_id
	_, msg, err := conn.ReadMessage()
	if err != nil {
		h.Logger.Error("Failed to read auth message", zap.Error(err))
		return
	}
	var auth AuthMessage
	if err := json.Unmarshal(msg, &auth); err != nil || auth.UserID == 0 {
		h.Logger.Error("Invalid auth message", zap.Error(err))
		conn.WriteMessage(websocket.TextMessage, []byte(`{"status":"invalid id"}`))
		return
	}

	// Если всё ок, отправляем "ok"
	conn.WriteMessage(websocket.TextMessage, []byte(`{"status":"ok"}`))

	// Сохраняем соединение
	h.connectionsM.Lock()
	h.connections[auth.UserID] = conn
	h.connectionsM.Unlock()
	h.Logger.Info("User connected", zap.Uint64("user_id", auth.UserID))

	unread, err := h.Service.GetUnreadWebNotifications(r.Context(), auth.UserID)
	if err != nil {
		h.Logger.Error("Failed to get unread notifications", zap.Error(err))
		conn.WriteMessage(websocket.TextMessage, []byte(`{"status":"error","error":"failed to get unread notifications"}`))
	} else {
		resp, _ := json.Marshal(map[string]interface{}{
			"action": "unread_notifications",
			"data":   unread,
		})
		conn.WriteMessage(websocket.TextMessage, resp)
	}
	// Пример: слушаем дальнейшие сообщения (можно убрать если не нужно)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			h.Logger.Info("User disconnected", zap.Uint64("user_id", auth.UserID))
			break
		}
		h.Logger.Info("Received WS message", zap.ByteString("msg", msg))
		var isReadMessage IsReadMessage
		if err := json.Unmarshal(msg, &isReadMessage); err != nil {
			h.Logger.Error("Failed to unmarshal message", zap.Error(err))
			continue
		}
		if err := h.Service.MarkNotificationAsRead(r.Context(), isReadMessage.Id); err != nil {
			h.Logger.Error("Failed to mark notification as read", zap.Error(err))
		}
	}

	// Удаляем соединение при отключении
	h.connectionsM.Lock()
	delete(h.connections, auth.UserID)
	h.connectionsM.Unlock()
}

// Метод для отправки уведомления пользователю
func (h *WSHandler) SendToUser(userID uint64, data []byte) error {
	h.connectionsM.RLock()
	conn, ok := h.connections[userID]
	h.connectionsM.RUnlock()
	if !ok {
		return nil // или ошибка, если нужно
	}
	return conn.WriteMessage(websocket.TextMessage, data)
}
