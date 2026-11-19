// apiClient duy nhất của FE: gắn "Authorization: Bearer <jwt>", đổi lỗi API thành ApiError.
// Gặp 401 trên API cần đăng nhập → xoá token, chuyển về /login.
import { session } from "./session";

export type User = { id: number; username: string; name: string; createdAt: string };

export class ApiError extends Error {
  constructor(
    public status: number,
    public code: string,
    message: string,
  ) {
    super(message);
  }
}

async function request<T>(path: string, init: RequestInit & { auth?: boolean } = {}): Promise<T> {
  const headers = new Headers(init.headers);
  headers.set("Content-Type", "application/json");
  const token = session.get();
  if (init.auth !== false && token) headers.set("Authorization", `Bearer ${token}`);
  const res = await fetch(path, { ...init, headers });
  if (res.ok) return (res.status === 204 ? undefined : await res.json()) as T;
  const body = await res.json().catch(() => ({}));
  const err = new ApiError(res.status, body?.error?.code ?? "INTERNAL", body?.error?.message ?? "Lỗi hệ thống");
  if (res.status === 401 && init.auth !== false) {
    session.clear();
    window.location.assign("/login");
  }
  throw err;
}

export const apiClient = {
  login: (username: string, password: string) =>
    request<{ token: string; user: User }>("/api/auth/login", {
      method: "POST",
      body: JSON.stringify({ username, password }),
      auth: false,
    }),
  me: () => request<User>("/api/me"),
};
