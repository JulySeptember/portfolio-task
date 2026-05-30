import {
  clearTokens,
  getTokens,
  isAccessTokenExpired,
} from "@/features/auth/lib/token-storage";

import { ApiError } from "./error";

type ApiClientOptions = RequestInit;

export async function apiClient<T>(
  input: RequestInfo | URL,
  options: ApiClientOptions = {},
): Promise<T> {
  const { headers, ...init } = options;

  const tokens = getTokens();

  // トークン期限切れなら即ログアウト
  if (tokens && isAccessTokenExpired()) {
    clearTokens();

    if (typeof window !== "undefined") {
      window.location.href = "/";
    }

    throw new ApiError(401, null, "Token expired");
  }

  const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}${input}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(tokens?.accessToken
        ? {
            Authorization: `Bearer ${tokens.accessToken}`,
          }
        : {}),
      ...headers,
    },
  });

  const contentType = response.headers.get("content-type");

  let body: unknown = null;

  try {
    if (contentType?.includes("application/json")) {
      body = await response.json();
    } else {
      body = await response.text();
    }
  } catch {
    body = null;
  }

  if (!response.ok) {
    if (typeof window !== "undefined" && response.status === 401) {
      clearTokens();
      window.location.href = "/";
    }

    throw new ApiError(response.status, body, "API Error");
  }

  return body as T;
}
