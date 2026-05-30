// src/features/auth/api/bootstrap.ts

import { apiClient } from "@/lib/api/client";

export async function bootstrapUser(): Promise<void> {
  await apiClient("/api/v1/auth/bootstrap", {
    method: "POST",
  });
}
