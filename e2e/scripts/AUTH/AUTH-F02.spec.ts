// @taf-source e2e/testcases/AUTH/AUTH-F02.md sha256:b60ecccfd4763623
import { test, expect, persona } from "../../support/fixtures";
import { serverDown, serverError, serverSlow } from "../../support/network";
import { LoginPage } from "../../pages/AUTH/LoginPage";
import { HomePage } from "../../pages/AUTH/HomePage";
import { ForgotPasswordPage, VERIFY_API } from "../../pages/AUTH/ForgotPasswordPage";
import { ResetPasswordPage, RESET_API } from "../../pages/AUTH/ResetPasswordPage";
import type { Page } from "@playwright/test";

const SUCCESS = "Đặt lại mật khẩu thành công. Đăng nhập bằng mật khẩu mới.";
const EXPIRED = "Phiên xác minh đã hết hạn, vui lòng xác minh lại";
const FORGOT_URL = /\/forgot-password$/;
const RESET_URL = /\/forgot-password\/reset$/;
const LOGIN_URL = /\/login$/;
const HOME_URL = /\/$/;

/** Tiền điều kiện "T đã xác minh xong, đang ở màn Đặt mật khẩu mới". */
async function verified(page: Page, u: { username: string; email: string | null }) {
  const forgot = new ForgotPasswordPage(page);
  await forgot.open();
  await forgot.verify(u.username, u.email ?? "");
  await expect(page).toHaveURL(RESET_URL);
  const reset = new ResetPasswordPage(page);
  // Màn đã sẵn sàng khi kiểm phiên xong và hiện "Tài khoản: …" (trước đó nút Đặt lại còn khoá).
  await expect(reset.account).toHaveText(`Tài khoản: ${u.username}`);
  return reset;
}

// ── AUTH-F02-US01 Xác minh tài khoản ──────────────────────────────────────────

test("TC-AUTH-F02-001 Quên mật khẩu, đặt mật khẩu mới rồi đăng nhập lại", async ({ page, fx }) => {
  const t = await fx.user();
  const login = new LoginPage(page);
  const forgot = new ForgotPasswordPage(page);
  const reset = new ResetPasswordPage(page);
  const home = new HomePage(page);

  await test.step("Bước 1: bấm Quên mật khẩu? → màn Quên mật khẩu có hướng dẫn", async () => {
    await login.open();
    await login.forgotLink.click();
    await expect(page).toHaveURL(FORGOT_URL);
    await expect(forgot.title).toHaveText("Quên mật khẩu");
    await expect(forgot.guide).toHaveText("Nhập tên đăng nhập và email của tài khoản để đặt lại mật khẩu.");
  });
  await test.step("Bước 2: xác minh tên + email → màn Đặt mật khẩu mới, Tài khoản, gợi ý", async () => {
    await forgot.verify(t.username, t.email!);
    await expect(page).toHaveURL(RESET_URL);
    await expect(reset.title).toHaveText("Đặt mật khẩu mới");
    await expect(reset.account).toHaveText(`Tài khoản: ${t.username}`);
    await expect(reset.hint).toHaveText("Từ 8 đến 16 ký tự");
  });
  await test.step("Bước 3: đặt mật khẩu Moi@12345 → về đăng nhập, thông báo thành công", async () => {
    await reset.setPassword("Moi@12345");
    await expect(page).toHaveURL(LOGIN_URL);
    await expect(login.successMessage).toHaveText(SUCCESS);
  });
  await test.step("Bước 4: tải lại trang → không còn thông báo thành công", async () => {
    await page.reload();
    await expect(login.successMessage).toHaveCount(0);
  });
  await test.step("Bước 5: đăng nhập bằng mật khẩu cũ → Sai tên đăng nhập hoặc mật khẩu", async () => {
    await login.login(t.username, "Test@1234");
    await expect(login.error).toHaveText("Sai tên đăng nhập hoặc mật khẩu");
  });
  await test.step("Bước 6: đăng nhập bằng mật khẩu mới → trang chủ, Xin chào", async () => {
    await login.login(t.username, "Moi@12345");
    await expect(page).toHaveURL(HOME_URL);
    await expect(home.greeting).toHaveText("Xin chào Người dùng E2E");
  });
});

