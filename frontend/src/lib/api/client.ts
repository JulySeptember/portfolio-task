// src/lib/api/client.ts

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

type ApiClientOptions = RequestInit & {
  token?: string;
};

export async function apiClient<T>(
  input: RequestInfo | URL,
  options: ApiClientOptions = {},
): Promise<T> {
  const { token, headers, ...init } = options;

  const response = await fetch(input, {
    ...init,

    headers: {
      "Content-Type": "application/json",

      ...(token && {
        Authorization: `Bearer ${token}`,
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
