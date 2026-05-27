// src/middleware.ts

import { NextRequest, NextResponse } from "next/server";

const PROTECTED_ROUTES = ["/tasks", "/settings"];

const AUTH_ROUTES = ["/login"];

export function middleware(request: NextRequest) {
  const accessToken = request.cookies.get("access_token")?.value;

  const pathname = request.nextUrl.pathname;

  const isProtectedRoute = PROTECTED_ROUTES.some((route) =>
    pathname.startsWith(route),
  );

  const isAuthRoute = AUTH_ROUTES.some((route) => pathname === route);

  if (isProtectedRoute && !accessToken) {
    return NextResponse.redirect(new URL("/login", request.url));
  }

  if (isAuthRoute && accessToken) {
    return NextResponse.redirect(new URL("/tasks", request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/tasks/:path*", "/settings/:path*", "/login"],
};
