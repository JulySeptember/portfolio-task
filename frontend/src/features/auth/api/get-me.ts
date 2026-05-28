import { apiClient } from "@/lib/api/client";

export type CurrentUser = {
  id: number;

  auth_user_id: string;

  email: string;

  created_at: string;

  updated_at: string;
};

export function getMe() {
  return apiClient<CurrentUser>("/api/v1/users/me");
}
