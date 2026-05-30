// src/app/auth/callback/page.tsx
"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

import { setTokens } from "@/features/auth/lib/token-storage";
import { bootstrapUser } from "@/features/auth/api/bootstrap";

function parseJwt(token: string) {
  const base64 = token.split(".")[1].replace(/-/g, "+").replace(/_/g, "/");

  const jsonPayload = decodeURIComponent(
    atob(base64)
      .split("")
      .map((c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
      .join(""),
  );

  return JSON.parse(jsonPayload);
}

export default function AuthCallbackPage() {
  const router = useRouter();

  useEffect(() => {
    let cancelled = false;

    const handleAuth = async () => {
      try {
        const hash = new URLSearchParams(window.location.hash.slice(1));

        const accessToken = hash.get("access_token");
        const idToken = hash.get("id_token");

        if (!accessToken || !idToken) {
          throw new Error("Tokens not found");
        }

        const payload = parseJwt(idToken);

        setTokens({
          accessToken,
          idToken,
          expiresAt: payload.exp,
        });

        // URLから #access_token=... を即削除
        window.history.replaceState(
          {},
          document.title,
          window.location.pathname,
        );

        await bootstrapUser();

        if (!cancelled) {
          router.replace("/tasks");
        }
      } catch (error) {
        console.error("Authentication failed", error);

        if (!cancelled) {
          router.replace("/");
        }
      }
    };

    void handleAuth();

    return () => {
      cancelled = true;
    };
  }, [router]);

  return <div>Logging in...</div>;
}
