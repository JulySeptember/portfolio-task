import { NextRequest } from "next/server";

import { proxyApi } from "@/lib/server/api-proxy";

export const dynamic = "force-dynamic";

export async function GET(request: NextRequest) {
  const search = request.nextUrl.search;

  return proxyApi({
    path: `/tasks${search}`,
    method: "GET",
  });
}

export async function POST(request: NextRequest) {
  let body: unknown = null;

  try {
    body = await request.json();
  } catch {
    body = null;
  }

  return proxyApi({
    path: "/tasks",
    method: "POST",
    body,
  });
}
