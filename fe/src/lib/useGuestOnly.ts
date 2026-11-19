"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { apiClient } from "./api";
import { session } from "./session";

/**
 * Trang dành cho người chưa đăng nhập (đăng nhập, quên mật khẩu): có phiên đăng nhập còn hiệu lực
 * (AUTH-BR03, AUTH-BR07) thì chuyển về trang chủ. Trả true khi đã chắc là khách — trang mới hiện nội dung.
 */
export function useGuestOnly(): boolean {
  const router = useRouter();
  const [guest, setGuest] = useState(false);
  useEffect(() => {
    if (!session.get()) {
      setGuest(true);
      return;
    }
    apiClient
      .meQuiet()
      .then(() => router.replace("/"))
      .catch(() => setGuest(true)); // token hết hạn / không hợp lệ: apiClient đã xoá token
  }, [router]);
  return guest;
}
