"use client";

import { useEffect } from "react";

import { Button } from "@/components/ui/button";

type Props = {
  error: Error & {
    digest?: string;
  };

  reset: () => void;
};

export default function ErrorPage({ error, reset }: Props) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <main className="flex min-h-[70vh] flex-col items-center justify-center gap-6 p-8 text-center">
      <div className="space-y-3">
        <h1 className="text-4xl font-bold">Something went wrong</h1>

        <p className="text-muted-foreground text-base">
          An unexpected error occurred.
        </p>
      </div>

      <Button onClick={() => reset()} className="h-11 px-6 text-base">
        Try Again
      </Button>
    </main>
  );
}
