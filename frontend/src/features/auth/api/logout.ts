"use client";

import { toast } from "sonner";

export function logout() {
  toast.success("Logged out");

  const isMockAuth = process.env.NEXT_PUBLIC_ENABLE_MOCK_AUTH === "true";

  if (isMockAuth) {
    window.location.href = "/login";

    return;
  }

  setTimeout(() => {
    window.location.href = "/api/auth/logout";
  }, 300);
}
