package httpapi

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"project1/be/internal/auth"
	"project1/be/internal/store"
)

const (
	msgVerifyBlank   = "Nhập tên đăng nhập và email"
	msgVerifyFailed  = "Tên đăng nhập hoặc email không đúng"
	msgResetBlank    = "Nhập mật khẩu mới và nhập lại mật khẩu"
	msgTooShort      = "Mật khẩu phải có ít nhất 8 ký tự"
	msgTooLong       = "Mật khẩu tối đa 16 ký tự"
	msgMismatch      = "Mật khẩu nhập lại không khớp"
	msgExpired       = "Phiên xác minh đã hết hạn, vui lòng xác minh lại"
	testOldPassword  = "Test@1234"
	testNewPassword  = "Moi@12345"
	testUserEmail    = "an@project1.local"
	testUsername     = "an"
	testNoEmailUser  = "khong-email"
	testNoEmailInput = "khong-email@project1.local"
)

func init() {
	// Handler ghi log mỗi request; tắt để output test gọn.
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
}

// newTestServer: store có "an" (email an@project1.local) và "khong-email" (chưa có email), cùng mật khẩu Test@1234.
func newTestServer(t *testing.T, testAPI bool) http.Handler {
	t.Helper()
	hash, err := auth.HashPassword(testOldPassword)
	if err != nil {
		t.Fatal(err)
	}
	st := store.New([]store.User{
		{Username: testUsername, Name: "An", Email: testUserEmail, PasswordHash: hash},
		{Username: testNoEmailUser, Name: "Không email", PasswordHash: hash},
	})
	api := &API{Auth: auth.NewService(st, "test-secret", time.Hour), Users: st, TestAPI: testAPI}
	return api.Routes()
}

type resp struct {
	Status int
	Body   []byte
}

func (r resp) json(t *testing.T) map[string]any {
	t.Helper()
	m := map[string]any{}
	if err := json.Unmarshal(r.Body, &m); err != nil {
		t.Fatalf("body không phải JSON: %q (%v)", r.Body, err)
	}
	return m
}

// wantError kiểm status + {"error": {"code", "message"}}.
func (r resp) wantError(t *testing.T, label string, status int, code, message string) {
	t.Helper()
	if r.Status != status {
		t.Errorf("%s: status = %d; muốn %d (body %s)", label, r.Status, status, r.Body)
		return
	}
	e, _ := r.json(t)["error"].(map[string]any)
	if e["code"] != code || e["message"] != message {
		t.Errorf("%s: error = %v; muốn code %q, message %q", label, e, code, message)
	}
}

func do(t *testing.T, h http.Handler, method, path, token string, body any) resp {
	t.Helper()
	var rd io.Reader
	switch b := body.(type) {
	case nil:
	case string:
		rd = strings.NewReader(b)
	default:
		buf, err := json.Marshal(b)
		if err != nil {
			t.Fatal(err)
		}
		rd = bytes.NewReader(buf)
	}
	req := httptest.NewRequest(method, path, rd)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return resp{Status: rec.Code, Body: rec.Body.Bytes()}
}

func verify(t *testing.T, h http.Handler, username, email string) resp {
	t.Helper()
	return do(t, h, http.MethodPost, "/api/auth/forgot-password/verify", "", map[string]string{"username": username, "email": email})
}

func mustResetToken(t *testing.T, h http.Handler) string {
	t.Helper()
	r := verify(t, h, testUsername, testUserEmail)
	tok, _ := r.json(t)["resetToken"].(string)
	if r.Status != http.StatusOK || tok == "" {
		t.Fatalf("verify: %d %s", r.Status, r.Body)
	}
	return tok
}

func reset(t *testing.T, h http.Handler, token, password, confirm string) resp {
	t.Helper()
	return do(t, h, http.MethodPost, "/api/auth/forgot-password/reset", "",
		map[string]string{"resetToken": token, "password": password, "confirm": confirm})
}

func check(t *testing.T, h http.Handler, token string) resp {
	t.Helper()
	return do(t, h, http.MethodPost, "/api/auth/forgot-password/check", "", map[string]string{"resetToken": token})
}

func login(t *testing.T, h http.Handler, username, password string) resp {
	t.Helper()
	return do(t, h, http.MethodPost, "/api/auth/login", "", map[string]string{"username": username, "password": password})
}

