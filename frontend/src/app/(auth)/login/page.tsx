"use client";

import { Button } from "@/components/ui/button";

import { buildLoginURL } from "@/features/auth/utils/hosted-ui";

import { saveTokens } from "@/features/auth/utils/token-storage";

export default function LoginPage() {
  const isMockAuth = process.env.NEXT_PUBLIC_ENABLE_MOCK_AUTH === "true";

  function handleMockLogin() {
    saveTokens({
      access_token: "mock-access",

      id_token: "mock-id",

      refresh_token: "mock-refresh",
    });

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
