import { NextRequest } from "next/server";

import { proxyApi } from "@/lib/server/api-proxy";

export async function GET(request: NextRequest) {
  const search = request.nextUrl.search;

  return proxyApi({
    path: `/tasks${search}`,
    method: "GET",
  });
}

export async function POST(request: NextRequest) {
  const body = await request.json();

  return proxyApi({
    path: "/tasks",
    method: "POST",
    body,
  });
}