func mustLoginToken(t *testing.T, h http.Handler, username, password string) string {
	t.Helper()
	r := login(t, h, username, password)
	tok, _ := r.json(t)["token"].(string)
	if r.Status != http.StatusOK || tok == "" {
		t.Fatalf("login(%q): %d %s", username, r.Status, r.Body)
	}
	return tok
}

// ── AUTH-F02-US01 Xác minh tài khoản ─────────────────────────────────────────

// AUTH-F02-US01, AUTH-BR01, AUTH-BR05: 200 {resetToken, username}; tên / email không phân biệt hoa thường, bỏ
// khoảng trắng đầu / cuối; username trả về là tên của tài khoản (hiện "Tài khoản: <tên>" ở AUTH-F02-SCR02).
func TestForgotVerify_OK(t *testing.T) {
	h := newTestServer(t, false)
	r := verify(t, h, "  AN ", " An@Project1.LOCAL  ")
	if r.Status != http.StatusOK {
		t.Fatalf("status = %d; muốn 200 (body %s)", r.Status, r.Body)
	}
	m := r.json(t)
	if tok, _ := m["resetToken"].(string); len(tok) < 32 {
		t.Errorf("resetToken = %q; muốn mã ngẫu nhiên đủ dài", tok)
	}
	if m["username"] != testUsername {
		t.Errorf("username = %v; muốn %q", m["username"], testUsername)
	}
	// Hậu điều kiện US01: mật khẩu chưa đổi.
	mustLoginToken(t, h, testUsername, testOldPassword)
}

// AUTH-F02-US01: mỗi lần xác minh mở phiên mới với mã khác nhau (không đoán được).
func TestForgotVerify_TokenUnique(t *testing.T) {
	h := newTestServer(t, false)
	a, b := mustResetToken(t, h), mustResetToken(t, h)
	if a == b {
		t.Errorf("hai lần xác minh trả cùng mã %q", a)
	}
}

// AUTH-F02-US01-A1: tên đăng nhập hoặc email trống / chỉ khoảng trắng → 400 INVALID_INPUT "Nhập tên đăng nhập và email".
// Trống được kiểm trước khớp (tên không tồn tại + email trống vẫn là A1).
func TestForgotVerify_Blank(t *testing.T) {
	h := newTestServer(t, false)
	for _, in := range [][2]string{
		{"", ""}, {"", testUserEmail}, {testUsername, ""}, {"   ", testUserEmail}, {testUsername, " \t "},
		{"khong-ton-tai", ""},
	} {
		verify(t, h, in[0], in[1]).wantError(t, "verify("+in[0]+", "+in[1]+")", http.StatusBadRequest, "INVALID_INPUT", msgVerifyBlank)
	}
	// Thiếu hẳn trường trong JSON cũng là trống.
	do(t, h, http.MethodPost, "/api/auth/forgot-password/verify", "", map[string]string{}).
		wantError(t, "body {}", http.StatusBadRequest, "INVALID_INPUT", msgVerifyBlank)
}

// AUTH-F02-US01-A2, AUTH-BR02, AUTH-BR05: tên không tồn tại / email sai / email của tài khoản khác / tài khoản chưa có
// email → cùng một phản hồi (status, code, message giống hệt nhau), không lộ tài khoản nào tồn tại.
func TestForgotVerify_MismatchSameResponse(t *testing.T) {
	h := newTestServer(t, true)
	// Tài khoản thứ hai có email, để thử email đúng nhưng thuộc tài khoản khác.
	if r := do(t, h, http.MethodPost, "/api/test/users", "", map[string]string{
		"username": "binh", "email": "binh@project1.local", "password": "Binh@1234",
	}); r.Status != http.StatusCreated {
		t.Fatalf("tạo binh: %d %s", r.Status, r.Body)
	}
	cases := [][2]string{
		{"khong-ton-tai", testUserEmail},      // tên không tồn tại
		{testUsername, "sai@project1.local"},  // email không đúng
		{testUsername, "binh@project1.local"}, // email của tài khoản khác
		{testNoEmailUser, testNoEmailInput},   // tài khoản chưa có email
		{testUsername, "an@project1.local.x"}, // gần giống
	}
	var first []byte
	for i, in := range cases {
		r := verify(t, h, in[0], in[1])
		r.wantError(t, "verify("+in[0]+", "+in[1]+")", http.StatusBadRequest, "VERIFY_FAILED", msgVerifyFailed)
		if i == 0 {
			first = r.Body
		} else if !bytes.Equal(r.Body, first) {
			t.Errorf("AUTH-BR02: phản hồi khác nhau giữa các trường hợp sai: %s vs %s", first, r.Body)
		}
		if bytes.Contains(r.Body, []byte("resetToken")) {
			t.Errorf("verify(%q, %q) sai mà vẫn trả resetToken", in[0], in[1])
		}
	}
}

