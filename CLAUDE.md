# CLAUDE.md — project1

Ứng dụng web nội bộ: BE Go (`be/`), FE Next.js (`fe/`), tài liệu yêu cầu ở `docs/`, test e2e ở `e2e/`.
Đọc file này trước khi viết hoặc sửa code.

## Tech stack

- **BE:** Go 1.23 (`go.mod` ghim toolchain), thư viện chuẩn `net/http` (router theo method của Go 1.22+).
  JWT: `github.com/golang-jwt/jwt/v5`. Băm mật khẩu: `golang.org/x/crypto/bcrypt`. Dữ liệu lưu trên RAM
  (`internal/store`, map + mutex) — khởi động lại là về dữ liệu mẫu. Không thêm thư viện khi chưa hỏi.
- **FE:** Next.js (App Router) + TypeScript (strict) + Tailwind CSS. Không thêm thư viện UI.
- **E2E:** Playwright (TypeScript) ở `e2e/` — quy ước trong `e2e/README.md`.

## Chạy

- Cả hai: `bash scripts/up.sh` (BE :8080 với `TEST_API=1`, FE :3000) — tắt: `bash scripts/down.sh`.
- Riêng: `cd be && go run ./cmd/server` · `cd fe && npm run dev`. FE chuyển `/api/*` sang BE (`fe/next.config.ts`).

## BE

- Chia lớp: `internal/httpapi` (handler) → `internal/auth`, … (service) → `internal/store` (lưu trữ).
  Handler chỉ parse / validate đầu vào, gọi service, trả JSON — không chứa logic nghiệp vụ.
- Lỗi nghiệp vụ là biến `Err…` trong package service; handler đổi sang HTTP status + mã lỗi.
  Lỗi trả về dạng `{"error": {"code": "INVALID_CREDENTIALS", "message": "…"}}`, message tiếng Việt cho người dùng.
- Endpoint cần đăng nhập bọc bằng `requireAuth` (đọc `Authorization: Bearer <jwt>`).
- API chỉ dùng cho test (`/api/test/*`) chỉ đăng ký khi `TEST_API=1`; không gọi từ FE.
- Bảo mật: mật khẩu chỉ lưu dạng bcrypt; không log body request, mật khẩu, token; khoá ký JWT lấy từ `JWT_SECRET`.
  Thông báo lỗi đăng nhập không được để lộ tên đăng nhập nào tồn tại.
- Trước khi commit: `cd be && gofmt -l . && go vet ./... && go build ./...` (gofmt không in gì).

## FE

- Trang ở `fe/src/app/`, component ở `fe/src/components/`, tiện ích ở `fe/src/lib/`.
- Gọi API **chỉ** qua `apiClient` (`fe/src/lib/api.ts`): tự gắn Bearer token, đổi lỗi thành `ApiError`, gặp 401 thì
  xoá token và chuyển về `/login`. Token đọc / ghi qua `session` (`fe/src/lib/session.ts`, localStorage).
- Component có state / sự kiện đánh dấu `"use client"`. Trang cần đăng nhập tự kiểm token, không có thì về `/login`.
- Style bằng class Tailwind, không viết file CSS riêng. Văn bản giao diện tiếng Việt.
- Form: `<label htmlFor>` cho mọi ô nhập; thông báo lỗi `role="alert"` kèm `data-testid` (Next.js đã có sẵn một phần tử
  `role="alert"` ẩn cho route announcer).
- Trước khi commit: `cd fe && npm run lint && npm run build`.

## Tài liệu

- Yêu cầu chức năng ở `docs/`: `docs/features/<chức năng>.md` (ticket: user story + tiêu chí chấp nhận) và
  `docs/SRS-*.md` (đặc tả theo phân hệ, mã `REQ-###`). Thêm chức năng mới thì viết tài liệu trước khi code.

## Git

- Commit: `feat(be): …`, `feat(fe): …`, `fix(be): …`, `test(e2e): …`, `docs: …`.
- Không commit `fe/.next/`, `node_modules/`, `.run/`, `e2e/test-results/`.
