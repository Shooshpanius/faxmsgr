package auth

import (
	"encoding/json"
	"net/http"
)

// requestCodeReq тело запроса /auth/request-code
type requestCodeReq struct {
	Phone string `json:"phone"`
}

// verifyCodeReq тело запроса /auth/verify-code
type verifyCodeReq struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

// refreshReq тело запроса /auth/refresh
type refreshReq struct {
	RefreshToken string `json:"refresh_token"`
}

// tokenResp ответ с парой токенов
type tokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// MakeRequestCodeHandler возвращает HTTP-обработчик POST /auth/request-code.
func MakeRequestCodeHandler(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req requestCodeReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Phone == "" {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}
		if err := svc.RequestCode(r.Context(), req.Phone); err != nil {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// MakeVerifyCodeHandler возвращает HTTP-обработчик POST /auth/verify-code.
func MakeVerifyCodeHandler(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req verifyCodeReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Phone == "" || req.Code == "" {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}
		access, refresh, err := svc.VerifyCode(r.Context(), req.Phone, req.Code)
		if err != nil {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokenResp{AccessToken: access, RefreshToken: refresh}) //nolint:errcheck
	}
}

// MakeRefreshHandler возвращает HTTP-обработчик POST /auth/refresh.
func MakeRefreshHandler(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req refreshReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}
		access, refresh, err := svc.Refresh(r.Context(), req.RefreshToken)
		if err != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokenResp{AccessToken: access, RefreshToken: refresh}) //nolint:errcheck
	}
}
