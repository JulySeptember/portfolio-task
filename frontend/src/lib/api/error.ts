// src/lib/api/error.ts

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
