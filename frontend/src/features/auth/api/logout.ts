"use client";

import { toast } from "sonner";

import { buildLogoutURL } from "../utils/hosted-ui";

import { clearTokens } from "../utils/token-storage";

export function logout() {
  clearTokens();

  toast.success("Logged out");

  const isMockAuth = process.env.NEXT_PUBLIC_ENABLE_MOCK_AUTH === "true";

  if (isMockAuth) {
    window.location.href = "/login";

    return;
  }

  setTimeout(() => {
    window.location.href = buildLogoutURL();
  }, 300);
}
