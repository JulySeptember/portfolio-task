// src/app/auth/callback/page.tsx

"use client";

import { useEffect } from "react";

import { useRouter } from "next/navigation";

export default function AuthCallbackPage() {
  const router = useRouter();

  useEffect(() => {
    async function authenticate() {
      const hash = window.location.hash;

      const params = new URLSearchParams(hash.replace("#", ""));

      const accessToken = params.get("access_token");

      const idToken = params.get("id_token");

      if (!accessToken || !idToken) {
        router.replace("/");

        return;
      }

      localStorage.setItem("access_token", accessToken);

      localStorage.setItem("id_token", idToken);

      const apiURL = process.env.NEXT_PUBLIC_API_URL;

      try {
        const meResponse = await fetch(`${apiURL}/api/v1/users/me`, {
          headers: {
            Authorization: `Bearer ${idToken}`,
          },
        });

        // 初回ログイン時のみ bootstrap
        if (meResponse.status === 404) {
          await fetch(`${apiURL}/api/v1/auth/bootstrap`, {
            method: "POST",

            headers: {
              Authorization: `Bearer ${idToken}`,
            },
          });
        }

        router.replace("/tasks");
      } catch (error) {
        console.error(error);

        router.replace("/");
      }
    }

    authenticate();
  }, [router]);

  return (
    <main className="flex min-h-screen items-center justify-center">
      <p className="text-sm text-muted-foreground">Signing in...</p>
    </main>
  );
}
