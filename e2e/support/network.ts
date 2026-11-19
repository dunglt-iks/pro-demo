// Giả lập phản hồi máy chủ cho một API (dựng tiền điều kiện "mất kết nối", "lỗi hệ thống", "phản hồi chậm").
//   await serverDown(page, "/api/auth/forgot-password/verify")    → request bị huỷ (lỗi mạng)
//   await serverError(page, "/api/auth/forgot-password/verify")   → 500 {"error":{"code":"INTERNAL","message":"Lỗi hệ thống"}}
//   const release = await serverSlow(page, "/api/…"); …; release() → request treo tới khi gọi release(), rồi đi tiếp
import type { Page } from "@playwright/test";

const glob = (path: string) => `**${path}`;

export async function serverDown(page: Page, path: string) {
  await page.route(glob(path), (route) => route.abort("failed"));
}

export async function serverError(page: Page, path: string) {
  await page.route(glob(path), (route) =>
    route.fulfill({
      status: 500,
      contentType: "application/json",
      body: JSON.stringify({ error: { code: "INTERNAL", message: "Lỗi hệ thống" } }),
    }),
  );
}

export async function serverSlow(page: Page, path: string): Promise<() => void> {
  let release!: () => void;
  const gate = new Promise<void>((r) => (release = r));
  await page.route(glob(path), async (route) => {
    await gate;
    await route.continue();
  });
  return release;
}
