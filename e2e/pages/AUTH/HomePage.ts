// AUTH-F01-SCR02 Trang chủ (/).
import type { Locator, Page } from "@playwright/test";

export class HomePage {
  readonly url = "/";
  readonly greeting: Locator; // 1
  readonly logoutButton: Locator; // 2

  constructor(private page: Page) {
    this.greeting = page.getByRole("heading", { level: 1 });
    this.logoutButton = page.getByRole("button", { name: "Đăng xuất" });
  }

  async open() {
    await this.page.goto(this.url);
  }
}
