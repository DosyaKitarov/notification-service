package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(restHandler *RestNotificationServiceHandler) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/notifications/email", restHandler.GetEmailNotifications).Methods("GET")
	return r
}
