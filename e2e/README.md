# Hướng dẫn viết test e2e — project1

Playwright (TypeScript). App phải đang chạy: `bash scripts/up.sh` (BE chạy với `TEST_API=1`). Chạy test: `npx
playwright test` ở gốc `project1` (địa chỉ app: `E2E_BASE_URL`, mặc định `http://localhost:3000`).

## Thư mục

| Thư mục | Nội dung |
|---|---|
| `support/` | `fixtures.ts` (test, expect, persona, fixture dữ liệu), `global-setup.ts` (đăng nhập persona), `personas.ts` |
| `testdata/` | `personas.yml` — tài khoản dùng trong test |
| `pages/` | Page object (hiện chưa có) — mỗi màn hình một class |
| `test-results/` | Kết quả chạy, không commit |

## Import

```ts
import { test, expect, persona } from "<đường dẫn tương đối>/support/fixtures";
```

Không import `@playwright/test` trực tiếp trong test.

## Đăng nhập — persona

Có sẵn `admin`, `user1` (xem `testdata/personas.yml`). `global-setup.ts` đăng nhập mỗi persona một lần qua API và lưu
JWT vào `localStorage` của storageState. Test cần đăng nhập:

```ts
test.use({ storageState: persona("user1") });
```

Test **về chính chức năng đăng nhập** thì không dùng persona — đăng nhập qua giao diện.

## Dữ liệu

- Tài khoản riêng cho test: `const u = await fx.user({ name: "…" })` → `{ username, password, name }`, tạo qua
  `/api/test/users`, tự xoá sau test. Không sửa / dùng mật khẩu của `admin`, `user1` cho test thay đổi dữ liệu.
- Tài khoản tạo **qua giao diện** (vd đăng ký mới): `fx.track(username)` để xoá sau test.
- Chuỗi duy nhất: `fx.uniq("label")`. Test chạy song song trên cùng BE, nên không giả định dữ liệu của test khác.

## Locator

- Ưu tiên `getByLabel` (ô nhập có `<label>`), `getByRole(…, { name })` (nút, tiêu đề, link).
- Thông báo lỗi: `getByTestId(…)` — **không** dùng `getByRole("alert")` trơn: Next.js có sẵn một phần tử `role="alert"`
  ẩn (route announcer) nên khớp 2 phần tử.
- Không `waitForTimeout`. Chờ bằng `expect(…)` (tự chờ) hoặc `waitForURL`.

## Kết quả mong đợi

- Kiểm đúng chữ trên màn hình: `toHaveText("Sai tên đăng nhập hoặc mật khẩu")`, `toHaveURL(/\/login$/)`.
- Logic thuần BE (ràng buộc dữ liệu, phân quyền API) kiểm qua `request` (fixture của Playwright), không cần mở trình duyệt.
