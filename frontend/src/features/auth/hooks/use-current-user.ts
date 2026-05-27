"use client";

import { useQuery } from "@tanstack/react-query";

import { getCurrentUser } from "../api/get-current-user";

import { getAccessToken } from "../utils/token-storage";

export function useCurrentUser() {
  const isMockAuth = process.env.NEXT_PUBLIC_ENABLE_MOCK_AUTH === "true";

  return useQuery({
    queryKey: ["current-user"],

    queryFn: async () => {
      if (isMockAuth) {
        return {
          id: 1,

          email: "mock@example.com",

          name: "Mock User",
        };
      }

      return getCurrentUser();
    },

    enabled: !!getAccessToken(),
  });
}