// AUTH-F02-US01: body không phải JSON → 400 (không 500, không mở phiên).
func TestForgotVerify_BadJSON(t *testing.T) {
	h := newTestServer(t, false)
	r := do(t, h, http.MethodPost, "/api/auth/forgot-password/verify", "", "{không phải json")
	if r.Status != http.StatusBadRequest {
		t.Errorf("status = %d; muốn 400", r.Status)
	}
}

// ── AUTH-F02-SCR02 mở trang: check ───────────────────────────────────────────

// AUTH-F02-SCR02 (quy tắc hiển thị), AUTH-BR06: phiên còn hiệu lực → 200 {username}; mã rỗng / không tồn tại /
// body sai → 400 RESET_EXPIRED kèm câu "Phiên xác minh đã hết hạn, vui lòng xác minh lại".
func TestForgotCheck(t *testing.T) {
	h := newTestServer(t, false)
	tok := mustResetToken(t, h)
	r := check(t, h, tok)
	if r.Status != http.StatusOK || r.json(t)["username"] != testUsername {
		t.Errorf("check phiên còn hiệu lực: %d %s; muốn 200 {username: an}", r.Status, r.Body)
	}
	for _, bad := range []string{"", "khong-co", tok + "x"} {
		check(t, h, bad).wantError(t, "check("+bad+")", http.StatusBadRequest, "RESET_EXPIRED", msgExpired)
	}
	do(t, h, http.MethodPost, "/api/auth/forgot-password/check", "", "rác").
		wantError(t, "check body sai", http.StatusBadRequest, "RESET_EXPIRED", msgExpired)
}

// ── AUTH-F02-US02 Đặt mật khẩu mới ───────────────────────────────────────────

// AUTH-F02-US02, AUTH-BR06, AUTH-BR07: 204; mật khẩu cũ hết dùng, mật khẩu mới đăng nhập được; phiên đặt lại đã đóng.
func TestForgotReset_OK(t *testing.T) {
	h := newTestServer(t, false)
	tok := mustResetToken(t, h)
	if r := reset(t, h, tok, testNewPassword, testNewPassword); r.Status != http.StatusNoContent {
		t.Fatalf("reset: %d %s; muốn 204", r.Status, r.Body)
	}
	login(t, h, testUsername, testOldPassword).
		wantError(t, "đăng nhập mật khẩu cũ", http.StatusUnauthorized, "INVALID_CREDENTIALS", "Sai tên đăng nhập hoặc mật khẩu")
	mustLoginToken(t, h, testUsername, testNewPassword)
	// Phiên đặt lại dùng một lần.
	check(t, h, tok).wantError(t, "check sau khi dùng", http.StatusBadRequest, "RESET_EXPIRED", msgExpired)
	reset(t, h, tok, "Khac@12345", "Khac@12345").
		wantError(t, "reset lần 2", http.StatusBadRequest, "RESET_EXPIRED", msgExpired)
}

