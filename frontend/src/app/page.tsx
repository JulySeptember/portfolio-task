// src/app/page.tsx

"use client";

import { Button } from "@/components/ui/button";
import { buildLoginURL } from "@/features/auth/lib/hosted-ui";
import { CheckCircle2, ShieldCheck, SquareCheckBig, Zap } from "lucide-react";

export default function HomePage() {
  const isMockAuth = process.env.NEXT_PUBLIC_ENABLE_MOCK_AUTH === "true";

  async function handleLogin() {
    try {
      if (isMockAuth) {
        localStorage.setItem("access_token", "local-dev-token");
        window.location.href = "/tasks";
        return;
      }

      const loginUrl = buildLoginURL();

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
            <SquareCheckBig className="text-primary h-8 w-8" />

            <span className="text-lg font-semibold tracking-tight">
              Task App
            </span>
          </div>
        </div>
      </header>

      {/* Hero */}

      <section className="mx-auto flex max-w-7xl flex-col items-center px-6 py-28 text-center">
        <h1 className="max-w-4xl text-5xl font-bold tracking-tight md:text-7xl">
          Task Management App
        </h1>

        <p className="text-muted-foreground mt-8 max-w-2xl text-lg leading-relaxed">
          Built with Next.js, TypeScript, Go, AWS, and Terraform.
        </p>

        <div className="mt-14 flex justify-center">
          <Button
            size="lg"
            className="h-24 rounded-3xl px-24 text-3xl font-bold shadow-xl"
            onClick={handleLogin}
          >
            Sign In
          </Button>
        </div>
      </section>

      {/* Features */}

      <section className="mx-auto grid max-w-7xl gap-6 px-6 pb-24 md:grid-cols-3">
        <div className="bg-card border-border rounded-2xl border p-6 shadow-sm transition-shadow hover:shadow-md">
          <div className="bg-primary/15 text-primary mb-4 flex h-12 w-12 items-center justify-center rounded-xl">
            <CheckCircle2 className="h-6 w-6" />
          </div>

          <h3 className="text-xl font-semibold">Task Management</h3>

          <p className="text-muted-foreground mt-3 text-sm leading-relaxed">
            Create, update, filter, and organize tasks with smooth UX and
            optimistic updates.
          </p>
        </div>

        <div className="bg-card border-border rounded-2xl border p-6 shadow-sm transition-shadow hover:shadow-md">
          <div className="bg-primary/15 text-primary mb-4 flex h-12 w-12 items-center justify-center rounded-xl">
            <Zap className="h-6 w-6" />
          </div>

          <h3 className="text-xl font-semibold">Lambda Backend</h3>

          <p className="text-muted-foreground mt-3 text-sm leading-relaxed">
            Powered by AWS Lambda, API Gateway, Amazon Cognito, Amazon RDS, and
            Terraform.
          </p>
        </div>

        <div className="bg-card border-border rounded-2xl border p-6 shadow-sm transition-shadow hover:shadow-md">
          <div className="bg-primary/15 text-primary mb-4 flex h-12 w-12 items-center justify-center rounded-xl">
            <ShieldCheck className="h-6 w-6" />
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
