"use client";

import { useEffect } from "react";

import { useRouter } from "next/navigation";

import { toast } from "sonner";

import { saveTokens } from "@/features/auth/utils/token-storage";

export default function AuthCallbackPage() {
  const router = useRouter();

  useEffect(() => {
    async function handleCallback() {
      try {
        const params = new URLSearchParams(window.location.search);

        const code = params.get("code");

        if (!code) {
          throw new Error("Authorization code not found");
        }

        const body = new URLSearchParams({
          grant_type: "authorization_code",

          client_id: process.env.NEXT_PUBLIC_COGNITO_CLIENT_ID!,

          code,

          redirect_uri: process.env.NEXT_PUBLIC_COGNITO_REDIRECT_URI!,
        });

        const response = await fetch(
          `https://${process.env.NEXT_PUBLIC_COGNITO_DOMAIN}/oauth2/token`,
          {
            method: "POST",

            headers: {
              "Content-Type": "application/x-www-form-urlencoded",
            },

            body,
          },
        );

        if (!response.ok) {
          throw new Error("Token exchange failed");
        }

        const tokens = await response.json();

        saveTokens(tokens);

        toast.success("Login successful");

        router.replace("/tasks");
      } catch (error) {
        console.error(error);

        toast.error("Login failed");

        router.replace("/login");
      }
    }

    void handleCallback();
  }, [router]);

  return (
    <div className="flex min-h-screen items-center justify-center">
      <p className="text-muted-foreground">Signing in...</p>
    </div>
  );
}