// AUTH-F02-US02-A1, A2, A5, A3, AUTH-BR04: mã lỗi + câu thông báo; lỗi nhập không đổi mật khẩu, không đóng phiên.
// Thứ tự: trống (bước 4) → độ dài (bước 4) → khớp (bước 5).
func TestForgotReset_InputErrors(t *testing.T) {
	cases := []struct {
		name, pw, confirm, code, msg string
	}{
		{"A1 cả hai trống", "", "", "INVALID_INPUT", msgResetBlank},
		{"A1 mật khẩu trống", "", testNewPassword, "INVALID_INPUT", msgResetBlank},
		{"A1 nhập lại trống", testNewPassword, "", "INVALID_INPUT", msgResetBlank},
		{"A1 trước A2", "Ab@1", "", "INVALID_INPUT", msgResetBlank},
		{"A2 7 ký tự", "Ab@1234", "Ab@1234", "PASSWORD_TOO_SHORT", msgTooShort},
		{"A2 trước A3", "Ab@1234", "khac", "PASSWORD_TOO_SHORT", msgTooShort},
		{"A2 7 ký tự có dấu (đếm ký tự, không đếm byte)", "ậậậậậậậ", "ậậậậậậậ", "PASSWORD_TOO_SHORT", msgTooShort},
		{"A5 17 ký tự", "Abcdefgh@12345678", "Abcdefgh@12345678", "PASSWORD_TOO_LONG", msgTooLong},
		{"A5 trước A3", "Abcdefgh@12345678", "khac", "PASSWORD_TOO_LONG", msgTooLong},
		{"A3 không khớp", testNewPassword, "Moi@12346", "PASSWORD_MISMATCH", msgMismatch},
		{"A3 khác khoảng trắng cuối", testNewPassword, testNewPassword + " ", "PASSWORD_MISMATCH", msgMismatch},
	}
	h := newTestServer(t, false)
	tok := mustResetToken(t, h)
	for _, c := range cases {
		reset(t, h, tok, c.pw, c.confirm).wantError(t, c.name, http.StatusBadRequest, c.code, c.msg)
	}
	// Hậu điều kiện thất bại: mật khẩu giữ nguyên, phiên đặt lại vẫn dùng được.
	mustLoginToken(t, h, testUsername, testOldPassword)
	if r := check(t, h, tok); r.Status != http.StatusOK {
		t.Errorf("sau lỗi nhập, check = %d; muốn phiên còn hiệu lực", r.Status)
	}
}

// AUTH-BR04: biên 8 và 16 ký tự (có dấu, khoảng trắng) hợp lệ; được trùng mật khẩu cũ.
func TestForgotReset_BoundaryLengths(t *testing.T) {
	for _, pw := range []string{"Abcd@123", "Abcdefgh@1234567", strings.Repeat("ậ", 16), "  pass  ", testOldPassword} {
		h := newTestServer(t, false)
		tok := mustResetToken(t, h)
		if r := reset(t, h, tok, pw, pw); r.Status != http.StatusNoContent {
			t.Errorf("reset(%q): %d %s; muốn 204", pw, r.Status, r.Body)
			continue
		}
		mustLoginToken(t, h, testUsername, pw)
	}
}

// AUTH-F02-US02-A4, AUTH-BR06: phiên hết hiệu lực được kiểm TRƯỚC mọi lỗi nhập (A1, A2, A3, A5); mật khẩu không đổi.
func TestForgotReset_ExpiredCheckedFirst(t *testing.T) {
	inputs := [][2]string{
		{"", ""},               // A1
		{"Ab@1234", "Ab@1234"}, // A2
		{"Abcdefgh@12345678", "Abcdefgh@12345678"}, // A5
		{testNewPassword, "khac@12345"},            // A3
		{testNewPassword, testNewPassword},         // hợp lệ
	}
	scenarios := map[string]func(t *testing.T, h http.Handler) string{
		"không có mã": func(t *testing.T, h http.Handler) string { return "" },
		"mã không tồn tại": func(t *testing.T, h http.Handler) string {
			return "khong-ton-tai"
		},
		"đã dùng": func(t *testing.T, h http.Handler) string {
			tok := mustResetToken(t, h)
			if r := reset(t, h, tok, "Dung@1234", "Dung@1234"); r.Status != http.StatusNoContent {
				t.Fatalf("reset: %d", r.Status)
			}
			return tok
		},
		"xác minh lại (phiên mới thay phiên cũ)": func(t *testing.T, h http.Handler) string {
			old := mustResetToken(t, h)
			mustResetToken(t, h)
			return old
		},
		"đã huỷ (A6)": func(t *testing.T, h http.Handler) string {
			tok := mustResetToken(t, h)
			do(t, h, http.MethodPost, "/api/auth/forgot-password/cancel", "", map[string]string{"resetToken": tok})
			return tok
		},
	}
	for name, setup := range scenarios {
		t.Run(name, func(t *testing.T) {
			h := newTestServer(t, false)
			tok := setup(t, h)
			for _, in := range inputs {
				reset(t, h, tok, in[0], in[1]).
					wantError(t, "reset("+in[0]+", "+in[1]+")", http.StatusBadRequest, "RESET_EXPIRED", msgExpired)
			}
		})
	}
}

