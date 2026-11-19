// Package httpapi: HTTP handler. Chỉ parse / validate đầu vào, gọi service, trả JSON.
package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

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
	mux.HandleFunc("POST /api/auth/forgot-password/verify", a.forgotVerify)
	mux.HandleFunc("POST /api/auth/forgot-password/check", a.forgotCheck)
	mux.HandleFunc("POST /api/auth/forgot-password/reset", a.forgotReset)
	mux.HandleFunc("POST /api/auth/forgot-password/cancel", a.forgotCancel)
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

// ── AUTH-F02 Quên mật khẩu ───────────────────────────────────────────────────

const msgResetExpired = "Phiên xác minh đã hết hạn, vui lòng xác minh lại"

type forgotVerifyRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (a *API) forgotVerify(w http.ResponseWriter, r *http.Request) {
	var req forgotVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_INPUT", "Dữ liệu gửi lên không hợp lệ")
		return
	}
	token, u, err := a.Auth.VerifyForReset(req.Username, req.Email)
	switch {
	case errors.Is(err, auth.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, "INVALID_INPUT", "Nhập tên đăng nhập và email")
	case errors.Is(err, auth.ErrVerifyFailed):
		writeError(w, http.StatusBadRequest, "VERIFY_FAILED", "Tên đăng nhập hoặc email không đúng")
	case err != nil:
		slog.Error("forgot verify", "err", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL", "Lỗi hệ thống")
	default:
		writeJSON(w, http.StatusOK, map[string]string{"resetToken": token, "username": u.Username})
	}
}

type resetTokenRequest struct {
	ResetToken string `json:"resetToken"`
}

func (a *API) forgotCheck(w http.ResponseWriter, r *http.Request) {
	var req resetTokenRequest
	_ = json.NewDecoder(r.Body).Decode(&req) // body sai → token rỗng → hết hạn
	u, err := a.Auth.CheckReset(req.ResetToken)
	if err != nil {
		writeError(w, http.StatusBadRequest, "RESET_EXPIRED", msgResetExpired)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"username": u.Username})
}

type forgotResetRequest struct {
	ResetToken string `json:"resetToken"`
	Password   string `json:"password"`
	Confirm    string `json:"confirm"`
}

func (a *API) forgotReset(w http.ResponseWriter, r *http.Request) {
	var req forgotResetRequest
	// Body sai coi như rỗng: service kiểm phiên (AUTH-F02-US02-A4) trước mọi lỗi nhập.
	_ = json.NewDecoder(r.Body).Decode(&req)
	err := a.Auth.ResetPassword(req.ResetToken, req.Password, req.Confirm)
	switch {
	case err == nil:
		w.WriteHeader(http.StatusNoContent)
	case errors.Is(err, auth.ErrResetExpired):
		writeError(w, http.StatusBadRequest, "RESET_EXPIRED", msgResetExpired)
	case errors.Is(err, auth.ErrPasswordEmpty):
		writeError(w, http.StatusBadRequest, "INVALID_INPUT", "Nhập mật khẩu mới và nhập lại mật khẩu")
	case errors.Is(err, auth.ErrPasswordTooShort):
		writeError(w, http.StatusBadRequest, "PASSWORD_TOO_SHORT", "Mật khẩu phải có ít nhất 8 ký tự")
	case errors.Is(err, auth.ErrPasswordTooLong):
		writeError(w, http.StatusBadRequest, "PASSWORD_TOO_LONG", "Mật khẩu tối đa 16 ký tự")
	case errors.Is(err, auth.ErrPasswordMismatch):
		writeError(w, http.StatusBadRequest, "PASSWORD_MISMATCH", "Mật khẩu nhập lại không khớp")
	default:
		slog.Error("forgot reset", "err", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL", "Lỗi hệ thống")
	}
}

func (a *API) forgotCancel(w http.ResponseWriter, r *http.Request) {
	var req resetTokenRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	a.Auth.CancelReset(req.ResetToken)
	w.WriteHeader(http.StatusNoContent)
}

// ── API cho e2e (chỉ khi TEST_API=1) ─────────────────────────────────────────

type createUserRequest struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"` // tuỳ chọn; bỏ trống = tài khoản chưa có email
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
	u, err := a.Users.Create(store.User{Username: req.Username, Name: req.Name, Email: strings.TrimSpace(req.Email), PasswordHash: hash})
	if errors.Is(err, store.ErrExists) {
		writeError(w, http.StatusConflict, "USERNAME_TAKEN", "Tên đăng nhập đã tồn tại")
		return
	}
	if errors.Is(err, store.ErrEmailExists) {
		writeError(w, http.StatusConflict, "EMAIL_TAKEN", "Email đã được dùng cho tài khoản khác")
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
