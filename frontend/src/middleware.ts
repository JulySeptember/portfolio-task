// src/middleware.ts

import type { NextRequest } from "next/server";

import { NextResponse } from "next/server";

const PROTECTED_ROUTES = ["/tasks", "/settings"];

export function middleware(request: NextRequest) {
  const pathname = request.nextUrl.pathname;

  const accessToken = request.cookies.get("access_token")?.value;

  const isProtectedRoute = PROTECTED_ROUTES.some((route) =>
    pathname.startsWith(route),
  );

  const isHomePage = pathname === "/";

  // 未ログインで protected route
  if (isProtectedRoute && !accessToken) {
    return NextResponse.redirect(new URL("/", request.url));
  }

  // ログイン済みで home
  if (isHomePage && accessToken) {
    return NextResponse.redirect(new URL("/tasks", request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/", "/tasks/:path*", "/settings/:path*"],
};
