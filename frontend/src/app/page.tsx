// src/app/page.tsx

"use client";

import { Button } from "@/components/ui/button";
import { buildLoginURL } from "@/features/auth/lib/hosted-ui";

export default function HomePage() {
  const isMockAuth = process.env.NEXT_PUBLIC_ENABLE_MOCK_AUTH === "true";

  async function handleLogin() {
    try {
      if (isMockAuth) {
        localStorage.setItem("access_token", "local-dev-token");
        window.location.href = "/tasks";
        return;
      }

      const loginUrl = await buildLoginURL();

      window.location.href = loginUrl;
    } catch (error) {
      console.error("Login redirect failed:", error);
    }
  }

  return (
    <main className="min-h-screen bg-background text-foreground">
      {/* Header */}

      <header className="border-border/80 border-b backdrop-blur">
        <div className="mx-auto flex h-16 max-w-7xl items-center px-6">
          <div className="flex items-center gap-3">
            <div className="bg-primary h-8 w-8 rounded-lg" />

            <span className="text-lg font-semibold tracking-tight">
              Task App
            </span>
          </div>
        </div>
      </header>

      {/* Hero */}

      <section className="mx-auto flex max-w-7xl flex-col items-center px-6 py-28 text-center">
        <div className="border-border bg-card/70 mb-8 rounded-full border px-4 py-1 text-sm backdrop-blur">
          Serverless Task Management App
        </div>

        <h1 className="max-w-4xl text-5xl font-bold tracking-tight md:text-7xl">
          Organize your tasks with a clean GitHub-style workflow.
        </h1>

        <p className="text-muted-foreground mt-8 max-w-2xl text-lg leading-relaxed">
          Built with Next.js, TypeScript, React Query, Tailwind CSS, AWS, Go,
          Terraform, and modern serverless architecture.
        </p>

        <div className="mt-14 flex justify-center">
          <Button
            size="lg"
            className="h-16 rounded-2xl px-14 text-xl font-semibold shadow-md"
            onClick={handleLogin}
          >
            Get Started
          </Button>
        </div>
      </section>

      {/* Features */}

      <section className="mx-auto grid max-w-7xl gap-6 px-6 pb-24 md:grid-cols-3">
        <div className="bg-card border-border rounded-2xl border p-6 shadow-sm transition-shadow hover:shadow-md">
          <div className="bg-primary/15 text-primary mb-4 flex h-12 w-12 items-center justify-center rounded-xl text-xl">
            ✓
          </div>

          <h3 className="text-xl font-semibold">Task Management</h3>

          <p className="text-muted-foreground mt-3 text-sm leading-relaxed">
            Create, update, filter, and organize tasks with smooth UX and
            optimistic updates.
          </p>
        </div>

        <div className="bg-card border-border rounded-2xl border p-6 shadow-sm transition-shadow hover:shadow-md">
          <div className="bg-primary/15 text-primary mb-4 flex h-12 w-12 items-center justify-center rounded-xl text-xl">
            ⚡
          </div>

          <h3 className="text-xl font-semibold">Serverless Backend</h3>

          <p className="text-muted-foreground mt-3 text-sm leading-relaxed">
            Powered by AWS Lambda, API Gateway, Cognito, Terraform, and RDS.
          </p>
        </div>

        <div className="bg-card border-border rounded-2xl border p-6 shadow-sm transition-shadow hover:shadow-md">
          <div className="bg-primary/15 text-primary mb-4 flex h-12 w-12 items-center justify-center rounded-xl text-xl">
            🔒
          </div>

          <h3 className="text-xl font-semibold">Secure Authentication</h3>

          <p className="text-muted-foreground mt-3 text-sm leading-relaxed">
            Authentication handled through Amazon Cognito Hosted UI and JWT
            authorization.
          </p>
        </div>
      </section>
    </main>
  );
}
