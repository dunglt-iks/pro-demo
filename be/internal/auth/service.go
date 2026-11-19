// Package auth: đăng nhập bằng username / password, cấp và kiểm JWT, quên mật khẩu (AUTH-F02).
package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"project1/be/internal/store"
)

// Lỗi nghiệp vụ — handler đổi sang mã lỗi API.
var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")

	ErrVerifyFailed     = errors.New("verify failed")      // AUTH-F02-US01-A2
	ErrResetExpired     = errors.New("reset expired")      // AUTH-F02-US02-A4
	ErrPasswordEmpty    = errors.New("password empty")     // AUTH-F02-US02-A1
	ErrPasswordTooShort = errors.New("password too short") // AUTH-F02-US02-A2
	ErrPasswordTooLong  = errors.New("password too long")  // AUTH-F02-US02-A5
	ErrPasswordMismatch = errors.New("password mismatch")  // AUTH-F02-US02-A3
)

// Giới hạn của AUTH-BR04, AUTH-BR06.
const (
	PasswordMinLen = 8
	PasswordMaxLen = 16
	ResetTTL       = 10 * time.Minute
)

// Service xử lý đăng nhập và đặt lại mật khẩu.
type Service struct {
	users  *store.Store
	secret []byte
	ttl    time.Duration
	now    func() time.Time // đổi được trong test (hạn phiên đặt lại)
}

// NewService tạo service; secret dùng ký JWT (HS256), ttl là thời hạn token.
func NewService(users *store.Store, secret string, ttl time.Duration) *Service {
	return &Service{users: users, secret: []byte(secret), ttl: ttl, now: time.Now}
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

// claims: JWT chuẩn + Ver = TokenVersion của tài khoản lúc cấp.
type claims struct {
	jwt.RegisteredClaims
	Ver int `json:"ver,omitempty"`
}

// UserFromToken kiểm JWT và trả tài khoản tương ứng (tài khoản đã bị xoá → token hết hiệu lực).
func (s *Service) UserFromToken(token string) (store.User, error) {
	c := claims{}
	_, err := jwt.ParseWithClaims(token, &c, func(t *jwt.Token) (any, error) { return s.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired(), jwt.WithTimeFunc(s.now))
	if err != nil {
		return store.User{}, ErrInvalidToken
	}
	u, err := s.users.FindByUsername(c.Subject)
	if err != nil {
		return store.User{}, ErrInvalidToken
	}
	// AUTH-BR07: mật khẩu đã được đặt lại sau khi cấp token → phiên cũ hết hiệu lực.
	if c.Ver != u.TokenVersion {
		return store.User{}, ErrInvalidToken
	}
	return u, nil
}

func (s *Service) issue(u store.User) (string, error) {
	now := s.now()
	c := claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   u.Username,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
		},
		Ver: u.TokenVersion,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(s.secret)
}

// ── AUTH-F02 Quên mật khẩu ──────────────────────────────────────────────────

// VerifyForReset xác minh tên đăng nhập + email (AUTH-F02-US01), mở phiên đặt lại và trả token của phiên.
func (s *Service) VerifyForReset(username, email string) (string, store.User, error) {
	// AUTH-BR01, AUTH-BR05: bỏ khoảng trắng đầu / cuối; ô chỉ có khoảng trắng tính là trống (US01-A1).
	username, email = strings.TrimSpace(username), strings.TrimSpace(email)
	if username == "" || email == "" {
		return "", store.User{}, ErrInvalidInput
	}
	// AUTH-BR05 + AUTH-BR02: tên và email phải thuộc cùng một tài khoản có email; mọi trường hợp sai trả cùng một lỗi.
	u, err := s.users.FindByUsername(username)
	if err != nil || u.Email == "" || !strings.EqualFold(u.Email, email) {
		return "", store.User{}, ErrVerifyFailed
	}
	token, err := newResetToken()
	if err != nil {
		return "", store.User{}, err
	}
	// AUTH-BR06: phiên hạn 10 phút; xác minh lại thì phiên mới thay phiên cũ (store huỷ phiên cũ).
	s.users.CreateReset(store.PasswordReset{Token: token, UserID: u.ID, Username: u.Username, ExpiresAt: s.now().Add(ResetTTL)})
	return token, u, nil
}

// CheckReset trả tài khoản của phiên đặt lại nếu phiên còn hiệu lực (dùng khi mở AUTH-F02-SCR02).
func (s *Service) CheckReset(token string) (store.User, error) {
	_, u, err := s.activeReset(token)
	return u, err
}

// ResetPassword đặt mật khẩu mới (AUTH-F02-US02). Thứ tự kiểm: phiên (A4) trước mọi lỗi nhập (A1, A2, A5, A3).
func (s *Service) ResetPassword(token, password, confirm string) error {
	r, u, err := s.activeReset(token)
	if err != nil {
		return err
	}
	if password == "" || confirm == "" { // AUTH-F02-US02-A1 (mật khẩu không bỏ khoảng trắng: khoảng trắng là ký tự)
		return ErrPasswordEmpty
	}
	// AUTH-BR04: 8–16 ký tự, đếm theo ký tự (không phải byte).
	switch n := utf8.RuneCountInString(password); {
	case n < PasswordMinLen:
		return ErrPasswordTooShort // AUTH-F02-US02-A2
	case n > PasswordMaxLen:
		return ErrPasswordTooLong // AUTH-F02-US02-A5
	}
	if password != confirm {
		return ErrPasswordMismatch // AUTH-F02-US02-A3
	}
	hash, err := HashPassword(password)
	if err != nil {
		return err
	}
	// Đóng phiên trước khi đổi mật khẩu: hai request song song cùng token thì chỉ một cái đổi được.
	if !s.users.CloseReset(r.Token, store.ResetUsed) {
		return ErrResetExpired
	}
	// AUTH-BR07: SetPassword tăng TokenVersion → mọi JWT cấp trước đó hết hiệu lực.
	if _, err := s.users.SetPassword(u.ID, u.Username, hash); err != nil {
		return ErrResetExpired
	}
	return nil
}

// CancelReset huỷ phiên đặt lại khi người dùng bỏ ngang ở bước 2 (AUTH-F02-US02-A6). Không có / đã đóng: bỏ qua.
func (s *Service) CancelReset(token string) {
	s.users.CloseReset(token, store.ResetCancelled)
}

// activeReset: AUTH-BR06 — phiên còn hiệu lực khi chưa dùng, chưa huỷ (kể cả bị phiên mới thay), chưa quá hạn,
// và tài khoản vẫn còn (đúng tài khoản lúc xác minh).
func (s *Service) activeReset(token string) (store.PasswordReset, store.User, error) {
	if token == "" {
		return store.PasswordReset{}, store.User{}, ErrResetExpired
	}
	r, err := s.users.FindReset(token)
	if err != nil || r.Status != store.ResetActive || !s.now().Before(r.ExpiresAt) {
		return store.PasswordReset{}, store.User{}, ErrResetExpired
	}
	u, err := s.users.FindByUsername(r.Username)
	if err != nil || u.ID != r.UserID {
		return store.PasswordReset{}, store.User{}, ErrResetExpired
	}
	return r, u, nil
}

func newResetToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
