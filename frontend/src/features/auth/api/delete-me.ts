import { apiClient } from "@/lib/api/client";

export function deleteMe() {
  return apiClient("/api/v1/users/me", {
    method: "DELETE",
  });
}
