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

// redirectOn401: false → gặp 401 vẫn xoá token nhưng không tự chuyển /login (trang khách tự xử lý).
async function request<T>(
  path: string,
  init: RequestInit & { auth?: boolean; redirectOn401?: boolean } = {},
): Promise<T> {
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
    if (init.redirectOn401 !== false) window.location.assign("/login");
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
  /** Kiểm phiên đăng nhập còn hiệu lực mà không tự chuyển trang (dùng ở trang dành cho người chưa đăng nhập). */
  meQuiet: () => request<User>("/api/me", { redirectOn401: false }),

  // AUTH-F02 Quên mật khẩu
  forgotVerify: (username: string, email: string) =>
    request<{ resetToken: string; username: string }>("/api/auth/forgot-password/verify", {
      method: "POST",
      body: JSON.stringify({ username, email }),
      auth: false,
    }),
  forgotCheck: (resetToken: string) =>
    request<{ username: string }>("/api/auth/forgot-password/check", {
      method: "POST",
      body: JSON.stringify({ resetToken }),
      auth: false,
    }),
  forgotReset: (resetToken: string, password: string, confirm: string) =>
    request<void>("/api/auth/forgot-password/reset", {
      method: "POST",
      body: JSON.stringify({ resetToken, password, confirm }),
      auth: false,
    }),
  forgotCancel: (resetToken: string) =>
    request<void>("/api/auth/forgot-password/cancel", {
      method: "POST",
      body: JSON.stringify({ resetToken }),
      auth: false,
    }),
};

/** Đổi lỗi gọi API thành câu hiển thị: lỗi từ máy chủ dùng message của nó, không gọi được máy chủ thì báo mất kết nối. */
export function errorMessage(err: unknown): string {
  return err instanceof ApiError ? err.message : "Không kết nối được máy chủ";
}