// AUTH-F02-US02-A4, AUTH-BR06: tài khoản bị xoá sau khi xác minh → phiên hết hiệu lực (qua API test).
func TestForgotReset_AccountDeleted(t *testing.T) {
	h := newTestServer(t, true)
	tok := mustResetToken(t, h)
	if r := do(t, h, http.MethodDelete, "/api/test/users/"+testUsername, "", nil); r.Status != http.StatusNoContent {
		t.Fatalf("xoá: %d", r.Status)
	}
	check(t, h, tok).wantError(t, "check", http.StatusBadRequest, "RESET_EXPIRED", msgExpired)
	reset(t, h, tok, "", "").wantError(t, "reset trống", http.StatusBadRequest, "RESET_EXPIRED", msgExpired)
}

// AUTH-F02-US02-A6, AUTH-BR06: cancel → 204 (kể cả mã không tồn tại / rỗng); phiên bị huỷ; mật khẩu không đổi.
func TestForgotCancel(t *testing.T) {
	h := newTestServer(t, false)
	for _, bad := range []string{"", "khong-co"} {
		if r := do(t, h, http.MethodPost, "/api/auth/forgot-password/cancel", "", map[string]string{"resetToken": bad}); r.Status != http.StatusNoContent {
			t.Errorf("cancel(%q) = %d; muốn 204", bad, r.Status)
		}
	}
	tok := mustResetToken(t, h)
	if r := do(t, h, http.MethodPost, "/api/auth/forgot-password/cancel", "", map[string]string{"resetToken": tok}); r.Status != http.StatusNoContent {
		t.Fatalf("cancel = %d; muốn 204", r.Status)
	}
	check(t, h, tok).wantError(t, "check sau huỷ", http.StatusBadRequest, "RESET_EXPIRED", msgExpired)
	mustLoginToken(t, h, testUsername, testOldPassword)
	// Huỷ lần 2 vẫn 204.
	if r := do(t, h, http.MethodPost, "/api/auth/forgot-password/cancel", "", map[string]string{"resetToken": tok}); r.Status != http.StatusNoContent {
		t.Errorf("cancel lần 2 = %d; muốn 204", r.Status)
	}
}

// AUTH-BR06: huỷ / dùng phiên của tài khoản này không ảnh hưởng phiên của tài khoản khác.
func TestForgotReset_SessionsPerAccount(t *testing.T) {
	h := newTestServer(t, true)
	do(t, h, http.MethodPost, "/api/test/users", "", map[string]string{"username": "binh", "email": "binh@project1.local", "password": "Binh@1234"})
	tokAn := mustResetToken(t, h)
	r := verify(t, h, "binh", "binh@project1.local")
	tokBinh, _ := r.json(t)["resetToken"].(string)
	if tokBinh == "" {
		t.Fatalf("verify binh: %d %s", r.Status, r.Body)
	}
	if r := check(t, h, tokAn); r.Status != http.StatusOK {
		t.Errorf("phiên an bị huỷ khi binh xác minh: %d", r.Status)
	}
	if r := reset(t, h, tokBinh, testNewPassword, testNewPassword); r.Status != http.StatusNoContent {
		t.Fatalf("reset binh: %d %s", r.Status, r.Body)
	}
	// Đặt lại của binh chỉ đổi mật khẩu binh.
	mustLoginToken(t, h, testUsername, testOldPassword)
	mustLoginToken(t, h, "binh", testNewPassword)
	if r := check(t, h, tokAn); r.Status != http.StatusOK || r.json(t)["username"] != testUsername {
		t.Errorf("phiên an sau khi binh đặt lại: %d %s", r.Status, r.Body)
	}
}

