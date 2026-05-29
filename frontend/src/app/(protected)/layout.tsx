"use client";

// src/app/(protected)/layout.tsx

import { useEffect, useState } from "react";

import { useRouter } from "next/navigation";

import { AppHeader } from "@/components/layout/app-header";

export default function ProtectedLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();

  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    async function verifyAuth() {
      const accessToken = localStorage.getItem("access_token");

      if (!accessToken) {
        router.replace("/");

        return;
      }

      try {
        const response = await fetch(
          `${process.env.NEXT_PUBLIC_API_URL}/api/v1/users/me`,
          {
            headers: {
              Authorization: `Bearer ${accessToken}`,
            },
          },
        );

        if (!response.ok) {
          localStorage.removeItem("access_token");

          localStorage.removeItem("id_token");

          router.replace("/");

          return;
        }

        setIsLoading(false);
      } catch (error) {
        console.error(error);

        router.replace("/");
      }
    }

    verifyAuth();
  }, [router]);

  if (isLoading) {
    return (
      <main className="flex min-h-screen items-center justify-center">
        <p className="text-sm text-muted-foreground">Loading...</p>
      </main>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      <AppHeader />

      <main>{children}</main>
    </div>
  );
}
