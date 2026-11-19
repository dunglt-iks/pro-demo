import { ResetPasswordForm } from "@/components/ResetPasswordForm";

// AUTH-F02-SCR02 Quên mật khẩu — đặt mật khẩu mới.
export default function ResetPasswordPage() {
  return (
    <main className="flex min-h-screen items-center justify-center p-4">
      <div className="w-full max-w-sm rounded-lg bg-white p-6 shadow">
        <h1 className="mb-2 text-xl font-semibold">Đặt mật khẩu mới</h1>
        <ResetPasswordForm />
      </div>
    </main>
  );
}
