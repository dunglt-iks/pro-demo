"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { ApiError, apiClient, errorMessage } from "@/lib/api";
import { flash, resetSession } from "@/lib/session";
import { useGuestOnly } from "@/lib/useGuestOnly";

const SUCCESS = "Đặt lại mật khẩu thành công. Đăng nhập bằng mật khẩu mới.";

// AUTH-F02-SCR02 — bước 2: đặt mật khẩu mới (AUTH-F02-US02).
export function ResetPasswordForm() {
  const router = useRouter();
  const guest = useGuestOnly();
  const [token, setToken] = useState<string | null>(null);
  const [account, setAccount] = useState("");
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  // Phiên đặt lại hết hiệu lực (AUTH-BR06, AUTH-F02-US02-A4) → về bước 1 kèm thông báo.
  function backToVerify(message: string) {
    flash.set("forgot-password", message);
    router.replace("/forgot-password");
  }

  // Khi mở trang: chưa xác minh → về bước 1, không thông báo; phiên hết hiệu lực → về bước 1 kèm thông báo.
  useEffect(() => {
    if (!guest) return;
    const t = resetSession.get();
    if (!t) {
      router.replace("/forgot-password");
      return;
    }
    apiClient
      .forgotCheck(t)
      .then(({ username }) => {
        setToken(t);
        setAccount(username);
      })
      .catch((err) => {
        if (err instanceof ApiError && err.code === "RESET_EXPIRED") return backToVerify(err.message);
        // Lỗi mạng / hệ thống khi kiểm: vẫn cho bấm Đặt lại — máy chủ tự kiểm phiên (A4) trước khi đổi mật khẩu.
        setToken(t);
        setError(errorMessage(err));
      });
  }, [guest]);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!token) return;
    setError("");
    setLoading(true);
    try {
      // Máy chủ kiểm theo thứ tự: phiên (A4) trước mọi lỗi nhập (A1, A2, A5, A3) — FE không tự kiểm trước.
      await apiClient.forgotReset(token, password, confirm);
      flash.set("login", SUCCESS);
      router.push("/login");
    } catch (err) {
      setLoading(false);
      if (err instanceof ApiError && err.code === "RESET_EXPIRED") backToVerify(err.message);
      else setError(errorMessage(err)); // A1, A2, A3, A5, E1, E2
    }
  }

  // AUTH-F02-US02-A6: bỏ đặt lại → huỷ phiên rồi mới về đăng nhập (để quay lại bằng Back thấy phiên đã hết hạn).
  async function cancel(e: React.MouseEvent) {
    e.preventDefault();
    const t = token ?? resetSession.get(); // huỷ được cả khi bước kiểm phiên chưa xong / bị lỗi
    if (t) await apiClient.forgotCancel(t).catch(() => {});
    router.push("/login");
  }

  if (!guest) return null;
  return (
    <>
      {account && <p className="mb-6 text-sm text-slate-600">Tài khoản: {account}</p>}
      <form onSubmit={onSubmit} className="space-y-4" noValidate>
        <div className="space-y-1">
          <label htmlFor="new-password" className="block text-sm font-medium">
            Mật khẩu mới
          </label>
          <input
            id="new-password"
            type="password"
            autoComplete="new-password"
            aria-describedby="new-password-hint"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full rounded-md border border-slate-300 px-3 py-2 focus:border-blue-500 focus:outline-none"
          />
          <p id="new-password-hint" className="text-xs text-slate-500">
            Từ 8 đến 16 ký tự
          </p>
        </div>
        <div className="space-y-1">
          <label htmlFor="confirm-password" className="block text-sm font-medium">
            Nhập lại mật khẩu mới
          </label>
          <input
            id="confirm-password"
            type="password"
            autoComplete="new-password"
            value={confirm}
            onChange={(e) => setConfirm(e.target.value)}
            className="w-full rounded-md border border-slate-300 px-3 py-2 focus:border-blue-500 focus:outline-none"
          />
        </div>
        {error && (
          <p role="alert" data-testid="reset-error" className="rounded-md bg-red-50 px-3 py-2 text-sm text-red-700">
            {error}
          </p>
        )}
        <button
          type="submit"
          disabled={loading || !token}
          className="w-full rounded-md bg-blue-600 px-3 py-2 font-medium text-white hover:bg-blue-700 disabled:opacity-60"
        >
          {loading ? "Đang lưu…" : "Đặt lại mật khẩu"}
        </button>
      </form>
      <Link href="/login" onClick={cancel} className="mt-4 inline-block text-sm text-blue-600 hover:underline">
        &lt; Về đăng nhập
      </Link>
    </>
  );
}