// AUTH-BR07, AUTH-F01-US02-A2: đặt lại mật khẩu → JWT cấp trước đó bị /api/me từ chối 401; JWT cấp sau dùng được;
// phiên của tài khoản khác không bị ảnh hưởng.
func TestForgotReset_InvalidatesLoginSessions(t *testing.T) {
	h := newTestServer(t, false)
	old := mustLoginToken(t, h, testUsername, testOldPassword)
	other := mustLoginToken(t, h, testNoEmailUser, testOldPassword)
	if r := do(t, h, http.MethodGet, "/api/me", old, nil); r.Status != http.StatusOK {
		t.Fatalf("me trước khi đặt lại: %d", r.Status)
	}
	tok := mustResetToken(t, h)
	if r := reset(t, h, tok, testNewPassword, testNewPassword); r.Status != http.StatusNoContent {
		t.Fatalf("reset: %d %s", r.Status, r.Body)
	}
	if r := do(t, h, http.MethodGet, "/api/me", old, nil); r.Status != http.StatusUnauthorized {
		t.Errorf("me bằng JWT cũ sau khi đặt lại = %d; muốn 401", r.Status)
	}
	fresh := mustLoginToken(t, h, testUsername, testNewPassword)
	if r := do(t, h, http.MethodGet, "/api/me", fresh, nil); r.Status != http.StatusOK {
		t.Errorf("me bằng JWT mới = %d; muốn 200", r.Status)
	}
	if r := do(t, h, http.MethodGet, "/api/me", other, nil); r.Status != http.StatusOK {
		t.Errorf("me của tài khoản khác = %d; muốn 200", r.Status)
	}
}

// AUTH-F02: các endpoint quên mật khẩu không cần đăng nhập và chỉ nhận POST.
func TestForgotEndpoints_MethodAndNoAuth(t *testing.T) {
	h := newTestServer(t, false)
	for _, p := range []string{"verify", "check", "reset", "cancel"} {
		if r := do(t, h, http.MethodGet, "/api/auth/forgot-password/"+p, "", nil); r.Status != http.StatusMethodNotAllowed {
			t.Errorf("GET %s = %d; muốn 405", p, r.Status)
		}
	}
}

// ── API test: /api/test/users nhận email (FSD §4) ────────────────────────────

// AUTH-F02: /api/test/users nhận email tuỳ chọn — có email thì xác minh được, không có thì A2; không trả passwordHash.
func TestTestAPI_CreateUserWithEmail(t *testing.T) {
	h := newTestServer(t, true)
	r := do(t, h, http.MethodPost, "/api/test/users", "", map[string]string{"username": "co-email", "email": "co@project1.local", "password": "Co@12345"})
	if r.Status != http.StatusCreated {
		t.Fatalf("tạo co-email: %d %s", r.Status, r.Body)
	}
	if m := r.json(t); m["email"] != "co@project1.local" {
		t.Errorf("email trả về = %v; muốn co@project1.local", m["email"])
	}
	if bytes.Contains(bytes.ToLower(r.Body), []byte("password")) {
		t.Errorf("phản hồi lộ mật khẩu / hash: %s", r.Body)
	}
	if r := verify(t, h, "co-email", "CO@project1.local"); r.Status != http.StatusOK {
		t.Errorf("verify tài khoản test có email: %d %s", r.Status, r.Body)
	}

	if r := do(t, h, http.MethodPost, "/api/test/users", "", map[string]string{"username": "chua-email", "password": "Chua@1234"}); r.Status != http.StatusCreated {
		t.Fatalf("tạo chua-email: %d %s", r.Status, r.Body)
	}
	verify(t, h, "chua-email", "chua-email@project1.local").
		wantError(t, "verify tài khoản chưa có email", http.StatusBadRequest, "VERIFY_FAILED", msgVerifyFailed)

	do(t, h, http.MethodPost, "/api/test/users", "", map[string]string{"username": "CO-EMAIL", "password": "x"}).
		wantError(t, "trùng username", http.StatusConflict, "USERNAME_TAKEN", "Tên đăng nhập đã tồn tại")
}

// API test chỉ đăng ký khi TEST_API=1 (CLAUDE.md).
func TestTestAPI_DisabledByDefault(t *testing.T) {
	h := newTestServer(t, false)
	r := do(t, h, http.MethodPost, "/api/test/users", "", map[string]string{"username": "x", "password": "x"})
	if r.Status == http.StatusCreated {
		t.Errorf("TEST_API tắt mà /api/test/users vẫn tạo được tài khoản")
	}
}

// ── AUTH-F01 hồi quy ─────────────────────────────────────────────────────────

