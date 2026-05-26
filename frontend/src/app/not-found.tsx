import Link from "next/link";

import { Button } from "@/components/ui/button";

export default function NotFoundPage() {
  return (
    <main className="flex min-h-[70vh] flex-col items-center justify-center gap-6 p-8 text-center">
      <div className="space-y-3">
        <h1 className="text-5xl font-bold">404</h1>

        <p className="text-muted-foreground text-base">Page not found.</p>
      </div>

      <Button asChild className="h-11 px-6 text-base">
        <Link href="/tasks">Back to Tasks</Link>
      </Button>
    </main>
  );
}
