"use client";

import { useEffect, useRef, useState } from "react";
import { flash, type FlashPage } from "./session";

/** Đọc (và xoá) thông báo một lần dành cho trang này. Chỉ đọc một lần kể cả khi effect chạy lại (StrictMode). */
export function useFlash(page: FlashPage): [string, (message: string) => void] {
  const [message, setMessage] = useState("");
  const taken = useRef(false);
  useEffect(() => {
    if (taken.current) return;
    taken.current = true;
    const m = flash.take(page);
    if (m) setMessage(m);
  }, [page]);
  return [message, setMessage];
}