const caseRows: { label: string; user: (u: string) => string; email: (e: string) => string }[] = [
  { label: "tên viết HOA", user: (u) => u.toUpperCase(), email: (e) => e },
  { label: "email viết HOA", user: (u) => u, email: (e) => e.toUpperCase() },
  { label: "có khoảng trắng đầu cuối", user: (u) => `  ${u}  `, email: (e) => `  ${e}  ` },
];
caseRows.forEach((row, i) => {
  test(`TC-AUTH-F02-002 Xác minh không phân biệt hoa thường, bỏ khoảng trắng đầu / cuối [${i + 1}] ${row.label}`, async ({ page, fx }) => {
    const t = await fx.user();
    const forgot = new ForgotPasswordPage(page);
    const reset = new ResetPasswordPage(page);
    await test.step("Bước 1: nhập theo dữ liệu, bấm Tiếp tục → sang màn Đặt mật khẩu mới", async () => {
      await forgot.open();
      await forgot.verify(row.user(t.username), row.email(t.email!));
      await expect(page).toHaveURL(RESET_URL);
      await expect(reset.title).toHaveText("Đặt mật khẩu mới");
    });
  });
});

const blankRows: { label: string; user: (u: string) => string; email: (e: string) => string }[] = [
  { label: "cả hai trống", user: () => "", email: () => "" },
  { label: "tên trống", user: () => "", email: (e) => e },
  { label: "email trống", user: (u) => u, email: () => "" },
  { label: "tên chỉ khoảng trắng", user: () => "   ", email: (e) => e },
  { label: "email chỉ khoảng trắng", user: (u) => u, email: () => "   " },
];
blankRows.forEach((row, i) => {
  test(`TC-AUTH-F02-003 Bỏ trống tên đăng nhập hoặc email [${i + 1}] ${row.label}`, async ({ page, fx }) => {
    const t = await fx.user();
    const forgot = new ForgotPasswordPage(page);
    await test.step("Bước 1: nhập theo dữ liệu, bấm Tiếp tục → Nhập tên đăng nhập và email", async () => {
      await forgot.open();
      await forgot.verify(row.user(t.username), row.email(t.email!));
      await expect(forgot.error).toHaveText("Nhập tên đăng nhập và email");
      await expect(page).toHaveURL(FORGOT_URL);
    });
  });
});

type Accounts = { t: { username: string; email: string | null }; u: { email: string | null }; k: { username: string } };
const mismatchRows: { label: string; input: (a: Accounts, missing: string) => [string, string] }[] = [
  { label: "email không tồn tại", input: (a) => [a.t.username, "khong-co@project1.local"] },
  { label: "tên không tồn tại", input: (a, missing) => [missing, a.t.email!] },
  { label: "email của tài khoản khác", input: (a) => [a.t.username, a.u.email!] },
  { label: "tài khoản không có email", input: (a) => [a.k.username, "k@project1.local"] },
];
mismatchRows.forEach((row, i) => {
  test(`TC-AUTH-F02-004 Tên đăng nhập và email không khớp [${i + 1}] ${row.label}`, async ({ page, fx }) => {
    const accounts: Accounts = { t: await fx.user(), u: await fx.user(), k: await fx.user({ email: null }) };
    const [username, email] = row.input(accounts, fx.uniq("khong-ton-tai"));
    const forgot = new ForgotPasswordPage(page);
    await test.step("Bước 1: nhập theo dữ liệu, bấm Tiếp tục → Tên đăng nhập hoặc email không đúng", async () => {
      await forgot.open();
      await forgot.verify(username, email);
      await expect(forgot.error).toHaveText("Tên đăng nhập hoặc email không đúng");
      await expect(page).toHaveURL(FORGOT_URL);
    });
  });
});

test("TC-AUTH-F02-005 Mất kết nối máy chủ khi xác minh", async ({ page, fx }) => {
  const t = await fx.user();
  const forgot = new ForgotPasswordPage(page);
  await forgot.open();
  await serverDown(page, VERIFY_API);
  await test.step("Bước 1: bấm Tiếp tục → Không kết nối được máy chủ", async () => {
    await forgot.verify(t.username, t.email!);
    await expect(forgot.error).toHaveText("Không kết nối được máy chủ");
    await expect(page).toHaveURL(FORGOT_URL);
  });
});

test("TC-AUTH-F02-006 Lỗi hệ thống khi xác minh", async ({ page, fx }) => {
  const t = await fx.user();
  const forgot = new ForgotPasswordPage(page);
  await forgot.open();
  await serverError(page, VERIFY_API);
  await test.step("Bước 1: bấm Tiếp tục → Lỗi hệ thống", async () => {
    await forgot.verify(t.username, t.email!);
    await expect(forgot.error).toHaveText("Lỗi hệ thống");
    await expect(page).toHaveURL(FORGOT_URL);
  });
});

