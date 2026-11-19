// AUTH-F02-SCR01 Quên mật khẩu — xác minh tài khoản (/forgot-password).
import type { Locator, Page } from "@playwright/test";

export const VERIFY_API = "/api/auth/forgot-password/verify";

export class ForgotPasswordPage {
  readonly url = "/forgot-password";
  readonly title: Locator; // 1
  readonly guide: Locator; // 2
  readonly username: Locator; // 3
  readonly email: Locator; // 4
  readonly error: Locator; // 5
  readonly submit: Locator; // 6
  readonly backToLogin: Locator; // 7

  constructor(private page: Page) {
    this.title = page.getByRole("heading", { name: "Quên mật khẩu" });
    this.guide = page.getByText("Nhập tên đăng nhập và email của tài khoản để đặt lại mật khẩu.");
    this.username = page.getByLabel("Tên đăng nhập");
    this.email = page.getByLabel("Email");
    this.error = page.getByTestId("forgot-error");
    this.submit = page.getByRole("button", { name: /Tiếp tục|Đang kiểm tra…/ });
    this.backToLogin = page.getByRole("link", { name: /Về đăng nhập/ });
  }

  async open() {
    await this.page.goto(this.url);
  }

  async verify(username: string, email: string) {
    await this.username.fill(username);
    await this.email.fill(email);
    await this.submit.click();
  }
}
