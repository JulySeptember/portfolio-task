import { getAccessToken } from "@/features/auth/utils/token-storage";

export function getAuthHeader() {
  const token = getAccessToken();

  if (!token) {
    return {};
  }

  return {
    Authorization: `Bearer ${token}`,
  };
}
