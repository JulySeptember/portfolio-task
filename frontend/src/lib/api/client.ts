// src/lib/api/client.ts

import { ApiError } from "./error";

type ApiClientOptions = RequestInit;

async function request<T>(
  input: RequestInfo | URL,
  options: ApiClientOptions = {},
): Promise<T> {
  const { headers, ...init } = options;

  const accessToken =
    typeof window !== "undefined" ? localStorage.getItem("access_token") : null;

  const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}${input}`, {
    ...init,

    headers: {
      "Content-Type": "application/json",

      ...(accessToken
        ? {
            Authorization: `Bearer ${accessToken}`,
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
      localStorage.removeItem("access_token");

      localStorage.removeItem("id_token");

      localStorage.removeItem("refresh_token");

      window.location.href = "/";
    }

    throw error;
  }
}
