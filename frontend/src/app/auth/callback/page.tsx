"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

import { exchangeCodeForTokens } from "@/features/auth/api/exchange-code";
import { bootstrapUser } from "@/features/auth/api/bootstrap";

export default function AuthCallbackPage() {
  const router = useRouter();

  useEffect(() => {
    let cancelled = false;

    const handleAuth = async () => {
      try {
        const params = new URLSearchParams(window.location.search);

        const code = params.get("code");

        if (!code) {
          throw new Error("Authorization code not found");
        }

        // code → token交換
        await exchangeCodeForTokens(code);

        // users table同期
        await bootstrapUser();

        if (!cancelled) {
          router.replace("/tasks");
        }
      } catch (error) {
        console.error("Authentication callback failed", error);
      }
    };

    void handleAuth();

    return () => {
      cancelled = true;
    };
  }, [router]);

  return <div>Logging in...</div>;
}
