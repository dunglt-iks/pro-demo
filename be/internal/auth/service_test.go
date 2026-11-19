package auth

import (
	"errors"
	"strings"
	"testing"
	"time"

	"project1/be/internal/store"
)

// newTestService: store có tài khoản "an" (email an@project1.local, mật khẩu Test@1234) và "khong-email" (chưa có
// email); đồng hồ giả bắt đầu ở clock[0], test tự đẩy.
func newTestService(t *testing.T) (*Service, *store.Store, *time.Time) {
	t.Helper()
	hash, err := HashPassword("Test@1234")
	if err != nil {
		t.Fatal(err)
	}
	st := store.New([]store.User{
		{Username: "an", Name: "An", Email: "an@project1.local", PasswordHash: hash},
		{Username: "khong-email", Name: "Không email", PasswordHash: hash},
	})
	svc := NewService(st, "test-secret", time.Hour)
	now := time.Date(2026, 11, 19, 9, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return now }
	return svc, st, &now
}

func mustVerify(t *testing.T, svc *Service) string {
	t.Helper()
	token, _, err := svc.VerifyForReset("an", "an@project1.local")
	if err != nil {
		t.Fatalf("VerifyForReset: %v", err)
	}
	return token
}

func mustLogin(t *testing.T, svc *Service, password string) string {
	t.Helper()
	token, _, err := svc.Login("an", password)
	if err != nil {
		t.Fatalf("Login(%q): %v", password, err)
	}
	return token
}

// AUTH-F02-US01, AUTH-BR01, AUTH-BR05: khớp tên + email, không phân biệt hoa thường, bỏ khoảng trắng đầu / cuối.
func TestVerifyForReset_OK(t *testing.T) {
	svc, _, _ := newTestService(t)
	for _, in := range [][2]string{
		{"an", "an@project1.local"},
		{"AN", "an@project1.local"},
		{"an", "AN@Project1.Local"},
		{"  an  ", "  an@project1.local  "},
	} {
		token, u, err := svc.VerifyForReset(in[0], in[1])
		if err != nil || token == "" || u.Username != "an" {
			t.Errorf("VerifyForReset(%q, %q) = %q, %q, %v; muốn token + tài khoản an", in[0], in[1], token, u.Username, err)
		}
	}
}

// AUTH-F02-US01-A1: trống hoặc chỉ khoảng trắng.
func TestVerifyForReset_Blank(t *testing.T) {
	svc, _, _ := newTestService(t)
	for _, in := range [][2]string{{"", ""}, {"", "an@project1.local"}, {"an", ""}, {"   ", "an@project1.local"}, {"an", "   "}} {
		if _, _, err := svc.VerifyForReset(in[0], in[1]); !errors.Is(err, ErrInvalidInput) {
			t.Errorf("VerifyForReset(%q, %q) err = %v; muốn ErrInvalidInput", in[0], in[1], err)
		}
	}
}

// AUTH-F02-US01-A2, AUTH-BR02: mọi trường hợp không khớp trả cùng một lỗi.
func TestVerifyForReset_Mismatch(t *testing.T) {
	svc, _, _ := newTestService(t)
	for _, in := range [][2]string{
		{"an", "sai@project1.local"},
		{"khong-ton-tai", "an@project1.local"},
		{"khong-email", "k@project1.local"},
	} {
		if _, _, err := svc.VerifyForReset(in[0], in[1]); !errors.Is(err, ErrVerifyFailed) {
			t.Errorf("VerifyForReset(%q, %q) err = %v; muốn ErrVerifyFailed", in[0], in[1], err)
		}
	}
}

// AUTH-F02-US02: đặt mật khẩu mới; mật khẩu cũ hết dùng được, mật khẩu mới đăng nhập được.
func TestResetPassword_OK(t *testing.T) {
	svc, _, _ := newTestService(t)
	token := mustVerify(t, svc)
	if err := svc.ResetPassword(token, "Moi@12345", "Moi@12345"); err != nil {
		t.Fatalf("ResetPassword: %v", err)
	}
	if _, _, err := svc.Login("an", "Test@1234"); !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("đăng nhập mật khẩu cũ err = %v; muốn ErrInvalidCredentials", err)
	}
	mustLogin(t, svc, "Moi@12345")
}

// AUTH-BR04: 8–16 ký tự, đếm theo ký tự (có dấu, khoảng trắng), được trùng mật khẩu cũ.
func TestResetPassword_ValidLengths(t *testing.T) {
	for _, pw := range []string{"Abcd@123", "Abcdefgh@1234567", "Mật khẩu 12", "Test@1234", strings.Repeat("ậ", 16), strings.Repeat(" ", 8)} {
		svc, _, _ := newTestService(t)
		token := mustVerify(t, svc)
		if err := svc.ResetPassword(token, pw, pw); err != nil {
			t.Errorf("ResetPassword(%q) err = %v; muốn thành công", pw, err)
			continue
		}
		mustLogin(t, svc, pw)
	}
}

