// src/app/api/auth/callback/route.ts

import { cookies } from "next/headers";

import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  try {
    const code = request.nextUrl.searchParams.get("code");

    if (!code) {
      return NextResponse.redirect(new URL("/", request.url));
    }

    const body = new URLSearchParams({
      grant_type: "authorization_code",

      client_id: process.env.NEXT_PUBLIC_COGNITO_CLIENT_ID!,

      code,

      redirect_uri: process.env.NEXT_PUBLIC_COGNITO_REDIRECT_URI!,
    });

    const response = await fetch(
      `https://${process.env.NEXT_PUBLIC_COGNITO_DOMAIN}/oauth2/token`,
      {
        method: "POST",

        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },

        body,
      },
    );

    if (!response.ok) {
      return NextResponse.redirect(new URL("/", request.url));
    }

    const tokens = await response.json();

    const cookieStore = await cookies();

    cookieStore.set("access_token", tokens.access_token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === "production",
      sameSite: "lax",
      path: "/",
      maxAge: 60 * 60,
    });

    cookieStore.set("id_token", tokens.id_token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === "production",
      sameSite: "lax",
      path: "/",
      maxAge: 60 * 60,
    });

    cookieStore.set("refresh_token", tokens.refresh_token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === "production",
      sameSite: "lax",
      path: "/",
      maxAge: 60 * 60 * 24 * 30,
    });

    return NextResponse.redirect(new URL("/tasks", request.url));
  } catch (error) {
    console.error(error);

    return NextResponse.redirect(new URL("/", request.url));
  }
}
