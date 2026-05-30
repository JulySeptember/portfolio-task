// src/components/layout/app-header.tsx

"use client";

import { useQueryClient } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { useCurrentUser } from "@/features/auth/hooks/use-current-user";
import { clearTokens } from "@/features/auth/lib/token-storage";
import { buildLogoutURL } from "@/features/auth/lib/hosted-ui";

export function AppHeader() {
  const queryClient = useQueryClient();
  const { data: user } = useCurrentUser();

  function handleLogout() {
    const accessToken = localStorage.getItem("access_token");
    const isMockAuth = accessToken === "local-dev-token";

    clearTokens();
    queryClient.clear();

    if (isMockAuth) {
      window.location.href = "/";
      return;
    }

    window.location.href = buildLogoutURL();
  }

  return (
    <header className="border-b bg-background">
      <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-8">
        {/* 左側：ユーザー名 / メール */}
        <div>
          <p className="truncate text-sm font-medium">
            {user?.email ?? "User"}
          </p>
        </div>

        {/* 右側：Logoutボタン */}
        <Button variant="outline" onClick={handleLogout}>
          Logout
        </Button>
      </div>
    </header>
  );
}
