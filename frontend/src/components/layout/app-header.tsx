// src/components/layout/app-header.tsx

import { Button } from "@/components/ui/button";

type Props = {
  user?: {
    name?: string;

    email?: string;
  } | null;
};

export function AppHeader({ user }: Props) {
  return (
    <header className="border-b bg-background">
      <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-8">
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-full border text-sm font-semibold">
            {user?.email?.[0]?.toUpperCase() ?? "U"}
          </div>

          <div className="min-w-0">
            <p className="truncate text-sm font-medium">
              {user?.name ?? "Mock User"}
            </p>

            <p className="text-muted-foreground truncate text-xs">
              {user?.email ?? "mock@example.com"}
            </p>
          </div>
        </div>

        <Button asChild variant="outline">
          <a href="/api/auth/logout">Logout</a>
        </Button>
      </div>
    </header>
  );
}
