import "./globals.css";

import type { Metadata } from "next";

import { Toaster } from "@/components/ui/sonner";
import { AuthProvider } from "@/providers/auth-provider";
import { QueryProvider } from "@/providers/query-provider";

export const metadata: Metadata = {
  title: "Task App",
  description: "Serverless Task Management App",
};

type Props = {
  children: React.ReactNode;
};

export default function RootLayout({ children }: Props) {
  return (
    <html lang="ja">
      <body className="bg-background text-foreground">
        <QueryProvider>
          <AuthProvider>
            {children}
            <Toaster />
          </AuthProvider>
        </QueryProvider>
      </body>
    </html>
  );
}