test("TC-AUTH-F02-007 Nút Tiếp tục khi đang chờ phản hồi", async ({ page, fx }) => {
  const t = await fx.user();
  const forgot = new ForgotPasswordPage(page);
  const reset = new ResetPasswordPage(page);
  await forgot.open();
  const release = await serverSlow(page, VERIFY_API);
  await test.step("Bước 1: bấm Tiếp tục → nút Đang kiểm tra… và bị khoá", async () => {
    await forgot.verify(t.username, t.email!);
    await expect(forgot.submit).toHaveText("Đang kiểm tra…");
    await expect(forgot.submit).toBeDisabled();
  });
  await test.step("Bước 2: máy chủ phản hồi → sang màn Đặt mật khẩu mới", async () => {
    release();
    await expect(page).toHaveURL(RESET_URL);
    await expect(reset.title).toHaveText("Đặt mật khẩu mới");
  });
});

// ── AUTH-F02-US02 Đặt mật khẩu mới ────────────────────────────────────────────

const emptyPwRows: { label: string; password: string; confirm: string }[] = [
  { label: "cả hai trống", password: "", confirm: "" },
  { label: "mật khẩu mới trống", password: "", confirm: "Moi@12345" },
  { label: "nhập lại trống", password: "Moi@12345", confirm: "" },
];
emptyPwRows.forEach((row, i) => {
  test(`TC-AUTH-F02-008 Bỏ trống ô mật khẩu [${i + 1}] ${row.label}`, async ({ page, fx }) => {
    const t = await fx.user();
    const reset = await verified(page, t);
    await test.step("Bước 1: bấm Đặt lại mật khẩu → Nhập mật khẩu mới và nhập lại mật khẩu", async () => {
      await reset.setPassword(row.password, row.confirm);
      await expect(reset.error).toHaveText("Nhập mật khẩu mới và nhập lại mật khẩu");
      await expect(page).toHaveURL(RESET_URL);
    });
    await test.step("Bước 2: đăng nhập bằng Test@1234 → trang chủ", async () => {
      const login = new LoginPage(page);
      await login.open();
      await login.login(t.username, "Test@1234");
      await expect(page).toHaveURL(HOME_URL);
      await expect(new HomePage(page).greeting).toHaveText("Xin chào Người dùng E2E");
    });
  });
});

test("TC-AUTH-F02-009 Mật khẩu ngắn hơn 8 ký tự", async ({ page, fx }) => {
  const t = await fx.user();
  const reset = await verified(page, t);
  await test.step("Bước 1: nhập Ab@1234 (7 ký tự) → Mật khẩu phải có ít nhất 8 ký tự", async () => {
    await reset.setPassword("Ab@1234");
    await expect(reset.error).toHaveText("Mật khẩu phải có ít nhất 8 ký tự");
    await expect(page).toHaveURL(RESET_URL);
  });
});

test("TC-AUTH-F02-010 Mật khẩu dài hơn 16 ký tự", async ({ page, fx }) => {
  const t = await fx.user();
  const reset = await verified(page, t);
  await test.step("Bước 1: nhập Abcdefgh@12345678 (17 ký tự) → Mật khẩu tối đa 16 ký tự", async () => {
    await reset.setPassword("Abcdefgh@12345678");
    await expect(reset.error).toHaveText("Mật khẩu tối đa 16 ký tự");
    await expect(page).toHaveURL(RESET_URL);
  });
});

const validPwRows: { label: string; password: string }[] = [
  { label: "đúng 8 ký tự", password: "Abcd@123" },
  { label: "đúng 16 ký tự", password: "Abcdefgh@1234567" },
  { label: "có dấu và khoảng trắng", password: "Mật khẩu 12" },
  { label: "trùng mật khẩu hiện tại", password: "Test@1234" },
];
validPwRows.forEach((row, i) => {
  test(`TC-AUTH-F02-011 Độ dài mật khẩu hợp lệ ở biên [${i + 1}] ${row.label}`, async ({ page, fx }) => {
    const t = await fx.user();
    const reset = await verified(page, t);
    const login = new LoginPage(page);
    await test.step("Bước 1: đặt mật khẩu theo dữ liệu → về đăng nhập, thông báo thành công", async () => {
      await reset.setPassword(row.password);
      await expect(page).toHaveURL(LOGIN_URL);
      await expect(login.successMessage).toHaveText(SUCCESS);
    });
    await test.step("Bước 2: đăng nhập bằng mật khẩu vừa đặt → trang chủ", async () => {
      await login.login(t.username, row.password);
      await expect(page).toHaveURL(HOME_URL);
      await expect(new HomePage(page).greeting).toHaveText("Xin chào Người dùng E2E");
    });
  });
});

