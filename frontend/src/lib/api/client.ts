import {
  clearTokens,
  getTokens,
  isAccessTokenExpired,
} from "@/features/auth/lib/token-storage";

import { refreshAccessToken } from "@/features/auth/api/refresh-token";

import { ApiError } from "./error";
type ApiClientOptions = RequestInit;

async function request<T>(
  input: RequestInfo | URL,
  options: ApiClientOptions = {},
) {
  const { headers, ...init } = options;

  let tokens = getTokens();

  if (tokens && isAccessTokenExpired()) {
    const refreshed = await refreshAccessToken();
    if (refreshed) tokens = refreshed;
  }

  const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}${input}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(tokens?.accessToken
        ? { Authorization: `Bearer ${tokens.accessToken}` }
        : {}),
      ...headers,
    },
  });

  const contentType = response.headers.get("content-type");

  let body: unknown = null;
  try {
    if (contentType?.includes("application/json")) body = await response.json();
    else body = await response.text();
  } catch {
    body = null;
  }

  if (!response.ok) throw new ApiError(response.status, body, "API Error");

  return body as T;
}

export async function apiClient<T>(
  input: RequestInfo | URL,
  options: ApiClientOptions = {},
): Promise<T> {
  try {
    return await request<T>(input, options);
  } catch (error) {
    if (
      typeof window !== "undefined" &&
      error instanceof ApiError &&
      error.status === 401
    ) {
      const refreshed = await refreshAccessToken();

      if (refreshed) {
        try {
          return await request<T>(input, options);
        } catch {
          clearTokens();
        }
      }

      clearTokens();
    }

    throw error;
  }
}
