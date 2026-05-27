// src/lib/server/api-proxy.ts

import { cookies } from "next/headers";

import { NextResponse } from "next/server";

const API_URL = process.env.NEXT_PUBLIC_API_URL!;

async function getAccessToken() {
  const cookieStore = await cookies();

  return cookieStore.get("access_token")?.value;
}

type ProxyOptions = {
  path: string;

  method: string;

  body?: unknown;
};

export async function proxyApi({ path, method, body }: ProxyOptions) {
  try {
    const accessToken = await getAccessToken();

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

    const response = await fetch(`${API_URL}${path}`, {
      method,

      headers: {
        "Content-Type": "application/json",

        Authorization: `Bearer ${accessToken}`,
      },

      body: body ? JSON.stringify(body) : undefined,

      cache: "no-store",
    });

    if (response.status === 204) {
      return new NextResponse(null, {
        status: 204,
      });
    }

    const contentType = response.headers.get("content-type");

    const isJson = contentType?.includes("application/json");

    const data = isJson ? await response.json() : null;

    return NextResponse.json(data, {
      status: response.status,
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
