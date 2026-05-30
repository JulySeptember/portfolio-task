// src/features/auth/hooks/use-current-user.ts

"use client";

import { useQuery } from "@tanstack/react-query";

import { getMe } from "../api/get-me";

export function useCurrentUser() {
  return useQuery({
    queryKey: ["me"],
    queryFn: getMe,
    staleTime: 5 * 60 * 1000,
    retry: false,
  });
}
