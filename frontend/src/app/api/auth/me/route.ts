import { cookies } from "next/headers";

import { NextResponse } from "next/server";

export async function GET() {
  const cookieStore = await cookies();

  const accessToken = cookieStore.get("access_token")?.value;

  if (!accessToken) {
    return NextResponse.json(
      {
        message: "Unauthorized",
      },
      {
        status: 401,
      },
    );
  }

  // mock auth
  if (process.env.NEXT_PUBLIC_ENABLE_MOCK_AUTH === "true") {
    return NextResponse.json({
      id: "mock-user",

      name: "Mock User",

      email: "mock@example.com",
    });
  }

  // TODO:
  // Cognito userinfo endpoint

  return NextResponse.json({
    id: "unknown",

    name: "Authenticated User",

    email: "user@example.com",
  });
}
