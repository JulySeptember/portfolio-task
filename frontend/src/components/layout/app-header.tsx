"use client";

import { Button } from "@/components/ui/button";

import { useAuth } from "@/providers/auth-provider";

export function AppHeader() {
  const { logout } = useAuth();

  return (
    <header className="flex items-center justify-between border-b px-6 py-4">
      <h1 className="text-xl font-bold">Task App</h1>

      <Button variant="outline" onClick={logout}>
        Logout
      </Button>
    </header>
  );
}
