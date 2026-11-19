"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { apiClient, errorMessage } from "@/lib/api";
import { resetSession } from "@/lib/session";
import { useFlash } from "@/lib/useFlash";
import { useGuestOnly } from "@/lib/useGuestOnly";

// AUTH-F02-SCR01 — bước 1: xác minh tên đăng nhập + email (AUTH-F02-US01).
export function ForgotPasswordForm() {
  const router = useRouter();
  const guest = useGuestOnly();
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [error, setError] = useFlash("forgot-password"); // "Phiên xác minh đã hết hạn…" chuyển từ bước 2
  const [loading, setLoading] = useState(false);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      // Kiểm trống / khớp ở máy chủ (AUTH-F02-US01-A1, A2 — AUTH-BR02, AUTH-BR05).
      const { resetToken } = await apiClient.forgotVerify(username, email);
      resetSession.set(resetToken);
      router.push("/forgot-password/reset");
    } catch (err) {
      setError(errorMessage(err)); // A1, A2, E1, E2
      setLoading(false);
    }
  }

  if (!guest) return null;
  return (
    <>
      <p className="mb-6 text-sm text-slate-600">Nhập tên đăng nhập và email của tài khoản để đặt lại mật khẩu.</p>
      <form onSubmit={onSubmit} className="space-y-4" noValidate>
        <div className="space-y-1">
          <label htmlFor="username" className="block text-sm font-medium">
            Tên đăng nhập
          </label>
          <input
            id="username"
            autoComplete="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="w-full rounded-md border border-slate-300 px-3 py-2 focus:border-blue-500 focus:outline-none"
          />
        </div>
        <div className="space-y-1">
          <label htmlFor="email" className="block text-sm font-medium">
            Email
          </label>
          <input
            id="email"
            type="email"
            autoComplete="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full rounded-md border border-slate-300 px-3 py-2 focus:border-blue-500 focus:outline-none"
          />
        </div>
        {error && (
          <p role="alert" data-testid="forgot-error" className="rounded-md bg-red-50 px-3 py-2 text-sm text-red-700">
            {error}
          </p>
        )}
        <button
          type="submit"
          disabled={loading}
          className="w-full rounded-md bg-blue-600 px-3 py-2 font-medium text-white hover:bg-blue-700 disabled:opacity-60"
        >
          {loading ? "Đang kiểm tra…" : "Tiếp tục"}
        </button>
      </form>
      <Link href="/login" className="mt-4 inline-block text-sm text-blue-600 hover:underline">
        &lt; Về đăng nhập
      </Link>
    </>
  );
}
