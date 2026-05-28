// src/app/api/auth/me/route.ts

import { cookies } from "next/headers";

import { NextResponse } from "next/server";

export async function GET() {
  try {
    const cookieStore = await cookies();

    const accessToken = cookieStore.get("access_token")?.value;

    if (!accessToken) {
      return NextResponse.json({ message: "Unauthorized" }, { status: 401 });
    }

    const response = await fetch(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/v1/users/me`,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },

        cache: "no-store",
      },
    );

    if (!response.ok) {
      return NextResponse.json(
        { message: "Failed to fetch current user" },
        { status: response.status },
      );
    }

    const user = await response.json();

    return NextResponse.json(user);
  } catch (error) {
    console.error(error);

    return NextResponse.json(
      { message: "Internal Server Error" },
      { status: 500 },
    );
  }
}
