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

- Chức năng mới: spec theo ai-kit (PRD/FSD) ở `docs/features/<MOD>-<tên>/`, mã `<MOD>-F##…` — viết bằng `/spec`,
  người duyệt bằng `kit approve` trước khi code. Ví dụ module đăng nhập là `AUTH` (quên mật khẩu, đăng ký: `AUTH-F02`, …).
- Tài liệu cũ, giữ nguyên làm tham khảo: `docs/features/login.md` (ticket AUTH-12, đã làm) và
  `docs/SRS-quan-ly-nhan-vien.md` (mã `REQ-###`, chưa làm — khi triển khai thì chuyển sang module mới qua `/spec`).

## Git

- Commit: `feat(be): …`, `feat(fe): …`, `fix(be): …`, `test(e2e): …`, `docs: …`.
- Không commit `fe/.next/`, `node_modules/`, `.run/`, `e2e/test-results/`.

<!-- ai-kit:start — do `kit init` quản lý; sửa ngoài khối này -->
## ai-kit — quy trình spec → prototype → dev → test → ship

Repo dùng ai-kit (gắn ở `.ai-kit/`, cấu hình `ai-kit.yml`). Có hai lớp quy định:

- **Kỹ thuật** — stack, framework, cấu trúc code, lint, bảo mật, test đơn vị: theo hướng dẫn của
  repo (phần ngoài khối này và các file trong `guides` của `ai-kit.yml`).
- **Quy trình** — tài liệu spec, mã yêu cầu, trạng thái duyệt, version, kế hoạch, test case/e2e,
  truy vết, PR: theo `.ai-kit/RULES.md` và `ai-kit.yml`.

Khi hai lớp đụng nhau: *code thế nào* theo repo; *spec, quy trình, e2e đi đường nào* theo ai-kit.
Mâu thuẫn thật thì dừng và hỏi người dùng, không tự chọn.

Trước mọi việc thuộc quy trình, chạy `node .ai-kit/cli/bin/kit.mjs context`. Chưa có `.ai-kit/` thì
gắn kit theo `kitRef` trong `ai-kit.yml`: `git clone --depth 1 -b <kitRef> https://github.com/dunglt-iks/ai-kit .ai-kit`
rồi `npm ci --prefix .ai-kit --omit=dev`. Duyệt spec / prototype / test case là việc của người, bằng
`kit approve <mã>` (thêm `--prototype` / `--testcase`) — AI không tự đặt trạng thái đã duyệt.

- `/spec` — làm rõ yêu cầu, ghi vào tài liệu spec (trạng thái nháp)
- `/prototype` — dựng HTML tĩnh từ wireframe trong spec
- `/testcase` — sinh/đồng bộ test case e2e từ spec + prototype đã duyệt
- `/e2e-script` — sinh/cập nhật script e2e từ test case đã duyệt, chạy thử
- `/e2e-triage` — phân loại test e2e fail (ở máy hoặc lần chạy CI), sửa script, ghi bug app
- `/dev` — kế hoạch → code → test → review theo spec đã duyệt; sửa bug từ issue
- `/ship` — kiểm tra cuối, mở PR, theo dõi CI
- `/kit-adopt` — đề xuất cấu hình ai-kit.yml khi mới nhúng kit
<!-- ai-kit:end -->
