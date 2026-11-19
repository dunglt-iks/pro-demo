// Fixture dùng chung. Script import { test, expect, persona } từ đây, không import @playwright/test trực tiếp.
//   fx.uniq("u")                      chuỗi duy nhất "e2e-u-3k9x…" (tên đăng nhập, tên hiển thị…)
//   fx.user({ password?, name? })     tạo tài khoản qua API test (BE chạy với TEST_API=1), tự xoá sau test
//   fx.track(username)                đăng ký xoá sau test một tài khoản tạo qua giao diện (vd đăng ký mới)
// Dọn lỗi không làm test fail nhưng ghi annotation "cleanup" để báo cáo thấy.
import { test as base, expect } from "@playwright/test";
import { persona } from "./personas";

export { expect, persona };

export type TestUser = { username: string; password: string; name: string };

export const test = base.extend<{ fx: Fixtures }>({
  fx: async ({ request }, use, testInfo) => {
    const created: string[] = [];
    const fx: Fixtures = {
      uniq: (label = "") => `e2e-${label ? `${label}-` : ""}${Date.now().toString(36)}${Math.random().toString(36).slice(2, 5)}`,
      async user(o = {}) {
        const u = { username: fx.uniq("u"), password: "Test@1234", name: "Người dùng E2E", ...o };
        const res = await request.post("/api/test/users", { data: u });
        expect(res.status(), `${await res.text()} — BE có chạy với TEST_API=1?`).toBe(201);
        created.push(u.username);
        return u;
      },
      track(username) {
        created.push(username);
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
};