// AUTH-F01-US01, AUTH-BR01: đăng nhập thành công, tên không phân biệt hoa thường, bỏ khoảng trắng; trả token + user,
// không lộ passwordHash.
func TestLogin_OK(t *testing.T) {
	h := newTestServer(t, false)
	r := login(t, h, "  AN ", testOldPassword)
	if r.Status != http.StatusOK {
		t.Fatalf("login: %d %s", r.Status, r.Body)
	}
	m := r.json(t)
	u, _ := m["user"].(map[string]any)
	if m["token"] == "" || u["username"] != testUsername || u["name"] != "An" {
		t.Errorf("login trả %s", r.Body)
	}
	if bytes.Contains(bytes.ToLower(r.Body), []byte("passwordhash")) {
		t.Errorf("phản hồi lộ passwordHash: %s", r.Body)
	}
}

// AUTH-F01-US01-A1: tên trống / chỉ khoảng trắng / mật khẩu trống → 400 INVALID_INPUT.
func TestLogin_Blank(t *testing.T) {
	h := newTestServer(t, false)
	for _, in := range [][2]string{{"", ""}, {"   ", testOldPassword}, {testUsername, ""}} {
		login(t, h, in[0], in[1]).wantError(t, "login("+in[0]+")", http.StatusBadRequest, "INVALID_INPUT", "Nhập tên đăng nhập và mật khẩu")
	}
}

// AUTH-F01-US01-A2, AUTH-BR02: tên không tồn tại và sai mật khẩu trả cùng một phản hồi.
func TestLogin_WrongSameResponse(t *testing.T) {
	h := newTestServer(t, false)
	a := login(t, h, "khong-ton-tai", testOldPassword)
	b := login(t, h, testUsername, "Sai@12345")
	a.wantError(t, "tên không tồn tại", http.StatusUnauthorized, "INVALID_CREDENTIALS", "Sai tên đăng nhập hoặc mật khẩu")
	if !bytes.Equal(a.Body, b.Body) || a.Status != b.Status {
		t.Errorf("AUTH-BR02: phản hồi khác nhau: %d %s vs %d %s", a.Status, a.Body, b.Status, b.Body)
	}
}

// AUTH-F01-US02, AUTH-F01-US02-A1: /api/me cần Bearer token hợp lệ.
func TestMe_RequireAuth(t *testing.T) {
	h := newTestServer(t, false)
	for _, tok := range []string{"", "rac"} {
		if r := do(t, h, http.MethodGet, "/api/me", tok, nil); r.Status != http.StatusUnauthorized {
			t.Errorf("me(token %q) = %d; muốn 401", tok, r.Status)
		}
	}
	tok := mustLoginToken(t, h, testUsername, testOldPassword)
	r := do(t, h, http.MethodGet, "/api/me", tok, nil)
	if r.Status != http.StatusOK || r.json(t)["name"] != "An" {
		t.Errorf("me = %d %s; muốn 200 name An", r.Status, r.Body)
	}
}

// AUTH-F02-US02-A4: body không phải JSON coi như rỗng → phiên (token rỗng) được kiểm trước mọi lỗi nhập.
func TestForgotReset_BadJSONIsExpired(t *testing.T) {
	h := newTestServer(t, false)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/forgot-password/reset", strings.NewReader("{không phải json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	resp{Status: rec.Code, Body: rec.Body.Bytes()}.wantError(t, "body sai", http.StatusBadRequest, "RESET_EXPIRED", msgExpired)
}

// FSD §3.1: email duy nhất nếu có (không phân biệt hoa thường); tài khoản không email thì không bị chặn.
func TestTestAPI_CreateUserEmailUnique(t *testing.T) {
	h := newTestServer(t, true)
	r := do(t, h, http.MethodPost, "/api/test/users", "", map[string]string{"username": "b", "password": "x", "email": "AN@project1.local"})
	r.wantError(t, "trùng email", http.StatusConflict, "EMAIL_TAKEN", "Email đã được dùng cho tài khoản khác")
	for _, u := range []string{"c", "d"} {
		if r := do(t, h, http.MethodPost, "/api/test/users", "", map[string]string{"username": u, "password": "x"}); r.Status != http.StatusCreated {
			t.Errorf("tạo %s không email: status = %d; muốn 201", u, r.Status)
		}
	}
}
