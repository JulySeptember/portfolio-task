import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  reactCompiler: true,

  devIndicators: false,

  output: "export",
  trailingSlash: true,
};

export default nextConfig;
