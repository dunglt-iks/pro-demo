// AUTH-F02-SCR02 Quên mật khẩu — đặt mật khẩu mới (/forgot-password/reset).
import type { Locator, Page } from "@playwright/test";

export const RESET_API = "/api/auth/forgot-password/reset";

export class ResetPasswordPage {
  readonly url = "/forgot-password/reset";
  readonly title: Locator; // 1
  readonly account: Locator; // 2
  readonly password: Locator; // 3
  readonly hint: Locator; // gợi ý dưới 3
  readonly confirm: Locator; // 4
  readonly error: Locator; // 5
  readonly submit: Locator; // 6
  readonly backToLogin: Locator; // 7

  constructor(private page: Page) {
    this.title = page.getByRole("heading", { name: "Đặt mật khẩu mới" });
    this.account = page.getByText(/^Tài khoản: /);
    this.password = page.getByLabel("Mật khẩu mới", { exact: true });
    this.hint = page.getByText("Từ 8 đến 16 ký tự");
    this.confirm = page.getByLabel("Nhập lại mật khẩu mới");
    this.error = page.getByTestId("reset-error");
    this.submit = page.getByRole("button", { name: /Đặt lại mật khẩu|Đang lưu…/ });
    this.backToLogin = page.getByRole("link", { name: /Về đăng nhập/ });
  }

  async open() {
    await this.page.goto(this.url);
  }

  async setPassword(password: string, confirm = password) {
    await this.password.fill(password);
    await this.confirm.fill(confirm);
    await this.submit.click();
  }
}
