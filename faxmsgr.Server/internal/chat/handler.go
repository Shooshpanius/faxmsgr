package chat

import (
	"encoding/json"
	"net/http"

	"faxmsgr/server/internal/middleware"
)

// MakeCreateChatHandler возвращает обработчик POST /chats.
func MakeCreateChatHandler(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.UserIDFromContext(r.Context())
		var req CreateChatReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}
		chat, err := svc.CreateChat(r.Context(), userID, req)
		if err != nil {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(chat) //nolint:errcheck
	}
}

// MakeListChatsHandler возвращает обработчик GET /chats.
func MakeListChatsHandler(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.UserIDFromContext(r.Context())
		chats, err := svc.ListChats(r.Context(), userID)
		if err != nil {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}
		if chats == nil {
			chats = []*Chat{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chats) //nolint:errcheck
	}
}
