// AUTH-F01-SCR01 Trang đăng nhập (/login).
import type { Locator, Page } from "@playwright/test";

export class LoginPage {
  readonly url = "/login";
  readonly successMessage: Locator; // 2
  readonly username: Locator; // 3
  readonly password: Locator; // 4
  readonly error: Locator; // 5
  readonly submit: Locator; // 6
  readonly forgotLink: Locator; // 7

  constructor(private page: Page) {
    this.successMessage = page.getByTestId("login-success");
    this.username = page.getByLabel("Tên đăng nhập");
    this.password = page.getByLabel("Mật khẩu", { exact: true });
    this.error = page.getByTestId("login-error");
    this.submit = page.getByRole("button", { name: /Đăng nhập|Đang đăng nhập…/ });
    this.forgotLink = page.getByRole("link", { name: "Quên mật khẩu?" });
  }

  async open() {
    await this.page.goto(this.url);
  }

  async login(username: string, password: string) {
    await this.username.fill(username);
    await this.password.fill(password);
    await this.submit.click();
  }
}
