// Package auth: đăng nhập bằng username / password, cấp và kiểm JWT.
package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"project1/be/internal/store"
)

// Lỗi nghiệp vụ — handler đổi sang mã lỗi API.
var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

// Service xử lý đăng nhập.
type Service struct {
	users  *store.Store
	secret []byte
	ttl    time.Duration
}

// NewService tạo service; secret dùng ký JWT (HS256), ttl là thời hạn token.
func NewService(users *store.Store, secret string, ttl time.Duration) *Service {
	return &Service{users: users, secret: []byte(secret), ttl: ttl}
}

// HashPassword băm mật khẩu bằng bcrypt.
func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

// Login kiểm username / password, trả JWT và thông tin tài khoản.
// Sai username hay sai mật khẩu đều trả cùng một lỗi, không để lộ username nào tồn tại.
func (s *Service) Login(username, password string) (string, store.User, error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return "", store.User{}, ErrInvalidInput
	}
	u, err := s.users.FindByUsername(username)
	if err != nil {
		return "", store.User{}, ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return "", store.User{}, ErrInvalidCredentials
	}
	token, err := s.issue(u)
	if err != nil {
		return "", store.User{}, err
	}
	return token, u, nil
}

// UserFromToken kiểm JWT và trả tài khoản tương ứng (tài khoản đã bị xoá → token hết hiệu lực).
func (s *Service) UserFromToken(token string) (store.User, error) {
	claims := jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) { return s.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())
	if err != nil {
		return store.User{}, ErrInvalidToken
	}
	u, err := s.users.FindByUsername(claims.Subject)
	if err != nil {
		return store.User{}, ErrInvalidToken
	}
	return u, nil
}

func (s *Service) issue(u store.User) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   u.Username,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}
