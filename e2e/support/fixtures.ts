// Fixture dùng chung. Script import { test, expect, persona } từ đây, không import @playwright/test trực tiếp.
//   fx.uniq("u")                      chuỗi duy nhất "e2e-u-3k9x…" (tên đăng nhập, tên hiển thị…)
//   fx.user({ password?, name?, email? })
//                                     tạo tài khoản qua API test (BE chạy với TEST_API=1), tự xoá sau test.
//                                     email mặc định "<username>@e2e.local"; email: null → tài khoản không có email
//                                     (API: POST /api/test/users {username, name, password, email?})
//   fx.track(username)                đăng ký xoá sau test một tài khoản tạo qua giao diện (vd đăng ký mới)
//   fx.deleteUser(username)           xoá tài khoản ngay (dựng trạng thái "tài khoản không còn" giữa test)
//                                     (API: DELETE /api/test/users/{username})
//   fx.signIn(page, user)             mở phiên đăng nhập cho user trên page: đăng nhập qua API, ghi JWT vào
//                                     localStorage "token" (API: POST /api/auth/login) — dùng cho tiền điều kiện
//                                     "đang đăng nhập", không phải để kiểm chức năng đăng nhập
// Dọn lỗi không làm test fail nhưng ghi annotation "cleanup" để báo cáo thấy.
import { test as base, expect, type Page } from "@playwright/test";
import { persona } from "./personas";

export { expect, persona };

export type TestUser = { username: string; password: string; name: string; email: string | null };

export const test = base.extend<{ fx: Fixtures }>({
  fx: async ({ request }, use, testInfo) => {
    const created: string[] = [];
    const fx: Fixtures = {
      uniq: (label = "") => `e2e-${label ? `${label}-` : ""}${Date.now().toString(36)}${Math.random().toString(36).slice(2, 5)}`,
      async user(o = {}) {
        const username = o.username ?? fx.uniq("u");
        const u: TestUser = {
          username,
          password: o.password ?? "Test@1234",
          name: o.name ?? "Người dùng E2E",
          email: o.email === undefined ? `${username}@e2e.local` : o.email,
        };
        const data: Record<string, string> = { username: u.username, password: u.password, name: u.name };
        if (u.email !== null) data.email = u.email;
        const res = await request.post("/api/test/users", { data });
        expect(res.status(), `${await res.text()} — BE có chạy với TEST_API=1?`).toBe(201);
        created.push(u.username);
        return u;
      },
      track(username) {
        created.push(username);
      },
      async deleteUser(username) {
        const res = await request.delete(`/api/test/users/${encodeURIComponent(username)}`);
        expect(res.status(), await res.text()).toBe(204);
        const i = created.indexOf(username);
        if (i >= 0) created.splice(i, 1);
      },
      async signIn(page, u) {
        const res = await request.post("/api/auth/login", { data: { username: u.username, password: u.password } });
        expect(res.status(), await res.text()).toBe(200);
        const { token } = await res.json();
        await page.goto("/login");
        await page.evaluate((t) => localStorage.setItem("token", t), token);
      },
    };
    await use(fx);
    for (const username of created.reverse()) {
      const res = await request.delete(`/api/test/users/${encodeURIComponent(username)}`).catch(() => null);
      if (!res || !res.ok())
        testInfo.annotations.push({ type: "cleanup", description: `Không xoá được tài khoản ${username} (${res?.status() ?? "lỗi mạng"})` });
    }
  },
});

type Fixtures = {
  uniq: (label?: string) => string;
  user: (o?: Partial<TestUser>) => Promise<TestUser>;
  track: (username: string) => void;
  deleteUser: (username: string) => Promise<void>;
  signIn: (page: Page, u: { username: string; password: string }) => Promise<void>;
};
