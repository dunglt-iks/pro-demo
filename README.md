# project1

Ứng dụng web nội bộ mẫu: **BE Go** (`be/`) + **FE Next.js** (`fe/`). Hiện có chức năng đăng nhập bằng tên đăng nhập /
mật khẩu (JWT). Tài liệu yêu cầu ở `docs/`, khung test e2e (Playwright) ở `e2e/` — chưa có test.

> Đây là repo mẫu **chưa nhúng ai-kit**, dùng để demo việc nhúng kit và làm chức năng mới (quên mật khẩu, đăng ký) theo
> quy trình của kit — xem `docs/12-user-guide.md` của ai-kit.

## Cần có

Go ≥ 1.22 (tự tải toolchain 1.23 theo `be/go.mod`), Node.js ≥ 20.9.

## Chạy

```bash
npm ci                              # Playwright (cho e2e)
(cd fe && npm ci)
bash scripts/up.sh                  # BE :8080 (TEST_API=1) + FE :3000, chờ tới khi sẵn sàng
open http://localhost:3000          # đăng nhập: admin / Admin@123 hoặc user1 / User@123
bash scripts/down.sh
```

E2E: `npx playwright install chromium` (một lần), rồi `npx playwright test` khi app đang chạy.

| Thư mục | Nội dung |
|---|---|
| `be/` | Go: `cmd/server`, `internal/{httpapi,auth,store}` |
| `fe/` | Next.js App Router: `src/app`, `src/components`, `src/lib` |
| `docs/` | `features/login.md` (ticket), `SRS-quan-ly-nhan-vien.md` (đặc tả, chưa triển khai) |
| `e2e/` | Khung Playwright: `support/`, `testdata/`, `pages/`, hướng dẫn `README.md` |
| `scripts/` | `up.sh`, `down.sh` |
