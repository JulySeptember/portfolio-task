// src/lib/api/client.ts

import {
  clearTokens,
  getAccessToken,
} from "@/features/auth/utils/token-storage";

import { toast } from "sonner";
import { refreshToken } from "@/features/auth/api/refresh-token";

export class ApiError extends Error {
  constructor(
    public status: number,
    public body?: unknown,
    message?: string,
  ) {
    super(message ?? "API Error");

    this.name = "ApiError";
  }
}

type ApiClientOptions = RequestInit;

let refreshPromise: Promise<string> | null = null;

async function request<T>(
  input: RequestInfo | URL,
  options: ApiClientOptions = {},
  accessToken?: string,
): Promise<T> {
  const { headers, ...init } = options;

  const response = await fetch(input, {
    ...init,

    headers: {
      "Content-Type": "application/json",

      ...(accessToken && {
        Authorization: `Bearer ${accessToken}`,
      }),

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
    throw new ApiError(
      response.status,
      body,
      typeof body === "object" &&
        body !== null &&
        "message" in body &&
        typeof body.message === "string"
        ? body.message
        : "API Error",
    );
  }

  return body as T;
}

export async function apiClient<T>(
  input: RequestInfo | URL,
  options: ApiClientOptions = {},
  isRetry = false,
): Promise<T> {
  const accessToken = getAccessToken();

  try {
    return await request<T>(input, options, accessToken ?? undefined);
  } catch (error) {
    if (error instanceof ApiError && error.status === 401 && !isRetry) {
      try {
        if (!refreshPromise) {
          refreshPromise = refreshToken();
        }

        const newAccessToken = await refreshPromise;

        refreshPromise = null;

        return await apiClient<T>(input, options, true);
      } catch {
        refreshPromise = null;

        clearTokens();

        toast.error("Session expired. Please login again.");

        setTimeout(() => {
          window.location.href = "/login";
        }, 500);

        throw error;
        throw error;
      }
    }

    throw error;
  }
}
