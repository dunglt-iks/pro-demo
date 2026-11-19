// Server API của project1. Biến môi trường:
//
//	PORT        cổng (mặc định 8080)
//	JWT_SECRET  khoá ký JWT (mặc định giá trị cho môi trường dev — đổi khi deploy)
//	TEST_API=1  bật /api/test/* cho e2e (tạo / xoá tài khoản test, reset dữ liệu)
package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"project1/be/internal/auth"
	"project1/be/internal/httpapi"
	"project1/be/internal/store"
)

func main() {
	port := env("PORT", "8080")
	secret := env("JWT_SECRET", "dev-secret-change-me")

	seed := []struct{ username, name, password string }{
		{"admin", "Quản trị viên", "Admin@123"},
		{"user1", "Nguyễn Văn A", "User@123"},
	}
	var users []store.User
	for _, s := range seed {
		hash, err := auth.HashPassword(s.password)
		if err != nil {
			slog.Error("hash seed", "err", err)
			os.Exit(1)
		}
		users = append(users, store.User{Username: s.username, Name: s.name, PasswordHash: hash})
	}
	st := store.New(users)

	api := &httpapi.API{
		Auth:    auth.NewService(st, secret, time.Hour),
		Users:   st,
		TestAPI: os.Getenv("TEST_API") == "1",
	}
	srv := &http.Server{Addr: ":" + port, Handler: api.Routes(), ReadHeaderTimeout: 5 * time.Second}
	slog.Info("listening", "port", port, "testAPI", api.TestAPI)
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("server", "err", err)
		os.Exit(1)
	}
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
