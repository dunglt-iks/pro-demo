// Trạng thái phía trình duyệt. Chỉ gọi ở client.
//   session       JWT đăng nhập ở localStorage (khoá "token").
//   resetSession  mã phiên đặt lại mật khẩu của tab này (sessionStorage) — AUTH-F02. Giữ lại cả khi phiên đã dùng /
//                 đã huỷ: mở lại màn đặt mật khẩu thì máy chủ báo hết hạn (khác "chưa xác minh" = không có mã).
//   flash         thông báo một lần cho trang kế tiếp (sessionStorage); đọc là xoá, nên tải lại trang thì mất.
const KEY = "token";
const RESET_KEY = "resetToken";
const FLASH_PREFIX = "flash:";

export const session = {
  get: (): string | null => (typeof window === "undefined" ? null : window.localStorage.getItem(KEY)),
  set: (token: string) => window.localStorage.setItem(KEY, token),
  clear: () => window.localStorage.removeItem(KEY),
};

export const resetSession = {
  get: (): string | null => window.sessionStorage.getItem(RESET_KEY),
  set: (token: string) => window.sessionStorage.setItem(RESET_KEY, token),
};

export type FlashPage = "login" | "forgot-password";

export const flash = {
  set: (page: FlashPage, message: string) => window.sessionStorage.setItem(FLASH_PREFIX + page, message),
  take: (page: FlashPage): string => {
    const message = window.sessionStorage.getItem(FLASH_PREFIX + page) ?? "";
    window.sessionStorage.removeItem(FLASH_PREFIX + page);
    return message;
  },
};
