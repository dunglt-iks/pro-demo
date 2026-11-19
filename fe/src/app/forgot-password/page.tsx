import { ForgotPasswordForm } from "@/components/ForgotPasswordForm";

// AUTH-F02-SCR01 Quên mật khẩu — xác minh tài khoản.
export default function ForgotPasswordPage() {
  return (
    <main className="flex min-h-screen items-center justify-center p-4">
      <div className="w-full max-w-sm rounded-lg bg-white p-6 shadow">
        <h1 className="mb-2 text-xl font-semibold">Quên mật khẩu</h1>
        <ForgotPasswordForm />
      </div>
    </main>
  );
}
