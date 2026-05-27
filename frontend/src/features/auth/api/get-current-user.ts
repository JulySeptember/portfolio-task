import { apiClient } from "@/lib/api/client";

import { currentUserSchema, type CurrentUser } from "../types/auth-user";

export async function getCurrentUser(): Promise<CurrentUser> {
  const data = await apiClient<unknown>(
    `${process.env.NEXT_PUBLIC_API_URL}/me`,
  );

  return currentUserSchema.parse(data);
}
