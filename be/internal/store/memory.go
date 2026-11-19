// Package store giữ dữ liệu trên RAM (map + mutex). Khởi động lại server là dữ liệu về như ban đầu.
package store

import (
	"errors"
	"strings"
	"sync"
	"time"
)

// ErrNotFound: không có bản ghi. ErrExists: trùng khoá. ErrEmailExists: trùng email.
var (
	ErrNotFound    = errors.New("not found")
	ErrExists      = errors.New("already exists")
	ErrEmailExists = errors.New("email already exists")
)

// User là tài khoản đăng nhập.
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Name         string    `json:"name"`
	Email        string    `json:"email,omitempty"` // rỗng = tài khoản chưa có email
	PasswordHash string    `json:"-"`
	TokenVersion int       `json:"-"` // tăng mỗi lần đặt lại mật khẩu; JWT mang số cũ thì hết hiệu lực
	CreatedAt    time.Time `json:"createdAt"`
}

// Trạng thái của phiên đặt lại mật khẩu. "Hết hạn" không lưu, suy ra từ ExpiresAt.
const (
	ResetActive    = "con_hieu_luc"
	ResetUsed      = "da_dung"
	ResetCancelled = "da_huy"
)

// PasswordReset là phiên đặt lại mật khẩu, mở ra khi xác minh tên đăng nhập + email thành công.
type PasswordReset struct {
	Token     string
	UserID    int
	Username  string
	ExpiresAt time.Time
	Status    string
}

// Store lưu tài khoản và phiên đặt lại mật khẩu trên RAM, an toàn khi nhiều request chạy song song.
type Store struct {
	mu          sync.RWMutex
	seq         int
	byName      map[string]*User
	resets      map[string]*PasswordReset // theo token
	activeReset map[string]string         // key(username) → token của phiên mới nhất
	seed        []User
}

// New tạo store với các tài khoản mẫu (PasswordHash đã băm sẵn).
func New(seed []User) *Store {
	s := &Store{seed: seed}
	s.Reset()
	return s
}

// Reset đưa dữ liệu về các tài khoản mẫu.
func (s *Store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seq = 0
	s.byName = map[string]*User{}
	s.resets = map[string]*PasswordReset{}
	s.activeReset = map[string]string{}
	for _, u := range s.seed {
		u := u
		s.seq++
		u.ID = s.seq
		if u.CreatedAt.IsZero() {
			u.CreatedAt = time.Now()
		}
		s.byName[key(u.Username)] = &u
	}
}

// FindByUsername tìm tài khoản, không phân biệt hoa thường.
func (s *Store) FindByUsername(username string) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.byName[key(username)]
	if !ok {
		return User{}, ErrNotFound
	}
	return *u, nil
}

// Create thêm tài khoản; trùng username → ErrExists; trùng email (nếu có, không phân biệt hoa thường) → ErrEmailExists.
func (s *Store) Create(u User) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byName[key(u.Username)]; ok {
		return User{}, ErrExists
	}
	if u.Email != "" {
		for _, other := range s.byName {
			if strings.EqualFold(other.Email, u.Email) {
				return User{}, ErrEmailExists
			}
		}
	}
	s.seq++
	u.ID = s.seq
	u.CreatedAt = time.Now()
	s.byName[key(u.Username)] = &u
	return u, nil
}

// Delete xoá tài khoản. Phiên đặt lại của tài khoản không cần đổi trạng thái: service coi phiên của tài khoản
// không còn là hết hiệu lực (AUTH-BR06).
func (s *Store) Delete(username string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byName[key(username)]; !ok {
		return ErrNotFound
	}
	delete(s.byName, key(username))
	return nil
}

// SetPassword đổi mật khẩu (đã băm) của đúng tài khoản id và tăng TokenVersion để các JWT cũ hết hiệu lực.
// Tài khoản không còn (hoặc đã bị xoá rồi tạo lại cùng tên) → ErrNotFound.
func (s *Store) SetPassword(id int, username, hash string) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.byName[key(username)]
	if !ok || u.ID != id {
		return User{}, ErrNotFound
	}
	u.PasswordHash = hash
	u.TokenVersion++
	return *u, nil
}

// CreateReset lưu phiên đặt lại mới; phiên còn hiệu lực trước đó của cùng tài khoản bị huỷ.
func (s *Store) CreateReset(r PasswordReset) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if old, ok := s.resets[s.activeReset[key(r.Username)]]; ok && old.Status == ResetActive {
		old.Status = ResetCancelled
	}
	r.Status = ResetActive
	s.resets[r.Token] = &r
	s.activeReset[key(r.Username)] = r.Token
}

// FindReset tìm phiên đặt lại theo token.
func (s *Store) FindReset(token string) (PasswordReset, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.resets[token]
	if !ok {
		return PasswordReset{}, ErrNotFound
	}
	return *r, nil
}

// CloseReset chuyển phiên đặt lại đang còn hiệu lực sang status (ResetUsed / ResetCancelled).
// Trả false nếu phiên không còn ở trạng thái còn hiệu lực (đã dùng / đã huỷ / không có).
func (s *Store) CloseReset(token, status string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.resets[token]
	if !ok || r.Status != ResetActive {
		return false
	}
	r.Status = status
	return true
}

func key(username string) string { return strings.ToLower(strings.TrimSpace(username)) }
