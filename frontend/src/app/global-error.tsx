"use client";

import { useEffect } from "react";

import { Button } from "@/components/ui/button";

type Props = {
  error: Error & {
    digest?: string;
  };

  reset: () => void;
};

export default function GlobalError({ error, reset }: Props) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <html>
      <body>
        <main className="flex min-h-screen flex-col items-center justify-center gap-6 p-8 text-center">
          <div className="space-y-3">
            <h1 className="text-5xl font-bold">Fatal Error</h1>

            <p className="text-muted-foreground text-base">
              The application crashed unexpectedly.
            </p>
          </div>

          <Button
            onClick={() => {
              reset();
            }}
          >
            Reload App
          </Button>
        </main>
      </body>
    </html>
  );
}
