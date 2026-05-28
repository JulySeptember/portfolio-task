import { cookies } from "next/headers";

import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  const cookieStore = await cookies();

  cookieStore.delete("access_token");

  cookieStore.delete("id_token");

  cookieStore.delete("refresh_token");

  // =========================
  // mock auth
  // =========================

  if (process.env.NEXT_PUBLIC_ENABLE_MOCK_AUTH === "true") {
    return NextResponse.redirect(new URL("/", request.url));
  }

  // =========================
  // cognito logout
  // =========================

  const logoutURL = new URL(
    `https://${process.env.NEXT_PUBLIC_COGNITO_DOMAIN}/logout`,
  );

  logoutURL.searchParams.set(
    "client_id",
    process.env.NEXT_PUBLIC_COGNITO_CLIENT_ID!,
  );

  logoutURL.searchParams.set(
    "logout_uri",
    process.env.NEXT_PUBLIC_COGNITO_LOGOUT_REDIRECT_URI!,
  );

  return NextResponse.redirect(logoutURL);
}
