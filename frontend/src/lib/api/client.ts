// src/lib/api/client.ts

import { ApiError } from "./error";

type ApiClientOptions = RequestInit;

async function request<T>(
  input: RequestInfo | URL,
  options: ApiClientOptions = {},
): Promise<T> {
  const { headers, ...init } = options;

  const response = await fetch(input, {
    ...init,

    credentials: "include",

    headers: {
      "Content-Type": "application/json",

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
): Promise<T> {
  try {
    return await request<T>(input, options);
  } catch (error) {
    if (
      typeof window !== "undefined" &&
      error instanceof ApiError &&
      error.status === 401
    ) {
      try {
        await fetch("/api/auth/refresh", {
          method: "POST",

          credentials: "include",
        });

        return await request<T>(input, options);
      } catch {
        window.location.href = "/";
      }
    }

    throw error;
  }
}
