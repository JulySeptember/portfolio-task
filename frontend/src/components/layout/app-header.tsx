// src/components/layout/app-header.tsx

"use client";

import { useEffect, useState } from "react";

import { Button } from "@/components/ui/button";

import { buildLogoutURL } from "@/features/auth/lib/hosted-ui";

type CurrentUser = {
  name?: string;

  email?: string;
};

export function AppHeader() {
  const [user, setUser] = useState<CurrentUser | null>(null);

  useEffect(() => {
    const accessToken = localStorage.getItem("access_token");

    // local dev mock auth
    if (accessToken === "local-dev-token") {
      setUser({
        name: "Dev User",

        email: "dev@example.com",
      });

      return;
    }

    async function fetchMe() {
      try {
        if (!accessToken) {
          return;
        }

        const response = await fetch(
          `${process.env.NEXT_PUBLIC_API_URL}/api/v1/users/me`,
          {
            headers: {
              Authorization: `Bearer ${accessToken}`,
            },
          },
        );

        if (!response.ok) {
          return;
        }

        const data = await response.json();

        setUser({
          name: data.name ?? data.email ?? "User",

          email: data.email,
        });
      } catch (error) {
        console.error(error);
      }
    }

    fetchMe();
  }, []);

  function handleLogout() {
    const accessToken = localStorage.getItem("access_token");

    const isMockAuth = accessToken === "local-dev-token";

    // clear auth
    localStorage.removeItem("access_token");

    localStorage.removeItem("id_token");

    // local dev logout
    if (isMockAuth) {
      window.location.href = "/";

      return;
    }

    // cognito logout
    window.location.href = buildLogoutURL();
  }

  return (
    <header className="border-b bg-background">
      <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-8">
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-full border text-sm font-semibold">
            {user?.email?.[0]?.toUpperCase() ?? "U"}
          </div>

          <div className="min-w-0">
            <p className="truncate text-sm font-medium">
              {user?.name ?? "User"}
            </p>

            <p className="text-muted-foreground truncate text-xs">
              {user?.email ?? "unknown@example.com"}
            </p>
          </div>
        </div>

        <Button variant="outline" onClick={handleLogout}>
          Logout
        </Button>
      </div>
    </header>
  );
}
