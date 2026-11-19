"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { apiClient, type User } from "@/lib/api";
import { session } from "@/lib/session";

// Trang chủ: cần đăng nhập. Chưa có token → /login; token hết hạn → apiClient tự chuyển /login.
export default function HomePage() {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    if (!session.get()) {
      router.replace("/login");
      return;
    }
    apiClient.me().then(setUser).catch(() => {});
  }, [router]);

  function logout() {
    session.clear();
    router.replace("/login");
  }

  if (!user) return null;
  return (
    <main className="mx-auto max-w-2xl p-6">
      <header className="flex items-center justify-between">
        <h1 className="text-xl font-semibold">Xin chào {user.name}</h1>
        <button onClick={logout} className="rounded-md border border-slate-300 px-3 py-1.5 text-sm hover:bg-slate-100">
          Đăng xuất
        </button>
      </header>
    </main>
  );
}