test("TC-AUTH-F02-012 Nhập lại mật khẩu không khớp", async ({ page, fx }) => {
  const t = await fx.user();
  const reset = await verified(page, t);
  await test.step("Bước 1: Moi@12345 / Moi@12346 → Mật khẩu nhập lại không khớp", async () => {
    await reset.setPassword("Moi@12345", "Moi@12346");
    await expect(reset.error).toHaveText("Mật khẩu nhập lại không khớp");
    await expect(page).toHaveURL(RESET_URL);
  });
});

test("TC-AUTH-F02-013 Phiên đặt lại hết hiệu lực khi bấm Đặt lại (xác minh lại ở tab khác)", async ({ page, fx }) => {
  const t = await fx.user();
  let resetA!: ResetPasswordPage;
  await test.step("Bước 1: tab A xác minh T → tab A ở màn Đặt mật khẩu mới", async () => {
    resetA = await verified(page, t);
    await expect(resetA.title).toHaveText("Đặt mật khẩu mới");
  });
  await test.step("Bước 2: tab B xác minh T lần nữa → tab B ở màn Đặt mật khẩu mới", async () => {
    const tabB = await page.context().newPage();
    const resetB = await verified(tabB, t);
    await expect(resetB.title).toHaveText("Đặt mật khẩu mới");
  });
  await test.step("Bước 3: tab A nhập mật khẩu quá ngắn, bấm Đặt lại → về Quên mật khẩu, phiên hết hạn", async () => {
    await resetA.setPassword("Ab@1");
    await expect(page).toHaveURL(FORGOT_URL);
    await expect(new ForgotPasswordPage(page).error).toHaveText(EXPIRED);
  });
  await test.step("Bước 4: đăng nhập bằng Test@1234 → trang chủ (mật khẩu không đổi)", async () => {
    const login = new LoginPage(page);
    await login.open();
    await login.login(t.username, "Test@1234");
    await expect(page).toHaveURL(HOME_URL);
    await expect(new HomePage(page).greeting).toHaveText("Xin chào Người dùng E2E");
  });
});

test("TC-AUTH-F02-014 Phiên đặt lại hết hiệu lực khi tài khoản bị xoá", async ({ page, fx }) => {
  const t = await fx.user();
  const reset = await verified(page, t);
  await fx.deleteUser(t.username);
  await test.step("Bước 1: bấm Đặt lại mật khẩu → về Quên mật khẩu, phiên hết hạn", async () => {
    await reset.setPassword("Moi@12345");
    await expect(page).toHaveURL(FORGOT_URL);
    await expect(new ForgotPasswordPage(page).error).toHaveText(EXPIRED);
  });
});

test("TC-AUTH-F02-015 Bấm Về đăng nhập ở bước 2 thì huỷ phiên đặt lại", async ({ page, fx }) => {
  const t = await fx.user();
  const reset = await verified(page, t);
  const login = new LoginPage(page);
  await test.step("Bước 1: bấm Về đăng nhập → trang đăng nhập, không có thông báo thành công", async () => {
    await reset.backToLogin.click();
    await expect(page).toHaveURL(LOGIN_URL);
    await expect(login.successMessage).toHaveCount(0);
  });
  await test.step("Bước 2: Back của trình duyệt → về Quên mật khẩu, phiên hết hạn", async () => {
    await page.goBack();
    await expect(page).toHaveURL(FORGOT_URL);
    await expect(new ForgotPasswordPage(page).error).toHaveText(EXPIRED);
  });
  await test.step("Bước 3: đăng nhập bằng Test@1234 → trang chủ (mật khẩu không đổi)", async () => {
    await login.open();
    await login.login(t.username, "Test@1234");
    await expect(page).toHaveURL(HOME_URL);
    await expect(new HomePage(page).greeting).toHaveText("Xin chào Người dùng E2E");
  });
});

test("TC-AUTH-F02-016 Mất kết nối máy chủ khi đặt lại", async ({ page, fx }) => {
  const t = await fx.user();
  const reset = await verified(page, t);
  await serverDown(page, RESET_API);
  await test.step("Bước 1: bấm Đặt lại mật khẩu → Không kết nối được máy chủ", async () => {
    await reset.setPassword("Moi@12345");
    await expect(reset.error).toHaveText("Không kết nối được máy chủ");
    await expect(page).toHaveURL(RESET_URL);
  });
});

test("TC-AUTH-F02-017 Lỗi hệ thống khi đặt lại", async ({ page, fx }) => {
  const t = await fx.user();
  const reset = await verified(page, t);
  await serverError(page, RESET_API);
  await test.step("Bước 1: bấm Đặt lại mật khẩu → Lỗi hệ thống", async () => {
    await reset.setPassword("Moi@12345");
    await expect(reset.error).toHaveText("Lỗi hệ thống");
    await expect(page).toHaveURL(RESET_URL);
  });
});

