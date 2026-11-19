// Package store giữ dữ liệu trên RAM (map + mutex). Khởi động lại server là dữ liệu về như ban đầu.
package store

import (
	"errors"
	"strings"
	"sync"
	"time"
)

// ErrNotFound: không có bản ghi. ErrExists: trùng khoá.
var (
	ErrNotFound = errors.New("not found")
	ErrExists   = errors.New("already exists")
)

// User là tài khoản đăng nhập.
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
}

// Store lưu tài khoản trên RAM, an toàn khi nhiều request chạy song song.
type Store struct {
	mu     sync.RWMutex
	seq    int
	byName map[string]*User
	seed   []User
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

// Create thêm tài khoản; trùng username → ErrExists.
func (s *Store) Create(u User) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byName[key(u.Username)]; ok {
		return User{}, ErrExists
	}
	s.seq++
	u.ID = s.seq
	u.CreatedAt = time.Now()
	s.byName[key(u.Username)] = &u
	return u, nil
}

// Delete xoá tài khoản.
func (s *Store) Delete(username string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byName[key(username)]; !ok {
		return ErrNotFound
	}
	delete(s.byName, key(username))
	return nil
}

func key(username string) string { return strings.ToLower(strings.TrimSpace(username)) }
