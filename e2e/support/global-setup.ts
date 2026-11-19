// Đăng nhập từng persona một lần qua API, lưu JWT vào localStorage của storageState
// (e2e/test-results/.auth/<persona>.json). Test dùng lại, không đăng nhập qua giao diện.
import fs from "node:fs";
import path from "node:path";
import { request, type FullConfig } from "@playwright/test";
import { AUTH_DIR, personas } from "./personas";

export default async function globalSetup(config: FullConfig) {
  const baseURL = config.projects[0].use.baseURL!;
  const api = await request.newContext({ baseURL });
  fs.mkdirSync(AUTH_DIR, { recursive: true });
  for (const [name, p] of Object.entries(personas())) {
    const res = await api.post("/api/auth/login", { data: { username: p.username, password: p.password } });
    if (!res.ok()) throw new Error(`Đăng nhập persona ${name} lỗi ${res.status()}: ${await res.text()} — app đã chạy chưa (${baseURL})?`);
    const { token } = await res.json();
    const state = { cookies: [], origins: [{ origin: new URL(baseURL).origin, localStorage: [{ name: "token", value: token }] }] };
    fs.writeFileSync(path.join(AUTH_DIR, `${name}.json`), JSON.stringify(state));
  }
  await api.dispose();
}