test("TC-AUTH-F02-018 Đặt lại mật khẩu làm phiên đăng nhập cũ hết hiệu lực", async ({ page, browser, baseURL, fx }) => {
  const t = await fx.user();
  const home = new HomePage(page);
  await fx.signIn(page, t); // trình duyệt X đang đăng nhập
  await home.open();
  await expect(home.greeting).toHaveText("Xin chào Người dùng E2E");
  const browserY = await browser.newContext({ baseURL });
  try {
    await test.step("Bước 1: trình duyệt Y xác minh T, đặt Moi@12345 → về đăng nhập, thông báo thành công", async () => {
      const pageY = await browserY.newPage();
      const reset = await verified(pageY, t);
      await reset.setPassword("Moi@12345");
      await expect(pageY).toHaveURL(LOGIN_URL);
      await expect(new LoginPage(pageY).successMessage).toHaveText(SUCCESS);
    });
    await test.step("Bước 2: trình duyệt X tải lại trang chủ → chuyển về đăng nhập", async () => {
      await page.reload();
      await expect(page).toHaveURL(LOGIN_URL);
    });
  } finally {
    await browserY.close();
  }
});

test("TC-AUTH-F02-019 Mở màn Đặt mật khẩu mới khi chưa xác minh", async ({ page }) => {
  const forgot = new ForgotPasswordPage(page);
  await test.step("Bước 1: mở thẳng màn Đặt mật khẩu mới → về Quên mật khẩu, không có lỗi", async () => {
    await new ResetPasswordPage(page).open();
    await expect(page).toHaveURL(FORGOT_URL);
    await expect(forgot.title).toHaveText("Quên mật khẩu");
    await expect(forgot.error).toHaveCount(0);
  });
});

test("TC-AUTH-F02-020 Mở lại màn Đặt mật khẩu mới sau khi đã đặt lại xong", async ({ page, fx }) => {
  const t = await fx.user();
  const reset = await verified(page, t);
  await reset.setPassword("Moi@12345");
  await expect(page).toHaveURL(LOGIN_URL);
  await test.step("Bước 1: mở lại màn Đặt mật khẩu mới → về Quên mật khẩu, phiên hết hạn", async () => {
    await reset.open();
    await expect(page).toHaveURL(FORGOT_URL);
    await expect(new ForgotPasswordPage(page).error).toHaveText(EXPIRED);
  });
});

test.describe("persona user1", () => {
  test.use({ storageState: persona("user1") });

  test("TC-AUTH-F02-021 Người đang đăng nhập mở màn Quên mật khẩu", async ({ page }) => {
    const home = new HomePage(page);
    await test.step("Bước 1: mở thẳng màn Quên mật khẩu → trang chủ, Xin chào Nguyễn Văn A", async () => {
      await new ForgotPasswordPage(page).open();
      await expect(page).toHaveURL(HOME_URL);
      await expect(home.greeting).toHaveText("Xin chào Nguyễn Văn A");
    });
    await test.step("Bước 2: mở thẳng màn Đặt mật khẩu mới → trang chủ", async () => {
      await new ResetPasswordPage(page).open();
      await expect(page).toHaveURL(HOME_URL);
      await expect(home.greeting).toHaveText("Xin chào Nguyễn Văn A");
    });
  });
});

// ── Tầng api ──────────────────────────────────────────────────────────────────

test("TC-AUTH-F02-022 Xác minh sai không để lộ tài khoản (API)", async ({ request, fx }) => {
  const t = await fx.user();
  let first!: { status: number; body: unknown };
  await test.step("Bước 1: tên không tồn tại + email bất kỳ → Tên đăng nhập hoặc email không đúng", async () => {
    const res = await request.post(VERIFY_API, { data: { username: fx.uniq("khong-ton-tai"), email: "bat-ky@e2e.local" } });
    first = { status: res.status(), body: await res.json() };
    expect(res.ok()).toBe(false);
    expect((first.body as { error: { message: string } }).error.message).toBe("Tên đăng nhập hoặc email không đúng");
  });
  await test.step("Bước 2: tên của T + email sai → cùng mã trạng thái, cùng nội dung lỗi", async () => {
    const res = await request.post(VERIFY_API, { data: { username: t.username, email: "sai@e2e.local" } });
    expect(res.status()).toBe(first.status);
    expect(await res.json()).toEqual(first.body);
  });
});
