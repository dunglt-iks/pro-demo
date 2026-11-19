// Lưu JWT ở localStorage (khoá "token"). Chỉ gọi ở client.
const KEY = "token";

export const session = {
  get: (): string | null => (typeof window === "undefined" ? null : window.localStorage.getItem(KEY)),
  set: (token: string) => window.localStorage.setItem(KEY, token),
  clear: () => window.localStorage.removeItem(KEY),
};
