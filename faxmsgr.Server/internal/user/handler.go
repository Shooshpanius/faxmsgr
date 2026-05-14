package user

import (
	"encoding/json"
	"net/http"

	"faxmsgr/server/internal/middleware"
)

// MakeGetProfileHandler возвращает обработчик GET /users/profile.
func MakeGetProfileHandler(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.UserIDFromContext(r.Context())
		profile, err := svc.GetProfile(r.Context(), userID)
		if err != nil {
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile) //nolint:errcheck
	}
}

// MakePutProfileHandler возвращает обработчик PUT /users/profile.
func MakePutProfileHandler(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.UserIDFromContext(r.Context())
		var req UpdateProfileReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}
		profile, err := svc.UpdateProfile(r.Context(), userID, req)
		if err != nil {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile) //nolint:errcheck
	}
}
