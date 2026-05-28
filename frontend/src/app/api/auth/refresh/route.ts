// src/app/api/auth/refresh/route.ts

import { cookies } from "next/headers";

import { NextResponse } from "next/server";

export async function POST() {
  try {
    const cookieStore = await cookies();

    const refreshToken = cookieStore.get("refresh_token")?.value;

    if (!refreshToken) {
      return NextResponse.json(
        {
          message: "No refresh token",
        },
        {
          status: 401,
        },
      );
    }

    const body = new URLSearchParams({
      grant_type: "refresh_token",

      client_id: process.env.NEXT_PUBLIC_COGNITO_CLIENT_ID!,

      refresh_token: refreshToken,
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
      return NextResponse.json(
        {
          message: "Refresh failed",
        },
        {
          status: 401,
        },
      );
    }

    const tokens = await response.json();

    cookieStore.set("access_token", tokens.access_token, {
      httpOnly: true,

      secure: process.env.NODE_ENV === "production",

      sameSite: "lax",

      path: "/",

      maxAge: 60 * 60,
    });

    if (tokens.id_token) {
      cookieStore.set("id_token", tokens.id_token, {
        httpOnly: true,

        secure: process.env.NODE_ENV === "production",

        sameSite: "lax",

        path: "/",

        maxAge: 60 * 60,
      });
    }

    return NextResponse.json({
      success: true,
    });
  } catch (error) {
    console.error(error);

    return NextResponse.json(
      {
        message: "Internal Server Error",
      },
      {
        status: 500,
      },
    );
  }
}
