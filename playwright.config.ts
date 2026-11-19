// Cấu hình Playwright của project1. App phải đang chạy (bash scripts/up.sh) — Playwright không tự bật app.
// Địa chỉ app: E2E_BASE_URL (mặc định http://localhost:3000).
import { defineConfig, devices } from "@playwright/test";

export default defineConfig({
  testDir: "./e2e",
  testMatch: /.*\.spec\.ts$/,
  outputDir: "./e2e/test-results/playwright",
  globalSetup: "./e2e/support/global-setup.ts",
  fullyParallel: true,
  retries: process.env.CI ? 1 : 0,
  reporter: [["list"]],
  use: {
    baseURL: process.env.E2E_BASE_URL || "http://localhost:3000",
    trace: "retain-on-failure",
    screenshot: "only-on-failure",
  },
  projects: [{ name: "chromium", use: { ...devices["Desktop Chrome"] } }],
});