// AUTH-F02-US02-A1, A2, A5, A3: lỗi nhập, phiên giữ nguyên để nhập lại.
func TestResetPassword_InputErrors(t *testing.T) {
	cases := []struct {
		name, pw, confirm string
		want              error
	}{
		{"A1 cả hai trống", "", "", ErrPasswordEmpty},
		{"A1 nhập lại trống", "Moi@12345", "", ErrPasswordEmpty},
		{"A2 7 ký tự", "Ab@1234", "Ab@1234", ErrPasswordTooShort},
		{"A2 7 khoảng trắng", strings.Repeat(" ", 7), strings.Repeat(" ", 7), ErrPasswordTooShort},
		{"A5 17 ký tự", "Abcdefgh@12345678", "Abcdefgh@12345678", ErrPasswordTooLong},
		{"A5 17 ký tự có dấu", strings.Repeat("ậ", 17), "x", ErrPasswordTooLong},
		{"A3 không khớp", "Moi@12345", "Moi@12346", ErrPasswordMismatch},
	}
	svc, _, _ := newTestService(t)
	token := mustVerify(t, svc)
	for _, c := range cases {
		if err := svc.ResetPassword(token, c.pw, c.confirm); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v; muốn %v", c.name, err, c.want)
		}
	}
	// Phiên vẫn còn hiệu lực sau các lỗi nhập.
	if err := svc.ResetPassword(token, "Moi@12345", "Moi@12345"); err != nil {
		t.Errorf("sau lỗi nhập, ResetPassword err = %v; muốn thành công", err)
	}
}

// AUTH-F02-US02-A4, AUTH-BR06: các trường hợp phiên hết hiệu lực; kiểm trước lỗi nhập.
func TestResetPassword_Expired(t *testing.T) {
	t.Run("quá 10 phút", func(t *testing.T) {
		svc, _, now := newTestService(t)
		token := mustVerify(t, svc)
		*now = now.Add(ResetTTL - time.Second)
		if _, err := svc.CheckReset(token); err != nil {
			t.Fatalf("trước 10 phút CheckReset err = %v; muốn còn hiệu lực", err)
		}
		*now = now.Add(time.Second)
		if err := svc.ResetPassword(token, "Ab@1", "khac"); !errors.Is(err, ErrResetExpired) {
			t.Errorf("đúng 10 phút, err = %v; muốn ErrResetExpired (trước lỗi nhập)", err)
		}
	})
	t.Run("đã dùng", func(t *testing.T) {
		svc, _, _ := newTestService(t)
		token := mustVerify(t, svc)
		if err := svc.ResetPassword(token, "Moi@12345", "Moi@12345"); err != nil {
			t.Fatal(err)
		}
		if err := svc.ResetPassword(token, "Khac@12345", "Khac@12345"); !errors.Is(err, ErrResetExpired) {
			t.Errorf("dùng lần 2 err = %v; muốn ErrResetExpired", err)
		}
	})
	t.Run("xác minh lại thay phiên cũ", func(t *testing.T) {
		svc, _, _ := newTestService(t)
		old := mustVerify(t, svc)
		fresh := mustVerify(t, svc)
		if err := svc.ResetPassword(old, "", ""); !errors.Is(err, ErrResetExpired) {
			t.Errorf("phiên cũ err = %v; muốn ErrResetExpired (trước lỗi trống)", err)
		}
		if _, err := svc.CheckReset(fresh); err != nil {
			t.Errorf("phiên mới CheckReset err = %v; muốn còn hiệu lực", err)
		}
	})
	t.Run("huỷ (A6)", func(t *testing.T) {
		svc, _, _ := newTestService(t)
		token := mustVerify(t, svc)
		svc.CancelReset(token)
		if _, err := svc.CheckReset(token); !errors.Is(err, ErrResetExpired) {
			t.Errorf("sau huỷ CheckReset err = %v; muốn ErrResetExpired", err)
		}
	})
	t.Run("tài khoản bị xoá", func(t *testing.T) {
		svc, st, _ := newTestService(t)
		token := mustVerify(t, svc)
		if err := st.Delete("an"); err != nil {
			t.Fatal(err)
		}
		if err := svc.ResetPassword(token, "Moi@12345", "Moi@12345"); !errors.Is(err, ErrResetExpired) {
			t.Errorf("err = %v; muốn ErrResetExpired", err)
		}
	})
	t.Run("tài khoản bị xoá rồi tạo lại cùng tên", func(t *testing.T) {
		svc, st, _ := newTestService(t)
		token := mustVerify(t, svc)
		_ = st.Delete("an")
		if _, err := st.Create(store.User{Username: "an", Email: "an@project1.local"}); err != nil {
			t.Fatal(err)
		}
		if _, err := svc.CheckReset(token); !errors.Is(err, ErrResetExpired) {
			t.Errorf("err = %v; muốn ErrResetExpired (khác tài khoản)", err)
		}
	})
	t.Run("token không tồn tại / rỗng", func(t *testing.T) {
		svc, _, _ := newTestService(t)
		for _, tok := range []string{"", "khong-co"} {
			if _, err := svc.CheckReset(tok); !errors.Is(err, ErrResetExpired) {
				t.Errorf("CheckReset(%q) err = %v; muốn ErrResetExpired", tok, err)
			}
		}
	})
}

// AUTH-BR07: đặt lại mật khẩu làm phiên đăng nhập cũ hết hiệu lực; phiên mới cấp sau đó vẫn dùng được.
func TestResetPassword_InvalidatesSessions(t *testing.T) {
	svc, _, _ := newTestService(t)
	old := mustLogin(t, svc, "Test@1234")
	if _, err := svc.UserFromToken(old); err != nil {
		t.Fatalf("token cũ trước khi đặt lại: %v", err)
	}
	token := mustVerify(t, svc)
	if err := svc.ResetPassword(token, "Moi@12345", "Moi@12345"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.UserFromToken(old); !errors.Is(err, ErrInvalidToken) {
		t.Errorf("token cũ sau khi đặt lại err = %v; muốn ErrInvalidToken", err)
	}
	fresh := mustLogin(t, svc, "Moi@12345")
	if _, err := svc.UserFromToken(fresh); err != nil {
		t.Errorf("token mới err = %v; muốn hợp lệ", err)
	}
}
