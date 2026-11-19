import type { NextConfig } from "next";

// /api/* chuyển sang BE (cùng origin với FE, không cần CORS). BE_URL mặc định http://localhost:8080.
const nextConfig: NextConfig = {
  async rewrites() {
    return [{ source: "/api/:path*", destination: `${process.env.BE_URL ?? "http://localhost:8080"}/api/:path*` }];
  },
};

export default nextConfig;
