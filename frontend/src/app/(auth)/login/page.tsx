// src/app/(auth)/login/page.tsx

"use client";

import { Button } from "@/components/ui/button";

import { buildLoginURL } from "@/features/auth/lib/hosted-ui";

export default function LoginPage() {
  const isMockAuth = process.env.NEXT_PUBLIC_ENABLE_MOCK_AUTH === "true";

  async function handleMockLogin() {
    const response = await fetch("/api/auth/mock-login", {
      method: "POST",
    });

    if (!response.ok) {
      console.error("mock login failed");

      return;
    }

    window.location.href = "/tasks";
  }

  return (
    <div className="flex min-h-screen items-center justify-center">
      <Button
        onClick={() => {
          if (isMockAuth) {
            handleMockLogin();

            return;
          }

          window.location.href = buildLoginURL();
        }}
      >
        Sign In
      </Button>
    </div>
  );
}
