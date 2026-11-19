Spec đã triển khai: docs/features/AUTH-xac-thuc/FSD-xac-thuc.md v1.0

# Kế hoạch — AUTH-F02 Quên mật khẩu

- **Spec:** [FSD §2 AUTH-F02](../features/AUTH-xac-thuc/FSD-xac-thuc.md#2-auth-f02--quên-mật-khẩu), BR01–BR07 ở
  [PRD §5](../features/AUTH-xac-thuc/PRD-xac-thuc.md#5-quy-tắc-nghiệp-vụ) — trạng thái: đã duyệt
- **Ước lượng:** vừa (một PR)
- **Phạm vi không làm:** gửi email / link / OTP; giới hạn số lần xác minh sai; màn nhập / sửa email; đổi mật khẩu khi
  đang đăng nhập (PRD §1).
- **Đụng AUTH-F01-SCR01** (đã duyệt cùng đợt): thêm link "Quên mật khẩu?" (7), thông báo thành công (2), và quy tắc
  "có phiên còn hiệu lực → về trang chủ". Đây là phần F01 mà spec F02 kéo theo, không thêm gì ngoài spec.

## 1. Tiêu chí chấp nhận

| Mã | Tiêu chí | Kiểm bằng |
|---|---|---|
| AUTH-F02-US01 | Tên + email khớp cùng tài khoản → sang bước 2 | Go unit + e2e TC-001, 002 |
| AUTH-F02-US01-A1 | Trống / chỉ khoảng trắng → "Nhập tên đăng nhập và email" | Go unit + e2e TC-003 |
| AUTH-F02-US01-A2 | Không khớp / không có email → "Tên đăng nhập hoặc email không đúng", không lộ tài khoản | Go unit + e2e TC-004, 022 |
| AUTH-F02-US01-E1, E2 | Mất kết nối / lỗi hệ thống → câu tương ứng | e2e TC-005, 006 |
| AUTH-F02-US02 | Lưu mật khẩu mới, về đăng nhập + thông báo thành công | Go unit + e2e TC-001, 011 |
| AUTH-F02-US02-A1, A2, A3, A5 | Trống / < 8 / không khớp / > 16 ký tự | Go unit + e2e TC-008–012 |
| AUTH-F02-US02-A4 | Phiên hết hiệu lực (quá 10 phút, đã dùng, bị thay, đã huỷ, tài khoản không còn), kiểm **trước** lỗi nhập | Go unit (kể cả 10 phút bằng đồng hồ giả) + e2e TC-013, 014, 020 |
| AUTH-F02-US02-A6 | Về đăng nhập ở bước 2 → huỷ phiên | Go unit + e2e TC-015 |
| AUTH-F02-US02-E1, E2 | Mất kết nối / lỗi hệ thống | e2e TC-016, 017 |
| AUTH-BR07 | Đặt lại xong → phiên đăng nhập cũ hết hiệu lực | Go unit + e2e TC-018 |
| AUTH-F02-SCR01/02 quy tắc hiển thị | Chưa xác minh → về bước 1 không thông báo; đang đăng nhập → trang chủ; nút khi chờ | e2e TC-007, 019, 020, 021 |

## 2. File tạo / sửa (theo thứ tự lớp: store → service → handler → FE)

1. `be/internal/store/memory.go` — `User` thêm `Email`, `TokenVersion`; `SetPassword` (đổi hash, tăng
   `TokenVersion`); kho **phiên đặt lại** trên RAM: tạo (thay phiên cũ của cùng tài khoản), lấy, đánh dấu đã dùng /
   đã huỷ; `Reset()` dọn cả phiên đặt lại.
2. `be/internal/auth/service.go` — JWT thêm claim `ver` (= `TokenVersion`), `UserFromToken` từ chối khi lệch
   (BR07); lỗi mới `ErrVerifyFailed`, `ErrResetExpired`, `ErrPasswordEmpty`, `ErrPasswordTooShort`,
   `ErrPasswordTooLong`, `ErrPasswordMismatch`; hàm `VerifyForReset`, `CheckReset`, `ResetPassword`, `CancelReset`.
   Đồng hồ inject được (`now func() time.Time`) để test hạn 10 phút.
3. `be/internal/auth/service_test.go` (mới) — unit test theo mã (US01/US02 và các nhánh, BR04–BR07).
4. `be/internal/httpapi/handler.go` — 4 route mới (bảng §3); `/api/test/users` nhận thêm `email`.
5. `be/cmd/server/main.go` — email cho tài khoản mẫu (`admin@project1.local`, `user1@project1.local`).
6. `fe/src/lib/api.ts` — `forgotVerify`, `forgotCheck`, `forgotReset`, `forgotCancel`; tuỳ chọn không tự chuyển
   `/login` khi 401 (dùng cho kiểm "đang đăng nhập" ở trang khách).
7. `fe/src/lib/session.ts` — thêm `resetSession` (sessionStorage: mã phiên đặt lại của tab) và cờ thông báo một lần
   (`flash`: thông báo thành công cho /login, thông báo hết hạn cho /forgot-password) — URL giữ sạch, tải lại là mất.
8. `fe/src/lib/useGuestOnly.ts` (mới) — hook: có phiên đăng nhập còn hiệu lực → `router.replace("/")`.
9. `fe/src/components/ForgotPasswordForm.tsx`, `ResetPasswordForm.tsx` (mới); `fe/src/app/forgot-password/page.tsx`,
   `fe/src/app/forgot-password/reset/page.tsx` (mới).
10. `fe/src/components/LoginForm.tsx`, `fe/src/app/login/page.tsx` — link "Quên mật khẩu?", thông báo thành công
    (`data-testid="login-success"`, ẩn khi bấm Đăng nhập), dùng `useGuestOnly`.

## 3. Hợp đồng API

Lỗi theo dạng chung `{"error": {"code", "message"}}`. Mọi endpoint không cần đăng nhập.

| Endpoint | Request | Response | Mã lỗi |
|---|---|---|---|
| `POST /api/auth/forgot-password/verify` | `{username, email}` | 200 `{resetToken, username}` | 400 `INVALID_INPUT` "Nhập tên đăng nhập và email"; 400 `VERIFY_FAILED` "Tên đăng nhập hoặc email không đúng" |
| `POST /api/auth/forgot-password/check` | `{resetToken}` | 200 `{username}` | 400 `RESET_EXPIRED` "Phiên xác minh đã hết hạn, vui lòng xác minh lại" |
| `POST /api/auth/forgot-password/reset` | `{resetToken, password, confirm}` | 204 | theo thứ tự: 400 `RESET_EXPIRED`; 400 `INVALID_INPUT` "Nhập mật khẩu mới và nhập lại mật khẩu"; 400 `PASSWORD_TOO_SHORT` "Mật khẩu phải có ít nhất 8 ký tự"; 400 `PASSWORD_TOO_LONG` "Mật khẩu tối đa 16 ký tự"; 400 `PASSWORD_MISMATCH` "Mật khẩu nhập lại không khớp" |
| `POST /api/auth/forgot-password/cancel` | `{resetToken}` | 204 (kể cả mã không tồn tại) | — |
| `POST /api/test/users` (chỉ `TEST_API=1`) | thêm `email` (tuỳ chọn) | 201 | như cũ |

- So với hợp đồng chốt ở `/e2e-script`: `reset` nhận thêm `confirm`, thêm `check` và `cancel`. Lý do: A4 phải được kiểm
  **trước** lỗi nhập (kể cả trống / không khớp) nên mọi kiểm tra chạy ở BE theo một thứ tự; FE không tự kiểm trước.
  Page object / script không đổi.
- `resetToken`: 32 byte ngẫu nhiên (`crypto/rand`), hex; không log.

## 4. Dữ liệu / migration

- Dữ liệu trên RAM, không có migration. Tài khoản mẫu có email; khởi động lại là về dữ liệu mẫu.
- Phiên đặt lại: `{token, username, expiresAt = lúc xác minh + 10 phút, status}`; "hết hạn" suy ra từ `expiresAt`.
- BR07 dùng `TokenVersion` trong JWT thay vì so thời điểm (tránh lệch cùng giây giữa đổi mật khẩu và đăng nhập lại).

## 5. E2E (TAF)

- Test case `e2e/testcases/AUTH/AUTH-F02.md` đã duyệt; script `e2e/scripts/AUTH/AUTH-F02.spec.ts` (36 test) dùng
  đúng URL / `data-testid` / API ở trên — không cần sửa. Chạy: `kit e2e AUTH-F02 --no-gate --base-url
  http://localhost:3000` sau `bash scripts/up.sh`.

## 6. Thư viện mới, rủi ro, điểm chưa chắc

- Không thêm thư viện.
- Đổi claim JWT (`ver`): token cấp trước khi deploy vẫn hợp lệ vì `ver` thiếu = 0 = `TokenVersion` ban đầu.
- Rủi ro đã chấp nhận ở PRD §1: không giới hạn số lần xác minh sai.
