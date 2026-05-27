import { apiClient } from "@/lib/api/client";

export type CurrentUser = {
  id?: string;

  name?: string;

  email?: string;
};

export async function getCurrentUser(): Promise<CurrentUser | null> {
  try {
    return await apiClient<CurrentUser>("/api/auth/me");
  } catch {
    return null;
  }
}
