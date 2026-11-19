// Đọc e2e/testdata/personas.yml (dạng đơn giản: tên → các trường) và đường dẫn storageState đã đăng nhập.
import fs from "node:fs";
import path from "node:path";

const E2E = path.resolve(__dirname, "..");
export const AUTH_DIR = path.join(E2E, "test-results", ".auth");

export type Persona = { username: string; password: string; name: string };

export function personas(): Record<string, Persona> {
  const out: Record<string, Record<string, string>> = {};
  let cur = "";
  for (const line of fs.readFileSync(path.join(E2E, "testdata", "personas.yml"), "utf8").split("\n")) {
    if (/^\s*(#|$)/.test(line)) continue;
    const top = /^([\w-]+):\s*$/.exec(line);
    const field = /^\s+(\w+):\s*(.*)$/.exec(line);
    if (top) out[(cur = top[1])] = {};
    else if (field && cur) out[cur][field[1]] = field[2].trim();
  }
  return out as Record<string, Persona>;
}

/** storageState (JWT trong localStorage) của persona — dùng: test.use({ storageState: persona("user1") }). */
export function persona(name: string): string {
  if (!personas()[name]) throw new Error(`Không có persona "${name}" trong e2e/testdata/personas.yml`);
  return path.join(AUTH_DIR, `${name}.json`);
}
