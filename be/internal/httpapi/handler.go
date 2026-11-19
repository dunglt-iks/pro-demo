// Package httpapi: HTTP handler. Chỉ parse / validate đầu vào, gọi service, trả JSON.
package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"project1/be/internal/auth"
	"project1/be/internal/store"
)

// API gom các phụ thuộc của handler.
type API struct {
	Auth    *auth.Service
	Users   *store.Store
	TestAPI bool // bật /api/test/* (chỉ dùng khi chạy e2e: TEST_API=1)
}

// Routes trả mux với mọi endpoint.
func (a *API) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("POST /api/auth/login", a.login)
	mux.Handle("GET /api/me", a.requireAuth(http.HandlerFunc(a.me)))
	if a.TestAPI {
		mux.HandleFunc("POST /api/test/users", a.testCreateUser)
		mux.HandleFunc("DELETE /api/test/users/{username}", a.testDeleteUser)
		mux.HandleFunc("POST /api/test/reset", a.testReset)
	}
	return logRequests(mux)
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *API) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_INPUT", "Dữ liệu gửi lên không hợp lệ")
		return
	}
	token, u, err := a.Auth.Login(req.Username, req.Password)
	switch {
	case errors.Is(err, auth.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, "INVALID_INPUT", "Nhập tên đăng nhập và mật khẩu")
	case errors.Is(err, auth.ErrInvalidCredentials):
		writeError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Sai tên đăng nhập hoặc mật khẩu")
	case err != nil:
		slog.Error("login", "err", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL", "Lỗi hệ thống")
	default:
		writeJSON(w, http.StatusOK, map[string]any{"token": token, "user": u})
	}
}

func (a *API) me(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, userFrom(r))
}

// ── API cho e2e (chỉ khi TEST_API=1) ─────────────────────────────────────────

type createUserRequest struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (a *API) testCreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Username == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "INVALID_INPUT", "Cần username, password")
		return
	}
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL", "Lỗi hệ thống")
		return
	}
	if req.Name == "" {
		req.Name = req.Username
	}
	u, err := a.Users.Create(store.User{Username: req.Username, Name: req.Name, PasswordHash: hash})
	if errors.Is(err, store.ErrExists) {
		writeError(w, http.StatusConflict, "USERNAME_TAKEN", "Tên đăng nhập đã tồn tại")
		return
	}
	writeJSON(w, http.StatusCreated, u)
}

func (a *API) testDeleteUser(w http.ResponseWriter, r *http.Request) {
	_ = a.Users.Delete(r.PathValue("username")) // xoá rồi / không có: vẫn coi là xong
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) testReset(w http.ResponseWriter, r *http.Request) {
	a.Users.Reset()
	w.WriteHeader(http.StatusNoContent)
}
