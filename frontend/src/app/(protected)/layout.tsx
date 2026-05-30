"use client";

// src/app/(protected)/layout.tsx

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

import { AppHeader } from "@/components/layout/app-header";
import { apiClient } from "@/lib/api/client";

export default function ProtectedLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;

    const verifyAuth = async () => {
      const accessToken = localStorage.getItem("access_token");

      if (!accessToken) {
        router.replace("/");
        return;
      }

      try {
        // apiClient で自動 Refresh 対応
        await apiClient("/api/v1/users/me");

        if (!cancelled) setIsLoading(false);
      } catch (error) {
        console.error("ProtectedLayout auth failed:", error);

        // Token削除
        localStorage.removeItem("access_token");
        localStorage.removeItem("id_token");
        localStorage.removeItem("refresh_token");

        if (!cancelled) router.replace("/");
      }
    };

    void verifyAuth();

    return () => {
      cancelled = true;
    };
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
