import { NextRequest } from "next/server";

import { proxyApi } from "@/lib/server/api-proxy";

type Params = {
  params: Promise<{
    id: string;
  }>;
};

export async function GET(_: NextRequest, { params }: Params) {
  const { id } = await params;

  return proxyApi({
    path: `/tasks/${id}`,
    method: "GET",
  });
}

export async function PUT(request: NextRequest, { params }: Params) {
  const { id } = await params;

  const body = await request.json();

  return proxyApi({
    path: `/tasks/${id}`,
    method: "PUT",
    body,
  });
}

export async function DELETE(_: NextRequest, { params }: Params) {
  const { id } = await params;

  return proxyApi({
    path: `/tasks/${id}`,
    method: "DELETE",
  });
}
