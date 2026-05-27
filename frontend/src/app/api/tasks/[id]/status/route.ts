import { NextRequest } from "next/server";

import { proxyApi } from "@/lib/server/api-proxy";

type Params = {
  params: Promise<{
    id: string;
  }>;
};

export async function PATCH(request: NextRequest, { params }: Params) {
  const { id } = await params;

  const body = await request.json();

  return proxyApi({
    path: `/tasks/${id}/status`,
    method: "PATCH",
    body,
  });
}
